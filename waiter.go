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
	"context"
	"fmt"
	"time"
)

// NewWaiter produces a Waiter based off an underlying BackOff
func NewWaiter(bo BackOff) (*Waiter, error) {
	if bo == nil {
		return nil, fmt.Errorf("backoff must be defined")
	}
	return &Waiter{bo: bo}, nil
}

// Waiter is a wrapper around a BackOff that will block
// execution for the amount of time dictated by that BackOff
type Waiter struct {
	bo BackOff
}

// Wait will interrogate the underlying BackOff for the expected
// duration, and will then block for that amount of time. You can send
// in a Context for early cancellation.
func (w Waiter) Wait(ctx context.Context, reset bool) error {
	dur, err := w.bo.Next(reset)
	if err != nil {
		return err
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(dur):
		// Happy path
	}
	return nil
}
