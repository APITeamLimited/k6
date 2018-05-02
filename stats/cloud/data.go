/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2018 Load Impact
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

package cloud

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/loadimpact/k6/lib/metrics"
	"github.com/loadimpact/k6/lib/netext"
	"github.com/loadimpact/k6/stats"
)

const DataTypeSingle = "Point"
const DataTypeMap = "Points"
const DataTypeAggregatedHTTPReqs = "AggregatedPoints"

// Timestamp is used for sending times encoded as microsecond UNIX timestamps to the cloud servers
type Timestamp time.Time

// MarshalJSON encodes the microsecond UNIX timestamps as strings because JavaScripts doesn't have actual integers and tends to round big numbers
func (ct Timestamp) MarshalJSON() ([]byte, error) ***REMOVED***
	return []byte(`"` + strconv.FormatInt(time.Time(ct).UnixNano()/1000, 10) + `"`), nil
***REMOVED***

// UnmarshalJSON decodes the string-enclosed microsecond timestamp back into the proper time.Time alias
func (ct *Timestamp) UnmarshalJSON(p []byte) error ***REMOVED***
	var s string
	if err := json.Unmarshal(p, &s); err != nil ***REMOVED***
		return err
	***REMOVED***
	microSecs, err := strconv.ParseInt(s, 10, 64)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	*ct = Timestamp(time.Unix(microSecs/1000000, (microSecs%1000000)*1000))
	return nil
***REMOVED***

// Sample is the generic struct that contains all types of data that we send to the cloud.
type Sample struct ***REMOVED***
	Type   string      `json:"type"`
	Metric string      `json:"metric"`
	Data   interface***REMOVED******REMOVED*** `json:"data"`
***REMOVED***

// UnmarshalJSON decodes the Data into the corresponding struct
func (ct *Sample) UnmarshalJSON(p []byte) error ***REMOVED***
	var tmpSample struct ***REMOVED***
		Type   string          `json:"type"`
		Metric string          `json:"metric"`
		Data   json.RawMessage `json:"data"`
	***REMOVED***
	if err := json.Unmarshal(p, &tmpSample); err != nil ***REMOVED***
		return err
	***REMOVED***
	s := Sample***REMOVED***
		Type:   tmpSample.Type,
		Metric: tmpSample.Metric,
	***REMOVED***

	switch tmpSample.Type ***REMOVED***
	case DataTypeSingle:
		s.Data = new(SampleDataSingle)
	case DataTypeMap:
		s.Data = new(SampleDataMap)
	case DataTypeAggregatedHTTPReqs:
		s.Data = new(SampleDataAggregatedHTTPReqs)
	default:
		return fmt.Errorf("Unknown sample type '%s'", tmpSample.Type)
	***REMOVED***

	if err := json.Unmarshal(tmpSample.Data, &s.Data); err != nil ***REMOVED***
		return err
	***REMOVED***

	*ct = s
	return nil
***REMOVED***

// SampleDataSingle is used in all simple un-aggregated single-value samples.
type SampleDataSingle struct ***REMOVED***
	Time  Timestamp         `json:"time"`
	Type  stats.MetricType  `json:"type"`
	Tags  *stats.SampleTags `json:"tags,omitempty"`
	Value float64           `json:"value"`
***REMOVED***

// SampleDataMap is used by samples that contain multiple values, currently
// that's only iteration metrics (`iter_li_all`) and unaggregated HTTP
// requests (`http_req_li_all`).
type SampleDataMap struct ***REMOVED***
	Time   Timestamp          `json:"time"`
	Type   stats.MetricType   `json:"type"`
	Tags   *stats.SampleTags  `json:"tags,omitempty"`
	Values map[string]float64 `json:"values,omitempty"`
***REMOVED***

// NewSampleFromTrail just creates a ready-to-send Sample instance
// directly from a netext.Trail.
func NewSampleFromTrail(trail *netext.Trail) *Sample ***REMOVED***
	return &Sample***REMOVED***
		Type:   DataTypeMap,
		Metric: "http_req_li_all",
		Data: &SampleDataMap***REMOVED***
			Time: Timestamp(trail.GetTime()),
			Tags: trail.GetTags(),
			Values: map[string]float64***REMOVED***
				metrics.HTTPReqs.Name:        1,
				metrics.HTTPReqDuration.Name: stats.D(trail.Duration),

				metrics.HTTPReqBlocked.Name:        stats.D(trail.Blocked),
				metrics.HTTPReqConnecting.Name:     stats.D(trail.Connecting),
				metrics.HTTPReqTLSHandshaking.Name: stats.D(trail.TLSHandshaking),
				metrics.HTTPReqSending.Name:        stats.D(trail.Sending),
				metrics.HTTPReqWaiting.Name:        stats.D(trail.Waiting),
				metrics.HTTPReqReceiving.Name:      stats.D(trail.Receiving),
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
***REMOVED***

