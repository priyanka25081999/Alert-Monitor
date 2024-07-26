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

func NewAlertMonitor() *AlertMonitor {
	return &AlertMonitor{
		eventLogs: make(map[string][]time.Time),
	}
}

func (am *AlertMonitor) RegisterAlertConfig(config types.AlertConfig) {
	am.mu.Lock()
    defer am.mu.Unlock()
	am.alertConfigs = append(am.alertConfigs, config)
}

// RecordEvent processes an event and returns an alert message if triggered
func (am *AlertMonitor) RecordEvent(event types.Event) (string, bool) {
	am.mu.Lock()
    defer am.mu.Unlock()

    key := fmt.Sprintf("%s-%s", event.Client, event.EventType)
    am.eventLogs[key] = append(am.eventLogs[key], event.Timestamp)

    var alertMessage string
    alertTriggered := false

    for _, config := range am.alertConfigs {
        if config.Client == event.Client && config.EventType == event.EventType {
            if am.checkThreshold(config, key) {
                alertMessage = fmt.Sprintf("Alert triggered for %s with event type %s", event.Client, event.EventType)
                alertTriggered = true
                am.dispatchAlert(config)
                break // Assuming you only need to dispatch the alert once per event
            }
        }
    }

    return alertMessage, alertTriggered
}

func (am *AlertMonitor) checkThreshold(config types.AlertConfig, key string) bool {
	events := am.eventLogs[key]
	now := time.Now()

	switch cfg := config.Config.(type) {
	case types.TumblingWindowConfig:
		windowStart := now.Truncate(time.Duration(cfg.WindowSizeInSecs) * time.Second)
		count := 0
		for _, eventTime := range events {
			if eventTime.After(windowStart) {
				count++
			}
		}
		return count >= cfg.Count

	case types.SlidingWindowConfig:
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

func (am *AlertMonitor) dispatchAlert(config types.AlertConfig) {
	for _, strategy := range config.DispatchStrategies {
		strategy.Dispatch(config.ConfigMessage())
	}
}