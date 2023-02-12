package service

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nimitsarup/restserver/config"
	"github.com/nimitsarup/restserver/db"
)

type ExternalServices struct {
	cfg        *config.Config
	db         db.UsersInMemoryDB
	httpServer HTTPServer
	router     *mux.Router
}

func NewServices(cfg *config.Config,
	router *mux.Router) (*ExternalServices, error) {
	e := &ExternalServices{
		cfg:    cfg,
		router: router,
	}

	return e, e.setup(cfg)
}

func (e *ExternalServices) GetHTTPServer() HTTPServer {
	return e.httpServer
}

func (e *ExternalServices) GetDB() db.UsersInMemoryDB {
	return e.db
}

func (e *ExternalServices) Shutdown(ctx context.Context) error {
	// any shutdown logic
	e.httpServer.Shutdown(ctx)
	return nil
}

func (e *ExternalServices) setup(cfg *config.Config) error {
	e.createDB()
	if !cfg.IsRunningInCloud {
		e.createHttpServer()
	}

	return nil
}

func (e *ExternalServices) createHttpServer() {
	s := &http.Server{
		Handler: e.router,
		Addr:    e.cfg.PortNumber,
	}
	e.httpServer = s
}

func (e *ExternalServices) createDB() {
	e.db = db.New()
}
