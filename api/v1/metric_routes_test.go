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

package v1

import (
	"encoding/json"
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/stats"
	"github.com/manyminds/api2go/jsonapi"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v3"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetMetrics(t *testing.T) ***REMOVED***
	engine, err := lib.NewEngine(nil, lib.Options***REMOVED******REMOVED***)
	assert.NoError(t, err)

	engine.Metrics = map[*stats.Metric]stats.Sink***REMOVED***
		&stats.Metric***REMOVED***
			Name:     "my_metric",
			Type:     stats.Trend,
			Contains: stats.Time,
			Tainted:  null.BoolFrom(true),
		***REMOVED***: &stats.TrendSink***REMOVED******REMOVED***,
	***REMOVED***

	rw := httptest.NewRecorder()
	NewHandler().ServeHTTP(rw, newRequestWithEngine(engine, "GET", "/v1/metrics", nil))
	res := rw.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	t.Run("document", func(t *testing.T) ***REMOVED***
		var doc jsonapi.Document
		assert.NoError(t, json.Unmarshal(rw.Body.Bytes(), &doc))
		if !assert.NotNil(t, doc.Data.DataArray) ***REMOVED***
			return
		***REMOVED***
		assert.Equal(t, "metrics", doc.Data.DataArray[0].Type)
	***REMOVED***)

	t.Run("metrics", func(t *testing.T) ***REMOVED***
		var metrics []Metric
		assert.NoError(t, jsonapi.Unmarshal(rw.Body.Bytes(), &metrics))
		if !assert.Len(t, metrics, 1) ***REMOVED***
			return
		***REMOVED***
		assert.Equal(t, "my_metric", metrics[0].Name)
		assert.True(t, metrics[0].Type.Valid)
		assert.Equal(t, stats.Trend, metrics[0].Type.Type)
		assert.True(t, metrics[0].Contains.Valid)
		assert.Equal(t, stats.Time, metrics[0].Contains.Type)
		assert.True(t, metrics[0].Tainted.Valid)
		assert.True(t, metrics[0].Tainted.Bool)
	***REMOVED***)
***REMOVED***

func TestGetMetric(t *testing.T) ***REMOVED***
	engine, err := lib.NewEngine(nil, lib.Options***REMOVED******REMOVED***)
	assert.NoError(t, err)

	engine.Metrics = map[*stats.Metric]stats.Sink***REMOVED***
		&stats.Metric***REMOVED***
			Name:     "my_metric",
			Type:     stats.Trend,
			Contains: stats.Time,
			Tainted:  null.BoolFrom(true),
		***REMOVED***: &stats.TrendSink***REMOVED******REMOVED***,
	***REMOVED***

	t.Run("nonexistent", func(t *testing.T) ***REMOVED***
		rw := httptest.NewRecorder()
		NewHandler().ServeHTTP(rw, newRequestWithEngine(engine, "GET", "/v1/metrics/notreal", nil))
		res := rw.Result()
		assert.Equal(t, http.StatusNotFound, res.StatusCode)
	***REMOVED***)

	t.Run("real", func(t *testing.T) ***REMOVED***
		rw := httptest.NewRecorder()
		NewHandler().ServeHTTP(rw, newRequestWithEngine(engine, "GET", "/v1/metrics/my_metric", nil))
		res := rw.Result()
		assert.Equal(t, http.StatusOK, res.StatusCode)

		t.Run("document", func(t *testing.T) ***REMOVED***
			var doc jsonapi.Document
			assert.NoError(t, json.Unmarshal(rw.Body.Bytes(), &doc))
			if !assert.NotNil(t, doc.Data.DataObject) ***REMOVED***
				return
			***REMOVED***
			assert.Equal(t, "metrics", doc.Data.DataObject.Type)
		***REMOVED***)

		t.Run("metric", func(t *testing.T) ***REMOVED***
			var metric Metric
			assert.NoError(t, jsonapi.Unmarshal(rw.Body.Bytes(), &metric))
			assert.Equal(t, "my_metric", metric.Name)
			assert.True(t, metric.Type.Valid)
			assert.Equal(t, stats.Trend, metric.Type.Type)
			assert.True(t, metric.Contains.Valid)
			assert.Equal(t, stats.Time, metric.Contains.Type)
			assert.True(t, metric.Tainted.Valid)
			assert.True(t, metric.Tainted.Bool)
		***REMOVED***)
	***REMOVED***)
***REMOVED***
