package pwsmodels

import (
    "sync"
    "sync/atomic"
)

const tagSize = 2

// baseModel "base class" for all items which are to be used
// with a ModelPool
type baseModel struct {
    self     any            // needed to enable calling release on the item inself
    owner    *sync.Pool     // the pool we've been allloc'd from
    refCount int64          // total number of outstanding references to us
    tags     [tagSize]int64 // user-defined metadata
}

// SetTag sets a user-defined value
func (m *baseModel) SetTag(idx uint8, value int64) {
    if idx > tagSize-1 {
        panic("bad tag index")
    }
    atomic.StoreInt64(&m.tags[idx], value)
}

// GetTag gets a user-defined value
func (m *baseModel) GetTag(idx uint8) int64 {
    if idx > tagSize-1 {
        panic("bad tag index")
    }
    return atomic.LoadInt64(&m.tags[idx])
}

// AutoRelease returns the item to the pool when its reference count reaches zero.
// Should only be called on instances retrieved with ModelPool.AutoAcquire
func (m *baseModel) AutoRelease() {
    switch atomic.AddInt64(&m.refCount, -1) {
    case 0:
        m.owner.Put(m.self)
        return
    case -1:
        panic("baseModel: AutoRelease called when model not AutoAcquired")
    }
}

// Release releases the model back to the Pool.
func (m *baseModel) Release() {
    m.refCount = 0
    m.owner.Put(m.self)
}

// register tracks which Pool created us and sets the reference to self
// set in ModelPool.<underlying_sync.Pool>.newFunc
func (m *baseModel) register(self any, owner *sync.Pool) {
    m.owner = owner
    m.self = self
}

// setReferenceCount sets our reference count to the specified value.
// A user who sets a reference count should only use a positive number,
// yet we also want to panic if our count goes below 0 to help with
// debugging. All this is faster than model.AutoRelease if error... blah.
// every time.
func (m *baseModel) setReferenceCount(rc uint64) {
    atomic.StoreInt64(&m.refCount, int64(rc))
}
