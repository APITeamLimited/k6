package function_auth_client

import (
	"context"
	"encoding/hex"
	"sync"

	functions "cloud.google.com/go/functions/apiv2"
	"github.com/APITeamLimited/globe-test/lib"
	"google.golang.org/api/option"
)

type FunctionAuthClient struct {
	functionClient     *functions.FunctionClient
	liveFunctions      []lib.LiveFunction
	liveFunctionsMutex sync.Mutex
	ctx                context.Context
}

var _ = lib.FunctionAuthClient(&FunctionAuthClient{})

func CreateFunctionAuthClient(ctx context.Context, funcMode bool) *FunctionAuthClient {
	if !funcMode {
		return nil
	}

	serviceAccountHex := lib.GetEnvVariable("SERVICE_ACCOUNT_KEY_HEX", "")

	// Convert hex to bytes
	serviceAccount, err := hex.DecodeString(serviceAccountHex)
	if err != nil {
		panic(err)
	}

	functionClient, err := functions.NewFunctionClient(ctx, option.WithCredentialsJSON([]byte(serviceAccount)))
	if err != nil {
		panic(err)
	}

	functionAuthClient := &FunctionAuthClient{
		functionClient:     functionClient,
		ctx:                ctx,
		liveFunctions:      []lib.LiveFunction{},
		liveFunctionsMutex: sync.Mutex{},
	}

	functionAuthClient.startAutoRefreshLiveFunctions()

	return functionAuthClient
}
