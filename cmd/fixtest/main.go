package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"thisgofix/internal/fix"

	"github.com/quickfixgo/quickfix"
)

// testDelay is the duration to pause between sending test orders.
// This can be adjusted to allow more time for the FIX engine to process messages
// and for execution reports to be received.
const testDelay = 1 * time.Second

func main() {
	cfgFileName := flag.String("cfg", "config.cfg", "Path to config file")
	flag.Parse()

	cfg, err := os.Open(*cfgFileName)
	if err != nil {
		log.Fatalf("Critical Error: Could not open config file '%s': %s", *cfgFileName, err)
	}
	defer cfg.Close()

	appSettings, err := quickfix.ParseSettings(cfg)
	if err != nil {
		log.Fatalf("Critical Error: Config file format is invalid: %s", err)
	}

	// 1. Setup Application in TEST MODE
	app := fix.NewApplication()
	app.IsTestMode = true

	messageStoreFactory := quickfix.NewMemoryStoreFactory()
	logFactory := &ConsoleLogFactory{}

	initiator, err := quickfix.NewInitiator(app, messageStoreFactory, appSettings, logFactory)
	if err != nil {
		log.Fatalf("Critical Error: Could not create Initiator: %s", err)
	}

	if err := initiator.Start(); err != nil {
		log.Fatalf("Critical Error: Failed to start application: %s", err)
	}

	// 2. Wait for Logon before starting tests
	log.Println("TestRunner: Waiting for FIX Logon...")
	select {
	case <-app.OnLogonChan:
		log.Println("TestRunner: Logged on. Starting scenarios...")
	case <-time.After(10 * time.Second):
		log.Fatalf("TestRunner: Timeout waiting for logon")
	}

	// 3. Execute Scenarios (Categorized)
	testSides(app)
	testOrderTypes(app)
	testExecInst(app)
	testSymbology(app)
	testTIF(app)
	testSettlement(app)
	testFractional(app)
	testNotional(app)
	testExtendedHours(app)
	testMisc(app)

	// 4. Keep running for a few seconds to see execution reports
	log.Println("TestRunner: All scenarios sent. Waiting for responses...")
	time.Sleep(10 * time.Second)

	log.Println("TestRunner: Stopping FIX Client...")
	initiator.Stop()
}

func testSides(app *fix.Application) {
	log.Println("--- Group: SIDES ---")
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeMarket, TIF: fix.TimeInForceDay})
	time.Sleep(testDelay)
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideSell, Qty: "1", OrdType: fix.OrdTypeLimit, LimitPrice: "300", TIF: fix.TimeInForceGTC})
	time.Sleep(testDelay)
	app.SendOrder(fix.OrderParams{Symbol: "QQQ", Side: fix.SideSellShort, Qty: "1", OrdType: fix.OrdTypeLimit, LimitPrice: "650", TIF: fix.TimeInForceDay})
	time.Sleep(testDelay)
}

func testOrderTypes(app *fix.Application) {
	log.Println("--- Group: ORDER TYPE ---")
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeMarket, TIF: fix.TimeInForceDay})
	time.Sleep(testDelay)
	app.SendOrder(fix.OrderParams{Symbol: "BRK.B", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeStop, StopPrice: "500", TIF: fix.TimeInForceDay})
	time.Sleep(testDelay)
	app.SendOrder(fix.OrderParams{Symbol: "BABA", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeStopLimit, StopPrice: "150", LimitPrice: "150", TIF: fix.TimeInForceDay})
	time.Sleep(testDelay)
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeMarketOnClose, TIF: fix.TimeInForceDay})
	time.Sleep(testDelay)
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeLimitOnClose, LimitPrice: "150", TIF: fix.TimeInForceDay})
	time.Sleep(testDelay)
}

