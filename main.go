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

var Version = "-dev"

// Program options
var inputAddress string
var address net.IP
var port uint
var callsign string
var comment string
var ssid string
var longitude float64
var latitude float64
var showVersion bool

func main() {
	os.Exit(body())
}

func body() int {
	parseArgs()

	if showVersion {
		fmt.Printf("v%s\n", Version)
		return 0
	}

	err := validate()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	err = server(address, port)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	return 0
}

func parseArgs() {
	flag.BoolVar(&showVersion, "version", false, "show version")
	flag.Float64Var(&longitude, "longitude", 0.0, "longitude (decimal)")
	flag.Float64Var(&latitude, "latitude", 0.0, "latitude (decimal)")
	flag.StringVar(&callsign, "callsign", "", "callsign")
	flag.StringVar(&comment, "comment", "", "comment")
	flag.StringVar(&ssid, "ssid", "15", "ssid")
	flag.UintVar(&port, "port", DEFAULT_PORT, "tcp port (automatic 0)")
	flag.StringVar(&inputAddress, "address", DEFAULT_ADDRESS_IPV4, "IP address")

	flag.Parse()
}

func validate() error {
	// longitude and latitude validation
	if longitude == 0.0 {
		return fmt.Errorf("invalid longitude")
	}
	if latitude == 0.0 {
		return fmt.Errorf("invalid latitude")
	}

	// address validation
	address = net.ParseIP(inputAddress)
	if address == nil {
		return fmt.Errorf("invalid address")
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
