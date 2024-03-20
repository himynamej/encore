// Package metrics provides support to update metric values.
package metrics

import (
	"expvar"
	"runtime"

	emetrics "encore.dev/metrics"
)

// Config lists the set of metrics that is tracked.
type Config struct {
	Goroutines *emetrics.Gauge[uint64]
	Requests   *emetrics.Counter[uint64]
	Failures   *emetrics.Counter[uint64]
	Panics     *emetrics.Counter[uint64]
}

// Values provides an api to work with metrics.
type Values struct {
	goroutines *emetrics.Gauge[uint64]
	requests   *emetrics.Counter[uint64]
	failures   *emetrics.Counter[uint64]
	panics     *emetrics.Counter[uint64]
	myRequests *expvar.Int
}

// New constructs a Values for working with metrics.
func New(cfg Config) *Values {
	return &Values{
		goroutines: cfg.Goroutines,
		requests:   cfg.Requests,
		failures:   cfg.Failures,
		panics:     cfg.Panics,
		myRequests: expvar.NewInt("requests"),
	}
}

// SetGoroutines updates the number of goroutines.
func (v *Values) SetGoroutines() {
	n := runtime.NumGoroutine()
	v.goroutines.Set(uint64(n))
}

// IncRequests increments the requests by 1.
func (v *Values) IncRequests() int64 {
	v.requests.Add(1)
	v.myRequests.Add(1)

	return v.myRequests.Value()
}

// IncFailures increments the failures by 1.
func (v *Values) IncFailures() {
	v.failures.Add(1)
}

// IncPanics increments the panics by 1.
func (v *Values) IncPanics() {
	v.panics.Add(1)
}