// SampleDataAggregatedHTTPReqs is used in aggregated samples for HTTP requests.
type SampleDataAggregatedHTTPReqs struct ***REMOVED***
	Time   Timestamp         `json:"time"`
	Type   string            `json:"type"`
	Count  uint64            `json:"count"`
	Tags   *stats.SampleTags `json:"tags,omitempty"`
	Values struct ***REMOVED***
		Duration       AggregatedMetric `json:"http_req_duration"`
		Blocked        AggregatedMetric `json:"http_req_blocked"`
		Connecting     AggregatedMetric `json:"http_req_connecting"`
		TLSHandshaking AggregatedMetric `json:"http_req_tls_handshaking"`
		Sending        AggregatedMetric `json:"http_req_sending"`
		Waiting        AggregatedMetric `json:"http_req_waiting"`
		Receiving      AggregatedMetric `json:"http_req_receiving"`
	***REMOVED*** `json:"values"`
***REMOVED***

// Add updates all agregated values with the supplied trail data
func (sdagg *SampleDataAggregatedHTTPReqs) Add(trail *netext.Trail) ***REMOVED***
	sdagg.Count++
	sdagg.Values.Duration.Add(trail.Duration)
	sdagg.Values.Blocked.Add(trail.Blocked)
	sdagg.Values.Connecting.Add(trail.Connecting)
	sdagg.Values.TLSHandshaking.Add(trail.TLSHandshaking)
	sdagg.Values.Sending.Add(trail.Sending)
	sdagg.Values.Waiting.Add(trail.Waiting)
	sdagg.Values.Receiving.Add(trail.Receiving)
***REMOVED***

// CalcAverages calculates and sets all `Avg` properties in the `Values` struct
func (sdagg *SampleDataAggregatedHTTPReqs) CalcAverages() ***REMOVED***
	count := float64(sdagg.Count)
	sdagg.Values.Duration.Calc(count)
	sdagg.Values.Blocked.Calc(count)
	sdagg.Values.Connecting.Calc(count)
	sdagg.Values.TLSHandshaking.Calc(count)
	sdagg.Values.Sending.Calc(count)
	sdagg.Values.Waiting.Calc(count)
	sdagg.Values.Receiving.Calc(count)
***REMOVED***

// AggregatedMetric is used to store aggregated information for a
// particular metric in an SampleDataAggregatedMap.
type AggregatedMetric struct ***REMOVED***
	// Used by Add() to keep working state
	minD time.Duration
	maxD time.Duration
	sumD time.Duration
	// Updated by Calc() and used in the JSON output
	Min float64 `json:"min"`
	Max float64 `json:"max"`
	Avg float64 `json:"avg"`
***REMOVED***

// Add the new duration to the internal sum and update Min and Max if necessary
func (am *AggregatedMetric) Add(t time.Duration) ***REMOVED***
	if am.sumD == 0 || am.minD > t ***REMOVED***
		am.minD = t
	***REMOVED***
	if am.maxD < t ***REMOVED***
		am.maxD = t
	***REMOVED***
	am.sumD += t
***REMOVED***

// Calc populates the float fields for min and max and calulates the average value
func (am *AggregatedMetric) Calc(count float64) ***REMOVED***
	am.Min = stats.D(am.minD)
	am.Max = stats.D(am.maxD)
	am.Avg = stats.D(am.sumD) / count
***REMOVED***

type aggregationBucket map[*stats.SampleTags][]*netext.Trail

type durations []time.Duration

func (d durations) Len() int           ***REMOVED*** return len(d) ***REMOVED***
func (d durations) Swap(i, j int)      ***REMOVED*** d[i], d[j] = d[j], d[i] ***REMOVED***
func (d durations) Less(i, j int) bool ***REMOVED*** return d[i] < d[j] ***REMOVED***
func (d durations) GetNormalBounds(iqrCoef float64) (min, max time.Duration) ***REMOVED***
	l := len(d)
	if l == 0 ***REMOVED***
		return
	***REMOVED***

	sort.Sort(d)
	var q1, q3 time.Duration
	if l%4 == 0 ***REMOVED***
		q1 = d[l/4]
		q3 = d[(l/4)*3]
	***REMOVED*** else ***REMOVED***
		q1 = (d[l/4] + d[(l/4)+1]) / 2
		q3 = (d[(l/4)*3] + d[(l/4)*3+1]) / 2
	***REMOVED***

	iqr := float64(q3 - q1)
	min = q1 - time.Duration(iqrCoef*iqr)
	max = q3 + time.Duration(iqrCoef*iqr)
	return
***REMOVED***
