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

var (
	signalChans      = make(map[string]SignalChan)
	signalChansMutex = &sync.Mutex{}

	timeout time.Duration
)

// SetTimeout sets a timeout for the actors to end during the exit process.
func SetTimeout(value time.Duration) {
	timeout = value
}

// NewSignalChan creates a new SignalChan and returns it.
func NewSignalChan(name string) (SignalChan, error) {
	signalChansMutex.Lock()
	defer signalChansMutex.Unlock()

	if _, ok := signalChans[name]; ok {
		return nil, ErrNameAlreadyExists
	}

	signalChan := make(SignalChan, 1)
	signalChans[name] = signalChan
	return signalChan, nil
}

// Exit sends an ErrChan through all the previously generated SignalChans
// and waits until all returned an error or nil. The received errors will be
// returned in an error report.
func Exit() *Report {
	signalChansMutex.Lock()
	defer signalChansMutex.Unlock()

	report := NewReport()
	wg := &sync.WaitGroup{}
	for name, signalChan := range signalChans {
		wg.Add(1)
		go func(name string, signalChan SignalChan) {
			if err := exit(name, signalChan); err != nil {
				report.Set(name, err)
			}
			wg.Done()
		}(name, signalChan)
		delete(signalChans, name)
	}
	wg.Wait()

	if report.Len() == 0 {
		return nil
	}
	return report
}

// Reset removes all created SignalChans without sending the exit signal.
func Reset() {
	signalChansMutex.Lock()
	defer signalChansMutex.Unlock()

	signalChans = make(map[string]SignalChan)
}

// ExitOn blocks until the process receives one of the provided signals and
// than calls Exit.
func ExitOn(osSignales ...os.Signal) *Report {
	osSignalChan := make(chan os.Signal)
	signal.Notify(osSignalChan, osSignales...)
	<-osSignalChan

	return Exit()
}

func exit(name string, signalChan SignalChan) error {
	errChan := make(ErrChan)
	signalChan <- errChan

	if timeout == 0 {
		return <-errChan
	}

	select {
	case err := <-errChan:
		return err
	case <-time.After(timeout):
		return ErrTimeout
	}
}
