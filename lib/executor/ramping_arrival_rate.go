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
	"context"
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/guregu/null.v3"

	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/metrics"
	"go.k6.io/k6/lib/types"
	"go.k6.io/k6/stats"
	"go.k6.io/k6/ui/pb"
)

const rampingArrivalRateType = "ramping-arrival-rate"

func init() ***REMOVED***
	lib.RegisterExecutorConfigType(
		rampingArrivalRateType,
		func(name string, rawJSON []byte) (lib.ExecutorConfig, error) ***REMOVED***
			config := NewRampingArrivalRateConfig(name)
			err := lib.StrictJSONUnmarshal(rawJSON, &config)
			return config, err
		***REMOVED***,
	)
***REMOVED***

// RampingArrivalRateConfig stores config for the ramping (i.e. variable)
// arrival-rate executor.
type RampingArrivalRateConfig struct ***REMOVED***
	BaseConfig
	StartRate null.Int           `json:"startRate"`
	TimeUnit  types.NullDuration `json:"timeUnit"`
	Stages    []Stage            `json:"stages"`

	// Initialize `PreAllocatedVUs` number of VUs, and if more than that are needed,
	// they will be dynamically allocated, until `MaxVUs` is reached, which is an
	// absolutely hard limit on the number of VUs the executor will use
	PreAllocatedVUs null.Int `json:"preAllocatedVUs"`
	MaxVUs          null.Int `json:"maxVUs"`
***REMOVED***

// NewRampingArrivalRateConfig returns a RampingArrivalRateConfig with default values
func NewRampingArrivalRateConfig(name string) *RampingArrivalRateConfig ***REMOVED***
	return &RampingArrivalRateConfig***REMOVED***
		BaseConfig: NewBaseConfig(name, rampingArrivalRateType),
		TimeUnit:   types.NewNullDuration(1*time.Second, false),
	***REMOVED***
***REMOVED***

// Make sure we implement the lib.ExecutorConfig interface
var _ lib.ExecutorConfig = &RampingArrivalRateConfig***REMOVED******REMOVED***

// GetPreAllocatedVUs is just a helper method that returns the scaled pre-allocated VUs.
func (varc RampingArrivalRateConfig) GetPreAllocatedVUs(et *lib.ExecutionTuple) int64 ***REMOVED***
	return et.Segment.Scale(varc.PreAllocatedVUs.Int64)
***REMOVED***

// GetMaxVUs is just a helper method that returns the scaled max VUs.
func (varc RampingArrivalRateConfig) GetMaxVUs(et *lib.ExecutionTuple) int64 ***REMOVED***
	return et.Segment.Scale(varc.MaxVUs.Int64)
***REMOVED***

// GetDescription returns a human-readable description of the executor options
func (varc RampingArrivalRateConfig) GetDescription(et *lib.ExecutionTuple) string ***REMOVED***
	// TODO: something better? always show iterations per second?
	maxVUsRange := fmt.Sprintf("maxVUs: %d", et.Segment.Scale(varc.PreAllocatedVUs.Int64))
	if varc.MaxVUs.Int64 > varc.PreAllocatedVUs.Int64 ***REMOVED***
		maxVUsRange += fmt.Sprintf("-%d", et.Segment.Scale(varc.MaxVUs.Int64))
	***REMOVED***
	maxUnscaledRate := getStagesUnscaledMaxTarget(varc.StartRate.Int64, varc.Stages)
	maxArrRatePerSec, _ := getArrivalRatePerSec(
		getScaledArrivalRate(et.Segment, maxUnscaledRate, time.Duration(varc.TimeUnit.Duration)),
	).Float64()

	return fmt.Sprintf("Up to %.2f iterations/s for %s over %d stages%s",
		maxArrRatePerSec, sumStagesDuration(varc.Stages),
		len(varc.Stages), varc.getBaseInfo(maxVUsRange))
***REMOVED***

