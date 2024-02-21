package server

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"time"

	"http-everything/httpe/pkg/actions/servedirectory"
	"http-everything/httpe/pkg/config"
	"http-everything/httpe/pkg/middleware"
	"http-everything/httpe/pkg/requesthandler"
	"http-everything/httpe/pkg/rules"
	"http-everything/httpe/pkg/share/logger"

	"github.com/gorilla/mux"

	"github.com/gorilla/handlers"
)

// Server is the core HTTP server for the HTTPE server
type Server struct {
	cfg   *config.Config
	rules *[]rules.Rule

	Handler http.Handler

	srv             *http.Server
	logger          *logger.Logger
	accessLogWriter io.Writer
}

// New creates a new Server. It will also create a new baseLogger which will be used to fork
// the loggers used by other packages.
func New(cfg *config.Config, rules *[]rules.Rule, baseLogger *logger.Logger, accessLogWriter io.Writer) (
	svr *Server, err error) {
	l := baseLogger.Fork("server")

	svr = &Server{
		cfg:             cfg,
		logger:          l,
		accessLogWriter: accessLogWriter,
		rules:           rules,
	}
	return svr, nil
}

// Setup creates the routes and sets up the http.Server
func (s *Server) Setup() {
	s.logger.Infof("setting up")
	r := mux.NewRouter()
	for _, rule := range *s.rules {
		h := requesthandler.Execute(rule, s.logger)
		m := middleware.New(rule, s.logger)
		if len(rule.On.Methods) == 0 {
			r.Handle(rule.On.Path, m.Collection(h))
		} else {
			for _, method := range rule.On.Methods {
				r.Handle(rule.On.Path, m.Collection(h)).Methods(method)
			}
		}
		if rule.Action() == rules.ServeDirectory {
			r.PathPrefix(rule.On.Path).Handler(m.Collection(servedirectory.Handle(rule.On.Path, rule.Do.ServeDirectory)))
		}
	}
	r.PathPrefix("/").Handler(http.HandlerFunc(s.catchAllHandler))

	if s.accessLogWriter != nil {
		accessLogHandler := handlers.CombinedLoggingHandler(s.accessLogWriter, r)
		s.Handler = accessLogHandler
	} else {
		s.Handler = r
	}

	var tlscfg *tls.Config
	if s.cfg.S.KeyFile != "" && s.cfg.S.CertFile != "" {
		tlscfg = &tls.Config{
			MinVersion:       tls.VersionTLS12,
			CurvePreferences: []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		}
	}

	s.srv = &http.Server{
		Addr:              s.cfg.S.Address,
		Handler:           s.Handler,
		IdleTimeout:       5 * time.Minute,
		ReadHeaderTimeout: time.Minute,
		TLSConfig:         tlscfg,
	}
}

// Serve starts a go routine with http.Server in http or https mode depending on the config settings.
// Will block waiting for ctrl+c or ctx done if requested by the withWait param.
func (s *Server) Serve(ctx context.Context, withWait bool) (err error) {
	stop := make(chan error)
	go func(stop chan error) {
		if s.cfg.S.KeyFile != "" && s.cfg.S.CertFile != "" {
			s.logger.Infof("listening on https://%s", s.cfg.S.Address)
			err = s.srv.ListenAndServeTLS(s.cfg.S.CertFile, s.cfg.S.KeyFile)
		} else {
			s.logger.Infof("listening on http://%s", s.cfg.S.Address)
			err = s.srv.ListenAndServe()
		}

		if !errors.Is(err, http.ErrServerClosed) {
			// stop with an err if we get here
			stop <- fmt.Errorf("http server stopped: %w", err)
		}
	}(stop)

	if withWait {
		err = s.WaitUntilDone(ctx, stop)
		if err != nil {
			return err
		}
	}
	return nil
}

// catchAllHandler returns unauthorised for all unknown routes
func (s *Server) catchAllHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	_, err := w.Write([]byte("not found\n"))
	if err != nil {
		s.logger.Errorf("unable to write response: %s", err)
	}
}

// Shutdown performs a clean shutdown of the http.Server
func (s *Server) Shutdown() {
	s.logger.Infof("shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	_ = s.srv.Shutdown(shutdownCtx)
}

// WaitUntilDone waits until ctrl+c, ctx done or stopped
func (s *Server) WaitUntilDone(ctx context.Context, stop chan error) (err error) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err = <-stop:
			return err
		case <-c:
			s.logger.Infof("ctrl+c received")
			return nil
		}
	}
}
