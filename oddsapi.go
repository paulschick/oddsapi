// Copyright (c) Paul Schick
// SPDX-License-Identifier: MPL-2.0

package oddsapi

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-retryablehttp"
	"golang.org/x/time/rate"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	DefaultBaseUrl    = "https://api.the-odds-api.com/"
	userAgent         = "go-oddsapi"
	DefaultRegion     = "us"
	DefaultDateFormat = "iso"
	DefaultOddsFormat = "decimal"
)

const (
	DefaultRateLimitPercent = 0.75
	DefaultBurstPercent     = 0.25
)

type Client struct {
	client            *retryablehttp.Client
	baseUrl           *url.URL
	apiToken          string
	UserAgent         string
	RateLimit         int
	limiter           *rate.Limiter
	maxPercentOfLimit float64
	limiterBurst      float64
	configureOnce     sync.Once

	SportsService *SportsService
}

func NewClient(apiToken string, rateLimitPerSec int, options ...ClientOption) (*Client, error) {
	c := &Client{
		apiToken:          apiToken,
		UserAgent:         userAgent,
		maxPercentOfLimit: DefaultRateLimitPercent,
		limiterBurst:      DefaultBurstPercent,
		RateLimit:         rateLimitPerSec,
	}
	err := c.setBaseUrl(DefaultBaseUrl)
	if err != nil {
		return nil, err
	}
	c.client = &retryablehttp.Client{
		CheckRetry: func(ctx context.Context, resp *http.Response, err error) (bool, error) {
			if ctx.Err() != nil {
				return false, ctx.Err()
			}
			if err != nil {
				return false, err
			}
			if resp.StatusCode == 429 || resp.StatusCode >= 500 {
				return true, nil
			}
			return false, nil
		},
		Backoff: func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
			min = 1 * time.Second
			max = 2 * time.Second
			return retryablehttp.LinearJitterBackoff(min, max, attemptNum, resp)
		},
		ErrorHandler: retryablehttp.PassthroughErrorHandler,
		HTTPClient:   cleanhttp.DefaultClient(),
		RetryWaitMin: 250 * time.Millisecond,
		RetryWaitMax: 1 * time.Second,
		RetryMax:     5,
	}

	c.SportsService = NewSportsService(c)

	err = c.applyOptions(options...)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Client) applyOptions(options ...ClientOption) error {
	for _, fn := range options {
		if fn != nil {
			err := fn(c)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Client) setBaseUrl(urlStr string) error {
	if !strings.HasSuffix(urlStr, "/") {
		urlStr += "/"
	}

	baseURL, err := url.Parse(urlStr)
	if err != nil {
		return err
	}

	c.baseUrl = baseURL
	return nil
}

func (c *Client) configureRateLimiter(ctx context.Context) {
	rl := float64(c.RateLimit)
	limit := rate.Limit(rl * c.maxPercentOfLimit)
	burst := 1
	if int(rl*c.limiterBurst) > 1 {
		burst = int(rl * c.limiterBurst)
	}
	c.limiter = rate.NewLimiter(limit, burst)

	_ = c.limiter.Wait(ctx)
}

func (c *Client) Do(req *retryablehttp.Request, data interface{}) (*Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = resp.Body.Close()
		_, _ = io.Copy(io.Discard, resp.Body)
	}()

	response := newResponse(resp)

	c.configureOnce.Do(func() { c.configureRateLimiter(req.Context()) })

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bodyBytes, data)
	if err != nil {
		return nil, err
	}

	return response, err
}

func (c *Client) GetBaseUrl() *url.URL {
	u := *c.baseUrl
	return &u
}

func (c *Client) NewGetRequest(requestUrl string, headers *map[string]string) (*retryablehttp.Request, error) {
	reqHeaders := make(http.Header)
	if c.UserAgent != "" {
		reqHeaders.Set("User-Agent", c.UserAgent)
	}

	if headers != nil {
		for k, v := range *headers {
			reqHeaders.Set(k, v)
		}
	}

	req, err := retryablehttp.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range reqHeaders {
		req.Header[k] = v
	}

	return req, nil
}

type Response struct {
	*http.Response
}

func newResponse(response *http.Response) *Response {
	r := &Response{Response: response}
	return r
}
