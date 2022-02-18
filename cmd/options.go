package cmd

import (
	"net"

	"github.com/ebarkie/aprs"
)

type options struct {
	aprsSource  aprs.Addr
	callsign    string
	comment     string
	dial        string
	ssid        string
	argAddress  string
	address     net.IP
	dialpass    int
	longitude   float64
	latitude    float64
	port        uint
	showVersion bool
}