func testExecInst(app *fix.Application) {
	log.Println("--- Group: EXEC INST ---")
	app.SendOrder(fix.OrderParams{Symbol: "BKNG", Side: fix.SideBuy, Qty: "100", OrdType: fix.OrdTypeMarket, TIF: fix.TimeInForceDay, ExecInst: fix.ExecInstNotHeld})
	time.Sleep(testDelay)
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeMarket, TIF: fix.TimeInForceDay, ExecInst: fix.ExecInstHeld})
	time.Sleep(testDelay)
}

func testSymbology(app *fix.Application) {
	log.Println("--- Group: SYMBOLOGY ---")
	app.SendOrder(fix.OrderParams{Symbol: "BRK.B", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeStop, StopPrice: "500", TIF: fix.TimeInForceDay})
	time.Sleep(testDelay)
}

func testTIF(app *fix.Application) {
	log.Println("--- Group: TIF ---")
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeMarket, TIF: fix.TimeInForceDay})
	time.Sleep(testDelay)
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideSell, Qty: "1", OrdType: fix.OrdTypeLimit, LimitPrice: "300", TIF: fix.TimeInForceGTC})
	time.Sleep(testDelay)
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeMarket, TIF: fix.TimeInForceIOC})
	time.Sleep(testDelay)
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeMarket, TIF: fix.TimeInForceFOK})
	time.Sleep(testDelay)
}

func testSettlement(app *fix.Application) {
	log.Println("--- Group: SETTLEMENT ---")
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeMarket, TIF: fix.TimeInForceDay, SettlTyp: fix.SettlmntTypRegular})
	time.Sleep(testDelay)
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeMarket, TIF: fix.TimeInForceDay, SettlTyp: fix.SettlmntTypCash})
	time.Sleep(testDelay)
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeMarket, TIF: fix.TimeInForceDay, SettlTyp: fix.SettlmntTypNextDay})
	time.Sleep(testDelay)
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeMarket, TIF: fix.TimeInForceDay, SettlTyp: fix.SettlmntTypTplus2})
	time.Sleep(testDelay)
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeMarket, TIF: fix.TimeInForceDay, SettlTyp: fix.SettlmntTypTplus3})
	time.Sleep(testDelay)
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeMarket, TIF: fix.TimeInForceDay, SettlTyp: fix.SettlmntTypTplus4})
	time.Sleep(testDelay)
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeMarket, TIF: fix.TimeInForceDay, SettlTyp: fix.SettlmntTypFuture})
	time.Sleep(testDelay)
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeMarket, TIF: fix.TimeInForceDay, SettlTyp: fix.SettlmntTypWhenIssued})
	time.Sleep(testDelay)
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeMarket, TIF: fix.TimeInForceDay, SettlTyp: fix.SettlmntTypSellersOption})
	time.Sleep(testDelay)
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeMarket, TIF: fix.TimeInForceDay, SettlTyp: fix.SettlmntTypTplus5})
	time.Sleep(testDelay)
}

func testFractional(app *fix.Application) {
	log.Println("--- Group: FRACTIONAL ---")
	log.Println("Scenario 38: Fractional MKT Buy .25 AMZN")
	app.SendOrder(fix.OrderParams{Symbol: "AMZN", Side: fix.SideBuy, Qty: ".25", OrdType: fix.OrdTypeMarket, IsFractional: true})
	time.Sleep(testDelay)

	log.Println("Scenario 39: Fractional LMT Buy .99 TSLA @ 400")
	app.SendOrder(fix.OrderParams{Symbol: "TSLA", Side: fix.SideBuy, Qty: ".99", OrdType: fix.OrdTypeLimit, LimitPrice: "400", IsFractional: true})
	time.Sleep(testDelay)

	log.Println("Scenario 40: Fractional GTC Buy .50 TSLA @ 401 LMT GTC")
	app.SendOrder(fix.OrderParams{Symbol: "TSLA", Side: fix.SideBuy, Qty: ".50", OrdType: fix.OrdTypeLimit, LimitPrice: "401", TIF: fix.TimeInForceGTC, IsFractional: true})
	time.Sleep(testDelay)

	log.Println("Scenario 41: Cancel Fractional Order")
	app.SendOrder(fix.OrderParams{Symbol: "TSLA", Side: fix.SideBuy, Qty: ".50", OrdType: fix.OrdTypeLimit, LimitPrice: "401", TIF: fix.TimeInForceGTC, IsFractional: true})
	time.Sleep(testDelay)
	app.SendCancelOrder()
	time.Sleep(testDelay)

	log.Println("Scenario 42: Replace Fractional Order")
	app.SendOrder(fix.OrderParams{Symbol: "TSLA", Side: fix.SideBuy, Qty: ".10", OrdType: fix.OrdTypeLimit, LimitPrice: "402", IsFractional: true})
	time.Sleep(testDelay)
	app.SendReplaceOrder(".10", fix.OrdTypeLimit, "403")
	time.Sleep(testDelay)
}

