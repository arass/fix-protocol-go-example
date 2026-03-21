package main

// -----------------------------------------------------------------------------
// MAIN ENTRY POINT
// -----------------------------------------------------------------------------
// This is the entry point for the application. It sets up the configuration
// and starts the FIX engine using the logic defined in the internal package.
// -----------------------------------------------------------------------------

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"thisgofix/internal/fix"

	"github.com/quickfixgo/quickfix"
)

func main() {
	// 1. Parse command line arguments
	//    We allow the user to specify a different config file if they want.
	//    Usage: go run cmd/thisgofix/main.go -cfg config.cfg
	//    (Note: relative path depends on where you run it from)
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

	// 4. Create our Application instance (logic defined in internal/fix)
	app := fix.NewApplication()

	// 5. Create Log Factory (Using Custom ConsoleLogFactory to avoid import issues)
	//    This will print all FIX messages to the console, which is great for debugging.
	logFactory := &ConsoleLogFactory{}

	// 6. Create Message Store Factory (Using MemoryStore to avoid compatibility issues)
	messageStoreFactory := quickfix.NewMemoryStoreFactory()

	// 7. Create the Initiator
	//    The Initiator connects to the server (Client -> Server).
	initiator, err := quickfix.NewInitiator(app, messageStoreFactory, appSettings, logFactory)
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

// -----------------------------------------------------------------------------
// CUSTOM LOG FACTORY
// -----------------------------------------------------------------------------
// We implement a custom logger to print to stdout, bypassing potential
// library version conflicts with built-in factories.

type ConsoleLogFactory struct{}

// Create creates a global log (not specific to a session, usually unused by Initiator)
func (f ConsoleLogFactory) Create() (quickfix.Log, error) {
	return &ConsoleLog{prefix: "GLOBAL"}, nil
}

// CreateSessionLog creates a log for a specific session
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
