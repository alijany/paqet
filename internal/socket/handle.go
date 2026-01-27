package socket

import (
	"fmt"
	"net"
	"paqet/internal/conf"
	"runtime"
	"strings"

	"github.com/gopacket/gopacket/pcap"
)

func newHandle(cfg *conf.Network) (*pcap.Handle, error) {
	// On Windows, pcap requires device names in \Device\NPF_{GUID} format
	// We need to find the pcap device name that matches our interface
	deviceName := cfg.Interface.Name
	if runtime.GOOS == "windows" {
		var err error
		deviceName, err = findWindowsPcapDevice(cfg.Interface)
		if err != nil {
			return nil, fmt.Errorf("failed to find pcap device for %s: %v", cfg.Interface.Name, err)
		}
	}

	inactive, err := pcap.NewInactiveHandle(deviceName)
	if err != nil {
		return nil, fmt.Errorf("failed to create inactive pcap handle for %s: %v", cfg.Interface.Name, err)
	}
	defer inactive.CleanUp()

	if err = inactive.SetBufferSize(cfg.PCAP.Sockbuf); err != nil {
		return nil, fmt.Errorf("failed to set pcap buffer size to %d: %v", cfg.PCAP.Sockbuf, err)
	}

	if err = inactive.SetSnapLen(65536); err != nil {
		return nil, fmt.Errorf("failed to set pcap snap length: %v", err)
	}
	if err = inactive.SetPromisc(true); err != nil {
		return nil, fmt.Errorf("failed to enable promiscuous mode: %v", err)
	}
	if err = inactive.SetTimeout(pcap.BlockForever); err != nil {
		return nil, fmt.Errorf("failed to set pcap timeout: %v", err)
	}
	if err = inactive.SetImmediateMode(true); err != nil {
		return nil, fmt.Errorf("failed to enable immediate mode: %v", err)
	}

	handle, err := inactive.Activate()
	if err != nil {
		return nil, fmt.Errorf("failed to activate pcap handle on %s: %v", cfg.Interface.Name, err)
	}

	return handle, nil
}

// findWindowsPcapDevice finds the pcap device name for a Windows network interface
// by matching the MAC address from the Go interface to pcap devices
func findWindowsPcapDevice(iface *net.Interface) (string, error) {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		return "", fmt.Errorf("failed to enumerate pcap devices: %v", err)
	}

	// Try to match by MAC address (hardware address)
	targetMAC := iface.HardwareAddr.String()
	
	for _, device := range devices {
		// Get addresses for this device to check if it matches
		for _, addr := range device.Addresses {
			// Try to match the interface by checking if any address matches
			if addr.IP != nil {
				// Get the Go interface for this IP to check MAC
				if goIface, err := net.InterfaceByName(iface.Name); err == nil {
					addrs, _ := goIface.Addrs()
					for _, a := range addrs {
						ipnet, ok := a.(*net.IPNet)
						if ok && ipnet.IP.Equal(addr.IP) {
							return device.Name, nil
						}
					}
				}
			}
		}
		
		// Alternative: Try matching by description containing the interface name
		// Windows pcap devices often have descriptions like "Intel(R) Wi-Fi 6E AX210 160MHz"
		if strings.Contains(strings.ToLower(device.Description), strings.ToLower(iface.Name)) {
			return device.Name, nil
		}
	}

	// If no match found, try a more flexible approach by checking all device names
	// that might contain the interface index
	ifaceIndex := fmt.Sprintf("%d", iface.Index)
	for _, device := range devices {
		if strings.Contains(device.Name, ifaceIndex) {
			return device.Name, nil
		}
	}

	return "", fmt.Errorf("no pcap device found matching interface %s (MAC: %s)", iface.Name, targetMAC)
}
