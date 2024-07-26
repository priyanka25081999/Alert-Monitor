package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/personal/Alert-Monitor/internal"
	"github.com/personal/Alert-Monitor/types"
)

// Define global alert monitor
var alertMonitor *internal.AlertMonitor

func init() {
    alertMonitor = internal.NewAlertMonitor()
}

// addAlertConfigHandler handles POST requests to add alert configurations
func addAlertConfigHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var config types.AlertConfig
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Received config: %+v\n", config)

	alertMonitor.RegisterAlertConfig(config)

	response := map[string]string{
		"message":   "Alert configuration added successfully",
		"client":    config.Client,
		"eventType": config.EventType,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Handler to simulate an event
func recordEventHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    var event types.Event
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&event); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Record the event and check for alerts
    alertResult, alertTriggered := alertMonitor.RecordEvent(event)

	fmt.Println("AlertResult: ", alertResult, " AlertTriggered: ", alertTriggered)

    response := map[string]interface{}{
        "message": "Event recorded successfully",
        "client":  event.Client,
        "eventType": event.EventType,
        "alert": alertResult,
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
}

func main() {
    http.HandleFunc("/add-config", addAlertConfigHandler)
    http.HandleFunc("/record-event", recordEventHandler)

    log.Println("Server starting on port 8080...")
    log.Fatal(http.ListenAndServe(":8080", nil))
}