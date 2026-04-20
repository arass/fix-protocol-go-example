package fix

import (
	"fmt"
	"github.com/quickfixgo/quickfix"
)

// -----------------------------------------------------------------------------
// HELPER FUNCTIONS (PRINTING)
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

// -----------------------------------------------------------------------------
// TRANSLATION UTILITIES
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
	case MsgTypeHeartBeat:
		return "Heartbeat"
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