func testNotional(app *fix.Application) {
	log.Println("--- Group: NOTIONAL ---")
	log.Println("Scenario 43: Notional MKT Buy $100 TSLA")
	app.SendOrder(fix.OrderParams{Symbol: "TSLA", Side: fix.SideBuy, Notional: "100", OrdType: fix.OrdTypeMarket})
	time.Sleep(testDelay)

	log.Println("Scenario 44: Notional LMT Buy $135 TSLA @ 500")
	app.SendOrder(fix.OrderParams{Symbol: "TSLA", Side: fix.SideBuy, Notional: "135", OrdType: fix.OrdTypeLimit, LimitPrice: "500"})
	time.Sleep(testDelay)

	log.Println("Scenario 45: Notional GTC Buy $150 TSLA @ 450 LMT")
	app.SendOrder(fix.OrderParams{Symbol: "TSLA", Side: fix.SideBuy, Notional: "150", OrdType: fix.OrdTypeLimit, LimitPrice: "450", TIF: fix.TimeInForceGTC})
	time.Sleep(testDelay)

	log.Println("Scenario 46: Cancel Notional Order")
	app.SendOrder(fix.OrderParams{Symbol: "TSLA", Side: fix.SideBuy, Notional: "175", OrdType: fix.OrdTypeMarket})
	time.Sleep(testDelay)
	app.SendCancelOrder()
	time.Sleep(testDelay)

	log.Println("Scenario 47: Replace Notional Order")
	app.SendOrder(fix.OrderParams{Symbol: "TSLA", Side: fix.SideBuy, Notional: "180", OrdType: fix.OrdTypeMarket})
	time.Sleep(testDelay)
	app.SendReplaceOrder("180", fix.OrdTypeMarket, "")
	time.Sleep(testDelay)
}

func testExtendedHours(app *fix.Application) {
	log.Println("--- Group: EXTENDED HOURS ---")
	log.Println("Scenario 48: AM Only DIA")
	app.SendOrder(fix.OrderParams{Symbol: "DIA", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeLimit, LimitPrice: "450", TradingSes: fix.TradingSessionAM})
	time.Sleep(testDelay)

	log.Println("Scenario 49: PM Only DIA")
	app.SendOrder(fix.OrderParams{Symbol: "DIA", Side: fix.SideBuy, Qty: "2", OrdType: fix.OrdTypeLimit, LimitPrice: "451", TradingSes: fix.TradingSessionPM})
	time.Sleep(testDelay)

	log.Println("Scenario 50: All Sessions DIA")
	app.SendOrder(fix.OrderParams{Symbol: "DIA", Side: fix.SideBuy, Qty: "3", OrdType: fix.OrdTypeLimit, LimitPrice: "452", TradingSes: fix.TradingSessionBoth})
	time.Sleep(testDelay)
}

