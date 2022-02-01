package cmd

import (
	"context"
	"errors"
	"net"
)

type optionKey string
type options struct {
	inputAddress string // BUG(low) remove
	address      net.IP
	port         uint
	callsign     string
	comment      string
	ssid         string
	longitude    float64
	latitude     float64
	showVersion  bool
}

func ctxWithOptions(ctx context.Context, opts options) context.Context {
	key := optionKey("options")
	return context.WithValue(context.Background(), key, opts)
}

func ctxOptions(ctx context.Context) (options, error) {
	opts, ok := ctx.Value(optionKey("options")).(options)
	if !ok {
		return opts, errors.New("unable to return options")
	}

	return opts, nil
}
