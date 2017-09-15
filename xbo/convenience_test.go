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

func TestConvenienceBackOffs(t *testing.T) {
	// Testing each of the convenience BackOff types
	testCases := []struct {
		bo  BackOff
		dur time.Duration
		err error
	}{
		{
			NewConstant(time.Minute),
			time.Minute,
			nil,
		},
		{
			NewZero(),
			ZeroDuration,
			nil,
		},
		{
			NewStop(),
			ZeroDuration,
			ErrStop,
		},
	}

	// We want to test that each of the convenience BackOff types
	// produce consistent output.
	iters := []int{11, 3, 7}
	for _, tc := range testCases {
		for _, iter := range iters {
			for ix := 0; ix < iter; ix++ {
				dur, err := tc.bo.Next(false)
				if dur != tc.dur {
					t.Errorf("expected %s: %s", tc.dur, dur)
				}
				if err != tc.err {
					t.Errorf("expected %v: %v", tc.err, err)
				}
			}

			// Also test that reset gets the expected standard results
			dur, err := tc.bo.Next(true)
			if dur != ZeroDuration {
				t.Errorf("expected %s: %s", ZeroDuration, dur)
			}
			if err != nil {
				t.Errorf("unexpected: %v", err)
			}
		}
	}
}
