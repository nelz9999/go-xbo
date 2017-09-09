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
	crand "crypto/rand"
	"encoding/binary"
	"fmt"
	mrand "math/rand"
	"time"
)

// JitterRand describes the one function that we need from the
// math/rand.Rand type.
type JitterRand interface {
	Int63n(n int64) int64
}

type jitter struct {
	bo    BackOff
	r     JitterRand
	under uint8
	over  uint8
}

// NewJitter creates a decoration around an underlying BackOff, which adds
// a bit of randomness to the durations returned. This helps systems avoid
// "thundering herd" problems that can happen when there is accidental
// synchronization
func NewJitter(bo BackOff, options ...JitterOption) (BackOff, error) {
	if bo == nil {
		return nil, fmt.Errorf("backoff source is required")
	}

	result := &jitter{bo: bo}
	for _, opt := range options {
		err := opt(result)
		if err != nil {
			return nil, err
		}
	}

	if result.under == 0 && result.over == 0 {
		return nil, fmt.Errorf("jitter over and under not defined")
	}

	// If no random source has been applied, create our own
	if result.r == nil {
		r, err := randomlySeededRand()
		if err != nil {
			return nil, err
		}
		result.r = r
	}

	return result, nil
}

// I've seen other utilities just use time.Now().UnixNano() to seed
// their random, but here we are using a randomly-generated seed,
// because the whole point of adding jitter is to reduce likelihood of
// accidental synchronization.
func randomlySeededRand() (*mrand.Rand, error) {
	b := make([]byte, 8)
	_, err := crand.Reader.Read(b)
	if err != nil {
		return nil, err
	}
	seed := int64(binary.BigEndian.Uint64(b))
	return mrand.New(mrand.NewSource(seed)), nil
}

func (j *jitter) Next(reset bool) (time.Duration, error) {
	// We don't short-circuit, we always need to know the underlying results
	dur, err := j.bo.Next(reset)

	// But we only have work to do if it's not reset, not an error,
	// and has a non-zero duration.
	if reset || err != nil || dur <= 0 {
		return dur, err
	}

	// Calculate the range of result
	min := (dur.Nanoseconds() * int64(100-j.under)) / 100
	max := (dur.Nanoseconds() * int64(100+j.over)) / 100

	// Add in a dash of randomness, et voila!
	offset := j.r.Int63n(max - min)
	return time.Duration(min + offset), nil
}

// JitterOption declares the functional options for changing behavior
type JitterOption func(*jitter) error

// JitterRandomizer gives the consumer the option of specifying the
// source of randomness for calculations. The JitterRand interface
// easily applies to the math/rand.Rand type.
func JitterRandomizer(r JitterRand) JitterOption {
	return JitterOption(func(j *jitter) error {
		if r == nil {
			return fmt.Errorf("nil randomizer")
		}
		j.r = r
		return nil
	})
}

// JitterUnder allows the consumer to define the maximum percent of reduction
// applied to the underlying duration.
func JitterUnder(percent uint8) JitterOption {
	return JitterOption(func(j *jitter) error {
		if percent > 100 {
			return fmt.Errorf("cannot jitter under 100 percent")
		}
		j.under = percent
		return nil
	})
}

// JitterOver allows the consumer to define the maximum percent of increase
// applied to the underlying duration.
func JitterOver(percent uint8) JitterOption {
	return JitterOption(func(j *jitter) error {
		if percent > 100 {
			return fmt.Errorf("cannot jitter over 100 percent")
		}
		j.over = percent
		return nil
	})
}
