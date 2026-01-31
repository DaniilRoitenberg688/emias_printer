package main

import (
	"fmt"
	"net"
	"strings"
	"time"
)

var (
	AvailableInterfaces = []string{"eth", "wlan", "wl", "Ethernet", "Wi-Fi"}
)

func IsNeededInterface(aifcs []string, ifc string) bool {
	for _, aifs := range aifcs {
		if strings.Contains(ifc, aifs) {
			return true
		}
	}
	return false
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func Hosts(cidr string) ([]string, int, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, 0, err
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}

	// remove network address and broadcast address
	lenIPs := len(ips)
	switch {
	case lenIPs < 2:
		return ips, lenIPs, nil

	default:
		return ips[1 : len(ips)-1], lenIPs - 2, nil
	}
}

func checkPort9100(ip string) bool {
	conn, err := net.DialTimeout("tcp", ip+":9100", 10*time.Millisecond)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func main() {
	ifaces, _ := net.Interfaces()
	for _, iface := range ifaces {

		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		if !IsNeededInterface(AvailableInterfaces, iface.Name) {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok || ipNet.IP.To4() == nil {
				continue
			}

			ips, _, _ := Hosts(ipNet.String())
			for _, ip := range ips {
				ok = checkPort9100(ip)
				if ok {
					fmt.Println(ip)
				}
			}

		}
	}
}
