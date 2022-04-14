package pwsmodels

// ModelPool a pool of models :)
type ModelPool[T pooledItem] struct {
    Pool[T]
}

// newModelFunc fenerics wrapper for calling new(some_kinda_item) as new(T)
// does not work as expected.
type newModelFunc[T pooledItem]func() T

// newModelPool returns a model pool which will automatically register
// items upon creation.
func newModelPool[T pooledItem](nf newModelFunc[T]) *ModelPool[T] {

    var n = new(ModelPool[T])
    n.pool.New = func() any {
        var o = nf()
        o.register(o, &n.pool)
        return o
    }

    return n
}

// AutoAcquire returns a Model which is reference counted.
func (this *ModelPool[T]) AutoAcquire(rc uint64) T {
    var n = this.pool.Get().(T)
    n.setReferenceCount(rc)
    return n
}
