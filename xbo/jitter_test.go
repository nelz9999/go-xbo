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

func TestNewJitterErrors(t *testing.T) {
	// This section tests errors thrown when under-specifying the inputs
	offs := []BackOff{
		nil,       // Nil underlying
		NewZero(), // No over/under specified
	}

	for _, off := range offs {
		bo, err := NewJitter(off)
		if err == nil {
			t.Errorf("expected error")
		}
		if bo != nil {
			t.Errorf("unexpected: %v", bo)
		}
	}

	// This section tests problems within the options themselves
	opts := []JitterOption{
		JitterRandomizer(nil),
		JitterUnder(101),
		JitterOver(101),
	}

	for _, opt := range opts {
		bo, err := NewJitter(NewZero(), opt)
		if err == nil {
			t.Errorf("expected error")
		}
		if bo != nil {
			t.Errorf("unexpected: %v", bo)
		}
	}
}

type randomFunc func(int64) int64

func (f randomFunc) Int63n(n int64) int64 {
	return f(n)
}

func Topper() JitterRand {
	return randomFunc(func(n int64) int64 {
		// Return range is [0,n)
		return n - 1
	})
}

func Bottom() JitterRand {
	return randomFunc(func(n int64) int64 {
		return 0
	})
}

func clean(bo BackOff, err error) BackOff {
	return bo
}

func TestJitterUnRandom(t *testing.T) {
	base := time.Second * 10
	bo := NewConstant(base)

	testCases := []struct {
		bo  BackOff
		dur time.Duration
		err error
	}{
		// All these have "randoms" that go to the bottom of the range
		{
			clean(NewJitter(bo,
				JitterUnder(10),
				JitterRandomizer(Bottom()),
			)),
			time.Second * 9,
			nil,
		},
		{
			clean(NewJitter(bo,
				JitterUnder(10),
				JitterOver(10),
				JitterRandomizer(Bottom()),
			)),
			time.Second * 9,
			nil,
		},
		{
			clean(NewJitter(bo,
				JitterOver(10),
				JitterRandomizer(Bottom()),
			)),
			time.Second * 10,
			nil,
		},
		// All these have "randoms" that go to the top of the range
		{
			clean(NewJitter(bo,
				JitterUnder(10),
				JitterRandomizer(Topper()),
			)),
			time.Second * 10,
			nil,
		},
		{
			clean(NewJitter(bo,
				JitterUnder(10),
				JitterOver(10),
				JitterRandomizer(Topper()),
			)),
			time.Second * 11,
			nil,
		},
		{
			clean(NewJitter(bo,
				JitterOver(10),
				JitterRandomizer(Topper()),
			)),
			time.Second * 11,
			nil,
		},
	}

	for _, testCase := range testCases {
		// Test a non-reset
		dur, err := testCase.bo.Next(false)
		if err != testCase.err {
			t.Errorf("expected %v: %v", testCase.err, err)
		}
		if dur != testCase.dur {
			t.Errorf("expected %v: %v", testCase.dur, dur)
		}

		// Also ensure the standard response for reset
		dur, err = testCase.bo.Next(true)
		if dur != ZeroDuration {
			t.Errorf("expected %v: %v", ZeroDuration, dur)
		}
		if err != nil {
			t.Errorf("unexpected: %v", err)
		}
	}
}

func TestJitterRandomish(t *testing.T) {
	base := NewConstant(time.Second * 10)
	now := time.Now().UnixNano()
	in1 := clean(NewJitter(base,
		JitterUnder(25),
		JitterOver(25),
		JitterRandomizer(rand.New(rand.NewSource(now))),
	))

	in2 := clean(NewJitter(base,
		JitterUnder(25),
		JitterOver(25),
		JitterRandomizer(rand.New(rand.NewSource(now))),
	))

	out := clean(NewJitter(base,
		JitterUnder(25),
		JitterOver(25),
	))

	allMatch := true
	for ix := 0; ix < 10; ix++ {
		d1, err := in1.Next(false)
		if err != nil {
			t.Errorf("unexpected: %v", err)
		}

		d2, err := in2.Next(false)
		if err != nil {
			t.Errorf("unexpected: %v", err)
		}

		do, err := out.Next(false)
		if err != nil {
			t.Errorf("unexpected: %v", err)
		}

		if d1 != d2 {
			t.Errorf("broken determinism: %s vs %s", d1, d2)
		}

		if d1 == do {
			t.Logf("unexpected match: %s vs %s", d1, do)
		} else {
			allMatch = false
		}
		// t.Logf("I: %s; O: %s\n", d1, do)
	}

	if allMatch {
		t.Errorf("Did not expect so many matches!")
	}
}
