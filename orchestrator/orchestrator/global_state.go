package orchestrator

import (
	"context"
	"encoding/json"
	"io"

	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/globe-test/orchestrator/orchMetrics"
	"github.com/APITeamLimited/redis/v9"
	"github.com/sirupsen/logrus"
)

type (
	consoleWriter struct ***REMOVED***
		gs *globalState
	***REMOVED***

	globalState struct ***REMOVED***
		ctx            context.Context
		logger         *logrus.Logger
		client         *redis.Client
		jobId          string
		orchestratorId string
		metricsStore   libOrch.BaseMetricsStore
		status         string
		childJobStates []libOrch.WorkerState
	***REMOVED***
)

var _ libOrch.BaseGlobalState = &globalState***REMOVED******REMOVED***

func (g *globalState) Ctx() context.Context ***REMOVED***
	return g.ctx
***REMOVED***

func (g *globalState) Logger() *logrus.Logger ***REMOVED***
	return g.logger
***REMOVED***

func (g *globalState) Client() *redis.Client ***REMOVED***
	return g.client
***REMOVED***

func (g *globalState) JobId() string ***REMOVED***
	return g.jobId
***REMOVED***

func (g *globalState) OrchestratorId() string ***REMOVED***
	return g.orchestratorId
***REMOVED***

func (g *globalState) MetricsStore() *libOrch.BaseMetricsStore ***REMOVED***
	return &g.metricsStore
***REMOVED***

func (g *globalState) GetStatus() string ***REMOVED***
	return g.status
***REMOVED***

func (g *globalState) SetStatus(status string) ***REMOVED***
	g.status = status
***REMOVED***

func (g *globalState) GetChildJobStates() []libOrch.WorkerState ***REMOVED***
	return g.childJobStates
***REMOVED***

func (g *globalState) SetChildJobState(workerId string, childJobId string, status string) ***REMOVED***
	newState := false
	foundCurrent := false

	for i, childJobState := range g.childJobStates ***REMOVED***
		if childJobState.ChildJobId == childJobId ***REMOVED***
			if g.childJobStates[i].Status != status ***REMOVED***
				newState = true
				g.childJobStates[i].Status = status
			***REMOVED***

			foundCurrent = true
			break
		***REMOVED***
	***REMOVED***

	if !foundCurrent ***REMOVED***
		g.childJobStates = append(g.childJobStates, libOrch.WorkerState***REMOVED***
			WorkerId: workerId,
			Status:   status,
		***REMOVED***)
		newState = true
	***REMOVED***

	// If worker state changes, broadcast a message
	if newState ***REMOVED***
		libOrch.DispatchWorkerMessage(g, workerId, childJobId, status, "STATUS")
	***REMOVED***
***REMOVED***

func NewGlobalState(ctx context.Context, client *redis.Client, jobId string, orchestratorId string) *globalState ***REMOVED***
	gs := &globalState***REMOVED***
		ctx:            ctx,
		client:         client,
		jobId:          jobId,
		orchestratorId: orchestratorId,
		childJobStates: []libOrch.WorkerState***REMOVED******REMOVED***,
	***REMOVED***

	gs.logger = &logrus.Logger***REMOVED***
		Out:       &consoleWriter***REMOVED***gs: gs***REMOVED***,
		Formatter: new(logrus.JSONFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.InfoLevel,
	***REMOVED***

	gs.metricsStore = orchMetrics.NewCachedMetricsStore(gs)

	return gs
***REMOVED***

var _ io.Writer = &consoleWriter***REMOVED******REMOVED***

func (w *consoleWriter) Write(p []byte) (n int, err error) ***REMOVED***
	origLen := len(p)

	// Intercept the write message so can assess log errors parse json
	parsed := make(map[string]interface***REMOVED******REMOVED***)
	if err := json.Unmarshal(p, &parsed); err != nil ***REMOVED***

		return origLen, err
	***REMOVED***

	// Check message level, if error then log error
	if parsed["level"] == "error" ***REMOVED***
		if parsed["error"] != nil ***REMOVED***
			libOrch.HandleStringError(w.gs, parsed["error"].(string))
		***REMOVED*** else ***REMOVED***
			libOrch.HandleStringError(w.gs, parsed["msg"].(string))
		***REMOVED***
		return
	***REMOVED***

	libOrch.DispatchMessage(w.gs, string(p), "STDOUT")

	return origLen, err
***REMOVED***
