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

package ws

import (
	"bytes"
	"context"
	"errors"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/dop251/goja"
	"github.com/gorilla/websocket"
	"github.com/loadimpact/k6/js/common"
	"github.com/loadimpact/k6/lib/metrics"
	"github.com/loadimpact/k6/stats"
)

type WS struct***REMOVED******REMOVED***

type Socket struct ***REMOVED***
	ctx           context.Context
	conn          *websocket.Conn
	eventHandlers map[string][]goja.Callable
	scheduled     chan goja.Callable
	done          chan struct***REMOVED******REMOVED***
***REMOVED***

const writeWait = 10 * time.Second

var (
	newline = []byte***REMOVED***'\n'***REMOVED***
	space   = []byte***REMOVED***' '***REMOVED***
)

func (*WS) Connect(ctx context.Context, url string, args ...goja.Value) *http.Response ***REMOVED***
	rt := common.GetRuntime(ctx)
	state := common.GetState(ctx)

	setupFn, userTags, header := parseArgs(rt, args)

	tags := map[string]string***REMOVED***
		"url":   url,
		"group": state.Group.Path,
	***REMOVED***
	// Merge with the user-provided tags
	for k, v := range userTags ***REMOVED***
		tags[k] = v
	***REMOVED***

	// Pass a custom net.Dial function to websocket.Dialer that will substitute
	// the underlying net.Conn with our own TraceConn
	var traceConn *TraceConn
	netDial := func(network, address string) (net.Conn, error) ***REMOVED***
		var d net.Dialer
		conn, err := d.Dial(network, address)
		traceConn = &TraceConn***REMOVED***conn, 0, 0***REMOVED***

		return traceConn, err
	***REMOVED***
	wsd := websocket.Dialer***REMOVED***NetDial: netDial, Proxy: http.ProxyFromEnvironment***REMOVED***

	connectionStart := time.Now()
	conn, response, connErr := wsd.Dial(url, header)
	connecting := float64(time.Since(connectionStart)) / float64(time.Millisecond)

	socket := Socket***REMOVED***
		ctx:           ctx,
		conn:          conn,
		eventHandlers: make(map[string][]goja.Callable),
		scheduled:     make(chan goja.Callable),
		done:          make(chan struct***REMOVED******REMOVED***),
	***REMOVED***

	// Run the user-provided set up function
	setupFn(goja.Undefined(), rt.ToValue(&socket))

	if connErr != nil ***REMOVED***
		// Pass the error to the user script before exiting immediately
		socket.handleEvent("error", rt.ToValue(connErr))
		return response
	***REMOVED***
	defer conn.Close()

	tags["status"] = strconv.Itoa(response.StatusCode)
	tags["subprotocol"] = response.Header.Get("Sec-WebSocket-Protocol")

	// The connection is now open, emit the event
	socket.handleEvent("open")

	// Pass ping/pong events through the main control loop
	pingPongChan := make(chan string)
	conn.SetPingHandler(func(string) error ***REMOVED*** pingPongChan <- "ping"; return nil ***REMOVED***)
	conn.SetPongHandler(func(string) error ***REMOVED*** pingPongChan <- "pong"; return nil ***REMOVED***)

	readDataChan := make(chan []byte)
	readErrChan := make(chan error)

	// Wraps a couple of channels around conn.ReadMessage
	go readPump(conn, readDataChan, readErrChan)

	// This is the main control loop. All JS code (including error handlers)
	// should only be executed by this thread to avoid race conditions
	for ***REMOVED***
		select ***REMOVED***
		case ev := <-pingPongChan:
			socket.handleEvent(ev)
		case readData := <-readDataChan:
			socket.handleEvent("message", rt.ToValue(string(readData)))
		case readErr := <-readErrChan:
			socket.handleEvent("error", rt.ToValue(readErr))

		case scheduledFn := <-socket.scheduled:
			scheduledFn(goja.Undefined())

		case <-ctx.Done():
			// This means that K6 is shutting down (e.g., during an interrupt)
			socket.handleEvent("close", rt.ToValue("Interrupt"))
			socket.closeConnection(websocket.CloseGoingAway)

		case <-socket.done:
			// This is the final exit point normally triggered by closeConnection
			duration := float64(time.Since(connectionStart)) / float64(time.Millisecond)
			end := time.Now()

			samples := []stats.Sample***REMOVED***
				***REMOVED***Metric: metrics.WSSessions, Time: end, Tags: tags, Value: 1***REMOVED***,
				***REMOVED***Metric: metrics.WSHandshaking, Time: end, Tags: tags, Value: connecting***REMOVED***,
				***REMOVED***Metric: metrics.WSSessionDuration, Time: end, Tags: tags, Value: duration***REMOVED***,
				***REMOVED***Metric: metrics.DataReceived, Time: end, Tags: tags, Value: float64(traceConn.BytesRead)***REMOVED***,
				***REMOVED***Metric: metrics.DataSent, Time: end, Tags: tags, Value: float64(traceConn.BytesWritten)***REMOVED***,
			***REMOVED***
			state.Samples = append(state.Samples, samples...)

			return response
		***REMOVED***
	***REMOVED***
***REMOVED***

func (s *Socket) On(event string, handler goja.Value) ***REMOVED***
	if handler, ok := goja.AssertFunction(handler); ok ***REMOVED***
		s.eventHandlers[event] = append(s.eventHandlers[event], handler)
	***REMOVED***
***REMOVED***

