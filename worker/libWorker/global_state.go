package libWorker

import (
	"context"

	"github.com/APITeamLimited/redis/v9"
)

type WorkerState struct ***REMOVED***
	Status   string `json:"status"`
	WorkerId string `json:"workerId"`
***REMOVED***

type BaseGlobalState interface ***REMOVED***
	Ctx() context.Context
	// The orchestrator client
	Client() *redis.Client
	JobId() string
	ChildJobId() string
	WorkerId() string
	GetWorkerStatus() string
	SetWorkerStatus(status string)
***REMOVED***
