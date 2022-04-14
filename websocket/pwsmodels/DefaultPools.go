package pwsmodels

var DefaultLiveEquityTradePool = newModelPool[*LiveEquityTrade](func() *LiveEquityTrade {
    return new(LiveEquityTrade)
})

var DefaultLiveEquityQuotePool = newModelPool[*LiveEquityQuote](func() *LiveEquityQuote {
    return new(LiveEquityQuote)
})
