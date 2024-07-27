package tests

import (
	"testing"
	"time"

	"github.com/personal/Alert-Monitor/internal"
	"github.com/personal/Alert-Monitor/types"
)

func TestAlertMonitor(t *testing.T) {
	alertMonitor := internal.NewAlertMonitor()

	// Define a test alert configuration
	alertConfig := types.AlertConfig{
		Client:    "TestClient",
		EventType: "TEST_EVENT",
		Config: types.TumblingWindowConfig{
			Count:            2,
			WindowSizeInSecs: 10,
		},
		DispatchStrategies: []types.DispatchStrategy{
			types.ConsoleDispatch{Message: "test event threshold breached"},
		},
	}

	alertMonitor.RegisterAlertConfig(alertConfig)

	// Record events and verify the alert
	alertMonitor.RecordEvent(types.Event{
		Client:    "TestClient",
		EventType: "TEST_EVENT",
		Timestamp: time.Now(),
	})
	alertMonitor.RecordEvent(types.Event{
		Client:    "TestClient",
		EventType: "TEST_EVENT",
		Timestamp: time.Now(),
	})

	time.Sleep(1 * time.Second)
}
