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

	outputs := make(map[string]interface{})
	outputs["myMessage"] = message
	outputs["object"] = map[string]interface{}{
		"somekey1": "value1",
		"somekey2": "value2",
	}
	headers := make(map[string]interface{})
	headers["header1"] = "header1Val"
	headers["header2"] = "header2Val"

	res := make(map[string]interface{})
	res["statusCode"] = "201"
	res["body"] = "my world"
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