// Validate makes sure all options are configured and valid
func (varc *RampingArrivalRateConfig) Validate() []error ***REMOVED***
	errors := varc.BaseConfig.Validate()

	if varc.StartRate.Int64 < 0 ***REMOVED***
		errors = append(errors, fmt.Errorf("the startRate value shouldn't be negative"))
	***REMOVED***

	if time.Duration(varc.TimeUnit.Duration) < 0 ***REMOVED***
		errors = append(errors, fmt.Errorf("the timeUnit should be more than 0"))
	***REMOVED***

	errors = append(errors, validateStages(varc.Stages)...)

	if !varc.PreAllocatedVUs.Valid ***REMOVED***
		errors = append(errors, fmt.Errorf("the number of preAllocatedVUs isn't specified"))
	***REMOVED*** else if varc.PreAllocatedVUs.Int64 < 0 ***REMOVED***
		errors = append(errors, fmt.Errorf("the number of preAllocatedVUs shouldn't be negative"))
	***REMOVED***

	if !varc.MaxVUs.Valid ***REMOVED***
		// TODO: don't change the config while validating
		varc.MaxVUs.Int64 = varc.PreAllocatedVUs.Int64
	***REMOVED*** else if varc.MaxVUs.Int64 < varc.PreAllocatedVUs.Int64 ***REMOVED***
		errors = append(errors, fmt.Errorf("maxVUs shouldn't be less than preAllocatedVUs"))
	***REMOVED***

	return errors
***REMOVED***

// GetExecutionRequirements returns the number of required VUs to run the
// executor for its whole duration (disregarding any startTime), including the
// maximum waiting time for any iterations to gracefully stop. This is used by
// the execution scheduler in its VU reservation calculations, so it knows how
// many VUs to pre-initialize.
func (varc RampingArrivalRateConfig) GetExecutionRequirements(et *lib.ExecutionTuple) []lib.ExecutionStep ***REMOVED***
	return []lib.ExecutionStep***REMOVED***
		***REMOVED***
			TimeOffset:      0,
			PlannedVUs:      uint64(et.Segment.Scale(varc.PreAllocatedVUs.Int64)),
			MaxUnplannedVUs: uint64(et.Segment.Scale(varc.MaxVUs.Int64 - varc.PreAllocatedVUs.Int64)),
		***REMOVED***,
		***REMOVED***
			TimeOffset:      sumStagesDuration(varc.Stages) + time.Duration(varc.GracefulStop.Duration),
			PlannedVUs:      0,
			MaxUnplannedVUs: 0,
		***REMOVED***,
	***REMOVED***
***REMOVED***

// NewExecutor creates a new RampingArrivalRate executor
func (varc RampingArrivalRateConfig) NewExecutor(
	es *lib.ExecutionState, logger *logrus.Entry,
) (lib.Executor, error) ***REMOVED***
	return RampingArrivalRate***REMOVED***
		BaseExecutor: NewBaseExecutor(&varc, es, logger),
		config:       varc,
	***REMOVED***, nil
***REMOVED***

// HasWork reports whether there is any work to be done for the given execution segment.
func (varc RampingArrivalRateConfig) HasWork(et *lib.ExecutionTuple) bool ***REMOVED***
	return varc.GetMaxVUs(et) > 0
***REMOVED***

// RampingArrivalRate tries to execute a specific number of iterations for a
// specific period.
// TODO: combine with the ConstantArrivalRate?
type RampingArrivalRate struct ***REMOVED***
	*BaseExecutor
	config RampingArrivalRateConfig
***REMOVED***

// Make sure we implement the lib.Executor interface.
var _ lib.Executor = &RampingArrivalRate***REMOVED******REMOVED***

