# Go FIX Client (Beginner Guide)

This is a simple trading application that connects to a server using the **FIX Protocol** (Financial Information eXchange). It is designed to be easy to read and understand for new developers.

## What does this app do?

1.  **Connects** to a trading server (like a stock exchange).
2.  **Logs in** automatically.
3.  **Sends a "Buy" order** for EUR/USD.
4.  **Prints the result** (Did we buy it? Was it rejected?) to the screen.

## Project Structure

We follow the standard Go project layout:

*   **`cmd/thisgofix/`**: Contains the main entry point (`main.go`).
*   **`internal/fix/`**: Contains the core application logic.
    *   `application.go`: Handles FIX events (`OnLogon`, `FromApp`).
    *   `constants.go`: Defines FIX tags and message types.
    *   `application_test.go`: Unit tests for the application logic.
    *   `utils_test.go`: Unit tests for helper functions.
*   **`config.cfg`**: The FIX session configuration.
*   **`go.mod` / `go.sum`**: Dependency management.

## Prerequisites

You need **Go** installed on your computer.

1.  Open your terminal/command prompt.
2.  Type `go version` to check if it's installed.

## How to Run

### Step 1: Download Dependencies (Use the script!)

Since you are seeing "missing go.sum entry" errors, you **MUST** run the provided script to fix it.

1.  Double-click `fix_dependencies.bat` (if on Windows).
2.  Or run `go mod tidy` in your terminal.

This will download the `quickfixgo` library and fix the `go.sum` file.

### Step 2: Start the App

Run the application with:

```bash
go run cmd/thisgofix/main.go
```

### Step 3: Watch the Output

You should see output like this:

```
System: Starting FIX Client...
Session Created: FIX.4.4:BIYACAP_UAT->RAPTOR_UAT
Connection: Logged On to FIX.4.4:BIYACAP_UAT->RAPTOR_UAT
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

### Step 4: Stop the App

Press **Ctrl+C** in your terminal to stop the application gracefully.

## How to Run Tests

To run the unit tests we added:

```bash
go test ./internal/fix
```

## Configuration (`config.cfg`)

If you need to change the server details, edit `config.cfg`.

*   **SocketConnectHost**: The IP address of the server (e.g., `10.10.70.60`).
*   **SocketConnectPort**: The port number (e.g., `7605`).
*   **SenderCompID**: Your ID.
*   **TargetCompID**: The server's ID.

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

## detailed Code Explanation

The `internal/fix/application.go` file contains the core logic. Open it in your editor and read through the `Application` struct methods (`OnLogon`, `FromApp`, etc.) to see exactly how we handle events.
