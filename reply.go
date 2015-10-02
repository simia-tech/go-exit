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

// Reply defines a channel of errors that can be used to deliver back
// an error after an actor has shut down.
type Reply chan error

// Err reports back the provided error.
func (r Reply) Err(err error) {
	r <- err
}

// Ok report back no error. Same as Err(nil).
func (r Reply) Ok() {
	r <- nil
}
