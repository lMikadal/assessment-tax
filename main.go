package main

import (
	"context"
	"crypto/subtle"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lMikadal/assessment-tax/postgres"
	"github.com/lMikadal/assessment-tax/tax"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	db, err := postgres.New()
	if err != nil {
		panic(err)
	}

	handler := tax.New(db)
	e := echo.New()
	e.POST("/tax/calculations", handler.TaxHandler)
	e.POST("tax/calculations/upload-csv", handler.UploadCSVHandler)

	a := e.Group("/admin")
	a.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if subtle.ConstantTimeCompare([]byte(username), []byte(os.Getenv("ADMIN_USERNAME"))) == 1 && subtle.ConstantTimeCompare([]byte(password), []byte(os.Getenv("ADMIN_PASSWORD"))) == 1 {
			return true, nil
		}

		return false, nil
	}))
	a.POST("/deductions/personal", handler.TaxDeducateHandler)
	a.POST("/deductions/k-receipt", handler.TaxDeducateKreceiptHandler)

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
