package fix

import (
	"fmt"
	"log"
	"os"
	"strings"
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
	LastClOrdID string             // Stores the ID of the last order sent (to allow cancelling/replacing/statusing it)
	IsTestMode  bool               // If true, disables the automatic "New -> Status -> Replace -> Cancel" flow
	OnLogonChan chan bool          // Channel to signal when logon is successful
}

// NewApplication creates a new instance of our FIX application logic.
func NewApplication() *Application {
	return &Application{
		IsTestMode:  false,
		OnLogonChan: make(chan bool, 1),
	}
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

	// Signal logon success
	select {
	case a.OnLogonChan <- true:
	default:
	}

	if a.IsTestMode {
		log.Printf("System: Application is in TEST MODE. Automatic flow disabled.")
		return
	}

	// We send an order in a separate thread (goroutine) so we don't block the
	// message processing loop. We wait 1 second to ensure everything is ready.
	go func() {
		time.Sleep(1 * time.Second)
		a.SendOrder(OrderParams{Symbol: "AAPL", Side: SideBuy, Qty: "100", OrdType: OrdTypeMarket, TIF: TimeInForceDay})

		// Wait 3 seconds, then request status
		time.Sleep(3 * time.Second)
		a.SendOrderStatusRequest()

		// Wait 3 seconds, then try to modify (replace) it
		time.Sleep(3 * time.Second)
		a.SendReplaceOrder("200", OrdTypeLimit, "150.50")

		// Wait 3 seconds and then try to cancel it
		time.Sleep(3 * time.Second)
		a.SendCancelOrder()
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

// OrderParams is a helper struct to avoid long parameter lists in SendOrder
type OrderParams struct {
	Symbol       string
	Side         string
	Qty          string
	OrdType      string
	LimitPrice   string
	StopPrice    string
	TIF          string
	ExecInst     string
	SettlTyp     string
	IsFractional bool
	Notional     string
	TradingSes   string
}

// SendOrder constructs and sends a "New Order Single" message to the server.
func (a *Application) SendOrder(p OrderParams) string {
	// 1. Create a new empty message
	msg := quickfix.NewMessage()

	// 2. Set the Header fields (Required)
	msg.Header.SetField(TagMsgType, quickfix.FIXString(MsgTypeNewOrderSingle))

	// 3. Set the Body fields

	account := os.Getenv("ACCOUNT_NUMBER")
	if account == "" {
		account = "FIX-TEST-ACCOUNT-1"
	}
	msg.Body.SetField(TagAccount, quickfix.FIXString(account))

	// All orders marked as Agency (Tag 47=A)
	msg.Body.SetField(TagRule80A, quickfix.FIXString(Rule80AAgency))

	clOrdID := fmt.Sprintf("ORD-%d", time.Now().UnixNano())
	msg.Body.SetField(TagClOrdID, quickfix.FIXString(clOrdID))
	a.LastClOrdID = clOrdID

	// Symbology handling for suffixes, e.g. BRK.B
	symbol, symbolSfx, hasSfx := strings.Cut(p.Symbol, ".")

	msg.Body.SetField(TagSymbol, quickfix.FIXString(symbol))

	if hasSfx && symbolSfx != "" {
		msg.Body.SetField(TagSymbolSfx, quickfix.FIXString(symbolSfx))
	}

	msg.Body.SetField(TagSide, quickfix.FIXString(p.Side))
	msg.Body.SetField(TagTransactTime, quickfix.FIXString(time.Now().Format("20060102-15:04:05.000")))

	// Handling Quantity vs Notional
	if p.Notional != "" {
		msg.Body.SetField(TagCashOrderQty, quickfix.FIXString(p.Notional))
	} else {
		msg.Body.SetField(TagOrderQty, quickfix.FIXString(p.Qty))
	}

	msg.Body.SetField(TagOrdType, quickfix.FIXString(p.OrdType))

	if p.LimitPrice != "" {
		msg.Body.SetField(TagPrice, quickfix.FIXString(p.LimitPrice))
	}

	if p.StopPrice != "" {
		msg.Body.SetField(TagStopPx, quickfix.FIXString(p.StopPrice))
	}

	if p.TIF != "" {
		msg.Body.SetField(TagTimeInForce, quickfix.FIXString(p.TIF))
	}

	if p.ExecInst != "" {
		msg.Body.SetField(TagExecInst, quickfix.FIXString(p.ExecInst))
	}

	if p.SettlTyp != "" {
		msg.Body.SetField(TagSettlmntTyp, quickfix.FIXString(p.SettlTyp))
	}

	if p.IsFractional {
		msg.Header.SetField(TagTargetSubID, quickfix.FIXString("FRAC"))
	}

	if p.TradingSes != "" {
		msg.Body.SetField(TagTradingSessionID, quickfix.FIXString(p.TradingSes))
	}

	// Special handling for Sell Short (Side=5)
	if p.Side == SideSellShort {
		msg.Body.SetField(TagLocateReqd, quickfix.FIXString("N"))
		msg.Body.SetField(TagLocateID, quickfix.FIXString("LOCATE-ID-123"))
	}

	// 4. Send the message to the session
	log.Printf("Action: Sending New Order Single (ID: %s, Symbol: %s, Side: %s, Type: %s)...", clOrdID, p.Symbol, p.Side, p.OrdType)
	err := quickfix.SendToTarget(msg, a.SessionID)
	if err != nil {
		log.Printf("Error: Failed to send order: %s", err)
	}
	return clOrdID
}

// SendOrderStatusRequest constructs and sends an "Order Status Request" message.
func (a *Application) SendOrderStatusRequest() {
	if a.LastClOrdID == "" {
		log.Printf("Action: No previous order to check status for.")
		return
	}

	// 1. Create a new empty message
	msg := quickfix.NewMessage()

	// 2. Set the Header fields
	//    MsgType 'H' tells the server this is an Order Status Request
	msg.Header.SetField(TagMsgType, quickfix.FIXString(MsgTypeOrderStatusRequest))

	// 3. Set the Body fields

	// ClOrdID: The ID of the order we want to check
	msg.Body.SetField(TagClOrdID, quickfix.FIXString(a.LastClOrdID))

	// Side: Must match the original order
	msg.Body.SetField(TagSide, quickfix.FIXString(SideBuy))

	// Symbol: Must match the original order
	msg.Body.SetField(TagSymbol, quickfix.FIXString("AAPL"))

	// 4. Send the message
	log.Printf("Action: Sending Order Status Request (ID: %s)...", a.LastClOrdID)
	err := quickfix.SendToTarget(msg, a.SessionID)
	if err != nil {
		log.Printf("Error: Failed to send status request: %s", err)
	}
}

// SendReplaceOrder constructs and sends an "Order Cancel/Replace Request" message.
func (a *Application) SendReplaceOrder(qty string, ordType string, price string) {
	if a.LastClOrdID == "" {
		log.Printf("Action: No previous order to replace.")
		return
	}

	// 1. Create a new empty message
	msg := quickfix.NewMessage()

	// 2. Set the Header fields
	//    MsgType 'G' tells the server this is an Order Cancel/Replace Request
	msg.Header.SetField(TagMsgType, quickfix.FIXString(MsgTypeOrderReplaceRequest))

	// 3. Set the Body fields

	// OrigClOrdID: The ID of the order we want to replace
	msg.Body.SetField(TagOrigClOrdID, quickfix.FIXString(a.LastClOrdID))

	// ClOrdID: A NEW unique ID for the replace request itself
	replaceClOrdID := fmt.Sprintf("REPL-%d", time.Now().UnixNano())
	msg.Body.SetField(TagClOrdID, quickfix.FIXString(replaceClOrdID))

	// Account: Same account
	account := os.Getenv("ACCOUNT_NUMBER")
	if account == "" {
		account = "FIX-TEST-ACCOUNT-1"
	}
	msg.Body.SetField(TagAccount, quickfix.FIXString(account))

	// All replaces marked as Agency (Tag 47=A)
	msg.Body.SetField(TagRule80A, quickfix.FIXString(Rule80AAgency))

	// Symbol: Must match the original order
	msg.Body.SetField(TagSymbol, quickfix.FIXString("AAPL"))

	// Side: Must match the original order
	msg.Body.SetField(TagSide, quickfix.FIXString(SideBuy))

	// TransactTime: Current UTC timestamp
	msg.Body.SetField(TagTransactTime, quickfix.FIXString(time.Now().Format("20060102-15:04:05.000")))

	// MODIFICATIONS
	msg.Body.SetField(TagOrderQty, quickfix.FIXString(qty))
	msg.Body.SetField(TagOrdType, quickfix.FIXString(ordType))
	if ordType == OrdTypeLimit && price != "" {
		msg.Body.SetField(TagPrice, quickfix.FIXString(price))
	}

	// 4. Send the message
	log.Printf("Action: Sending Order Replace Request (Target ID: %s, New ID: %s)...", a.LastClOrdID, replaceClOrdID)

	// Crucial: Update our last ID to the new one
	a.LastClOrdID = replaceClOrdID

	err := quickfix.SendToTarget(msg, a.SessionID)
	if err != nil {
		log.Printf("Error: Failed to send replace request: %s", err)
	}
}

// SendCancelOrder constructs and sends an "Order Cancel Request" message.
func (a *Application) SendCancelOrder() {
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

	// OrigClOrdID: The ID of the order we want to cancel
	msg.Body.SetField(TagOrigClOrdID, quickfix.FIXString(a.LastClOrdID))

	// ClOrdID: A NEW unique ID for the cancel request itself
	cancelClOrdID := fmt.Sprintf("CAN-%d", time.Now().UnixNano())
	msg.Body.SetField(TagClOrdID, quickfix.FIXString(cancelClOrdID))

	// Account: Same account
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
