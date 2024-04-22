// Copyright (c) Paul Schick
// SPDX-License-Identifier: MPL-2.0

package oddsapi

import "net/url"

type Params interface {
	BuildPath(baseUrl *url.URL) (string, error)
}
