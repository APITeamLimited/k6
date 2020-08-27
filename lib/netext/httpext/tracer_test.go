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

package httpext

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/http/httptest"
	"net/http/httptrace"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/mccutchen/go-httpbin/httpbin"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/metrics"
	"github.com/loadimpact/k6/lib/netext"
	"github.com/loadimpact/k6/stats"
)

func TestTracer(t *testing.T) ***REMOVED***
	if runtime.GOOS == "windows" ***REMOVED***
		t.Skip()
	***REMOVED***
	t.Parallel()
	srv := httptest.NewTLSServer(httpbin.New().Handler())
	defer srv.Close()

	transport, ok := srv.Client().Transport.(*http.Transport)
	assert.True(t, ok)
	transport.DialContext = netext.NewDialer(
		net.Dialer***REMOVED******REMOVED***,
		netext.NewResolver(net.LookupIP, 0, lib.DNSFirst),
	).DialContext

	var prev int64
	assertLaterOrZero := func(t *testing.T, val int64, canBeZero bool) ***REMOVED***
		if canBeZero && val == 0 ***REMOVED***
			return
		***REMOVED***
		if prev > val ***REMOVED***
			_, file, line, _ := runtime.Caller(1)
			t.Errorf("Expected %d to be greater or equal to %d (from %s:%d)", val, prev, file, line)
			return
		***REMOVED***
		prev = val
	***REMOVED***

	for tnum, isReuse := range []bool***REMOVED***false, true, true***REMOVED*** ***REMOVED***
		t.Run(fmt.Sprintf("Test #%d", tnum), func(t *testing.T) ***REMOVED***
			// Do not enable parallel testing, test relies on sequential execution
			tracer := &Tracer***REMOVED******REMOVED***
			req, err := http.NewRequest("GET", srv.URL+"/get", nil)
			require.NoError(t, err)

			res, err := transport.RoundTrip(req.WithContext(httptrace.WithClientTrace(context.Background(), tracer.Trace())))
			require.NoError(t, err)

			_, err = io.Copy(ioutil.Discard, res.Body)
			assert.NoError(t, err)
			assert.NoError(t, res.Body.Close())
			trail := tracer.Done()
			trail.SaveSamples(stats.IntoSampleTags(&map[string]string***REMOVED***"tag": "value"***REMOVED***))
			samples := trail.GetSamples()

			assertLaterOrZero(t, tracer.getConn, isReuse)
			assertLaterOrZero(t, tracer.connectStart, isReuse)
			assertLaterOrZero(t, tracer.connectDone, isReuse)
			assertLaterOrZero(t, tracer.tlsHandshakeStart, isReuse)
			assertLaterOrZero(t, tracer.tlsHandshakeDone, isReuse)
			assertLaterOrZero(t, tracer.gotConn, false)
			assertLaterOrZero(t, tracer.wroteRequest, false)
			assertLaterOrZero(t, tracer.gotFirstResponseByte, false)
			assertLaterOrZero(t, now(), false)

			assert.Equal(t, strings.TrimPrefix(srv.URL, "https://"), trail.ConnRemoteAddr.String())

			assert.Len(t, samples, 8)
			seenMetrics := map[*stats.Metric]bool***REMOVED******REMOVED***
			for i, s := range samples ***REMOVED***
				assert.NotContains(t, seenMetrics, s.Metric)
				seenMetrics[s.Metric] = true

				assert.False(t, s.Time.IsZero())
				assert.Equal(t, map[string]string***REMOVED***"tag": "value"***REMOVED***, s.Tags.CloneTags())

				switch s.Metric ***REMOVED***
				case metrics.HTTPReqs:
					assert.Equal(t, 1.0, s.Value)
					assert.Equal(t, 0, i, "`HTTPReqs` is reported before the other HTTP metrics")
				case metrics.HTTPReqConnecting, metrics.HTTPReqTLSHandshaking:
					if isReuse ***REMOVED***
						assert.Equal(t, 0.0, s.Value)
						break
					***REMOVED***
					fallthrough
				case metrics.HTTPReqDuration, metrics.HTTPReqBlocked, metrics.HTTPReqSending, metrics.HTTPReqWaiting, metrics.HTTPReqReceiving:
					assert.True(t, s.Value > 0.0, "%s is <= 0", s.Metric.Name)
				default:
					t.Errorf("unexpected metric: %s", s.Metric.Name)
				***REMOVED***
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

type failingConn struct ***REMOVED***
	net.Conn
***REMOVED***

var failOnConnWrite = false

func (c failingConn) Write(b []byte) (int, error) ***REMOVED***
	if failOnConnWrite ***REMOVED***
		failOnConnWrite = false
		return 0, errors.New("write error")
	***REMOVED***

	return c.Conn.Write(b)
***REMOVED***

func TestTracerNegativeHttpSendingValues(t *testing.T) ***REMOVED***
	if runtime.GOOS == "windows" ***REMOVED***
		t.Skip()
	***REMOVED***
	t.Parallel()
	srv := httptest.NewTLSServer(httpbin.New().Handler())
	defer srv.Close()

	transport, ok := srv.Client().Transport.(*http.Transport)
	assert.True(t, ok)

	dialer := &net.Dialer***REMOVED******REMOVED***
	transport.DialContext = func(ctx context.Context, proto, addr string) (net.Conn, error) ***REMOVED***
		conn, err := dialer.DialContext(ctx, proto, addr)
		return failingConn***REMOVED***conn***REMOVED***, err
	***REMOVED***

	req, err := http.NewRequest("GET", srv.URL+"/get", nil)
	require.NoError(t, err)

	***REMOVED***
		tracer := &Tracer***REMOVED******REMOVED***
		res, err := transport.RoundTrip(req.WithContext(httptrace.WithClientTrace(context.Background(), tracer.Trace())))
		require.NoError(t, err)
		_, err = io.Copy(ioutil.Discard, res.Body)
		assert.NoError(t, err)
		assert.NoError(t, res.Body.Close())
		tracer.Done()
	***REMOVED***

	// make the next connection write fail
	failOnConnWrite = true

	***REMOVED***
		tracer := &Tracer***REMOVED******REMOVED***
		res, err := transport.RoundTrip(req.WithContext(httptrace.WithClientTrace(context.Background(), tracer.Trace())))
		require.NoError(t, err)
		_, err = io.Copy(ioutil.Discard, res.Body)
		assert.NoError(t, err)
		assert.NoError(t, res.Body.Close())
		trail := tracer.Done()
		trail.SaveSamples(nil)

		require.True(t, trail.Sending > 0)
	***REMOVED***
***REMOVED***

func TestTracerError(t *testing.T) ***REMOVED***
	t.Parallel()
	srv := httptest.NewTLSServer(httpbin.New().Handler())
	defer srv.Close()

	tracer := &Tracer***REMOVED******REMOVED***
	req, err := http.NewRequest("GET", srv.URL+"/get", nil)
	require.NoError(t, err)

	_, err = http.DefaultTransport.RoundTrip(
		req.WithContext(
			httptrace.WithClientTrace(
				context.Background(),
				tracer.Trace())))

	assert.Error(t, err)
***REMOVED***

func TestCancelledRequest(t *testing.T) ***REMOVED***
	t.Parallel()
	srv := httptest.NewTLSServer(httpbin.New().Handler())
	defer srv.Close()

	cancelTest := func(t *testing.T) ***REMOVED***
		t.Parallel()
		tracer := &Tracer***REMOVED******REMOVED***
		req, err := http.NewRequestWithContext(context.Background(), "GET", srv.URL+"/delay/1", nil)
		require.NoError(t, err)

		ctx, cancel := context.WithCancel(httptrace.WithClientTrace(req.Context(), tracer.Trace()))
		req = req.WithContext(ctx)
		go func() ***REMOVED***
			time.Sleep(time.Duration(rand.Int31n(50)) * time.Millisecond) //nolint:gosec
			cancel()
		***REMOVED***()

		resp, err := srv.Client().Transport.RoundTrip(req) //nolint:bodyclose
		_ = tracer.Done()
		if resp == nil && err == nil ***REMOVED***
			t.Errorf("Expected either a RoundTrip response or error but got %#v and %#v", resp, err)
		***REMOVED***
	***REMOVED***

	// This Run will not return until the parallel subtests complete.
	t.Run("group", func(t *testing.T) ***REMOVED***
		for i := 0; i < 200; i++ ***REMOVED***
			t.Run(fmt.Sprintf("TestCancelledRequest_%d", i), cancelTest)
		***REMOVED***
	***REMOVED***)
***REMOVED***
