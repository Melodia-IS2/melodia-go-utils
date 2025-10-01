package http

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Estado del circuit breaker
type CircuitBreakerState int

const (
	StateClosed CircuitBreakerState = iota
	StateOpen
	StateHalfOpen
)

type CircuitBreakerConfig struct {
	MaxFailures    int
	ResetTimeout   time.Duration
	RequestTimeout time.Duration
}

type CircuitBreaker struct {
	config      CircuitBreakerConfig
	state       CircuitBreakerState
	failures    int
	lastFailure time.Time
	mutex       sync.RWMutex
}

func NewCircuitBreaker(config CircuitBreakerConfig) *CircuitBreaker {
	return &CircuitBreaker{
		config: config,
		state:  StateClosed,
	}
}

func (cb *CircuitBreaker) CanExecute() bool {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()

	if cb.state == StateClosed {
		return true
	}

	if cb.state == StateOpen {
		return time.Since(cb.lastFailure) >= cb.config.ResetTimeout
	}

	return true
}

func (cb *CircuitBreaker) OnSuccess() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	cb.failures = 0
	cb.state = StateClosed
}

func (cb *CircuitBreaker) OnFailure() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	cb.failures++
	cb.lastFailure = time.Now()

	if cb.failures >= cb.config.MaxFailures {
		cb.state = StateOpen
	} else if cb.state == StateHalfOpen {
		cb.state = StateOpen
	}
}

func (cb *CircuitBreaker) GetState() CircuitBreakerState {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	return cb.state
}

type RobustHTTPClient struct {
	client         *http.Client
	circuitBreaker *CircuitBreaker
	retryAttempts  int
	retryDelay     time.Duration
	baseURL        string
}

type HTTPClientConfig struct {
	BaseURL       string
	Timeout       time.Duration
	MaxFailures   int
	ResetTimeout  time.Duration
	RetryAttempts int
	RetryDelay    time.Duration
}

func NewRobustHTTPClient(config HTTPClientConfig) *RobustHTTPClient {
	cbConfig := CircuitBreakerConfig{
		MaxFailures:    config.MaxFailures,
		ResetTimeout:   config.ResetTimeout,
		RequestTimeout: config.Timeout,
	}

	return &RobustHTTPClient{
		client: &http.Client{
			Timeout: config.Timeout,
			Transport: &http.Transport{
				MaxIdleConns:       10,
				IdleConnTimeout:    30 * time.Second,
				DisableCompression: true,
			},
		},
		circuitBreaker: NewCircuitBreaker(cbConfig),
		retryAttempts:  config.RetryAttempts,
		retryDelay:     config.RetryDelay,
		baseURL:        config.BaseURL,
	}
}

func (c *RobustHTTPClient) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	if !c.circuitBreaker.CanExecute() {
		state := c.circuitBreaker.GetState()
		if state == StateOpen {
			return nil, fmt.Errorf("Service is currently unavailable")
		}
	}

	var lastErr error

	for attempt := 0; attempt <= c.retryAttempts; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(c.retryDelay):
			}
		}

		resp, err := c.client.Do(req)

		if err == nil && c.isSuccessStatusCode(resp.StatusCode) {
			c.circuitBreaker.OnSuccess()
			return resp, nil
		}

		if err != nil {
			lastErr = err
		} else {
			lastErr = fmt.Errorf("HTTP %d", resp.StatusCode)
			resp.Body.Close()
		}

		if resp != nil && !c.shouldRetry(resp.StatusCode) {
			break
		}
	}

	c.circuitBreaker.OnFailure()

	return nil, fmt.Errorf("request failed after %d attempts: %w", c.retryAttempts+1, lastErr)
}

func (c *RobustHTTPClient) isSuccessStatusCode(code int) bool {
	return code >= 200 && code < 300
}

func (c *RobustHTTPClient) shouldRetry(statusCode int) bool {
	return (statusCode == http.StatusTooManyRequests || statusCode == http.StatusRequestTimeout) || statusCode >= 500
}

func (c *RobustHTTPClient) GetBaseURL() string {
	return c.baseURL
}

func (c *RobustHTTPClient) GetCircuitBreakerState() CircuitBreakerState {
	return c.circuitBreaker.GetState()
}
