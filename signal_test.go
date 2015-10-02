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

package exit_test

import (
	"testing"
	"time"

	"github.com/simia-tech/go-exit"
)

func TestSignalExitWithoutTimeout(t *testing.T) {
	exitSignal := exit.NewSignal("one")
	go func() {
		reply := <-exitSignal.Chan
		reply.Ok()
	}()

	err := exitSignal.Exit()
	assertNil(t, err)
}

func TestSignalExitWithTimeout(t *testing.T) {
	exitSignal := exit.NewSignal("one")
	exitSignal.SetTimeout(50 * time.Millisecond)
	go func() {
		reply := <-exitSignal.Chan
		time.Sleep(100 * time.Millisecond)
		reply.Ok()
	}()

	err := exitSignal.Exit()
	assertEqual(t, exit.ErrTimeout, err)
}
