package pwsmodels

// ControlMessage a control message sent by polysocket.
type ControlMessage struct {
    Ev      string `json:"ev,omitempty"`
    Status  string `json:"status,omitempty"`
    Message string `json:"Message,omitempty"`
    Action  string `json:"action,omitempty"`
    Params  string `json:"params,omitempty"`
}
