package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"thisgofix/internal/fix"

	"github.com/quickfixgo/quickfix"
)

// testDelay is the duration to pause between sending test orders in automatic mode.
const testDelay = 2 * time.Second

// Global variable for interactive mode, set via flag.
var interactiveMode bool

func main() {
	// 1. Define command line flags
	cfgFileName := flag.String("cfg", "config.cfg", "Path to config file")
	flag.BoolVar(&interactiveMode, "interactive", false, "Wait for ENTER between each test case")
	flag.Parse()

	// 2. Open the configuration file
	cfg, err := os.Open(*cfgFileName)
	if err != nil {
		log.Fatalf("Critical Error: Could not open config file '%s': %s", *cfgFileName, err)
	}
	defer func() {
		if closeErr := cfg.Close(); closeErr != nil {
			log.Printf("Warning: Failed to close config file: %s", closeErr)
		}
	}()

	appSettings, err := quickfix.ParseSettings(cfg)
	if err != nil {
		log.Fatalf("Critical Error: Config file format is invalid: %s", err)
	}

	// 3. Setup Application in TEST MODE
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

	// 4. Wait for Logon before starting tests
	log.Println("TestRunner: Waiting for FIX Logon...")
	select {
	case <-app.OnLogonChan:
		log.Println("TestRunner: Logged on. Starting scenarios...")
	case <-time.After(10 * time.Second):
		log.Fatalf("TestRunner: Timeout waiting for logon")
	}

	// 5. Execute Scenarios (Categorized)
	reader := bufio.NewReader(os.Stdin)

	testSides(app, reader)
	testOrderTypes(app, reader)
	testExecInst(app, reader)
	testSymbology(app, reader)
	testTIF(app, reader)
	testSettlement(app, reader)
	testFractional(app, reader)
	testNotional(app, reader)
	testExtendedHours(app, reader)
	testOptions(app, reader) // New options test group
	testOptionsRQD(app, reader)
	testOptionsComplexRQD(app, reader)
	testMisc(app, reader)

	// 6. Keep running for a few seconds to see final responses
	log.Println("TestRunner: All scenarios sent. Waiting for responses...")
	time.Sleep(10 * time.Second)

	log.Println("TestRunner: Stopping FIX Client...")
	initiator.Stop()
}

// waitNext handles the pause between cases, supporting both automatic delay and interactive mode.
func waitNext(reader *bufio.Reader, scenarioName string) {
	if interactiveMode {
		log.Println(scenarioName)
		fmt.Printf("[INTERACTIVE] Press ENTER to run next step...")
		_, _ = reader.ReadString('\n')
	} else {
		time.Sleep(testDelay)
		log.Println(scenarioName)
	}
}

func testSides(app *fix.Application, r *bufio.Reader) {
	waitNext(r, "--- Group: SIDES ---")
	waitNext(r, "Scenario 1: Buy AAPL MKT")
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeMarket, TIF: fix.TimeInForceDay})

	waitNext(r, "Scenario 2: Sell AAPL 300 GTC")
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideSell, Qty: "1", OrdType: fix.OrdTypeLimit, LimitPrice: "300", TIF: fix.TimeInForceGTC})

	waitNext(r, "Scenario 3: Sell Short QQQ 650")
	app.SendOrder(fix.OrderParams{Symbol: "QQQ", Side: fix.SideSellShort, Qty: "1", OrdType: fix.OrdTypeLimit, LimitPrice: "650", TIF: fix.TimeInForceDay})
}

func testOrderTypes(app *fix.Application, r *bufio.Reader) {
	waitNext(r, "--- Group: ORDER TYPE ---")
	waitNext(r, "Scenario 4: MKT - Buy 1 AAPL")
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeMarket, TIF: fix.TimeInForceDay})

	waitNext(r, "Scenario 5: Stop - Buy 1 BRK.B @ 500")
	app.SendOrder(fix.OrderParams{Symbol: "BRK.B", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeStop, StopPrice: "500", TIF: fix.TimeInForceDay})

	waitNext(r, "Scenario 6: Stop Limit - Buy 1 BABA @ 150 STP 150 LMT")
	app.SendOrder(fix.OrderParams{Symbol: "BABA", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeStopLimit, StopPrice: "150", LimitPrice: "150", TIF: fix.TimeInForceDay})

	waitNext(r, "Scenario 7: Market on close")
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeMarketOnClose, TIF: fix.TimeInForceDay})

	waitNext(r, "Scenario 8: Limit on close")
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeLimitOnClose, LimitPrice: "150", TIF: fix.TimeInForceDay})
}

