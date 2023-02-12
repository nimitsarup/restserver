package service

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/nimitsarup/restserver/api"
	"github.com/nimitsarup/restserver/config"
	"github.com/nimitsarup/restserver/db"
	"github.com/nimitsarup/restserver/handlers"
	"github.com/pkg/errors"
)

type Service struct {
	Server      HTTPServer
	Router      *mux.Router
	ServiceList ServiceContainer
	DB          db.UsersInMemoryDB
}

func Run(ctx context.Context,
	services ServiceContainer,
	svcErrors chan error,
	cfg *config.Config, r *mux.Router) (*Service, error) {

	log.Println("starting service")

	db := services.GetDB()
	api := &api.API{Handlers: &handlers.Handlers{DB: db}}

	r.Path("/users").HandlerFunc(api.AddUser).Methods(http.MethodPost)
	r.Path("/users").HandlerFunc(api.GetAllUsers).Methods(http.MethodGet)
	r.Path("/users/{id}").HandlerFunc(api.GetUser).Methods(http.MethodGet)

	s := services.GetHTTPServer()

	// Run the http server async
	go func() {
		log.Printf("listening at port [%s]", cfg.PortNumber)
		if err := s.ListenAndServe(); err != nil {
			svcErrors <- errors.Wrap(err, "failure in http listen and serve")
		}
	}()

	return &Service{Server: s, Router: r, ServiceList: services, DB: db}, nil
}

func (svc *Service) Close(ctx context.Context, timeout time.Duration) error {
	log.Println("commencing shutdown")
	ctx, cancel := context.WithTimeout(ctx, timeout)
	var err error
	go func() {
		defer cancel()
		err = svc.ServiceList.Shutdown(ctx)
	}()

	// wait for shutdown
	<-ctx.Done()

	// timeout expired
	if ctx.Err() == context.DeadlineExceeded {
		log.Println("shutdown timed out")
		return ctx.Err()
	}
	// other error
	if err != nil {
		log.Println(ctx, "failed to shutdown gracefully ", err)
		return err
	}

	log.Println("shutdown successful")
	return nil
}
