package store

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/puppe1990/yt-resumos/internal/models"
)

var ErrVideoNotFound = errors.New("video not found")

const videoColumns = `id, youtube_id, slug, title, channel, url, thumbnail,
	COALESCE(published_at, ''), imported_at, tags, search_text, payload`

// UpsertVideo imports or refreshes one resumo; keyed by youtube_id so the
// boot import is idempotent.
func (s *SQLiteStore) UpsertVideo(v models.Video) error {
	_, err := s.db.Raw().Exec(`
		INSERT INTO videos (youtube_id, slug, title, channel, url, thumbnail, published_at, tags, search_text, payload)
		VALUES (?, ?, ?, ?, ?, ?, NULLIF(?, ''), ?, ?, ?)
		ON CONFLICT(youtube_id) DO UPDATE SET
			slug = excluded.slug,
			title = excluded.title,
			channel = excluded.channel,
			url = excluded.url,
			thumbnail = excluded.thumbnail,
			published_at = excluded.published_at,
			tags = excluded.tags,
			search_text = excluded.search_text,
			payload = excluded.payload`,
		v.YouTubeID, v.Slug, v.Title, v.Channel, v.URL, v.Thumbnail, v.PublishedAt,
		strings.Join(v.Tags, ","), v.SearchText, v.Payload)
	if err != nil {
		return fmt.Errorf("upsert video %s: %w", v.YouTubeID, err)
	}
	return nil
}

// ListVideos returns every video, most recent first — the "próximo" (newest)
// is always index 0 (project invariant). q filters search_text, case-insensitive.
func (s *SQLiteStore) ListVideos(q string) ([]models.Video, error) {
	query := `SELECT ` + videoColumns + ` FROM videos`
	args := []any{}
	if q = strings.TrimSpace(strings.ToLower(q)); q != "" {
		query += ` WHERE search_text LIKE ?`
		args = append(args, "%"+q+"%")
	}
	query += ` ORDER BY COALESCE(NULLIF(published_at, ''), imported_at) DESC, id DESC`

	rows, err := s.db.Raw().Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("list videos: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var videos []models.Video
	for rows.Next() {
		v, err := scanVideo(rows)
		if err != nil {
			return nil, fmt.Errorf("scan video: %w", err)
		}
		videos = append(videos, v)
	}
	return videos, rows.Err()
}

// FindVideoBySlug returns one video or ErrVideoNotFound.
func (s *SQLiteStore) FindVideoBySlug(slug string) (models.Video, error) {
	row := s.db.Raw().QueryRow(`SELECT `+videoColumns+` FROM videos WHERE slug = ?`, slug)
	v, err := scanVideo(row)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Video{}, ErrVideoNotFound
	}
	if err != nil {
		return models.Video{}, fmt.Errorf("find video %s: %w", slug, err)
	}
	return v, nil
}

type videoScanner interface{ Scan(dest ...any) error }

func scanVideo(row videoScanner) (models.Video, error) {
	var v models.Video
	var tags string
	err := row.Scan(&v.ID, &v.YouTubeID, &v.Slug, &v.Title, &v.Channel, &v.URL,
		&v.Thumbnail, &v.PublishedAt, &v.ImportedAt, &tags, &v.SearchText, &v.Payload)
	if err != nil {
		return models.Video{}, err
	}
	if tags != "" {
		v.Tags = strings.Split(tags, ",")
	}
	return v, nil
}
