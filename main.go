package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
)

const DEFAULT_PORT uint = 8080 // BUG(high) move

var ipv4 bool
var ipv6 bool
var port uint
var address net.IP

func errorHandler(w http.ResponseWriter, req *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		fmt.Fprint(w, "custom 404")
	}

	fmt.Println(req)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}
	fmt.Fprint(w, "welcome home")
}

func helloHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println(req)
}

func listenNetwork(ipv4 bool, ipv6 bool) string {
	if ipv4 {
		return "tcp4"
	} else if ipv6 {
		return "tcp6"
	}

	return "tcp"
}

func main() {
	os.Exit(body())
}

func body() int {
	err := parseArgs()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	err = server(listenNetwork(ipv4, ipv6), listenAddress(address, port))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	return 0
}

func server(network, address string) error {
	listener, err := net.Listen(network, address)
	if err != nil {
		return err
	}

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/data", helloHandler)

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
	return fmt.Sprintf("%s:%d", address.String(), port)
}

func parseArgs() error {
	var inputAddress string
	flag.BoolVar(&ipv4, "4", false, "IPv4 only")
	flag.BoolVar(&ipv6, "6", false, "IPv6 only")
	flag.UintVar(&port, "port", DEFAULT_PORT, "tcp port (automatic 0)")
	flag.StringVar(&inputAddress, "address", "", "IP address")

	flag.Parse()

	// ipv4/ipv6 validation
	if ipv4 && ipv6 {
		return fmt.Errorf("ipv4 and ipv6 mutually exclusive")
	}

	// port validation
	const max_port = 65535
	if port > max_port {
		return fmt.Errorf("max port (%d)", max_port)
	}

	// address validation
	address = net.ParseIP(inputAddress)
	if address == nil {
		return fmt.Errorf("invalid address")
	}

	return nil
}
