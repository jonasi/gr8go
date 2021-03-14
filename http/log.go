package gr8http

import (
	"context"
	"net"

	gr8log "github.com/jonasi/gr8go/log"
)

func logServerStarting(ctx context.Context) {
	gr8log.Info(ctx, "Starting server")
}

func logServerListening(ctx context.Context, addr net.Addr) {
	gr8log.Info(ctx, "Server listening", gr8log.Args{
		"addr": addr.String(),
	})
}

func logServerListeningCleanupErr(ctx context.Context, addr net.Addr, err error) {
	gr8log.Error(ctx, "Attempting to cleanup listeners, but encountered error for listener", gr8log.Args{
		"addr":  addr.String(),
		"error": err,
	})
}
