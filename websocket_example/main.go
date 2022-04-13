package main

import (
    polygonws "github.com/polygon-io/client-go/websocket"
    log "github.com/sirupsen/logrus"
    "os"
    "sync"
)

func main() {

    var opt polygonws.Options
    opt.APIKey = os.Getenv("POLY_API_KEY")
    opt.ClusterType = polygonws.CT_Stocks

    pws, err := polygonws.NewClient(opt)
    if err != nil {
        panic(err)
    }

    pws.SetOnStateChangedHandler(func(state polygonws.State) {
        log.Println("State Changed", state)
    })

    err = pws.Connect()
    if err != nil {
        panic(err)
    }

    var wg sync.WaitGroup
    wg.Add(1)
    wg.Wait()

}
