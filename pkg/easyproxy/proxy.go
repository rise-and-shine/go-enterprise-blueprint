package easyproxy

import (
	"context"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sort"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/observability/logger"
)

type EasyProxy struct {
	server *http.Server
	mux    *http.ServeMux
}

func New(cfg Config) (*EasyProxy, error) {
	mux := http.NewServeMux()

	routes := make([]Route, len(cfg.Routes))
	copy(routes, cfg.Routes)
	sort.Slice(routes, func(i, j int) bool {
		return len(routes[i].Path) > len(routes[j].Path)
	})

	for _, route := range routes {
		target, err := url.Parse(route.Target)
		if err != nil {
			return nil, errx.Wrap(err)
		}

		proxy := httputil.NewSingleHostReverseProxy(target)

		originalDirector := proxy.Director
		proxy.Director = func(req *http.Request) {
			originalDirector(req)
			req.Host = target.Host
		}

		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			logger.With("error", err, "path", r.URL.Path, "target", target.String()).
				Error("Proxy error")
			w.WriteHeader(http.StatusBadGateway)
			_, _ = w.Write([]byte(`{"error": "Bad gateway"}`))
		}

		logger.With("path", route.Path, "target", route.Target).Info("Registered route")
		mux.Handle(route.Path, proxy)
	}

	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status": "OK"}`))
	})

	server := &http.Server{
		Addr:    cfg.Listen,
		Handler: mux,
	}

	return &EasyProxy{
		server: server,
		mux:    mux,
	}, nil
}

func (p *EasyProxy) Run(ctx context.Context) error {
	errCh := make(chan error, 1)

	go func() {
		logger.With("addr", p.server.Addr).Info("Starting proxy server")
		if err := p.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- errx.Wrap(err)
		}
		close(errCh)
	}()

	select {
	case <-ctx.Done():
		logger.Info("Shutting down proxy server")
		if err := p.server.Shutdown(context.Background()); err != nil {
			return errx.Wrap(err)
		}
		return nil
	case err := <-errCh:
		return err
	}
}
