package main
////////////test commit
import (
	"log"
	"net"
)

func ErrorHandler(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func GetIPv4Address() string{
	addrs, err := net.InterfaceAddrs()
    ErrorHandler(err)
	for _, addr := range addrs {
		if ip_net, ok := addr.(*net.IPNet); ok && !ip_net.IP.IsLoopback() {
			if ip_net.IP.To4() != nil {
				return ip_net.IP.String() 
			}
		}
	}
	return "127.0.0.1"
}