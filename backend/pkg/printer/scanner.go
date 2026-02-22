package printer

import (
	"fmt"
	"net"
	"sync"
	"time"
)

var (
	AvailableInterfaces = []string{"eth", "wlan", "wl", "Ethernet", "Wi-Fi"}
)

var (
	NoPrinterFound         = fmt.Errorf("no printer found")
	CannotConnectToPrinter = fmt.Errorf("cannot connect to printer")
)

type PrinterManipulator struct {
	PingDelay int
}

func NewPrinterManipulator(delay int) *PrinterManipulator {
	return &PrinterManipulator{
		PingDelay: delay,
	}
}

func (p *PrinterManipulator) Scan() ([]string, error) {
	ifaces, _ := net.Interfaces()
	var results []string
	var wg sync.WaitGroup
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		if !IsNeededInterface(AvailableInterfaces, iface.Name) {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, NoPrinterFound
		}
		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok || ipNet.IP.To4() == nil {
				continue
			}

			ips, _, _ := hosts(ipNet.String())
			for _, ip := range ips {
				wg.Go(func() {
					ok := checkPort9100(ip, p.PingDelay)
					if ok {
						results = append(results, ip)
					}
				})
			}
		}
	}
	wg.Wait()

	if len(results) == 0 {
		return nil, NoPrinterFound
	}
	return results, nil
}

func (p *PrinterManipulator) SendRequest(text string, ip string, port int) error {
	address := fmt.Sprintf("%v:%v", ip, port)
	fmt.Println(address)
	conn, err := net.DialTimeout("tcp", address, time.Duration(p.PingDelay)*time.Millisecond)
	if err != nil {
		return CannotConnectToPrinter
	}
	defer conn.Close()
	_, err = conn.Write([]byte(text))
	return err
}

func (p *PrinterManipulator) CheckPrinterIp(ip string) bool {
	conn, err := net.DialTimeout("tcp", ip+":9100", time.Duration(p.PingDelay)*time.Millisecond)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}
