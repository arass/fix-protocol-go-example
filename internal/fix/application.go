package fix

import (
	"fmt"
	"log"
	"time"

	"github.com/quickfixgo/quickfix"
)

// -----------------------------------------------------------------------------
// APPLICATION LOGIC
// -----------------------------------------------------------------------------

// FIXApplication is our custom struct that implements the quickfix.Application interface.
// The FIX engine will call methods on this struct when events happen.
type Application struct {
	SessionID quickfix.SessionID // Stores the ID of the current connection
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

// -----------------------------------------------------------------------------
// HELPER FUNCTIONS
// -----------------------------------------------------------------------------

// printExecutionReport parses and prints the details of an Execution Report.
func printExecutionReport(msg *quickfix.Message) {
	// Extract fields from the message body
	orderID, _ := msg.Body.GetString(TagOrderID)
	execType, _ := msg.Body.GetString(TagExecType)
	ordStatus, _ := msg.Body.GetString(TagOrdStatus)
	filledQty, _ := msg.Body.GetString(TagCumQty)
	leavesQty, _ := msg.Body.GetString(TagLeavesQty)
	avgPx, _ := msg.Body.GetString(TagAvgPx)

	// Print a clean, human-readable table
	fmt.Printf("\n--- [Execution Report] ---\n")
	fmt.Printf("Server Order ID: %s\n", orderID)
	fmt.Printf("Event Type:      %s\n", translateExecType(execType))
	fmt.Printf("Current Status:  %s\n", translateOrdStatus(ordStatus))
	fmt.Printf("Filled Qty:      %s\n", filledQty)
	fmt.Printf("Remaining Qty:   %s\n", leavesQty)
	fmt.Printf("Average Price:   %s\n", avgPx)
	fmt.Printf("--------------------------\n")
}

// printOrderCancelReject parses and prints the details of a Cancel Reject.
func printOrderCancelReject(msg *quickfix.Message) {
	orderID, _ := msg.Body.GetString(TagOrderID)
	origClOrdID, _ := msg.Body.GetString(TagOrigClOrdID)
	ordStatus, _ := msg.Body.GetString(TagOrdStatus)
	text, _ := msg.Body.GetString(TagText)

	fmt.Printf("\n--- [Order Cancel Reject] ---\n")
	fmt.Printf("Reason:          %s\n", text)
	fmt.Printf("Server Order ID: %s\n", orderID)
	fmt.Printf("Orig Client ID:  %s\n", origClOrdID)
	fmt.Printf("Order Status:    %s\n", translateOrdStatus(ordStatus))
	fmt.Printf("-----------------------------\n")
}

// printReject parses and prints the details of a Session Level Reject (MsgType 3).
func printReject(msg *quickfix.Message) {
	refSeqNum, _ := msg.Body.GetString(TagRefSeqNum)
	refMsgType, _ := msg.Body.GetString(TagRefMsgType)
	text, _ := msg.Body.GetString(TagText)

	// This happens when the server rejects a message before it even reaches the matching engine
	fmt.Printf("\n--- [Session Level Reject (3)] ---\n")
	fmt.Printf("Reason:          %s\n", text)
	fmt.Printf("Rejected Msg:    %s (MsgType: %s)\n", translateMsgType(refMsgType), refMsgType)
	fmt.Printf("Ref Seq Number:  %s\n", refSeqNum)
	fmt.Printf("----------------------------------\n")
}

// sendOrder constructs and sends a "New Order Single" message to the server.
func (a *Application) sendOrder() {
	// 1. Create a new empty message
	msg := quickfix.NewMessage()

	// 2. Set the Header fields (Required)
	//    MsgType 'D' tells the server this is a New Order Single
	msg.Header.SetField(TagMsgType, quickfix.FIXString(MsgTypeNewOrderSingle))

	// 3. Set the Body fields (The details of the order)

	// ClOrdID: A unique ID WE generate to track this order
	clOrdID := fmt.Sprintf("ORD-%d", time.Now().UnixNano())
	msg.Body.SetField(TagClOrdID, quickfix.FIXString(clOrdID))

	// Symbol: What we want to trade
	msg.Body.SetField(TagSymbol, quickfix.FIXString("EUR/USD"))

	// Side: Buy (1)
	msg.Body.SetField(TagSide, quickfix.FIXString(SideBuy))

	// TransactTime: Current UTC timestamp
	msg.Body.SetField(TagTransactTime, quickfix.FIXString(time.Now().Format("20060102-15:04:05.000")))

	// OrderQty: How much we want
	msg.Body.SetField(TagOrderQty, quickfix.FIXString("100"))

	// OrdType: Market (1) - execute immediately at best price
	msg.Body.SetField(TagOrdType, quickfix.FIXString(OrdTypeMarket))

	// 4. Send the message to the session
	log.Printf("Action: Sending New Order Single (ID: %s)...", clOrdID)
	err := quickfix.SendToTarget(msg, a.SessionID)
	if err != nil {
		log.Printf("Error: Failed to send order: %s", err)
	}
}

// -----------------------------------------------------------------------------
// UTILITIES
// -----------------------------------------------------------------------------

// translateExecType converts FIX codes (0, 1, 2) to human names (New, Partial, Fill)
func translateExecType(val string) string {
	switch val {
	case ExecTypeNew:
		return "New (Accepted)"
	case ExecTypePartialFill:
		return "Partially Filled"
	case ExecTypeFill:
		return "Filled (Complete)"
	case ExecTypeCanceled:
		return "Canceled"
	case ExecTypeRejected:
		return "Rejected"
	default:
		return val // Return raw value if unknown
	}
}

// translateOrdStatus converts FIX status codes to human names
func translateOrdStatus(val string) string {
	switch val {
	case "0":
		return "New"
	case "1":
		return "Partially Filled"
	case "2":
		return "Filled"
	case "4":
		return "Canceled"
	case "8":
		return "Rejected"
	default:
		return val
	}
}

// translateMsgType converts FIX MsgType codes to human-readable names
func translateMsgType(msgType string) string {
	switch msgType {
	case MsgTypeLogon:
		return "Logon"
	case MsgTypeReject:
		return "Session Reject"
	case MsgTypeExecutionReport:
		return "Execution Report"
	case MsgTypeOrderCancelReject:
		return "Order Cancel Reject"
	case MsgTypeNewOrderSingle:
		return "New Order Single"
	case MsgTypeOrderCancelRequest:
		return "Order Cancel Request"
	default:
		return fmt.Sprintf("Unknown Type (%s)", msgType)
	}
}
