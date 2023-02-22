package aggregator

import (
	"fmt"
	"regexp"

	"github.com/APITeamLimited/globe-test/lib"
)

// Aggregates intervals and corrects rates
func (aggregator *aggregator) aggregateIntervals(iwsfs []*intervalWithSubfraction, includeGlobalLocation bool) (*Interval, error) {
	// Get all possible locations
	locations := make([]string, 0)
	for _, iwsf := range iwsfs {
		containsLocation := false
		for _, location := range locations {
			if location == iwsf.location {
				containsLocation = true
				break
			}
		}

		if !containsLocation {
			locations = append(locations, iwsf.location)
		}
	}

	if includeGlobalLocation {
		locations = append(locations, "global")
	}

	iwsfByLocation := make([]*intervalWithSubfraction, len(locations))

	for i, location := range locations {
		intevalsToCombine := make([]*intervalWithSubfraction, 0)

		if location == "global" {
			intevalsToCombine = iwsfs
		} else {
			for _, iwsf := range iwsfs {
				if iwsf.location == location {
					intevalsToCombine = append(intevalsToCombine, iwsf)
				}
			}
		}

		combinedInterval, err := combinePeriodIntervals(intevalsToCombine, location)
		if err != nil {
			return nil, err
		}

		iwsfByLocation[i] = combinedInterval
	}

	sinkCount := 0
	for _, iwsf := range iwsfByLocation {
		sinkCount += len(iwsf.interval.Sinks)
	}

	unifiedInterval := Interval{
		Period: iwsfByLocation[0].interval.Period,
		Sinks:  make(map[string]*Sink, sinkCount),
	}

	for _, iwsf := range iwsfByLocation {
		for sinkName, sink := range iwsf.interval.Sinks {
			unifiedInterval.Sinks[sinkName] = sink
		}
	}

	return aggregator.correctIntervalRates(&unifiedInterval), nil
}

// Combines intervals frrom the same period into a single interval
func combinePeriodIntervals(iwsfs []*intervalWithSubfraction, sinkPrefix string) (*intervalWithSubfraction, error) {
	totalFraction := 0.0
	for _, iwsf := range iwsfs {
		totalFraction += iwsf.subFraction
	}

	newIwsf := intervalWithSubfraction{
		subFraction: totalFraction,
		interval: &Interval{
			Period: iwsfs[0].interval.Period,
			// Init sinks with the size of the first interval to minimise allocations
			// Future intervals will still require a reallocation if they have more sinks
			Sinks: make(map[string]*Sink, len(iwsfs[0].interval.Sinks)),
		},
	}

	var aggregatedSubfraction = 0.0

	for _, iwf := range iwsfs {
		for sinkName, sink := range iwf.interval.Sinks {
			prefixedName := fmt.Sprintf("%s::%s", sinkPrefix, sinkName)

			// Add sink if it doesn't exist
			aggregatedSink, ok := newIwsf.interval.Sinks[prefixedName]
			if !ok {
				aggregatedSink = &Sink{
					Labels: make(map[string]float64, len(sink.Labels)),
				}

				newIwsf.interval.Sinks[prefixedName] = aggregatedSink
			}

			for key, value := range sink.Labels {
				aggregatedValue, ok := aggregatedSink.Labels[key]
				if !ok {
					aggregatedValue = 0
				}

				newValue, err := combineSinkValues(key, aggregatedValue, aggregatedSubfraction, value, iwf.subFraction)
				if err != nil {
					// Unknown sink key, just ignore it as not critical
					continue
				}

				aggregatedSink.Labels[key] = newValue
			}
		}

		aggregatedSubfraction += iwf.subFraction
	}

	return &newIwsf, nil
}

var (
	sumKeys         = []string{"count", "rate"}
	meanKeys        = []string{"avg", "med"}
	percentileRegex = regexp.MustCompile(`p\([1-9][0-9]?|100\)`)
)

// Combines sink values from multiple intervals
func combineSinkValues(key string, value1, fraction1, value2, fraction2 float64) (float64, error) {
	// If the key is a sum key, then we can just add the values
	if lib.StringInSlice(sumKeys, key) {
		return value1 + value2, nil
	} else if lib.StringInSlice(meanKeys, key) || percentileRegex.MatchString(key) {
		return ((value1 * fraction1) + (value2 * fraction2)) / (fraction1 + fraction2), nil
	} else if key == "max" {
		if value1 > value2 {
			return value1, nil
		} else {
			return value2, nil
		}
	} else if key == "min" {
		if value1 < value2 {
			return value1, nil
		} else {
			return value2, nil
		}
	} else {
		// Unknown key just ignore it
		return 0, fmt.Errorf("Unknown key: '%s'" + key)
	}
}

// Correct rates, look through and find rates in unified sinks, if found, look up in the previous interval
// then calculate rate based on the previous interval's value and the current interval's value
func (aggregator *aggregator) correctIntervalRates(interval *Interval) *Interval {
	defer func() {
		// Add interval to previous intervals and remove the first added interval if there are more than intervalMaxLagPeriods
		aggregator.previousIntervals = append(aggregator.previousIntervals, interval)
		if len(aggregator.previousIntervals) > intervalMaxLagPeriods {
			aggregator.previousIntervals = aggregator.previousIntervals[1:]
		}
	}()

	// Determine if period exists in previous intervals
	var previousInterval *Interval
	previousPeriod := interval.Period - 1

	for _, iwsf := range aggregator.previousIntervals {
		if iwsf.Period == previousPeriod {
			previousInterval = iwsf
			break
		}
	}

	if previousInterval == nil {
		// Rates are best guess if the period doesn't exist in previous intervals
		return interval
	}

	// Correct rates in sinks
	for sinkName, sink := range interval.Sinks {
		for label := range sink.Labels {
			if label == "rate" {
				previousSink, ok := previousInterval.Sinks[sinkName]
				if !ok {
					continue
				}

				previousCount, ok := previousSink.Labels["count"]
				if !ok {
					continue
				}

				currentCount, ok := sink.Labels["count"]
				if !ok {
					continue
				}

				// These periods have a second difference, so no need to divide by the period
				sink.Labels[label] = (currentCount - previousCount)
			}
		}
	}

	return interval
}
