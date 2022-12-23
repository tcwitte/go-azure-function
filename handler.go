// Example from https://learn.microsoft.com/en-us/azure/azure-functions/create-first-function-vs-code-other?tabs=go%2Clinux
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/microsoft/ApplicationInsights-Go/appinsights"
)

var telemetryClient = appinsights.NewTelemetryClient(os.Getenv("APPINSIGHTS_INSTRUMENTATIONKEY"))
var cosmosClient *azcosmos.Client

func init() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}
	// Create a CosmosDB client
	cosmosUri := fmt.Sprintf("https://%s.documents.azure.com:443/", os.Getenv("CosmosAccount"))
	cosmosClient, err = azcosmos.NewClient(cosmosUri, cred, nil)
	if err != nil {
		log.Fatal("Failed to create Azure Cosmos DB client: ", err)
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	defer func(start time.Time) {
		telemetryClient.TrackRequest(r.Method, r.RequestURI, time.Since(start), "200")
	}(time.Now())

	message := "This HTTP triggered function executed successfully. Pass a name in the query string for a personalized response.\n"
	name := r.URL.Query().Get("name")
	if name != "" {
		message = fmt.Sprintf("Hello, %s. This HTTP triggered function executed successfully.\n", name)
	}
	fmt.Fprint(w, message)
}

func main() {
	defer appinsights.TrackPanic(telemetryClient, false)
	createDatabase("myfirstdatabase")

	listenAddr := ":8080"
	if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
		listenAddr = ":" + val
	}
	http.HandleFunc("/api/HttpExample", helloHandler)
	log.Printf("About to listen on %s. Go to https://127.0.0.1%s/", listenAddr, listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}

func createDatabase(databaseName string) error {
	databaseProperties := azcosmos.DatabaseProperties{ID: databaseName}

	ctx := context.TODO()
	_, err := cosmosClient.CreateDatabase(ctx, databaseProperties, nil)
	return err
}
