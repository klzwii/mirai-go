package util

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"sync/atomic"
	"testing"
)

func TestEventCenter(t *testing.T) {
	center := New(1000)
	type temp struct {
		id uint32
		ch chan *Result
	}
	ch := make(chan *temp)
	wg := &sync.WaitGroup{}
	ck := make([]*atomic.Bool, 200000)
	for i := 0; i < 200000; i++ {
		ck[i] = &atomic.Bool{}
		ck[i].Store(false)
	}
	go func() {
		for i := 0; i < 200000; i++ {
			id, resCh := center.RegisterEvent()
			ch <- &temp{id, resCh}
			wg.Add(1)
			go func(id uint32, resCh chan *Result) {
				c := <-resCh
				ck[id-1].Store(true)
				assert.Equal(t, id, c.Data.(uint32))
				wg.Done()
			}(id, resCh)
		}
		close(ch)
	}()
	for i := 0; i < 200000; i++ {
		ca := <-ch
		err := center.Notify(ca.id, ca.id)
		assert.Nil(t, err)
	}
	wg.Wait()
	for _, a := range ck {
		assert.Equal(t, a.Load(), true)
	}
}
