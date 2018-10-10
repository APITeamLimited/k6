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

package dummy

import (
	"context"

	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/stats"
	log "github.com/sirupsen/logrus"
)

// Collector implements the lib.Collector interface and should be used only for testing
type Collector struct ***REMOVED***
	RunStatus lib.RunStatus

	SampleContainers []stats.SampleContainer
	Samples          []stats.Sample
***REMOVED***

// Verify that Collector implements lib.Collector
var _ lib.Collector = &Collector***REMOVED******REMOVED***

// Init does nothing, it's only included to satisfy the lib.Collector interface
func (c *Collector) Init() error ***REMOVED*** return nil ***REMOVED***

// MakeConfig does nothing, it's only included to satisfy the lib.Collector interface
func (c *Collector) MakeConfig() interface***REMOVED******REMOVED*** ***REMOVED*** return nil ***REMOVED***

// Run just blocks until the context is done
func (c *Collector) Run(ctx context.Context) ***REMOVED***
	<-ctx.Done()
	log.Debugf("finished status: %d", c.RunStatus)
***REMOVED***

// Collect just appends all of the samples passed to it to the internal sample slice.
// According to the the lib.Collector interface, it should never be called concurrently,
// so there's no locking on purpose - that way Go's race condition detector can actually
// detect incorrect usage.
// Also, theoretically the collector doesn't have to actually Run() before samples start
// being collected, it only has to be initialized.
func (c *Collector) Collect(scs []stats.SampleContainer) ***REMOVED***
	for _, sc := range scs ***REMOVED***
		c.SampleContainers = append(c.SampleContainers, sc)
		c.Samples = append(c.Samples, sc.GetSamples()...)
	***REMOVED***
***REMOVED***

// Link returns a dummy string, it's only included to satisfy the lib.Collector interface
func (c *Collector) Link() string ***REMOVED***
	return "http://example.com/"
***REMOVED***

// GetRequiredSystemTags returns which sample tags are needed by this collector
func (c *Collector) GetRequiredSystemTags() lib.TagSet ***REMOVED***
	return lib.TagSet***REMOVED******REMOVED*** // There are no required tags for this collector
***REMOVED***

// SetRunStatus just saves the passed status for later inspection
func (c *Collector) SetRunStatus(status lib.RunStatus) ***REMOVED***
	c.RunStatus = status
***REMOVED***
