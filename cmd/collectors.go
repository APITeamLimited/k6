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

package cmd

import (
	"fmt"
	"strings"

	"gopkg.in/guregu/null.v3"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"

	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/consts"
	"github.com/loadimpact/k6/loader"
	"github.com/loadimpact/k6/stats"
	"github.com/loadimpact/k6/stats/cloud"
	"github.com/loadimpact/k6/stats/csv"
	"github.com/loadimpact/k6/stats/datadog"
	"github.com/loadimpact/k6/stats/influxdb"
	jsonc "github.com/loadimpact/k6/stats/json"
	"github.com/loadimpact/k6/stats/kafka"
	"github.com/loadimpact/k6/stats/statsd"
)

const (
	collectorInfluxDB = "influxdb"
	collectorJSON     = "json"
	collectorKafka    = "kafka"
	collectorCloud    = "cloud"
	collectorStatsD   = "statsd"
	collectorDatadog  = "datadog"
	collectorCSV      = "csv"
)

func parseCollector(s string) (t, arg string) ***REMOVED***
	parts := strings.SplitN(s, "=", 2)
	switch len(parts) ***REMOVED***
	case 0:
		return "", ""
	case 1:
		return parts[0], ""
	default:
		return parts[0], parts[1]
	***REMOVED***
***REMOVED***

// TODO: totally refactor this...
func getCollector(
	logger logrus.FieldLogger,
	collectorName, arg string, src *loader.SourceData, conf Config, executionPlan []lib.ExecutionStep,
) (lib.Collector, error) ***REMOVED***
	switch collectorName ***REMOVED***
	case collectorJSON:
		return jsonc.New(logger, afero.NewOsFs(), arg)
	case collectorInfluxDB:
		config := influxdb.NewConfig().Apply(conf.Collectors.InfluxDB)
		if err := envconfig.Process("", &config); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		urlConfig, err := influxdb.ParseURL(arg)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		config = config.Apply(urlConfig)

		return influxdb.New(logger, config)
	case collectorCloud:
		config := cloud.NewConfig().Apply(conf.Collectors.Cloud)
		if err := envconfig.Process("", &config); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if arg != "" ***REMOVED***
			config.Name = null.StringFrom(arg)
		***REMOVED***

		return cloud.New(logger, config, src, conf.Options, executionPlan, consts.Version)
	case collectorKafka:
		config := kafka.NewConfig().Apply(conf.Collectors.Kafka)
		if err := envconfig.Process("", &config); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if arg != "" ***REMOVED***
			cmdConfig, err := kafka.ParseArg(arg)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			config = config.Apply(cmdConfig)
		***REMOVED***

		return kafka.New(logger, config)
	case collectorStatsD:
		config := statsd.NewConfig().Apply(conf.Collectors.StatsD)
		if err := envconfig.Process("k6_statsd", &config); err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		return statsd.New(logger, config)
	case collectorDatadog:
		config := datadog.NewConfig().Apply(conf.Collectors.Datadog)
		if err := envconfig.Process("k6_datadog", &config); err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		return datadog.New(logger, config)
	case collectorCSV:
		config := csv.NewConfig().Apply(conf.Collectors.CSV)
		if err := envconfig.Process("", &config); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if arg != "" ***REMOVED***
			cmdConfig, err := csv.ParseArg(arg)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***

			config = config.Apply(cmdConfig)
		***REMOVED***

		return csv.New(logger, afero.NewOsFs(), conf.SystemTags.Map(), config)

	default:
		return nil, errors.Errorf("unknown output type: %s", collectorName)
	***REMOVED***
***REMOVED***

func newCollector(
	logger logrus.FieldLogger,
	collectorName, arg string, src *loader.SourceData, conf Config, executionPlan []lib.ExecutionStep,
) (lib.Collector, error) ***REMOVED***
	collector, err := getCollector(logger, collectorName, arg, src, conf, executionPlan)
	if err != nil ***REMOVED***
		return collector, err
	***REMOVED***

	// Check if all required tags are present
	missingRequiredTags := []string***REMOVED******REMOVED***
	requiredTags := collector.GetRequiredSystemTags()
	for _, tag := range stats.SystemTagSetValues() ***REMOVED***
		if requiredTags.Has(tag) && !conf.SystemTags.Has(tag) ***REMOVED***
			missingRequiredTags = append(missingRequiredTags, tag.String())
		***REMOVED***
	***REMOVED***
	if len(missingRequiredTags) > 0 ***REMOVED***
		return collector, fmt.Errorf(
			"the specified collector '%s' needs the following system tags enabled: %s",
			collectorName,
			strings.Join(missingRequiredTags, ", "),
		)
	***REMOVED***

	return collector, nil
***REMOVED***
