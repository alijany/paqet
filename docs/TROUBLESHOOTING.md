# Troubleshooting Guide

Solutions to common issues when running paqet.

## Table of Contents

1. [Quick Diagnostics](#quick-diagnostics)
2. [Installation Issues](#installation-issues)
3. [Connection Problems](#connection-problems)
4. [Performance Issues](#performance-issues)
5. [Configuration Errors](#configuration-errors)
6. [Platform-Specific Issues](#platform-specific-issues)
7. [Debugging Tools](#debugging-tools)

## Quick Diagnostics

### Check Installation

```bash
# Verify paqet is installed
./paqet version

# Check if libpcap is available (Linux)
ldconfig -p | grep pcap

# List network interfaces
./paqet iface
```

### Test Configuration

```bash
# Validate configuration file
./paqet run -c config.yaml --dry-run  # If available

# Enable debug logging
# Set log.level: "debug" in config.yaml
```

### Test Connectivity

```bash
# Ping test (if implemented)
sudo ./paqet ping -c config.yaml

# Network reachability
ping <server_ip>
traceroute <server_ip>
```

## Installation Issues

### Error: `libpcap not found`

**Symptoms**:
```
error: pcap.h: No such file or directory
```

**Solution**:

**Linux (Debian/Ubuntu)**:
```bash
sudo apt-get update
sudo apt-get install libpcap-dev
```

**Linux (RHEL/CentOS/Fedora)**:
```bash
sudo yum install libpcap-devel
# or
sudo dnf install libpcap-devel
```

**macOS**:
```bash
xcode-select --install
```

**Windows**:
- Download and install [Npcap](https://npcap.com/)
- Enable "WinPcap API-compatible Mode" during installation

### Error: `go: directive requires go 1.25`

**Symptoms**:
```
go: paqet requires go >= 1.25
```

**Solution**:
```bash
# Check current Go version
go version

# Update Go
# Download from https://golang.org/dl/
# Or use version manager like gvm or asdf
```

### Build Errors on Windows

**Symptoms**:
```
undefined: pcap.OpenLive
```

**Solution**:
1. Install Npcap with SDK
2. Set CGO environment:
   ```powershell
   $env:CGO_ENABLED = "1"
   ```
3. Ensure gcc is in PATH (install MinGW)

## Connection Problems

### Error: `Permission denied`

**Symptoms**:
```
failed to create raw packet conn: permission denied
```

**Solution**:

Raw sockets require elevated privileges.

**Linux/macOS**:
```bash
sudo ./paqet run -c config.yaml
```

**Windows**:
- Run PowerShell or CMD as Administrator
- Right-click → "Run as administrator"

**Alternative (Linux only)**:
```bash
# Give binary raw socket capabilities
sudo setcap cap_net_raw,cap_net_admin=eip ./paqet

# Now run without sudo
./paqet run -c config.yaml
```

### Error: `No such device`

**Symptoms**:
```
failed to create receive handle on eth0: no such device
```

**Solution**:

Interface name is incorrect.

```bash
# List all interfaces
./paqet iface

# Update config.yaml with correct interface name
network:
  interface: "correct-name-here"
```

**Common Interface Names**:
- Linux: `eth0`, `ens3`, `wlan0`, `wlp3s0`
- macOS: `en0`, `en1`
- Windows: `Ethernet`, `Wi-Fi`

### Error: `Connection timeout`

**Symptoms**:
- Client can't connect to server
- No response from server

**Troubleshooting Steps**:

1. **Verify server is running**:
   ```bash
   # On server machine
   ps aux | grep paqet
   ```

2. **Check network reachability**:
   ```bash
   # From client machine
   ping <server_ip>
   nc -zv <server_ip> <port>
   ```

3. **Check firewall rules**:
   ```bash
   # Linux - allow port
   sudo ufw allow 9999
   sudo iptables -A INPUT -p tcp --dport 9999 -j ACCEPT
   
   # Check if port is listening
   sudo netstat -tlnp | grep 9999
   ```

4. **Verify configuration**:
   - Client `server.addr` matches server IP/port
   - Server `listen.addr` port matches `network.ipv4.addr` port
   - Both use same `transport.kcp.key`

### Error: `Invalid MAC address`

**Symptoms**:
```
invalid MAC address format
```

**Solution**:

```bash
# Find correct gateway MAC
# Linux
ip r | grep default  # Get gateway IP
arp -n <gateway_ip>  # Get gateway MAC

# macOS
netstat -rn | grep default  # Get gateway IP
arp -n <gateway_ip>         # Get gateway MAC

# Windows
ipconfig /all              # Get gateway IP (Default Gateway)
arp -a <gateway_ip>        # Get gateway MAC

# Update config.yaml
network:
  ipv4:
    router_mac: "aa:bb:cc:dd:ee:ff"  # Use lowercase, colon-separated
```

### KCP Connection Fails

**Symptoms**:
- Connection established but no data transfer
- Immediate disconnection

**Solutions**:

1. **Verify encryption key matches**:
   ```yaml
   # Client and server MUST have identical keys
   transport:
     kcp:
       key: "exact-same-key-on-both-sides"
   ```

2. **Check MTU settings**:
   ```yaml
   # Try smaller MTU
   kcp:
     mtu: 1200  # Reduce from default 1350
   ```

3. **Verify cipher compatibility**:
   ```yaml
   # Use simple cipher for testing
   kcp:
     block: "aes"  # Same on both sides
   ```

### SOCKS5 Connection Refused

**Symptoms**:
```
curl: (7) Failed to connect to 127.0.0.1 port 1080: Connection refused
```

**Solution**:

1. **Verify client is running**:
   ```bash
   ps aux | grep paqet
   ```

2. **Check SOCKS5 listening**:
   ```bash
   netstat -tln | grep 1080
   # or
   lsof -i :1080
   ```

3. **Verify configuration**:
   ```yaml
   socks5:
     - listen: "127.0.0.1:1080"  # Make sure this is configured
   ```

4. **Check logs**:
   ```yaml
   log:
     level: "debug"  # Enable debug logging
   ```

## Performance Issues

### Slow Transfer Speed

**Symptoms**:
- Much slower than direct connection
- High latency

**Solutions**:

1. **Increase connection count**:
   ```yaml
   transport:
     conn: 4  # Or 8, 16 for high concurrency
   ```

2. **Use faster KCP mode**:
   ```yaml
   kcp:
     mode: "fast2"  # Or "fast3" for lowest latency
   ```

3. **Increase window sizes**:
   ```yaml
   kcp:
     rcvwnd: 2048
     sndwnd: 2048
   ```

4. **Optimize MTU**:
   ```yaml
   kcp:
     mtu: 1400  # Increase if network supports
   ```

5. **Use faster encryption**:
   ```yaml
   kcp:
     block: "aes-128"  # Faster than aes-256
   ```

### High Packet Loss

**Symptoms**:
- Frequent reconnections
- Choppy performance
- Data corruption

**Solutions**:

1. **Enable FEC**:
   ```yaml
   kcp:
     datashard: 10
     parityshard: 3
   ```

2. **Increase buffer sizes**:
   ```yaml
   pcap:
     sockbuf: 8388608  # Increase to 8MB
   ```

3. **Use more conservative mode**:
   ```yaml
   kcp:
     mode: "normal"  # More reliable
   ```

4. **Reduce MTU**:
   ```yaml
   kcp:
     mtu: 1200  # Smaller packets
   ```

### High CPU Usage

**Symptoms**:
- 100% CPU utilization
- System slowdown

**Solutions**:

1. **Reduce connection count**:
   ```yaml
   transport:
     conn: 2  # Fewer connections
   ```

2. **Use more efficient mode**:
   ```yaml
   kcp:
     mode: "normal"  # Less aggressive
   ```

3. **Disable FEC if not needed**:
   ```yaml
   kcp:
     datashard: 0
     parityshard: 0
   ```

4. **Lighter encryption**:
   ```yaml
   kcp:
     block: "salsa20"  # Very fast
   ```

### High Memory Usage

**Symptoms**:
- Excessive RAM consumption
- Out of memory errors

**Solutions**:

1. **Reduce buffer sizes**:
   ```yaml
   pcap:
     sockbuf: 2097152  # 2MB instead of 4MB
   ```

2. **Reduce window sizes**:
   ```yaml
   kcp:
     rcvwnd: 256
     sndwnd: 256
   ```

3. **Reduce connection count**:
   ```yaml
   transport:
     conn: 1  # Single connection
   ```

## Configuration Errors

### Error: `role must be 'client' or 'server'`

**Solution**:
```yaml
role: "client"  # Must be exactly "client" or "server"
```

### Error: `failed to parse config`

**Symptoms**:
```
yaml: unmarshal errors: ...
```

**Solutions**:

1. **Check YAML syntax**:
   - Correct indentation (2 spaces)
   - No tabs (use spaces only)
   - Proper quotes around strings

2. **Validate YAML**:
   ```bash
   # Use online YAML validator
   # Or install yamllint
   yamllint config.yaml
   ```

3. **Common mistakes**:
   ```yaml
   # Wrong (tabs)
   network:
   	interface: "eth0"
   
   # Correct (spaces)
   network:
     interface: "eth0"
   
   # Wrong (invalid value)
   kcp:
     mode: fast  # Missing quotes
   
   # Correct
   kcp:
     mode: "fast"
   ```

### Error: `validation failed`

**Symptoms**:
- Config loads but validation fails
- Specific field errors

**Solutions**:

Check each configuration section:

1. **Listen port mismatch (server)**:
   ```yaml
   # Wrong
   listen:
     addr: ":9999"
   network:
     ipv4:
       addr: "10.0.0.100:8888"  # Different port!
   
   # Correct
   listen:
     addr: ":9999"
   network:
     ipv4:
       addr: "10.0.0.100:9999"  # Same port
   ```

2. **Invalid mode**:
   ```yaml
   # Wrong
   kcp:
     mode: "ultra-fast"
   
   # Correct (must be: normal, fast, fast2, fast3)
   kcp:
     mode: "fast3"
   ```

3. **Invalid cipher**:
   ```yaml
   # Check spelling
   kcp:
     block: "aes"  # Not "AES" or "aes-256"
   ```

## Platform-Specific Issues

### Linux Issues

**Error: `Operation not permitted`**

Solution:
```bash
# Run with sudo
sudo ./paqet run -c config.yaml

# Or set capabilities
sudo setcap cap_net_raw,cap_net_admin=eip ./paqet
```

**Error: `RTNETLINK answers: Operation not supported`**

Solution:
- Update kernel
- Check interface exists: `ip link show`

### macOS Issues

**Error: `No suitable device found`**

Solution:
```bash
# Grant terminal/app Full Disk Access
# System Preferences → Security & Privacy → Privacy → Full Disk Access
# Add Terminal.app or iTerm.app

# Run with sudo
sudo ./paqet run -c config.yaml
```

**Error: `kern.boottime failed`**

Solution:
- Update macOS
- Reinstall Xcode Command Line Tools

### Windows Issues

**Error: `NPF driver isn't running`**

Solution:
1. Install Npcap from https://npcap.com/
2. During installation, enable "WinPcap API-compatible Mode"
3. Restart computer
4. Verify service is running:
   ```powershell
   Get-Service npcap
   # Should show "Running"
   ```

**Error: `The handle is invalid`**

Solution:
```yaml
# Add NPF GUID to config
network:
  guid: "{YOUR-INTERFACE-GUID}"

# Find GUID:
# Get-NetAdapter | Select-Object Name, InterfaceGuid
```

**Permissions Error on Windows**

Solution:
- Run PowerShell as Administrator
- Run CMD as Administrator
- Check User Account Control (UAC) settings

## Debugging Tools

### Enable Debug Logging

```yaml
log:
  level: "debug"
```

This shows detailed information about:
- Packet send/receive
- Connection establishment
- Protocol messages
- Error details

### Packet Capture

**Using paqet built-in dump**:
```bash
sudo ./paqet dump -i eth0
```

**Using tcpdump**:
```bash
# Capture all traffic on port 9999
sudo tcpdump -i eth0 port 9999 -X

# Save to file
sudo tcpdump -i eth0 port 9999 -w capture.pcap

# Analyze with Wireshark
wireshark capture.pcap
```

### Network Interface Info

```bash
# List all interfaces with details
./paqet iface

# Show detailed interface info
# Linux
ip addr show
ip link show

# macOS
ifconfig -a

# Windows
ipconfig /all
Get-NetAdapter | Format-List
```

### Check Open Ports

```bash
# Linux
sudo netstat -tlnp
sudo ss -tlnp

# macOS
netstat -an | grep LISTEN
lsof -nP -iTCP -sTCP:LISTEN

# Windows
netstat -ano
Get-NetTCPConnection | Where-Object {$_.State -eq "Listen"}
```

### Test Encryption

```bash
# Generate a test secret key
./paqet secret

# Verify key is consistent on both sides
# Keys must match exactly (case-sensitive)
```

### Monitor System Resources

```bash
# Linux
top -p $(pgrep paqet)
htop -p $(pgrep paqet)

# macOS
top -pid $(pgrep paqet)

# Windows
Get-Process paqet | Format-List
```

### Check Firewall Rules

**Linux (iptables)**:
```bash
sudo iptables -L -n -v
sudo iptables -L INPUT -v
```

**Linux (ufw)**:
```bash
sudo ufw status verbose
```

**macOS**:
```bash
sudo /usr/libexec/ApplicationFirewall/socketfilterfw --getglobalstate
```

**Windows**:
```powershell
Get-NetFirewallRule | Where-Object {$_.Enabled -eq 'True'}
netsh advfirewall show allprofiles
```

## Common Error Messages

### `failed to establish connection`

**Causes**:
- Server not running
- Incorrect server address
- Firewall blocking
- Network unreachable

**Solutions**: See [Connection Problems](#connection-problems)

### `invalid packet: checksum mismatch`

**Causes**:
- Network corruption
- Encryption key mismatch
- MTU too large

**Solutions**:
1. Verify encryption keys match
2. Reduce MTU
3. Enable FEC

### `context deadline exceeded`

**Causes**:
- Operation timeout
- Network latency too high
- Server not responding

**Solutions**:
1. Check network connectivity
2. Increase timeout (if configurable)
3. Verify server is running

### `address already in use`

**Causes**:
- Port already bound
- Previous instance still running
- Another application using port

**Solutions**:
```bash
# Find process using port
# Linux/macOS
lsof -i :1080
sudo kill <PID>

# Windows
netstat -ano | findstr :1080
taskkill /PID <PID> /F
```

## Getting Help

If you can't resolve your issue:

1. **Check Existing Issues**: Search [GitHub Issues](https://github.com/hanselime/paqet/issues)

2. **Gather Information**:
   - Operating system and version
   - paqet version (`./paqet version`)
   - Go version (`go version`)
   - Configuration file (remove sensitive data)
   - Full error messages
   - Debug logs

3. **Create an Issue**:
   - Describe the problem clearly
   - Include steps to reproduce
   - Attach relevant logs
   - Specify your environment

4. **Community Help**:
   - GitHub Discussions
   - Related forums or chat channels

---

**Remember**: Most issues are configuration-related. Double-check:
- Interface names
- IP addresses
- MAC addresses
- Port numbers
- Encryption keys

Run with `log.level: "debug"` for detailed troubleshooting information.
