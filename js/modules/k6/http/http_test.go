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

package http

import (
	"testing"

	"github.com/dop251/goja"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.k6.io/k6/js/common"
	"go.k6.io/k6/lib/netext/httpext"
)

func TestTagURL(t *testing.T) ***REMOVED***
	rt := goja.New()
	rt.SetFieldNameMapper(common.FieldNameMapper***REMOVED******REMOVED***)
	rt.Set("http", common.Bind(rt, new(GlobalHTTP).NewModuleInstancePerVU(), nil))

	testdata := map[string]struct***REMOVED*** u, n string ***REMOVED******REMOVED***
		`http://localhost/anything/`:               ***REMOVED***"http://localhost/anything/", "http://localhost/anything/"***REMOVED***,
		`http://localhost/anything/$***REMOVED***1+1***REMOVED***`:         ***REMOVED***"http://localhost/anything/2", "http://localhost/anything/$***REMOVED******REMOVED***"***REMOVED***,
		`http://localhost/anything/$***REMOVED***1+1***REMOVED***/`:        ***REMOVED***"http://localhost/anything/2/", "http://localhost/anything/$***REMOVED******REMOVED***/"***REMOVED***,
		`http://localhost/anything/$***REMOVED***1+1***REMOVED***/$***REMOVED***1+2***REMOVED***`:  ***REMOVED***"http://localhost/anything/2/3", "http://localhost/anything/$***REMOVED******REMOVED***/$***REMOVED******REMOVED***"***REMOVED***,
		`http://localhost/anything/$***REMOVED***1+1***REMOVED***/$***REMOVED***1+2***REMOVED***/`: ***REMOVED***"http://localhost/anything/2/3/", "http://localhost/anything/$***REMOVED******REMOVED***/$***REMOVED******REMOVED***/"***REMOVED***,
	***REMOVED***
	for expr, data := range testdata ***REMOVED***
		expr, data := expr, data
		t.Run("expr="+expr, func(t *testing.T) ***REMOVED***
			tag, err := httpext.NewURL(data.u, data.n)
			require.NoError(t, err)
			v, err := rt.RunString("http.url`" + expr + "`")
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, tag, v.Export())
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***
