package main

import (
	"time"

	"github.com/personal/Alert-Monitor/internal"
	"github.com/personal/Alert-Monitor/types"
)

func main() {
	alertMonitor := internal.NewAlertMonitor()

	// Define alert configurations
	alertConfigList := []types.AlertConfig{
		{
            Client:    "X",
            EventType: "PAYMENT_EXCEPTION",
            Config: types.TumblingWindowConfig{
                Count:            10,
                WindowSizeInSecs: 10,
            },
            DispatchStrategies: []types.DispatchStrategy{
                types.ConsoleDispatch{Message: "issue in payment"},
                types.EmailDispatch{Subject: "payment exception threshold breached"},
            },
        },
		{
			Client:     "X",
			EventType:  "USERSERVICE_EXCEPTION",
			Config:      types.SlidingWindowConfig{Count: 10, WindowSizeInSecs: 10},
			DispatchStrategies: []types.DispatchStrategy{
				types.ConsoleDispatch{Message: "issue in user service"},
			},
		},
	}

	// Register configurations
	for _, config := range alertConfigList {
		alertMonitor.RegisterAlertConfig(config)
	}

	// Simulate events
	go func() {
		for i := 0; i < 20; i++ {
			alertMonitor.RecordEvent(types.Event{
				Client:    "X",
				EventType: "PAYMENT_EXCEPTION",
				Timestamp: time.Now(),
			})
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		for i := 0; i < 20; i++ {
			alertMonitor.RecordEvent(types.Event{
				Client:    "X",
				EventType: "USERSERVICE_EXCEPTION",
				Timestamp: time.Now(),
			})
			time.Sleep(2 * time.Second)
		}
	}()

	// Keep the program running to allow for event processing
	select {}
}
