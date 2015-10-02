// Copyright 2015 Philipp Brüll <bruell@simia.tech>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package exit

import (
	"os"
	"os/signal"
	"sync"
	"time"
)

// Exit defines an exit that contains multiple SignalChans.
type Exit struct {
	Name string

	signals      []*Signal
	signalsMutex sync.Mutex

	timeout     time.Duration
	afterEachFn func(error)
	afterAllFn  func(*Report)
}

// New returns a new exit with the provided name.
func New(name string) *Exit {
	return &Exit{
		Name:    name,
		signals: []*Signal{},
	}
}

// SetTimeout sets a timeout for the actors to end during the exit process.
func (e *Exit) SetTimeout(value time.Duration) {
	e.timeout = value
}

// HasTimeout returns true if a timeout is set.
func (e *Exit) HasTimeout() bool {
	return e.timeout != 0
}

// AfterEach registers a callback function that is called after each actor
// that has exited.
func (e *Exit) AfterEach(fn func(error)) {
	e.afterEachFn = fn
}

// AfterAll registers a callback function that is called after all actors
// has exited.
func (e *Exit) AfterAll(fn func(*Report)) {
	e.afterAllFn = fn
}

// NewSignal creates a new Signal, attaches it to the exit and returns it.
func (e *Exit) NewSignal(name string) *Signal {
	e.signalsMutex.Lock()
	defer e.signalsMutex.Unlock()

	signal := NewSignal(name)
	e.signals = append(e.signals, signal)
	return signal
}

// Exit sends an ErrChan through all the previously generated SignalChans
// and waits until all returned an error or nil. The received errors will be
// returned in an error report.
func (e *Exit) Exit() *Report {
	e.signalsMutex.Lock()
	defer e.signalsMutex.Unlock()

	report := NewReport(e.Name)
	wg := &sync.WaitGroup{}
	for _, signal := range e.signals {
		wg.Add(1)
		go func(signal *Signal) {
			if e.HasTimeout() && !signal.HasTimeout() {
				signal.SetTimeout(e.timeout)
			}
			err := signal.Exit()
			if err != nil {
				report.Set(signal.Name, err)
			}
			if e.afterEachFn != nil {
				e.afterEachFn(err)
			}
			wg.Done()
		}(signal)
	}
	wg.Wait()
	e.signals = []*Signal{}

	if e.afterAllFn != nil {
		e.afterAllFn(report)
	}

	if report.Len() == 0 {
		return nil
	}
	return report
}

// ExitOn blocks until the process receives one of the provided signals and
// than calls Exit.
func (e *Exit) ExitOn(osSignales ...os.Signal) *Report {
	osSignalChan := make(chan os.Signal)
	signal.Notify(osSignalChan, osSignales...)
	<-osSignalChan

	return e.Exit()
}
