package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/puppe1990/amarra-cais/pkg/cais"

	appi18n "github.com/puppe1990/yt-summaries/internal/i18n"
	"github.com/puppe1990/yt-summaries/internal/store"
)

func newVideo(t *testing.T, s store.Store) *VideoHandler {
	t.Helper()
	return NewVideoHandler(setupTestViews(t), s, testSite(), appi18n.DefaultCatalog(), cais.Config{})
}

func getVideoBody(t *testing.T, h *VideoHandler, slug string) (int, string) {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/videos/"+slug, nil)
	req.SetPathValue("slug", slug)
	rr := httptest.NewRecorder()
	h.Route()(rr, req)
	return rr.Code, rr.Body.String()
}

func TestVideo_Found(t *testing.T) {
	s := setupTestStore(t)
	seedResumo(t, s, "a1", "antigo", "Vídeo Antigo", "2026-10-01")

	code, body := getVideoBody(t, newVideo(t, s), "antigo")
	if code != http.StatusOK {
		t.Fatalf("status = %d, want 200", code)
	}
	if !strings.Contains(body, "Vídeo Antigo") {
		t.Errorf("página deveria conter o título do vídeo")
	}
	if !strings.Contains(body, "resumo de Vídeo Antigo") {
		t.Errorf("página deveria conter o TLDR")
	}
}

func TestVideo_LegacyRedirect(t *testing.T) {
	s := setupTestStore(t)
	seedResumo(t, s, "a1", "antigo", "Vídeo Antigo", "2026-10-01")

	req := httptest.NewRequest(http.MethodGet, "/videos/antigo", nil)
	req.SetPathValue("slug", "antigo")
	rr := httptest.NewRecorder()
	newVideo(t, s).LegacyRedirect(rr, req)

	if rr.Code != http.StatusMovedPermanently {
		t.Fatalf("status = %d, want 301", rr.Code)
	}
	if loc := rr.Header().Get("Location"); loc != "/resumos/antigo" {
		t.Errorf("Location = %q, want /resumos/antigo", loc)
	}
}

func TestVideo_NotFound(t *testing.T) {
	s := setupTestStore(t)
	code, body := getVideoBody(t, newVideo(t, s), "nao-existe")
	if code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", code)
	}
	if !strings.Contains(body, "não está no cais") {
		t.Errorf("esperado a página 404 desenhada, body começa com: %.200s", body)
	}
}
