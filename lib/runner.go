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
	"github.com/loadimpact/k6/stats"
)

// Ensure RunnerFunc conforms to Runner.
var _ Runner = RunnerFunc(nil)

// A Runner is a factory for VUs.
type Runner interface ***REMOVED***
	// Creates a new VU. As much as possible should be precomputed here, to allow a pool
	// of prepared VUs to be used to quickly scale up and down.
	NewVU() (VU, error)

	// Returns the default (root) group.
	GetDefaultGroup() *Group

	// Returns the option set.
	GetOptions() Options

	// Applies a set of options.
	ApplyOptions(opts Options)
***REMOVED***

// A VU is a Virtual User.
type VU interface ***REMOVED***
	// Runs the VU once. An iteration should be completely self-contained, and no state
	// or open connections should carry over from one iteration to the next.
	RunOnce(ctx context.Context) ([]stats.Sample, error)

	// Called when the VU's identity changes.
	Reconfigure(id int64) error
***REMOVED***

// RunnerFunc adapts a function to be used as both a runner and a VU. NewVU() returns copies of
// the function itself. Mainly useful for testing.
type RunnerFunc func(ctx context.Context) ([]stats.Sample, error)

func (fn RunnerFunc) NewVU() (VU, error) ***REMOVED***
	return fn, nil
***REMOVED***

func (fn RunnerFunc) GetDefaultGroup() *Group ***REMOVED***
	return nil
***REMOVED***

func (fn RunnerFunc) GetOptions() Options ***REMOVED***
	return Options***REMOVED******REMOVED***
***REMOVED***

func (fn RunnerFunc) ApplyOptions(opts Options) ***REMOVED***
***REMOVED***

func (fn RunnerFunc) RunOnce(ctx context.Context) ([]stats.Sample, error) ***REMOVED***
	return fn(ctx)
***REMOVED***

func (fn RunnerFunc) Reconfigure(id int64) error ***REMOVED***
	return nil
***REMOVED***
