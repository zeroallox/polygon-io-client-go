package polygonws

type State uint8

const (
    // ST_Invalid Zero value
    ST_Invalid State = 0
    // ST_Error Client encountered kind of unrecoverable error.
    ST_Error State = 1
    // ST_Closed Client was closed by calling Close()
    ST_Closed State = 2
    // ST_Disconnected Client::Disconnect() was called.
    ST_Disconnected State = 3
    // ST_Unconnected Client object was just initialized and no operations have
    // been performed.
    ST_Unconnected State = 4
    // ST_Connecting Client between when Connect() was called and
    // the server sending a confirmation "connected" message.
    ST_Connecting State = 5
    // ST_Connected Client is connected but not yet authenticated
    ST_Connected State = 6
    // ST_Ready Client is connected and ready to subscribe.
    ST_Ready State = 7
    // ST_Reconnecting Client is in the process of reconnecting.
    ST_Reconnecting State = 8
)

func (cs State) String() string {
    switch cs {
    case ST_Error:
        return "ST_Error"
    case ST_Unconnected:
        return "ST_Unconnected"
    case ST_Closed:
        return "ST_Closed"
    case ST_Connecting:
        return "ST_Connecting"
    case ST_Connected:
        return "ST_Connected"
    case ST_Ready:
        return "ST_Ready"
    case ST_Reconnecting:
        return "ST_Reconnecting"
    default:
        return "ST_Invalid"
    }
}
