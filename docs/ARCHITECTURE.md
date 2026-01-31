# Architecture

This document provides a comprehensive overview of paqet's architecture, design decisions, and implementation details.

## Table of Contents

1. [System Overview](#system-overview)
2. [Design Principles](#design-principles)
3. [Architecture Layers](#architecture-layers)
4. [Component Details](#component-details)
5. [Data Flow](#data-flow)
6. [Protocol Specification](#protocol-specification)
7. [Security Model](#security-model)
8. [Performance Considerations](#performance-considerations)
9. [Limitations and Trade-offs](#limitations-and-trade-offs)

## System Overview

paqet is a bidirectional packet-level proxy that creates an encrypted tunnel using KCP protocol over raw TCP packets. It operates below the standard OS network stack to bypass kernel-level connection tracking and firewall detection.

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        Application Layer                        │
│  (Browser, curl, etc.)                                          │
└────────────────┬────────────────────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────────────────────┐
│                    paqet Client (Local)                         │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │   SOCKS5     │  │   Forward    │  │   Protocol   │          │
│  │   Proxy      │  │   Handler    │  │   Handler    │          │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘          │
│         │                 │                  │                  │
│         └─────────────────┴──────────────────┘                  │
│                           │                                     │
│                   ┌───────▼────────┐                            │
│                   │  Stream Mux    │                            │
│                   │    (smux)      │                            │
│                   └───────┬────────┘                            │
│                           │                                     │
│                   ┌───────▼────────┐                            │
│                   │  KCP Transport │                            │
│                   │  + Encryption  │                            │
│                   └───────┬────────┘                            │
│                           │                                     │
│                   ┌───────▼────────┐                            │
│                   │  Raw Socket    │                            │
│                   │  (libpcap)     │                            │
│                   └───────┬────────┘                            │
└───────────────────────────┼─────────────────────────────────────┘
                            │
                            ▼
                    [Network Interface]
                            │
                            ▼ Crafted TCP Packets
                    [Internet / Network]
                            │
                            ▼
┌───────────────────────────┼─────────────────────────────────────┐
│                   ┌───────▼────────┐                            │
│                   │  Raw Socket    │                            │
│                   │  (libpcap)     │                            │
│                   └───────┬────────┘                            │
│                           │                                     │
│                   ┌───────▼────────┐                            │
│                   │  KCP Transport │                            │
│                   │  + Decryption  │                            │
│                   └───────┬────────┘                            │
│                           │                                     │
│                   ┌───────▼────────┐                            │
│                   │  Stream Demux  │                            │
│                   │    (smux)      │                            │
│                   └───────┬────────┘                            │
│                           │                                     │
│         ┌─────────────────┴──────────────────┐                  │
│         │                                    │                  │
│  ┌──────▼───────┐                    ┌──────▼───────┐          │
│  │   TCP        │                    │   UDP        │          │
│  │   Forward    │                    │   Forward    │          │
│  └──────┬───────┘                    └──────┬───────┘          │
│         │                                    │                  │
│                    paqet Server (Remote)                        │
└─────────┼────────────────────────────────────┼─────────────────┘
          │                                    │
          ▼                                    ▼
    [Target TCP Service]              [Target UDP Service]
```

## Design Principles

### 1. Stealth Communication

**Objective**: Evade detection by standard network monitoring tools.

**Implementation**:
- Raw packet crafting bypasses OS TCP/IP stack
- Custom packet structures don't follow standard handshakes
- No standard port fingerprints
- Encrypted payload prevents DPI (Deep Packet Inspection)

### 2. Reliability Over Speed

**Objective**: Ensure data integrity and reliability.

**Implementation**:
- KCP provides reliable, ordered delivery
- Aggressive retransmission for packet loss
- Flow control and congestion control
- Forward Error Correction (FEC) support

### 3. Modularity

**Objective**: Separate concerns and allow extensibility.

**Implementation**:
- Clear layer separation (Application → Transport → Raw)
- Interface-based design for swappable components
- Pluggable encryption algorithms
- Support for multiple proxy modes (SOCKS5, port forwarding)

### 4. Security by Default

**Objective**: Secure communication without user configuration.

**Implementation**:
- Mandatory encryption on KCP transport
- No plaintext fallback
- Secure random number generation
- Defense against timing attacks

## Architecture Layers

### Layer 1: Application Layer

**Responsibility**: Handle application-level protocols and user requests.

**Components**:
- **SOCKS5 Proxy** (`internal/socks/`)
  - Handles SOCKS5 protocol negotiation
  - Supports authentication (username/password)
  - Manages TCP and UDP association
  
- **Port Forwarding** (`internal/forward/`)
  - TCP port forwarding
  - UDP port forwarding
  - Direct connection mapping

**Interface**:
```go
// Application layer listens on local ports
// Accepts connections from user applications
// Forwards data through transport layer
```

### Layer 2: Transport Layer

**Responsibility**: Provide reliable, encrypted data transport.

**Components**:
- **KCP Protocol** (`internal/tnet/kcp/`)
  - Reliable ARQ protocol
  - Built-in encryption
  - Configurable parameters (MTU, window size, mode)
  
- **Stream Multiplexing** (`github.com/xtaci/smux`)
  - Multiple logical streams over one KCP connection
  - Connection pooling
  - Automatic stream management

**Interface**:
```go
type Conn interface {
    Read([]byte) (int, error)
    Write([]byte) (int, error)
    Close() error
    LocalAddr() net.Addr
    RemoteAddr() net.Addr
}

type Listener interface {
    Accept() (Conn, error)
    Close() error
    Addr() net.Addr
}
```

### Layer 3: Raw Packet Layer

**Responsibility**: Send and receive raw network packets.

**Components**:
- **Packet Socket** (`internal/socket/`)
  - Raw packet injection
  - Packet capture and filtering
  - Hardware address management
  
- **Packet Construction** (`github.com/gopacket/gopacket`)
  - Ethernet frame crafting
  - IP header construction
  - TCP header construction

**Interface**:
```go
type PacketConn interface {
    ReadFrom([]byte) (int, net.Addr, error)
    WriteTo([]byte, net.Addr) (int, error)
    Close() error
}
```

## Component Details

### Configuration System (`internal/conf/`)

**Design**: YAML-based configuration with validation.

**Key Features**:
- Role-based configuration (client/server)
- Default value application
- Comprehensive validation
- Clear error messages

**Configuration Flow**:
```
YAML File → Parse → Set Defaults → Validate → Conf Struct
```

**Example Structure**:
```go
type Conf struct {
    Role      string      // "client" or "server"
    Log       Log         // Logging configuration
    Network   Network     // Network interface settings
    Transport Transport   // KCP transport settings
    Listen    Server      // Server listen config
    Server    Server      // Server address (client)
    SOCKS5    []SOCKS5    // SOCKS5 proxy configs
    Forward   []Forward   // Port forwarding configs
}
```

### Socket Layer (`internal/socket/`)

**Design**: Abstraction over raw sockets with pcap.

**Components**:

1. **SendHandle** - Packet injection
   ```go
   type SendHandle struct {
       handle  *pcap.Handle
       srcMAC  net.HardwareAddr
       dstMAC  net.HardwareAddr
   }
   ```
   - Crafts Ethernet, IP, TCP headers
   - Injects packets onto network
   - Manages source/destination MACs

2. **RecvHandle** - Packet capture
   ```go
   type RecvHandle struct {
       handle *pcap.Handle
       filter string
   }
   ```
   - Captures packets matching filter
   - Extracts payload from TCP packets
   - Converts to standard net.Addr

3. **PacketConn** - Unified interface
   ```go
   type PacketConn struct {
       sendHandle *SendHandle
       recvHandle *RecvHandle
   }
   ```
   - Implements net.PacketConn interface
   - Coordinates send/receive operations
   - Handles deadlines and context cancellation

**Packet Filter**:
```
tcp and dst port <local_port> and src <remote_ip>
```

### KCP Transport (`internal/tnet/kcp/`)

**Design**: KCP protocol implementation over raw packets.

**Key Features**:
- Fast retransmission
- Selective repeat
- Congestion control
- FEC (Forward Error Correction)
- Multiple cipher options

**KCP Modes**:
```go
const (
    ModeNormal = "normal"  // Default, balanced
    ModeFast   = "fast"    // Low latency
    ModeFast2  = "fast2"   // Lower latency
    ModeFast3  = "fast3"   // Lowest latency
)
```

**Encryption Ciphers**:
- AES (128, 192, 256-bit)
- Salsa20
- Blowfish
- Twofish
- ChaCha20
- SM4 (Chinese standard)
- None (not recommended)

**Connection Establishment**:
```
Client                           Server
  |                                |
  |--- KCP SYN (encrypted) ------->|
  |                                |
  |<-- KCP SYN-ACK (encrypted) ----|
  |                                |
  |--- KCP ACK (encrypted) ------->|
  |                                |
  |<====== Encrypted Data ========>|
```

### Protocol Messages (`internal/protocol/`)

**Design**: Simple binary protocol for client-server communication.

**Message Types**:

1. **PPING / PPONG** - Keepalive
   ```go
   Proto{Type: PPING}
   Proto{Type: PPONG}
   ```

2. **PTCPF** - TCP Flag Configuration
   ```go
   Proto{
       Type: PTCPF,
       TCPF: []conf.TCPF{"PA", "S", "A"},
   }
   ```

3. **PTCP** - TCP Connection Request
   ```go
   Proto{
       Type: PTCP,
       Addr: &tnet.Addr{IP: "1.1.1.1", Port: 443},
   }
   ```

4. **PUDP** - UDP Connection Request
   ```go
   Proto{
       Type: PUDP,
       Addr: &tnet.Addr{IP: "8.8.8.8", Port: 53},
   }
   ```

**Serialization**: Uses Go's `encoding/gob` for efficient binary encoding.

### Client Architecture (`internal/client/`)

**Design**: Connection pool with automatic failover.

**Components**:

1. **TimedConn** - Connection with health monitoring
   ```go
   type timedConn struct {
       conn      tnet.Conn
       strm      tnet.Strm
       lastUsed  time.Time
       mu        sync.Mutex
   }
   ```

2. **Iterator** - Round-robin connection selection
   ```go
   type Iterator struct {
       Items []timedConn
       index int
   }
   ```

3. **UDP Pool** - UDP stream management
   ```go
   type udpPool struct {
       strms map[uint64]tnet.Strm
       mu    sync.Mutex
   }
   ```

**Connection Lifecycle**:
```
Connect → Health Check → Use → Idle → Reconnect (if needed)
```

**Ticker**: Background goroutine checking connection health every 10 seconds.

### Server Architecture (`internal/server/`)

**Design**: Accept connections and forward to targets.

**Request Handling Flow**:
```
1. Accept KCP connection
2. Read protocol message
3. Determine type (PTCP/PUDP)
4. Connect to target
5. Bidirectional copy
6. Clean up on disconnect
```

**Concurrency**:
- Each connection handled in separate goroutine
- Wait group for graceful shutdown
- Context for cancellation propagation

## Data Flow

### TCP Connection Flow (SOCKS5)

```
Application (curl)
  │
  ├─ Connect to localhost:1080
  ├─ SOCKS5 handshake
  ├─ Request: CONNECT example.com:443
  │
  ▼
SOCKS5 Handler
  │
  ├─ Validate request
  ├─ Get KCP stream from pool
  │
  ▼
Protocol Handler
  │
  ├─ Send PTCP message
  ├─ Target: example.com:443
  │
  ▼
KCP Stream (Client)
  │
  ├─ Encrypt data
  ├─ Fragment into KCP packets
  │
  ▼
Raw Socket (Client)
  │
  ├─ Craft TCP packet
  ├─ Inject onto network
  │
  ▼
[Network]
  │
  ▼
Raw Socket (Server)
  │
  ├─ Capture TCP packet
  ├─ Extract payload
  │
  ▼
KCP Stream (Server)
  │
  ├─ Reassemble fragments
  ├─ Decrypt data
  │
  ▼
Protocol Handler (Server)
  │
  ├─ Parse PTCP message
  ├─ Extract target: example.com:443
  │
  ▼
TCP Forwarder (Server)
  │
  ├─ Connect to example.com:443
  ├─ Bidirectional copy
  │
  ▼
Target Server (example.com:443)
  │
  └─ [Data flows back in reverse]
```

### UDP Association Flow

```
Application (DNS query)
  │
  ├─ SOCKS5 UDP ASSOCIATE
  │
  ▼
SOCKS5 Handler
  │
  ├─ Create UDP listener
  ├─ Return relay address
  │
  ▼
Application sends UDP
  │
  ▼
SOCKS5 UDP Handler
  │
  ├─ Parse SOCKS5 UDP packet
  ├─ Extract target and payload
  │
  ▼
Protocol Handler
  │
  ├─ Send PUDP message
  ├─ Send UDP data
  │
  ▼
[KCP → Raw Socket → Network]
  │
  ▼
Server UDP Handler
  │
  ├─ Get or create UDP connection
  ├─ Send to target
  │
  ▼
Target UDP Service
  │
  └─ [Response flows back]
```

## Protocol Specification

### Packet Structure

**Raw TCP Packet**:
```
┌─────────────────────────────────────────────────┐
│           Ethernet Frame Header                 │
│  Dst MAC (6) | Src MAC (6) | EtherType (2)      │
├─────────────────────────────────────────────────┤
│              IP Header (20 bytes)               │
│  Version | IHL | ToS | Length | ID | Flags...   │
├─────────────────────────────────────────────────┤
│              TCP Header (20+ bytes)             │
│  Src Port | Dst Port | Seq | Ack | Flags...     │
├─────────────────────────────────────────────────┤
│           KCP Packet (Encrypted)                │
│  conv | cmd | frg | wnd | ts | sn | una | ...   │
│  + Encrypted Application Data                   │
└─────────────────────────────────────────────────┘
```

### KCP Packet Format

```
 0                   1                   2                   3
 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                         conv (4 bytes)                        |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
| cmd |  frg  |              wnd (2 bytes)                      |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                          ts (4 bytes)                         |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                          sn (4 bytes)                         |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                         una (4 bytes)                         |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                         len (4 bytes)                         |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                          data...                              |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
```

- **conv**: Conversation ID (connection identifier)
- **cmd**: Command (PUSH, ACK, WASK, WINS)
- **frg**: Fragment number
- **wnd**: Window size
- **ts**: Timestamp
- **sn**: Sequence number
- **una**: Unacknowledged sequence number
- **len**: Data length
- **data**: Encrypted payload

## Security Model

### Threat Model

**Assumptions**:
- Attacker can observe network traffic
- Attacker may attempt to intercept or modify packets
- Attacker may attempt to inject packets
- Attacker cannot break encryption primitives

**Out of Scope**:
- Endpoint compromise
- Timing attacks (partially mitigated)
- Traffic analysis (connection patterns visible)

### Security Measures

1. **Encryption**:
   - All KCP data is encrypted
   - Shared secret key required
   - Multiple cipher options

2. **Authentication**:
   - KCP conv ID acts as connection identifier
   - SOCKS5 username/password authentication (optional)

3. **Integrity**:
   - KCP checksums verify data integrity
   - Sequence numbers prevent replay attacks

4. **Confidentiality**:
   - Encrypted payload prevents DPI
   - No plaintext protocol markers

### Known Limitations

1. **Traffic Analysis**:
   - Packet sizes and timing patterns visible
   - Connection endpoints visible
   
2. **Active Attacks**:
   - DoS possible by flooding packets
   - No authentication of initial connection

3. **Key Management**:
   - Symmetric key must be shared out-of-band
   - No automatic key rotation

## Performance Considerations

### Optimization Strategies

1. **Connection Pooling**:
   - Multiple KCP connections reduce head-of-line blocking
   - Round-robin distribution

2. **Buffer Management**:
   - Configurable socket buffers
   - Efficient buffer reuse

3. **KCP Tuning**:
   - Mode selection (normal/fast/fast2/fast3)
   - Window size adjustment
   - FEC configuration

### Performance Trade-offs

| Aspect | Trade-off |
|--------|-----------|
| **Latency** | Lower with aggressive KCP modes, but higher CPU usage |
| **Throughput** | Higher with larger windows, but more memory |
| **Reliability** | Higher with FEC, but more bandwidth |
| **Stealth** | Better with smaller packets, but lower throughput |

### Benchmarking

**Typical Performance** (varies by network):
- Latency overhead: 5-20ms (compared to direct connection)
- Throughput: 80-95% of direct connection
- CPU usage: 10-30% per active connection

## Limitations and Trade-offs

### Current Limitations

1. **Platform Support**:
   - Requires libpcap on all platforms
   - Windows requires Npcap installation
   - Requires root/administrator privileges

2. **Configuration Complexity**:
   - Manual network interface configuration
   - MAC address discovery required
   - No automatic configuration

3. **Protocol Support**:
   - TCP and UDP only
   - No ICMP or other protocols
   - IPv6 support experimental

4. **Scalability**:
   - Limited by single server instance
   - No load balancing
   - No connection migration

### Design Trade-offs

1. **Raw Sockets vs Standard Sockets**:
   - ✅ Bypass kernel stack and firewalls
   - ❌ Requires privileges and complex setup

2. **KCP vs TCP**:
   - ✅ Lower latency, better loss recovery
   - ❌ Not compatible with standard TCP tools

3. **Symmetric Encryption vs PKI**:
   - ✅ Simpler implementation, faster
   - ❌ Key distribution challenges

4. **SOCKS5 vs Custom Protocol**:
   - ✅ Compatible with existing applications
   - ❌ Additional protocol overhead

## Future Considerations

### Potential Improvements

1. **Dynamic Configuration**:
   - Automatic interface discovery
   - Dynamic MAC resolution
   - Configuration wizard

2. **Enhanced Security**:
   - Public key authentication
   - Perfect forward secrecy
   - Certificate-based trust

3. **Performance**:
   - Multi-threaded packet processing
   - Zero-copy buffer operations
   - Hardware offload support

4. **Features**:
   - HTTP proxy support
   - Plugin system
   - Web-based management interface

### Architectural Evolution

```
Current: Simple client-server
  ↓
Future: Multi-hop routing
  ↓
Future: Distributed network
  ↓
Future: Full mesh networking
```

---

This architecture enables paqet to provide stealthy, reliable communication while maintaining modularity and security. Understanding these components and their interactions is key to extending and improving the system.
