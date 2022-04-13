package polygonws

import "errors"

// ErrClientNotReady is returned when user tried to perform an action
// and the client is not in the "CSReady" state.
// An example is trying to subscribe while disconnected.
var ErrClientNotReady = errors.New("client not ready")

// ErrUnsupportedTopic is returned if the ClusterType the Client is connected
// to does not support the topic. For example: You can not request Crypto Quotes
// while connected to the Stocks cluster.
var ErrUnsupportedTopic = errors.New("topic does not match cluster type")

// ErrNoSymbols is returned if a no symbols are specified during
// subscription operations.
var ErrNoSymbols = errors.New("no symbols")

// ErrClientClosed Client::Connect() called after Client::Close()
var ErrClientClosed = errors.New("client closed")

// ErrAuthenticationFailed is returned when on authentication failure.
var ErrAuthenticationFailed = errors.New("authentication failed")

// ErrUnhandledStatusMessage is returned when the client receives an unknown
// status message.
var ErrUnhandledStatusMessage = errors.New("unhandled status message")
