package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/horiondreher/go-parking-lot/parkings/internal/adapters/grpc"
	"github.com/horiondreher/go-parking-lot/parkings/internal/adapters/http/httpv1"
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

	httpServer, err := httpv1.NewHTTPAdapter()
	if err != nil {
		log.Err(err).Msg("error creating server")
		stop()
	}

	gRPCServer := grpc.NewAdapter()

	// starts the server in a goroutine to let the main goroutine listen for the interrupt signal
	go func() {
		if err := httpServer.Start(); err != nil && err != http.ErrServerClosed {
			log.Err(err).Msg("error starting http server")
		}
	}()

	go func() {
		if err := gRPCServer.Start(); err != nil {
			log.Err(err).Msg("error starting grpc server")
		}
	}()

	<-ctx.Done()

	// gracefully shutdown the server
	httpServer.Shutdown()

	log.Info().Msg("server stopped")
}
