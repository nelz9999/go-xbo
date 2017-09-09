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
	"math"
	"testing"
	"time"
)

func TestExponentialHappyPath(t *testing.T) {
	count := uint32(4)
	expected := []time.Duration{1, 2, 4, 8}
	x, err := NewExponentialBackOff(
		100*time.Millisecond,
		1.0,
		ExponentialStop(count),
	)
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}

	for ix := 0; ix < 2; ix++ {
		for _, expect := range expected {
			dur, xerr := x.Next(false)
			if xerr != nil {
				t.Errorf("unexpected: %v", xerr)
			}

			std := expect * 100 * time.Millisecond
			if std != dur {
				t.Errorf("expected %s: %s", std, dur)
			}
		}
		// Make sure reset starts the whole thing over again
		if ix == 0 {
			dur, xerr := x.Next(true)
			if xerr != nil {
				t.Errorf("unexpected: %v", xerr)
			}
			if dur != 0 {
				t.Errorf("expected 0: %s", dur)
			}
		}
	}

	_, err = x.Next(false)
	if err != ErrStop {
		t.Errorf("expected [%v]: %v", ErrStop, err)
	}
}

func TestExponentialCheckInitial(t *testing.T) {
	bads := []int64{-1, 0}
	for _, bad := range bads {
		b, err := NewExponentialBackOff(time.Duration(bad), 1.0)
		if err == nil {
			t.Errorf("expected error")
		}
		if b != nil {
			t.Errorf("expected nil: %v", b)
		}
	}
}

func TestExponentialCheckIncrease(t *testing.T) {
	bads := []float64{-1.0, math.NaN()}
	for _, bad := range bads {
		b, err := NewExponentialBackOff(1, bad)
		if err == nil {
			t.Errorf("expected error")
		}
		if b != nil {
			t.Errorf("expected nil: %v", b)
		}
	}
}
