package logging

import (
	"errors"
	"math"
	"math/rand"
	"time"
)

const (
	UnlimitedDeadline time.Duration = time.Duration(math.MaxInt64)
	UnlimitedDelay                  = time.Duration(math.MaxInt64)
)

var (
	ForceRetryError  error = errors.New("force to retry")
	RetryFailedError       = errors.New("retry failed")
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

type Retry interface {
	Do(func() error) error
}

type NTimesRetry struct {
	sleepFunc               func(time.Duration)
	maxTimes                uint32
	sleepTimeBetweenRetries time.Duration
}

func NewNTimesRetry(
	sleepFunc func(time.Duration),
	maxTimes uint32,
	sleepTimeBetweenRetries time.Duration) *NTimesRetry {

	return &NTimesRetry{
		sleepFunc:               sleepFunc,
		maxTimes:                maxTimes,
		sleepTimeBetweenRetries: sleepTimeBetweenRetries,
	}
}

func (self *NTimesRetry) Do(fn func() error) error {
	var err error
	for i := uint32(0); i < self.maxTimes; i++ {
		if err = fn(); err != nil {
			self.sleepFunc(self.sleepTimeBetweenRetries)
			continue
		}
		return nil
	}
	return err
}

type OnceRetry struct {
	*NTimesRetry
}

func NewOnceRetry(
	sleepFunc func(time.Duration),
	sleepTimeBetweenRetries time.Duration) *OnceRetry {

	return &OnceRetry{
		NTimesRetry: NewNTimesRetry(sleepFunc, uint32(1), sleepTimeBetweenRetries),
	}
}

type UntilElapsedRetry struct {
	sleepFunc               func(time.Duration)
	sleepTimeBetweenRetries time.Duration
	maxElapsedTime          time.Duration
}

func NewUntilElapsedRetry(
	sleepFunc func(time.Duration),
	sleepTimeBetweenRetries time.Duration,
	maxElapsedTime time.Duration) *UntilElapsedRetry {

	return &UntilElapsedRetry{
		sleepFunc:               sleepFunc,
		sleepTimeBetweenRetries: sleepTimeBetweenRetries,
		maxElapsedTime:          maxElapsedTime,
	}
}

func (self *UntilElapsedRetry) Do(fn func() error) error {
	startTime := time.Now()
	for {
		if err := fn(); err != nil {
			self.sleepFunc(self.sleepTimeBetweenRetries)
			if time.Since(startTime) >= self.maxElapsedTime {
				return err
			}
			continue
		}
		return nil
	}
}

type ExponentialBackoffRetry struct {
	sleepFunc     func(time.Duration)
	baseSleepTime time.Duration
	maxSleepTime  time.Duration
}

func NewExponentialBackoffRetry(
	sleepFunc func(time.Duration),
	baseSleepTime time.Duration,
	maxSleepTime time.Duration) *ExponentialBackoffRetry {

	return &ExponentialBackoffRetry{
		sleepFunc:     sleepFunc,
		baseSleepTime: baseSleepTime,
		maxSleepTime:  maxSleepTime,
	}
}

func (self *ExponentialBackoffRetry) Do(fn func() error) error {
	var err error
	sleepTime := self.baseSleepTime
	for {
		if err = fn(); err != nil {
			if sleepTime > self.maxSleepTime {
				sleepTime = self.maxSleepTime
			}
			self.sleepFunc(sleepTime)
			sleepTime = time.Duration(2 * int64(sleepTime))
			continue
		}
		return nil
	}
}

type BoundedExponentialBackoffRetry struct {
	sleepFunc     func(time.Duration)
	maxTries      uint32
	baseSleepTime time.Duration
	maxSleepTime  time.Duration
}

func NewBoundedExponentialBackoffRetry(
	sleepFunc func(time.Duration),
	maxTries uint32,
	baseSleepTime time.Duration,
	maxSleepTime time.Duration) *BoundedExponentialBackoffRetry {

	return &BoundedExponentialBackoffRetry{
		sleepFunc:     sleepFunc,
		maxTries:      maxTries,
		baseSleepTime: baseSleepTime,
		maxSleepTime:  maxSleepTime,
	}
}

func (self *BoundedExponentialBackoffRetry) Do(fn func() error) error {
	var err error
	sleepTime := self.baseSleepTime
	for i := uint32(0); i < self.maxTries; i++ {
		if err = fn(); err != nil {
			sleepTime = time.Duration(2 * int64(sleepTime))
			if sleepTime > self.maxSleepTime {
				sleepTime = self.maxSleepTime
			}
			self.sleepFunc(sleepTime)
			continue
		}
		return nil
	}
	return err
}

type ErrorRetry struct {
	sleepFunc   func(time.Duration)
	maxTries    int
	delay       time.Duration
	backoff     uint32
	maxJitter   float32
	maxDelay    time.Duration
	deadline    time.Duration
	retryErrors *ListSet
}

func NewErrorRetry() *ErrorRetry {
	set := NewListSet()
	set.SetAdd(ForceRetryError)
	return &ErrorRetry{
		sleepFunc:   time.Sleep,
		maxTries:    -1,
		delay:       time.Millisecond * 10,
		backoff:     2,
		maxJitter:   0.1,
		maxDelay:    UnlimitedDelay,
		deadline:    UnlimitedDeadline,
		retryErrors: set,
	}
}

func (self *ErrorRetry) SleepFunc(fn func(time.Duration)) *ErrorRetry {
	self.sleepFunc = fn
	return self
}

func (self *ErrorRetry) MaxTries(maxTries int) *ErrorRetry {
	self.maxTries = maxTries
	return self
}

func (self *ErrorRetry) Delay(delay time.Duration) *ErrorRetry {
	self.delay = delay
	return self
}

func (self *ErrorRetry) Backoff(backoff uint32) *ErrorRetry {
	self.backoff = backoff
	return self
}

func (self *ErrorRetry) MaxJitter(maxJitter float32) *ErrorRetry {
	self.maxJitter = maxJitter
	return self
}

func (self *ErrorRetry) MaxDelay(maxDelay time.Duration) *ErrorRetry {
	self.maxDelay = maxDelay
	return self
}

func (self *ErrorRetry) Deadline(deadline time.Duration) *ErrorRetry {
	self.deadline = deadline
	return self
}

func (self *ErrorRetry) OnError(err error) *ErrorRetry {
	self.retryErrors.SetAdd(err)
	return self
}

func (self *ErrorRetry) Copy() *ErrorRetry {
	return &ErrorRetry{
		sleepFunc:   self.sleepFunc,
		maxTries:    self.maxTries,
		delay:       self.delay,
		backoff:     self.backoff,
		maxJitter:   self.maxJitter,
		maxDelay:    self.maxDelay,
		deadline:    self.deadline,
		retryErrors: self.retryErrors.SetClone(),
	}
}

func RandIntN(n int) int {
	return rand.Intn(n)
}

func (self *ErrorRetry) jitterDelay(delay time.Duration) time.Duration {
	jitter := float64(RandIntN(int(self.maxJitter*100))) / 100
	return time.Duration(int64(float64(int64(delay)) * (1 + jitter)))
}

func (self *ErrorRetry) backoffDelay(delay time.Duration) time.Duration {
	return time.Duration(
		Min(int64(delay)*int64(self.backoff), int64(self.maxDelay)))
}

func Min(a, b int64) int64 {
	if a <= b {
		return a
	}
	return b
}

func (self *ErrorRetry) Do(fn func() error) error {
	latestDelay := self.delay
	startTime := time.Now()
	var err error
	for attempt := 0; attempt != self.maxTries; attempt++ {
		if err = fn(); err != nil {
			if self.retryErrors.SetContains(err) {
				sleepTime := self.jitterDelay(latestDelay)
				if self.deadline != UnlimitedDeadline {
					if (time.Since(startTime) + sleepTime) >= self.deadline {
						return RetryFailedError
					}
				}
				self.sleepFunc(sleepTime)
				latestDelay = self.backoffDelay(latestDelay)
				continue
			}
			return err
		}
		return nil
	}
	return RetryFailedError
}
