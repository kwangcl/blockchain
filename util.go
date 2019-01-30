package main

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

/*
func GetOutboundIP() string {
    conn, err := net.Dial("udp", "8.8.8.8:80")
    HandleError("net.Dial: ",err)
    defer conn.Close()
    localAddr := conn.LocalAddr().String()
    idx := strings.LastIndex(localAddr, ":")
    return localAddr[0:idx]
}*/
