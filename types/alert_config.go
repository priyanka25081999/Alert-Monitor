// types/alert_config.go

package types

import "fmt"

// AlertConfig defines the configuration for alerts
type AlertConfig struct {
    Client             string
    EventType          string
    Config             WindowConfig
    DispatchStrategies []DispatchStrategy
}

// WindowConfig interface
type WindowConfig interface {
    isWindowConfig()
}

// TumblingWindowConfig defines a tumbling window configuration
type TumblingWindowConfig struct {
    Count            int
    WindowSizeInSecs int
}

func (TumblingWindowConfig) isWindowConfig() {}

// SlidingWindowConfig defines a sliding window configuration
type SlidingWindowConfig struct {
    Count            int
    WindowSizeInSecs int
}

func (SlidingWindowConfig) isWindowConfig() {}

// ConfigMessage generates a message for the alert configuration
func (cfg AlertConfig) ConfigMessage() string {
    return fmt.Sprintf("%s %s threshold breached", cfg.Client, cfg.EventType)
}
