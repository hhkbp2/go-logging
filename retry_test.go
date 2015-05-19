package logging

import (
	"errors"
	"github.com/hhkbp2/testify/require"
	"testing"
	"time"
)

func TestNTimesRetry(t *testing.T) {
	sleepCount := 0
	sleepFunc := func(d time.Duration) {
		sleepCount++
	}
	maxTimes := 3
	sleepTime := time.Second * 1
	retry := NewNTimesRetry(sleepFunc, uint32(maxTimes), sleepTime)
	e := errors.New("test error")
	fnCount := 0
	fn := func() error {
		fnCount++
		return e
	}
	err := retry.Do(fn)
	require.Equal(t, err, e)
	require.Equal(t, fnCount, maxTimes)
	require.Equal(t, sleepCount, maxTimes)
	fnCount = 0
	fn = func() error {
		fnCount++
		return nil
	}
	sleepCount = 0
	err = retry.Do(fn)
	require.Equal(t, err, nil)
	require.Equal(t, fnCount, 1)
	require.Equal(t, sleepCount, 0)
}

func TestOnceRetry(t *testing.T) {
	sleepCount := 0
	sleepFunc := func(d time.Duration) {
		sleepCount++
	}
	sleepTime := time.Second * 1
	retry := NewOnceRetry(sleepFunc, sleepTime)
	e := errors.New("test error")
	fnCount := 0
	fn := func() error {
		fnCount++
		return e
	}
	err := retry.Do(fn)
	require.Equal(t, err, e)
	require.Equal(t, fnCount, 1)
	require.Equal(t, sleepCount, 1)
	fnCount = 0
	fn = func() error {
		fnCount++
		return nil
	}
	sleepCount = 0
	err = retry.Do(fn)
	require.Equal(t, err, nil)
	require.Equal(t, fnCount, 1)
	require.Equal(t, sleepCount, 0)
}

func TestUntilElapsedRetry(t *testing.T) {
	sleepCount := 0
	sleepFunc := func(d time.Duration) {
		sleepCount++
		time.Sleep(d)
	}
	sleepTime := time.Second * 1
	maxElapsedTime := time.Second * 3
	retry := NewUntilElapsedRetry(sleepFunc, sleepTime, maxElapsedTime)
	e := errors.New("test error")
	fnCount := 0
	fn := func() error {
		fnCount++
		return e
	}
	startTime := time.Now()
	err := retry.Do(fn)
	require.Equal(t, err, e)
	count := int(maxElapsedTime / sleepTime)
	require.Equal(t, fnCount, count)
	require.Equal(t, sleepCount, count)
	require.True(t, time.Since(startTime) > maxElapsedTime)
	fnCount = 0
	fn = func() error {
		fnCount++
		return nil
	}
	sleepCount = 0
	startTime = time.Now()
	err = retry.Do(fn)
	require.Equal(t, err, nil)
	require.Equal(t, fnCount, 1)
	require.Equal(t, sleepCount, 0)
	require.True(t, time.Since(startTime) < maxElapsedTime)
}

func TestExponentialBackoffRetry(t *testing.T) {
	sleepCount := 0
	sleepFunc := func(d time.Duration) {
		sleepCount++
		time.Sleep(d)
	}
	times := 3
	baseSleepTime := time.Second * 1
	maxSleepTime := time.Duration(int64(time.Second) * int64(times))
	retry := NewExponentialBackoffRetry(sleepFunc, baseSleepTime, maxSleepTime)
	e := errors.New("test error")
	fnCount := 0
	fn := func() error {
		fnCount++
		if fnCount > times {
			return nil
		}
		return e
	}
	err := retry.Do(fn)
	require.Equal(t, err, nil)
	require.Equal(t, fnCount, times+1)
	require.Equal(t, sleepCount, times)
	fnCount = 0
	fn = func() error {
		fnCount++
		return nil
	}
	sleepCount = 0
	err = retry.Do(fn)
	require.Equal(t, err, nil)
	require.Equal(t, fnCount, 1)
	require.Equal(t, sleepCount, 0)
}

