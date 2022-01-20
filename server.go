package main

import (
	"fmt"
	"net"
	"net/http"
)

func server(address net.IP, port uint) error {
	// tcp6 was attempted but wasn't supported by weather station
	listener, err := net.Listen("tcp4", listenAddress(address, port))
	if err != nil {
		return err
	}

	// Routes
	http.HandleFunc("/", catchall)
	http.HandleFunc("/wxigate/awp/v1", v1)

	_, err = fmt.Printf("%s-%s %f %f listening on: %s\n", callsign, ssid, longitude, latitude, listener.Addr().String())
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
