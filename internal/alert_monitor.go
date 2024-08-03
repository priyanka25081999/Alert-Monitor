package internal

import (
	"fmt"
	"sync"
	"time"

	"github.com/personal/Alert-Monitor/types"
)

type AlertMonitor struct {
	alertConfigs []types.AlertConfig
	eventLogs    map[string][]time.Time
	mu           sync.Mutex
}

// The purpose of this function is to create and return a new AlertMonitor instance with its internal state (eventLogs)
// properly initialized. This initializes the eventLogs field of the AlertMonitor struct. eventLogs is a map where
// the keys are strings, and the values are slices of time.Time objects.
// The make function is used to create the map, ensuring that it is properly initialized before being used.
func NewAlertMonitor() *AlertMonitor {
	return &AlertMonitor{
		eventLogs: make(map[string][]time.Time),
	}
}

func (am *AlertMonitor) RegisterAlertConfig(config types.AlertConfig) {
	// This locks the mu mutex in the AlertMonitor struct, ensuring that only one goroutine can modify the shared resource
	// (alertConfigs slice) at a time. This is crucial for thread-safe operations in concurrent environments.
	am.mu.Lock()
	defer am.mu.Unlock()

	// Appends the new alert configuration (config) to the alertConfigs slice of the AlertMonitor instance.
	// The alertConfigs slice stores all the alert configurations that the AlertMonitor is managing.
	am.alertConfigs = append(am.alertConfigs, config)
}

// RecordEvent processes an event and returns an alert message if triggered
func (am *AlertMonitor) RecordEvent(event types.Event) (string, bool) {
	am.mu.Lock()
	defer am.mu.Unlock()

	// Generates a unique key for the event based on the client and event type, used to categorize events in the eventLogs.
	key := fmt.Sprintf("%s-%s", event.Client, event.EventType)

	// Adds the event's timestamp to the list of events under the generated key in eventLogs.
	am.eventLogs[key] = append(am.eventLogs[key], event.Timestamp)

	// Prepares variables to store the alert message and a boolean indicating whether an alert was triggered.
	var alertMessage string
	alertTriggered := false

	// Iterates through all alert configurations and checks if the event matches the configuration's client and event type.
	// If a match is found, it calls checkThreshold to see if the threshold is met. If the threshold is met,
	// an alert message is set, alertTriggered is set to true, and dispatchAlert is called to dispatch the alert.
	// The loop breaks after the first triggered alert.
	for _, config := range am.alertConfigs {
		if config.Client == event.Client && config.EventType == event.EventType {
			if am.checkThreshold(config, key) {
				alertMessage = fmt.Sprintf("Alert triggered for %s with event type %s", event.Client, event.EventType)
				alertTriggered = true
				am.dispatchAlert(config)
				break
			}
		}
	}

	return alertMessage, alertTriggered
}

func (am *AlertMonitor) checkThreshold(config types.AlertConfig, key string) bool {
	// Retrieves the list of events for the specified key and the current time.
	events := am.eventLogs[key]
	now := time.Now()

	// Determines the type of alert configuration (TumblingWindowConfig or SlidingWindowConfig) and applies the corresponding logic.
	switch cfg := config.Config.(type) {
	// Window Start: Calculates the start of the current time window.
	// Count Events: Counts the events that occurred after the window start.
	// Returns Threshold Check: Returns true if the event count meets or exceeds the threshold.
	case types.TumblingWindowConfig:
		// If now is 2024-08-12 10:23:45.678, truncating to the nearest hour would give 2024-08-12 10:00:00.000.
		windowStart := now.Truncate(time.Duration(cfg.WindowSizeInSecs) * time.Second)
		count := 0
		for _, eventTime := range events {
			if eventTime.After(windowStart) {
				count++
			}
		}
		return count >= cfg.Count

	case types.SlidingWindowConfig:
		// If now is 2024-08-12 10:23:45.678, adding 2 hours would give 2024-08-12 12:23:45.678.
		windowStart := now.Add(-time.Duration(cfg.WindowSizeInSecs) * time.Second)
		count := 0

		for _, eventTime := range events {
			if eventTime.After(windowStart) {
				count++
			}
		}
		return count >= cfg.Count
	}
	return false
}

// This method handles the dispatch of alerts based on the defined dispatch strategies.
func (am *AlertMonitor) dispatchAlert(config types.AlertConfig) {
	// Iterates over all the dispatch strategies defined in the alert configuration and calls the Dispatch method
	// on each strategy with the alert's configuration message.
	for _, strategy := range config.DispatchStrategies {
		strategy.Dispatch(config.ConfigMessage())
	}
}
