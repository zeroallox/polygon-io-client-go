package polygonws

import (
    "fmt"
    "github.com/polygon-io/client-go/websocket/pwsmodels"
)

func makeAuthMessage(apiKey string) pwsmodels.Message {
    return pwsmodels.Message{
        Action: "auth",
        Params: apiKey,
    }
}

func makeSubscribeMessage(topic Topic, tickers ...string) pwsmodels.Message {
    return pwsmodels.Message{
        Action: "subscribe",
        Params: generateSubListString(topic, tickers),
    }
}

func makeUnsubscribeMessage(topic Topic, tickers ...string) pwsmodels.Message {
    return pwsmodels.Message{
        Action: "subscribe",
        Params: generateSubListString(topic, tickers),
    }
}

func generateSubListString(topic Topic, symbols []string) string {

    var tp = topic.subscriptionPrefix()
    var str string
    str = str + fmt.Sprintf("%v.%v", tp, symbols[0])

    for _, cSymbol := range symbols[1:] {
        str = str + fmt.Sprintf(",%v.%v", tp, cSymbol)
    }

    return str
}
