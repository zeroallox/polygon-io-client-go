package polygonws

type ClusterType uint8

const (
    CT_Invalid       ClusterType = 0
    CT_Stocks        ClusterType = 1
    CT_StocksDelayed ClusterType = 2
    CT_Forex         ClusterType = 3
    CT_Crypto        ClusterType = 4
)

func (ct ClusterType) endpoint() string {
    switch ct {
    case CT_Stocks:
        return "wss://socket.polygon.io/stocks"
    case CT_StocksDelayed:
        return "wss://delayed.polygon.io/stocks"
    case CT_Forex:
        return "wss://socket.polygon.io/forex"
    case CT_Crypto:
        return "wss://socket.polygon.io/crypto"
    default:
        panic("should never happen")
    }
}

func (ct ClusterType) supportsTopic(topic Topic) bool {

    switch ct {
    case CT_Stocks, CT_StocksDelayed:
        if topic > tp_EquityMin && topic < tp_EquityMax {
            return true
        }
    case CT_Forex:
        if topic > tp_ForexMin && topic < tp_ForexMax {
            return true
        }
    case CT_Crypto:
        if topic > tp_CryptoMin && topic < tp_CryptoMax {
            return true
        }
    }

    return false
}
