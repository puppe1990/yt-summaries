package db

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"sort"

	"github.com/puppe1990/yt-resumos/internal/models"
	"github.com/puppe1990/yt-resumos/internal/store"
)

// ImportVideos upserts every embedded data/videos/*.json into SQLite.
// Idempotent: safe to run on every boot.
func ImportVideos(s store.Store, fsys fs.FS) (int, error) {
	names, err := fs.Glob(fsys, "videos/*.json")
	if err != nil {
		return 0, fmt.Errorf("list resumo files: %w", err)
	}
	sort.Strings(names)
	for i, name := range names {
		if err := importVideoFile(s, fsys, name); err != nil {
			return i, err
		}
	}
	return len(names), nil
}

func importVideoFile(s store.Store, fsys fs.FS, name string) error {
	raw, err := fs.ReadFile(fsys, name)
	if err != nil {
		return fmt.Errorf("read %s: %w", name, err)
	}
	var r models.Resumo
	if err := json.Unmarshal(raw, &r); err != nil {
		return fmt.Errorf("parse %s: %w", name, err)
	}
	if r.YouTubeID == "" || r.Slug == "" || r.Title == "" {
		return fmt.Errorf("resumo %s: youtube_id, slug and title are required", name)
	}
	v := models.Video{
		YouTubeID:   r.YouTubeID,
		Slug:        r.Slug,
		Title:       r.Title,
		Channel:     r.Channel,
		URL:         r.URL,
		Thumbnail:   r.Thumbnail,
		PublishedAt: r.PublishedAt,
		Tags:        r.Tags,
		SearchText:  r.SearchText(),
		Payload:     string(raw),
	}
	return s.UpsertVideo(v)
}
