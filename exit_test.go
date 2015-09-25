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
	"testing"
	"time"

	"github.com/simia-tech/go-exit"
)

func TestTwoSignalChansWithSameName(t *testing.T) {
	defer exit.Reset()

	_, err := exit.NewSignalChan("one")
	assertNil(t, err)
	_, err = exit.NewSignalChan("one")
	assertEqual(t, exit.ErrNameAlreadyExists, err)
}

func TestExitWithoutError(t *testing.T) {
	exitSignalChan, err := exit.NewSignalChan("one")
	assertNil(t, err)
	go func() {
		errChan := <-exitSignalChan
		errChan <- nil
	}()

	report := exit.Exit()
	assertNil(t, report)
}

func TestExitOfTwoGoroutines(t *testing.T) {
	exitSignalChanOne, err := exit.NewSignalChan("one")
	assertNil(t, err)
	go func() {
		errChan := <-exitSignalChanOne
		errChan <- fmt.Errorf("err one")
	}()

	exitSignalChanTwo, err := exit.NewSignalChan("two")
	assertNil(t, err)
	go func() {
		errChan := <-exitSignalChanTwo
		errChan <- fmt.Errorf("err two")
	}()

	report := exit.Exit()
	assertEqual(t, 2, report.Len())
	assertEqual(t, "err one", report.Get("one").Error())
	assertEqual(t, "err two", report.Get("two").Error())
}

func TestExitWithTimeout(t *testing.T) {
	exit.SetTimeout(100 * time.Millisecond)
	defer exit.SetTimeout(0)

	exitSignalChan, err := exit.NewSignalChan("one")
	assertNil(t, err)
	go func() {
		<-exitSignalChan
	}()
	exit.NewSignalChan("two")

	report := exit.Exit()
	assertEqual(t, 2, report.Len())
	assertEqual(t, exit.ErrTimeout, report.Get("one"))
	assertEqual(t, exit.ErrTimeout, report.Get("two"))
}
