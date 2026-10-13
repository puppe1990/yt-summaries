package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/puppe1990/amarra-cais/pkg/cais"

	appi18n "github.com/puppe1990/yt-summaries/internal/i18n"
	"github.com/puppe1990/yt-summaries/internal/models"
	"github.com/puppe1990/yt-summaries/internal/store"
)

func seedResumo(t *testing.T, s store.Store, videoID, slug, title, published string) {
	t.Helper()
	r := models.Resumo{
		YouTubeID:   videoID,
		Slug:        slug,
		Title:       title,
		Channel:     "Canal " + videoID,
		URL:         "https://youtu.be/" + videoID,
		PublishedAt: published,
		Tags:        []string{"tag-" + videoID},
		TLDR:        []string{"resumo de " + title},
	}
	payload, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	v := models.Video{
		YouTubeID:   videoID,
		Slug:        slug,
		Title:       title,
		PublishedAt: published,
		SearchText:  r.SearchText(),
		Payload:     string(payload),
	}
	if err := s.UpsertVideo(v); err != nil {
		t.Fatal(err)
	}
}

func newHome(t *testing.T, s store.Store) *HomeHandler {
	t.Helper()
	return NewHomeHandler(setupTestViews(t), s, testSite(), appi18n.DefaultCatalog(), cais.Config{})
}

func getBody(t *testing.T, h http.Handler, target string) (int, string) {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, target, nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	return rr.Code, rr.Body.String()
}

func TestHome_IndexPutsNextFirst(t *testing.T) {
	s := setupTestStore(t)
	seedResumo(t, s, "a1", "antigo", "Vídeo Antigo", "2026-10-01")
	seedResumo(t, s, "a2", "novo", "Vídeo Novo", "2026-10-12")

	code, body := getBody(t, newHome(t, s), "/")
	if code != http.StatusOK {
		t.Fatalf("status = %d, want 200", code)
	}
	iNovo := strings.Index(body, "/videos/novo")
	iAntigo := strings.Index(body, "/videos/antigo")
	if iNovo == -1 || iAntigo == -1 {
		t.Fatalf("cards ausentes no body (novo=%d, antigo=%d)", iNovo, iAntigo)
	}
	if iNovo > iAntigo {
		t.Errorf("o próximo (novo) deve vir primeiro: novo@%d antigo@%d", iNovo, iAntigo)
	}
	flag := strings.Index(body, `data-proximo="true"`)
	if flag == -1 || flag > iAntigo {
		t.Errorf("o destaque ▶ deveria estar no card do próximo (flag@%d, antigo@%d)", flag, iAntigo)
	}
}

func TestHome_SearchFilters(t *testing.T) {
	s := setupTestStore(t)
	seedResumo(t, s, "a1", "antigo", "Vídeo Antigo", "2026-10-01")
	seedResumo(t, s, "a2", "novo", "Vídeo Novo", "2026-10-12")

	code, body := getBody(t, newHome(t, s), "/?q=antigo")
	if code != http.StatusOK {
		t.Fatalf("status = %d, want 200", code)
	}
	if !strings.Contains(body, "/videos/antigo") {
		t.Errorf("busca por 'antigo' deveria manter o card antigo")
	}
	if strings.Contains(body, "/videos/novo") {
		t.Errorf("busca por 'antigo' não deveria mostrar o card novo")
	}
}

func TestHome_EmptyState(t *testing.T) {
	s := setupTestStore(t)
	code, body := getBody(t, newHome(t, s), "/")
	if code != http.StatusOK {
		t.Fatalf("status = %d, want 200", code)
	}
	if !strings.Contains(body, "Ainda não há resumos") {
		t.Errorf("esperado empty state, body começa com: %.200s", body)
	}
}

func TestHome_RedirectsLegacyVParam(t *testing.T) {
	s := setupTestStore(t)
	seedResumo(t, s, "a1", "antigo", "Vídeo Antigo", "2026-10-01")

	req := httptest.NewRequest(http.MethodGet, "/?v=antigo", nil)
	rr := httptest.NewRecorder()
	newHome(t, s).ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Fatalf("status = %d, want 303", rr.Code)
	}
	if loc := rr.Header().Get("Location"); loc != "/videos/antigo" {
		t.Errorf("Location = %q, want /videos/antigo", loc)
	}
}
