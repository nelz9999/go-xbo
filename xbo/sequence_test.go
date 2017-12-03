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
	"testing"
	"time"
)

func TestSequenceBackOffs(t *testing.T) {
	// Testing each of the convenience BackOff types
	testCases := []struct {
		iterations int
		bos        []BackOff
	}{
		{
			5,
			[]BackOff{
				NewConstant(time.Second),
				NewLoop([]time.Duration{time.Second}, true),
				NewLoop([]time.Duration{time.Second}, false),
				NewEcho([]time.Duration{time.Second}, true),
				NewEcho([]time.Duration{time.Second}, false),
				NewLimit([]time.Duration{
					time.Second, time.Second, time.Second, time.Second, time.Second,
				}, true),
				NewLimit([]time.Duration{
					time.Second, time.Second, time.Second, time.Second, time.Second,
				}, false),
			},
		},
		{
			6,
			[]BackOff{
				NewLoop([]time.Duration{
					time.Millisecond, time.Second, time.Minute,
				}, true),
				NewLoop([]time.Duration{
					time.Millisecond, time.Second, time.Minute,
				}, false),
				NewEcho([]time.Duration{
					time.Millisecond, time.Second, time.Minute,
					time.Millisecond, time.Second, time.Minute,
				}, true),
				NewEcho([]time.Duration{
					time.Millisecond, time.Second, time.Minute,
					time.Millisecond, time.Second, time.Minute,
				}, false),
				NewLimit([]time.Duration{
					time.Millisecond, time.Second, time.Minute,
					time.Millisecond, time.Second, time.Minute,
				}, true),
				NewLimit([]time.Duration{
					time.Millisecond, time.Second, time.Minute,
					time.Millisecond, time.Second, time.Minute,
				}, false),
			},
		},
		{
			10,
			[]BackOff{
				MaxAttempts(NewConstant(time.Minute), 4, true),
				NewLimit([]time.Duration{
					time.Minute, time.Minute, time.Minute, time.Minute,
				}, true),
				NewLimit([]time.Duration{
					time.Minute, time.Minute, time.Minute, time.Minute,
				}, false),
			},
		},
	}

	// Each testCase should hold a set of BackOffs that will return
	// equivalent results, for the duration of specified iterations.
	for _, tc := range testCases {
		for ct := 0; ct < tc.iterations; ct++ {
			var xDur time.Duration
			var xErr error
			for ix, bo := range tc.bos {
				// We'll use the first BackOff as the standard for this iteration
				if ix == 0 {
					xDur, xErr = bo.Next(false)
					continue
				}
				dur, err := bo.Next(false)
				if dur != xDur {
					t.Errorf("expected %s: %s", xDur, dur)
				}
				if err != xErr {
					t.Errorf("expected %v: %v", xErr, err)
				}
			}
		}

		for _, bo := range tc.bos {
			// Also test that reset gets the expected standard results
			dur, err := bo.Next(true)
			if dur != ZeroDuration {
				t.Errorf("expected %s: %s", ZeroDuration, dur)
			}
			if err != nil {
				t.Errorf("unexpected: %v", err)
			}
		}
	}
}
