package server

import (
	"backend/internal/app"
	"backend/internal/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Run(a *app.App) error {
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.CORS())

	service.NewGreeter(a).Mount(e) // one line per service

	return e.Start(":" + a.Config.Port)
}
