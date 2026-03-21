package fix

import (
	"testing"

	"github.com/quickfixgo/quickfix"
)

func TestNewApplication(t *testing.T) {
	app := NewApplication()
	if app == nil {
		t.Fatal("NewApplication returned nil")
	}
}

// Example of how you might test message handling in the future
func TestFromApp_ExecutionReport(t *testing.T) {
	app := NewApplication()
	sessionID := quickfix.SessionID{}

	// Create a dummy Execution Report message
	msg := quickfix.NewMessage()
	msg.Header.SetField(TagMsgType, quickfix.FIXString(MsgTypeExecutionReport))
	msg.Body.SetField(TagExecType, quickfix.FIXString(ExecTypeNew))
	msg.Body.SetField(TagOrdStatus, quickfix.FIXString("0")) // New

	// Simulate receiving it
	err := app.FromApp(msg, sessionID)
	if err != nil {
		t.Errorf("FromApp returned error: %v", err)
	}

	// Since FromApp just prints to stdout, we can't easily assert the output here without capturing stdout.
	// But passing without error is a good start.
}
