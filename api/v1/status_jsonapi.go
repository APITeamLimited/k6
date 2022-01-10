/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2021 Load Impact
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

package v1

import (
	"go.k6.io/k6/core"
)

// StatusJSONAPI is JSON API envelop for metrics
type StatusJSONAPI struct ***REMOVED***
	Data statusData `json:"data"`
***REMOVED***

// NewStatusJSONAPI creates the JSON API status envelop
func NewStatusJSONAPI(s Status) StatusJSONAPI ***REMOVED***
	return StatusJSONAPI***REMOVED***
		Data: statusData***REMOVED***
			ID:         "default",
			Type:       "status",
			Attributes: s,
		***REMOVED***,
	***REMOVED***
***REMOVED***

// Status extract the v1.Status from the JSON API envelop
func (s StatusJSONAPI) Status() Status ***REMOVED***
	return s.Data.Attributes
***REMOVED***

type statusData struct ***REMOVED***
	Type       string `json:"type"`
	ID         string `json:"id"`
	Attributes Status `json:"attributes"`
***REMOVED***

func newStatusJSONAPIFromEngine(engine *core.Engine) StatusJSONAPI ***REMOVED***
	return NewStatusJSONAPI(NewStatus(engine))
***REMOVED***
