package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/gorilla/mux"
	"github.com/nimitsarup/restserver/config"
	"github.com/nimitsarup/restserver/service"
	"github.com/pkg/errors"
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		log.Fatal(nil, "fatal runtime error", err)
	}
}

func run(ctx context.Context) error {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill)

	// Read config
	cfg, err := config.Get()
	if err != nil {
		return errors.Wrap(err, "error getting configuration")
	}

	// Run the service, providing an error channel for fatal errors
	svcErrors := make(chan error, 1)
	r := mux.NewRouter().StrictSlash(true)

	svcList, err := service.NewServices(cfg, r)
	if err != nil {
		return errors.Wrap(err, "initialising services failed")

	}

	// Start service
	svc, err := service.Run(ctx, svcList, svcErrors, cfg, r)
	if err != nil {
		return errors.Wrap(err, "running service failed")
	}

	// blocks until an os interrupt or a fatal error occurs
	select {
	case err := <-svcErrors:
		return errors.Wrap(err, "service error received")
	case sig := <-signals:
		log.Println(fmt.Printf("os signal received %s recieved", sig.String()))
	}
	return svc.Close(ctx, cfg.ShutdownTimeout)
}
