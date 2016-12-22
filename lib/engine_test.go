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

package lib

import (
	"context"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v3"
	"runtime"
	"testing"
	"time"
)

// Helper for asserting the number of active/dead VUs.
func assertActiveVUs(t *testing.T, e *Engine, active, dead int) ***REMOVED***
	var numActive, numDead int
	var lastWasDead bool
	for _, vu := range e.vuEntries ***REMOVED***
		if vu.Cancel != nil ***REMOVED***
			numActive++
			assert.False(t, lastWasDead, "living vu in dead zone")
		***REMOVED*** else ***REMOVED***
			numDead++
			lastWasDead = true
		***REMOVED***
	***REMOVED***
	assert.Equal(t, active, numActive, "wrong number of active vus")
	assert.Equal(t, dead, numDead, "wrong number of dead vus")
***REMOVED***

func TestNewEngine(t *testing.T) ***REMOVED***
	_, err := NewEngine(nil, Options***REMOVED******REMOVED***)
	assert.NoError(t, err)
***REMOVED***

func TestNewEngineOptions(t *testing.T) ***REMOVED***
	t.Run("VUsMax", func(t *testing.T) ***REMOVED***
		t.Run("not set", func(t *testing.T) ***REMOVED***
			e, err := NewEngine(nil, Options***REMOVED******REMOVED***)
			assert.NoError(t, err)
			assert.Equal(t, int64(0), e.GetVUsMax())
			assert.Equal(t, int64(0), e.GetVUs())
		***REMOVED***)
		t.Run("set", func(t *testing.T) ***REMOVED***
			e, err := NewEngine(nil, Options***REMOVED***
				VUsMax: null.IntFrom(10),
			***REMOVED***)
			assert.NoError(t, err)
			assert.Equal(t, int64(10), e.GetVUsMax())
			assert.Equal(t, int64(0), e.GetVUs())
		***REMOVED***)
	***REMOVED***)
	t.Run("VUs", func(t *testing.T) ***REMOVED***
		t.Run("no max", func(t *testing.T) ***REMOVED***
			_, err := NewEngine(nil, Options***REMOVED***
				VUs: null.IntFrom(10),
			***REMOVED***)
			assert.EqualError(t, err, "more vus than allocated requested")
		***REMOVED***)
		t.Run("max too low", func(t *testing.T) ***REMOVED***
			_, err := NewEngine(nil, Options***REMOVED***
				VUsMax: null.IntFrom(1),
				VUs:    null.IntFrom(10),
			***REMOVED***)
			assert.EqualError(t, err, "more vus than allocated requested")
		***REMOVED***)
		t.Run("max higher", func(t *testing.T) ***REMOVED***
			e, err := NewEngine(nil, Options***REMOVED***
				VUsMax: null.IntFrom(10),
				VUs:    null.IntFrom(1),
			***REMOVED***)
			assert.NoError(t, err)
			assert.Equal(t, int64(10), e.GetVUsMax())
			assert.Equal(t, int64(1), e.GetVUs())
		***REMOVED***)
		t.Run("max just right", func(t *testing.T) ***REMOVED***
			e, err := NewEngine(nil, Options***REMOVED***
				VUsMax: null.IntFrom(10),
				VUs:    null.IntFrom(10),
			***REMOVED***)
			assert.NoError(t, err)
			assert.Equal(t, int64(10), e.GetVUsMax())
			assert.Equal(t, int64(10), e.GetVUs())
		***REMOVED***)
	***REMOVED***)
	t.Run("Paused", func(t *testing.T) ***REMOVED***
		t.Run("not set", func(t *testing.T) ***REMOVED***
			e, err := NewEngine(nil, Options***REMOVED******REMOVED***)
			assert.NoError(t, err)
			assert.False(t, e.IsPaused())
		***REMOVED***)
		t.Run("false", func(t *testing.T) ***REMOVED***
			e, err := NewEngine(nil, Options***REMOVED***
				Paused: null.BoolFrom(false),
			***REMOVED***)
			assert.NoError(t, err)
			assert.False(t, e.IsPaused())
		***REMOVED***)
		t.Run("true", func(t *testing.T) ***REMOVED***
			e, err := NewEngine(nil, Options***REMOVED***
				Paused: null.BoolFrom(true),
			***REMOVED***)
			assert.NoError(t, err)
			assert.True(t, e.IsPaused())
		***REMOVED***)
	***REMOVED***)
***REMOVED***

func TestEngineRun(t *testing.T) ***REMOVED***
	t.Run("exits with context", func(t *testing.T) ***REMOVED***
		startTime := time.Now()
		duration := 100 * time.Millisecond
		e, err := NewEngine(nil, Options***REMOVED******REMOVED***)
		assert.NoError(t, err)

		ctx, _ := context.WithTimeout(context.Background(), duration)
		assert.NoError(t, e.Run(ctx))
		assert.WithinDuration(t, startTime.Add(duration), time.Now(), 100*time.Millisecond)
	***REMOVED***)
	t.Run("terminates subctx", func(t *testing.T) ***REMOVED***
		e, err := NewEngine(nil, Options***REMOVED******REMOVED***)
		assert.NoError(t, err)

		subctx := e.subctx
		select ***REMOVED***
		case <-subctx.Done():
			assert.Fail(t, "context is already terminated")
		default:
		***REMOVED***

		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		assert.NoError(t, e.Run(ctx))

		assert.NotEqual(t, subctx, e.subctx, "subcontext not changed")
		select ***REMOVED***
		case <-subctx.Done():
		default:
			assert.Fail(t, "context was not terminated")
		***REMOVED***
	***REMOVED***)
	t.Run("updates AtTime", func(t *testing.T) ***REMOVED***
		e, err := NewEngine(nil, Options***REMOVED******REMOVED***)
		assert.NoError(t, err)

		d := 50 * time.Millisecond
		ctx, _ := context.WithTimeout(context.Background(), d)
		startTime := time.Now()
		assert.NoError(t, e.Run(ctx))
		assert.WithinDuration(t, startTime.Add(d), startTime.Add(e.AtTime()), 1*TickRate)
	***REMOVED***)
