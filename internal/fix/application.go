package fix

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/quickfixgo/quickfix"
)

// -----------------------------------------------------------------------------
// APPLICATION LOGIC
// -----------------------------------------------------------------------------

// FIXApplication is our custom struct that implements the quickfix.Application interface.
// The FIX engine will call methods on this struct when events happen.
type Application struct {
	SessionID   quickfix.SessionID // Stores the ID of the current connection
	LastClOrdID string             // Stores the ID of the last order sent (to allow cancelling it)
}

// NewApplication creates a new instance of our FIX application logic.
func NewApplication() *Application {
	return &Application{}
}

// OnCreate is called when a FIX session is created (before connecting).
// We use this to save the SessionID for later use.
func (a *Application) OnCreate(sessionID quickfix.SessionID) {
	a.SessionID = sessionID
	log.Printf("Setup: Session created for %s", sessionID)
}

// OnLogon is called when we successfully connect and authenticate with the server.
// This is where we start our business logic (sending orders).
func (a *Application) OnLogon(sessionID quickfix.SessionID) {
	log.Printf("Connection: Logged On to %s", sessionID)

	// We send an order in a separate thread (goroutine) so we don't block the
	// message processing loop. We wait 1 second to ensure everything is ready.
	go func() {
		time.Sleep(1 * time.Second)
		a.sendOrder()

		// Wait 5 seconds and then try to cancel it
		time.Sleep(5 * time.Second)
		a.sendCancelOrder()
	}()
}

// OnLogout is called when we disconnect from the server.
func (a *Application) OnLogout(sessionID quickfix.SessionID) {
	log.Printf("Connection: Logged Out from %s", sessionID)
}

// ToAdmin is called before sending an administrative message (Logon, Heartbeat).
// You can use this to add Username/Password to the Logon message if needed.
func (a *Application) ToAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) {
	msgType, err := msg.Header.GetString(TagMsgType)

	// If it's a Logon message (MsgType="A"), we log it.
	if err == nil && msgType == MsgTypeLogon {
		log.Printf("Admin: Sending Logon message...")
	}
}

// ToApp is called before sending an application message (like an Order).
// It's a last chance to modify the message or log it.
func (a *Application) ToApp(msg *quickfix.Message, sessionID quickfix.SessionID) error {
	log.Printf("Outbound: Sending Application Message: %s", msg.String())
	return nil
}

// FromAdmin is called when we receive an administrative message (Heartbeat, Logon confirmation).
func (a *Application) FromAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	// 1. Get the Message Type to decide how to handle it
	msgType, err := msg.Header.GetString(TagMsgType)
	if err != nil {
		return nil
	}

	// 2. Switch based on the message type
	switch msgType {
	case MsgTypeReject:
		// Tag 35 = 3: The server is rejecting a session-level message (like an unauthorized order)
		printReject(msg)
	case MsgTypeHeartBeat:
		// Tag 35 = 0: Heartbeat. We do nothing and don't log it to avoid clutter.
		return nil
	default:
		log.Printf("Admin: Received Admin Message: %s", msg.String())
	}
	return nil
}

// FromApp is called when we receive a business message (ExecutionReport, etc.).
// This is the core handler for incoming data.
func (a *Application) FromApp(msg *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	// 1. Get the Message Type to decide how to handle it
	msgType, err := msg.Header.GetString(TagMsgType)
	if err != nil {
		log.Printf("Error: Could not read MsgType")
		return nil
	}

	// 2. Switch based on the message type
	switch msgType {
	case MsgTypeExecutionReport:
		// Tag 35 = 8: The server is updating us on an order
		printExecutionReport(msg)
	case MsgTypeOrderCancelReject:
		// Tag 35 = 9: The server rejected a cancel request
		printOrderCancelReject(msg)
	default:
		// Any other message we don't explicitly handle
		log.Printf("Info: Received generic message: %s", msg.String())
	}
	return nil
}

