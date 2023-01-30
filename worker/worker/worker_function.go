package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/APITeamLimited/globe-test/lib"
	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func RunWorkerServer() {
	port := lib.GetEnvVariableRaw("WORKER_SERVER_PORT", "8090", true)
	fmt.Printf("Starting worker server on port %s\n", port)

	http.HandleFunc("/", runWorkerFunction)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func runWorkerFunction(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	workerId := uuid.NewString()

	upgrader := websocket.Upgrader{}

	// Upgrade the connection to a websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Error upgrading connection to websocket: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	var childJob *libOrch.ChildJob

	connWriteMutex := &sync.Mutex{}
	connReadMutex := &sync.Mutex{}

	for {
		eventMessage := lib.EventMessage{}

		// Read the message from the connection
		connReadMutex.Lock()
		// Okay here as at the start
		err := conn.ReadJSON(&eventMessage)
		connReadMutex.Unlock()
		if err != nil {
			// If websocket is closed, return
			fmt.Printf("Error reading message from connection: %s", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		if eventMessage.Variant == lib.CHILD_JOB_INFO {
			err := json.Unmarshal([]byte(eventMessage.Data), &childJob)
			if err != nil {
				fmt.Printf("Error unmarshalling child job: %s", err.Error())
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}
			break
		}
	}

	creditsClient := lib.GetCreditsClient(true)

	fmt.Printf("Worker %s executing child job %s\n", workerId, childJob.ChildJobId)

	successfullExecution := handleExecution(ctx, conn, childJob, workerId, creditsClient, true, connReadMutex, connWriteMutex)

	// Close the connection gracefully

	connWriteMutex.Lock()
	conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	connWriteMutex.Unlock()

	fmt.Printf("Worker %s finished executing child job %s with success: %t\n", workerId, childJob.ChildJobId, successfullExecution)
}
