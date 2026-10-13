package handlers

import (
	"errors"
	"net/http"

	"github.com/puppe1990/amarra-cais/pkg/amarra/view"
	"github.com/puppe1990/amarra-cais/pkg/cais"
	"github.com/puppe1990/amarra-cais/pkg/cais/i18n"
	"github.com/puppe1990/amarra-cais/pkg/cais/meta"

	"github.com/puppe1990/yt-summaries/internal/models"
	"github.com/puppe1990/yt-summaries/internal/store"
)

type VideoHandler struct {
	views   *view.Renderer
	store   store.Store
	site    meta.Site
	catalog *i18n.Catalog
	cfg     cais.Config
}

func NewVideoHandler(views *view.Renderer, s store.Store, site meta.Site, catalog *i18n.Catalog, cfg cais.Config) *VideoHandler {
	return &VideoHandler{views: views, store: s, site: site, catalog: catalog, cfg: cfg}
}

// Route wraps the handler with the {slug} path param.
func (h *VideoHandler) Route() http.HandlerFunc {
	return cais.StringParam("slug", func(w http.ResponseWriter, r *http.Request, slug string) {
		h.serve(w, r, slug)
	})
}

func (h *VideoHandler) serve(w http.ResponseWriter, r *http.Request, slug string) {
	v, err := h.store.FindVideoBySlug(slug)
	if errors.Is(err, store.ErrVideoNotFound) {
		h.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, "erro ao carregar o vídeo", http.StatusInternalServerError)
		return
	}
	rv, err := models.ParseResumo(v.Payload)
	if err != nil {
		http.Error(w, "erro ao carregar o resumo", http.StatusInternalServerError)
		return
	}
	writeView(w, r, h.views, h.cfg, "app", "video", amarraData(r, h.site, map[string]any{
		"Title":  v.Title,
		"Resumo": rv,
	}), http.StatusOK)
}

// NotFound also serves unmatched routes (registered via r.NotFound).
func (h *VideoHandler) NotFound(w http.ResponseWriter, r *http.Request) {
	writeView(w, r, h.views, h.cfg, "app", "notfound", amarraData(r, h.site, map[string]any{
		"Title": "Não encontrado",
	}), http.StatusNotFound)
}
