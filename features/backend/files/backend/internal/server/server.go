package server

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"backend/internal/app"
	"backend/internal/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	// scaffold:region:server-imports:start
	// scaffold:region:server-imports:end
)

func Run(a *app.App) error {
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.CORS())
	// scaffold:region:server-middleware:start
	// scaffold:region:server-middleware:end

	service.NewGreeter(a).Mount(e) // one line per service

	go func() {
		if err := e.Start(":" + a.Config.Port); err != nil && !errors.Is(err, http.ErrServerClosed) {
			e.Logger.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return e.Shutdown(ctx)
}
