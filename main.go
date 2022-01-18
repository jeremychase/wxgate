package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

const DEFAULT_PORT uint = 8080 // BUG(high) move
const DEFAULT_ADDRESS_IPV4 string = "0.0.0.0"
const DEFAULT_ADDRESS_IPV6 string = "::"

var address net.IP
var ipv4 bool
var ipv6 bool
var port uint
var callsign string
var ssid string

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
	flag.StringVar(&callsign, "callsign", "", "callsign")
	flag.StringVar(&ssid, "ssid", "15", "ssid")
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

	// ssid validation
	if len(ssid) > 2 {
		return fmt.Errorf("ssid too long")
	} else if len(ssid) < 1 {
		return fmt.Errorf("ssid empty")
	}

	// callsign validation
	if len(callsign) == 0 {
		return fmt.Errorf("missing callsign")
	}
	if len(callsign) > 8 { // BUG(medium) fix
		return fmt.Errorf("callsign too long")
	} else if len(callsign) < 3 { // BUG(medium) fix
		return fmt.Errorf("callsign too short")
	}

	callsign = strings.Trim(callsign, "-")
	ssid = strings.Trim(ssid, "-")

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
