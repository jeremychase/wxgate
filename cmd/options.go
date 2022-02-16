package cmd

import (
	"net"
)

type options struct {
	argAddress  string
	callsign    string
	comment     string
	ssid        string
	address     net.IP
	longitude   float64
	latitude    float64
	port        uint
	showVersion bool
}
