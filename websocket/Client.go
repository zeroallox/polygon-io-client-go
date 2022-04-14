package polygonws

import (
    "bytes"
    "context"
    "errors"
    jsoniter "github.com/json-iterator/go"
    log "github.com/sirupsen/logrus"
    "github.com/valyala/bytebufferpool"
    "github.com/zeroallox/polygon-io-client-go/websocket/pwsmodels"
    "io"
    "nhooyr.io/websocket"
    "sync"
    "time"
)

const readBufferSize = 1000000 // 10 Megabytes
const defaultConnectionInterval = time.Second * 3

type Client struct {
    ws             *websocket.Conn
    opt            Options
    autoReconnect  bool
    mtxc           sync.Mutex
    msgQueue       []*bytebufferpool.ByteBuffer
    cond           sync.Cond
    stop           bool
    state          State
    onStateChanged OnConnectionStateChangedFunc
}

// NewClient initializes a new Client configured with Options.
func NewClient(options Options) (*Client, error) {

    var n = new(Client)
    n.cond = sync.Cond{L: &n.mtxc}
    n.state = ST_Unconnected

    var err = validateOptions(&options)
    if err != nil {
        return nil, err
    }

    n.opt = options

    return n, nil
}

type OnConnectionStateChangedFunc func(state State)

// State returns the Clients current State
func (cli *Client) State() State {
    cli.cond.L.Lock()
    defer cli.cond.L.Unlock()

    return cli.state
}

// setState sets the Clients current State and calls
// the onStateChangedHandler if one was set by the user.
func (cli *Client) setState(cs State) {
    cli.cond.L.Lock()

    if cli.state == cs {
        cli.cond.L.Unlock()
        return
    }
    cli.state = cs

    cli.cond.L.Unlock()

    if cli.onStateChanged != nil {
        cli.onStateChanged(cs)
    }
}

func (cli *Client) SetOnStateChangedHandler(onStateChanged OnConnectionStateChangedFunc) {
    cli.onStateChanged = onStateChanged
}

// Connect opens a connection to the websocket server.
func (cli *Client) Connect() error {

    if cli.State() == ST_Closed {
        return ErrClientClosed
    }

    var err = cli.closeWSConnection()
    if err != nil {
        return err
    }

    cli.setState(ST_Connecting)
    return cli.connect()
}

// connect opens the raw websocket connection and starts the
// message read thread.
func (cli *Client) connect() error {

    ctx, cancel := context.WithTimeout(context.Background(), cli.opt.AutoReconnectInterval)
    defer cancel()

    var ep = cli.opt.ClusterType.endpoint()

    conn, res, err := websocket.Dial(ctx, ep, nil)
    if err != nil {
        return err
    }
    conn.SetReadLimit(readBufferSize)

    if res.StatusCode != 101 { // Switching Protocols
        return errors.New("not switching protocols")
    }

    cli.ws = conn
    cli.msgQueue = cli.msgQueue[:0]

    go cli.beginProcessMessages()
    go cli.beginReading()
    go cli.beginPinging()

    return nil
}

// Disconnect gracefully closes the websocket connection but does
// NOT kill the message dispatch thread. This is useful if you'd
// like to reuse the client later and simply want to stop receiving
// messages.
//
// Note: AutoReconnect will be disabled.
func (cli *Client) Disconnect() error {
    cli.autoReconnect = false
    cli.setState(ST_Disconnected)

    return cli.closeWSConnection()
}

// Close gracefully closes the underlying websocket connection AND kills
// the message processor thread. Call this when you are completely done
// with the Client. Closed clients may not be reused.
func (cli *Client) Close() error {
    var err = cli.closeWSConnection()
    if err != nil {
        return err
    }

    cli.shutdownProcessor()
    cli.setState(ST_Closed)

    return nil
}

// Subscribe subscribes to the specified topic and symbols.
func (cli *Client) Subscribe(topic Topic, symbols ...string) error {
    return nil
}

