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
	"log"
	"os"
	"syscall"
	"time"

	"github.com/simia-tech/go-exit"
)

func main() {
	exit.SetTimeout(2 * time.Second)

	counterExitSignalChan, err := exit.NewSignalChan("counter")
	if err != nil {
		log.Fatal(err)
	}
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

		switch {
		case counter%5 == 0:
			// Don't send a return via errChan to simulate
			// an infinite running go routine. The timeout
			// should be hit in this case.
		case counter%2 == 1:
			errChan <- fmt.Errorf("exit on the odd counter %d", counter)
		default:
			errChan <- nil
		}
	}()

	if report := exit.ExitOn(syscall.SIGINT); report != nil {
		fmt.Println()
		report.WriteTo(os.Stderr)
		os.Exit(-1)
	}
	fmt.Println()
}
