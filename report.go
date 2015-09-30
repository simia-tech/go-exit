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
	"fmt"
	"io"
	"sort"
	"strings"
	"sync"
)

// Report defines the report of the exit process. It contains map of all
// signal names with thier returned error.
type Report struct {
	errors map[string]error
	mutex  *sync.RWMutex
}

// NewReport returns a new initialized Report.
func NewReport() *Report {
	return &Report{
		errors: make(map[string]error),
		mutex:  &sync.RWMutex{},
	}
}

// Set the provided error for the provided name.
func (r *Report) Set(name string, err error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.errors[name] = err
}

// Get return the error for the provided name.
func (r *Report) Get(name string) error {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return r.errors[name]
}

// Len returns the number of errors in the report.
func (r *Report) Len() int {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return len(r.errors)
}

// WriteTo prints the report to the provided io.Writer.
func (r *Report) WriteTo(w io.Writer) (int64, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	total := int64(0)
	for _, name := range r.sortedNames() {
		n, err := fmt.Fprintf(w, "%s: %v\n", name, r.errors[name])
		if err != nil {
			return total, err
		}
		total += int64(n)
	}
	return total, nil
}

// Error implements the error interface.
func (r *Report) Error() string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var parts []string
	for _, name := range r.sortedNames() {
		parts = append(parts, fmt.Sprintf("%s: %v", name, r.errors[name]))
	}
	return strings.Join(parts, " / ")
}

func (r *Report) sortedNames() []string {
	names := []string{}
	for name := range r.errors {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
