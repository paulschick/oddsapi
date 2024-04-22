// Copyright (c) Paul Schick
// SPDX-License-Identifier: MPL-2.0

package oddsapi

import (
	"github.com/hashicorp/go-retryablehttp"
	"net/url"
)

type testRequestClient struct {
	ApiToken string
	BaseURL  *url.URL

	newGetRequest func(requestUrl string, headers *map[string]string) (*retryablehttp.Request, error)
	do            func(req *retryablehttp.Request, data interface{}) (*Response, error)
}

func newTestRequestClient(apiToken string) *testRequestClient {
	bURL, _ := url.Parse(DefaultBaseUrl)
	return &testRequestClient{ApiToken: apiToken, BaseURL: bURL}
}

func (t *testRequestClient) GetApiToken() string {
	return t.ApiToken
}

func (t *testRequestClient) GetBaseUrl() *url.URL {
	u := *t.BaseURL
	return &u
}

func (t *testRequestClient) SetNewGetRequest(fn func(requestUrl string, headers *map[string]string) (*retryablehttp.Request, error)) {
	t.newGetRequest = fn
}

func (t *testRequestClient) NewGetRequest(requestUrl string, headers *map[string]string) (*retryablehttp.Request, error) {
	return t.newGetRequest(requestUrl, headers)
}

func (t *testRequestClient) SetDo(fn func(req *retryablehttp.Request, data interface{}) (*Response, error)) {
	t.do = fn
}

func (t *testRequestClient) Do(req *retryablehttp.Request, data interface{}) (*Response, error) {
	return t.do(req, data)
}