func testMisc(app *fix.Application) {
	log.Println("--- Group: MISC SCENARIOS ---")
	log.Println("Scenario 27: Partial Fill & Cancel - Buy 1600 LCID MKT")
	app.SendOrder(fix.OrderParams{Symbol: "LCID", Side: fix.SideBuy, Qty: "1600", OrdType: fix.OrdTypeMarket})
	time.Sleep(testDelay)
	app.SendCancelOrder()
	time.Sleep(testDelay)

	log.Println("Scenario 28: Full Fill")
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeMarket})
	time.Sleep(testDelay)

	log.Println("Scenario 29: Cancel an Acknowledged Order")
	app.SendOrder(fix.OrderParams{Symbol: "BRK.B", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeStop, StopPrice: "500"})
	time.Sleep(testDelay)
	app.SendCancelOrder()
	time.Sleep(testDelay)

	log.Println("Scenario 30: Increase Quantity and Price")
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideBuy, Qty: "2", OrdType: fix.OrdTypeLimit, LimitPrice: "450"})
	time.Sleep(testDelay)
	app.SendReplaceOrder("3", fix.OrdTypeLimit, "450")
	time.Sleep(testDelay)

	log.Println("Scenario 31: Reject an order")
	app.SendOrder(fix.OrderParams{Symbol: "LCID", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeLimit, LimitPrice: "2500"})
	time.Sleep(testDelay)

	log.Println("Scenario 32: Reject Cxl/Replace Request")
	app.SendOrder(fix.OrderParams{Symbol: "LCID", Side: fix.SideBuy, Qty: "3", OrdType: fix.OrdTypeLimit, LimitPrice: "5"})
	time.Sleep(testDelay)
	app.SendReplaceOrder("3", fix.OrdTypeLimit, "5")
	time.Sleep(testDelay)

	log.Println("Scenario 33: Market to Limit")
	app.SendOrder(fix.OrderParams{Symbol: "LCID", Side: fix.SideBuy, Qty: "2600", OrdType: fix.OrdTypeMarket})
	time.Sleep(testDelay)
	app.SendReplaceOrder("2600", fix.OrdTypeLimit, "3.00")
	time.Sleep(testDelay)

	log.Println("Scenario 34: Limit to Market")
	app.SendOrder(fix.OrderParams{Symbol: "LCID", Side: fix.SideBuy, Qty: "5", OrdType: fix.OrdTypeLimit, LimitPrice: "5.02"})
	time.Sleep(testDelay)
	app.SendReplaceOrder("5", fix.OrdTypeMarket, "")
	time.Sleep(testDelay)

	log.Println("Scenario 35: Done For Day")
	app.SendOrder(fix.OrderParams{Symbol: "GRO", Side: fix.SideBuy, Qty: "2600", OrdType: fix.OrdTypeLimit, LimitPrice: "3.00"})
	time.Sleep(testDelay)

	log.Println("Scenario 36: Unsolicited Cancel")
	app.SendOrder(fix.OrderParams{Symbol: "GRO", Side: fix.SideBuy, Qty: "2700", OrdType: fix.OrdTypeLimit, LimitPrice: "2.99"})
	time.Sleep(testDelay)

	log.Println("Scenario 37: ExecTransType New (Tag 20=0)")
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeMarket})
	time.Sleep(testDelay)
}

// -----------------------------------------------------------------------------
// REUSE LOG FACTORY FROM MAIN
// -----------------------------------------------------------------------------

type ConsoleLogFactory struct{}

func (f ConsoleLogFactory) Create() (quickfix.Log, error) {
	return &ConsoleLog{prefix: "GLOBAL"}, nil
}

func (f ConsoleLogFactory) CreateSessionLog(sessionID quickfix.SessionID) (quickfix.Log, error) {
	return &ConsoleLog{prefix: sessionID.String()}, nil
}

type ConsoleLog struct {
	prefix string
}

func (l *ConsoleLog) OnIncoming(msg []byte) {
	log.Printf("[%s] < %s", l.prefix, string(msg))
}

func (l *ConsoleLog) OnOutgoing(msg []byte) {
	log.Printf("[%s] > %s", l.prefix, string(msg))
}

func (l *ConsoleLog) OnEvent(msg string) {
	log.Printf("[%s] EVENT: %s", l.prefix, msg)
}

func (l *ConsoleLog) OnEventf(format string, a ...interface{}) {
	log.Printf("[%s] EVENT: %s", l.prefix, fmt.Sprintf(format, a...))
}
