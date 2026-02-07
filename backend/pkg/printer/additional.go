package printer

import (
    "net"
    "strings"
	"time"
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

func hosts(cidr string) ([]string, int, error) {
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

func checkPort9100(ip string, delay int) bool {
	conn, err := net.DialTimeout("tcp", ip+":9100", time.Duration(delay)*time.Millisecond)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}
