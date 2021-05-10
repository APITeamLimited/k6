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

package encoding

import (
	"context"
	"encoding/base64"

	"go.k6.io/k6/js/common"
)

type Encoding struct***REMOVED******REMOVED***

func New() *Encoding ***REMOVED***
	return &Encoding***REMOVED******REMOVED***
***REMOVED***

// B64encode returns the base64 encoding of input as a string.
// The data type of input can be a string, []byte or ArrayBuffer.
func (e *Encoding) B64encode(ctx context.Context, input interface***REMOVED******REMOVED***, encoding string) string ***REMOVED***
	data, err := common.ToBytes(input)
	if err != nil ***REMOVED***
		common.Throw(common.GetRuntime(ctx), err)
	***REMOVED***
	switch encoding ***REMOVED***
	case "rawstd":
		return base64.StdEncoding.WithPadding(base64.NoPadding).EncodeToString(data)
	case "std":
		return base64.StdEncoding.EncodeToString(data)
	case "rawurl":
		return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(data)
	case "url":
		return base64.URLEncoding.EncodeToString(data)
	default:
		return base64.StdEncoding.EncodeToString(data)
	***REMOVED***
***REMOVED***

// B64decode returns the decoded data of the base64 encoded input string using
// the given encoding. If format is "s" it returns the data as a string,
// otherwise as an ArrayBuffer.
func (e *Encoding) B64decode(ctx context.Context, input, encoding, format string) interface***REMOVED******REMOVED*** ***REMOVED***
	var output []byte
	var err error

	switch encoding ***REMOVED***
	case "rawstd":
		output, err = base64.StdEncoding.WithPadding(base64.NoPadding).DecodeString(input)
	case "std":
		output, err = base64.StdEncoding.DecodeString(input)
	case "rawurl":
		output, err = base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(input)
	case "url":
		output, err = base64.URLEncoding.DecodeString(input)
	default:
		output, err = base64.StdEncoding.DecodeString(input)
	***REMOVED***

	rt := common.GetRuntime(ctx) //nolint: ifshort
	if err != nil ***REMOVED***
		common.Throw(rt, err)
	***REMOVED***

	var out interface***REMOVED******REMOVED***
	if format == "s" ***REMOVED***
		out = string(output)
	***REMOVED*** else ***REMOVED***
		ab := rt.NewArrayBuffer(output)
		out = &ab
	***REMOVED***

	return out
***REMOVED***
