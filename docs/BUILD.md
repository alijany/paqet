# Building from Source

Complete guide to building paqet from source code.

## Prerequisites

### Required Software

1. **Go 1.25 or later**
   - Download from [golang.org](https://golang.org/dl/)
   - Verify: `go version`

2. **Git**
   - Download from [git-scm.com](https://git-scm.com/)
   - Verify: `git --version`

3. **C Compiler** (for CGO)
   - **Linux**: gcc (usually pre-installed)
   - **macOS**: Xcode Command Line Tools
   - **Windows**: MinGW or TDM-GCC

4. **libpcap Development Libraries**
   - **Linux (Debian/Ubuntu)**:
     ```bash
     sudo apt-get install libpcap-dev
     ```
   - **Linux (RHEL/CentOS)**:
     ```bash
     sudo yum install libpcap-devel
     ```
   - **macOS**:
     ```bash
     xcode-select --install
     ```
   - **Windows**:
     - Install [Npcap](https://npcap.com/) with SDK

## Quick Build

```bash
# Clone repository
git clone https://github.com/hanselime/paqet.git
cd paqet

# Download dependencies
go mod download

# Build
go build -o paqet cmd/main.go

# Verify
./paqet version
```

## Detailed Build Instructions

### Linux

```bash
# Install dependencies
sudo apt-get update
sudo apt-get install -y build-essential libpcap-dev git

# Install Go (if not installed)
wget https://go.dev/dl/go1.25.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.25.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# Clone and build
git clone https://github.com/hanselime/paqet.git
cd paqet
go mod download
go build -o paqet cmd/main.go

# Optional: Install to system
sudo cp paqet /usr/local/bin/
sudo setcap cap_net_raw,cap_net_admin=eip /usr/local/bin/paqet
```

### macOS

```bash
# Install Xcode Command Line Tools
xcode-select --install

# Install Go (using Homebrew)
brew install go

# Or download from https://go.dev/dl/

# Clone and build
git clone https://github.com/hanselime/paqet.git
cd paqet
go mod download
go build -o paqet cmd/main.go

# Optional: Install to system
sudo cp paqet /usr/local/bin/
```

### Windows

**Using PowerShell**:

```powershell
# Install chocolatey (if not installed)
# See https://chocolatey.org/install

# Install dependencies
choco install -y golang git mingw

# Install Npcap
# Download from https://npcap.com/ and install

# Clone and build
git clone https://github.com/hanselime/paqet.git
cd paqet
go mod download
go build -o paqet.exe cmd/main.go

# Verify
.\paqet.exe version
```

## Build Options

### Standard Build

```bash
go build -o paqet cmd/main.go
```

### Optimized Build

```bash
# Smaller binary, stripped symbols
go build -ldflags="-s -w" -o paqet cmd/main.go
```

### Static Binary (Linux)

```bash
# Fully static binary (no shared libraries)
CGO_ENABLED=1 go build -ldflags="-linkmode external -extldflags -static" -o paqet cmd/main.go
```

### Cross-Compilation

**Linux → Windows**:
```bash
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go build -o paqet.exe cmd/main.go
```

**macOS → Linux**:
```bash
GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o paqet cmd/main.go
```

> **Note**: Cross-compilation with CGO is complex and may require additional setup.

### Build with Version Information

```bash
VERSION=$(git describe --tags --always)
BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
go build -ldflags="-X main.Version=$VERSION -X main.BuildDate=$BUILD_DATE" -o paqet cmd/main.go
```

## Development Build

### With Debug Symbols

```bash
# Keep debug symbols for debugging
go build -gcflags="all=-N -l" -o paqet cmd/main.go
```

### With Race Detector

```bash
# Detect race conditions (slower, for testing)
go build -race -o paqet cmd/main.go
```

## Build for Multiple Platforms

### Build Script

Create `build.sh`:
```bash
#!/bin/bash
VERSION=$(git describe --tags --always)

# Linux AMD64
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o dist/paqet-linux-amd64 cmd/main.go

# Linux ARM64
GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o dist/paqet-linux-arm64 cmd/main.go

# macOS AMD64
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o dist/paqet-darwin-amd64 cmd/main.go

# macOS ARM64 (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o dist/paqet-darwin-arm64 cmd/main.go

# Windows AMD64
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o dist/paqet-windows-amd64.exe cmd/main.go

echo "Build complete: $VERSION"
```

Run:
```bash
chmod +x build.sh
./build.sh
```

## Testing the Build

### Basic Test

```bash
# Show version
./paqet version

# List interfaces
./paqet iface

# Generate secret
./paqet secret
```

### Integration Test

```bash
# Create test configuration
cp example/server.yaml.example test-server.yaml
cp example/client.yaml.example test-client.yaml

# Edit configurations with your network details
# Run server (in one terminal)
sudo ./paqet run -c test-server.yaml

# Run client (in another terminal)
sudo ./paqet run -c test-client.yaml

# Test SOCKS5 proxy
curl -x socks5://127.0.0.1:1080 https://api.ipify.org
```

## Troubleshooting Build Issues

### Error: `pcap.h: No such file or directory`

**Solution**: Install libpcap development libraries

```bash
# Linux
sudo apt-get install libpcap-dev

# macOS
xcode-select --install

# Windows
# Install Npcap with SDK from https://npcap.com/
```

### Error: `gcc: command not found`

**Solution**: Install C compiler

```bash
# Linux
sudo apt-get install build-essential

# macOS
xcode-select --install

# Windows
# Install MinGW from https://www.mingw-w64.org/
```

### Error: Module Issues

**Solution**: Clean and re-download modules

```bash
go clean -modcache
go mod download
go mod verify
```

### Error: CGO Disabled

**Solution**: Enable CGO

```bash
export CGO_ENABLED=1
go build -o paqet cmd/main.go
```

## Installation

### System-Wide Installation (Linux/macOS)

```bash
# Copy binary to system path
sudo cp paqet /usr/local/bin/

# Set capabilities (Linux only, avoids need for sudo)
sudo setcap cap_net_raw,cap_net_admin=eip /usr/local/bin/paqet

# Verify
paqet version
```

### User Installation

```bash
# Create local bin directory
mkdir -p ~/bin

# Copy binary
cp paqet ~/bin/

# Add to PATH (add to ~/.bashrc or ~/.zshrc)
export PATH=$PATH:~/bin

# Reload shell
source ~/.bashrc
```

### Windows Installation

```powershell
# Copy to a directory in PATH
Copy-Item paqet.exe C:\Windows\System32\

# Or create a dedicated directory
New-Item -Path "C:\Program Files\paqet" -ItemType Directory
Copy-Item paqet.exe "C:\Program Files\paqet\"

# Add to PATH
[Environment]::SetEnvironmentVariable("Path", $env:Path + ";C:\Program Files\paqet", "Machine")
```

## Development Setup

### Clone for Development

```bash
# Fork the repository on GitHub first
git clone https://github.com/YOUR_USERNAME/paqet.git
cd paqet

# Add upstream remote
git remote add upstream https://github.com/hanselime/paqet.git

# Create feature branch
git checkout -b feature/my-feature
```

### Install Development Tools

```bash
# Go tools
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Optional: Air for hot reload
go install github.com/cosmtrek/air@latest
```

### Code Formatting

```bash
# Format all Go files
go fmt ./...

# Or use goimports
goimports -w .
```

### Linting

```bash
# Run golangci-lint
golangci-lint run

# Fix auto-fixable issues
golangci-lint run --fix
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./internal/socket/

# Run with verbose output
go test -v ./...
```

## Dependency Management

### Update Dependencies

```bash
# Update all dependencies to latest minor/patch versions
go get -u ./...

# Update specific dependency
go get -u github.com/xtaci/kcp-go/v5

# Tidy dependencies
go mod tidy
```

### Verify Dependencies

```bash
# Verify checksums
go mod verify

# View dependency graph
go mod graph

# Check for updates
go list -u -m all
```

## Creating a Release

### Tag Version

```bash
# Create and push tag
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

### Build Release Binaries

```bash
# Build for all platforms
./build.sh

# Create archives
cd dist
tar czf paqet-linux-amd64.tar.gz paqet-linux-amd64
tar czf paqet-darwin-amd64.tar.gz paqet-darwin-amd64
zip paqet-windows-amd64.zip paqet-windows-amd64.exe
```

### Generate Checksums

```bash
cd dist
sha256sum * > SHA256SUMS
```

## Next Steps

- Read [DEVELOPER_GUIDE.md](DEVELOPER_GUIDE.md) for architecture overview
- Read [CONTRIBUTING.md](CONTRIBUTING.md) for contribution guidelines
- Check [TROUBLESHOOTING.md](TROUBLESHOOTING.md) if you encounter issues

---

For questions about building, open an issue on GitHub.
