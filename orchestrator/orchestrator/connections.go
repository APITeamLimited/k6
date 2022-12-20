package orchestrator

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"

	"github.com/APITeamLimited/globe-test/lib"
	"github.com/APITeamLimited/globe-test/lib/agent"
	"github.com/APITeamLimited/globe-test/orchestrator/libOrch"
	"github.com/APITeamLimited/redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func tryGetClient(currentIndex int) *libOrch.NamedClient {
	host := lib.GetEnvVariableRaw(fmt.Sprintf("WORKER_%d_HOST", currentIndex), "NONE", true)
	port := lib.GetEnvVariableRaw(fmt.Sprintf("WORKER_%d_PORT", currentIndex), "NONE", true)
	password := lib.GetEnvVariableRaw(fmt.Sprintf("WORKER_%d_PASSWORD", currentIndex), "NONE", true)
	displayName := lib.GetEnvVariableRaw(fmt.Sprintf("WORKER_%d_DISPLAY_NAME", currentIndex), "NONE", true)

	if host == "NONE" || port == "NONE" || password == "NONE" || displayName == "NONE" {
		return nil
	}

	options := &redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       0,
	}

	isSecure := lib.GetEnvVariableRaw(fmt.Sprintf("WORKER_%d_IS_SECURE", currentIndex), "false", true) == "true"

	if isSecure {
		clientCert := lib.GetHexEnvVariable(fmt.Sprintf("WORKER_%d_CERT_HEX", currentIndex), "")
		clientKey := lib.GetHexEnvVariable(fmt.Sprintf("WORKER_%d_KEY_HEX", currentIndex), "")

		cert, err := tls.X509KeyPair(clientCert, clientKey)
		if err != nil {
			panic(fmt.Errorf("error loading orchestrator cert: %s", err))
		}

		// Load CA cert
		caCertPool := x509.NewCertPool()
		caCert := lib.GetHexEnvVariable(fmt.Sprintf("WORKER_%d_CA_CERT_HEX", currentIndex), "")
		ok := caCertPool.AppendCertsFromPEM(caCert)
		if !ok {
			panic("failed to parse root certificate")
		}

		options.TLSConfig = &tls.Config{
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: lib.GetEnvVariable(fmt.Sprintf("WORKER_%d_INSECURE_SKIP_VERIFY", currentIndex), "false") == "true",
			Certificates:       []tls.Certificate{cert},
			RootCAs:            caCertPool,
		}
	}

	client := redis.NewClient(options)

	// Ensure that the client is connected
	if err := client.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}

	return &libOrch.NamedClient{
		Name:   displayName,
		Client: client,
	}
}

func connectWorkerClients(ctx context.Context, standalone bool) libOrch.WorkerClients {
	workerClients := libOrch.WorkerClients{
		Clients: make(map[string]*libOrch.NamedClient),
	}

	if !standalone {
		// Just get a single worker client for the agent
		workerClients.Clients[agent.AgentWorkerName] = &libOrch.NamedClient{
			Name: agent.AgentWorkerName,
			Client: redis.NewClient(&redis.Options{
				Addr: fmt.Sprintf("%s:%s", agent.WorkerRedisHost, agent.WorkerRedisPort),
			}),
		}
		workerClients.DefaultClient = workerClients.Clients[agent.AgentWorkerName]

		return workerClients
	}

	currentIndex := 0

	for {
		namedClient := tryGetClient(currentIndex)

		if namedClient == nil {
			if currentIndex == 0 {
				panic("At least one worker client must be defined")
			}

			break
		}

		if currentIndex == 0 {
			workerClients.DefaultClient = namedClient
		}

		workerClients.Clients[namedClient.Name] = namedClient

		currentIndex++
	}

	return workerClients
}

func getStoreMongoDB(ctx context.Context, standalone bool) *mongo.Database {
	if !standalone {
		return nil
	}

	storeMongoUser := lib.GetEnvVariable("STORE_MONGO_USER", "")
	storeMongoPassword := lib.GetEnvVariable("STORE_MONGO_PASSWORD", "")
	storeMongoHost := lib.GetEnvVariable("STORE_MONGO_HOST", "")
	storeMongoPort := lib.GetEnvVariable("STORE_MONGO_PORT", "")
	storeMongoDatabase := lib.GetEnvVariable("STORE_MONGO_DATABASE", "")

	storeURI := fmt.Sprintf("mongodb://%s:%s@%s:%s", storeMongoUser, storeMongoPassword, storeMongoHost, storeMongoPort)

	client, err := mongo.NewClient(options.Client().ApplyURI(storeURI))
	if err != nil {
		panic(err)
	}

	if err := client.Connect(ctx); err != nil {
		panic(err)
	}

	mongoDB := client.Database(storeMongoDatabase)

	return mongoDB
}
