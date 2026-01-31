# Configuration Guide

Complete guide to configuring paqet for client and server deployments.

## Table of Contents

1. [Configuration Overview](#configuration-overview)
2. [Network Configuration](#network-configuration)
3. [Client Configuration](#client-configuration)
4. [Server Configuration](#server-configuration)
5. [Transport Settings](#transport-settings)
6. [Advanced Options](#advanced-options)
7. [Examples](#examples)

## Configuration Overview

paqet uses YAML configuration files with role-based settings. Each instance must be configured as either a `client` or `server`.

### Basic Structure

```yaml
role: "client"  # or "server"
log:
  level: "info"
network:
  # Interface and packet settings
transport:
  # KCP protocol settings
# Role-specific settings below
```

### Configuration File Location

Specify config file with `-c` or `--config` flag:
```bash
./paqet run -c /path/to/config.yaml
```

## Network Configuration

Network configuration is the most critical part and requires careful setup.

### Finding Network Information

**Linux:**
```bash
# Find interface and IP
ip a

# Find gateway IP
ip r | grep default

# Find gateway MAC
arp -n <gateway_ip>
```

**macOS:**
```bash
# Find interface and IP
ifconfig

# Find gateway IP
netstat -rn | grep default

# Find gateway MAC
arp -n <gateway_ip>
```

**Windows:**
```powershell
# Find interface and IP
ipconfig /all

# List interfaces
netsh interface show interface

# Find NPF GUID
Get-NetAdapter | Select-Object Name, InterfaceGuid

# Find gateway MAC
arp -a <gateway_ip>
```

### Network Configuration Block

```yaml
network:
  interface: "eth0"           # Interface name
  guid: "{...}"               # Windows only: NPF GUID
  
  ipv4:
    addr: "192.168.1.100:0"   # Local IP:Port (0 = random)
    router_mac: "aa:bb:cc:dd:ee:ff"  # Gateway MAC
  
  ipv6:  # Optional
    addr: "[::1]:0"
    router_mac: "aa:bb:cc:dd:ee:ff"
  
  tcp:
    local_flag: ["PA"]        # TCP flags for local packets
    remote_flag: ["PA"]       # TCP flags for remote packets
  
  pcap:
    sockbuf: 4194304          # Socket buffer size (bytes)
```

### Interface Selection

Choose the correct network interface:

- **Wired Ethernet**: Usually `eth0`, `ens3`, `Ethernet`
- **WiFi**: Usually `wlan0`, `wlp3s0`, `Wi-Fi`
- **VPN**: Use the VPN interface if routing through VPN

**Verify with:**
```bash
./paqet iface
```

### IP Address Configuration

**Format**: `<IP>:<PORT>`

- **Client**: Use your local IP address
- **Server**: Use your public or local IP address
- **Port**: 
  - Use `0` for random port assignment
  - Or specify a fixed port (must match in all places)

**Examples**:
```yaml
# Random port (recommended for client)
addr: "192.168.1.100:0"

# Fixed port
addr: "10.0.0.50:9999"

# IPv6
addr: "[2001:db8::1]:9999"
```

### MAC Address Configuration

The `router_mac` is the hardware address of your default gateway (router).

**Common Issues**:
- Using local interface MAC instead of gateway MAC ❌
- Using uppercase/lowercase inconsistently (use lowercase)
- Missing colons or using hyphens

**Correct Format**:
```yaml
router_mac: "aa:bb:cc:dd:ee:ff"  # Lowercase, colon-separated
```

### TCP Flags

Configure TCP flags for packet crafting:

```yaml
tcp:
  local_flag: ["PA"]    # Push + Ack
  remote_flag: ["PA"]   # Push + Ack
```

**Available Flags**:
- `S` - SYN
- `A` - ACK
- `P` - PUSH
- `F` - FIN
- `R` - RST
- `PA` - PUSH + ACK (default)

**When to Change**:
- Usually keep default `["PA"]`
- Modify if firewall blocks specific flag combinations
- Use `["S", "A"]` to mimic SYN-ACK packets

### PCAP Buffer Size

Controls the packet capture buffer:

```yaml
pcap:
  sockbuf: 4194304  # 4MB (client default)
  sockbuf: 8388608  # 8MB (server default)
```

**Guidelines**:
- **Client**: 4MB usually sufficient
- **Server**: 8MB+ for high traffic
- **High throughput**: Increase to 16MB or 32MB
- **Low memory**: Decrease to 2MB

## Client Configuration

### SOCKS5 Proxy Mode

Run paqet as a SOCKS5 proxy:

```yaml
role: "client"

socks5:
  - listen: "127.0.0.1:1080"
    username: ""        # Optional authentication
    password: ""        # Optional authentication

server:
  addr: "10.0.0.100:9999"  # paqet server address
```

**Using the SOCKS5 Proxy**:
```bash
# curl with SOCKS5
curl -x socks5://127.0.0.1:1080 https://api.ipify.org

# Browser: Configure SOCKS5 proxy to 127.0.0.1:1080

# SSH through SOCKS5
ssh -o ProxyCommand='nc -X 5 -x 127.0.0.1:1080 %h %p' user@host
```

### Port Forwarding Mode

Forward specific ports through the tunnel:

```yaml
role: "client"

forward:
  - listen: "127.0.0.1:8080"    # Local port
    target: "192.168.1.50:80"   # Remote target (via server)
    protocol: "tcp"             # tcp or udp

  - listen: "127.0.0.1:5353"    # DNS forwarding
    target: "8.8.8.8:53"
    protocol: "udp"

server:
  addr: "10.0.0.100:9999"
```

**Usage**:
```bash
# Access forwarded HTTP
curl http://127.0.0.1:8080

# Use forwarded DNS
dig @127.0.0.1 -p 5353 example.com
```

### Multiple SOCKS5 Listeners

Run multiple SOCKS5 proxies simultaneously:

```yaml
socks5:
  - listen: "127.0.0.1:1080"
    username: "user1"
    password: "pass1"
  
  - listen: "127.0.0.1:1081"
    username: ""
    password: ""
  
  - listen: "0.0.0.0:1082"  # Listen on all interfaces
    username: "admin"
    password: "secret"
```

### Complete Client Example

```yaml
role: "client"

log:
  level: "info"

socks5:
  - listen: "127.0.0.1:1080"

network:
  interface: "en0"
  ipv4:
    addr: "192.168.1.100:0"
    router_mac: "aa:bb:cc:dd:ee:ff"
  tcp:
    local_flag: ["PA"]
    remote_flag: ["PA"]
  pcap:
    sockbuf: 4194304

server:
  addr: "10.0.0.100:9999"

transport:
  protocol: "kcp"
  conn: 4
  kcp:
    mode: "fast"
    mtu: 1350
    block: "aes"
    key: "your-secret-key-here"
```

## Server Configuration

### Basic Server Setup

```yaml
role: "server"

listen:
  addr: ":9999"  # Listen port

network:
  interface: "eth0"
  ipv4:
    addr: "10.0.0.100:9999"  # Must match listen port
    router_mac: "aa:bb:cc:dd:ee:ff"
  pcap:
    sockbuf: 8388608

transport:
  protocol: "kcp"
  kcp:
    mode: "fast"
    block: "aes"
    key: "your-secret-key-here"  # Must match client
```

**Important**:
- `listen.addr` port must match `network.ipv4.addr` port
- `transport.kcp.key` must match client exactly
- Use server's public IP if accessible from internet

### Listen Address Formats

```yaml
# Listen on all interfaces, port 9999
listen:
  addr: ":9999"

# Listen on specific IP, port 9999
listen:
  addr: "10.0.0.100:9999"

# Listen on IPv6
listen:
  addr: "[::]:9999"
```

### Complete Server Example

```yaml
role: "server"

log:
  level: "info"

listen:
  addr: ":9999"

network:
  interface: "eth0"
  ipv4:
    addr: "10.0.0.100:9999"
    router_mac: "aa:bb:cc:dd:ee:ff"
  tcp:
    local_flag: ["PA"]
  pcap:
    sockbuf: 8388608

transport:
  protocol: "kcp"
  conn: 1
  kcp:
    mode: "fast"
    mtu: 1350
    rcvwnd: 1024
    sndwnd: 1024
    block: "aes"
    key: "your-secret-key-here"
```

## Transport Settings

### KCP Configuration

```yaml
transport:
  protocol: "kcp"
  conn: 4  # Number of parallel connections (1-256)
  
  kcp:
    mode: "fast"       # normal, fast, fast2, fast3
    mtu: 1350          # Maximum transmission unit (50-1500)
    rcvwnd: 512        # Receive window (client: 512, server: 1024)
    sndwnd: 512        # Send window (client: 512, server: 1024)
    datashard: 10      # FEC data shards
    parityshard: 3     # FEC parity shards
    block: "aes"       # Encryption algorithm
    key: "secret"      # Shared secret key
```

### KCP Modes

| Mode | Latency | CPU | Bandwidth | Use Case |
|------|---------|-----|-----------|----------|
| `normal` | Highest | Low | Efficient | Stable networks |
| `fast` | Medium | Medium | Balanced | General use |
| `fast2` | Low | High | Higher | Gaming, VoIP |
| `fast3` | Lowest | Highest | Highest | Real-time apps |

**Recommendation**: Start with `fast`, adjust based on needs.

### MTU Configuration

Maximum transmission unit affects packet size:

```yaml
mtu: 1350  # Default, safe for most networks
```

**Guidelines**:
- **Standard networks**: 1350-1400
- **VPN/tunnel**: 1200-1300 (account for overhead)
- **Large packets**: 1450-1500 (if network supports)
- **High latency**: 1000-1200 (smaller packets)

**Test MTU**:
```bash
# Linux/macOS
ping -M do -s 1322 <server_ip>

# Windows
ping -f -l 1322 <server_ip>
```

### Window Sizes

Control send and receive buffers:

```yaml
kcp:
  rcvwnd: 512   # Receive window
  sndwnd: 512   # Send window
```

**Guidelines**:
- **Low latency**: 256-512
- **High throughput**: 1024-2048
- **Client default**: 512
- **Server default**: 1024

**Formula**: `Throughput ≈ (Window * MTU) / RTT`

### Forward Error Correction (FEC)

Recover lost packets without retransmission:

```yaml
kcp:
  datashard: 10     # Data shards
  parityshard: 3    # Parity shards
```

**Trade-offs**:
- **No FEC** (0/0): No overhead, relies on retransmission
- **Light FEC** (10/3): 30% overhead, good loss recovery
- **Heavy FEC** (10/10): 100% overhead, high loss tolerance

**When to use**:
- **No FEC**: Stable networks, low loss rate
- **Light FEC**: General use, moderate loss
- **Heavy FEC**: Unreliable networks, high loss rate

### Encryption Algorithms

```yaml
kcp:
  block: "aes"  # Encryption cipher
```

**Available Ciphers**:

| Cipher | Security | Speed | Key Size |
|--------|----------|-------|----------|
| `aes` | High | Fast | 256-bit |
| `aes-128` | High | Faster | 128-bit |
| `aes-192` | High | Fast | 192-bit |
| `chacha20` | High | Fast | 256-bit |
| `salsa20` | Medium | Very Fast | 256-bit |
| `blowfish` | Medium | Fast | Variable |
| `twofish` | High | Medium | 256-bit |
| `sm4` | High | Medium | 128-bit (Chinese) |
| `none` | ⚠️ None | Fastest | N/A |

**Recommendations**:
- **General use**: `aes` (good balance)
- **Maximum speed**: `aes-128` or `salsa20`
- **Maximum security**: `aes` or `chacha20`
- **Never use**: `none` (no encryption!)

### Connection Count

```yaml
transport:
  conn: 4  # Number of parallel KCP connections
```

**Guidelines**:
- **conn: 1**: Single connection, simple
- **conn: 2-4**: Load balancing, good for multiple streams
- **conn: 8-16**: High concurrency, many simultaneous connections
- **conn: 32+**: Very high load (requires more resources)

**Benefits**:
- Parallel connections reduce head-of-line blocking
- Better utilization of bandwidth
- Automatic failover if one connection fails

**Costs**:
- Higher memory usage
- More CPU for encryption
- Additional management overhead

## Advanced Options

### Logging

```yaml
log:
  level: "info"  # none, debug, info, warn, error, fatal
```

**Levels**:
- `debug`: All messages (verbose)
- `info`: Normal operation (default)
- `warn`: Warnings only
- `error`: Errors only
- `fatal`: Critical errors only
- `none`: No logging

### IPv6 Support

```yaml
network:
  ipv6:
    addr: "[2001:db8::1]:9999"
    router_mac: "aa:bb:cc:dd:ee:ff"
```

**Notes**:
- IPv6 address must be in brackets
- Gateway MAC same as IPv4 gateway (usually)
- Experimental support

### Windows-Specific: NPF GUID

```yaml
network:
  guid: "{12345678-1234-1234-1234-123456789012}"
```

**Find GUID**:
```powershell
Get-NetAdapter | Select-Object Name, InterfaceGuid
```

Required for Windows packet capture.

## Examples

### Example 1: Basic SOCKS5 Proxy

**Client**:
```yaml
role: "client"
log:
  level: "info"
socks5:
  - listen: "127.0.0.1:1080"
network:
  interface: "eth0"
  ipv4:
    addr: "192.168.1.100:0"
    router_mac: "aa:bb:cc:dd:ee:ff"
server:
  addr: "server.example.com:9999"
transport:
  protocol: "kcp"
  kcp:
    mode: "fast"
    block: "aes"
    key: "my-secret-key"
```

**Server**:
```yaml
role: "server"
log:
  level: "info"
listen:
  addr: ":9999"
network:
  interface: "eth0"
  ipv4:
    addr: "10.0.0.100:9999"
    router_mac: "aa:bb:cc:dd:ee:ff"
transport:
  protocol: "kcp"
  kcp:
    mode: "fast"
    block: "aes"
    key: "my-secret-key"
```

### Example 2: High-Performance Setup

```yaml
transport:
  conn: 8  # 8 parallel connections
  kcp:
    mode: "fast2"        # Low latency mode
    mtu: 1400            # Larger packets
    rcvwnd: 2048         # Large windows
    sndwnd: 2048
    datashard: 0         # No FEC (stable network)
    parityshard: 0
    block: "aes-128"     # Fast encryption
```

### Example 3: Unreliable Network

```yaml
transport:
  conn: 4
  kcp:
    mode: "normal"       # More conservative
    mtu: 1200            # Smaller packets
    rcvwnd: 512
    sndwnd: 512
    datashard: 10        # FEC enabled
    parityshard: 5       # 50% overhead
    block: "aes"
```

### Example 4: Port Forwarding

```yaml
role: "client"
forward:
  - listen: "127.0.0.1:3306"
    target: "192.168.1.100:3306"
    protocol: "tcp"
  - listen: "127.0.0.1:5432"
    target: "192.168.1.101:5432"
    protocol: "tcp"
```

---

For troubleshooting configuration issues, see [TROUBLESHOOTING.md](TROUBLESHOOTING.md).
