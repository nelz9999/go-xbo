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
	"time"
)

// NewSequenceStop returns a discrete progression of durations. Once
// the number of attempts eclipses the number of elements in the
// sequence, then it returns ErrStop.
func NewSequenceStop(durs ...time.Duration) BackOff {
	size := len(durs)
	if size == 0 {
		return NewStop()
	}
	v := validator()
	return BackOffFunc(func(ctx context.Context, attempt int) (time.Duration, error) {
		z, err := v.Calc(ctx, attempt)
		if err != nil {
			return z, err
		}
		if attempt >= size {
			return ZeroDuration, ErrStop
		}
		return durs[attempt], nil
	})
}

// NewSequenceLoop returns a discrete progression of durations. Once
// the number of attempts eclipses the number of elements in the
// sequence, then continually wraps around and starts back at
// the beginning of the sequence.
func NewSequenceLoop(durs ...time.Duration) BackOff {
	size := len(durs)
	if size == 0 {
		return NewStop()
	}
	v := validator()
	return BackOffFunc(func(ctx context.Context, attempt int) (time.Duration, error) {
		z, err := v.Calc(ctx, attempt)
		if err != nil {
			return z, err
		}
		offset := attempt % size
		return durs[offset], nil
	})
}

// NewSequenceEcho returns a discrete progression of durations. Once
// the number of attempts eclipses the number of elements in the
// sequence, then it will just continually return the last
// duration in the sequence forevermore.
func NewSequenceEcho(durs ...time.Duration) BackOff {
	size := len(durs)
	if size == 0 {
		return NewStop()
	}
	v := validator()
	return BackOffFunc(func(ctx context.Context, attempt int) (time.Duration, error) {
		z, err := v.Calc(ctx, attempt)
		if err != nil {
			return z, err
		}
		offset := attempt
		if offset >= size {
			offset = size - 1
		}
		return durs[offset], nil
	})
}
