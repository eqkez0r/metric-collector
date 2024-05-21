package result

import "sync/atomic"

type Result struct {
	all, errors int64
}

func New() *Result {
	return &Result{
		all:    0,
		errors: 0,
	}
}

func (r *Result) All() int64 {
	return atomic.LoadInt64(&r.all)
}

func (r *Result) Errors() int64 {
	return atomic.LoadInt64(&r.errors)
}

func (r *Result) IncAll() {
	atomic.AddInt64(&r.all, 1)
}

func (r *Result) IncErrors() {
	atomic.AddInt64(&r.errors, 1)
}

func (r *Result) Reset() {
	r.all = 0
	r.errors = 0
}
