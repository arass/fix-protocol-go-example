# Go FIX Client (Trading & Certification Tool)

This application is a FIX Protocol client built with Go and `quickfixgo`. It serves two purposes:
1.  **Main App**: A simple trading interface for sending manual/automated orders.
2.  **Certification Tool**: A test runner (`fixtest`) designed to execute 50+ certification scenarios required for exchange sign-off.

## Some Features
*   **Sign-off Test Runner**: Automated and interactive execution of 50+ test cases.
*   **Order Complexity**: Support for Market, Limit, Stop, Stop Limit, MOC, and LOC orders.
*   **Fractional & Notional**: Support for fractional share quantities (Tag 57=FRAC) and dollar-based notional orders (Tag 152).
*   **Advanced Symbology**: Automatic suffix handling (e.g., `BRK.B` splits to Tag 55=BRK, Tag 65=B).
*   **Compliance**: All orders are automatically marked as **Agency** (Tag 47=A).
*   **Interactive Mode**: Step-by-step test execution for manual verification.

## Project Structure

*   **`cmd/thisgofix/`**: Main application entry point.
*   **`cmd/fixtest/`**: Certification test runner entry point.
*   **`internal/fix/`**: Core logic and utilities.
    *   `application.go`: Implements `quickfix.Application` and high-level order methods.
    *   `constants.go`: Centralized FIX tags, message types, and enum values.
    *   `utils.go`: Helper functions for human-readable console output and translations.
*   **`run-fix-cert.sh`**: Helper script to run the certification suite with arguments.
*   **`config.cfg`**: The FIX session configuration.
*   **`go.mod` / `go.sum`**: Dependency management.

## Prerequisites

You need **Go** installed on your computer.

1.  Open your terminal/command prompt.
2.  Type `go version` to check if it's installed.

### Download Dependencies (Use the script!)

Since you are seeing "missing go.sum entry" errors, you **MUST** run the provided script to fix it.

1.  Double-click `fix_dependencies.bat` (if on Windows).
2.  Or run `go mod tidy` in your terminal.

This will download the `quickfixgo` library and fix the `go.sum` file.

## Configuration (`config.cfg`)

If you need to change the server details, edit `config.cfg`.

*   **SocketConnectHost**: The IP address of the server (e.g., `10.10.70.60`).
*   **SocketConnectPort**: The port number (e.g., `7605`).
*   **SenderCompID**: Your ID.
*   **TargetCompID**: The server's ID.


## How to Run

### Certification Test Runner
To run the full suite of 50 scenarios:

```bash
# Automatic mode (1s delay between cases)
./run-fix-cert.sh

# Interactive mode (Waits for ENTER between every case)
./run-fix-cert.sh -interactive

```

### Main Application Runner
To run the standard application:

```bash
go run cmd/thisgofix/main.go
```

### Watching the Output

You should see output like this:

```
System: Starting FIX Client...
Session Created: FIX.4.4:MY_UAT->SERVER_UAT
Connection: Logged On to FIX.4.4:MY_UAT->SERVER_UAT
Action: Sending New Order Single (ID: ORD-167890000)...
...
--- [Execution Report] ---
Server Order ID: 12345
Event Type:      New (Accepted)
Current Status:  New
Filled Qty:      0
Remaining Qty:   100
...
```


### How to Stop the App

Press **Ctrl+C** in your terminal to stop the application gracefully.

### Certification Application
To run the certification application:

```bash
go run cmd/fixtest/main.go
```

## Certification Scenarios Covered

The `fixtest` runner executes scenarios categorized as follows:
1.  **SIDES**: Buy, Sell, Sell Short (with automatic Locate tags).
2.  **ORDER TYPES**: Market, Limit, Stop, Stop Limit, Market on Close, Limit on Close.
3.  **EXEC INST**: Not Held vs Held instructions.
4.  **TIF**: Day, GTC, IOC, FOK.
5.  **SETTLEMENT**: Regular through T+5, Cash, Next Day, Future.
6.  **FRACTIONAL**: Fractional quantity orders using the `FRAC` indicator.
7.  **NOTIONAL**: Dollar-based orders using Tag 152.
8.  **EXTENDED HOURS**: AM Only, PM Only, and All Sessions.
9.  **MISC**: Partial fills, cancels, price/qty modifications, and rejections.


### How to Stop the App

Press **Ctrl+C** in your terminal to stop the application gracefully.

## How to Run Tests

To run the unit tests we added:

```bash
go test ./internal/fix
```


## Troubleshooting / Common Errors

### Error: `missing go.sum entry`
**Reason**: The library `quickfixgo` has not been downloaded to your computer yet.
**Fix**: Run `fix_dependencies.bat`.

### Error: `Connection refused`
**Reason**: The server might be down, or your firewall is blocking the connection.
**Fix**: Check if the Host and Port in `config.cfg` are correct.

### Error: `MsgSeqNum too low`
**Reason**: This means your `store/` folder is out of sync with the server (e.g., the server was reset but your local files weren't).
**Fix**: Delete the `store/` folder and try again. This resets your sequence numbers.

### Tip: Understanding Order Status
**Always look at the ExecutionReport.**
Check `ExecType` (Tag 150), not just the message type, to know the true status of your orders.
*   `ExecType=0 (New)`: Order accepted.
*   `ExecType=8 (Rejected)`: Order failed.
*   `ExecType=F (Trade)`: Trade executed.

## Detailed Code Explanation

The `internal/fix/application.go` file contains the core logic. Open it in your editor and read through the `Application` struct methods (`OnLogon`, `FromApp`, etc.) to see exactly how we handle events.
