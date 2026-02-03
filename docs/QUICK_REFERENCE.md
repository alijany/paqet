# Quick Reference

Fast reference guide for common paqet operations.

## Installation

```bash
# Download binary
wget https://github.com/hanselime/paqet/releases/latest/download/paqet-linux-amd64.tar.gz
tar xzf paqet-linux-amd64.tar.gz
sudo mv paqet /usr/local/bin/

# Set capabilities (Linux)
sudo setcap cap_net_raw,cap_net_admin=eip /usr/local/bin/paqet
```

## Basic Commands

```bash
# Show version
paqet version

# List network interfaces
paqet iface

# Generate secret key
paqet secret

# Run client
sudo paqet run -c client.yaml

# Run server
sudo paqet run -c server.yaml

# Dump packets
sudo paqet dump -i eth0

# Test connectivity
sudo paqet ping -c client.yaml
```

## Finding Network Information

```bash
# Linux - interface and IP
ip addr show

# Linux - gateway
ip route | grep default

# Linux - gateway MAC
arp -n <gateway_ip>

# macOS - interface and IP
ifconfig

# macOS - gateway  
netstat -rn | grep default

# Windows - all info
ipconfig /all

# Windows - interface GUID
Get-NetAdapter | Select-Object Name, InterfaceGuid
```

## Minimal Client Config

```yaml
role: "client"
socks5:
  - listen: "127.0.0.1:1080"
network:
  interface: "eth0"
  ipv4:
    addr: "192.168.1.100:0"
    router_mac: "aa:bb:cc:dd:ee:ff"
server:
  addr: "10.0.0.100:9999"
transport:
  protocol: "kcp"
  conn: 1
  kcp:
    mode: "fast"
    key: "your-secret-key"
    # Optional: block: "aes"
```

## Minimal Server Config

```yaml
role: "server"
listen:
  addr: ":9999"
network:
  interface: "eth0"
  ipv4:
    addr: "10.0.0.100:9999"
    router_mac: "aa:bb:cc:dd:ee:ff"
transport:
  protocol: "kcp"
  conn: 1
  kcp:
    mode: "fast"
    key: "your-secret-key"
    # Optional: block: "aes"
```

## Using SOCKS5 Proxy

```bash
# curl
curl -x socks5://127.0.0.1:1080 https://api.ipify.org

# wget
wget -e use_proxy=yes -e socks_proxy=127.0.0.1:1080 https://example.com

# SSH
ssh -o ProxyCommand='nc -X 5 -x 127.0.0.1:1080 %h %p' user@host

# Firefox: Settings → Network Settings → Manual proxy
# SOCKS Host: 127.0.0.1, Port: 1080, SOCKS v5
```

## Configuration Examples

### High Performance

```yaml
transport:
  conn: 8
  kcp:
    mode: "fast2"
    mtu: 1400
    rcvwnd: 2048
    sndwnd: 2048
    block: "aes-128"
```

### Lossy Network

```yaml
transport:
  conn: 4
  kcp:
    mode: "normal"
    mtu: 1200
    datashard: 10
    parityshard: 3
```

### Port Forwarding

```yaml
forward:
  - listen: "127.0.0.1:3306"
    target: "192.168.1.100:3306"
    protocol: "tcp"
  - listen: "127.0.0.1:5353"
    target: "8.8.8.8:53"
    protocol: "udp"
```

## Troubleshooting

```bash
# Check if running
ps aux | grep paqet

# Check listening ports
netstat -tlnp | grep paqet

# Enable debug logging
# In config: log.level: "debug"

# Test with tcpdump
sudo tcpdump -i eth0 port 9999 -X

# Check firewall
sudo ufw status
sudo iptables -L -n

# Verify config
cat config.yaml
```

## Common Errors

| Error | Solution |
|-------|----------|
| Permission denied | Run with `sudo` or set capabilities |
| No such device | Check interface name: `paqet iface` |
| Connection timeout | Verify server running and reachable |
| Invalid MAC address | Use `arp -n <gateway_ip>` |
| Address in use | Kill existing process: `sudo killall paqet` |

## KCP Modes

| Mode | Latency | Use Case |
|------|---------|----------|
| normal | High | Stable networks |
| fast | Medium | General use |
| fast2 | Low | Gaming, VoIP |
| fast3 | Lowest | Real-time apps |

## Encryption Ciphers

| Cipher | Security | Speed |
|--------|----------|-------|
| aes | High | Fast |
| aes-128 | High | Faster |
| chacha20 | High | Fast |
| salsa20 | Medium | Very Fast |
| none | ⚠️ None | Fastest |

## systemd Service

```ini
[Unit]
Description=paqet Client
After=network.target

[Service]
Type=simple
User=root
ExecStart=/usr/local/bin/paqet run -c /etc/paqet/client.yaml
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl enable paqet-client
sudo systemctl start paqet-client
sudo systemctl status paqet-client
sudo journalctl -u paqet-client -f
```

## Performance Tuning

```yaml
# Increase connections
transport:
  conn: 8  # More parallel connections

# Larger windows
kcp:
  rcvwnd: 2048
  sndwnd: 2048

# Larger MTU (if supported)
kcp:
  mtu: 1400

# Larger buffers
pcap:
  sockbuf: 8388608  # 8MB
```

## Port Checks

```bash
# Check if port is open
nc -zv <host> <port>

# Check listening ports
# Linux
sudo ss -tlnp
sudo netstat -tlnp

# macOS
netstat -an | grep LISTEN
lsof -nP -iTCP -sTCP:LISTEN

# Windows
netstat -ano
```

## Useful Links

- [Documentation Index](README.md)
- [Developer Guide](DEVELOPER_GUIDE.md)
- [Configuration Guide](CONFIGURATION.md)
- [Troubleshooting](TROUBLESHOOTING.md)
- [GitHub Repository](https://github.com/hanselime/paqet)

## Quick Debugging

1. Enable debug logging: `log.level: "debug"`
2. Check interface: `paqet iface`
3. Test network: `ping <server_ip>`
4. Check firewall: `sudo iptables -L`
5. Verify config: Keys match, ports match, IPs correct
6. Check logs: `journalctl -u paqet-client -f`

---

For detailed information, see the full documentation in the [docs/](.) folder.
