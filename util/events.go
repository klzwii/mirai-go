package util

import (
	"errors"
	"runtime"
	"sync"
	"sync/atomic"
)

type EventCenter interface {
	// RegisterEvent to event center, return unique identifier for current event
	RegisterEvent() (uint32, chan *Result)
	// Notify the event with id, and pass param in to it
	Notify(id uint32, in any, err error) error
}

type Result struct {
	Data any
	Err  error
}

var pool = &sync.Pool{
	New: func() any {
		return &event{ch: make(chan *Result, 1)}
	},
}

type event struct {
	ch chan *Result
	id uint32
}

type state struct {
	head uint32
	size uint32
}

var (
	placeHolderDelete    = &event{}
	placeHolderAdd       = &event{}
	EventOutOfRangeError = errors.New("event out of range")
	EventHandleTwice     = errors.New("event handled twice")
)

type eventCenterImp struct {
	mu     *sync.Mutex
	cond   *sync.Cond
	cap    uint32
	events []atomic.Pointer[event]
	state  atomic.Pointer[state]
}

func (e2 *eventCenterImp) RegisterEvent() (uint32, chan *Result) {
	curState := e2.state.Load()
	nSTate := &state{
		head: curState.head,
		size: curState.size + 1,
	}
	for curState.size == e2.cap || !e2.state.CompareAndSwap(curState, nSTate) {
		if curState.size == e2.cap {
			for curState.size == e2.cap {
				e2.notifySlow()
				curState = e2.state.Load()
				runtime.Gosched()
			}
		} else {
			curState = e2.state.Load()
		}
		nSTate.head = curState.head
		nSTate.size = curState.size + 1
	}
	id := nSTate.head + nSTate.size - 1
	newEvent := pool.Get().(*event)
	retCh := newEvent.ch
	newEvent.id = id
	oldValue := e2.events[id%e2.cap].Swap(newEvent)
	if oldValue != placeHolderAdd {
		panic("consistency check fail, all new value should be put on placeholder add")
	}
	return id, retCh
}

func (e2 *eventCenterImp) Notify(id uint32, in any, err error) error {
	curState := e2.state.Load()
	if id < curState.head || id >= curState.head+curState.size {
		return EventOutOfRangeError
	}
	t := e2.events[id%e2.cap].Swap(placeHolderDelete)
	if t == placeHolderDelete {
		return EventHandleTwice
	}
	if t == placeHolderAdd {
		panic("consistency check fail, element in queue should not be placeholder add")
	}
	t.ch <- &Result{
		Data: in,
		Err:  err,
	}
	close(t.ch)
	t.ch = make(chan *Result, 1)
	pool.Put(t)
	return nil
}

func (e2 *eventCenterImp) notifySlow() uint32 {
	if !e2.mu.TryLock() {
		return 0
	}
	defer e2.mu.Unlock()
	curState := e2.state.Load()
	var eraseSize = uint32(0)
	for curState.size > eraseSize && e2.events[(curState.head+eraseSize)%e2.cap].CompareAndSwap(placeHolderDelete, placeHolderAdd) {
		eraseSize += 1
	}
	nState := &state{
		head: curState.head + eraseSize,
		size: curState.size - eraseSize,
	}
	for !e2.state.CompareAndSwap(curState, nState) {
		curState = e2.state.Load()
		for curState.size > eraseSize && e2.events[(curState.head+eraseSize)%e2.cap].CompareAndSwap(placeHolderDelete, placeHolderAdd) {
			eraseSize += 1
		}
		nState.size = curState.size - eraseSize
		nState.head = curState.head + eraseSize
	}
	e2.cond.Broadcast()
	return eraseSize
}

func New(cap uint32) EventCenter {
	center := &eventCenterImp{
		mu:     &sync.Mutex{},
		cond:   sync.NewCond(&sync.Mutex{}),
		cap:    cap,
		events: make([]atomic.Pointer[event], cap),
		state:  atomic.Pointer[state]{},
	}
	center.state.Swap(&state{
		head: 1,
		size: 0,
	})
	for i := uint32(0); i < cap; i++ {
		center.events[i].Store(placeHolderAdd)
	}
	return center
}
