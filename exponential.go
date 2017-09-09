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
	"math"
	"sync/atomic"
	"time"
)

type exponential struct {
	count  int32
	seed   float64
	factor float64
	stop   uint32
}

// NewExponentialBackOff defines a backoff that will increase the suggested
// wait time with each subsequent attempt.
func NewExponentialBackOff(initial time.Duration, increase float64, options ...ExponentialOption) (BackOff, error) {
	if initial <= 0 {
		return nil, fmt.Errorf("initial must be greater than zero: %v", initial)
	}
	if math.IsNaN(increase) || increase <= 0.0 {
		return nil, fmt.Errorf("increase must be a real number greater than zero: %f", increase)
	}

	result := &exponential{
		count:  -1,
		seed:   float64(initial),
		factor: 1.0 + increase,
		stop:   0,
	}

	for _, opt := range options {
		err := opt(result)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (x *exponential) Next(reset bool) (time.Duration, error) {
	if reset {
		atomic.StoreInt32(&x.count, -1)
		return 0, nil
	}

	// Which attempt is this?
	count := atomic.AddInt32(&x.count, 1)

	// User may want us to stop after too many attempts
	if x.stop > 0 && count >= int32(x.stop) {
		return 0, ErrStop
	}

	// seed * (factor**count)
	// TODO: Check for, and error, if we are overflowing int64
	multiplier := math.Pow(x.factor, float64(count))
	result := time.Duration(x.seed * multiplier)

	return result, nil
}

// ExponentialOption declares the functional options for changing behavior
type ExponentialOption func(*exponential) error

// ExponentialStop says the backoff should signal ErrStop after a sequence of
// the given number of (un-reset) attempts.
func ExponentialStop(attempts uint32) ExponentialOption {
	return ExponentialOption(func(x *exponential) error {
		x.stop = attempts
		return nil
	})
}
