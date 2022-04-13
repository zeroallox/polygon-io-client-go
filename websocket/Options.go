package polygonws

import (
    "errors"
    "time"
)

type Options struct {
    AutoReconnect         bool
    AutoReconnectInterval time.Duration
    APIKey                string
    ClusterType           ClusterType
}

func validateOptions(options *Options) error {

    if options.ClusterType == CT_Invalid {
        return errors.New("ClusterType Invalid")
    }

    if len(options.APIKey) == 0 {
        return errors.New("invalid APIKey")
    }

    if options.AutoReconnectInterval == 0 {
        options.AutoReconnectInterval = defaultConnectionInterval
    }

    return nil
}