func (s *Socket) handleEvent(event string, args ...goja.Value) ***REMOVED***
	if handlers, ok := s.eventHandlers[event]; ok ***REMOVED***
		for _, handler := range handlers ***REMOVED***
			handler(goja.Undefined(), args...)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (s *Socket) Send(message string) ***REMOVED***
	// NOTE: No binary message support for the time being since goja doesn't
	// support typed arrays.
	rt := common.GetRuntime(s.ctx)

	writeData := []byte(message)
	if err := s.conn.WriteMessage(websocket.TextMessage, writeData); err != nil ***REMOVED***
		s.handleEvent("error", rt.ToValue(err))
	***REMOVED***
***REMOVED***

func (s *Socket) Ping() ***REMOVED***
	rt := common.GetRuntime(s.ctx)
	deadline := time.Now().Add(writeWait)

	err := s.conn.WriteControl(websocket.PingMessage, []byte***REMOVED******REMOVED***, deadline)
	if err != nil ***REMOVED***
		s.handleEvent("error", rt.ToValue(err))
	***REMOVED***
***REMOVED***

func (s *Socket) SetTimeout(fn goja.Callable, timeoutMs int) ***REMOVED***
	// Starts a goroutine, blocks once on the timeout and pushes the callable
	// back to the main loop through the scheduled channel
	go func() ***REMOVED***
		select ***REMOVED***
		case <-time.After(time.Duration(timeoutMs) * time.Millisecond):
			s.scheduled <- fn

		case <-s.done:
			return
		***REMOVED***
	***REMOVED***()
***REMOVED***

func (s *Socket) SetInterval(fn goja.Callable, intervalMs int) ***REMOVED***
	// Starts a goroutine, blocks forever on the ticker and pushes the callable
	// back to the main loop through the scheduled channel
	go func() ***REMOVED***
		ticker := time.NewTicker(time.Duration(intervalMs) * time.Millisecond)
		defer ticker.Stop()

		for ***REMOVED***
			select ***REMOVED***
			case <-ticker.C:
				s.scheduled <- fn

			case <-s.done:
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***()
***REMOVED***

func (s *Socket) Close(args ...goja.Value) ***REMOVED***
	code := websocket.CloseGoingAway
	if len(args) > 0 && !goja.IsUndefined(args[0]) && !goja.IsNull(args[0]) ***REMOVED***
		code = int(args[0].ToInteger())
	***REMOVED***

	s.closeConnection(code)
***REMOVED***

func (s *Socket) closeConnection(code int) error ***REMOVED***
	// Attempts to close the websocket gracefully

	select ***REMOVED***
	case <-s.done:
		// If the done channel is closed, this means someone has called this
		// function already
		return nil

	case <-time.After(time.Second):
		rt := common.GetRuntime(s.ctx)

		err := s.conn.WriteControl(websocket.CloseMessage,
			websocket.FormatCloseMessage(code, ""),
			time.Now().Add(writeWait),
		)
		if err != nil ***REMOVED***
			s.handleEvent("error", rt.ToValue(err))
			// Just call the handler, we'll try to close the connection anyway
		***REMOVED***
		s.conn.Close()

		close(s.done)
		return err
	***REMOVED***
***REMOVED***

// Wraps conn.ReadMessage in a channel
func readPump(conn *websocket.Conn, readChan chan []byte, errorChan chan error) ***REMOVED***
	defer conn.Close()

	for ***REMOVED***
		_, message, err := conn.ReadMessage()
		if err != nil ***REMOVED***
			// Only emit the error if we didn't close the socket ourselves
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) ***REMOVED***
				errorChan <- err
			***REMOVED***

			return
		***REMOVED***

		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		readChan <- message
	***REMOVED***
***REMOVED***

func parseArgs(rt *goja.Runtime, args []goja.Value) (goja.Callable, map[string]string, http.Header) ***REMOVED***
	var callableV goja.Value
	var paramsV goja.Value

	// The params argument is optional
	if len(args) == 2 ***REMOVED***
		paramsV = args[0]
		callableV = args[1]
	***REMOVED*** else if len(args) == 1 ***REMOVED***
		paramsV = goja.Undefined()
		callableV = args[0]
	***REMOVED*** else ***REMOVED***
		common.Throw(rt, errors.New("Invalid number of arguments to ws.connect"))
		return nil, nil, nil
	***REMOVED***

	// Get the callable (required)
	var callable goja.Callable
	var isFunc bool
	if callable, isFunc = goja.AssertFunction(callableV); !isFunc ***REMOVED***
		common.Throw(rt, errors.New("Last argument to ws.connect must be a function"))
		return nil, nil, nil
	***REMOVED***

	// Leave header to nil by default so we can pass it directly to the Dialer
	var header http.Header
	tags := map[string]string***REMOVED******REMOVED***

	if !goja.IsUndefined(paramsV) && !goja.IsNull(paramsV) ***REMOVED***
		params := paramsV.ToObject(rt)
		for _, k := range params.Keys() ***REMOVED***
			switch k ***REMOVED***
			case "headers":
				header = http.Header***REMOVED******REMOVED***
				headersV := params.Get(k)
				if goja.IsUndefined(headersV) || goja.IsNull(headersV) ***REMOVED***
					continue
				***REMOVED***
				headersObj := headersV.ToObject(rt)
				if headersObj == nil ***REMOVED***
					continue
				***REMOVED***
				for _, key := range headersObj.Keys() ***REMOVED***
					header.Set(key, headersObj.Get(key).String())
				***REMOVED***
			case "tags":
				tagsV := params.Get(k)
				if goja.IsUndefined(tagsV) || goja.IsNull(tagsV) ***REMOVED***
					continue
				***REMOVED***
				tagObj := tagsV.ToObject(rt)
				if tagObj == nil ***REMOVED***
					continue
				***REMOVED***
				for _, key := range tagObj.Keys() ***REMOVED***
					tags[key] = tagObj.Get(key).String()
				***REMOVED***
			***REMOVED***
		***REMOVED***

	***REMOVED***

	return callable, tags, header
***REMOVED***
