package utils

import (
	"bytes"
	"net"
	"os/exec"
	"strings"
)

func FindNetworkInterface() (string, error) {
	cmd := exec.Command("ip", "link", "show")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	output := stdout.String()
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "state UP") {
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				return strings.TrimSpace(parts[1]), nil
			}
		}
	}

	return "", nil
}

func FindServerIP(intf string) (string, error) {
	
	iface, err := net.InterfaceByName(intf)
	if err != nil {
		return "", err
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return "", err
	}

	var ip net.IP
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			ip = ipNet.IP.To4()
			break
		}
	}

	return ip.String(), nil
}
