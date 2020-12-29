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

func logRouterHandlingRoute(ctx context.Context, method string, path string, mws []string) {
	if method == "*" {
		gr8log.Info(ctx, "Handling all methods for route", gr8log.Args{
			"path":        path,
			"middlewares": mws,
		})
	} else {
		gr8log.Info(ctx, "Handling route", gr8log.Args{
			"method":      method,
			"path":        path,
			"middlewares": mws,
		})
	}
}

func logRouterHandlingNotFound(ctx context.Context, middlewares []string) {
	gr8log.Info(ctx, "Handling not found", gr8log.Args{
		"middlewares": middlewares,
	})
}

func logTemplateHandlerError(ctx context.Context, err error) {
	gr8log.Error(ctx, "Template render error", gr8log.Args{
		"error": err,
	})
}

func logSPAIndexHandlerError(ctx context.Context, err error) {
	gr8log.Error(ctx, "Error making index handler", gr8log.Args{
		"error": err,
	})
}
