// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package handlers

import (
	"context"
	"sync"

	"github.com/Eqke/metric-collector/pkg/metric"
)

// Ensure, that NewJSONMetricProviderMock does implement NewJSONMetricProvider.
// If this is not the case, regenerate this file with moq.
var _ NewJSONMetricProvider = &NewJSONMetricProviderMock{}

// NewJSONMetricProviderMock is a mock implementation of NewJSONMetricProvider.
//
//	func TestSomethingThatUsesNewJSONMetricProvider(t *testing.T) {
//
//		// make and configure a mocked NewJSONMetricProvider
//		mockedNewJSONMetricProvider := &NewJSONMetricProviderMock{
//			SetMetricFunc: func(contextMoqParam context.Context, metrics metric.Metrics) error {
//				panic("mock out the SetMetric method")
//			},
//		}
//
//		// use mockedNewJSONMetricProvider in code that requires NewJSONMetricProvider
//		// and then make assertions.
//
//	}
type NewJSONMetricProviderMock struct {
	// SetMetricFunc mocks the SetMetric method.
	SetMetricFunc func(contextMoqParam context.Context, metrics metric.Metrics) error

	// calls tracks calls to the methods.
	calls struct {
		// SetMetric holds details about calls to the SetMetric method.
		SetMetric []struct {
			// ContextMoqParam is the contextMoqParam argument value.
			ContextMoqParam context.Context
			// Metrics is the metrics argument value.
			Metrics metric.Metrics
		}
	}
	lockSetMetric sync.RWMutex
}

// SetMetric calls SetMetricFunc.
func (mock *NewJSONMetricProviderMock) SetMetric(contextMoqParam context.Context, metrics metric.Metrics) error {
	if mock.SetMetricFunc == nil {
		panic("NewJSONMetricProviderMock.SetMetricFunc: method is nil but NewJSONMetricProvider.SetMetric was just called")
	}
	callInfo := struct {
		ContextMoqParam context.Context
		Metrics         metric.Metrics
	}{
		ContextMoqParam: contextMoqParam,
		Metrics:         metrics,
	}
	mock.lockSetMetric.Lock()
	mock.calls.SetMetric = append(mock.calls.SetMetric, callInfo)
	mock.lockSetMetric.Unlock()
	return mock.SetMetricFunc(contextMoqParam, metrics)
}

// SetMetricCalls gets all the calls that were made to SetMetric.
// Check the length with:
//
//	len(mockedNewJSONMetricProvider.SetMetricCalls())
func (mock *NewJSONMetricProviderMock) SetMetricCalls() []struct {
	ContextMoqParam context.Context
	Metrics         metric.Metrics
} {
	var calls []struct {
		ContextMoqParam context.Context
		Metrics         metric.Metrics
	}
	mock.lockSetMetric.RLock()
	calls = mock.calls.SetMetric
	mock.lockSetMetric.RUnlock()
	return calls
}