package app

import (
	"github.com/puppe1990/amarra-cais/pkg/amarra/live"
	"github.com/puppe1990/amarra-cais/pkg/cais"
	"github.com/puppe1990/amarra-cais/pkg/cais/middleware"

	"github.com/puppe1990/yt-summaries/internal/handlers"
)

func registerRoutes(r *cais.Router, deps Deps, cfg cais.Config) {
	home := handlers.NewHomeHandler(deps.Views, deps.Store, deps.Site, deps.Catalog, cfg)
	video := handlers.NewVideoHandler(deps.Views, deps.Store, deps.Site, deps.Catalog, cfg)
	api := handlers.NewAPIHandler(deps.Store, cfg)

	r.Get("/", home.ServeHTTP)
	r.Get("/resumos/{slug}", video.Route())
	r.Get("/videos/{slug}", video.LegacyRedirect)
	r.Post("/api/resumos", middleware.Protect(cfg, api.CreateResumo))
	r.NotFound(video.NotFound)
}

func registerLiveViews(hub *live.Hub, deps Deps) {
	_ = hub
	_ = deps
	// cais:live-views
}
