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
	"math/rand"
	"testing"
	"time"
)

func TestBoundaryErrors(t *testing.T) {
	// Testing each of the bounded BackOff types
	offs := []BackOff{
		MaxAttempts(NewStop(), 0, false),
		MaxAttempts(NewStop(), 0, true),
		Ceiling(NewStop(), 0),
		Elapsed(NewStop(), 0),
	}
	resets := []bool{true, false}

	for ct, off := range offs {
		for _, reset := range resets {
			dur, err := off.Next(reset)
			if dur != ZeroDuration {
				t.Errorf("%d expected %s: %s", ct, ZeroDuration, dur)
			}
			if err != ErrLowBound {
				t.Errorf("%d expected %v: %v", ct, ErrLowBound, err)
			}
		}
	}
}

func TestMaxAttempts(t *testing.T) {
	max := 4 + rand.Intn(10)

	expected := time.Minute
	safes := []bool{true, false}
	for _, safe := range safes {
		bo := MaxAttempts(NewConstant(expected), uint32(max), safe)

		// Have to do it at least 2 times to make sure the
		// reset takes
		cycles := 2 + rand.Intn(3)
		for ix := 0; ix < cycles; ix++ {
			// Get us up to, but not over the max
			for jx := 0; jx < max; jx++ {
				dur, err := bo.Next(false)
				if dur != expected {
					t.Errorf("expected %v: %v", expected, dur)
				}
				if err != nil {
					t.Errorf("unexpected: %v", err)
				}
			}

			// We expect to be at the max, even for several attempts
			for jx := 0; jx < 3; jx++ {
				dur, err := bo.Next(false)
				if dur != ZeroDuration {
					t.Errorf("expected %v: %v", ZeroDuration, dur)
				}
				if err != ErrStop {
					t.Errorf("expected %v: %v", ErrStop, err)
				}
			}

			// Call for a reset, even a couple of times. This
			// sets us up for the next cycle to start fresh
			for jx := 0; jx < 3; jx++ {
				dur, err := bo.Next(true)
				if dur != ZeroDuration {
					t.Errorf("expected %v: %v", ZeroDuration, dur)
				}
				if err != nil {
					t.Errorf("unexpected: %v", err)
				}
			}
		}
	}
}

func TestCeiling(t *testing.T) {
	under, err := NewExponential(time.Millisecond, 1.0)
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}

	top := time.Millisecond * 10
	bo := Ceiling(under, top)

	// Several cycles to prove reset works
	cycles := 2 + rand.Intn(3)
	for ix := 0; ix < cycles; ix++ {

		// This should get us 1ms, 2ms, 4ms, 8ms
		for jx := 0; jx < 4; jx++ {
			dur, err := bo.Next(false)
			if dur >= top {
				t.Errorf("unexpected: %v", dur)
			}
			if err != nil {
				t.Errorf("unexpected: %v", err)
			}
		}

		// Next in the sequence should be
		// 16ms, 32ms, 64ms, but these should all be
		// governed down to the 10ms ceiling
		for jx := 0; jx < 3; jx++ {
			dur, err := bo.Next(false)
			if dur != top {
				t.Errorf("expected %v: %v", top, dur)
			}
			if err != nil {
				t.Errorf("unexpected: %v", err)
			}
		}

		// Call for a reset, even a couple of times. This
		// sets us up for the next cycle to start fresh
		for jx := 0; jx < 3; jx++ {
			dur, err := bo.Next(true)
			if dur != ZeroDuration {
				t.Errorf("expected %v: %v", ZeroDuration, dur)
			}
			if err != nil {
				t.Errorf("unexpected: %v", err)
			}
		}
	}
}

func TestElapsed(t *testing.T) {
	expected := time.Minute
	under := NewConstant(expected)

	window := time.Millisecond * 200
	bo := Elapsed(under, window)

	// Several cycles to prove reset works
	cycles := 2 + rand.Intn(3)
	for ix := 0; ix < cycles; ix++ {
		// Since we're not over time, we should just get
		// what the underlying says
		for jx := 0; jx < 10; jx++ {
			dur, err := bo.Next(false)
			if dur != expected {
				t.Errorf("expected %v: %v", expected, dur)
			}
			if err != nil {
				t.Errorf("unexpected: %v", err)
			}
		}

		// Wait until we've gone over time
		time.Sleep(window)

		// We expect to be told to stop, even for several attempts
		for jx := 0; jx < 3; jx++ {
			dur, err := bo.Next(false)
			if dur != ZeroDuration {
				t.Errorf("expected %v: %v", ZeroDuration, dur)
			}
			if err != ErrStop {
				t.Errorf("expected %v: %v", ErrStop, err)
			}
		}

		// Call for a reset, even a couple of times. This
		// sets us up for the next cycle to start fresh
		for jx := 0; jx < 3; jx++ {
			dur, err := bo.Next(true)
			if dur != ZeroDuration {
				t.Errorf("expected %v: %v", ZeroDuration, dur)
			}
			if err != nil {
				t.Errorf("unexpected: %v", err)
			}
		}
	}
}

// TODO: more tests
// TestMaxAttemptsSafe
// TestMaxAttemptsShortCircuit
// TestElapsedShortCircuit