func testExecInst(app *fix.Application, r *bufio.Reader) {
	waitNext(r, "--- Group: EXEC INST ---")
	waitNext(r, "Scenario 10: Not Held - Buy 100 BKNG MKT")
	app.SendOrder(fix.OrderParams{Symbol: "BKNG", Side: fix.SideBuy, Qty: "100", OrdType: fix.OrdTypeMarket, TIF: fix.TimeInForceDay, ExecInst: fix.ExecInstNotHeld})

	waitNext(r, "Scenario 11: Held - Buy 1 AAPL MKT")
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeMarket, TIF: fix.TimeInForceDay, ExecInst: fix.ExecInstHeld})
}

func testSymbology(app *fix.Application, r *bufio.Reader) {
	waitNext(r, "--- Group: SYMBOLOGY ---")
	waitNext(r, "Scenario 12: Symbol Suffix - Buy 1 BRK.B @ 500 STP")
	app.SendOrder(fix.OrderParams{Symbol: "BRK.B", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeStop, StopPrice: "500", TIF: fix.TimeInForceDay})
}

func testTIF(app *fix.Application, r *bufio.Reader) {
	waitNext(r, "--- Group: TIF ---")
	waitNext(r, "Scenario 13: TIF Day - Buy 1 AAPL MKT")
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeMarket, TIF: fix.TimeInForceDay})

	waitNext(r, "Scenario 14: TIF GTC - Sell 1 AAPL 300 GTC")
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideSell, Qty: "1", OrdType: fix.OrdTypeLimit, LimitPrice: "300", TIF: fix.TimeInForceGTC})

	waitNext(r, "Scenario 15: TIF IOC - Buy 1 AAPL MKT IOC")
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeMarket, TIF: fix.TimeInForceIOC})

	waitNext(r, "Scenario 16: TIF FOK - Buy 1 AAPL MKT FOK")
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeMarket, TIF: fix.TimeInForceFOK})
}

func testSettlement(app *fix.Application, r *bufio.Reader) {
	waitNext(r, "--- Group: SETTLEMENT ---")
	waitNext(r, "Scenario 17: Regular")
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeMarket, TIF: fix.TimeInForceDay, SettlTyp: fix.SettlmntTypRegular})

	waitNext(r, "Scenario 18: Cash (Expected: REJECT)")
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeMarket, TIF: fix.TimeInForceDay, SettlTyp: fix.SettlmntTypCash})

	waitNext(r, "Scenario 19: Next Day (Expected: REJECT)")
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeMarket, TIF: fix.TimeInForceDay, SettlTyp: fix.SettlmntTypNextDay})

	waitNext(r, "Scenario 20: T+2 (Expected: REJECT)")
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeMarket, TIF: fix.TimeInForceDay, SettlTyp: fix.SettlmntTypTplus2})

	waitNext(r, "Scenario 21: T+3 (Expected: REJECT)")
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeMarket, TIF: fix.TimeInForceDay, SettlTyp: fix.SettlmntTypTplus3})

	waitNext(r, "Scenario 22: T+4 (Expected: REJECT)")
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeMarket, TIF: fix.TimeInForceDay, SettlTyp: fix.SettlmntTypTplus4})

	waitNext(r, "Scenario 23: Future (Expected: REJECT)")
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeMarket, TIF: fix.TimeInForceDay, SettlTyp: fix.SettlmntTypFuture})

	waitNext(r, "Scenario 24: When Issued (Expected: REJECT)")
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeMarket, TIF: fix.TimeInForceDay, SettlTyp: fix.SettlmntTypWhenIssued})

	waitNext(r, "Scenario 25: Sellers Option (Expected: REJECT)")
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeMarket, TIF: fix.TimeInForceDay, SettlTyp: fix.SettlmntTypSellersOption})

	waitNext(r, "Scenario 26: T+5 (Expected: REJECT)")
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeMarket, TIF: fix.TimeInForceDay, SettlTyp: fix.SettlmntTypTplus5})
}