func TestBoundedExponentialBackoffRetry(t *testing.T) {
	sleepCount := 0
	sleepFunc := func(d time.Duration) {
		sleepCount++
		time.Sleep(d)
	}
	times := 3
	maxTries := 2
	baseSleepTime := time.Second * 1
	maxSleepTime := time.Duration(int64(time.Second) * int64(times))
	retry := NewBoundedExponentialBackoffRetry(
		sleepFunc, uint32(maxTries), baseSleepTime, maxSleepTime)
	e := errors.New("test error")
	fnCount := 0
	fn := func() error {
		fnCount++
		return e
	}
	startTime := time.Now()
	err := retry.Do(fn)
	require.Equal(t, err, e)
	require.Equal(t, fnCount, maxTries)
	require.Equal(t, sleepCount, maxTries)
	require.True(t, time.Since(startTime) > (baseSleepTime*(1+2)))
	fnCount = 0
	fn = func() error {
		fnCount++
		return nil
	}
	sleepCount = 0
	err = retry.Do(fn)
	require.Equal(t, err, nil)
	require.Equal(t, fnCount, 1)
	require.Equal(t, sleepCount, 0)
}

func TestErrorRetryMaxTries(t *testing.T) {
	sleepCount := 0
	sleepFunc := func(d time.Duration) {
		sleepCount++
	}
	maxTries := 3
	retry := NewErrorRetry().MaxTries(maxTries).SleepFunc(sleepFunc)
	fnCount := 0
	fn := func() error {
		fnCount++
		return ForceRetryError
	}
	err := retry.Do(fn)
	require.Equal(t, err, RetryFailedError)
	require.Equal(t, fnCount, maxTries)
	require.Equal(t, sleepCount, maxTries)
	sleepCount = 0
	fnCount = 0
	fn = func() error {
		fnCount++
		return nil
	}
	err = retry.Do(fn)
	require.Equal(t, err, nil)
	require.Equal(t, fnCount, 1)
	require.Equal(t, sleepCount, 0)
}

func TestErrorRetryDelay(t *testing.T) {
	sleepCount := 0
	sleepFunc := func(d time.Duration) {
		sleepCount++
		time.Sleep(d)
	}
	maxTries := 2
	delay := time.Second * 2
	maxDelay := time.Second * 3
	retry := NewErrorRetry().
		MaxTries(maxTries).
		SleepFunc(sleepFunc).
		Delay(delay).
		MaxDelay(maxDelay)
	fnCount := 0
	fn := func() error {
		fnCount++
		return ForceRetryError
	}
	startTime := time.Now()
	err := retry.Do(fn)
	require.Equal(t, err, RetryFailedError)
	require.Equal(t, fnCount, 2)
	require.Equal(t, sleepCount, 2)
	require.True(t, time.Since(startTime) > (delay+maxDelay))
}

func TestErrorRetryJitter(t *testing.T) {
	maxTries := 2
	delay := time.Second * 1
	retry := NewErrorRetry().MaxTries(maxTries).Delay(delay).MaxJitter(0.5)
	fn := func() error {
		return ForceRetryError
	}
	startTime := time.Now()
	err := retry.Do(fn)
	require.Equal(t, err, RetryFailedError)
	require.True(
		t, int64(time.Since(startTime)) > (int64(delay)*int64(maxTries)))
}

func TestErrorRetryBackoff(t *testing.T) {
	maxTries := 2
	delay := time.Second * 1
	backoff := 3
	retry := NewErrorRetry().
		MaxTries(maxTries).
		Delay(delay).
		Backoff(uint32(backoff))
	fn := func() error {
		return ForceRetryError
	}
	startTime := time.Now()
	err := retry.Do(fn)
	require.Equal(t, err, RetryFailedError)
	require.True(
		t, int64(time.Since(startTime)) > (int64(delay)*int64(1+backoff)))
}

func TestErrorRetryDeadline(t *testing.T) {
	delay := time.Second * 1
	deadline := time.Second * 4
	retry := NewErrorRetry().Delay(delay).Deadline(deadline)
	fn := func() error {
		return ForceRetryError
	}
	startTime := time.Now()
	err := retry.Do(fn)
	require.Equal(t, err, RetryFailedError)
	require.True(t, time.Since(startTime) <= deadline)
}

func TestErrorRetryErrors(t *testing.T) {
	sleepCount := 0
	sleepFunc := func(d time.Duration) {
		sleepCount++
	}
	e1 := errors.New("test error 1")
	e2 := errors.New("test error 2")
	e3 := errors.New("test error 3")
	retry := NewErrorRetry().
		SleepFunc(sleepFunc).
		Delay(time.Second).
		OnError(e1).
		OnError(e2)
	triesBeforeSuccess := 3
	fnCount := 0
	fn := func() error {
		fnCount++
		if fnCount > triesBeforeSuccess {
			return e3
		}
		switch {
		case (fnCount / 2) == 0:
			return e1
		default:
			return e2
		}
	}
	err := retry.Do(fn)
	require.Equal(t, err, e3)
	require.Equal(t, sleepCount, triesBeforeSuccess)
	require.Equal(t, fnCount, triesBeforeSuccess+1)
}
