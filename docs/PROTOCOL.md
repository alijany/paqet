# Protocol Specification

Technical specification of the paqet protocol and packet formats.

## Overview

paqet implements a custom network protocol stack consisting of three layers:
1. Raw Packet Layer (Ethernet/IP/TCP)
2. Transport Layer (KCP)
3. Application Layer (Protocol Messages)

## Raw Packet Layer

### Ethernet Frame

Standard Ethernet II frame format:

```
 0                   1                   2                   3
 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                    Destination MAC Address                    |
|                         (6 bytes)                             |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                      Source MAC Address                       |
|                         (6 bytes)                             |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|         EtherType            |                                |
|          (0x0800 for IPv4)   |                                |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+                                +
|                          Payload                              |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
```

**Fields**:
- **Destination MAC**: Router/Gateway MAC address
- **Source MAC**: Local interface MAC address  
- **EtherType**: 0x0800 (IPv4) or 0x86DD (IPv6)

### IPv4 Header

Standard IPv4 header (20 bytes minimum):

```
 0                   1                   2                   3
 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|Version|  IHL  |Type of Service|          Total Length         |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|         Identification        |Flags|      Fragment Offset    |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|  Time to Live |    Protocol   |         Header Checksum       |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                       Source Address                          |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                    Destination Address                        |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
```

**Key Fields**:
- **Version**: 4 (IPv4)
- **IHL**: 5 (20 bytes, no options)
- **Protocol**: 6 (TCP)
- **Source Address**: Local IP address
- **Destination Address**: Server IP address

### TCP Header

Standard TCP header (20 bytes minimum):

```
 0                   1                   2                   3
 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|          Source Port          |       Destination Port        |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                        Sequence Number                        |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                    Acknowledgment Number                      |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|  Data |       |C|E|U|A|P|R|S|F|                               |
| Offset| Rsrvd |W|C|R|C|S|S|Y|I|            Window             |
|       |       |R|E|G|K|H|T|N|N|                               |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|           Checksum            |         Urgent Pointer        |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                          Payload                              |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
```

**Key Fields**:
- **Source Port**: Local port (random or configured)
- **Destination Port**: Server port (configured)
- **Sequence/Ack**: Managed by raw socket layer
- **Flags**: Configurable (default: PSH+ACK)
- **Payload**: KCP packets (encrypted)

**Default TCP Flags**: `PSH + ACK` (0x18)
- PSH: Indicates data should be delivered immediately
- ACK: Acknowledges receipt of data

## Transport Layer (KCP)

### KCP Packet Format

KCP is a reliable UDP-like protocol. Each KCP packet has:

```
 0                   1                   2                   3
 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                       Conversation ID                         |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|     CMD       |     FRG       |           Window              |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                         Timestamp                             |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                       Sequence Number                         |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                   Unacknowledged Number                       |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                            Length                             |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                        Data (encrypted)                       |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
```

**Fields**:
- **Conversation ID** (4 bytes): Connection identifier
- **CMD** (1 byte): Command type
  - 81 (PUSH): Data packet
  - 82 (ACK): Acknowledgment
  - 83 (WASK): Window size query
  - 84 (WINS): Window size response
- **FRG** (1 byte): Fragment number (for large packets)
- **Window** (2 bytes): Receive window size
- **Timestamp** (4 bytes): Sending timestamp
- **Sequence Number** (4 bytes): Packet sequence
- **Una** (4 bytes): Unacknowledged sequence number
- **Length** (4 bytes): Data length
- **Data**: Encrypted application data

### KCP Connection Establishment

Unlike TCP, KCP doesn't have a formal handshake. Connection is established implicitly:

```
Client                                    Server
  |                                          |
  |--- PUSH (conv=X, sn=0) ------------------>|
  |                                          |
  |<-- ACK (conv=X, una=1) -------------------|
  |                                          |
  |<======== Data Exchange ==================>|
```

### KCP Modes

Different modes configure KCP parameters for different scenarios:

**Normal Mode**:
```
nodelay=0, interval=40ms, resend=2, nc=1
```

**Fast Mode**:
```
nodelay=0, interval=30ms, resend=2, nc=1
```

**Fast2 Mode**:
```
nodelay=1, interval=20ms, resend=2, nc=1
```