func testFractional(app *fix.Application, r *bufio.Reader) {
	waitNext(r, "--- Group: FRACTIONAL ---")
	waitNext(r, "Scenario 38: Fractional MKT Buy .25 AMZN")
	app.SendOrder(fix.OrderParams{Symbol: "AMZN", Side: fix.SideBuy, Qty: ".25", OrdType: fix.OrdTypeMarket, IsFractional: true})

	waitNext(r, "Scenario 39: Fractional LMT Buy .99 TSLA @ 400")
	app.SendOrder(fix.OrderParams{Symbol: "TSLA", Side: fix.SideBuy, Qty: ".99", OrdType: fix.OrdTypeLimit, LimitPrice: "400", IsFractional: true})

	waitNext(r, "Scenario 40: Fractional GTC Buy .50 TSLA @ 401 LMT GTC")
	app.SendOrder(fix.OrderParams{Symbol: "TSLA", Side: fix.SideBuy, Qty: ".50", OrdType: fix.OrdTypeLimit, LimitPrice: "401", TIF: fix.TimeInForceGTC, IsFractional: true})

	waitNext(r, "Scenario 41: Cancel Fractional Order")
	app.SendOrder(fix.OrderParams{Symbol: "TSLA", Side: fix.SideBuy, Qty: ".50", OrdType: fix.OrdTypeLimit, LimitPrice: "401", TIF: fix.TimeInForceGTC, IsFractional: true})
	waitNext(r, "Confirm Cancel Fractional")
	app.SendCancelOrder()

	waitNext(r, "Scenario 42: Replace Fractional Order")
	app.SendOrder(fix.OrderParams{Symbol: "TSLA", Side: fix.SideBuy, Qty: ".10", OrdType: fix.OrdTypeLimit, LimitPrice: "402", IsFractional: true})
	waitNext(r, "Confirm Replace Fractional")
	app.SendReplaceOrder(".10", fix.OrdTypeLimit, "403")
}

func testNotional(app *fix.Application, r *bufio.Reader) {
	waitNext(r, "--- Group: NOTIONAL ---")
	waitNext(r, "Scenario 43: Notional MKT Buy $100 TSLA")
	app.SendOrder(fix.OrderParams{Symbol: "TSLA", Side: fix.SideBuy, Notional: "100", OrdType: fix.OrdTypeMarket})

	waitNext(r, "Scenario 44: Notional LMT Buy $135 TSLA @ 500")
	app.SendOrder(fix.OrderParams{Symbol: "TSLA", Side: fix.SideBuy, Notional: "135", OrdType: fix.OrdTypeLimit, LimitPrice: "500"})

	waitNext(r, "Scenario 45: Notional GTC Buy $150 TSLA @ 450 LMT")
	app.SendOrder(fix.OrderParams{Symbol: "TSLA", Side: fix.SideBuy, Notional: "150", OrdType: fix.OrdTypeLimit, LimitPrice: "450", TIF: fix.TimeInForceGTC})

	waitNext(r, "Scenario 46: Cancel Notional Order")
	app.SendOrder(fix.OrderParams{Symbol: "TSLA", Side: fix.SideBuy, Notional: "175", OrdType: fix.OrdTypeMarket})
	waitNext(r, "Confirm Cancel Notional")
	app.SendCancelOrder()

	waitNext(r, "Scenario 47: Replace Notional Order")
	app.SendOrder(fix.OrderParams{Symbol: "TSLA", Side: fix.SideBuy, Notional: "180", OrdType: fix.OrdTypeMarket})
	waitNext(r, "Confirm Replace Notional")
	app.SendReplaceOrder("180", fix.OrdTypeMarket, "")
}

func testExtendedHours(app *fix.Application, r *bufio.Reader) {
	waitNext(r, "--- Group: EXTENDED HOURS ---")
	waitNext(r, "Scenario 48: AM Only DIA")
	app.SendOrder(fix.OrderParams{Symbol: "DIA", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeLimit, LimitPrice: "450", TradingSes: fix.TradingSessionAM})

	waitNext(r, "Scenario 49: PM Only DIA")
	app.SendOrder(fix.OrderParams{Symbol: "DIA", Side: fix.SideBuy, Qty: "2", OrdType: fix.OrdTypeLimit, LimitPrice: "451", TradingSes: fix.TradingSessionPM})

	waitNext(r, "Scenario 50: All Sessions DIA")
	app.SendOrder(fix.OrderParams{Symbol: "DIA", Side: fix.SideBuy, Qty: "3", OrdType: fix.OrdTypeLimit, LimitPrice: "452", TradingSes: fix.TradingSessionBoth})
}

func testOptions(app *fix.Application, r *bufio.Reader) {
	waitNext(r, "--- Group: OPTIONS ---")
	waitNext(r, "Scenario 51: Buy META Call Option (20260522, Strike 595)")
	app.SendOrder(fix.OrderParams{
		Symbol:             "META",
		SecurityType:       fix.SecurityTypeOption,
		MaturityMonthYear:  "202605",
		MaturityDay:        "22",
		PutOrCall:          fix.PutOrCallCall,
		StrikePrice:        "595",
		ContractMultiplier: "100",
		Side:               fix.SideBuy,
		Qty:                "1",
		OrdType:            fix.OrdTypeMarket,
		TIF:                fix.TimeInForceDay,
		TradingSes:         fix.TradingSessionBoth,
		OpenClose:          fix.OpenCloseOpen,
	})
	waitNext(r, "Scenario 52: Sell META Put Option (20260522, Strike 595)")
	app.SendOrder(fix.OrderParams{
		Symbol:             "META",
		SecurityType:       fix.SecurityTypeOption,
		MaturityMonthYear:  "202605",
		MaturityDay:        "22",
		PutOrCall:          fix.PutOrCallPut,
		StrikePrice:        "595",
		ContractMultiplier: "100",
		Side:               fix.SideSell,
		Qty:                "1",
		OrdType:            fix.OrdTypeLimit,
		LimitPrice:         "1.50",
		TIF:                fix.TimeInForceDay,
		TradingSes:         fix.TradingSessionBoth,
		OpenClose:          fix.OpenCloseClose,
	})
}

