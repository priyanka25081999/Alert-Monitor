package types

import "fmt"

// DispatchStrategy interface
type DispatchStrategy interface {
    Dispatch(alertMessage string)
}

// ConsoleDispatch implements DispatchStrategy for console output
type ConsoleDispatch struct {
	Type    string `json:"type"`
    Message string `json:"message"`
}

func (c ConsoleDispatch) Dispatch(alertMessage string) {
    fmt.Printf("[WARN] Alert: `%s`\n", c.Message)
}

// EmailDispatch implements DispatchStrategy for email output
type EmailDispatch struct {
	Type    string `json:"type"`
    Subject string `json:"subject"`
}

func (e EmailDispatch) Dispatch(alertMessage string) {
    fmt.Printf("[INFO] AlertingService: Dispatching an Email with subject: `%s`\n", e.Subject)
}
