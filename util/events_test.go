package util

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"runtime"
	"sort"
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
		err := center.Notify(ca.id, ca.id, nil)
		assert.Nil(t, err)
	}
	wg.Wait()
	for _, a := range ck {
		assert.Equal(t, a.Load(), true)
	}
}

type PlainEventCenter struct {
	data sync.Map
	cnt  atomic.Uint32
}

func (p *PlainEventCenter) RegisterEvent() (uint32, chan *Result) {
	ch := make(chan *Result)
	id := p.cnt.Add(1)
	p.data.Store(id, ch)
	return id, ch
}

func (p *PlainEventCenter) Notify(id uint32, in any, _ error) error {
	val, ok := p.data.LoadAndDelete(id)
	result := &Result{
		Data: in,
	}
	if ok {
		val.(chan *Result) <- result
	}
	return nil
}

type operation struct {
	ty bool
	id uint32
}

type generateTemp struct {
	id    uint32
	order int
}

type sortContainer []*generateTemp

func (s sortContainer) Len() int {
	return len(s)
}

func (s sortContainer) Less(i, j int) bool {
	return s[i].order < s[j].order
}

func (s sortContainer) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func generateTestData(n int) []*operation {
	var k sortContainer
	for i := 0; i < n; i++ {
		k = append(k, &generateTemp{
			id:    uint32(i + 1),
			order: rand.Int()%(n-i) + i,
		})
	}
	sort.Sort(k)
	var ret []*operation
	j := 0
	for i := 0; i < n; i++ {
		ret = append(ret, &operation{
			ty: false,
			id: uint32(i + 1),
		})
		for ; j < n; j++ {
			if k[j].order > i {
				break
			}
			ret = append(ret, &operation{
				ty: true,
				id: k[j].id,
			})
		}
	}
	return ret
}

func generateSequentialTestData(n int) []*operation {
	var ret []*operation
	for i := 0; i < n; i++ {
		ret = append(ret, &operation{
			ty: false,
			id: uint32(i + 1),
		})
		ret = append(ret, &operation{
			ty: true,
			id: uint32(i + 1),
		})
	}
	return ret
}

func innerBench(b *testing.B, e EventCenter, generator func(n int) []*operation) {
	b.StopTimer()
	ope := generator(20000)
	b.StartTimer()
	for _, o := range ope {
		if o.ty {
			_ = e.Notify(o.id, o.id, nil)
		} else {
			id, ch := e.RegisterEvent()
			assert.Equal(b, id, o.id)
			//fmt.Println(id, o.id)
			go func(ch chan *Result, id uint32) {
				res := <-ch
				assert.Equal(b, id, res.Data.(uint32))
			}(ch, id)
		}
	}
}

func innerBenchParallel(b *testing.B, e EventCenter, generator func(n int) []*operation) {
	b.StopTimer()
	ope := generator(20000)
	var mus [20000]atomic.Bool
	b.StartTimer()
	for _, o := range ope {
		go func(o *operation) {
			if o.ty {
				for !mus[o.id-1].CompareAndSwap(true, false) {
					runtime.Gosched()
				}
				_ = e.Notify(o.id, o.id, nil)
			} else {
				id, ch := e.RegisterEvent()
				mus[id-1].Store(true)
				go func(ch chan *Result, id uint32) {
					res := <-ch
					assert.Equal(b, id, res.Data.(uint32))
					//println(id)
				}(ch, id)
			}
		}(o)
	}
}

func BenchmarkEventCenter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		e := New(20000)
		innerBench(b, e, generateTestData)
	}
}

func BenchmarkPlain(b *testing.B) {
	for i := 0; i < b.N; i++ {
		e := &PlainEventCenter{}
		innerBench(b, e, generateTestData)
	}
}

func BenchmarkEventCenterSequential(b *testing.B) {
	for i := 0; i < b.N; i++ {
		e := New(2000)
		innerBench(b, e, generateSequentialTestData)
	}
}

func BenchmarkPlainSequential(b *testing.B) {
	for i := 0; i < b.N; i++ {
		e := &PlainEventCenter{}
		innerBench(b, e, generateSequentialTestData)
	}
}

func BenchmarkEventCenterParallel(b *testing.B) {
	for i := 0; i < b.N; i++ {
		e := New(20000)
		innerBenchParallel(b, e, generateSequentialTestData)
	}
}

func BenchmarkPlainParallel(b *testing.B) {
	for i := 0; i < b.N; i++ {
		e := &PlainEventCenter{}
		innerBenchParallel(b, e, generateSequentialTestData)
	}
}
