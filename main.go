package main

import (
	"flag"
	"github.com/pterm/pterm"
	"net"
	"time"
)

func main() {

	ipAddress := GetLocalIP()
	ip := flag.String("ip", ipAddress, "Server IP address")
	port := flag.Int("port", 9001, "Server Port Address")
	flag.Parse()
	ScreenLogo()

	serverInfo, _ := pterm.DefaultSpinner.Start("Server Starting ...")
	time.Sleep(time.Second * 2)
	CreateAndStartServer(*ip, *port, serverInfo)

}

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
