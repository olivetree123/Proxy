package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"proxy"
)

var (
	localAddr  = flag.String("l", ":9999", "local address")
	remoteAddr = flag.String("r", "localhost:80", "remote address")
)

func main() {
	flag.Parse()
	log.Println(fmt.Sprintf("Proxying from %v to %v", *localAddr, *remoteAddr))
	laddr, err := net.ResolveTCPAddr("tcp", *localAddr)
	if err != nil {
		log.Println(fmt.Sprintf("Failed to resolve local address: %s", err))
		os.Exit(1)
	}
	raddr, err := net.ResolveTCPAddr("tcp", *remoteAddr)
	if err != nil {
		log.Println(fmt.Sprintf("Failed to resolve remote address: %s", err))
		os.Exit(1)
	}
	listener, _ := net.ListenTCP("tcp", laddr)
	if err != nil {
		log.Println(fmt.Sprintf("Failed to open local port to listen: %s", err))
		os.Exit(1)
	}
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Println(fmt.Sprintf("Failed to accept connection '%s'", err))
			continue
		}
		var p *proxy.Proxy
		p = proxy.New(conn, laddr, raddr)
		go p.Start()
	}
}
