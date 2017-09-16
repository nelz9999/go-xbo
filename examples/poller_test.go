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

package examples

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/nelz9999/go-xbo/xbo"
)

func TestExamplePoller(t *testing.T) {
	fib := []int{1, 1, 2, 3, 5, 8, 13}
	index := 0
	counter := 0
	start := time.Now()
	underlyingError := fmt.Errorf("forced error")

	f := func(ctx context.Context) error {
		var err error
		fail := "Y"
		counter++
		if counter > fib[index] {
			fail = "N"
			counter = 0
			index++
		} else {
			err = underlyingError
		}
		now := time.Now()
		diff := now.Sub(start)
		out := diff.String()
		if diff < time.Millisecond {
			out = "n/a"
		}

		// This should show us that we are resetting
		// once per every Fibonnacci number of failures
		t.Logf("Fail: %s; Elapsed: %s\n", fail, out)

		start = now
		return err
	}

	// We'll wait 50+ milliseconds after each error from `f`
	bo, err := xbo.NewExponential(time.Millisecond*50, 0.5)
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}

	// But stop if there are more than 10 errors in a row
	bo = xbo.MaxAttempts(bo, 10, false)

	// This convenience method will make sure to listen to a
	// Context for cancellation as well.
	w, err := xbo.NewWaiter(bo)
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}

	// create the example type
	p := &poller{f: f, w: w}

	// Continually do the work until something stops it
	perr := p.poll(context.Background())
	if perr != underlyingError {
		t.Errorf("expected setinel error of [%v]: %v", underlyingError, perr)
	}

}

// poller is an example of a type that repeatedly talks to
// *something* (in this case a function), which we may want
// to give some "breathing room" (back off) to when it errors out.
type poller struct {
	f func(context.Context) error
	w *xbo.Waiter
}

// poll repeated interacts with some other resource (f).
func (p *poller) poll(ctx context.Context) error {
	err := p.w.Wait(ctx, true)
	if err != nil {
		// Usually won't happen, unless a major config error
		return err
	}

	for {
		// Doing this function repeatedly is the point of the Poller
		err = p.f(ctx)

		// But if an error happens, we may want to wait a while before
		// we retry it. (We only signal to reset the BackOff when we've had
		// a successful interaction with the polled function)
		werr := p.w.Wait(ctx, (err == nil))
		if werr != nil {
			if werr == xbo.ErrStop {
				// If there's been just too many errors in a row,
				// we might want to choose to surface the error
				// that was returned from the resource being polled.
				return err
			}
			return werr
		}
	}
}
