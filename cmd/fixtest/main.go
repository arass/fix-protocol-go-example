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

	// 3. Execute Scenarios
	runScenarios(app)

	// 4. Keep running for a few seconds to see execution reports
	log.Println("TestRunner: All scenarios sent. Waiting for responses...")
	time.Sleep(10 * time.Second)

	log.Println("TestRunner: Stopping FIX Client...")
	initiator.Stop()
}

func runScenarios(app *fix.Application) {
	// --- SIDES BATCH ---
	log.Println("--- Group: SIDES ---")

	log.Println("Scenario 1: Buy AAPL MKT")
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeMarket, TIF: fix.TimeInForceDay})
	time.Sleep(1 * time.Second)

	log.Println("Scenario 2: Sell AAPL 300 GTC")
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideSell, Qty: "1", OrdType: fix.OrdTypeLimit, LimitPrice: "300", TIF: fix.TimeInForceGTC})
	time.Sleep(1 * time.Second)

	log.Println("Scenario 3: Sell Short QQQ 650")
	app.SendOrder(fix.OrderParams{Symbol: "QQQ", Side: fix.SideSellShort, Qty: "1", OrdType: fix.OrdTypeLimit, LimitPrice: "650", TIF: fix.TimeInForceDay})
	time.Sleep(1 * time.Second)

	// --- ORDER TYPE BATCH ---
	log.Println("--- Group: ORDER TYPE ---")

	log.Println("Scenario 4: MKT - Buy 1 AAPL")
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeMarket, TIF: fix.TimeInForceDay})
	time.Sleep(1 * time.Second)

	log.Println("Scenario 5: Limit - Sell Short 1 QQQ @ 650")
	app.SendOrder(fix.OrderParams{Symbol: "QQQ", Side: fix.SideSellShort, Qty: "1", OrdType: fix.OrdTypeLimit, LimitPrice: "650", TIF: fix.TimeInForceDay})
	time.Sleep(1 * time.Second)

	log.Println("Scenario 6: Stop - Buy 1 BRK.B @ 500")
	app.SendOrder(fix.OrderParams{Symbol: "BRK.B", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeStop, StopPrice: "500", TIF: fix.TimeInForceDay})
	time.Sleep(1 * time.Second)

	log.Println("Scenario 7: Stop Limit - Buy 1 BABA @ 150 STP 150 LMT")
	app.SendOrder(fix.OrderParams{Symbol: "BABA", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeStopLimit, StopPrice: "150", LimitPrice: "150", TIF: fix.TimeInForceDay})
	time.Sleep(1 * time.Second)

	log.Println("Scenario 8: Market on close")
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeMarketOnClose, TIF: fix.TimeInForceDay})
	time.Sleep(1 * time.Second)

	log.Println("Scenario 9: Limit on close")
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeLimitOnClose, LimitPrice: "150", TIF: fix.TimeInForceDay})
	time.Sleep(1 * time.Second)

	// --- EXEC INST BATCH ---
	log.Println("--- Group: EXEC INST ---")

	log.Println("Scenario 10: Not Held - Buy 100 BKNG MKT")
	app.SendOrder(fix.OrderParams{Symbol: "BKNG", Side: fix.SideBuy, Qty: "100", OrdType: fix.OrdTypeMarket, TIF: fix.TimeInForceDay, ExecInst: fix.ExecInstNotHeld})
	time.Sleep(1 * time.Second)

	log.Println("Scenario 11: Held - Buy 1 AAPL MKT")
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeMarket, TIF: fix.TimeInForceDay, ExecInst: fix.ExecInstHeld})
	time.Sleep(1 * time.Second)

	// --- SYMBOLOGY BATCH ---
	log.Println("--- Group: SYMBOLOGY ---")

	log.Println("Scenario 12: Symbol Suffix - Buy 1 BRK.B @ 500 STP")
	app.SendOrder(fix.OrderParams{Symbol: "BRK.B", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeStop, StopPrice: "500", TIF: fix.TimeInForceDay})
	time.Sleep(1 * time.Second)

	// --- TIF BATCH ---
	log.Println("--- Group: TIF ---")

	log.Println("Scenario 13: TIF Day - Buy 1 AAPL MKT")
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeMarket, TIF: fix.TimeInForceDay})
	time.Sleep(1 * time.Second)

	log.Println("Scenario 14: TIF GTC - Sell 1 AAPL 300 GTC")
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideSell, Qty: "1", OrdType: fix.OrdTypeLimit, LimitPrice: "300", TIF: fix.TimeInForceGTC})
	time.Sleep(1 * time.Second)

	log.Println("Scenario 15: TIF IOC - Buy 1 AAPL MKT IOC")
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeMarket, TIF: fix.TimeInForceIOC})
	time.Sleep(1 * time.Second)

	log.Println("Scenario 16: TIF FOK - Buy 1 AAPL MKT FOK")
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeMarket, TIF: fix.TimeInForceFOK})
	time.Sleep(1 * time.Second)

	// --- SETTLEMENT BATCH ---
	log.Println("--- Group: SETTLEMENT ---")

	log.Println("Scenario 17: Regular")
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeMarket, TIF: fix.TimeInForceDay, SettlTyp: fix.SettlmntTypRegular})
	time.Sleep(1 * time.Second)

	log.Println("Scenario 18: Cash (Expected: REJECT)")
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeMarket, TIF: fix.TimeInForceDay, SettlTyp: fix.SettlmntTypCash})
	time.Sleep(1 * time.Second)

	log.Println("Scenario 19: Next Day (Expected: REJECT)")
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeMarket, TIF: fix.TimeInForceDay, SettlTyp: fix.SettlmntTypNextDay})
	time.Sleep(1 * time.Second)

	log.Println("Scenario 20: T+2 (Expected: REJECT)")
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeMarket, TIF: fix.TimeInForceDay, SettlTyp: fix.SettlmntTypTplus2})
	time.Sleep(1 * time.Second)

	log.Println("Scenario 21: T+3 (Expected: REJECT)")
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeMarket, TIF: fix.TimeInForceDay, SettlTyp: fix.SettlmntTypTplus3})
	time.Sleep(1 * time.Second)

	log.Println("Scenario 22: T+4 (Expected: REJECT)")
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeMarket, TIF: fix.TimeInForceDay, SettlTyp: fix.SettlmntTypTplus4})
	time.Sleep(1 * time.Second)

	log.Println("Scenario 23: Future (Expected: REJECT)")
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeMarket, TIF: fix.TimeInForceDay, SettlTyp: fix.SettlmntTypFuture})
	time.Sleep(1 * time.Second)

	log.Println("Scenario 24: When Issued (Expected: REJECT)")
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeMarket, TIF: fix.TimeInForceDay, SettlTyp: fix.SettlmntTypWhenIssued})
	time.Sleep(1 * time.Second)

	log.Println("Scenario 25: Sellers Option (Expected: REJECT)")
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeMarket, TIF: fix.TimeInForceDay, SettlTyp: fix.SettlmntTypSellersOption})
	time.Sleep(1 * time.Second)

	log.Println("Scenario 26: T+5 (Expected: REJECT)")
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeMarket, TIF: fix.TimeInForceDay, SettlTyp: fix.SettlmntTypTplus5})
	time.Sleep(1 * time.Second)

	// --- MISC BATCH ---
	log.Println("--- Group: MISC SCENARIOS ---")

	log.Println("Scenario 27: Partial Fill & Cancel - Buy 1600 LCID MKT")
	app.SendOrder(fix.OrderParams{Symbol: "LCID", Side: fix.SideBuy, Qty: "1600", OrdType: fix.OrdTypeMarket})
	time.Sleep(1 * time.Second)
	app.SendCancelOrder()
	time.Sleep(1 * time.Second)

	log.Println("Scenario 28: Full Fill - Buy 1 AAPL MKT")
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeMarket})
	time.Sleep(1 * time.Second)

	log.Println("Scenario 29: Cancel an Acknowledged Order - Buy 1 BRK.B 500 STP")
	app.SendOrder(fix.OrderParams{Symbol: "BRK.B", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeStop, StopPrice: "500"})
	time.Sleep(1 * time.Second)
	app.SendCancelOrder()
	time.Sleep(1 * time.Second)

	log.Println("Scenario 30: Increase Quantity and Price - Replace to 3 shares @ 450")
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideBuy, Qty: "2", OrdType: fix.OrdTypeLimit, LimitPrice: "450"})
	time.Sleep(1 * time.Second)
	app.SendReplaceOrder("3", fix.OrdTypeLimit, "450")
	time.Sleep(1 * time.Second)

	log.Println("Scenario 31: Reject an order - Buy 1 LCID @ 2500 LMT")
	app.SendOrder(fix.OrderParams{Symbol: "LCID", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeLimit, LimitPrice: "2500"})
	time.Sleep(1 * time.Second)

	log.Println("Scenario 32: Reject Cxl/Replace Request - Buy 3 LCID @ 5 LMT")
	app.SendOrder(fix.OrderParams{Symbol: "LCID", Side: fix.SideBuy, Qty: "3", OrdType: fix.OrdTypeLimit, LimitPrice: "5"})
	time.Sleep(1 * time.Second)
	app.SendReplaceOrder("3", fix.OrdTypeLimit, "5")
	time.Sleep(1 * time.Second)

	log.Println("Scenario 33: Market to Limit - Replace to 3.00 LMT")
	app.SendOrder(fix.OrderParams{Symbol: "LCID", Side: fix.SideBuy, Qty: "2600", OrdType: fix.OrdTypeMarket})
	time.Sleep(1 * time.Second)
	app.SendReplaceOrder("2600", fix.OrdTypeLimit, "3.00")
	time.Sleep(1 * time.Second)

	log.Println("Scenario 34: Limit to Market - Replace to MKT")
	app.SendOrder(fix.OrderParams{Symbol: "LCID", Side: fix.SideBuy, Qty: "5", OrdType: fix.OrdTypeLimit, LimitPrice: "5.02"})
	time.Sleep(1 * time.Second)
	app.SendReplaceOrder("5", fix.OrdTypeMarket, "")
	time.Sleep(1 * time.Second)

	log.Println("Scenario 35: Done For Day - Coordinate with RQD")
	app.SendOrder(fix.OrderParams{Symbol: "GRO", Side: fix.SideBuy, Qty: "2600", OrdType: fix.OrdTypeLimit, LimitPrice: "3.00"})
	time.Sleep(1 * time.Second)

	log.Println("Scenario 36: Unsolicited Cancel - Coordinate with RQD")
	app.SendOrder(fix.OrderParams{Symbol: "GRO", Side: fix.SideBuy, Qty: "2700", OrdType: fix.OrdTypeLimit, LimitPrice: "2.99"})
	time.Sleep(1 * time.Second)

	log.Println("Scenario 37: ExecTransType New (Tag 20=0)")
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeMarket})
	time.Sleep(1 * time.Second)

	// --- FRACTIONAL BATCH ---
	log.Println("--- Group: FRACTIONAL ---")

	log.Println("Scenario 38: Fractional MKT Buy .25 AMZN")
	app.SendOrder(fix.OrderParams{Symbol: "AMZN", Side: fix.SideBuy, Qty: ".25", OrdType: fix.OrdTypeMarket, IsFractional: true})
	time.Sleep(1 * time.Second)

	log.Println("Scenario 39: Fractional LMT Buy .99 TSLA @ 400")
	app.SendOrder(fix.OrderParams{Symbol: "TSLA", Side: fix.SideBuy, Qty: ".99", OrdType: fix.OrdTypeLimit, LimitPrice: "400", IsFractional: true})
	time.Sleep(1 * time.Second)

	log.Println("Scenario 40: Fractional GTC Buy .50 TSLA @ 401 LMT GTC")
	app.SendOrder(fix.OrderParams{Symbol: "TSLA", Side: fix.SideBuy, Qty: ".50", OrdType: fix.OrdTypeLimit, LimitPrice: "401", TIF: fix.TimeInForceGTC, IsFractional: true})
	time.Sleep(1 * time.Second)

	log.Println("Scenario 41: Cancel Fractional Order")
	app.SendOrder(fix.OrderParams{Symbol: "TSLA", Side: fix.SideBuy, Qty: ".50", OrdType: fix.OrdTypeLimit, LimitPrice: "401", TIF: fix.TimeInForceGTC, IsFractional: true})
	time.Sleep(1 * time.Second)
	app.SendCancelOrder()
	time.Sleep(1 * time.Second)

	log.Println("Scenario 42: Replace Fractional Order")
	app.SendOrder(fix.OrderParams{Symbol: "TSLA", Side: fix.SideBuy, Qty: ".10", OrdType: fix.OrdTypeLimit, LimitPrice: "402", IsFractional: true})
	time.Sleep(1 * time.Second)
	app.SendReplaceOrder(".10", fix.OrdTypeLimit, "403")
	time.Sleep(1 * time.Second)

	// --- NOTIONAL BATCH ---
	log.Println("--- Group: NOTIONAL ---")

	log.Println("Scenario 43: Notional MKT Buy $100 TSLA")
	app.SendOrder(fix.OrderParams{Symbol: "TSLA", Side: fix.SideBuy, Notional: "100", OrdType: fix.OrdTypeMarket})
	time.Sleep(1 * time.Second)

	log.Println("Scenario 44: Notional LMT Buy $135 TSLA @ 500")
	app.SendOrder(fix.OrderParams{Symbol: "TSLA", Side: fix.SideBuy, Notional: "135", OrdType: fix.OrdTypeLimit, LimitPrice: "500"})
	time.Sleep(1 * time.Second)

	log.Println("Scenario 45: Notional GTC Buy $150 TSLA @ 450 LMT")
	app.SendOrder(fix.OrderParams{Symbol: "TSLA", Side: fix.SideBuy, Notional: "150", OrdType: fix.OrdTypeLimit, LimitPrice: "450", TIF: fix.TimeInForceGTC})
	time.Sleep(1 * time.Second)

	log.Println("Scenario 46: Cancel Notional Order")
	app.SendOrder(fix.OrderParams{Symbol: "TSLA", Side: fix.SideBuy, Notional: "175", OrdType: fix.OrdTypeMarket})
	time.Sleep(1 * time.Second)
	app.SendCancelOrder()
	time.Sleep(1 * time.Second)

	log.Println("Scenario 47: Replace Notional Order")
	app.SendOrder(fix.OrderParams{Symbol: "TSLA", Side: fix.SideBuy, Notional: "180", OrdType: fix.OrdTypeMarket})
	time.Sleep(1 * time.Second)
	app.SendReplaceOrder("180", fix.OrdTypeMarket, "")
	time.Sleep(1 * time.Second)

	// --- EXTENDED HOURS BATCH ---
	log.Println("--- Group: EXTENDED HOURS ---")

	log.Println("Scenario 48: AM Only DIA")
	app.SendOrder(fix.OrderParams{Symbol: "DIA", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeLimit, LimitPrice: "450", TradingSes: fix.TradingSessionAM})
	time.Sleep(1 * time.Second)

	log.Println("Scenario 49: PM Only DIA")
	app.SendOrder(fix.OrderParams{Symbol: "DIA", Side: fix.SideBuy, Qty: "2", OrdType: fix.OrdTypeLimit, LimitPrice: "451", TradingSes: fix.TradingSessionPM})
	time.Sleep(1 * time.Second)

	log.Println("Scenario 50: All Sessions DIA")
	app.SendOrder(fix.OrderParams{Symbol: "DIA", Side: fix.SideBuy, Qty: "3", OrdType: fix.OrdTypeLimit, LimitPrice: "452", TradingSes: fix.TradingSessionBoth})

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