// cal calculates the  transtitions between stages and gives the next full value produced by the
// stages. In this explanation we are talking about events and in practice those events are starting
// of an iteration, but could really be anything that needs to occur at a constant or linear rate.
//
// The basic idea is that we make a graph with the X axis being time and the Y axis being
// events/s we know that the area of the figure between the graph and the X axis is equal to the
// amount of events done - we multiply time by events per time so we get events ...
// Mathematics :).
//
// Lets look at a simple example - lets say we start with 2 events and the first stage is 5
// seconds to 2 events/s and then we have a second stage for 5 second that goes up to 3 events
// (using small numbers because ... well it is easier :D). This will look something like:
//  ^
// 7|
// 6|
// 5|
// 4|
// 3|       ,-+
// 2|----+-'  |
// 1|    |    |
//  +----+----+---------------------------------->
//  0s   5s   10s
// TODO: bigger and more stages
//
// Now the question is when(where on the graph) does the first event happen? Well in this simple
// case it is easy it will be at 0.5 seconds as we are doing 2 events/s. If we want to know when
// event n will happen we need to calculate n = 2 * x, where x is the time it will happen, so we
// need to calculate x = n/2as we are interested in the time, x.
// So if we just had a constant function for each event n we can calculate n/2 and find out when
// it needs to start.
// As we can see though the graph changes as stages change. But we can calculate how many events
// each stage will have, again it is the area from the start of the stage to it's end and between
// the graph and the X axis. So in this case we know that the first stage will have 10 full events
// in it and no more or less. So we are trying to find out when the 12 event will happen the answer
// will be after the 5th second.
//
// The graph doesn't show this well but we are ramping up linearly (we could possibly add
// other ramping up/down functions later). So at 7.5 seconds for example we should be doing 2.5
// events/s. You could start slicing the graph constantly and in this way to represent the ramping
// up/down as a multiple constant functions, and you will get mostly okayish results. But here is
// where calculus comes into play. Calculus gives us a way of exactly calculate the area for any
// given function and linear ramp up/downs just happen to be pretty easy(actual math prove in
// https://github.com/loadimpact/k6/issues/1299#issuecomment-575661084).
//
// One tricky last point is what happens if stage only completes 9.8 events? Let's say that the
// first stage above was 4.9 seconds long 2 * 4.9 is 9.8, we have 9 events and .8 of an event, what
// do with do with that? Well the 10th even will happen in the next stage (if any) and will happen
// when the are from the start till time x is 0.2 (instead of 1) as 0.2 + 0.8 is 10. So the 12th for
// example will be when the area is 2.2 as 9.8+2.2. So we just carry this around.
//
// So in the end what calis doing is to get formulas which will tell it when
// a given event n in order will happen. It helps itself by knowing that in a given
// stage will do some given amount (the area of the stage) events and if we past that one we
// know we are not in that stage.
//
// The specific implementation here can only go forward and does incorporate
// the striping algorithm from the lib.ExecutionTuple for additional speed up but this could
// possibly be refactored if need for this arises.
func (varc RampingArrivalRateConfig) cal(et *lib.ExecutionTuple, ch chan<- time.Duration) ***REMOVED***
	start, offsets, _ := et.GetStripedOffsets()
	li := -1
	// TODO: move this to a utility function, or directly what GetStripedOffsets uses once we see everywhere we will use it
	next := func() int64 ***REMOVED***
		li++
		return offsets[li%len(offsets)]
	***REMOVED***
	defer close(ch) // TODO: maybe this is not a good design - closing a channel we get
	var (
		stageStart                   time.Duration
		timeUnit                     = float64(varc.TimeUnit.Duration)
		doneSoFar, endCount, to, dur float64
		from                         = float64(varc.StartRate.ValueOrZero()) / timeUnit
		// start .. starts at 0 but the algorithm works with area so we need to start from 1 not 0
		i = float64(start + 1)
	)

	for _, stage := range varc.Stages ***REMOVED***
		to = float64(stage.Target.ValueOrZero()) / timeUnit
		dur = float64(stage.Duration.Duration)
		if from != to ***REMOVED*** // ramp up/down
			endCount += dur * ((to-from)/2 + from)
			for ; i <= endCount; i += float64(next()) ***REMOVED***
				// TODO: try to twist this in a way to be able to get i (the only changing part)
				// somewhere where it is less in the middle of the equation
				x := (from*dur - math.Sqrt(dur*(from*from*dur+2*(i-doneSoFar)*(to-from)))) / (from - to)

				ch <- time.Duration(x) + stageStart
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			endCount += dur * to
			for ; i <= endCount; i += float64(next()) ***REMOVED***
				ch <- time.Duration((i-doneSoFar)/to) + stageStart
			***REMOVED***
		***REMOVED***
		doneSoFar = endCount
		from = to
		stageStart += time.Duration(stage.Duration.Duration)
	***REMOVED***
***REMOVED***

