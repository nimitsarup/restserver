package steps

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/nimitsarup/restserver/config"
	"github.com/nimitsarup/restserver/db"
	"github.com/nimitsarup/restserver/service"
	sMock "github.com/nimitsarup/restserver/service/mock"
)

type UsersApiComponent struct {
	svc            *service.Service
	errorChan      chan error
	Config         *config.Config
	HTTPServer     *http.Server
	HttpResponse   *http.Response
	ServiceRunning bool
	DB             db.UsersInMemoryDB
}

func NewUsersApiComponent() (*UsersApiComponent, error) {
	c := &UsersApiComponent{
		HTTPServer:     &http.Server{},
		errorChan:      make(chan error),
		ServiceRunning: false,
	}

	c.DB = db.New()
	var err error
	c.Config, err = config.Get()
	if err != nil {
		return nil, err
	}

	log.Printf("configuration for component test [%v]", c.Config)

	return c, nil
}

func (c *UsersApiComponent) Initialiser() (http.Handler, error) {
	r := &mux.Router{}
	cfg, _ := config.Get()
	svcList := &sMock.ServiceContainerMock{
		GetDBFunc:         func() db.UsersInMemoryDB { return c.DB },
		GetHTTPServerFunc: func() service.HTTPServer { return c.HTTPServer },
	} //service.NewServices(cfg, r)
	c.svc, _ = service.Run(context.Background(), svcList, c.errorChan, cfg, r)
	c.svc.DB = c.DB

	return c.svc.Router, nil
}

func (c *UsersApiComponent) Reset() error {
	return nil
}

func (c *UsersApiComponent) Close() error {
	ctx := context.Background()
	if c.svc != nil && c.ServiceRunning {
		if err := c.svc.Close(ctx, 10*time.Second); err != nil {
			return err
		}
		c.ServiceRunning = false
	}
	return nil
}