***REMOVED***

func TestEngineIsRunning(t *testing.T) ***REMOVED***
	ctx, cancel := context.WithCancel(context.Background())
	e, err := NewEngine(nil, Options***REMOVED******REMOVED***)
	assert.NoError(t, err)

	go func() ***REMOVED*** assert.NoError(t, e.Run(ctx)) ***REMOVED***()
	runtime.Gosched()
	assert.True(t, e.IsRunning())

	cancel()
	runtime.Gosched()
	assert.False(t, e.IsRunning())
***REMOVED***

func TestEngineSetPaused(t *testing.T) ***REMOVED***
	e, err := NewEngine(nil, Options***REMOVED******REMOVED***)
	assert.NoError(t, err)
	assert.False(t, e.IsPaused())

	e.SetPaused(true)
	assert.True(t, e.IsPaused())

	e.SetPaused(false)
	assert.False(t, e.IsPaused())
***REMOVED***

func TestEngineSetVUsMax(t *testing.T) ***REMOVED***
	t.Run("not set", func(t *testing.T) ***REMOVED***
		e, err := NewEngine(nil, Options***REMOVED******REMOVED***)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), e.GetVUsMax())
		assert.Len(t, e.vuEntries, 0)
	***REMOVED***)
	t.Run("set", func(t *testing.T) ***REMOVED***
		e, err := NewEngine(nil, Options***REMOVED******REMOVED***)
		assert.NoError(t, err)
		assert.NoError(t, e.SetVUsMax(10))
		assert.Equal(t, int64(10), e.GetVUsMax())
		assert.Len(t, e.vuEntries, 10)
		for _, vu := range e.vuEntries ***REMOVED***
			assert.Nil(t, vu.Cancel)
		***REMOVED***

		t.Run("higher", func(t *testing.T) ***REMOVED***
			assert.NoError(t, e.SetVUsMax(15))
			assert.Equal(t, int64(15), e.GetVUsMax())
			assert.Len(t, e.vuEntries, 15)
			for _, vu := range e.vuEntries ***REMOVED***
				assert.Nil(t, vu.Cancel)
			***REMOVED***
		***REMOVED***)

		t.Run("lower", func(t *testing.T) ***REMOVED***
			assert.NoError(t, e.SetVUsMax(5))
			assert.Equal(t, int64(5), e.GetVUsMax())
			assert.Len(t, e.vuEntries, 5)
			for _, vu := range e.vuEntries ***REMOVED***
				assert.Nil(t, vu.Cancel)
			***REMOVED***
		***REMOVED***)
	***REMOVED***)
	t.Run("set negative", func(t *testing.T) ***REMOVED***
		e, err := NewEngine(nil, Options***REMOVED******REMOVED***)
		assert.NoError(t, err)
		assert.EqualError(t, e.SetVUsMax(-1), "vus-max can't be negative")
		assert.Len(t, e.vuEntries, 0)
	***REMOVED***)
	t.Run("set too low", func(t *testing.T) ***REMOVED***
		e, err := NewEngine(nil, Options***REMOVED***
			VUsMax: null.IntFrom(10),
			VUs:    null.IntFrom(10),
		***REMOVED***)
		assert.NoError(t, err)
		assert.EqualError(t, e.SetVUsMax(5), "can't reduce vus-max below vus")
		assert.Len(t, e.vuEntries, 10)
	***REMOVED***)
***REMOVED***

func TestEngineSetVUs(t *testing.T) ***REMOVED***
	t.Run("not set", func(t *testing.T) ***REMOVED***
		e, err := NewEngine(nil, Options***REMOVED******REMOVED***)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), e.GetVUsMax())
		assert.Equal(t, int64(0), e.GetVUs())
	***REMOVED***)
	t.Run("set", func(t *testing.T) ***REMOVED***
		e, err := NewEngine(nil, Options***REMOVED***VUsMax: null.IntFrom(15)***REMOVED***)
		assert.NoError(t, err)
		assert.NoError(t, e.SetVUs(10))
		assert.Equal(t, int64(10), e.GetVUs())
		assertActiveVUs(t, e, 10, 5)

		t.Run("negative", func(t *testing.T) ***REMOVED***
			assert.EqualError(t, e.SetVUs(-1), "vus can't be negative")
			assert.Equal(t, int64(10), e.GetVUs())
			assertActiveVUs(t, e, 10, 5)
		***REMOVED***)

		t.Run("too high", func(t *testing.T) ***REMOVED***
			assert.EqualError(t, e.SetVUs(20), "more vus than allocated requested")
			assert.Equal(t, int64(10), e.GetVUs())
			assertActiveVUs(t, e, 10, 5)
		***REMOVED***)

		t.Run("lower", func(t *testing.T) ***REMOVED***
			assert.NoError(t, e.SetVUs(5))
			assert.Equal(t, int64(5), e.GetVUs())
			assertActiveVUs(t, e, 5, 10)
		***REMOVED***)

		t.Run("higher", func(t *testing.T) ***REMOVED***
			assert.NoError(t, e.SetVUs(15))
			assert.Equal(t, int64(15), e.GetVUs())
			assertActiveVUs(t, e, 15, 0)
		***REMOVED***)
	***REMOVED***)
***REMOVED***
