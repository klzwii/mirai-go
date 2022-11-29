package util

import (
	"errors"
	"sync"
	"sync/atomic"
)

type EventCenter interface {
	// RegisterEvent to event center, return unique identifier for current event
	RegisterEvent() (uint32, chan *Result)
	// Notify the event with id, and pass param in to it
	Notify(id uint32, in any) error
}

type Result struct {
	Data any
	Err  error
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
		curState = e2.state.Load()
		nSTate.head = curState.head
		nSTate.size = curState.size + 1
	}
	t := e2.state.Load()
	println(t)
	id := nSTate.head + nSTate.size - 1
	retCh := make(chan *Result, 1)
	oldValue := e2.events[id%e2.cap].Swap(&event{
		ch: retCh,
		id: id,
	})
	if oldValue != placeHolderAdd {
		panic("consistency check fail, all new value should be put on placeholder add")
	}
	return id, retCh
}

func (e2 *eventCenterImp) Notify(id uint32, in any) error {
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
		Err:  nil,
	}
	close(t.ch)
	if e2.events[curState.head%e2.cap].Load() == placeHolderDelete {
		e2.mu.Lock()
		defer e2.mu.Unlock()
		curState = e2.state.Load()
		nState := &state{
			head: curState.head,
			size: curState.size,
		}
		var eraseSize = uint32(0)
		for nState.size > 0 && e2.events[nState.head%e2.cap].CompareAndSwap(placeHolderDelete, placeHolderAdd) {
			eraseSize += 1
		}
		nState.head += eraseSize
		nState.size -= eraseSize
		for !e2.state.CompareAndSwap(curState, nState) {
			curState = e2.state.Load()
			nState.head = curState.head + eraseSize
			nState.size = curState.size - eraseSize
		}
	}
	return nil
}

func New(cap uint32) EventCenter {
	center := &eventCenterImp{
		mu:     &sync.Mutex{},
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
