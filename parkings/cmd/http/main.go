package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	httpV1 "github.com/horiondreher/go-parking-lot/parkings/internal/adapters/http/v1"
	"github.com/horiondreher/go-parking-lot/parkings/internal/utils"

	"github.com/rs/zerolog/log"
)

var interruptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func main() {
	os.Setenv("TZ", "UTC")

	utils.StartLogger()

	// creates a new context with a cancel function that is called when the interrupt signal is received
	ctx, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	defer stop()

	server, err := httpV1.NewHTTPAdapter()
	if err != nil {
		log.Err(err).Msg("error creating server")
		stop()
	}

	// starts the server in a goroutine to let the main goroutine listen for the interrupt signal
	go func() {
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			log.Err(err).Msg("error starting server")
		}
	}()

	<-ctx.Done()

	// gracefully shutdown the server
	server.Shutdown()

	log.Info().Msg("server stopped")
}
