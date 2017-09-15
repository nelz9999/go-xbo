// Copyright © 2017 Nelz
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package xbo

import (
	"fmt"
	"time"
)

// ErrStop is the sentinel value that is returned when the
// calculations say no further attempts should be made
var ErrStop = fmt.Errorf("stop any further attempts")

// ZeroDuration is what is returned when requesting a reset
var ZeroDuration = time.Duration(0)

// BackOff defines objects that will tell you how long
// to wait between attempts
type BackOff interface {
	// Next returns the amount of time that should be
	// waited until the next attempt.
	// If an xbo.ErrStop is returned, that is the signal
	// that no further attempts should be made
	// Sending a reset value of true means you want to
	// start the sequence of values over again from the
	// beginning. When being reset, it is customary to
	// return ZeroDuration.
	Next(reset bool) (time.Duration, error)
}

// BackOffFunc is an adapter so a function can implement
// the BackOff interface
type BackOffFunc func(bool) (time.Duration, error)

// Next fulfills the BackOff interface
func (f BackOffFunc) Next(reset bool) (time.Duration, error) {
	return f(reset)
}
