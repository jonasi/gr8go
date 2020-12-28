package gr8http

import (
	"context"
	"net"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"

	gr8service "github.com/jonasi/gr8go/service"
)

type Server struct {
	Addr    []string
	Handler http.Handler

	onceInit sync.Once
	service  gr8service.Service
	server   *http.Server
	started  int32
}

func (s *Server) init(ctx context.Context) {
	s.onceInit.Do(func() {
		s.server = &http.Server{
			Handler: s.Handler,
			BaseContext: func(_ net.Listener) context.Context {
				return ctx
			},
		}

		s.service = gr8service.FromBlocking(s.start, s.stop)
	})
}

// AddListenAddr adds a new address to listen to when the server starts
func (s *Server) AddListenAddr(addr string) {
	if atomic.LoadInt32(&s.started) == 1 {
		panic("Attempting to add a listen addr after the server has started")
	}

	s.Addr = append(s.Addr, addr)
}

func (s *Server) start(ctx context.Context) error {
	atomic.StoreInt32(&s.started, 1)
	s.init(ctx)
	logServerStarting(ctx)

	return s.initListeners(ctx)
}

func (s *Server) initListeners(ctx context.Context) error {
	ls := []net.Listener{}

	for _, addr := range s.Addr {
		var (
			l   net.Listener
			err error
		)

		switch {
		case strings.HasPrefix(addr, "unix://"):
			l, err = net.Listen("unix", addr[7:])
		default:
			l, err = net.Listen("tcp", addr)
		}
		if err != nil {
			for _, l := range ls {
				if err := l.Close(); err != nil {
					logServerListeningCleanupErr(ctx, l.Addr(), err)
				}
			}
			return err
		}

		ls = append(ls, l)
	}

	var (
		serveCh   = make([]chan struct{}, len(ls))
		serveErrs = make([]error, len(ls))
		shutCh    = make(chan struct{})
	)

	for i, l := range ls {
		serveCh[i] = make(chan struct{})
		go func(l net.Listener, i int) {
			logServerListening(ctx, l.Addr())
			err := s.server.Serve(l)

			// normal shutdown
			if err == http.ErrServerClosed {
				err = nil
			}
			serveErrs[i] = err
			close(serveCh[i])
		}(l, i)
	}

	for i := 0; i < len(ls); i++ {
		<-serveCh[i]
		if err := serveErrs[i]; err != nil {
			return err
		}
	}

	<-shutCh
	return nil
}

// Stop stops the server
func (s *Server) stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Server) Start(ctx context.Context) error {
	s.init(ctx)
	return s.service.Start(ctx)
}

func (s *Server) Stop(ctx context.Context) error {
	s.init(ctx)
	return s.service.Stop(ctx)
}

func (s *Server) Wait(ctx context.Context) error {
	s.init(ctx)
	return s.service.Wait(ctx)
}
