package types

import (
	"encoding/json"
	"fmt"
)

// AlertConfig defines the configuration for alerts
type AlertConfig struct {
	Client             string             `json:"client"`
	EventType          string             `json:"eventType"`
	Config             WindowConfig       `json:"config"`
	DispatchStrategies []DispatchStrategy `json:"dispatchStrategies"`
}

// UnmarshalJSON custom unmarshals AlertConfig
func (cfg *AlertConfig) UnmarshalJSON(data []byte) error {
	type Alias AlertConfig
	aux := &struct {
		Config             json.RawMessage   `json:"config"`
		DispatchStrategies []json.RawMessage `json:"dispatchStrategyList"`
		*Alias
	}{
		Alias: (*Alias)(cfg),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Unmarshal Config based on type
	var configMap map[string]interface{}
	if err := json.Unmarshal(aux.Config, &configMap); err != nil {
		return err
	}

	switch configMap["type"] {
	case "TUMBLING_WINDOW":
		var twc TumblingWindowConfig
		if err := json.Unmarshal(aux.Config, &twc); err != nil {
			return err
		}
		cfg.Config = twc
	case "SLIDING_WINDOW":
		var swc SlidingWindowConfig
		if err := json.Unmarshal(aux.Config, &swc); err != nil {
			return err
		}
		cfg.Config = swc
	default:
		return fmt.Errorf("unknown config type")
	}

	// Unmarshal DispatchStrategies based on type
	for _, ds := range aux.DispatchStrategies {
		var strategyMap map[string]interface{}
		if err := json.Unmarshal(ds, &strategyMap); err != nil {
			return err
		}

		switch strategyMap["type"] {
		case "CONSOLE":
			var consoleDispatch ConsoleDispatch
			if err := json.Unmarshal(ds, &consoleDispatch); err != nil {
				return err
			}
			cfg.DispatchStrategies = append(cfg.DispatchStrategies, consoleDispatch)
		case "EMAIL":
			var emailDispatch EmailDispatch
			if err := json.Unmarshal(ds, &emailDispatch); err != nil {
				return err
			}
			cfg.DispatchStrategies = append(cfg.DispatchStrategies, emailDispatch)
		default:
			return fmt.Errorf("unknown dispatch strategy type")
		}
	}

	return nil
}

// WindowConfig interface
type WindowConfig interface {
	// marker method - marker interface pattern
	isWindowConfig()
}

// TumblingWindowConfig defines a tumbling window configuration
type TumblingWindowConfig struct {
	Type             string `json:"type"`
	Count            int    `json:"count"`
	WindowSizeInSecs int    `json:"windowSizeInSecs"`
}

// The func (TumblingWindowConfig) isWindowConfig() {} method doesn't need to perform any operations. Its sole purpose is to indicate
// that TumblingWindowConfig satisfies the WindowConfig interface, allowing it to be treated as a WindowConfig type. This pattern is particularly
// useful for categorizing types without enforcing specific behaviors, allowing for flexible and type-safe handling of different configurations.
func (TumblingWindowConfig) isWindowConfig() {}

// SlidingWindowConfig defines a sliding window configuration
type SlidingWindowConfig struct {
	Type             string `json:"type"`
	Count            int    `json:"count"`
	WindowSizeInSecs int    `json:"windowSizeInSecs"`
}

func (SlidingWindowConfig) isWindowConfig() {}

// ConfigMessage generates a message for the alert configuration
func (cfg AlertConfig) ConfigMessage() string {
	return fmt.Sprintf("%s %s threshold breached", cfg.Client, cfg.EventType)
}