**Fast3 Mode**:
```
nodelay=1, interval=10ms, resend=2, nc=0
```

Parameters:
- **nodelay**: 1 = no delay ACK, 0 = delay ACK
- **interval**: Internal update interval
- **resend**: Fast retransmission trigger (2 = after 2 ACKs)
- **nc**: 1 = enable flow control, 0 = disable

### Encryption

KCP data is encrypted before transmission. Supported ciphers:

| Cipher | Block Size | Key Size | Mode |
|--------|------------|----------|------|
| AES | 16 bytes | 32 bytes | CFB |
| AES-128 | 16 bytes | 16 bytes | CFB |
| AES-192 | 16 bytes | 24 bytes | CFB |
| Salsa20 | Stream | 32 bytes | Stream |
| ChaCha20 | Stream | 32 bytes | Stream |
| Blowfish | 8 bytes | Variable | CFB |
| Twofish | 16 bytes | 32 bytes | CFB |
| SM4 | 16 bytes | 16 bytes | CFB |

**Encryption Process**:
1. Derive key from shared secret
2. Generate IV (initialization vector)
3. Encrypt KCP payload
4. Prepend IV to ciphertext

### Forward Error Correction (FEC)

Optional FEC adds redundancy for packet loss recovery:

**FEC Block Structure**:
```
Data Shards (10):    [D0] [D1] [D2] ... [D9]
Parity Shards (3):   [P0] [P1] [P2]
```

**Properties**:
- Can recover from losing up to `parityshard` packets
- Adds `parityshard/datashard` bandwidth overhead
- Example: 10/3 adds 30% overhead

**FEC Packet Header**:
```
 0                   1                   2                   3
 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|F|    Shard ID        |        Sequence Number                |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                          KCP Packet                           |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
```

- **F**: Flag (0=data shard, 1=parity shard)
- **Shard ID**: Index within FEC group
- **Sequence Number**: FEC group sequence

## Application Layer

### Protocol Messages

paqet defines a simple binary protocol using Go's `encoding/gob`:

**Message Structure**:
```go
type Proto struct {
    Type byte           // Message type
    Addr *tnet.Addr    // Target address (optional)
    TCPF []conf.TCPF   // TCP flags (optional)
}
```

### Message Types

**PPING (0x01)** - Ping Request:
```
┌─────┬─────────────┐
│Type │ Addr  TCPF  │
│0x01 │ nil   nil   │
└─────┴─────────────┘
```

**PPONG (0x02)** - Ping Response:
```
┌─────┬─────────────┐
│Type │ Addr  TCPF  │
│0x02 │ nil   nil   │
└─────┴─────────────┘
```

**PTCPF (0x03)** - TCP Flag Configuration:
```
┌─────┬──────────────────────┐
│Type │ Addr  TCPF           │
│0x03 │ nil   ["PA"]         │
└─────┴──────────────────────┘
```

**PTCP (0x04)** - TCP Connection Request:
```
┌─────┬────────────────────────────┐
│Type │ Addr                 TCPF  │
│0x04 │ "example.com:443"    nil   │
└─────┴────────────────────────────┘
```

**PUDP (0x05)** - UDP Association Request:
```
┌─────┬────────────────────────────┐
│Type │ Addr                 TCPF  │
│0x05 │ "8.8.8.8:53"         nil   │
└─────┴────────────────────────────┘
```

### Protocol Flow - TCP Connection

```
Client                                              Server
  |                                                    |
  | 1. Application connects to SOCKS5                  |
  |    Request: CONNECT example.com:443                |
  |                                                    |
  | 2. Send PTCP message                               |
  |--- PTCP {Type:0x04, Addr:"example.com:443"} ------>|
  |                                                    |
  |                                  3. Server connects|
  |                                     to example.com |
  |                                                    |
  | 4. Bidirectional data flow                         |
  |<=============== Data ==============================>|
  |                                                    |
  | 5. Connection close                                |
  |--- FIN ------------------------------------------->|
  |<-- FIN-ACK ----------------------------------------|
```

### Protocol Flow - UDP Association

