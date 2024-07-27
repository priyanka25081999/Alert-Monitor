package tests

import (
	"testing"

	"github.com/personal/Alert-Monitor/types"
)

func TestConsoleDispatch(t *testing.T) {
	dispatch := types.ConsoleDispatch{Message: "test console dispatch"}
	dispatch.Dispatch("test console dispatch")
}

func TestEmailDispatch(t *testing.T) {
	dispatch := types.EmailDispatch{Subject: "test email dispatch"}
	dispatch.Dispatch("test email dispatch")
}
