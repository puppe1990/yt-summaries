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

func TestVideo_RendersPlanosTable(t *testing.T) {
	s := setupTestStore(t)
	resumo := models.Resumo{
		SourceID: "artigo-planos",
		Slug:     "tabela-planos",
		Title:    "Planos",
		TLDR:     []string{"tldr dos planos"},
		Planos: models.Planos{
			Titulo:   "O que tem em cada plano",
			Legenda:  "Pool semanal compartilhado.",
			Colunas:  []string{"Free", "Plus", "Premium"},
			Precos:   []string{"US$ 0", "US$ 8", "US$ 30"},
			Destaque: "Premium",
			Grupos: []models.PlanosGrupo{{
				Nome: "Grok",
				Linhas: []models.PlanosLinha{{
					Recurso: "Grok Bot",
					Celulas: []string{"—", "—", "sim"},
				}},
			}},
		},
	}
	payload, err := json.Marshal(resumo)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.UpsertVideo(models.Video{
		YouTubeID:  resumo.SourceID,
		Slug:       resumo.Slug,
		Title:      resumo.Title,
		SearchText: resumo.SearchText(),
		Payload:    string(payload),
	}); err != nil {
		t.Fatal(err)
	}

	code, body := getVideoBody(t, newVideo(t, s), "tabela-planos")
	if code != http.StatusOK {
		t.Fatalf("status = %d, want 200", code)
	}
	for _, want := range []string{"sr-planos", "O que tem em cada plano", "Grok Bot", "US$ 30", "Pool semanal compartilhado."} {
		if !strings.Contains(body, want) {
			t.Errorf("página deveria conter %q", want)
		}
	}
	if !strings.Contains(body, "is-hot") {
		t.Errorf("coluna em destaque deveria marcar is-hot")
	}
}

func TestVideo_HeadingAnchors(t *testing.T) {
	s := setupTestStore(t)
	resumo := models.Resumo{
		SourceID: "artigo-ancoras",
		Slug:     "ancoras",
		Title:    "Âncoras",
		TLDR:     []string{"tldr"},
		Planos: models.Planos{
			Colunas: []string{"Free"},
			Precos:  []string{"US$ 0"},
		},
		Secoes: []models.Secao{{
			N: "01", Tag: "Anúncio", Titulo: "Uma thread",
		}},
		Numeros:     []models.Numero{{Label: "n", Valor: "1"}},
		Citacoes:    []models.Citacao{{Texto: "fala", Fonte: "fonte"}},
		Aplicacao:   []models.Item{{Titulo: "casa", Texto: "faça"}},
		Comentarios: models.Comentarios{Threads: 3},
	}
	payload, err := json.Marshal(resumo)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.UpsertVideo(models.Video{
		YouTubeID:  resumo.SourceID,
		Slug:       resumo.Slug,
		Title:      resumo.Title,
		SearchText: resumo.SearchText(),
		Payload:    string(payload),
	}); err != nil {
		t.Fatal(err)
	}

	code, body := getVideoBody(t, newVideo(t, s), "ancoras")
	if code != http.StatusOK {
		t.Fatalf("status = %d, want 200", code)
	}
	for _, id := range []string{
		"em-30-segundos", "planos", "01-anuncio",
		"numeros", "citacoes", "levar-para-casa", "a-plateia",
	} {
		if !strings.Contains(body, `id="`+id+`"`) {
			t.Errorf("h2 deveria ter id=%q", id)
		}
		if !strings.Contains(body, `href="#`+id+`"`) {
			t.Errorf("h2 deveria ligar href=#%s", id)
		}
	}
}
