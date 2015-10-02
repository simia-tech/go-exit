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

// ErrChan defines a channel of errors that can be used to deliver back
// an error after an actor has shut down.
type ErrChan chan error

// SignalChan defines a channel of ErrChan that is used to signal an
// actor to shut down.
type SignalChan chan ErrChan

// Exit defines an exit that contains multiple SignalChans.
type Exit struct {
	Name string

	signalChans      map[string]SignalChan
	signalChansMutex sync.Mutex

	timeout time.Duration
}

// New returns a new exit with the provided name.
func New(name string) *Exit {
	return &Exit{
		Name:        name,
		signalChans: make(map[string]SignalChan),
	}
}

// SetTimeout sets a timeout for the actors to end during the exit process.
func (e *Exit) SetTimeout(value time.Duration) {
	e.timeout = value
}

// NewSignalChan creates a new SignalChan and returns it.
func (e *Exit) NewSignalChan(name string) (SignalChan, error) {
	e.signalChansMutex.Lock()
	defer e.signalChansMutex.Unlock()

	if _, ok := e.signalChans[name]; ok {
		return nil, ErrNameAlreadyExists
	}

	signalChan := make(SignalChan, 1)
	e.signalChans[name] = signalChan
	return signalChan, nil
}

// Exit sends an ErrChan through all the previously generated SignalChans
// and waits until all returned an error or nil. The received errors will be
// returned in an error report.
func (e *Exit) Exit() *Report {
	e.signalChansMutex.Lock()
	defer e.signalChansMutex.Unlock()

	report := NewReport(e.Name)
	wg := &sync.WaitGroup{}
	for name, signalChan := range e.signalChans {
		wg.Add(1)
		go func(name string, signalChan SignalChan) {
			if err := e.exit(name, signalChan); err != nil {
				report.Set(name, err)
			}
			wg.Done()
		}(name, signalChan)
		delete(e.signalChans, name)
	}
	wg.Wait()

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

func (e *Exit) exit(name string, signalChan SignalChan) error {
	errChan := make(ErrChan)
	signalChan <- errChan

	if e.timeout == 0 {
		return <-errChan
	}

	select {
	case err := <-errChan:
		return err
	case <-time.After(e.timeout):
		return ErrTimeout
	}
}
