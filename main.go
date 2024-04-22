package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lMikadal/assessment-tax/postgres"
	"github.com/lMikadal/assessment-tax/tax"
	"github.com/labstack/echo/v4"
)

func main() {
	db, err := postgres.New()
	if err != nil {
		panic(err)
	}

	e := echo.New()
	handler := tax.New(db)
	e.POST("/tax/calculations", handler.TaxHandler)

	// Start server
	go func() {
		post := ":" + os.Getenv("PORT")
		if err := e.Start(post); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	<-shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
