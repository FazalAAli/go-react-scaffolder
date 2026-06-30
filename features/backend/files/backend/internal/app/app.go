package app

import (
	"backend/internal/config"

	"github.com/labstack/echo/v5"
	// scaffold:region:app-imports:start
	// scaffold:region:app-imports:end
)

// App is the dependency container. Features attach as fields via the
// app-fields region and are constructed in New via the app-init region.
type App struct {
	Config *config.Config
	Mounts []func(*echo.Echo)
	// scaffold:region:app-fields:start
	// scaffold:region:app-fields:end
}

func New(cfg *config.Config) (*App, error) {
	a := &App{Config: cfg}
	// scaffold:region:app-init:start
	// scaffold:region:app-init:end
	return a, nil
}

// Close releases resources held by features. Features register cleanup via the
// app-close region; it runs on graceful shutdown.
func (a *App) Close() {
	// scaffold:region:app-close:start
	// scaffold:region:app-close:end
}
