package polygonws

type Topic uint8

const (
    tp_EquityMin        Topic = 10
    TP_EquityTrades     Topic = 11
    TP_EquityQuotes     Topic = 12
    TP_EquityAggMin     Topic = 13
    TP_EquityAggSec     Topic = 14
    TP_EquityLULD       Topic = 15
    TP_EquityImbalances Topic = 16
    tp_EquityMax        Topic = 17

    tp_ForexMin    Topic = 50
    TP_ForexTrades Topic = 51
    TP_ForexAggMin Topic = 52
    tp_ForexMax    Topic = 53

    tp_CryptoMin    Topic = 60
    TP_CryptoTrades Topic = 61
    TP_CryptoQuotes Topic = 62
    TP_CryptoBook   Topic = 63
    TP_CryptoAggMin Topic = 64
    tp_CryptoMax    Topic = 65
)

func (tp Topic) subscriptionPrefix() string {
    switch tp {
    case TP_EquityTrades:
        return "T"
    case TP_EquityQuotes:
        return "Q"
    case TP_EquityAggMin:
        return "AM"
    case TP_EquityAggSec:
        return "A"
    case TP_EquityLULD:
        return "LULD"
    case TP_EquityImbalances:
        return "NOI"
    case TP_ForexTrades:
        return "C"
    case TP_ForexAggMin:
        return "CA"
    case TP_CryptoTrades:
        return "XT"
    case TP_CryptoQuotes:
        return "XQ"
    case TP_CryptoBook:
        return "XL2"
    case TP_CryptoAggMin:
        return "XA"
    default:
        return ""
    }
}
