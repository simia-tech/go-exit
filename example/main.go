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

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/simia-tech/go-exit"
)

func main() {
	counterExitSignalChan := exit.Signal("counter")
	go func() {
		counter := 0

		var errChan exit.ErrChan
		for errChan == nil {
			select {
			case errChan = <-counterExitSignalChan:
				break
			case <-time.After(1 * time.Second):
				counter++
				fmt.Printf("%d ", counter)
			}
		}

		if counter%2 == 0 {
			errChan <- nil
		} else {
			errChan <- fmt.Errorf("exit on the odd counter %d", counter)
		}
	}()

	time.Sleep(6 * time.Second)
	fmt.Println()
	if report := exit.Exit(); report != nil {
		report.Fprint(os.Stderr)
		os.Exit(-1)
	}
}
