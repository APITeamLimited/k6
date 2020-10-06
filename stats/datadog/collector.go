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

package datadog

import (
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/guregu/null.v3"

	"github.com/loadimpact/k6/lib/types"
	"github.com/loadimpact/k6/stats"
	"github.com/loadimpact/k6/stats/statsd/common"
)

type tagHandler stats.TagSet

func (t tagHandler) processTags(tags map[string]string) []string ***REMOVED***
	var res []string
	for key, value := range tags ***REMOVED***
		if value != "" && !t[key] ***REMOVED***
			res = append(res, key+":"+value)
		***REMOVED***
	***REMOVED***
	return res
***REMOVED***

// Config defines the Datadog configuration.
type Config struct ***REMOVED***
	Addr         null.String        `json:"addr,omitempty" envconfig:"K6_DATADOG_ADDR"`
	BufferSize   null.Int           `json:"bufferSize,omitempty" envconfig:"K6_DATADOG_BUFFER_SIZE"`
	Namespace    null.String        `json:"namespace,omitempty" envconfig:"K6_DATADOG_NAMESPACE"`
	PushInterval types.NullDuration `json:"pushInterval,omitempty" envconfig:"K6_DATADOG_PUSH_INTERVAL"`
	TagBlacklist stats.TagSet       `json:"tagBlacklist,omitempty" envconfig:"K6_DATADOG_TAG_BLACKLIST"`
***REMOVED***

// GetAddr returns the address of the DogStatsD service.
func (c Config) GetAddr() null.String ***REMOVED***
	return c.Addr
***REMOVED***

// GetBufferSize returns the size of the commands buffer.
func (c Config) GetBufferSize() null.Int ***REMOVED***
	return c.BufferSize
***REMOVED***

// GetNamespace returns the namespace prepended to all statsd calls.
func (c Config) GetNamespace() null.String ***REMOVED***
	return c.Namespace
***REMOVED***

// GetPushInterval returns the time interval between outgoing data batches.
func (c Config) GetPushInterval() types.NullDuration ***REMOVED***
	return c.PushInterval
***REMOVED***

var _ common.Config = &Config***REMOVED******REMOVED***

// Apply saves config non-zero config values from the passed config in the receiver.
func (c Config) Apply(cfg Config) Config ***REMOVED***
	if cfg.Addr.Valid ***REMOVED***
		c.Addr = cfg.Addr
	***REMOVED***
	if cfg.BufferSize.Valid ***REMOVED***
		c.BufferSize = cfg.BufferSize
	***REMOVED***
	if cfg.Namespace.Valid ***REMOVED***
		c.Namespace = cfg.Namespace
	***REMOVED***
	if cfg.PushInterval.Valid ***REMOVED***
		c.PushInterval = cfg.PushInterval
	***REMOVED***
	if cfg.TagBlacklist != nil ***REMOVED***
		c.TagBlacklist = cfg.TagBlacklist
	***REMOVED***

	return c
***REMOVED***

// NewConfig creates a new Config instance with default values for some fields.
func NewConfig() Config ***REMOVED***
	return Config***REMOVED***
		Addr:         null.NewString("localhost:8125", false),
		BufferSize:   null.NewInt(20, false),
		Namespace:    null.NewString("k6.", false),
		PushInterval: types.NewNullDuration(1*time.Second, false),
		TagBlacklist: stats.TagSet***REMOVED******REMOVED***,
	***REMOVED***
***REMOVED***

// New creates a new Datadog connector client
func New(logger logrus.FieldLogger, conf Config) (*common.Collector, error) ***REMOVED***
	return &common.Collector***REMOVED***
		Config:      conf,
		Type:        "datadog",
		ProcessTags: tagHandler(conf.TagBlacklist).processTags,
		Logger:      logger,
	***REMOVED***, nil
***REMOVED***
