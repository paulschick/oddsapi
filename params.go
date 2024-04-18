// Copyright (c) Paul Schick
// SPDX-License-Identifier: MPL-2.0

package oddsapi

import "net/url"

type Params interface {
	GetEncoded() (string, error)
	BuildPath(baseUrl *url.URL) (string, error)
}
