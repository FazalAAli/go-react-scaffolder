package app

import (
	"backend/internal/config"
	// scaffold:region:app-imports:start
	// scaffold:region:app-imports:end
)

// App is the dependency container. Features attach as fields via the
// app-fields region and are constructed in New via the app-init region.
type App struct {
	Config *config.Config
	// scaffold:region:app-fields:start
	// scaffold:region:app-fields:end
}

func New(cfg *config.Config) (*App, error) {
	a := &App{Config: cfg}
	// scaffold:region:app-init:start
	// scaffold:region:app-init:end
	return a, nil
}
