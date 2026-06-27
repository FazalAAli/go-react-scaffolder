package app

import "backend/internal/config"

// App is the dependency container. Future features attach as fields here
// (e.g. DB *gorm.DB) and are constructed in New.
type App struct {
	Config *config.Config
}

func New(cfg *config.Config) (*App, error) {
	return &App{Config: cfg}, nil
}
