/*
Copyright 2023 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package login

import (
	"encoding/json"
	"net/http"
	"net/url"
)

func orgCheck(name string, issuer string) (verified bool, err error) {
	var endpoint string
	{
		u, err := url.Parse(issuer)
		if err != nil {
			return false, err
		}
		u.Path = `/orgcheck`

		query := make(url.Values, 1)
		query.Add("name", name)
		u.RawQuery = query.Encode()

		endpoint = u.String()
	}

	//#nosec G107 HTTP request with variable input not a risk here. The result
	// is from a public endpoint and only used as a convenience to fail early on
	// misconfigured organizations
	resp, err := http.Get(endpoint)
	if err != nil {
		return false, err
	}

	err = json.NewDecoder(resp.Body).Decode(&verified)
	if err != nil {
		return false, err
	}
	return verified, nil
}
