package promise

import "sync/atomic"

type Error = Promise[error]

func NewError() *Error {
	return New[error]()
}

type Promise[T any] struct {
	awaiter chan struct{}
	value   atomic.Pointer[T]
}

func New[T any]() *Promise[T] {
	return &Promise[T]{
		awaiter: make(chan struct{}),
	}
}

func (p *Promise[T]) Awaiter() <-chan struct{} {
	return p.awaiter
}

func (p *Promise[T]) Await() {
	<-p.awaiter
}

func (p *Promise[T]) Get() (T, bool) {
	value := p.value.Load()
	if value == nil {
		return *new(T), false
	}

	return *value, true
}

func (p *Promise[T]) MustGet() T {
	value := p.value.Load()
	if value == nil {
		panic("promise: value is not awaited")
	}

	return *value
}

func (p *Promise[T]) AwaitAndGet() T {
	p.Await()
	return *p.value.Load()
}

func (p *Promise[T]) Set(value T) bool {
	if !p.value.CompareAndSwap(nil, &value) {
		return false
	}

	close(p.awaiter)
	return true
}
