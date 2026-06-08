package fix

import (
	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/tag"
)

// This file keeps short, project-friendly names for the FIX tags and enum
// values used by the app. Standard FIX values come from quickfixgo/tag and
// quickfixgo/enum. Literal values below are local RQD/Raptor conventions.

const (
	TagMsgType               = tag.MsgType
	TagClOrdID               = tag.ClOrdID
	TagSymbol                = tag.Symbol
	TagSymbolSfx             = tag.SymbolSfx
	TagSecurityType          = tag.SecurityType
	TagMaturityMonthYear     = tag.MaturityMonthYear
	TagMaturityDay           = tag.MaturityDay
	TagPutOrCall             = tag.PutOrCall
	TagStrikePrice           = tag.StrikePrice
	TagContractMultiplier    = tag.ContractMultiplier
	TagSide                  = tag.Side
	TagTransactTime          = tag.TransactTime
	TagOrderQty              = tag.OrderQty
	TagOrdType               = tag.OrdType
	TagPrice                 = tag.Price
	TagStopPx                = tag.StopPx
	TagOrderID               = tag.OrderID
	TagRule80A               = tag.Rule80A
	TagSettlmntTyp           = tag.SettlmntTyp
	TagTargetSubID           = tag.TargetSubID
	TagOpenClose             = tag.OpenClose
	TagExecInst              = tag.ExecInst
	TagExecType              = tag.ExecType
	TagOrdStatus             = tag.OrdStatus
	TagCumQty                = tag.CumQty
	TagLeavesQty             = tag.LeavesQty
	TagAvgPx                 = tag.AvgPx
	TagText                  = tag.Text
	TagOrigClOrdID           = tag.OrigClOrdID
	TagRefSeqNum             = tag.RefSeqNum
	TagRefMsgType            = tag.RefMsgType
	TagAccount               = tag.Account
	TagTimeInForce           = tag.TimeInForce
	TagLocateReqd            = tag.LocateReqd
	TagCashOrderQty          = tag.CashOrderQty
	TagTradingSessionID      = tag.TradingSessionID
	TagNoLegs                = tag.NoLegs
	TagLegPositionEffect     = tag.LegPositionEffect
	TagLegPrice              = tag.LegPrice
	TagLegSymbol             = tag.LegSymbol
	TagLegCFICode            = tag.LegCFICode
	TagLegSecurityType       = tag.LegSecurityType
	TagLegMaturityMonthYear  = tag.LegMaturityMonthYear
	TagLegMaturityDate       = tag.LegMaturityDate
	TagLegStrikePrice        = tag.LegStrikePrice
	TagLegContractMultiplier = tag.LegContractMultiplier
	TagLegRatioQty           = tag.LegRatioQty
	TagLegSide               = tag.LegSide
	TagLegRefID              = tag.LegRefID
	TagLegQty                = tag.LegQty

	TagLocateID               quickfix.Tag = 5700  // RQD locate ID for short-sale orders
	TagTargetRaptorFractional quickfix.Tag = 20038 // Raptor fractional quantity tag
)

const (
	MsgTypeLogon               = string(enum.MsgType_LOGON)
	MsgTypeReject              = string(enum.MsgType_REJECT)
	MsgTypeExecutionReport     = string(enum.MsgType_EXECUTION_REPORT)
	MsgTypeOrderCancelReject   = string(enum.MsgType_ORDER_CANCEL_REJECT)
	MsgTypeNewOrderSingle      = string(enum.MsgType_ORDER_SINGLE)
	MsgTypeOrderCancelRequest  = string(enum.MsgType_ORDER_CANCEL_REQUEST)
	MsgTypeOrderReplaceRequest = string(enum.MsgType_ORDER_CANCEL_REPLACE_REQUEST)
	MsgTypeOrderStatusRequest  = string(enum.MsgType_ORDER_STATUS_REQUEST)
	MsgTypeNewOrderMultileg    = string(enum.MsgType_NEW_ORDER_MULTILEG)
	MsgTypeHeartBeat           = string(enum.MsgType_HEARTBEAT)
)

const (
	SideBuy       = string(enum.Side_BUY)
	SideSell      = string(enum.Side_SELL)
	SideSellShort = string(enum.Side_SELL_SHORT)

	OrdTypeMarket        = string(enum.OrdType_MARKET)
	OrdTypeLimit         = string(enum.OrdType_LIMIT)
	OrdTypeStop          = string(enum.OrdType_STOP)
	OrdTypeStopLimit     = string(enum.OrdType_STOP_LIMIT)
	OrdTypeMarketOnClose = string(enum.OrdType_MARKET_ON_CLOSE)
	OrdTypeLimitOnClose  = string(enum.OrdType_LIMIT_ON_CLOSE)

	SecurityTypeOption   = string(enum.SecurityType_OPTION)
	SecurityTypeCommon   = string(enum.SecurityType_COMMON_STOCK)
	SecurityTypeMultileg = string(enum.SecurityType_MULTILEG_INSTRUMENT)

	PutOrCallPut  = string(enum.PutOrCall_PUT)
	PutOrCallCall = string(enum.PutOrCall_CALL)

	OpenCloseOpen  = string(enum.OpenClose_OPEN)
	OpenCloseClose = string(enum.OpenClose_CLOSE)

	Rule80AAgency = string(enum.Rule80A_AGENCY_SINGLE_ORDER)

	SettlmntTypRegular       = string(enum.SettlmntTyp_REGULAR)
	SettlmntTypCash          = string(enum.SettlmntTyp_CASH)
	SettlmntTypNextDay       = string(enum.SettlmntTyp_NEXT_DAY)
	SettlmntTypTplus2        = string(enum.SettlmntTyp_T_PLUS_2)
	SettlmntTypTplus3        = string(enum.SettlmntTyp_T_PLUS_3)
	SettlmntTypTplus4        = string(enum.SettlmntTyp_T_PLUS_4)
	SettlmntTypFuture        = string(enum.SettlmntTyp_FUTURE)
	SettlmntTypWhenIssued    = string(enum.SettlmntTyp_WHEN_AND_IF_ISSUED)
	SettlmntTypSellersOption = string(enum.SettlmntTyp_SELLERS_OPTION)
	SettlmntTypTplus5        = string(enum.SettlmntTyp_T_PLUS_5)

	ExecInstNotHeld = string(enum.ExecInst_NOT_HELD)
	ExecInstHeld    = string(enum.ExecInst_HELD)

	TimeInForceDay = string(enum.TimeInForce_DAY)
	TimeInForceGTC = string(enum.TimeInForce_GOOD_TILL_CANCEL)
	TimeInForceIOC = string(enum.TimeInForce_IMMEDIATE_OR_CANCEL)
	TimeInForceFOK = string(enum.TimeInForce_FILL_OR_KILL)

	ExecTypeNew         = string(enum.ExecType_NEW)
	ExecTypePartialFill = string(enum.ExecType_PARTIAL_FILL)
	ExecTypeFill        = string(enum.ExecType_FILL)
	ExecTypeCanceled    = string(enum.ExecType_CANCELED)
	ExecTypeRejected    = string(enum.ExecType_REJECTED)
)

const (
	LegCFICodeCall = "OC" // Option call; used because FIX 4.4 has no direct LegPutOrCall field here.
	LegCFICodePut  = "OP" // Option put; used because FIX 4.4 has no direct LegPutOrCall field here.

	TradingSessionAM   = "1" // RQD AM only
	TradingSessionPM   = "2" // RQD PM only
	TradingSessionBoth = "3" // RQD all sessions
)
