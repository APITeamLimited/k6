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

package csv

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"gopkg.in/guregu/null.v3"

	"github.com/mstoykov/envconfig"
	"github.com/sirupsen/logrus"
	"go.k6.io/k6/lib/types"
)

// Config is the config for the csv output
type Config struct ***REMOVED***
	// Samples.
	FileName     null.String        `json:"file_name" envconfig:"K6_CSV_FILENAME"`
	SaveInterval types.NullDuration `json:"save_interval" envconfig:"K6_CSV_SAVE_INTERVAL"`
	TimeFormat   TimeFormat         `json:"time_format" envconfig:"K6_CSV_TIME_FORMAT"`
***REMOVED***

// NewConfig creates a new Config instance with default values for some fields.
func NewConfig() Config ***REMOVED***
	return Config***REMOVED***
		FileName:     null.StringFrom("file.csv"),
		SaveInterval: types.NullDurationFrom(1 * time.Second),
		TimeFormat:   Unix,
	***REMOVED***
***REMOVED***

// Apply merges two configs by overwriting properties in the old config
func (c Config) Apply(cfg Config) Config ***REMOVED***
	if cfg.FileName.Valid ***REMOVED***
		c.FileName = cfg.FileName
	***REMOVED***
	if cfg.SaveInterval.Valid ***REMOVED***
		c.SaveInterval = cfg.SaveInterval
	***REMOVED***
	c.TimeFormat = cfg.TimeFormat
	return c
***REMOVED***

// ParseArg takes an arg string and converts it to a config
func ParseArg(arg string, logger *logrus.Logger) (Config, error) ***REMOVED***
	c := Config***REMOVED******REMOVED***

	if !strings.Contains(arg, "=") ***REMOVED***
		c.FileName = null.StringFrom(arg)
		c.SaveInterval = types.NullDurationFrom(1 * time.Second)
		c.TimeFormat = Unix
		return c, nil
	***REMOVED***

	pairs := strings.Split(arg, ",")
	for _, pair := range pairs ***REMOVED***
		r := strings.SplitN(pair, "=", 2)
		if len(r) != 2 ***REMOVED***
			return c, fmt.Errorf("couldn't parse %q as argument for csv output", arg)
		***REMOVED***
		switch r[0] ***REMOVED***
		case "save_interval":
			logger.Warnf("CSV output argument '%s' is deprecated, please use 'saveInterval' instead.", r[0])
			fallthrough
		case "saveInterval":
			err := c.SaveInterval.UnmarshalText([]byte(r[1]))
			if err != nil ***REMOVED***
				return c, err
			***REMOVED***
		case "file_name":
			logger.Warnf("CSV output argument '%s' is deprecated, please use 'fileName' instead.", r[0])
			fallthrough
		case "fileName":
			c.FileName = null.StringFrom(r[1])
		case "timeFormat":
			c.TimeFormat = TimeFormat(r[1])
			if !c.TimeFormat.IsValid() ***REMOVED***
				return c, fmt.Errorf("unknown value %q as argument for csv output timeFormat, expected 'unix' or 'rfc3399'", arg)
			***REMOVED***

		default:
			return c, fmt.Errorf("unknown key %q as argument for csv output", r[0])
		***REMOVED***
	***REMOVED***

	return c, nil
***REMOVED***

// GetConsolidatedConfig combines ***REMOVED***default config values + JSON config +
// environment vars + arg config values***REMOVED***, and returns the final result.
func GetConsolidatedConfig(
	jsonRawConf json.RawMessage, env map[string]string, arg string, logger *logrus.Logger,
) (Config, error) ***REMOVED***
	result := NewConfig()
	if jsonRawConf != nil ***REMOVED***
		jsonConf := Config***REMOVED******REMOVED***
		if err := json.Unmarshal(jsonRawConf, &jsonConf); err != nil ***REMOVED***
			return result, err
		***REMOVED***
		result = result.Apply(jsonConf)
	***REMOVED***

	envConfig := Config***REMOVED******REMOVED***
	if err := envconfig.Process("", &envConfig, func(key string) (string, bool) ***REMOVED***
		v, ok := env[key]
		return v, ok
	***REMOVED***); err != nil ***REMOVED***
		// TODO: get rid of envconfig and actually use the env parameter...
		return result, err
	***REMOVED***
	result = result.Apply(envConfig)

	if arg != "" ***REMOVED***
		urlConf, err := ParseArg(arg, logger)
		if err != nil ***REMOVED***
			return result, err
		***REMOVED***
		result = result.Apply(urlConf)
	***REMOVED***

	return result, nil
***REMOVED***
