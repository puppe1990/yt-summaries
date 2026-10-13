package handlers

import (
	"net/http"
	"strings"

	"github.com/puppe1990/amarra-cais/pkg/amarra/view"
	"github.com/puppe1990/amarra-cais/pkg/cais"
	"github.com/puppe1990/amarra-cais/pkg/cais/i18n"
	"github.com/puppe1990/amarra-cais/pkg/cais/meta"

	"github.com/puppe1990/yt-resumos/internal/models"
	"github.com/puppe1990/yt-resumos/internal/store"
)

type HomeHandler struct {
	views   *view.Renderer
	store   store.Store
	site    meta.Site
	catalog *i18n.Catalog
	cfg     cais.Config
}

func NewHomeHandler(views *view.Renderer, s store.Store, site meta.Site, catalog *i18n.Catalog, cfg cais.Config) *HomeHandler {
	return &HomeHandler{views: views, store: s, site: site, catalog: catalog, cfg: cfg}
}

func (h *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	videos, err := h.store.ListVideos("")
	if err != nil {
		http.Error(w, "erro ao listar os resumos", http.StatusInternalServerError)
		return
	}

	q := strings.TrimSpace(r.URL.Query().Get("q"))
	filtered := filterVideos(videos, q)
	active, ok := pickActive(filtered, r.URL.Query().Get("v"))

	var resumo *resumoView
	if ok {
		rv, err := buildResumoView(active)
		if err != nil {
			http.Error(w, "erro ao carregar o resumo", http.StatusInternalServerError)
			return
		}
		resumo = &rv
	}

	view.Write(w, r, h.views, view.Page{
		Layout: "app",
		Name:   "home",
		Data: amarraData(r, h.site, map[string]any{
			"Q":      q,
			"Tabs":   tabsFrom(filtered, active.Slug),
			"Resumo": resumo,
			"Total":  len(videos),
		}),
	}, h.cfg)
}

func filterVideos(videos []models.Video, q string) []models.Video {
	if q == "" {
		return videos
	}
	q = strings.ToLower(q)
	out := make([]models.Video, 0, len(videos))
	for _, v := range videos {
		if strings.Contains(v.SearchText, q) {
			out = append(out, v)
		}
	}
	return out
}

// pickActive prefers ?v=<slug>; otherwise the first of the list — the "próximo"
// is always the first item on the page (project invariant).
func pickActive(videos []models.Video, want string) (models.Video, bool) {
	if len(videos) == 0 {
		return models.Video{}, false
	}
	for _, v := range videos {
		if v.Slug == want {
			return v, true
		}
	}
	return videos[0], true
}

func tabsFrom(videos []models.Video, activeSlug string) []tabView {
	tabs := make([]tabView, 0, len(videos))
	for i, v := range videos {
		tabs = append(tabs, tabView{
			Slug:    v.Slug,
			Title:   v.Title,
			Channel: v.Channel,
			Active:  v.Slug == activeSlug,
			Proximo: i == 0,
		})
	}
	return tabs
}
