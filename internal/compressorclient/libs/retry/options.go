package retry

import (
    "math/rand"
    "time"
)

// Function signature of retry if function
type RetryIfFunc func(error) bool

// Function signature of OnRetry function
// n = count of attempts
type OnRetryFunc func(n uint, err error)

type DelayTypeFunc func(n uint, config *Config) time.Duration

type Config struct {
    attempts      uint
    delay         time.Duration
    maxDelay      time.Duration
    maxJitter     time.Duration
    onRetry       OnRetryFunc
    retryIf       RetryIfFunc
    delayType     DelayTypeFunc
    lastErrorOnly bool
}

// Option represents an option for retry.
type Option func(*Config)

// return the direct last error that came from the retried function
// default is false (return wrapped errors with everything)
func LastErrorOnly(lastErrorOnly bool) Option {
    return func(c *Config) {
        c.lastErrorOnly = lastErrorOnly
    }
}

// Attempts set count of retry
// default is 10
func Attempts(attempts uint) Option {
    return func(c *Config) {
        c.attempts = attempts
    }
}

// Delay set delay between retry
// default is 100ms
func Delay(delay time.Duration) Option {
    return func(c *Config) {
        c.delay = delay
    }
}

// MaxDelay set maximum delay between retry
// does not apply by default
func MaxDelay(maxDelay time.Duration) Option {
    return func(c *Config) {
        c.maxDelay = maxDelay
    }
}

// MaxJitter sets the maximum random Jitter between retries for RandomDelay
func MaxJitter(maxJitter time.Duration) Option {
    return func(c *Config) {
        c.maxJitter = maxJitter
    }
}

// DelayType set type of the delay between retries
// default is BackOff
func DelayType(delayType DelayTypeFunc) Option {
    return func(c *Config) {
        c.delayType = delayType
    }
}

// BackOffDelay is a DelayType which increases delay between consecutive retries
func BackOffDelay(n uint, config *Config) time.Duration {
    return config.delay * (1 << n)
}

// FixedDelay is a DelayType which keeps delay the same through all iterations
func FixedDelay(_ uint, config *Config) time.Duration {
    return config.delay
}

// RandomDelay is a DelayType which picks a random delay up to config.maxJitter
func RandomDelay(_ uint, config *Config) time.Duration {
    return time.Duration(rand.Int63n(int64(config.maxJitter)))
}

// CombineDelay is a DelayType the combines all of the specified delays into a new DelayTypeFunc
func CombineDelay(delays ...DelayTypeFunc) DelayTypeFunc {
    return func(n uint, config *Config) time.Duration {
        var total time.Duration
        for _, delay := range delays {
            total += delay(n, config)
        }
        return total
    }
}

func OnRetry(onRetry OnRetryFunc) Option {
    return func(c *Config) {
        c.onRetry = onRetry
    }
}

func RetryIf(retryIf RetryIfFunc) Option {
    return func(c *Config) {
        c.retryIf = retryIf
    }
}
