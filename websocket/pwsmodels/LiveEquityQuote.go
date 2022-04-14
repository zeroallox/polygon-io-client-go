package pwsmodels

type LiveEquityQuote struct {
    baseModel

    Ticker      string  `json:"sym"`
    BidExchange int     `json:"bx"`
    BidPrice    float64 `json:"bp"`
    BidSize     int     `json:"bs"`
    AskExchange int     `json:"ax"`
    AskPrice    float64 `json:"ap"`
    AskSize     int     `json:"as"`
    Condition   int     `json:"c"`
    Timestamp   int64   `json:"t"`
    Tape        int     `json:"z"`
}
