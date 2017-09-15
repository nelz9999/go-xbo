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
	"sync/atomic"
	"time"
)

// ErrLowBound is the sentinel error returned when a boundary decorator was
// configured with a value that is so low as to be non-sensical.
var ErrLowBound = fmt.Errorf("boundary condition too low")

// MaxAttempts is a BackOff decorator that will short-circuit the underlying
// BackOff if there have been too many un-reset requests in a row.
// This can be made concurrent-safe by setting the safe boolean to true
func MaxAttempts(bo BackOff, bound uint32, safe bool) BackOff {
	count := uint32(0)
	return BackOffFunc(func(reset bool) (time.Duration, error) {
		// Check for non-sensical boundary condition
		if bound < 1 {
			return ZeroDuration, ErrLowBound
		}

		// Reset is pretty easy
		if reset {
			if safe {
				atomic.StoreUint32(&count, 0)
			} else {
				count = 0
			}
			return bo.Next(reset)
		}

		// Calculate how many sequential attempts have been made
		next := count + 1
		if safe {
			next = atomic.AddUint32(&count, 1)
		} else {
			count = next
		}

		// We've maxed out the attempts, tell them to stop
		if next > bound {
			return ZeroDuration, ErrStop
		}

		// Fall back to the underlying BackOff
		return bo.Next(reset)
	})
}

// Ceiling is a BackOff decorator that limits the maximum duration the consumer
// will be told to wait.
func Ceiling(bo BackOff, bound time.Duration) BackOff {
	return BackOffFunc(func(reset bool) (time.Duration, error) {
		// Check for non-sensical boundary condition
		if bound < 1 {
			return ZeroDuration, ErrLowBound
		}

		// Find out what the underlying BackOff says
		dur, err := bo.Next(reset)

		// We only interject for non-reset, non-error conditions
		if err == nil && !reset && dur > bound {
			return bound, nil
		}

		// Otherwise we let the underlying BackOff stand
		return dur, err
	})
}

// Elapsed is a BackOff decorator that will short-circuit the underlying
// BackOff if too much time has elapsed since the last reset, and will
// return ErrStop if that is the case.
func Elapsed(bo BackOff, bound time.Duration) BackOff {
	start := time.Now()
	return BackOffFunc(func(reset bool) (time.Duration, error) {
		// Check for non-sensical boundary condition
		if bound < 1 {
			return ZeroDuration, ErrLowBound
		}

		// Restart the clock on reset
		if reset {
			start = time.Now()
			return bo.Next(reset)
		}

		// Check elapsed before delegating, for short-circuit
		if time.Now().Sub(start) > bound {
			return ZeroDuration, ErrStop
		}

		// Fall back to the underlying BackOff
		return bo.Next(reset)
	})
}
