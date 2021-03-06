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

import "time"

// NewConstant creates a BackOff that will always return the configured
// duration (except during resets, where the convention is to return
// ZeroDuration).
//
// This is useful for testing.
func NewConstant(d time.Duration) BackOff {
	return BackOffFunc(func(reset bool) (time.Duration, error) {
		if reset {
			return ZeroDuration, nil
		}
		return d, nil
	})
}

// NewZero creates a BackOff that will always return 0 durations.
//
// This is useful for testing.
func NewZero() BackOff {
	return NewConstant(0)
}

// NewStop creates a BackOff that will return an ErrStop for all
// non-reset cases.
//
// This is useful for testing.
func NewStop() BackOff {
	return BackOffFunc(func(reset bool) (time.Duration, error) {
		if reset {
			return ZeroDuration, nil
		}
		return ZeroDuration, ErrStop
	})
}
