package util

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var TimeOutError = errors.New("time out")

type executor[T any] interface {
	Process(resolve func(T), reject func(err error))
}

type promise[T any] struct {
	executor executor[T]
	cond     *sync.Cond
	finished bool
	res      T
	err      error
}

type executorImp[T any] struct {
	process func(resolve func(T), reject func(err error))
}

func (e *executorImp[T]) Process(resolve func(T), reject func(err error)) {
	e.process(resolve, reject)
}

func FutureFunc[T any](process func(resolve func(T), reject func(error))) *promise[T] {
	return &promise[T]{
		executor: &executorImp[T]{
			process: process,
		},
		cond: sync.NewCond(&sync.Mutex{}),
		err:  nil,
	}
}

func (p *promise[T]) Start() {
	p.StartWithTimeOut(time.Hour * 120)
}

func (p *promise[T]) StartWithTimeOut(timeOut time.Duration) {
	go func() {
		timer := time.NewTimer(timeOut)
		p.cond.L.Lock()
		defer func() {
			err := recover()
			if err != nil {
				p.err = fmt.Errorf("%v", err)
			}
			p.cond.L.Unlock()
		}()
		if p.finished {
			return
		}
		resChannel := make(chan T, 1)
		errChannel := make(chan error, 1)
		go func() {
			p.executor.Process(func(res T) {
				resChannel <- res
			}, func(err error) {
				errChannel <- err
			})
		}()
		select {
		case res := <-resChannel:
			p.res = res
		case err := <-errChannel:
			p.err = err
		case <-timer.C:
			p.err = TimeOutError
		}
		p.finished = true
		p.cond.Broadcast()
	}()
}

// Await wait until function is finished and return
func (p *promise[T]) Await() (T, error) {
	p.cond.L.Lock()
	defer p.cond.L.Unlock()
	for p.finished != true {
		p.cond.Wait()
	}
	return p.res, p.err
}
