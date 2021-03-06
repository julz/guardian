// This file was generated by counterfeiter
package fakes

import (
	"sync"

	"github.com/cloudfoundry-incubator/garden"
	"github.com/cloudfoundry-incubator/guardian/rundmc"
)

type FakeActualContainer struct {
	RunStub        func(spec garden.ProcessSpec, io garden.ProcessIO) (garden.Process, error)
	runMutex       sync.RWMutex
	runArgsForCall []struct {
		spec garden.ProcessSpec
		io   garden.ProcessIO
	}
	runReturns struct {
		result1 garden.Process
		result2 error
	}
}

func (fake *FakeActualContainer) Run(spec garden.ProcessSpec, io garden.ProcessIO) (garden.Process, error) {
	fake.runMutex.Lock()
	fake.runArgsForCall = append(fake.runArgsForCall, struct {
		spec garden.ProcessSpec
		io   garden.ProcessIO
	}{spec, io})
	fake.runMutex.Unlock()
	if fake.RunStub != nil {
		return fake.RunStub(spec, io)
	} else {
		return fake.runReturns.result1, fake.runReturns.result2
	}
}

func (fake *FakeActualContainer) RunCallCount() int {
	fake.runMutex.RLock()
	defer fake.runMutex.RUnlock()
	return len(fake.runArgsForCall)
}

func (fake *FakeActualContainer) RunArgsForCall(i int) (garden.ProcessSpec, garden.ProcessIO) {
	fake.runMutex.RLock()
	defer fake.runMutex.RUnlock()
	return fake.runArgsForCall[i].spec, fake.runArgsForCall[i].io
}

func (fake *FakeActualContainer) RunReturns(result1 garden.Process, result2 error) {
	fake.RunStub = nil
	fake.runReturns = struct {
		result1 garden.Process
		result2 error
	}{result1, result2}
}

var _ rundmc.ActualContainer = new(FakeActualContainer)
