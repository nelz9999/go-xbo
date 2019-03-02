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
	"testing"
	"time"
)

func TestSequences(t *testing.T) {
	// Testing each of the convenience BackOff types
	testCases := []struct {
		name    string
		bo      BackOff
		expects map[int]time.Duration
	}{
		{
			"stop",
			NewSequenceStop(time.Millisecond, time.Second, time.Minute),
			map[int]time.Duration{
				0:   time.Millisecond,
				1:   time.Second,
				2:   time.Minute,
				3:   ZeroDuration,
				4:   ZeroDuration,
				11:  ZeroDuration,
				17:  ZeroDuration,
				229: ZeroDuration,
			},
		},
		{
			"loop",
			NewSequenceLoop(time.Millisecond, time.Second, time.Minute),
			map[int]time.Duration{
				0:   time.Millisecond,
				1:   time.Second,
				2:   time.Minute,
				3:   time.Millisecond,
				4:   time.Second,
				5:   time.Minute,
				11:  time.Minute,
				19:  time.Second,
				24:  time.Millisecond,
				229: time.Second,
			},
		},
		{
			"echo",
			NewSequenceEcho(time.Millisecond, time.Second, time.Minute),
			map[int]time.Duration{
				0:   time.Millisecond,
				1:   time.Second,
				2:   time.Minute,
				3:   time.Minute,
				4:   time.Minute,
				5:   time.Minute,
				11:  time.Minute,
				19:  time.Minute,
				229: time.Minute,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for attempt, exp := range tc.expects {
				d, err := tc.bo.Calc(context.Background(), attempt)
				if err != nil {
					// we do expect ErrStop when maxing out Stop
					if exp != 0 {
						t.Errorf("%d - unexpected: %v\n", attempt, err)
					}
				}
				if d != exp {
					t.Errorf("%d - expected %q; got %q\n", attempt, exp, d)
				}
			}
		})
	}
}
