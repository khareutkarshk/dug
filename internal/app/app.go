package app

import (
	"github.com/khareutkarshk/dug/internal/config"
	"github.com/khareutkarshk/dug/internal/logger"
	"github.com/khareutkarshk/dug/internal/router"
	"github.com/khareutkarshk/dug/internal/server"
)

type App struct {
	Server  *server.Server
	Manager *router.Manager
}

func New(cfg *config.Config) (*App, error) {

	r, err := router.NewRouter(cfg)
	if err != nil {
		return nil, err
	}

	manager := router.NewManager(r)

	return &App{
		Server:  server.New(cfg, manager),
		Manager: manager,
	}, nil
}

func (a *App) EnableConfigReload(path string) error {

	return config.Watch(path, func() {

		cfg, err := config.Load(path)
		if err != nil {
			logger.Log.Error("reload failed", "error", err)
			return
		}

		newRouter, err := router.NewRouter(cfg)
		if err != nil {
			logger.Log.Error("router build failed", "error", err)
			return
		}

		a.Manager.Update(newRouter)

		logger.Log.Info("configuration reloaded")
	})
}
