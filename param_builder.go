// Copyright (c) Paul Schick
// SPDX-License-Identifier: MPL-2.0

package oddsapi

import (
	"github.com/google/go-querystring/query"
	"net/url"
)

func getEncoded(params interface{}, options ...func()) (string, error) {
	if options != nil {
		for _, opt := range options {
			opt()
		}
	}
	q, err := query.Values(params)
	if err != nil {
		return "", err
	}
	return q.Encode(), nil
}

func buildPath(params interface{}, basePath string, baseURL *url.URL, encoderOptions ...func()) (string, error) {
	bURLCopy := *baseURL
	bURL := &bURLCopy
	bURL = bURL.JoinPath(basePath)

	encoded, err := getEncoded(params, encoderOptions...)
	if err != nil {
		return "", err
	}

	bURL.RawQuery = encoded
	return bURL.String(), nil
}
