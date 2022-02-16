package cmd

import (
	"fmt"
	"net"
	"net/http"
)

func server(opts options) error {
	// tcp6 was attempted but wasn't supported by weather station
	listener, err := net.Listen("tcp4", listenAddress(opts.address, opts.port))
	if err != nil {
		return err
	}

	// Routes
	http.HandleFunc("/", defaultHandler)
	http.Handle("/wxigate/awp/v1", awpHandlerV1(opts))

	_, err = fmt.Printf("%s-%s %f %f listening on: %s\n", opts.callsign, opts.ssid, opts.longitude, opts.latitude, listener.Addr().String())
	if err != nil {
		return err
	}

	// blocks until err
	err = http.Serve(listener, nil)
	if err != nil {
		return err
	}

	return nil
}

func listenAddress(address net.IP, port uint) string {
	return fmt.Sprintf("[%s]:%d", address.String(), port)
}
