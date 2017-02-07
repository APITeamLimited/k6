/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2016 Load Impact
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

package js

import (
	"context"
	"github.com/loadimpact/k6/lib"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRunner(t *testing.T) ***REMOVED***
	if testing.Short() ***REMOVED***
		return
	***REMOVED***

	rt, err := New()
	assert.NoError(t, err)
	srcdata := &lib.SourceData***REMOVED***
		Filename: "test.js",
		Data:     []byte("export default function() ***REMOVED******REMOVED***"),
	***REMOVED***
	exp, err := rt.load(srcdata.Filename, srcdata.Data)
	assert.NoError(t, err)
	r, err := NewRunner(rt, exp)
	assert.NoError(t, err)
	if !assert.NotNil(t, r) ***REMOVED***
		return
	***REMOVED***

	t.Run("GetGroups", func(t *testing.T) ***REMOVED***
		g := r.GetGroups()
		assert.Len(t, g, 1)
		assert.Equal(t, r.DefaultGroup, g[0])
	***REMOVED***)

	t.Run("GetTests", func(t *testing.T) ***REMOVED***
		assert.Len(t, r.GetChecks(), 0)
	***REMOVED***)

	t.Run("VU", func(t *testing.T) ***REMOVED***
		vu_, err := r.NewVU()
		assert.NoError(t, err)
		vu := vu_.(*VU)

		t.Run("Reconfigure", func(t *testing.T) ***REMOVED***
			assert.NoError(t, vu.Reconfigure(12345))
			assert.Equal(t, int64(12345), vu.ID)
		***REMOVED***)

		t.Run("RunOnce", func(t *testing.T) ***REMOVED***
			_, err := vu.RunOnce(context.Background())
			assert.NoError(t, err)
		***REMOVED***)
	***REMOVED***)
***REMOVED***
