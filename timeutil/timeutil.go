package timeutil

import (
	"context"
	"time"
)

// Sleep is similar to time.Sleep, but can be canceled with the provided context.
// If ctx has deadline or timeout, Sleep will return after ctx is done.
func Sleep(ctx context.Context, d time.Duration) (canceled bool) {
	t := time.NewTimer(d)
	defer t.Stop()

	select {
	case <-ctx.Done():
		canceled = true
		return
	case <-t.C:
		return
	}
}
