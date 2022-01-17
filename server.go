package main

import (
	"fmt"
	"net"
	"net/http"
)

func server(network string, address net.IP, port uint) error {
	listener, err := net.Listen(network, listenAddress(address, port))
	if err != nil {
		return err
	}

	// Routes
	http.HandleFunc("/", catchall)
	http.HandleFunc("/wxigate/v1", v1)

	_, err = fmt.Printf("Listening on: %s\n", listener.Addr().String())
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
