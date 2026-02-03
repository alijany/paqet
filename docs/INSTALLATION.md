# Installation Guide

Complete installation instructions for all supported platforms.

## Quick Install

### Pre-compiled Binaries (Recommended)

1. Download the latest release for your platform from the [Releases page](https://github.com/hanselime/paqet/releases)

2. Extract the archive:
   ```bash
   tar xzf paqet-linux-amd64.tar.gz
   # or
   unzip paqet-windows-amd64.zip
   ```

3. Move to system path (optional):
   ```bash
   # Linux/macOS
   sudo mv paqet /usr/local/bin/
   
   # Windows
   # Move paqet.exe to C:\Windows\System32\
   # Or add to PATH
   ```

4. Verify installation:
   ```bash
   paqet version
   ```

## Platform-Specific Installation

### Linux (Debian/Ubuntu)

**Prerequisites**:
```bash
sudo apt-get update
sudo apt-get install libpcap0.8
```

**Install from Binary**:
```bash
# Download
wget https://github.com/hanselime/paqet/releases/latest/download/paqet-linux-amd64.tar.gz

# Extract
tar xzf paqet-linux-amd64.tar.gz

# Install
sudo mv paqet /usr/local/bin/

# Set capabilities (avoids needing sudo every time)
sudo setcap cap_net_raw,cap_net_admin=eip /usr/local/bin/paqet

# Verify
paqet version
```

**Install from Source**:
```bash
# Install build dependencies
sudo apt-get install -y build-essential libpcap-dev git golang-go

# Clone and build
git clone https://github.com/hanselime/paqet.git
cd paqet
go build -o paqet cmd/main.go

# Install
sudo cp paqet /usr/local/bin/
sudo setcap cap_net_raw,cap_net_admin=eip /usr/local/bin/paqet
```

### Linux (RHEL/CentOS/Fedora)

**Prerequisites**:
```bash
sudo yum install libpcap
# or
sudo dnf install libpcap
```

**Install from Binary**:
```bash
# Download
wget https://github.com/hanselime/paqet/releases/latest/download/paqet-linux-amd64.tar.gz

# Extract
tar xzf paqet-linux-amd64.tar.gz

# Install
sudo mv paqet /usr/local/bin/
sudo setcap cap_net_raw,cap_net_admin=eip /usr/local/bin/paqet

# Verify
paqet version
```

### macOS

**Prerequisites**:

libpcap is included with macOS, but you need Xcode Command Line Tools:
```bash
xcode-select --install
```

**Install from Binary**:
```bash
# Download (Intel)
curl -LO https://github.com/hanselime/paqet/releases/latest/download/paqet-darwin-amd64.tar.gz

# Or for Apple Silicon (M1/M2)
curl -LO https://github.com/hanselime/paqet/releases/latest/download/paqet-darwin-arm64.tar.gz

# Extract
tar xzf paqet-darwin-*.tar.gz

# Install
sudo mv paqet /usr/local/bin/

# Verify
paqet version
```

**Using Homebrew** (if formula available):
```bash
brew install paqet
```

**Install from Source**:
```bash
# Install Go (if not installed)
brew install go

# Clone and build
git clone https://github.com/hanselime/paqet.git
cd paqet
go build -o paqet cmd/main.go

# Install
sudo cp paqet /usr/local/bin/
```

### Windows

**Prerequisites**:

1. **Install Npcap**:
   - Download from https://npcap.com/
   - Run installer
   - Enable "WinPcap API-compatible Mode"
   - Restart computer

**Install from Binary**:

1. Download `paqet-windows-amd64.zip` from [Releases](https://github.com/hanselime/paqet/releases)

2. Extract to desired location (e.g., `C:\Program Files\paqet\`)

3. Add to PATH (optional):
   ```powershell
   # PowerShell (as Administrator)
   $path = [Environment]::GetEnvironmentVariable("Path", "Machine")
   [Environment]::SetEnvironmentVariable("Path", "$path;C:\Program Files\paqet", "Machine")
   ```

4. Verify (in new terminal):
   ```powershell
   paqet version
   ```

**Install from Source**:

1. Install prerequisites:
   ```powershell
   # Using Chocolatey
   choco install -y golang git mingw
   ```

2. Install Npcap (see above)

3. Clone and build:
   ```powershell
   git clone https://github.com/hanselime/paqet.git
   cd paqet
   go build -o paqet.exe cmd/main.go
   ```

4. Copy to system path or use from current directory

### Docker

**Run as Container** (if Dockerfile available):

```bash
# Build image
docker build -t paqet .

# Run client
docker run --rm --net=host --cap-add=NET_ADMIN paqet run -c /config/client.yaml

# Run server
docker run -d --net=host --cap-add=NET_ADMIN paqet run -c /config/server.yaml
```

**Note**: Requires `--net=host` and `--cap-add=NET_ADMIN` for raw socket access.

## Post-Installation Setup

### Linux: Avoid Using sudo

Set capabilities to run without sudo:

```bash
sudo setcap cap_net_raw,cap_net_admin=eip /usr/local/bin/paqet
```

**Capabilities Explained**:
- `cap_net_raw`: Allows raw socket creation
- `cap_net_admin`: Allows network administration

**Security Note**: This grants network privileges to the binary. Only do this for trusted binaries.

**Alternative**: Create a systemd service to run as root (see below).

### macOS: Grant Permissions

macOS requires elevated privileges for raw socket access:

1. **Option 1**: Always run with `sudo`:
   ```bash
   sudo paqet run -c config.yaml
   ```

2. **Option 2**: Grant Full Disk Access:
   - System Preferences → Security & Privacy → Privacy
   - Full Disk Access → Add Terminal or iTerm

### Windows: Administrator Rights

Windows requires administrator privileges:

1. **Option 1**: Run Command Prompt/PowerShell as Administrator

2. **Option 2**: Set paqet.exe to always run as administrator:
   - Right-click paqet.exe → Properties
   - Compatibility tab
   - Check "Run this program as an administrator"
   - Apply

### Verify Npcap (Windows)

```powershell
# Check if Npcap service is running
Get-Service npcap

# Should show:
# Status   Name               DisplayName
# ------   ----               -----------
# Running  npcap              Npcap Packet Driver (NPCAP)
```

If not running:
```powershell
Start-Service npcap
```

## Configuration

### Get Example Configurations

```bash
# From source repository
wget https://raw.githubusercontent.com/hanselime/paqet/master/example/client.yaml.example
wget https://raw.githubusercontent.com/hanselime/paqet/master/example/server.yaml.example

# Or clone the repository
git clone https://github.com/hanselime/paqet.git
cd paqet/example
```

### Create Your Configuration

1. Copy example config:
   ```bash
   cp client.yaml.example my-client.yaml
   ```

2. Edit with your network details:
   ```yaml
   network:
     interface: "eth0"          # Your interface name
     ipv4:
       addr: "192.168.1.100:0"  # Your local IP
       router_mac: "..."        # Your gateway MAC
   server:
     addr: "server.example.com:9999"  # Server address
   transport:
     protocol: "kcp"
     conn: 1
     kcp:
       mode: "fast"
       key: "your-secret-key"   # Shared secret
   ```

3. See [CONFIGURATION.md](CONFIGURATION.md) for detailed guide

## Running paqet

### Client Mode

```bash
# Linux/macOS
sudo paqet run -c client.yaml

# Windows (as Administrator)
paqet.exe run -c client.yaml
```

### Server Mode

```bash
# Linux/macOS
sudo paqet run -c server.yaml

# Windows (as Administrator)
paqet.exe run -c server.yaml
```

### Check if Running

```bash
# Linux/macOS
ps aux | grep paqet

# Windows
Get-Process paqet
```

## Running as a Service

### Linux: systemd

**Client Service**:

Create `/etc/systemd/system/paqet-client.service`:
```ini
[Unit]
Description=paqet Client
After=network.target

[Service]
Type=simple
User=root
ExecStart=/usr/local/bin/paqet run -c /etc/paqet/client.yaml
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
```

**Server Service**:

Create `/etc/systemd/system/paqet-server.service`:
```ini
[Unit]
Description=paqet Server
After=network.target

[Service]
Type=simple
User=root
ExecStart=/usr/local/bin/paqet run -c /etc/paqet/server.yaml
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
```

**Enable and Start**:
```bash
# Copy config to /etc/paqet/
sudo mkdir -p /etc/paqet
sudo cp client.yaml /etc/paqet/

# Reload systemd
sudo systemctl daemon-reload

# Enable service
sudo systemctl enable paqet-client

# Start service
sudo systemctl start paqet-client

# Check status
sudo systemctl status paqet-client

# View logs
sudo journalctl -u paqet-client -f
```

### macOS: LaunchAgent

Create `~/Library/LaunchAgents/com.paqet.client.plist`:
```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.paqet.client</string>
    <key>ProgramArguments</key>
    <array>
        <string>/usr/local/bin/paqet</string>
        <string>run</string>
        <string>-c</string>
        <string>/etc/paqet/client.yaml</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardOutPath</key>
    <string>/var/log/paqet.log</string>
    <key>StandardErrorPath</key>
    <string>/var/log/paqet.err</string>
</dict>
</plist>
```

Load:
```bash
launchctl load ~/Library/LaunchAgents/com.paqet.client.plist
```

### Windows: Service

**Using NSSM** (Non-Sucking Service Manager):

1. Download NSSM from https://nssm.cc/

2. Install service:
   ```powershell
   # As Administrator
   nssm install paqet "C:\Program Files\paqet\paqet.exe" "run -c C:\paqet\client.yaml"
   ```

3. Start service:
   ```powershell
   Start-Service paqet
   ```

**Using Task Scheduler**:

1. Open Task Scheduler
2. Create Basic Task
3. Trigger: At startup
4. Action: Start a program
   - Program: `C:\Program Files\paqet\paqet.exe`
   - Arguments: `run -c C:\paqet\client.yaml`
5. Check "Run with highest privileges"

## Upgrading

### Upgrade Binary

```bash
# Stop running instance
sudo systemctl stop paqet-client  # Linux
# or
sudo killall paqet

# Download new version
wget https://github.com/hanselime/paqet/releases/latest/download/paqet-linux-amd64.tar.gz

# Extract and replace
tar xzf paqet-linux-amd64.tar.gz
sudo mv paqet /usr/local/bin/
sudo setcap cap_net_raw,cap_net_admin=eip /usr/local/bin/paqet

# Restart
sudo systemctl start paqet-client

# Verify new version
paqet version
```

### Configuration Migration

Check release notes for configuration changes. Generally:
1. Backup current config
2. Update config with new fields
3. Test configuration
4. Restart service

## Uninstallation

### Linux

```bash
# Stop service
sudo systemctl stop paqet-client
sudo systemctl disable paqet-client

# Remove binary
sudo rm /usr/local/bin/paqet

# Remove config
sudo rm -rf /etc/paqet

# Remove service file
sudo rm /etc/systemd/system/paqet-client.service
sudo systemctl daemon-reload
```

### macOS

```bash
# Unload service
launchctl unload ~/Library/LaunchAgents/com.paqet.client.plist

# Remove files
rm ~/Library/LaunchAgents/com.paqet.client.plist
sudo rm /usr/local/bin/paqet
sudo rm -rf /etc/paqet
```

### Windows

```powershell
# Stop and remove service (if using NSSM)
nssm stop paqet
nssm remove paqet confirm

# Remove binary
Remove-Item "C:\Program Files\paqet" -Recurse

# Uninstall Npcap (if desired)
# Control Panel → Programs → Uninstall Npcap
```

## Troubleshooting

See [TROUBLESHOOTING.md](TROUBLESHOOTING.md) for common issues and solutions.

### Quick Checks

```bash
# Check if paqet is installed
which paqet
paqet version

# Check if libpcap is available
# Linux
ldconfig -p | grep pcap

# macOS  
ls /usr/lib/libpcap*

# Windows
Get-Service npcap
```

### Common Issues

1. **Permission denied**: Run with sudo/administrator
2. **No such device**: Check interface name with `paqet iface`
3. **Connection timeout**: Verify server is running and reachable

## Getting Help

- Read the [documentation](README.md)
- Check [troubleshooting guide](TROUBLESHOOTING.md)
- Search [existing issues](https://github.com/hanselime/paqet/issues)
- Create a [new issue](https://github.com/hanselime/paqet/issues/new)

## Next Steps

After installation:
1. Configure paqet: [CONFIGURATION.md](CONFIGURATION.md)
2. Run and test your setup
3. Set up as a service (optional)
4. Read the [developer guide](DEVELOPER_GUIDE.md) to understand how it works

---

Successfully installed? Star the project on [GitHub](https://github.com/hanselime/paqet)! ⭐
