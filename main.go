package main

// -----------------------------------------------------------------------------
// MAIN SOURCE FILE: Go FIX Client
// -----------------------------------------------------------------------------
//
// Purpose:
// This file implements a basic FIX (Financial Information eXchange) client.
// It connects to a trading server, sends an order, and prints the results.
//
// How it works:
// 1. We define a 'FIXApplication' struct that handles FIX events (Logon, Messages).
// 2. We configure the connection (Host, Port, IDs) using a config file.
// 3. We start an 'Initiator' which manages the network connection.
// 4. When we log on, we send a 'NewOrderSingle' message.
// 5. When we receive messages (like Execution Reports), we print them.
//
// -----------------------------------------------------------------------------

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/quickfixgo/quickfix"
)

// -----------------------------------------------------------------------------
// CONSTANTS & TAGS
// -----------------------------------------------------------------------------
// FIX messages use integer "Tags" to identify fields (e.g., Tag 35 is MsgType).
// We define them here to avoid needing complex external imports and to make
// the code easier to read.
// -----------------------------------------------------------------------------

const (
	// Standard FIX Tags (Field IDs)
	TagMsgType      quickfix.Tag = 35  // Identifies the type of message (e.g., Order, Logon)
	TagClOrdID      quickfix.Tag = 11  // Client Order ID (Unique ID we assign)
	TagSymbol       quickfix.Tag = 55  // The thing we are trading (e.g., EUR/USD)
	TagSide         quickfix.Tag = 54  // Buy (1) or Sell (2)
	TagTransactTime quickfix.Tag = 60  // Time the order was sent
	TagOrderQty     quickfix.Tag = 38  // How many units to buy/sell
	TagOrdType      quickfix.Tag = 40  // Market (1) or Limit (2)
	TagOrderID      quickfix.Tag = 37  // Server's Order ID
	TagExecType     quickfix.Tag = 150 // What happened to the order? (New, Filled, etc.)
	TagOrdStatus    quickfix.Tag = 39  // Current status of the order
	TagCumQty       quickfix.Tag = 14  // Total quantity filled so far
	TagLeavesQty    quickfix.Tag = 151 // Quantity remaining to be filled
	TagAvgPx        quickfix.Tag = 6   // Average price of fills
	TagText         quickfix.Tag = 58  // Text description / reason
	TagOrigClOrdID  quickfix.Tag = 41  // Original Order ID (for cancels)

	// Message Types (Values for Tag 35)
	MsgTypeLogon             = "A" // Connection established
	MsgTypeExecutionReport   = "8" // Server telling us about an order change
	MsgTypeOrderCancelReject = "9" // Server rejected our request to cancel
	MsgTypeNewOrderSingle    = "D" // We are sending a new order

	// Field Values
	SideBuy       = "1" // Value '1' means Buy
	OrdTypeMarket = "1" // Value '1' means Market Order

	// Execution Types (Values for Tag 150) - What happened?
	ExecTypeNew         = "0" // Order accepted
	ExecTypePartialFill = "1" // Part of the order filled
	ExecTypeFill        = "2" // Entire order filled
	ExecTypeCanceled    = "4" // Order canceled
	ExecTypeRejected    = "8" // Order rejected
)

// -----------------------------------------------------------------------------
// APPLICATION LOGIC
// -----------------------------------------------------------------------------

// FIXApplication is our custom struct that implements the quickfix.Application interface.
// The FIX engine will call methods on this struct when events happen.
type FIXApplication struct {
	SessionID quickfix.SessionID // Stores the ID of the current connection
}

// OnCreate is called when a FIX session is created (before connecting).
// We use this to save the SessionID for later use.
func (a *FIXApplication) OnCreate(sessionID quickfix.SessionID) {
	a.SessionID = sessionID
	log.Printf("Setup: Session created for %s", sessionID)
}

