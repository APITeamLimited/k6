package orchMetrics

import (
	"time"

	"github.com/APITeamLimited/globe-test/worker/output/globetest"
)

func calculateMean(wrappedEnvelopes []*wrappedEnvelope, envelope *wrappedEnvelope, currentTime time.Time, vuCount int) []*wrappedEnvelope ***REMOVED***
	// Zero weighting has no effect on the rate
	if vuCount == 0 ***REMOVED***
		return wrappedEnvelopes
	***REMOVED***

	// Check if wrapped envelope exists in the slice
	metricName := envelope.Metric.Name
	accEnvelope := &wrappedEnvelope***REMOVED******REMOVED***

	for _, wrappedEnvelopeCurrent := range wrappedEnvelopes ***REMOVED***
		if wrappedEnvelopeCurrent.Metric.Name == metricName ***REMOVED***
			accEnvelope = wrappedEnvelopeCurrent
			break
		***REMOVED***
	***REMOVED***

	// If wrapped envelope does not exist, create a new one
	if accEnvelope.Metric == nil ***REMOVED***
		accEnvelope := wrappedEnvelope***REMOVED***
			SampleEnvelope: globetest.SampleEnvelope***REMOVED***
				Type: "Point",
				Data: globetest.SampleData***REMOVED***
					Time:  currentTime,
					Value: 0,
				***REMOVED***,
				Metric: envelope.Metric,
			***REMOVED***,
			workerId: envelope.workerId,
			location: envelope.location,
		***REMOVED***

		wrappedEnvelopes = append(wrappedEnvelopes, &accEnvelope)
	***REMOVED***

	// Rate tracks the percentage of values that are non-zero
	// Perform weighted average

	weightingOld := accEnvelope.weighting
	weightingEnvelope := envelope.weighting

	rateOld := accEnvelope.Data.Value
	rateEnvelope := envelope.Data.Value

	weightingNew := weightingOld + weightingEnvelope

	// Cannot divide by zero
	if weightingNew == 0 ***REMOVED***
		return wrappedEnvelopes
	***REMOVED***

	rateNew := (rateOld*float64(weightingOld) + rateEnvelope*float64(weightingEnvelope)) / float64(weightingNew)

	accEnvelope.weighting = weightingNew
	accEnvelope.Data.Value = rateNew

	return wrappedEnvelopes
***REMOVED***
