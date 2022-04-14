package pwsmodels

type LiveEquityTrade struct {
    baseModel

    Ticker     string  `json:"sym"`
    Exchange   int     `json:"x"`
    TradeID    string  `json:"i"`
    Tape       int     `json:"z"`
    Price      float64 `json:"p"`
    Volume     int     `json:"s"`
    Conditions []int   `json:"c"`
    Timestamp  int64   `json:"t"`
}
