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
	"sync/atomic"
	"time"
)

// NewLoop allows the client to define a discrete list of durations that
// will be continually looped over, without end. Resets do put the iteration
// back to the first item.
//
// Misconfiguraiton (0-length slice) will create a Backoff that always
// returns ErrStop (when not being reset).
func NewLoop(durs []time.Duration, safe bool) BackOff {
	return newSequence(durs, safe, true, false)
}

// NewLimit allows the client to define a discrete list of durations that
// will be returned until the list runs out, at which time the BackOff will
// just return ErrStop, until reset.
//
// Misconfiguraiton (0-length slice) will create a Backoff that always
// returns ErrStop (when not being reset).
func NewLimit(durs []time.Duration, safe bool) BackOff {
	return newSequence(durs, safe, false, false)
}

// NewEcho allows the client to define a discrete list of durations that
// will be returned until the list runs out, at which time the last duration
// in the list will be continually returned, until a reset.
//
// Misconfiguraiton (0-length slice) will create a Backoff that always
// returns ErrStop (when not being reset).
func NewEcho(durs []time.Duration, safe bool) BackOff {
	return newSequence(durs, safe, false, true)
}

func newSequence(durs []time.Duration, safe bool, loop bool, echo bool) BackOff {
	// Since we are trying to protect some underlying resource, if the user
	// specified an empty (nonsensical) slice, then default to stopping
	// any retries
	size := uint32(len(durs))
	if size == 0 {
		return NewStop()
	}

	count := uint32(0)
	return BackOffFunc(func(reset bool) (time.Duration, error) {
		// Reset is pretty easy
		if reset {
			if safe {
				atomic.StoreUint32(&count, 0)
			} else {
				count = 0
			}
			return ZeroDuration, nil
		}

		// After calculating the current duration we will increase the count
		// for the next time around
		defer func() {
			next := count + 1
			if safe {
				atomic.AddUint32(&count, 1)
			} else {
				count = next
			}
		}()

		offset := count
		if loop {
			offset = count % size
		}

		// Short-circuit if we're at the max size
		if offset >= size {
			if echo {
				// Just echo the last entry
				return durs[size-1], nil
			}
			return ZeroDuration, ErrStop
		}

		return durs[offset], nil
	})
}
