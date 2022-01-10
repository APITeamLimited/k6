/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2017 Load Impact
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package client

import (
	"context"
	"net/http"
	"net/url"

	v1 "go.k6.io/k6/api/v1"
)

// Metrics returns the current metrics summary.
func (c *Client) Metrics(ctx context.Context) (ret []v1.Metric, err error) ***REMOVED***
	var resp v1.MetricsJSONAPI

	err = c.CallAPI(ctx, http.MethodGet, &url.URL***REMOVED***Path: "/v1/metrics"***REMOVED***, nil, &resp)
	if err != nil ***REMOVED***
		return ret, err
	***REMOVED***

	return resp.Metrics(), nil
***REMOVED***