// OnLogon is called when we successfully connect and authenticate with the server.
// This is where we start our business logic (sending orders).
func (a *FIXApplication) OnLogon(sessionID quickfix.SessionID) {
	log.Printf("Connection: Logged On to %s", sessionID)

	// We send an order in a separate thread (goroutine) so we don't block the
	// message processing loop. We wait 1 second to ensure everything is ready.
	go func() {
		time.Sleep(1 * time.Second)
		a.sendOrder()
	}()
}

// OnLogout is called when we disconnect from the server.
func (a *FIXApplication) OnLogout(sessionID quickfix.SessionID) {
	log.Printf("Connection: Logged Out from %s", sessionID)
}

// ToAdmin is called before sending an administrative message (Logon, Heartbeat).
// You can use this to add Username/Password to the Logon message if needed.
func (a *FIXApplication) ToAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) {
	msgType, err := msg.Header.GetString(TagMsgType)

	// If it's a Logon message (MsgType="A"), we log it.
	if err == nil && msgType == MsgTypeLogon {
		log.Printf("Admin: Sending Logon message...")
	}
}

// ToApp is called before sending an application message (like an Order).
// It's a last chance to modify the message or log it.
func (a *FIXApplication) ToApp(msg *quickfix.Message, sessionID quickfix.SessionID) error {
	log.Printf("Outbound: Sending Application Message: %s", msg.String())
	return nil
}

// FromAdmin is called when we receive an administrative message (Heartbeat, Logon confirmation).
func (a *FIXApplication) FromAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	log.Printf("Admin: Received Admin Message: %s", msg.String())
	return nil
}

// FromApp is called when we receive a business message (ExecutionReport, etc.).
// This is the core handler for incoming data.
func (a *FIXApplication) FromApp(msg *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
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

// sendOrder constructs and sends a "New Order Single" message to the server.
func (a *FIXApplication) sendOrder() {
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

// -----------------------------------------------------------------------------
// MAIN ENTRY POINT
// -----------------------------------------------------------------------------

func main() {
	// 1. Parse command line arguments
	//    We allow the user to specify a different config file if they want.
	//    Usage: go run main.go -cfg myconfig.cfg
	cfgFileName := flag.String("cfg", "config.cfg", "Path to config file")
	flag.Parse()

	// 2. Open the configuration file
	cfg, err := os.Open(*cfgFileName)
	if err != nil {
		log.Fatalf("Critical Error: Could not open config file '%s': %s", *cfgFileName, err)
	}
	// Ensure the file is closed when the function exits
	defer func() { _ = cfg.Close() }()

	// 3. Parse the configuration
	appSettings, err := quickfix.ParseSettings(cfg)
	if err != nil {
		log.Fatalf("Critical Error: Config file format is invalid: %s", err)
	}

	// 4. Create our Application instance (logic defined above)
	app := &FIXApplication{}

	// 5. Create File Log Factory
	//    This tells the engine to write logs to the 'log/' directory defined in config.
	fileLogFactory, err := quickfix.NewFileLogFactory(appSettings)
	if err != nil {
		log.Fatalf("Critical Error: Could not create log factory: %s", err)
	}

	// 6. Create Message Store Factory (Using MemoryStore to avoid compatibility issues)
	//    In a production app, you would use NewFileStoreFactory to save session state.
	messageStoreFactory := quickfix.NewMemoryStoreFactory()

	// 7. Create the Initiator
	//    The Initiator connects to the server (Client -> Server).
	initiator, err := quickfix.NewInitiator(app, messageStoreFactory, appSettings, fileLogFactory)
	if err != nil {
		log.Fatalf("Critical Error: Could not create Initiator: %s", err)
	}

	// 8. Start the Initiator
	//    This opens the network connection and starts the FIX session.
	log.Println("System: Starting FIX Client...")
	if err := initiator.Start(); err != nil {
		log.Fatalf("Critical Error: Failed to start application: %s", err)
	}

	// 9. Keep running until the user hits Ctrl+C
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// Block here until a signal is received
	<-interrupt

	// 10. Clean shutdown
	log.Println("System: Stopping FIX Client...")
	initiator.Stop()
}
