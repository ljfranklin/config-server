// This file was generated by counterfeiter
package fakes

import (
	"config_server/types"
	"crypto/rsa"
	"crypto/x509"
	"sync"
)

type FakeCertsLoader struct {
	LoadCertsStub        func(certFile, keyFile string) (*x509.Certificate, *rsa.PrivateKey, error)
	loadCertsMutex       sync.RWMutex
	loadCertsArgsForCall []struct {
		certFile string
		keyFile  string
	}
	loadCertsReturns struct {
		result1 *x509.Certificate
		result2 *rsa.PrivateKey
		result3 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeCertsLoader) LoadCerts(certFile string, keyFile string) (*x509.Certificate, *rsa.PrivateKey, error) {
	fake.loadCertsMutex.Lock()
	fake.loadCertsArgsForCall = append(fake.loadCertsArgsForCall, struct {
		certFile string
		keyFile  string
	}{certFile, keyFile})
	fake.recordInvocation("LoadCerts", []interface{}{certFile, keyFile})
	fake.loadCertsMutex.Unlock()
	if fake.LoadCertsStub != nil {
		return fake.LoadCertsStub(certFile, keyFile)
	} else {
		return fake.loadCertsReturns.result1, fake.loadCertsReturns.result2, fake.loadCertsReturns.result3
	}
}

func (fake *FakeCertsLoader) LoadCertsCallCount() int {
	fake.loadCertsMutex.RLock()
	defer fake.loadCertsMutex.RUnlock()
	return len(fake.loadCertsArgsForCall)
}

func (fake *FakeCertsLoader) LoadCertsArgsForCall(i int) (string, string) {
	fake.loadCertsMutex.RLock()
	defer fake.loadCertsMutex.RUnlock()
	return fake.loadCertsArgsForCall[i].certFile, fake.loadCertsArgsForCall[i].keyFile
}

func (fake *FakeCertsLoader) LoadCertsReturns(result1 *x509.Certificate, result2 *rsa.PrivateKey, result3 error) {
	fake.LoadCertsStub = nil
	fake.loadCertsReturns = struct {
		result1 *x509.Certificate
		result2 *rsa.PrivateKey
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeCertsLoader) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.loadCertsMutex.RLock()
	defer fake.loadCertsMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeCertsLoader) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ types.CertsLoader = new(FakeCertsLoader)