func testOptionsRQD(app *fix.Application, r *bufio.Reader) {
	waitNext(r, "--- Group: OPTIONS RQD ---")

	waitNext(r, "Scenario 53 / S-001: BTO AAPL 20260717 270 Call MKT DAY")
	app.SendOrder(fix.OrderParams{
		Symbol:             "AAPL",
		SecurityType:       fix.SecurityTypeOption,
		MaturityMonthYear:  "202607",
		MaturityDay:        "17",
		PutOrCall:          fix.PutOrCallCall,
		StrikePrice:        "270",
		ContractMultiplier: "100",
		Side:               fix.SideBuy,
		Qty:                "1",
		OrdType:            fix.OrdTypeMarket,
		TIF:                fix.TimeInForceDay,
		TradingSes:         fix.TradingSessionBoth,
		OpenClose:          fix.OpenCloseOpen,
	})

	waitNext(r, "Scenario 54 / S-002: BTO AAPL 20260717 275 Call LMT 4.20 DAY")
	app.SendOrder(fix.OrderParams{
		Symbol:             "AAPL",
		SecurityType:       fix.SecurityTypeOption,
		MaturityMonthYear:  "202607",
		MaturityDay:        "17",
		PutOrCall:          fix.PutOrCallCall,
		StrikePrice:        "275",
		ContractMultiplier: "100",
		Side:               fix.SideBuy,
		Qty:                "1",
		OrdType:            fix.OrdTypeLimit,
		LimitPrice:         "4.20",
		TIF:                fix.TimeInForceDay,
		TradingSes:         fix.TradingSessionBoth,
		OpenClose:          fix.OpenCloseOpen,
	})

	waitNext(r, "Scenario 55 / S-003: BTO NVDA 20260717 215 Call LMT 8.50 GTC")
	app.SendOrder(fix.OrderParams{
		Symbol:             "NVDA",
		SecurityType:       fix.SecurityTypeOption,
		MaturityMonthYear:  "202607",
		MaturityDay:        "17",
		PutOrCall:          fix.PutOrCallCall,
		StrikePrice:        "215",
		ContractMultiplier: "100",
		Side:               fix.SideBuy,
		Qty:                "5",
		OrdType:            fix.OrdTypeLimit,
		LimitPrice:         "8.50",
		TIF:                fix.TimeInForceGTC,
		TradingSes:         fix.TradingSessionBoth,
		OpenClose:          fix.OpenCloseOpen,
	})

	waitNext(r, "Scenario 56 / S-004: BTO NVDA 20260717 215 Call STP 19 DAY")
	app.SendOrder(fix.OrderParams{
		Symbol:             "NVDA",
		SecurityType:       fix.SecurityTypeOption,
		MaturityMonthYear:  "202607",
		MaturityDay:        "17",
		PutOrCall:          fix.PutOrCallCall,
		StrikePrice:        "215",
		ContractMultiplier: "100",
		Side:               fix.SideBuy,
		Qty:                "1",
		OrdType:            fix.OrdTypeStop,
		StopPrice:          "19",
		TIF:                fix.TimeInForceDay,
		TradingSes:         fix.TradingSessionBoth,
		OpenClose:          fix.OpenCloseOpen,
	})

	waitNext(r, "Scenario 57 / S-005: BTO TSLA 20260717 380 Call STPLMT 6.50/6.75 DAY")
	app.SendOrder(fix.OrderParams{
		Symbol:             "TSLA",
		SecurityType:       fix.SecurityTypeOption,
		MaturityMonthYear:  "202607",
		MaturityDay:        "17",
		PutOrCall:          fix.PutOrCallCall,
		StrikePrice:        "380",
		ContractMultiplier: "100",
		Side:               fix.SideBuy,
		Qty:                "1",
		OrdType:            fix.OrdTypeStopLimit,
		StopPrice:          "6.50",
		LimitPrice:         "6.75",
		TIF:                fix.TimeInForceDay,
		TradingSes:         fix.TradingSessionBoth,
		OpenClose:          fix.OpenCloseOpen,
	})

	waitNext(r, "Scenario 58 / S-006: BTO AAPL 20260717 270 Call MKT DAY Qty 2")
	app.SendOrder(fix.OrderParams{
		Symbol:             "AAPL",
		SecurityType:       fix.SecurityTypeOption,
		MaturityMonthYear:  "202607",
		MaturityDay:        "17",
		PutOrCall:          fix.PutOrCallCall,
		StrikePrice:        "270",
		ContractMultiplier: "100",
		Side:               fix.SideBuy,
		Qty:                "2",
		OrdType:            fix.OrdTypeMarket,
		TIF:                fix.TimeInForceDay,
		TradingSes:         fix.TradingSessionBoth,
		OpenClose:          fix.OpenCloseOpen,
	})

	waitNext(r, "Scenario 59 / S-007: BTO NVDA 20260717 220 Call LMT 5 GTC Qty 10")
	app.SendOrder(fix.OrderParams{
		Symbol:             "NVDA",
		SecurityType:       fix.SecurityTypeOption,
		MaturityMonthYear:  "202607",
		MaturityDay:        "17",
		PutOrCall:          fix.PutOrCallCall,
		StrikePrice:        "220",
		ContractMultiplier: "100",
		Side:               fix.SideBuy,
		Qty:                "10",
		OrdType:            fix.OrdTypeLimit,
		LimitPrice:         "5",
		TIF:                fix.TimeInForceGTC,
		TradingSes:         fix.TradingSessionBoth,
		OpenClose:          fix.OpenCloseOpen,
	})

	waitNext(r, "Scenario 60 / S-008: BTO TSLA 20260717 380 Call LMT 0.25 DAY then cancel")
	app.SendOrder(fix.OrderParams{
		Symbol:             "TSLA",
		SecurityType:       fix.SecurityTypeOption,
		MaturityMonthYear:  "202607",
		MaturityDay:        "17",
		PutOrCall:          fix.PutOrCallCall,
		StrikePrice:        "380",
		ContractMultiplier: "100",
		Side:               fix.SideBuy,
		Qty:                "1",
		OrdType:            fix.OrdTypeLimit,
		LimitPrice:         "0.25",
		TIF:                fix.TimeInForceDay,
		TradingSes:         fix.TradingSessionBoth,
		OpenClose:          fix.OpenCloseOpen,
	})
	waitNext(r, "Scenario 60 / S-008: Cancel TSLA option order")
	app.SendCancelOrder()

	waitNext(r, "Scenario 61 / S-009: STO NVDA 20260717 215 Put MKT DAY")
	app.SendOrder(fix.OrderParams{
		Symbol:             "NVDA",
		SecurityType:       fix.SecurityTypeOption,
		MaturityMonthYear:  "202607",
		MaturityDay:        "17",
		PutOrCall:          fix.PutOrCallPut,
		StrikePrice:        "215",
		ContractMultiplier: "100",
		Side:               fix.SideSell,
		Qty:                "1",
		OrdType:            fix.OrdTypeMarket,
		TIF:                fix.TimeInForceDay,
		TradingSes:         fix.TradingSessionBoth,
		OpenClose:          fix.OpenCloseOpen,
	})

	waitNext(r, "Scenario 62 / S-010: STO NVDA 20260717 215 Put LMT 4 DAY")
	app.SendOrder(fix.OrderParams{
		Symbol:             "NVDA",
		SecurityType:       fix.SecurityTypeOption,
		MaturityMonthYear:  "202607",
		MaturityDay:        "17",
		PutOrCall:          fix.PutOrCallPut,
		StrikePrice:        "215",
		ContractMultiplier: "100",
		Side:               fix.SideSell,
		Qty:                "1",
		OrdType:            fix.OrdTypeLimit,
		LimitPrice:         "4",
		TIF:                fix.TimeInForceDay,
		TradingSes:         fix.TradingSessionBoth,
		OpenClose:          fix.OpenCloseOpen,
	})

	waitNext(r, "Scenario 63 / S-011: STO TSLA 20260717 370 Put LMT 35 GTC")
	app.SendOrder(fix.OrderParams{
		Symbol:             "TSLA",
		SecurityType:       fix.SecurityTypeOption,
		MaturityMonthYear:  "202607",
		MaturityDay:        "17",
		PutOrCall:          fix.PutOrCallPut,
		StrikePrice:        "370",
		ContractMultiplier: "100",
		Side:               fix.SideSell,
		Qty:                "2",
		OrdType:            fix.OrdTypeLimit,
		LimitPrice:         "35",
		TIF:                fix.TimeInForceGTC,
		TradingSes:         fix.TradingSessionBoth,
		OpenClose:          fix.OpenCloseOpen,
	})

	waitNext(r, "Scenario 64 / S-012: STO SPY 20260717 705 Put MKT DAY Qty 5")
	app.SendOrder(fix.OrderParams{
		Symbol:             "SPY",
		SecurityType:       fix.SecurityTypeOption,
		MaturityMonthYear:  "202607",
		MaturityDay:        "17",
		PutOrCall:          fix.PutOrCallPut,
		StrikePrice:        "705",
		ContractMultiplier: "100",
		Side:               fix.SideSell,
		Qty:                "5",
		OrdType:            fix.OrdTypeMarket,
		TIF:                fix.TimeInForceDay,
		TradingSes:         fix.TradingSessionBoth,
		OpenClose:          fix.OpenCloseOpen,
	})

	waitNext(r, "Scenario 65 / S-013: BTC NVDA 20260717 215 Put MKT DAY")
	app.SendOrder(fix.OrderParams{
		Symbol:             "NVDA",
		SecurityType:       fix.SecurityTypeOption,
		MaturityMonthYear:  "202607",
		MaturityDay:        "17",
		PutOrCall:          fix.PutOrCallPut,
		StrikePrice:        "215",
		ContractMultiplier: "100",
		Side:               fix.SideBuy,
		Qty:                "1",
		OrdType:            fix.OrdTypeMarket,
		TIF:                fix.TimeInForceDay,
		TradingSes:         fix.TradingSessionBoth,
		OpenClose:          fix.OpenCloseClose,
	})

	waitNext(r, "Scenario 66 / S-014: STC AAPL 20260717 270 Call MKT DAY")
	app.SendOrder(fix.OrderParams{
		Symbol:             "AAPL",
		SecurityType:       fix.SecurityTypeOption,
		MaturityMonthYear:  "202607",
		MaturityDay:        "17",
		PutOrCall:          fix.PutOrCallCall,
		StrikePrice:        "270",
		ContractMultiplier: "100",
		Side:               fix.SideSell,
		Qty:                "1",
		OrdType:            fix.OrdTypeMarket,
		TIF:                fix.TimeInForceDay,
		TradingSes:         fix.TradingSessionBoth,
		OpenClose:          fix.OpenCloseClose,
	})
}

