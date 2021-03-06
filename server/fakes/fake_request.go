// This file was generated by counterfeiter
package fakes

import (
	"sync"

	"github.com/Bo0mer/os-agent/server"
)

type FakeRequest struct {
	BodyStub        func() []byte
	bodyMutex       sync.RWMutex
	bodyArgsForCall []struct{}
	bodyReturns     struct {
		result1 []byte
	}
	ParamValuesStub        func(string) ([]string, bool)
	paramValuesMutex       sync.RWMutex
	paramValuesArgsForCall []struct {
		arg1 string
	}
	paramValuesReturns struct {
		result1 []string
		result2 bool
	}
}

func (fake *FakeRequest) Body() []byte {
	fake.bodyMutex.Lock()
	fake.bodyArgsForCall = append(fake.bodyArgsForCall, struct{}{})
	fake.bodyMutex.Unlock()
	if fake.BodyStub != nil {
		return fake.BodyStub()
	} else {
		return fake.bodyReturns.result1
	}
}

func (fake *FakeRequest) BodyCallCount() int {
	fake.bodyMutex.RLock()
	defer fake.bodyMutex.RUnlock()
	return len(fake.bodyArgsForCall)
}

func (fake *FakeRequest) BodyReturns(result1 []byte) {
	fake.BodyStub = nil
	fake.bodyReturns = struct {
		result1 []byte
	}{result1}
}

func (fake *FakeRequest) ParamValues(arg1 string) ([]string, bool) {
	fake.paramValuesMutex.Lock()
	fake.paramValuesArgsForCall = append(fake.paramValuesArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.paramValuesMutex.Unlock()
	if fake.ParamValuesStub != nil {
		return fake.ParamValuesStub(arg1)
	} else {
		return fake.paramValuesReturns.result1, fake.paramValuesReturns.result2
	}
}

func (fake *FakeRequest) ParamValuesCallCount() int {
	fake.paramValuesMutex.RLock()
	defer fake.paramValuesMutex.RUnlock()
	return len(fake.paramValuesArgsForCall)
}

func (fake *FakeRequest) ParamValuesArgsForCall(i int) string {
	fake.paramValuesMutex.RLock()
	defer fake.paramValuesMutex.RUnlock()
	return fake.paramValuesArgsForCall[i].arg1
}

func (fake *FakeRequest) ParamValuesReturns(result1 []string, result2 bool) {
	fake.ParamValuesStub = nil
	fake.paramValuesReturns = struct {
		result1 []string
		result2 bool
	}{result1, result2}
}

var _ server.Request = new(FakeRequest)
