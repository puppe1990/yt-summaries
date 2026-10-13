package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/puppe1990/amarra-cais/pkg/cais"
	"github.com/puppe1990/amarra-cais/pkg/cais/middleware"
)

func postResumo(t *testing.T, handler http.HandlerFunc, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/api/resumos", strings.NewReader(body))
	rr := httptest.NewRecorder()
	handler(rr, req)
	return rr
}

func TestAPI_CreateResumoUpserts(t *testing.T) {
	s := setupTestStore(t)
	api := NewAPIHandler(s, cais.Config{Env: "test"})

	body := `{"source_id":"reel-123","source":"@perfil","source_type":"instagram","slug":"reel-do-perfil","title":"Reel do perfil","tldr":["primeiro bullet"]}`
	rr := postResumo(t, api.CreateResumo, body)
	if rr.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201 (body: %s)", rr.Code, rr.Body.String())
	}

	got, err := s.FindVideoBySlug("reel-do-perfil")
	if err != nil {
		t.Fatal(err)
	}
	if got.YouTubeID != "reel-123" {
		t.Errorf("key = %q, want reel-123", got.YouTubeID)
	}
	if !strings.Contains(got.SearchText, "instagram") {
		t.Errorf("search_text deveria incluir o source_type: %q", got.SearchText)
	}
}

func TestAPI_LegacyPayloadFields(t *testing.T) {
	s := setupTestStore(t)
	api := NewAPIHandler(s, cais.Config{Env: "test"})

	body := `{"youtube_id":"vid-9","channel":"Canal X","slug":"video-legado","title":"Vídeo legado"}`
	rr := postResumo(t, api.CreateResumo, body)
	if rr.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201 (body: %s)", rr.Code, rr.Body.String())
	}
	got, err := s.FindVideoBySlug("video-legado")
	if err != nil {
		t.Fatal(err)
	}
	if got.YouTubeID != "vid-9" || got.Channel != "Canal X" {
		t.Errorf("fallback legado falhou: %+v", got)
	}
}

func TestAPI_Validation(t *testing.T) {
	s := setupTestStore(t)
	api := NewAPIHandler(s, cais.Config{Env: "test"})

	cases := []struct {
		name string
		body string
		want int
	}{
		{"not json", `{`, http.StatusBadRequest},
		{"missing key", `{"slug":"sem-fonte","title":"x"}`, http.StatusUnprocessableEntity},
		{"missing title", `{"source_id":"a1","slug":"sem-titulo"}`, http.StatusUnprocessableEntity},
		{"bad slug", `{"source_id":"a1","slug":"Slug Ruim","title":"x"}`, http.StatusUnprocessableEntity},
	}
	for _, tc := range cases {
		rr := postResumo(t, api.CreateResumo, tc.body)
		if rr.Code != tc.want {
			t.Errorf("%s: status = %d, want %d", tc.name, rr.Code, tc.want)
		}
	}
}

func TestAPI_AdminAuth(t *testing.T) {
	s := setupTestStore(t)
	api := NewAPIHandler(s, cais.Config{Env: "production", AdminToken: "s3cret"})
	guarded := middleware.Protect(cais.Config{Env: "production", AdminToken: "s3cret"}, api.CreateResumo)

	valid := `{"source_id":"a1","slug":"protegido","title":"Protegido"}`

	rr := postResumo(t, guarded, valid)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("sem token: status = %d, want 401", rr.Code)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/resumos", strings.NewReader(valid))
	req.Header.Set("Authorization", "Bearer wrong")
	rr = httptest.NewRecorder()
	guarded(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("token errado: status = %d, want 401", rr.Code)
	}

	req = httptest.NewRequest(http.MethodPost, "/api/resumos", strings.NewReader(valid))
	req.Header.Set("Authorization", "Bearer s3cret")
	rr = httptest.NewRecorder()
	guarded(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("token certo: status = %d, want 201 (body: %s)", rr.Code, rr.Body.String())
	}
}