func testOptionsComplexRQD(app *fix.Application, r *bufio.Reader) {
	waitNext(r, "--- Group: OPTIONS COMPLEX RQD ---")

	waitNext(r, "Scenario 67: NVDA Debit Call Spread - Buy 250 Call / Sell 275 Call for 3.00 debit")
	app.SendMultilegOrder(fix.MultilegOrderParams{
		Symbol:     "NVDA",
		Side:       fix.SideBuy,
		Qty:        "1",
		OrdType:    fix.OrdTypeLimit,
		Price:      "3.00",
		TIF:        fix.TimeInForceDay,
		TradingSes: fix.TradingSessionBoth,
		Text:       "RQD Debit Spread: BTO 1 NVDA 20260717 250C, STO 1 NVDA 20260717 275C, net debit 3.00",
		Legs: []fix.MultilegLegParams{
			{
				LegRefID:           "BTO-250C",
				Symbol:             "NVDA",
				SecurityType:       fix.SecurityTypeOption,
				CFICode:            fix.LegCFICodeCall,
				MaturityMonthYear:  "202607",
				MaturityDate:       "20260717",
				StrikePrice:        "250",
				ContractMultiplier: "100",
				RatioQty:           "1",
				Qty:                "1",
				Side:               fix.SideBuy,
				PositionEffect:     fix.OpenCloseOpen,
			},
			{
				LegRefID:           "STO-275C",
				Symbol:             "NVDA",
				SecurityType:       fix.SecurityTypeOption,
				CFICode:            fix.LegCFICodeCall,
				MaturityMonthYear:  "202607",
				MaturityDate:       "20260717",
				StrikePrice:        "275",
				ContractMultiplier: "100",
				RatioQty:           "1",
				Qty:                "1",
				Side:               fix.SideSell,
				PositionEffect:     fix.OpenCloseOpen,
			},
		},
	})

	waitNext(r, "Scenario 68: AMZN Credit Put Spread - Sell 250 Put / Buy 240 Put for 4.00 credit")
	app.SendMultilegOrder(fix.MultilegOrderParams{
		Symbol:     "AMZN",
		Side:       fix.SideSell,
		Qty:        "1",
		OrdType:    fix.OrdTypeLimit,
		Price:      "4.00",
		TIF:        fix.TimeInForceDay,
		TradingSes: fix.TradingSessionBoth,
		Text:       "RQD Credit Spread: STO 1 AMZN 20260717 250P, BTO 1 AMZN 20260717 240P, net credit 4.00",
		Legs: []fix.MultilegLegParams{
			{
				LegRefID:           "STO-250P",
				Symbol:             "AMZN",
				SecurityType:       fix.SecurityTypeOption,
				CFICode:            fix.LegCFICodePut,
				MaturityMonthYear:  "202607",
				MaturityDate:       "20260717",
				StrikePrice:        "250",
				ContractMultiplier: "100",
				RatioQty:           "1",
				Qty:                "1",
				Side:               fix.SideSell,
				PositionEffect:     fix.OpenCloseOpen,
			},
			{
				LegRefID:           "BTO-240P",
				Symbol:             "AMZN",
				SecurityType:       fix.SecurityTypeOption,
				CFICode:            fix.LegCFICodePut,
				MaturityMonthYear:  "202607",
				MaturityDate:       "20260717",
				StrikePrice:        "240",
				ContractMultiplier: "100",
				RatioQty:           "1",
				Qty:                "1",
				Side:               fix.SideBuy,
				PositionEffect:     fix.OpenCloseOpen,
			},
		},
	})

	waitNext(r, "Scenario 69: XLY Covered Call - Buy 100 XLY / Sell 1 117 Call MKT")
	app.SendMultilegOrder(fix.MultilegOrderParams{
		Symbol:     "XLY",
		Side:       fix.SideBuy,
		Qty:        "1",
		OrdType:    fix.OrdTypeMarket,
		TIF:        fix.TimeInForceDay,
		TradingSes: fix.TradingSessionBoth,
		Text:       "RQD Covered Call: Buy 100 XLY, STO 1 XLY 20260717 117C at market",
		Legs: []fix.MultilegLegParams{
			{
				LegRefID:       "BUY-100-SH",
				Symbol:         "XLY",
				SecurityType:   fix.SecurityTypeCommon,
				RatioQty:       "100",
				Qty:            "100",
				Side:           fix.SideBuy,
				PositionEffect: fix.OpenCloseOpen,
			},
			{
				LegRefID:           "STO-117C",
				Symbol:             "XLY",
				SecurityType:       fix.SecurityTypeOption,
				CFICode:            fix.LegCFICodeCall,
				MaturityMonthYear:  "202607",
				MaturityDate:       "20260717",
				StrikePrice:        "117",
				ContractMultiplier: "100",
				RatioQty:           "1",
				Qty:                "1",
				Side:               fix.SideSell,
				PositionEffect:     fix.OpenCloseOpen,
			},
		},
	})

	waitNext(r, "Scenario 70: XLY Covered Call Roll - Buy close Jul 117 Call / Sell open Sep 117 Call for 4.50 debit")
	app.SendMultilegOrder(fix.MultilegOrderParams{
		Symbol:     "XLY",
		Side:       fix.SideBuy,
		Qty:        "1",
		OrdType:    fix.OrdTypeLimit,
		Price:      "4.50",
		TIF:        fix.TimeInForceDay,
		TradingSes: fix.TradingSessionBoth,
		Text:       "RQD Roll: BTC 1 XLY 20260717 117C, STO 1 XLY 20260918 117C, net debit 4.50",
		Legs: []fix.MultilegLegParams{
			{
				LegRefID:           "BTC-JUL-117C",
				Symbol:             "XLY",
				SecurityType:       fix.SecurityTypeOption,
				CFICode:            fix.LegCFICodeCall,
				MaturityMonthYear:  "202607",
				MaturityDate:       "20260717",
				StrikePrice:        "117",
				ContractMultiplier: "100",
				RatioQty:           "1",
				Qty:                "1",
				Side:               fix.SideBuy,
				PositionEffect:     fix.OpenCloseClose,
			},
			{
				LegRefID:           "STO-SEP-117C",
				Symbol:             "XLY",
				SecurityType:       fix.SecurityTypeOption,
				CFICode:            fix.LegCFICodeCall,
				MaturityMonthYear:  "202609",
				MaturityDate:       "20260918",
				StrikePrice:        "117",
				ContractMultiplier: "100",
				RatioQty:           "1",
				Qty:                "1",
				Side:               fix.SideSell,
				PositionEffect:     fix.OpenCloseOpen,
			},
		},
	})

	waitNext(r, "Scenario 71: EWJ Married Put - Buy 100 EWJ / Buy 1 92 Put for 96.00 debit")
	app.SendMultilegOrder(fix.MultilegOrderParams{
		Symbol:     "EWJ",
		Side:       fix.SideBuy,
		Qty:        "1",
		OrdType:    fix.OrdTypeLimit,
		Price:      "96.00",
		TIF:        fix.TimeInForceDay,
		TradingSes: fix.TradingSessionBoth,
		Text:       "RQD Married Put: Buy 100 EWJ, BTO 1 EWJ 20260717 92P, net debit 96.00",
		Legs: []fix.MultilegLegParams{
			{
				LegRefID:       "BUY-100-SH",
				Symbol:         "EWJ",
				SecurityType:   fix.SecurityTypeCommon,
				CFICode:        fix.LegCFICodeEquity,
				RatioQty:       "100",
				Qty:            "100",
				Side:           fix.SideBuy,
				PositionEffect: fix.OpenCloseOpen,
			},
			{
				LegRefID:           "BTO-92P",
				Symbol:             "EWJ",
				SecurityType:       fix.SecurityTypeOption,
				CFICode:            fix.LegCFICodePut,
				MaturityMonthYear:  "202607",
				MaturityDate:       "20260717",
				StrikePrice:        "92",
				ContractMultiplier: "100",
				RatioQty:           "1",
				Qty:                "1",
				Side:               fix.SideBuy,
				PositionEffect:     fix.OpenCloseOpen,
			},
		},
	})
}

