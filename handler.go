// Example from https://learn.microsoft.com/en-us/azure/azure-functions/create-first-function-vs-code-other?tabs=go%2Clinux
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/microsoft/ApplicationInsights-Go/appinsights"
)

type InvokeResponse struct {
	Outputs     map[string]interface{}
	Logs        []string
	ReturnValue interface{}
}

var telemetryClient = appinsights.NewTelemetryClient(os.Getenv("APPINSIGHTS_INSTRUMENTATIONKEY"))

func helloHandler(w http.ResponseWriter, r *http.Request) {
	defer func(start time.Time) {
		telemetryClient.TrackRequest(r.Method, r.RequestURI, time.Since(start), "200")
	}(time.Now())

	message := "This HTTP triggered function executed successfully. Pass a name in the query string for a personalized response.\n"
	name := r.URL.Query().Get("name")
	if name != "" {
		message = fmt.Sprintf("Hello, %s. This HTTP triggered function executed successfully.\n", name)
	}

	returnValue := 100
	outputs := make(map[string]interface{})
	outputs["output1"] = message

	invokeResponse := InvokeResponse{outputs, []string{"test log1", "test log2"}, returnValue}

	js, err := json.Marshal(invokeResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func main() {
	defer appinsights.TrackPanic(telemetryClient, false)

	listenAddr := ":8080"
	if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
		listenAddr = ":" + val
	}
	http.HandleFunc("/api/HttpExample", helloHandler)
	log.Printf("About to listen on %s. Go to https://127.0.0.1%s/", listenAddr, listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
