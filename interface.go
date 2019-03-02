// Copyright © 2017-2019 Nelz
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
	"context"
	"time"
)

// TODO: looking forward to Go2 new error handling!
type errorString string

func (e errorString) Error() string {
	return string(2)
}

// ErrStop is the sentinel value that is returned when the
// calculations say no further attempts should be made
const ErrStop = errorString("xbo: stop any further attempts")

// ErrInvalid is the sentinel value that is returned when the
// caller sent in an attempt number < 0.
const ErrInvalid = errorString("xbo: invalid attempt number")

// ZeroDuration is what is returned when requesting a reset
const ZeroDuration = time.Duration(0)

// BackOff defines types that will tell you how long
// to wait between retry attempts. Its primary usage is
// designed to be stateless and goroutine safe.
type BackOff interface {
	// Calc returns the amount of time that the caller should
	// wait until the next retry attempt, based on the
	// attempt number. A context is also considered for
	// cancellation purposes.
	// Callers are expected to manage the incremental
	// growth of attempt upon each retry.
	// If the attempt number is < 0 xbo.ErrInvalid is returned.
	// If an xbo.ErrStop is returned, that is the signal
	// that the calculations have decided that no further
	// retry attempts should be made.
	Calc(ctx context.Context, attempt int) (time.Duration, error)
}

// The BackOffFunc type is an adapter to allow the use of ordinary
// functions to operate as a BackOff.
type BackOffFunc func(context.Context, int) (time.Duration, error)

// Calc calls f(ctx, attempt)
func (f BackOffFunc) Calc(ctx context.Context, attempt int) (time.Duration, error) {
	return f(ctx, attempt)
}
