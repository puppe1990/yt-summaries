package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"

	"github.com/puppe1990/amarra-cais/pkg/cais"

	"github.com/puppe1990/yt-summaries/internal/models"
	"github.com/puppe1990/yt-summaries/internal/store"
)

// slugPattern keeps slugs URL-safe: kebab-case letters, digits and dashes.
var slugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

// maxResumoBodyBytes caps the API body — a super resumo payload is a few KB.
const maxResumoBodyBytes = 1 << 20

type APIHandler struct {
	store store.Store
	cfg   cais.Config
}

func NewAPIHandler(s store.Store, cfg cais.Config) *APIHandler {
	return &APIHandler{store: s, cfg: cfg}
}

// CreateResumo upserts one super resumo sent by a machine client. The route is
// wrapped in middleware.Protect (Bearer ADMIN_TOKEN), and the CSRF middleware
// already steps aside for Bearer clients (#331).
func (h *APIHandler) CreateResumo(w http.ResponseWriter, r *http.Request) {
	var resumo models.Resumo
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxResumoBodyBytes))
	if err := dec.Decode(&resumo); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid JSON body: "+err.Error())
		return
	}
	if err := validateResumo(resumo); err != nil {
		writeJSONError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	payload, err := json.Marshal(resumo)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "encode resumo: "+err.Error())
		return
	}
	video := models.Video{
		YouTubeID:   resumo.Key(),
		Slug:        resumo.Slug,
		Title:       resumo.Title,
		Channel:     resumo.SourceName(),
		URL:         resumo.URL,
		Thumbnail:   resumo.Thumbnail,
		PublishedAt: resumo.PublishedAt,
		Tags:        resumo.Tags,
		SearchText:  resumo.SearchText(),
		Payload:     string(payload),
	}
	if err := h.store.UpsertVideo(video); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"slug":   resumo.Slug,
		"key":    resumo.Key(),
	})
}

func validateResumo(r models.Resumo) error {
	if r.Key() == "" {
		return errors.New("source_id (or legacy youtube_id) is required")
	}
	if r.Title == "" {
		return errors.New("title is required")
	}
	if !slugPattern.MatchString(r.Slug) {
		return errors.New("slug must be kebab-case (a-z, 0-9, single dashes)")
	}
	return nil
}

func writeJSONError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "error", "error": msg})
}