```
Client                                              Server
  |                                                    |
  | 1. SOCKS5 UDP ASSOCIATE                            |
  |                                                    |
  | 2. Send PUDP message                               |
  |--- PUDP {Type:0x05, Addr:"8.8.8.8:53"} ----------->|
  |                                                    |
  |                             3. Server creates UDP  |
  |                                association         |
  |                                                    |
  | 4. Send UDP data                                   |
  |--- UDP data ---------------------------------------->|
  |                                                    |
  |                                    5. Forward to   |
  |                                       target       |
  |                                                    |
  | 6. Receive UDP response                            |
  |<-- UDP data ----------------------------------------|
```

## Stream Multiplexing

paqet uses smux for stream multiplexing over KCP connections.

### Smux Frame Format

```
 0                   1                   2                   3
 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|    Version    |      CMD      |            Length             |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                           Stream ID                           |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                             Data                              |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
```

**Commands**:
- **SYN** (0): Stream open
- **FIN** (1): Stream close
- **PSH** (2): Data push
- **NOP** (3): Keepalive

## Connection Lifecycle

### Client Connection Lifecycle

```
1. Initialize
   ├─ Load configuration
   ├─ Create raw packet connection
   └─ Establish KCP connections (x N)

2. Listen for Applications
   ├─ Start SOCKS5 listener(s)
   └─ Start port forwarding listener(s)

3. Handle Requests
   ├─ Accept application connection
   ├─ Open smux stream
   ├─ Send protocol message (PTCP/PUDP)
   └─ Bidirectional data copy

4. Health Monitoring
   ├─ Ticker checks connections every 10s
   ├─ Send PPING, expect PPONG
   └─ Reconnect if unhealthy

5. Shutdown
   ├─ Close application listeners
   ├─ Close all streams
   ├─ Close KCP connections
   └─ Close raw packet connection
```

### Server Connection Lifecycle

```
1. Initialize
   ├─ Load configuration
   ├─ Create raw packet connection
   └─ Start KCP listener

2. Accept Connections
   ├─ Wait for KCP connection
   ├─ Accept smux session
   └─ Spawn handler goroutine

3. Handle Streams
   ├─ Accept smux stream
   ├─ Read protocol message
   ├─ Connect to target (PTCP/PUDP)
   └─ Bidirectional data copy

4. Connection Management
   ├─ Handle PPING → respond PPONG
   ├─ Handle stream close
   └─ Clean up resources

5. Shutdown
   ├─ Stop accepting new connections
   ├─ Wait for active connections
   └─ Close raw packet connection
```

## Security Considerations

### Threat Model

**Protected Against**:
- Passive eavesdropping (encryption)
- Packet modification (integrity checks)
- Replay attacks (sequence numbers)

**Not Protected Against**:
- Traffic analysis (packet sizes/timing visible)
- Active attacks (packet injection possible)
- Endpoint compromise
- Timing attacks

### Security Features

1. **Encryption**: All KCP data encrypted
2. **Authentication**: Conversation ID prevents hijacking
3. **Integrity**: KCP checksums detect corruption
4. **Confidentiality**: Encrypted payload prevents DPI

### Known Weaknesses

1. **No Perfect Forward Secrecy**: Same key used for all sessions
2. **No Public Key Auth**: Symmetric key only
3. **Traffic Patterns**: Connection metadata visible
4. **DoS Vulnerability**: No rate limiting on packets

## Performance Characteristics

### Latency

**Total Latency** = Base + KCP + Raw Socket

- **Base Network**: Varies by connection
- **KCP Overhead**: 5-20ms (mode dependent)
- **Raw Socket**: 1-5ms

**Typical**: 10-30ms additional latency

### Throughput

**Theoretical Max** = Window × MTU / RTT

Example:
- Window: 1024
- MTU: 1350 bytes
- RTT: 50ms
- Max: (1024 × 1350) / 0.05 = 27.6 MB/s

**Actual**: 80-95% of theoretical (overhead, loss, etc.)

### Overhead

| Component | Overhead |
|-----------|----------|
| Ethernet | 14 bytes |
| IP | 20 bytes |
| TCP | 20 bytes |
| KCP Header | 24 bytes |
| Encryption IV | 16 bytes (typical) |
| FEC (10/3) | 30% bandwidth |
| **Total** | ~100 bytes + 30% FEC |

---

This protocol specification is subject to change. Refer to the source code for the authoritative implementation.
