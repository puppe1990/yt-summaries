package app

import (
	"github.com/puppe1990/amarra-cais/pkg/amarra/live"
	"github.com/puppe1990/amarra-cais/pkg/cais"

	"github.com/puppe1990/yt-resumos/internal/handlers"
)

func registerRoutes(r *cais.Router, deps Deps, cfg cais.Config) {
	home := handlers.NewHomeHandler(deps.Views, deps.Store, deps.Site, deps.Catalog, cfg)
	video := handlers.NewVideoHandler(deps.Views, deps.Store, deps.Site, deps.Catalog, cfg)

	r.Get("/", home.ServeHTTP)
	r.Get("/videos/{slug}", video.Route())
	r.NotFound(video.NotFound)
}

func registerLiveViews(hub *live.Hub, deps Deps) {
	_ = hub
	_ = deps
	// cais:live-views
}
