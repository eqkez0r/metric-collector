// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package handlers

import (
	"context"
	"sync"

	store "github.com/Eqke/metric-collector/internal/storage"
)

// Ensure, that RootMetricsProviderMock does implement RootMetricsProvider.
// If this is not the case, regenerate this file with moq.
var _ RootMetricsProvider = &RootMetricsProviderMock{}

// RootMetricsProviderMock is a mock implementation of RootMetricsProvider.
//
//	func TestSomethingThatUsesRootMetricsProvider(t *testing.T) {
//
//		// make and configure a mocked RootMetricsProvider
//		mockedRootMetricsProvider := &RootMetricsProviderMock{
//			GetMetricsFunc: func(contextMoqParam context.Context) (map[string][]store.Metric, error) {
//				panic("mock out the GetMetrics method")
//			},
//		}
//
//		// use mockedRootMetricsProvider in code that requires RootMetricsProvider
//		// and then make assertions.
//
//	}
type RootMetricsProviderMock struct {
	// GetMetricsFunc mocks the GetMetrics method.
	GetMetricsFunc func(contextMoqParam context.Context) (map[string][]store.Metric, error)

	// calls tracks calls to the methods.
	calls struct {
		// GetMetrics holds details about calls to the GetMetrics method.
		GetMetrics []struct {
			// ContextMoqParam is the contextMoqParam argument value.
			ContextMoqParam context.Context
		}
	}
	lockGetMetrics sync.RWMutex
}

// GetMetrics calls GetMetricsFunc.
func (mock *RootMetricsProviderMock) GetMetrics(contextMoqParam context.Context) (map[string][]store.Metric, error) {
	if mock.GetMetricsFunc == nil {
		panic("RootMetricsProviderMock.GetMetricsFunc: method is nil but RootMetricsProvider.GetMetrics was just called")
	}
	callInfo := struct {
		ContextMoqParam context.Context
	}{
		ContextMoqParam: contextMoqParam,
	}
	mock.lockGetMetrics.Lock()
	mock.calls.GetMetrics = append(mock.calls.GetMetrics, callInfo)
	mock.lockGetMetrics.Unlock()
	return mock.GetMetricsFunc(contextMoqParam)
}

// GetMetricsCalls gets all the calls that were made to GetMetrics.
// Check the length with:
//
//	len(mockedRootMetricsProvider.GetMetricsCalls())
func (mock *RootMetricsProviderMock) GetMetricsCalls() []struct {
	ContextMoqParam context.Context
} {
	var calls []struct {
		ContextMoqParam context.Context
	}
	mock.lockGetMetrics.RLock()
	calls = mock.calls.GetMetrics
	mock.lockGetMetrics.RUnlock()
	return calls
}