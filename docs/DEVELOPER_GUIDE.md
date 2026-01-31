# Developer Guide

Welcome to the paqet developer guide! This document will help you understand the project structure, architecture, and how to contribute effectively.

## Table of Contents

1. [Project Overview](#project-overview)
2. [Technology Stack](#technology-stack)
3. [Project Structure](#project-structure)
4. [Core Concepts](#core-concepts)
5. [Development Setup](#development-setup)
6. [Code Organization](#code-organization)
7. [Key Components](#key-components)
8. [Development Workflow](#development-workflow)
9. [Debugging Tips](#debugging-tips)
10. [Next Steps](#next-steps)

## Project Overview

**paqet** is a bidirectional packet-level proxy that operates using raw sockets to completely bypass the host operating system's TCP/IP stack. It uses KCP (a fast and reliable ARQ protocol) for secure, reliable transport over crafted raw TCP packets.

### What Makes paqet Unique?

- **Packet-Level Operation**: Works below the OS network stack using raw sockets
- **Custom Protocol**: Uses KCP over crafted TCP packets for encryption and reliability
- **Stealth Communication**: Bypasses standard firewall detection by not using standard handshakes
- **Multiple Modes**: Supports SOCKS5 proxy and port forwarding

### Use Cases

- Network security research and penetration testing
- Bypassing restrictive firewalls
- Data exfiltration scenarios
- Educational purposes for low-level network programming

## Technology Stack

### Core Technologies

- **Language**: Go 1.25+
- **Packet Manipulation**: `gopacket` library for packet crafting/parsing
- **Packet Capture**: `libpcap` for capturing network packets
- **Transport Protocol**: KCP (reliable UDP-like protocol) over raw TCP
- **Multiplexing**: `smux` for stream multiplexing
- **Encryption**: Built-in KCP encryption with multiple cipher options

### Key Dependencies

```go
github.com/gopacket/gopacket   // Packet manipulation
github.com/xtaci/kcp-go/v5     // KCP protocol implementation
github.com/xtaci/smux          // Stream multiplexing
github.com/txthinking/socks5   // SOCKS5 server implementation
github.com/spf13/cobra         // CLI framework
golang.org/x/crypto            // Cryptographic primitives
```

## Project Structure

```
paqet/
├── cmd/                    # Command-line interface
│   ├── main.go            # Entry point, CLI setup
│   ├── run/               # Run command (client/server)
│   ├── dump/              # Packet dump utility
│   ├── ping/              # Ping utility
│   ├── secret/            # Secret key generation
│   ├── iface/             # Network interface info
│   └── version/           # Version command
├── internal/              # Internal packages (not importable)
│   ├── client/           # Client implementation
│   ├── server/           # Server implementation
│   ├── conf/             # Configuration management
│   ├── socket/           # Raw socket handling
│   ├── protocol/         # Protocol definitions
│   ├── tnet/             # Transport network abstraction
│   │   └── kcp/         # KCP transport implementation
│   ├── socks/            # SOCKS5 proxy implementation
│   ├── forward/          # Port forwarding implementation
│   ├── flog/             # Logging utilities
│   └── pkg/              # Shared internal packages
│       ├── buffer/       # Buffer management
│       ├── errors/       # Error types
│       ├── hash/         # Hashing utilities
│       └── iterator/     # Iterator pattern
├── example/              # Configuration examples
├── docs/                 # Documentation
└── tmp/                  # Temporary/development files
```

### Directory Responsibilities

- **`cmd/`**: All user-facing commands and CLI logic
- **`internal/`**: Core implementation, not meant to be imported by external projects
- **`internal/client/`**: Client-side connection management, SOCKS5/forwarding
- **`internal/server/`**: Server-side listener and connection handling
- **`internal/socket/`**: Low-level raw packet send/receive
- **`internal/tnet/`**: Transport network abstraction layer
- **`internal/conf/`**: Configuration parsing and validation

## Core Concepts

### 1. Raw Socket Communication

Unlike normal network applications that use the OS TCP/IP stack, paqet:
- Opens raw network interfaces using libpcap
- Crafts TCP packets manually with custom headers
- Injects packets directly onto the network
- Captures incoming packets before the OS sees them

**Why?** This bypasses kernel-level connection tracking and firewall hooks.

### 2. KCP Transport Layer

KCP (KCP ARQ Protocol) provides:
- Reliable, ordered delivery (like TCP)
- Lower latency than TCP (aggressive retransmission)
- Built-in encryption
- Flow control and congestion control

paqet encapsulates KCP inside raw TCP packets, creating a "KCP over raw TCP" transport.

### 3. Three-Layer Architecture

```
Application Layer (SOCKS5/Forward)
         ↓
Transport Layer (KCP + Encryption)
         ↓
Raw Packet Layer (Crafted TCP Packets)
```

### 4. Client-Server Model

**Client Side:**
- Accepts connections from applications (SOCKS5 or port forwarding)
- Multiplexes connections over KCP transport to server
- Handles local packet injection and capture

**Server Side:**
- Listens for KCP connections over raw packets
- Demultiplexes streams and connects to target services
- Handles remote packet injection and capture

## Development Setup

### Prerequisites

1. **Install Go 1.25+**
   ```bash
   go version  # Should be 1.25 or higher
   ```

2. **Install libpcap**
   
   - **Linux (Debian/Ubuntu):**
     ```bash
     sudo apt-get install libpcap-dev
     ```
   
   - **Linux (RHEL/CentOS):**
     ```bash
     sudo yum install libpcap-devel
     ```
   
   - **macOS:**
     ```bash
     xcode-select --install
     ```
   
   - **Windows:**
     - Install [Npcap](https://npcap.com/)
     - Ensure "WinPcap API-compatible Mode" is enabled

3. **Clone the Repository**
   ```bash
   git clone https://github.com/hanselime/paqet.git
   cd paqet
   ```

4. **Install Dependencies**
   ```bash
   go mod download
   ```

5. **Build the Project**
   ```bash
   go build -o paqet cmd/main.go
   ```

### Running Locally

**Test Network Interface Discovery:**
```bash
./paqet iface
```

**Generate Secret Key:**
```bash
./paqet secret
```

**Run Client:**
```bash
sudo ./paqet run -c example/client.yaml.example
```

**Run Server:**
```bash
sudo ./paqet run -c example/server.yaml.example
```

> **Note**: Raw socket operations require root/administrator privileges.

## Code Organization

### Configuration Flow

1. User provides YAML config file
2. `conf.LoadFromFile()` reads and parses YAML
3. `setDefaults()` applies default values
4. `validate()` checks for errors
5. Returns `*conf.Conf` struct

### Client Initialization Flow

```
main.go
  └─> run.Cmd (client mode)
      └─> client.New(cfg)
          ├─> socket.New() - creates raw packet connection
          ├─> newTimedConn() - establishes KCP connections
          └─> starts SOCKS5/forward listeners
```

### Server Initialization Flow

```
main.go
  └─> run.Cmd (server mode)
      └─> server.New(cfg)
          ├─> socket.New() - creates raw packet connection
          ├─> kcp.Listen() - starts KCP listener
          └─> handleConn() - handles incoming connections
```

### Packet Flow (Client → Server)

```
Application (e.g., curl)
  ↓
SOCKS5 Handler (internal/socks/)
  ↓
Protocol Message (internal/protocol/)
  ↓
KCP Stream (internal/tnet/kcp/)
  ↓
Raw Packet Send (internal/socket/)
  ↓
Network Interface
```

## Key Components

### 1. Socket Layer (`internal/socket/`)

**Purpose**: Handle raw packet injection and capture

**Key Files:**
- `socket.go` - Main PacketConn implementation
- `send_handle.go` - Packet injection (writing)
- `recv_handle.go` - Packet capture (reading)

**Key Operations:**
```go
// Create raw socket connection
conn, err := socket.New(ctx, &cfg.Network)

// Read incoming packet
n, addr, err := conn.ReadFrom(buffer)

// Send outgoing packet
n, err := conn.WriteTo(data, addr)
```

### 2. Protocol Layer (`internal/protocol/`)

**Purpose**: Define message types for client-server communication

**Message Types:**
- `PPING` / `PPONG` - Ping/pong for keepalive
- `PTCPF` - TCP flag configuration
- `PTCP` - TCP data forwarding
- `PUDP` - UDP data forwarding

**Usage:**
```go
proto := &protocol.Proto{
    Type: protocol.PTCP,
    Addr: targetAddr,
}
proto.Write(stream)
```

### 3. KCP Transport (`internal/tnet/kcp/`)

**Purpose**: Provide KCP-based reliable transport over raw packets

**Key Operations:**
```go
// Server: Listen for connections
listener, err := kcp.Listen(cfg, packetConn)

// Server: Accept connections
conn, err := listener.Accept()

// Client: Dial server
conn, err := kcp.Dial(cfg, serverAddr, packetConn)
```

### 4. Client (`internal/client/`)

**Purpose**: Manage client-side connections and proxying

**Key Components:**
- `client.go` - Main client struct
- `tcp.go` - TCP connection handling
- `udp.go` - UDP connection handling
- `dial.go` - Connection establishment
- `ticker.go` - Connection health monitoring

### 5. Server (`internal/server/`)

**Purpose**: Accept connections and forward to targets

**Key Components:**
- `server.go` - Main server struct
- `handle.go` - Connection handling logic
- `tcp.go` - TCP forwarding
- `udp.go` - UDP forwarding
- `ping.go` - Ping/pong handling

### 6. SOCKS5 (`internal/socks/`)

**Purpose**: Implement SOCKS5 proxy protocol

**Key Operations:**
```go
// Start SOCKS5 server
server := socks.NewServer(cfg)
server.ListenAndServe()
```

### 7. Configuration (`internal/conf/`)

**Purpose**: Parse, validate, and manage configuration

**Key Structures:**
```go
type Conf struct {
    Role      string      // "client" or "server"
    Log       Log         // Logging config
    Network   Network     // Interface, IPs, MACs
    Transport Transport   // KCP settings
    SOCKS5    []SOCKS5    // SOCKS5 configs
    Forward   []Forward   // Port forwarding configs
}
```

## Development Workflow

### 1. Making Changes

1. Create a feature branch:
   ```bash
   git checkout -b feature/my-feature
   ```

2. Make your changes

3. Test locally:
   ```bash
   go build -o paqet cmd/main.go
   sudo ./paqet run -c test-config.yaml
   ```

4. Run tests (if available):
   ```bash
   go test ./...
   ```

5. Format code:
   ```bash
   go fmt ./...
   ```

6. Commit and push:
   ```bash
   git commit -m "Add feature X"
   git push origin feature/my-feature
   ```

### 2. Adding a New Command

Example: Adding a `status` command

1. Create `cmd/status/status.go`:
   ```go
   package status

   import (
       "github.com/spf13/cobra"
   )

   var Cmd = &cobra.Command{
       Use:   "status",
       Short: "Show connection status",
       Run: func(cmd *cobra.Command, args []string) {
           // Implementation
       },
   }
   ```

2. Register in `cmd/main.go`:
   ```go
   import "paqet/cmd/status"
   
   func main() {
       rootCmd.AddCommand(status.Cmd)
       // ...
   }
   ```

### 3. Adding Configuration Options

1. Add field to config struct in `internal/conf/`:
   ```go
   type Transport struct {
       Protocol string
       NewOption int `yaml:"new_option"`
   }
   ```

2. Add validation in `validate()` method

3. Add default in `setDefaults()` method

4. Update example YAML files

### 4. Adding a New Protocol Type

1. Define constant in `internal/protocol/protocol.go`:
   ```go
   const PNEWTYPE PType = 0x06
   ```

2. Handle in client/server:
   ```go
   case protocol.PNEWTYPE:
       // Handle new type
   ```

## Debugging Tips

### 1. Enable Debug Logging

Set `log.level: "debug"` in your config file to see detailed logs.

### 2. Packet Capture

Use the built-in dump command:
```bash
sudo ./paqet dump -i eth0
```

Or use tcpdump:
```bash
sudo tcpdump -i eth0 port 9999 -X
```

### 3. Check Interface Configuration

```bash
./paqet iface
```

Lists all network interfaces with IPs and MACs.

### 4. Test KCP Connectivity

Use the ping command:
```bash
sudo ./paqet ping -c client.yaml
```

### 5. Common Issues

**"Permission denied" when running:**
- Solution: Run with `sudo` or as administrator

**"No such device" error:**
- Solution: Check interface name with `./paqet iface`

**"Invalid MAC address" error:**
- Solution: Use `arp -a` to find correct gateway MAC

**KCP connection timeout:**
- Check firewall rules
- Verify server is running
- Check network configuration (IPs, MACs, ports)

## Next Steps

### For New Contributors

1. Read [ARCHITECTURE.md](ARCHITECTURE.md) for detailed system design
2. Review [CONTRIBUTING.md](CONTRIBUTING.md) for contribution guidelines
3. Check [open issues](https://github.com/hanselime/paqet/issues) labeled `good first issue`
4. Join discussions and ask questions

### Learning Resources

- [Go by Example](https://gobyexample.com/)
- [gopacket Documentation](https://pkg.go.dev/github.com/gopacket/gopacket)
- [KCP Protocol](https://github.com/skywind3000/kcp)
- [SOCKS5 RFC](https://datatracker.ietf.org/doc/html/rfc1928)

### Suggested First Contributions

- Improve documentation
- Add unit tests
- Fix bugs in open issues
- Add platform-specific installation guides
- Improve error messages
- Add configuration validation

## Questions?

If you have questions or need help:
- Check the [Troubleshooting Guide](TROUBLESHOOTING.md)
- Open a GitHub issue
- Review existing discussions

Happy coding! 🚀
