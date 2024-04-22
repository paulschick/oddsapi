// Copyright (c) Paul Schick
// SPDX-License-Identifier: MPL-2.0

package oddsapi

func requestHandler[T any](params Params, client BaseRequestClient, data T) (T, *Response, error) {
	reqUrl, err := params.BuildPath(client.GetBaseUrl())
	if err != nil {
		return data, nil, err
	}

	req, err := client.NewGetRequest(reqUrl, nil)
	if err != nil {
		return data, nil, err
	}

	resp, err := client.Do(req, &data)
	if err != nil {
		return data, resp, err
	}

	return data, resp, nil
}