// Unsubscribe unsubscribes from the specified topic and symbols.
func (cli *Client) Unsubscribe(topic Topic, symbols ...string) error {
    return nil
}

// Gracefully closes the websocket connection.
func (cli *Client) closeWSConnection() error {
    cli.cond.L.Lock()
    defer cli.cond.L.Unlock()

    if cli.ws == nil {
        return nil
    }

    return cli.ws.Close(websocket.StatusNormalClosure, "")
}

// shutdownProcessor wakes up and kills the processing thread.
func (cli *Client) shutdownProcessor() {
    cli.cond.L.Lock()
    defer cli.cond.L.Unlock()

    cli.stop = true
    cli.cond.Signal()
}

func (cli *Client) isStopped() bool {
    cli.cond.L.Lock()
    defer cli.cond.L.Unlock()

    return cli.stop
}

// beginReading reads from the websocket and submits messages
// to the processing queue.
func (cli *Client) beginReading() {

    defer func() {
        if cli.autoReconnect == true {
            time.Sleep(cli.opt.AutoReconnectInterval)
            go cli.reconnect()
        }
    }()

    var buff [readBufferSize]byte

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    for {

        msgType, reader, err := cli.ws.Reader(ctx)
        if err != nil {
            log.WithError(err).Error("polyws: get ws reader failed")
            return
        }

        if msgType != websocket.MessageText {
            err = errors.New("polyws: expected text message")
            return
        }

        var cursor = 0
        for {
            bytesRead, err := reader.Read(buff[cursor : readBufferSize-cursor])
            if err != nil {
                if err == io.EOF {
                    break
                }
                log.WithError(err).Error("polyws: reader error")
                return
            }

            cursor = cursor + bytesRead
        }

        var data = buff[0:cursor]
        var bb = bytebufferpool.Get()
        bb.Set(data)

        cli.cond.L.Lock()
        cli.msgQueue = append(cli.msgQueue, bb)
        cli.cond.L.Unlock()

        cli.cond.Signal()
    }

}

// beginPinging maintains the ping / pong with the server.
//Ensures we track catch disconnects.
func (cli *Client) beginPinging() {

    for {

        if cli.isStopped() == true {
            return
        }

        var err = cli.ws.Ping(context.Background())
        if err != nil {
            log.WithError(err).Debug("pinging failed")
            return
        }

        time.Sleep(3 * time.Second)
    }
}

// reconnect starts the reconnection cycle.
func (cli *Client) reconnect() {
    cli.setState(ST_Reconnecting)

    var err = cli.closeWSConnection()
    if err != nil {
        cli.setState(ST_Error)
    }

    err = cli.connect()
    if err != nil {
        cli.setState(ST_Error)
    }
}

// beginProcessMessages handles the decoding and dispatching of
// received messages.
func (cli *Client) beginProcessMessages() {

    log.Debug("polyws: Start Processing Thread")
    defer func() {
        log.Debug("polyws: Exit Processing Thread")
    }()

    var localQueue []*bytebufferpool.ByteBuffer

    for {

        cli.cond.L.Lock()
        for len(cli.msgQueue) == 0 {

            if cli.stop == true {
                cli.cond.L.Unlock()
                return
            }

            cli.cond.Wait()
        }

        localQueue = append(localQueue, cli.msgQueue...)
        cli.msgQueue = cli.msgQueue[:0]
        cli.cond.L.Unlock()

        cli.processMessageQueue(localQueue)

        localQueue = localQueue[:0]
    }
}

var eqQuotePrefix = []byte("{\"ev\":\"Q\"")
var eqTradePrefix = []byte("{\"ev\":\"T\"")
var statusPrefix = []byte("{\"ev\":\"status\"")

