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
	TagMsgType                quickfix.Tag = 35    // Identifies the type of message (e.g., Order, Logon)
	TagClOrdID                quickfix.Tag = 11    // Client Order ID (Unique ID we assign)
	TagSymbol                 quickfix.Tag = 55    // The thing we are trading (e.g., EUR/USD)
	TagSymbolSfx              quickfix.Tag = 65    // Suffix for the symbol
	TagSide                   quickfix.Tag = 54    // Buy (1) or Sell (2)
	TagTransactTime           quickfix.Tag = 60    // Time the order was sent
	TagOrderQty               quickfix.Tag = 38    // How many units to buy/sell
	TagOrdType                quickfix.Tag = 40    // Market (1) or Limit (2)
	TagPrice                  quickfix.Tag = 44    // Limit price
	TagStopPx                 quickfix.Tag = 99    // Stop price
	TagOrderID                quickfix.Tag = 37    // Server's Order ID
	TagRule80A                quickfix.Tag = 47    // Account Type (Agency, Principal, etc. - Tag 47)
	TagSettlmntTyp            quickfix.Tag = 63    // Settlement Type (Tag 63)
	TagTargetSubID            quickfix.Tag = 57    // Target Sub ID (Used for FRAC indicator)
	TagTargetRaptorFractional quickfix.Tag = 20038 // Raptor wants 20038 to provide actual fractional shares
	TagExecInst               quickfix.Tag = 18    // Instructions for order handling (Not Held, etc.)
	TagExecTransType          quickfix.Tag = 20    // Execution Transaction Type (Tag 20)
	TagExecType               quickfix.Tag = 150   // What happened to the order? (New, Filled, etc.)
	TagOrdStatus              quickfix.Tag = 39    // Current status of the order
	TagCumQty                 quickfix.Tag = 14    // Total quantity filled so far
	TagLeavesQty              quickfix.Tag = 151   // Quantity remaining to be filled
	TagAvgPx                  quickfix.Tag = 6     // Average price of fills
	TagText                   quickfix.Tag = 58    // Text description / reason
	TagOrigClOrdID            quickfix.Tag = 41    // Original Order ID (for cancels/modifies)
	TagRefSeqNum              quickfix.Tag = 45    // Reference message sequence number (for Rejects)
	TagRefMsgType             quickfix.Tag = 372   // Reference message type (for Rejects)
	TagAccount                quickfix.Tag = 1     // Account ID (e.g., Trading Account)
	TagTimeInForce            quickfix.Tag = 59    // How long the order stays active
	TagLocateReqd             quickfix.Tag = 114   // Locate Required (for Short Sales)
	TagLocateID               quickfix.Tag = 5700  // Locate ID (Custom/Specific tag often used for SS)
	TagCashOrderQty           quickfix.Tag = 152   // Notional Order Quantity (Tag 152)
	TagTradingSessionID       quickfix.Tag = 336   // Trading Session ID (Tag 336)

	// Message Types (Values for Tag 35)
	MsgTypeLogon               = "A" // Connection established
	MsgTypeReject              = "3" // Session-level reject
	MsgTypeExecutionReport     = "8" // Server telling us about an order change
	MsgTypeOrderCancelReject   = "9" // Server rejected our request to cancel
	MsgTypeNewOrderSingle      = "D" // We are sending a new order
	MsgTypeOrderCancelRequest  = "F" // Request to cancel an existing order
	MsgTypeOrderReplaceRequest = "G" // Request to modify an existing order
	MsgTypeOrderStatusRequest  = "H" // Request for an order's current status
	MsgTypeHeartBeat           = "0" // Heartbeat

	// Field Values (Tag 54 - Side)
	SideBuy       = "1" // Value '1' means Buy
	SideSell      = "2" // Value '2' means Sell
	SideSellShort = "5" // Value '5' means Sell Short

	// Field Values (Tag 40 - OrdType)
	OrdTypeMarket        = "1" // Value '1' means Market Order
	OrdTypeLimit         = "2" // Value '2' means Limit Order
	OrdTypeStop          = "3" // Value '3' means Stop Order
	OrdTypeStopLimit     = "4" // Value '4' means Stop Limit Order
	OrdTypeMarketOnClose = "5" // Value '5' means Market On Close
	OrdTypeLimitOnClose  = "B" // Value 'B' means Limit On Close

	// Rule80A Values (Tag 47 - Account Type)
	Rule80AAgency    = "A" // Agency
	Rule80APrincipal = "P" // Principal

	// SettlmntTyp Values (Tag 63)
	SettlmntTypRegular       = "0" // Regular
	SettlmntTypCash          = "1" // Cash
	SettlmntTypNextDay       = "2" // Next Day
	SettlmntTypTplus2        = "3" // T+2
	SettlmntTypTplus3        = "4" // T+3
	SettlmntTypTplus4        = "5" // T+4
	SettlmntTypFuture        = "6" // Future
	SettlmntTypWhenIssued    = "7" // When Issued
	SettlmntTypSellersOption = "8" // Sellers Option
	SettlmntTypTplus5        = "9" // T+5

	// ExecInst Values (Tag 18)
	ExecInstNotHeld = "1" // Broker is not held to immediate execution
	ExecInstHeld    = "5" // Broker is held to immediate execution

	// ExecTransType Values (Tag 20)
	ExecTransTypeNew     = "0" // New
	ExecTransTypeCancel  = "1" // Cancel (Bust)
	ExecTransTypeCorrect = "2" // Correct (Price Correct)

	// TradingSessionID Values (Tag 336)
	TradingSessionAM   = "1" // AM Only
	TradingSessionPM   = "2" // PM Only
	TradingSessionBoth = "3" // All Sessions

	// Time In Force Values (Tag 59)
	TimeInForceDay = "0" // Active for the trading day
	TimeInForceGTC = "1" // Good Till Cancelled
	TimeInForceIOC = "3" // Immediate Or Cancel
	TimeInForceFOK = "4" // Fill Or Kill

	// Execution Types (Values for Tag 150) - What happened?
	ExecTypeNew         = "0" // Order accepted
	ExecTypePartialFill = "1" // Part of the order filled
	ExecTypeFill        = "2" // Entire order filled
	ExecTypeCanceled    = "4" // Order canceled
	ExecTypeRejected    = "8" // Order rejected
)
