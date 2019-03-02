// Copyright © 2017-2019 Nelz
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

func TestValidation(t *testing.T) {
	testCases := []struct {
		name string
		bo   BackOff
	}{
		{"const", NewConstant(time.Minute)},
		{"zero", NewZero()},
		{"stop", NewStop()},
		{"stopSeq0", NewSequenceStop()},
		{"stopSeq1", NewSequenceStop(time.Minute)},
		{"loopSeq0", NewSequenceLoop()},
		{"loopSeq1", NewSequenceLoop(time.Minute)},
		{"echoSeq0", NewSequenceEcho()},
		{"echoSeq1", NewSequenceEcho(time.Minute)},
	}

	bg := context.Background()
	ctx, cxl := context.WithCancel(bg)
	cxl()

	var err error
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Attempt < 0 always is invalid
			_, err = tc.bo.Calc(bg, -1)
			if err != ErrInvalid {
				t.Errorf("expected %v; got %v\n", ErrInvalid, err)
			}

			// If the context is cancelled, we return that error
			_, err = tc.bo.Calc(ctx, 0)
			if err != context.Canceled {
				t.Errorf("expected %v; got %v\n", context.Canceled, err)
			}
		})
	}
}