func testMisc(app *fix.Application, r *bufio.Reader) {
	waitNext(r, "--- Group: MISC SCENARIOS ---")
	waitNext(r, "Scenario 27: Partial Fill & Cancel - Buy 1600 LCID MKT")
	app.SendOrder(fix.OrderParams{Symbol: "LCID", Side: fix.SideBuy, Qty: "1600", OrdType: fix.OrdTypeMarket})
	waitNext(r, "Confirm Cancel Partial")
	app.SendCancelOrder()

	waitNext(r, "Scenario 28: Full Fill")
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeMarket})

	waitNext(r, "Scenario 29: Cancel an Acknowledged Order")
	app.SendOrder(fix.OrderParams{Symbol: "BRK.B", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeStop, StopPrice: "500"})
	waitNext(r, "Confirm Cancel Ack'd")
	app.SendCancelOrder()

	waitNext(r, "Scenario 30: Increase Quantity and Price")
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideBuy, Qty: "2", OrdType: fix.OrdTypeLimit, LimitPrice: "450"})
	waitNext(r, "Confirm Replace (Increase)")
	app.SendReplaceOrder("3", fix.OrdTypeLimit, "450")

	waitNext(r, "Scenario 31: Reject an order")
	app.SendOrder(fix.OrderParams{Symbol: "LCID", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeLimit, LimitPrice: "2500"})

	waitNext(r, "Scenario 32: Reject Cxl/Replace Request")
	app.SendOrder(fix.OrderParams{Symbol: "LCID", Side: fix.SideBuy, Qty: "3", OrdType: fix.OrdTypeLimit, LimitPrice: "5"})
	waitNext(r, "Confirm Replace (Expect Reject)")
	app.SendReplaceOrder("3", fix.OrdTypeLimit, "5")

	waitNext(r, "Scenario 33: Market to Limit")
	app.SendOrder(fix.OrderParams{Symbol: "LCID", Side: fix.SideBuy, Qty: "2600", OrdType: fix.OrdTypeMarket})
	waitNext(r, "Confirm Replace (Mkt to Lmt)")
	app.SendReplaceOrder("2600", fix.OrdTypeLimit, "3.00")

	waitNext(r, "Scenario 34: Limit to Market")
	app.SendOrder(fix.OrderParams{Symbol: "LCID", Side: fix.SideBuy, Qty: "5", OrdType: fix.OrdTypeLimit, LimitPrice: "5.02"})
	waitNext(r, "Confirm Replace (Lmt to Mkt)")
	app.SendReplaceOrder("5", fix.OrdTypeMarket, "")

	waitNext(r, "Scenario 35: Done For Day")
	app.SendOrder(fix.OrderParams{Symbol: "GRO", Side: fix.SideBuy, Qty: "2600", OrdType: fix.OrdTypeLimit, LimitPrice: "3.00"})

	waitNext(r, "Scenario 36: Unsolicited Cancel")
	app.SendOrder(fix.OrderParams{Symbol: "GRO", Side: fix.SideBuy, Qty: "2700", OrdType: fix.OrdTypeLimit, LimitPrice: "2.99"})

	waitNext(r, "Scenario 37: ExecTransType New (Tag 20=0)")
	app.SendOrder(fix.OrderParams{Symbol: "AAPL", Side: fix.SideBuy, Qty: "1", OrdType: fix.OrdTypeMarket})
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
