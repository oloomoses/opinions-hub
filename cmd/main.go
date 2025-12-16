package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/oloomoses/opinions-hub/internal/config"
	"github.com/oloomoses/opinions-hub/internal/router"
	"github.com/oloomoses/opinions-hub/internal/server"
)

func main() {

	if err := config.Load(); err != nil {
		log.Fatal("Failed to load config", err)
	}

	sconf := config.LoadServerConfig()

	r := router.New()

	svr := server.New(sconf.Addr, r)

	svr.Start()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()
	svr.ShutDown(ctx)

}
