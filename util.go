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
	var str string
	for _, addr := range addrs {
		if ip_net, ok := addr.(*net.IPNet); ok && !ip_net.IP.IsLoopback() {
			if ip_net.IP.To4() != nil {

				str =  ip_net.IP.String()
				log.Println(str)
			}
		}
	}
	return "127.0.0.1"
}
