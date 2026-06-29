package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"backend/internal/app"
	"backend/internal/service"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	// scaffold:region:server-imports:start
	// scaffold:region:server-imports:end
)

func Run(a *app.App) error {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	// scaffold:region:server-middleware:start
	// scaffold:region:server-middleware:end

	service.NewGreeter(a).Mount(e) // one line per service

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	sc := echo.StartConfig{
		Address:         ":" + a.Config.Port,
		HideBanner:      true,
		GracefulTimeout: 10 * time.Second,
	}
	return sc.Start(ctx, e)
}
