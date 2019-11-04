/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2019 Load Impact
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

package executor

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/loadimpact/k6/lib/types"
	null "gopkg.in/guregu/null.v3"
)

var executorNameWhitelist = regexp.MustCompile(`^[0-9a-zA-Z_-]+$`) //nolint:gochecknoglobals
const executorNameErr = "the executor name should contain only numbers, latin letters, underscores, and dashes"

// BaseConfig contains the common config fields for all executors
type BaseConfig struct ***REMOVED***
	Name         string             `json:"-"` // set via the JS object key
	Type         string             `json:"type"`
	StartTime    types.NullDuration `json:"startTime"`
	GracefulStop types.NullDuration `json:"gracefulStop"`
	Env          map[string]string  `json:"env"`
	Exec         null.String        `json:"exec"` // function name, externally validated

	//TODO: future extensions like tags, distribution, others?
***REMOVED***

// NewBaseConfig returns a default base config with the default values
func NewBaseConfig(name, configType string) BaseConfig ***REMOVED***
	return BaseConfig***REMOVED***
		Name:         name,
		Type:         configType,
		GracefulStop: types.NewNullDuration(30*time.Second, false),
	***REMOVED***
***REMOVED***

// Validate checks some basic things like present name, type, and a positive start time
func (bc BaseConfig) Validate() (errors []error) ***REMOVED***
	// Some just-in-case checks, since those things are likely checked in other places or
	// even assigned by us:
	if bc.Name == "" ***REMOVED***
		errors = append(errors, fmt.Errorf("executor name shouldn't be empty"))
	***REMOVED***
	if !executorNameWhitelist.MatchString(bc.Name) ***REMOVED***
		errors = append(errors, fmt.Errorf(executorNameErr))
	***REMOVED***
	if bc.Exec.Valid && bc.Exec.String == "" ***REMOVED***
		errors = append(errors, fmt.Errorf("exec value cannot be empty"))
	***REMOVED***
	if bc.Type == "" ***REMOVED***
		errors = append(errors, fmt.Errorf("missing or empty type field"))
	***REMOVED***
	// The actually reasonable checks:
	if bc.StartTime.Duration < 0 ***REMOVED***
		errors = append(errors, fmt.Errorf("the startTime can't be negative"))
	***REMOVED***
	if bc.GracefulStop.Duration < 0 ***REMOVED***
		errors = append(errors, fmt.Errorf("the gracefulStop timeout can't be negative"))
	***REMOVED***
	return errors
***REMOVED***

// GetName returns the name of the executor.
func (bc BaseConfig) GetName() string ***REMOVED***
	return bc.Name
***REMOVED***

// GetType returns the executor's type as a string ID.
func (bc BaseConfig) GetType() string ***REMOVED***
	return bc.Type
***REMOVED***

// GetStartTime returns the starting time, relative to the beginning of the
// actual test, that this executor is supposed to execute.
func (bc BaseConfig) GetStartTime() time.Duration ***REMOVED***
	return time.Duration(bc.StartTime.Duration)
***REMOVED***

// GetGracefulStop returns how long k6 is supposed to wait for any still
// running iterations to finish executing at the end of the normal executor
// duration, before it actually kills them.
//
// Of course, that doesn't count when the user manually interrupts the test,
// then iterations are immediately stopped.
func (bc BaseConfig) GetGracefulStop() time.Duration ***REMOVED***
	return time.Duration(bc.GracefulStop.Duration)
***REMOVED***

// GetEnv returns any specific environment key=value pairs that
// are configured for the executor.
func (bc BaseConfig) GetEnv() map[string]string ***REMOVED***
	return bc.Env
***REMOVED***

// GetExec returns the configured custom exec value, if any.
func (bc BaseConfig) GetExec() null.String ***REMOVED***
	return bc.Exec
***REMOVED***

// IsDistributable returns true since by default all executors could be run in
// a distributed manner.
func (bc BaseConfig) IsDistributable() bool ***REMOVED***
	return true
***REMOVED***

// getBaseInfo is a helper method for the "parent" String methods.
func (bc BaseConfig) getBaseInfo(facts ...string) string ***REMOVED***
	if bc.Exec.Valid ***REMOVED***
		facts = append(facts, fmt.Sprintf("exec: %s", bc.Exec.String))
	***REMOVED***
	if bc.StartTime.Duration > 0 ***REMOVED***
		facts = append(facts, fmt.Sprintf("startTime: %s", bc.StartTime.Duration))
	***REMOVED***
	if bc.GracefulStop.Duration > 0 ***REMOVED***
		facts = append(facts, fmt.Sprintf("gracefulStop: %s", bc.GracefulStop.Duration))
	***REMOVED***
	if len(facts) == 0 ***REMOVED***
		return ""
	***REMOVED***
	return " (" + strings.Join(facts, ", ") + ")"
***REMOVED***
