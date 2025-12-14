package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/oloomoses/opinions-hub/internal/config"
	"github.com/oloomoses/opinions-hub/internal/database"
)

func main() {

	if err := config.Load(); err != nil {
		log.Fatal("Failed to load config", err)
	}

	db, err := database.Connect()

	if err != nil {
		log.Fatal("db error: ", err)
	}

	defer db.Close()

	srv := &http.Server{
		Addr: ":8080",
	}

	go func() {
		log.Println("server started")

		if err := srv.ListenAndServe(); err != nil {
			log.Print(err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Println("shutting down......")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	srv.Shutdown(ctx)
}
