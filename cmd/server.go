package cmd

import (
	"fmt"
	"net"
	"net/http"
	"wxigate/requestlog"
)

func server(opts options) error {
	// tcp6 was attempted but wasn't supported by weather station
	listener, err := net.Listen("tcp4", listenAddress(opts.address, opts.port))
	if err != nil {
		return err
	}

	mux := http.NewServeMux()

	// Routes
	mux.HandleFunc("/", defaultHandler)
	mux.Handle("/wxigate/awp/v1", awpHandlerV1(opts))

	_, err = fmt.Printf("%s-%s %f %f listening on: %s\n", opts.callsign, opts.ssid, opts.longitude, opts.latitude, listener.Addr().String())
	if err != nil {
		return err
	}

	var h http.Handler

	// Request logging middleware
	if opts.requestlog != nil {
		h = requestlog.NewReqLog(mux, *opts.requestlog)
	} else {
		h = mux
	}

	// blocks until err
	err = http.Serve(listener, h)
	if err != nil {
		return err
	}

	return nil
}

func listenAddress(address net.IP, port uint) string {
	return fmt.Sprintf("[%s]:%d", address.String(), port)
}
