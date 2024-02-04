package common

import (
	"context"
	"fmt"
	"github.com/sebasttiano/Blackbird.git/internal/logger"
	"time"
)

type RetryError struct {
	Err error
}

func (re *RetryError) Error() string {
	return fmt.Sprintf("%v", re.Err)
}

func (re *RetryError) Unwrap() error {
	return re.Err
}

func Retry[T any](ctx context.Context, retryDelays []uint, f func(context.Context, T) (T, error), arg T) (T, error) {

	var retries = len(retryDelays)
	for _, delay := range retryDelays {
		select {
		case <-ctx.Done():
			return arg, ctx.Err()
		default:
			result, err := f(ctx, arg)
			retries -= 1
			if err != nil {
				logger.Log.Error(fmt.Sprintf("Request to server failed. retrying in %d seconds... Retries left %d\n", delay, retries))
				time.Sleep(time.Duration(delay) * time.Second)
				if retries == 0 {
					return arg, &RetryError{Err: err}
				}
			} else {
				return result, nil
			}
		}
	}
	return arg, nil
}
