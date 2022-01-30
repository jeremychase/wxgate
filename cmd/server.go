package cmd

import (
	"context"
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
	http.HandleFunc("/", catchall)
	http.HandleFunc("/wxigate/awp/v1", awpHandlerV1)

	_, err = fmt.Printf("%s-%s %f %f listening on: %s\n", opts.callsign, opts.ssid, opts.longitude, opts.latitude, listener.Addr().String())
	if err != nil {
		return err
	}

	srvr := &http.Server{
		BaseContext: func(_ net.Listener) context.Context {
			return ctxWithOptions(context.Background(), opts)
		},
	}

	// blocks until err
	err = srvr.Serve(listener)
	if err != nil {
		return err
	}

	return nil
}

func listenAddress(address net.IP, port uint) string {
	return fmt.Sprintf("[%s]:%d", address.String(), port)
}
