package timeutil

import (
	"context"
	"strconv"
	"testing"
	"time"
)

func TestSleep(t *testing.T) {
	cases := [...]struct {
		Sleep            time.Duration
		Timeout          time.Duration
		ShouldBeCanceled bool
	}{
		{
			Sleep:            10 * time.Millisecond,
			Timeout:          0,
			ShouldBeCanceled: false,
		},
		{
			Sleep:            20 * time.Millisecond,
			Timeout:          10 * time.Millisecond,
			ShouldBeCanceled: true,
		},
		{
			Sleep:            10 * time.Millisecond,
			Timeout:          20 * time.Millisecond,
			ShouldBeCanceled: false,
		},
	}

	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			start := time.Now()

			var ctx context.Context
			if c.Timeout == 0 {
				ctx = context.Background()
			} else {
				c, cancel := context.WithTimeout(context.Background(), c.Timeout)
				defer cancel()
				ctx = c
			}

			canceled := Sleep(ctx, c.Sleep)

			slept := time.Since(start)

			if c.ShouldBeCanceled && slept < c.Timeout {
				t.Errorf("cancel: slept = %v, expected = %v", slept, c.Sleep)
			}

			if !c.ShouldBeCanceled && slept < c.Sleep {
				t.Errorf("sleep: slept = %v, expected = %v", slept, c.Sleep)
			}

			if canceled != c.ShouldBeCanceled {
				t.Errorf("canceled = %v, expected = %v", canceled, c.ShouldBeCanceled)
			}
		})
	}
}
