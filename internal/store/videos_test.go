package store

import (
	"errors"
	"strings"
	"testing"

	"github.com/puppe1990/yt-summaries/internal/models"
)

func newVideosStore(t *testing.T) *SQLiteStore {
	t.Helper()
	s, err := NewSQLiteStore(":memory:", "test")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = s.Close() })
	return s
}

func fixtureVideo(videoID, slug, title, published string) models.Video {
	return models.Video{
		YouTubeID:   videoID,
		Slug:        slug,
		Title:       title,
		PublishedAt: published,
		Tags:        []string{"tag"},
		SearchText:  strings.ToLower(title + " tag"),
		Payload:     `{"youtube_id":"` + videoID + `"}`,
	}
}

func TestUpsertVideo_IsIdempotentAndUpdates(t *testing.T) {
	s := newVideosStore(t)
	v := fixtureVideo("v1", "primeiro", "Vídeo Primeiro", "2026-10-01")
	if err := s.UpsertVideo(v); err != nil {
		t.Fatal(err)
	}
	v.Title = "Vídeo Primeiro (editado)"
	if err := s.UpsertVideo(v); err != nil {
		t.Fatal(err)
	}

	got, err := s.FindVideoBySlug("primeiro")
	if err != nil {
		t.Fatal(err)
	}
	if got.Title != "Vídeo Primeiro (editado)" {
		t.Errorf("title = %q, want edited title", got.Title)
	}
}

func TestListVideos_MostRecentFirst(t *testing.T) {
	s := newVideosStore(t)
	if err := s.UpsertVideo(fixtureVideo("v1", "antigo", "Antigo", "2026-10-01")); err != nil {
		t.Fatal(err)
	}
	if err := s.UpsertVideo(fixtureVideo("v2", "novo", "Novo", "2026-10-12")); err != nil {
		t.Fatal(err)
	}

	list, err := s.ListVideos("")
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 2 {
		t.Fatalf("len = %d, want 2", len(list))
	}
	if list[0].Slug != "novo" {
		t.Errorf("first = %q, want novo (most recent)", list[0].Slug)
	}
}

func TestListVideos_FiltersByQuery(t *testing.T) {
	s := newVideosStore(t)
	if err := s.UpsertVideo(fixtureVideo("v1", "antigo", "Cavalo Antigo", "2026-10-01")); err != nil {
		t.Fatal(err)
	}
	if err := s.UpsertVideo(fixtureVideo("v2", "novo", "Drone Novo", "2026-10-12")); err != nil {
		t.Fatal(err)
	}

	list, err := s.ListVideos("cavalo")
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 || list[0].Slug != "antigo" {
		t.Errorf("filter = %+v, want only antigo", list)
	}
}

func TestFindVideoBySlug_NotFound(t *testing.T) {
	s := newVideosStore(t)
	_, err := s.FindVideoBySlug("nao-existe")
	if !errors.Is(err, ErrVideoNotFound) {
		t.Errorf("err = %v, want ErrVideoNotFound", err)
	}
}
