package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/puppe1990/amarra-cais/pkg/cais"

	appi18n "github.com/puppe1990/yt-resumos/internal/i18n"
	"github.com/puppe1990/yt-resumos/internal/models"
	"github.com/puppe1990/yt-resumos/internal/store"
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

func TestHome_PutsNextVideoFirst(t *testing.T) {
	s := setupTestStore(t)
	seedResumo(t, s, "a1", "antigo", "Vídeo Antigo", "2026-10-01")
	seedResumo(t, s, "a2", "novo", "Vídeo Novo", "2026-10-12")

	code, body := getBody(t, newHome(t, s), "/")
	if code != http.StatusOK {
		t.Fatalf("status = %d, want 200", code)
	}
	iNovo := strings.Index(body, "v=novo")
	iAntigo := strings.Index(body, "v=antigo")
	if iNovo == -1 || iAntigo == -1 {
		t.Fatalf("abas ausentes no body (novo=%d, antigo=%d)", iNovo, iAntigo)
	}
	if iNovo > iAntigo {
		t.Errorf("o vídeo mais novo (próximo) deve vir primeiro: novo@%d antigo@%d", iNovo, iAntigo)
	}
	if !strings.Contains(body, "▶") {
		t.Errorf("faltou o marcador ▶ de próximo")
	}
	if !strings.Contains(body, "resumo de Vídeo Novo") {
		t.Errorf("a aba ativa deveria ser o próximo (Vídeo Novo)")
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
	if !strings.Contains(body, "v=antigo") {
		t.Errorf("busca por 'antigo' deveria manter a aba antigo")
	}
	if strings.Contains(body, "v=novo") {
		t.Errorf("busca por 'antigo' não deveria mostrar a aba novo")
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
