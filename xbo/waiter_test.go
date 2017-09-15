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

func TestNewWaiterError(t *testing.T) {
	_, err := NewWaiter(nil)
	if err == nil {
		t.Errorf("expected error")
	}
}

func TestWaiterDeadline(t *testing.T) {
	w, err := NewWaiter(NewConstant(time.Millisecond * 100))
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}

	ctx, cxl := context.WithTimeout(
		context.Background(),
		time.Millisecond*250)
	defer cxl()

	attempts := 0
	for err == nil {
		attempts++
		err = w.Wait(ctx, false)
	}

	if err != context.DeadlineExceeded {
		t.Errorf("unexpected: %v", err)
	}

	if attempts != 3 {
		t.Errorf("expected 3 attempts: %d", attempts)
	}
}

func TestWaiterStop(t *testing.T) {
	bo := MaxAttempts(NewConstant(time.Millisecond*25), 4, false)
	w, err := NewWaiter(bo)
	if err != nil {
		t.Errorf("unexpected: %v", err)
	}

	attempts := 0
	for err == nil {
		attempts++
		err = w.Wait(context.Background(), false)
	}

	if err != ErrStop {
		t.Errorf("unexpected: %v", err)
	}

	if attempts != 5 {
		t.Errorf("expected 3 attempts: %d", attempts)
	}
}
