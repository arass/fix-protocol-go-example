package fix

import "github.com/quickfixgo/quickfix"

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
	TagRefSeqNum    quickfix.Tag = 45  // Reference message sequence number (for Rejects)
	TagRefMsgType   quickfix.Tag = 372 // Reference message type (for Rejects)

	// Message Types (Values for Tag 35)
	MsgTypeLogon              = "A" // Connection established
	MsgTypeReject             = "3" // Session-level reject
	MsgTypeExecutionReport    = "8" // Server telling us about an order change
	MsgTypeOrderCancelReject  = "9" // Server rejected our request to cancel
	MsgTypeNewOrderSingle     = "D" // We are sending a new order
	MsgTypeOrderCancelRequest = "F" // Request to cancel an existing order
	MsgTypeHeartBeat          = "0" // Request to cancel an existing order

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