// sendOrder constructs and sends a "New Order Single" message to the server.
func (a *Application) sendOrder() {
	// 1. Create a new empty message
	msg := quickfix.NewMessage()

	// 2. Set the Header fields (Required)
	//    MsgType 'D' tells the server this is a New Order Single
	msg.Header.SetField(TagMsgType, quickfix.FIXString(MsgTypeNewOrderSingle))

	// 3. Set the Body fields (The details of the order)

	// Account: The trading account ID (Sourced from env or default)
	account := os.Getenv("ACCOUNT_NUMBER")
	if account == "" {
		account = "FIX-TEST-ACCOUNT-1"
	}
	msg.Body.SetField(TagAccount, quickfix.FIXString(account))

	// ClOrdID: A unique ID WE generate to track this order
	clOrdID := fmt.Sprintf("ORD-%d", time.Now().UnixNano())
	msg.Body.SetField(TagClOrdID, quickfix.FIXString(clOrdID))

	// Store the ID so we can cancel it later
	a.LastClOrdID = clOrdID

	// Symbol: What we want to trade
	msg.Body.SetField(TagSymbol, quickfix.FIXString("AAPL"))

	// Side: Buy (1)
	msg.Body.SetField(TagSide, quickfix.FIXString(SideBuy))

	// TransactTime: Current UTC timestamp
	msg.Body.SetField(TagTransactTime, quickfix.FIXString(time.Now().Format("20060102-15:04:05.000")))

	// OrderQty: How much we want
	msg.Body.SetField(TagOrderQty, quickfix.FIXString("100"))

	// OrdType: Market (1) - execute immediately at best price
	msg.Body.SetField(TagOrdType, quickfix.FIXString(OrdTypeMarket))

	// TimeInForce: Day (0) - order expires at end of session
	msg.Body.SetField(TagTimeInForce, quickfix.FIXString(TimeInForceDay))

	// 4. Send the message to the session
	log.Printf("Action: Sending New Order Single (ID: %s)...", clOrdID)
	err := quickfix.SendToTarget(msg, a.SessionID)
	if err != nil {
		log.Printf("Error: Failed to send order: %s", err)
	}
}

// sendCancelOrder constructs and sends an "Order Cancel Request" message.
func (a *Application) sendCancelOrder() {
	if a.LastClOrdID == "" {
		log.Printf("Action: No previous order to cancel.")
		return
	}

	// 1. Create a new empty message
	msg := quickfix.NewMessage()

	// 2. Set the Header fields
	//    MsgType 'F' tells the server this is an Order Cancel Request
	msg.Header.SetField(TagMsgType, quickfix.FIXString(MsgTypeOrderCancelRequest))

	// 3. Set the Body fields

	// OrigClOrdID: The ID of the order we want to cancel (the one we stored earlier)
	msg.Body.SetField(TagOrigClOrdID, quickfix.FIXString(a.LastClOrdID))

	// ClOrdID: A NEW unique ID for the cancel request itself
	cancelClOrdID := fmt.Sprintf("CAN-%d", time.Now().UnixNano())
	msg.Body.SetField(TagClOrdID, quickfix.FIXString(cancelClOrdID))

	// Account: Same account as the original order
	account := os.Getenv("ACCOUNT_NUMBER")
	if account == "" {
		account = "FIX-TEST-ACCOUNT-1"
	}
	msg.Body.SetField(TagAccount, quickfix.FIXString(account))

	// Symbol: Must match the original order
	msg.Body.SetField(TagSymbol, quickfix.FIXString("AAPL"))

	// Side: Must match the original order
	msg.Body.SetField(TagSide, quickfix.FIXString(SideBuy))

	// TransactTime: Current UTC timestamp
	msg.Body.SetField(TagTransactTime, quickfix.FIXString(time.Now().Format("20060102-15:04:05.000")))

	// 4. Send the message
	log.Printf("Action: Sending Order Cancel Request (Target ID: %s, Cancel ID: %s)...", a.LastClOrdID, cancelClOrdID)
	err := quickfix.SendToTarget(msg, a.SessionID)
	if err != nil {
		log.Printf("Error: Failed to send cancel request: %s", err)
	}
}
