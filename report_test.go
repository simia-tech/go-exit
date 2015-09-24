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
	"bytes"
	"testing"

	"github.com/simia-tech/go-exit"
)

func TestReportFprint(t *testing.T) {
	report := exit.NewReport()
	report.Set("one", exit.ErrTimeout)

	buffer := &bytes.Buffer{}
	report.Fprint(buffer)
	assertEqual(t, "one: timeout\n", buffer.String())
}

func TestReportErrorInterface(t *testing.T) {
	report := exit.NewReport()
	report.Set("one", exit.ErrTimeout)
	report.Set("two", exit.ErrTimeout)

	var err error = report
	assertEqual(t, "one: timeout / two: timeout", err.Error())
}