// Run executes a variable number of iterations per second.
//
// TODO: Split this up and make an independent component that can be reused
// between the constant and ramping arrival rate executors - that way we can
// keep the complexity in one well-architected part (with short methods and few
// lambdas :D), while having both config frontends still be present for maximum
// UX benefits. Basically, keep the progress bars and scheduling (i.e. at what
// time should iteration X begin) different, but keep everyhing else the same.
// This will allow us to implement https://github.com/loadimpact/k6/issues/1386
// and things like all of the TODOs below in one place only.
//nolint:funlen,gocognit
func (varr RampingArrivalRate) Run(parentCtx context.Context, out chan<- stats.SampleContainer) (err error) ***REMOVED***
	segment := varr.executionState.ExecutionTuple.Segment
	gracefulStop := varr.config.GetGracefulStop()
	duration := sumStagesDuration(varr.config.Stages)
	preAllocatedVUs := varr.config.GetPreAllocatedVUs(varr.executionState.ExecutionTuple)
	maxVUs := varr.config.GetMaxVUs(varr.executionState.ExecutionTuple)

	// TODO: refactor and simplify
	timeUnit := time.Duration(varr.config.TimeUnit.Duration)
	startArrivalRate := getScaledArrivalRate(segment, varr.config.StartRate.Int64, timeUnit)
	maxUnscaledRate := getStagesUnscaledMaxTarget(varr.config.StartRate.Int64, varr.config.Stages)
	maxArrivalRatePerSec, _ := getArrivalRatePerSec(getScaledArrivalRate(segment, maxUnscaledRate, timeUnit)).Float64()
	startTickerPeriod := getTickerPeriod(startArrivalRate)

	// Make sure the log and the progress bar have accurate information
	varr.logger.WithFields(logrus.Fields***REMOVED***
		"maxVUs": maxVUs, "preAllocatedVUs": preAllocatedVUs, "duration": duration, "numStages": len(varr.config.Stages),
		"startTickerPeriod": startTickerPeriod.Duration, "type": varr.config.GetType(),
	***REMOVED***).Debug("Starting executor run...")

	activeVUsWg := &sync.WaitGroup***REMOVED******REMOVED***

	returnedVUs := make(chan struct***REMOVED******REMOVED***)
	startTime, maxDurationCtx, regDurationCtx, cancel := getDurationContexts(parentCtx, duration, gracefulStop)

	defer func() ***REMOVED***
		// Make sure all VUs aren't executing iterations anymore, for the cancel()
		// below to deactivate them.
		<-returnedVUs
		cancel()
		activeVUsWg.Wait()
	***REMOVED***()
	activeVUs := make(chan lib.ActiveVU, maxVUs)
	activeVUsCount := uint64(0)

	returnVU := func(u lib.InitializedVU) ***REMOVED***
		varr.executionState.ReturnVU(u, true)
		activeVUsWg.Done()
	***REMOVED***
	activateVU := func(initVU lib.InitializedVU) lib.ActiveVU ***REMOVED***
		activeVUsWg.Add(1)
		activeVU := initVU.Activate(getVUActivationParams(maxDurationCtx, varr.config.BaseConfig, returnVU))
		varr.executionState.ModCurrentlyActiveVUsCount(+1)
		atomic.AddUint64(&activeVUsCount, 1)
		return activeVU
	***REMOVED***

	remainingUnplannedVUs := maxVUs - preAllocatedVUs
	makeUnplannedVUCh := make(chan struct***REMOVED******REMOVED***)
	defer close(makeUnplannedVUCh)
	go func() ***REMOVED***
		defer close(returnedVUs)
		defer func() ***REMOVED***
			// this is done here as to not have an unplannedVU in the middle of initialization when
			// starting to return activeVUs
			for i := uint64(0); i < atomic.LoadUint64(&activeVUsCount); i++ ***REMOVED***
				<-activeVUs
			***REMOVED***
		***REMOVED***()
		for range makeUnplannedVUCh ***REMOVED***
			varr.logger.Debug("Starting initialization of an unplanned VU...")
			initVU, err := varr.executionState.GetUnplannedVU(maxDurationCtx, varr.logger)
			if err != nil ***REMOVED***
				// TODO figure out how to return it to the Run goroutine
				varr.logger.WithError(err).Error("Error while allocating unplanned VU")
			***REMOVED*** else ***REMOVED***
				varr.logger.Debug("The unplanned VU finished initializing successfully!")
				activeVUs <- activateVU(initVU)
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	// Get the pre-allocated VUs in the local buffer
	for i := int64(0); i < preAllocatedVUs; i++ ***REMOVED***
		initVU, err := varr.executionState.GetPlannedVU(varr.logger, false)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		activeVUs <- activateVU(initVU)
	***REMOVED***

	tickerPeriod := int64(startTickerPeriod.Duration)

	vusFmt := pb.GetFixedLengthIntFormat(maxVUs)
	itersFmt := pb.GetFixedLengthFloatFormat(maxArrivalRatePerSec, 0) + " iters/s"

	progressFn := func() (float64, []string) ***REMOVED***
		currActiveVUs := atomic.LoadUint64(&activeVUsCount)
		currentTickerPeriod := atomic.LoadInt64(&tickerPeriod)
		vusInBuffer := uint64(len(activeVUs))
		progVUs := fmt.Sprintf(vusFmt+"/"+vusFmt+" VUs",
			currActiveVUs-vusInBuffer, currActiveVUs)

		itersPerSec := 0.0
		if currentTickerPeriod > 0 ***REMOVED***
			itersPerSec = float64(time.Second) / float64(currentTickerPeriod)
		***REMOVED***
		progIters := fmt.Sprintf(itersFmt, itersPerSec)

		right := []string***REMOVED***progVUs, duration.String(), progIters***REMOVED***

		spent := time.Since(startTime)
		if spent > duration ***REMOVED***
			return 1, right
		***REMOVED***

		spentDuration := pb.GetFixedLengthDuration(spent, duration)
		progDur := fmt.Sprintf("%s/%s", spentDuration, duration)
		right[1] = progDur

		return math.Min(1, float64(spent)/float64(duration)), right
	***REMOVED***

	varr.progress.Modify(pb.WithProgress(progressFn))
	go trackProgress(parentCtx, maxDurationCtx, regDurationCtx, varr, progressFn)

	regDurationDone := regDurationCtx.Done()
	runIterationBasic := getIterationRunner(varr.executionState, varr.logger)
	runIteration := func(vu lib.ActiveVU) ***REMOVED***
		runIterationBasic(maxDurationCtx, vu)
		activeVUs <- vu
	***REMOVED***

	timer := time.NewTimer(time.Hour)
	start := time.Now()
	ch := make(chan time.Duration, 10) // buffer 10 iteration times ahead
	var prevTime time.Duration
	shownWarning := false
	metricTags := varr.getMetricTags(nil)
	go varr.config.cal(varr.executionState.ExecutionTuple, ch)
	for nextTime := range ch ***REMOVED***
		select ***REMOVED***
		case <-regDurationDone:
			return nil
		default:
		***REMOVED***
		atomic.StoreInt64(&tickerPeriod, int64(nextTime-prevTime))
		prevTime = nextTime
		b := time.Until(start.Add(nextTime))
		if b > 0 ***REMOVED*** // TODO: have a minimal ?
			timer.Reset(b)
			select ***REMOVED***
			case <-timer.C:
			case <-regDurationDone:
				return nil
			***REMOVED***
		***REMOVED***

		select ***REMOVED***
		case vu := <-activeVUs: // ideally, we get the VU from the buffer without any issues
			go runIteration(vu) //TODO: refactor so we dont spin up a goroutine for each iteration
			continue
		default: // no free VUs currently available
		***REMOVED***
		// Since there aren't any free VUs available, consider this iteration
		// dropped - we aren't going to try to recover it, but

		stats.PushIfNotDone(parentCtx, out, stats.Sample***REMOVED***
			Value: 1, Metric: metrics.DroppedIterations,
			Tags: metricTags, Time: time.Now(),
		***REMOVED***)

		// We'll try to start allocating another VU in the background,
		// non-blockingly, if we have remainingUnplannedVUs...
		if remainingUnplannedVUs == 0 ***REMOVED***
			if !shownWarning ***REMOVED***
				varr.logger.Warningf("Insufficient VUs, reached %d active VUs and cannot initialize more", maxVUs)
				shownWarning = true
			***REMOVED***
			continue
		***REMOVED***

		select ***REMOVED***
		case makeUnplannedVUCh <- struct***REMOVED******REMOVED******REMOVED******REMOVED***: // great!
			remainingUnplannedVUs--
		default: // we're already allocating a new VU
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***
