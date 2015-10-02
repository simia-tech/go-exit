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

import "time"

// Signal defines an exit signal signal for a single actor.
type Signal struct {
	Name string
	Chan chan Reply

	timeout time.Duration
}

// NewSignal returns a new initialized signal with the provided name.
func NewSignal(name string) *Signal {
	return &Signal{
		Name: name,
		Chan: make(chan Reply),
	}
}

// SetTimeout sets the timeout for this specific exit signal to complete. If no
// or zero timeout is set, the exit process can last forever.
func (s *Signal) SetTimeout(value time.Duration) {
	s.timeout = value
}

// HasTimeout returns true if a timeout is set.
func (s *Signal) HasTimeout() bool {
	return s.timeout != 0
}

// Exit performs the exit process for this specific signal.
func (s *Signal) Exit() error {
	reply := make(Reply)

	if !s.HasTimeout() {
		s.Chan <- reply
		return <-reply
	}

	select {
	case s.Chan <- reply:
	case <-time.After(s.timeout):
		return ErrTimeout
	}

	select {
	case err := <-reply:
		return err
	case <-time.After(s.timeout):
		return ErrTimeout
	}
}
