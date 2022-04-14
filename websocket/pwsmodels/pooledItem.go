package pwsmodels

import "sync"

// pooledItem interface for
type pooledItem interface {
    register(self any, owner *sync.Pool)
    setReferenceCount(rc uint64)
}
