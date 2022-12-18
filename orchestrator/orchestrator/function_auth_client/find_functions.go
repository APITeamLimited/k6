package function_auth_client

import (
	"fmt"
	"time"

	"cloud.google.com/go/functions/apiv2/functionspb"
	"github.com/APITeamLimited/globe-test/lib"
)

func (config *FunctionAuthClient) Functions() []lib.LiveFunction {
	config.liveFunctionsMutex.Lock()
	defer config.liveFunctionsMutex.Unlock()
	return config.liveFunctions
}

func (config *FunctionAuthClient) startAutoRefreshLiveFunctions() {
	// Get a new token straight away
	err := config.updateLiveFunctions()
	if err != nil {
		fmt.Println(err)
	}

	go func() {
		for {
			time.Sleep(time.Second * 5)
			err := config.updateLiveFunctions()
			if err != nil {
				fmt.Println(err)
			}
		}
	}()
}

// Queries the google cloud functions API to get the list of function URLs
func (config *FunctionAuthClient) updateLiveFunctions() error {
	functionsIterator := config.functionClient.ListFunctions(config.ctx, &functionspb.ListFunctionsRequest{
		Parent: fmt.Sprintf("projects/apiteam-%s/locations/-", lib.GetEnvVariable("ENVIRONMENT", "")),
	})

	var functions []lib.LiveFunction

	config.liveFunctionsMutex.Lock()
	defer config.liveFunctionsMutex.Unlock()

	for {
		function, err := functionsIterator.Next()
		if err != nil {
			break
		}

		functions = append(functions, lib.LiveFunction{
			Location: parseLocation(function.Name),
			Uri:      function.ServiceConfig.Uri,
			State:    function.State,
		})
	}

	config.liveFunctions = functions

	return nil
}

func parseLocation(functionName string) string {
	// worker-{location} format
	return functionName[7:]
}
