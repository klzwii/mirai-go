package util

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"sync"
	"testing"
	"time"
)

func TestFuturePromise(t *testing.T) {
	suite.Run(t, new(FuturePromiseTestSuite))
}

type FuturePromiseTestSuite struct {
	suite.Suite
}

var testErr = errors.New("test error")

func (f *FuturePromiseTestSuite) TestNormalFuturePromise() {
	future := FutureFunc(func(resolve func(ret int), reject func(error)) {
		time.Sleep(time.Second * 2)
		resolve(10)
	})
	future.Start()
	wg := &sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			now := time.Now()
			ans, err := future.Await()
			cost := time.Now().Sub(now).Seconds()
			wg.Done()
			assert.Nil(f.T(), err)
			assert.Equal(f.T(), 10, ans)
			assert.Less(f.T(), 2.0, cost)
			assert.Greater(f.T(), 2.1, cost)
		}()
	}
	wg.Wait()
}

func (f *FuturePromiseTestSuite) TestFailFast() {
	future := FutureFunc(func(resolve func(ret int), reject func(error)) {
		go func() {
			time.Sleep(time.Second)
			reject(testErr)
		}()
		time.Sleep(time.Second * 2)
		resolve(10)
	})
	future.Start()
	wg := &sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			now := time.Now()
			_, err := future.Await()
			cost := time.Now().Sub(now).Seconds()
			wg.Done()
			assert.ErrorIs(f.T(), testErr, err)
			assert.Less(f.T(), 1.0, cost)
			assert.Greater(f.T(), 1.1, cost)
		}()
	}
	wg.Wait()
}

func (f *FuturePromiseTestSuite) TestFailWithTimeOut() {
	future := FutureFunc(func(resolve func(ret int), reject func(error)) {
		time.Sleep(time.Second * 100)
	})
	future.StartWithTimeOut(time.Second * 1)
	wg := &sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			now := time.Now()
			_, err := future.Await()
			cost := time.Now().Sub(now).Seconds()
			wg.Done()
			assert.ErrorIs(f.T(), TimeOutError, err)
			assert.Less(f.T(), 1.0, cost)
			assert.Greater(f.T(), 1.1, cost)
		}()
	}
	wg.Wait()
}
