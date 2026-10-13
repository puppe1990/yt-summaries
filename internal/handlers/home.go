package handlers

import (
	"net/http"
	"strings"

	"github.com/puppe1990/amarra-cais/pkg/amarra/view"
	"github.com/puppe1990/amarra-cais/pkg/cais"
	"github.com/puppe1990/amarra-cais/pkg/cais/i18n"
	"github.com/puppe1990/amarra-cais/pkg/cais/meta"

	"github.com/puppe1990/yt-summaries/internal/models"
	"github.com/puppe1990/yt-summaries/internal/store"
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
	if slug := strings.TrimSpace(r.URL.Query().Get("v")); slug != "" {
		http.Redirect(w, r, "/resumos/"+slug, http.StatusSeeOther)
		return
	}

	videos, err := h.store.ListVideos("")
	if err != nil {
		http.Error(w, "erro ao listar os resumos", http.StatusInternalServerError)
		return
	}

	// The "próximo" is the first of the unfiltered list — it keeps the ▶ badge
	// even while searching (project invariant).
	proximoSlug := ""
	if len(videos) > 0 {
		proximoSlug = videos[0].Slug
	}

	q := strings.TrimSpace(r.URL.Query().Get("q"))
	cards := indexCards(filterVideos(videos, q), proximoSlug)

	view.Write(w, r, h.views, view.Page{
		Layout: "app",
		Name:   "home",
		Data: amarraData(r, h.site, map[string]any{
			"Q":     q,
			"Cards": cards,
			"Total": len(videos),
		}),
	}, h.cfg)
}

func indexCards(videos []models.Video, proximoSlug string) []indexCard {
	cards := make([]indexCard, 0, len(videos))
	for _, v := range videos {
		card := indexCard{
			Slug:    v.Slug,
			Channel: v.Channel,
			Date:    v.PublishedAt,
			Proximo: v.Slug == proximoSlug,
		}
		if r, err := models.ParseResumo(v.Payload); err == nil {
			card.Kicker = r.Kicker
			card.TitleA = r.H1A
			card.TitleEm = r.H1Em
			card.Lede = r.Lede
			card.Duration = r.DurationLabel()
			card.Threads = r.Comentarios.Threads
			card.Tags = r.Tags
			if card.Channel == "" {
				card.Channel = r.Channel
			}
		}
		if card.TitleA == "" {
			card.TitleA = v.Title
		}
		cards = append(cards, card)
	}
	return cards
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
