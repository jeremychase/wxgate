package main

import (
	"flag"
	"fmt"
	"net"
	"os"
)

const DEFAULT_PORT uint = 8080 // BUG(high) move
const DEFAULT_ADDRESS_IPV4 string = "0.0.0.0"
const DEFAULT_ADDRESS_IPV6 string = "::"

var ipv4 bool
var ipv6 bool
var port uint
var address net.IP

func main() {
	os.Exit(body())
}

func body() int {
	err := parseArgs()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	err = server(listenNetwork(ipv4, ipv6), address, port)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	return 0
}

func parseArgs() error {
	var inputAddress string
	flag.BoolVar(&ipv4, "4", false, "IPv4 only")
	flag.BoolVar(&ipv6, "6", false, "IPv6 only")
	flag.UintVar(&port, "port", DEFAULT_PORT, "tcp port (automatic 0)")
	flag.StringVar(&inputAddress, "address", DEFAULT_ADDRESS_IPV4, "IP address")

	flag.Parse()

	// ipv4/ipv6 validation
	if ipv4 && ipv6 {
		return fmt.Errorf("ipv4 and ipv6 mutually exclusive")
	}

	// default to ipv4
	if !ipv6 {
		ipv4 = true
	}

	// address validation
	address = net.ParseIP(inputAddress)
	if address == nil {
		return fmt.Errorf("invalid address")
	}

	// fix default address for ipv6
	if ipv6 && net.ParseIP(DEFAULT_ADDRESS_IPV4).Equal(address) {
		address = net.ParseIP(DEFAULT_ADDRESS_IPV6)
		if address == nil {
			return fmt.Errorf("invalid default ipv6 address")
		}
	}

	// port validation
	const max_port = 65535
	if port > max_port {
		return fmt.Errorf("max port (%d)", max_port)
	}

	return nil
}

func listenNetwork(ipv4 bool, ipv6 bool) string {
	if ipv4 {
		return "tcp4"
	} else if ipv6 {
		return "tcp6"
	}

	return "tcp"
}