// processMessageQueue routes each message to the correct handler. We specifically
// do not break and return on error as it is possible for a single message
// to be malformed.
func (cli *Client) processMessageQueue(msgs []*bytebufferpool.ByteBuffer) {

    var jMessages []jsoniter.RawMessage

    for _, cRawMessage := range msgs {

        var err = json.Unmarshal(cRawMessage.B, &jMessages)
        if err != nil {
            log.WithError(err).Error("polyws: failed to unmarshal root message")
            bytebufferpool.Put(cRawMessage)
            continue
        }

        for _, cMessage := range jMessages {

            switch {
            case bytes.HasPrefix(cMessage, eqTradePrefix):
                err = cli.handleLiveEquityTrade(cMessage)
                break
            case bytes.HasPrefix(cMessage, eqQuotePrefix):
                err = cli.handleLiveEquityQuote(cMessage)
                break
            case bytes.HasPrefix(cMessage, statusPrefix):
                err = cli.handleStatusMessage(cMessage)
                break
            }

            if err != nil {
                log.WithError(err).Error(string(cMessage))
            }
        }

        bytebufferpool.Put(cRawMessage)
    }
}

// statusConnected is sent when the ws connection is opened and protocol has settled
const statusConnected = "connected"

// statusAuthSuccess is sent when api key validation passes
const statusAuthSuccess = "auth_success"

// statusAuthFailed is sent when api key validation fails
const statusAuthFailed = "auth_failed"

// statusMaxConnections is sent when you have exceeded your per-cluster connection entitlements
const statusMaxConnections = "max_connections"

// statusSuccess is sent un successful subscription
const statusSuccess = "success"

// handleStatusMessage handles status / control messages sent by the server.
func (cli *Client) handleStatusMessage(msg jsoniter.RawMessage) error {

    var sm pwsmodels.ControlMessage

    var err = json.Unmarshal(msg, &sm)
    if err != nil {
        return err
    }

    switch sm.Status {
    case statusConnected:
        cli.setState(ST_Connected)
        return cli.sendAuthRequest()
    case statusAuthSuccess:
        cli.setState(ST_Ready)
        return nil
    case statusAuthFailed:
        cli.setState(ST_Error)
        _ = cli.closeWSConnection()
        return ErrAuthenticationFailed
    case statusSuccess:
        // TODO: Don't need to do anything?
        return nil
    case statusMaxConnections:
        cli.reconnect()
        return nil
    default:
        cli.setState(ST_Error)
        _ = cli.closeWSConnection()
        return ErrUnhandledStatusMessage
    }
}

// handleLiveEquityQuote decodes and dispatches Quote messages
func (cli *Client) handleLiveEquityQuote(msg jsoniter.RawMessage) error {

    //if cli.onDataReceived == nil {
    //    return nil
    //}

    //var quote = polymodels.DefaultLiveEquityQuotePool.Acquire()
    //var err = json.Unmarshal(msg, &quote)
    //if err != nil {
    //    return err
    //}
    //
    //cli.onDataReceived(quote)

    return nil
}

// handleLiveEquityTrade decodes and dispatches Trade messages
func (cli *Client) handleLiveEquityTrade(msg jsoniter.RawMessage) error {

    //if cli.onDataReceived == nil {
    //    return nil
    //}

    //var trade = polymodels.DefaultLiveEquityTradePool.Acquire()
    //var err = json.Unmarshal(msg, &trade)
    //if err != nil {
    //    return err
    //}
    //
    //cli.onDataReceived(trade)

    return nil
}

// sendAuthRequest creates an auth message and sends to the websocket server
func (cli *Client) sendAuthRequest() error {

    var msg = makeAuthMessage(cli.opt.APIKey)

    jData, err := json.Marshal(&msg)
    if err != nil {
        return err
    }

    return cli.writeMessage(jData)
}

// writeMessage sends the message data to the websocket server. We do not
// need to lock as the ws lib does this for us.
func (cli *Client) writeMessage(data []byte) error {

    var state = cli.State()
    if state != ST_Ready && state != ST_Connected {
        return ErrClientNotReady
    }

    return cli.ws.Write(context.Background(), websocket.MessageText, data)
}
