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
	"fmt"
	"syscall"
	"testing"
	"time"

	"github.com/simia-tech/go-exit"
)

func TestExitWithoutError(t *testing.T) {
	e := exit.New("test")

	exitSignal := e.NewSignal("one")
	go func() {
		reply := <-exitSignal.Chan
		reply <- nil
	}()

	report := e.Exit()
	assertNil(t, report)
}

func TestExitOfTwoGoroutines(t *testing.T) {
	e := exit.New("test")

	exitSignalOne := e.NewSignal("one")
	go func() {
		reply := <-exitSignalOne.Chan
		reply <- fmt.Errorf("err one")
	}()

	exitSignalTwo := e.NewSignal("two")
	go func() {
		reply := <-exitSignalTwo.Chan
		reply <- fmt.Errorf("err two")
	}()

	report := e.Exit()
	assertEqual(t, 2, report.Len())
	assertEqual(t, "err one", report.Get("one").Error())
	assertEqual(t, "err two", report.Get("two").Error())
}

func TestExitWithTimeout(t *testing.T) {
	e := exit.New("test")
	e.SetTimeout(100 * time.Millisecond)

	exitSignalOne := e.NewSignal("one")
	go func() {
		reply := <-exitSignalOne.Chan
		reply <- nil
	}()
	exitSignalTwo := e.NewSignal("two")
	go func() {
		<-exitSignalTwo.Chan
	}()
	e.NewSignal("three")

	report := e.Exit()
	assertEqual(t, 2, report.Len())
	assertEqual(t, exit.ErrTimeout, report.Get("two"))
	assertEqual(t, exit.ErrTimeout, report.Get("three"))
}

func TestExitOnSignal(t *testing.T) {
	e := exit.New("test")

	exitSignal := e.NewSignal("one")
	go func() {
		reply := <-exitSignal.Chan
		reply <- nil
	}()

	go func() {
		syscall.Kill(syscall.Getpid(), syscall.SIGHUP)
	}()
	report := e.ExitOn(syscall.SIGHUP)
	assertNil(t, report)
}
