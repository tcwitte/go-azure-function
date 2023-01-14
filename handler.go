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

type InvokeRequest struct {
	Data     map[string]interface{}
	Metadata map[string]interface{}
}

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

	var invokeReq InvokeRequest
	d := json.NewDecoder(r.Body)
	decodeErr := d.Decode(&invokeReq)
	if decodeErr != nil {
		http.Error(w, decodeErr.Error(), http.StatusBadRequest)
		return
	}

	outputs := make(map[string]interface{})
	outputs["document"] = map[string]interface{}{
		"data": invokeReq.Data,
		"azureFunctionsInvocationId": r.Header.Get("X-Azure-Functions-InvocationId"),
	}
	headers := make(map[string]interface{})
	headers["Content-Type"] = "text/plain"

	res := make(map[string]interface{})
	res["statusCode"] = "201"
	res["body"] = invokeReq.Data
	res["headers"] = headers
	outputs["res"] = res
	invokeResponse := InvokeResponse{outputs, []string{"test log1", "test log2"}, "Hello,World"}

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
	http.HandleFunc("/HttpExample", helloHandler)
	log.Printf("About to listen on %s. Go to https://127.0.0.1%s/", listenAddr, listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
