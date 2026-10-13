package models

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Video is one summarized YouTube video as stored in SQLite. Payload holds the
// full resumo JSON produced by the yt-resumos-zap skill.
type Video struct {
	ID          int64
	YouTubeID   string
	Slug        string
	Title       string
	Channel     string
	URL         string
	Thumbnail   string
	PublishedAt string
	ImportedAt  string
	Tags        []string
	SearchText  string
	Payload     string
}

// Resumo mirrors data/videos/<id>.json — the shape the app renders.
type Resumo struct {
	YouTubeID   string      `json:"youtube_id"`
	Slug        string      `json:"slug"`
	URL         string      `json:"url"`
	Title       string      `json:"title"`
	Channel     string      `json:"channel"`
	Thumbnail   string      `json:"thumbnail"`
	PublishedAt string      `json:"published_at"`
	DurationS   int         `json:"duration_s"`
	Tags        []string    `json:"tags"`
	TLDR        []string    `json:"tldr"`
	Secoes      []Secao     `json:"secoes"`
	Numeros     []Numero    `json:"numeros"`
	Citacoes    []string    `json:"citacoes"`
	Comentarios Comentarios `json:"comentarios"`
}

type Secao struct {
	T       string   `json:"t"`
	Titulo  string   `json:"titulo"`
	Bullets []string `json:"bullets"`
}

type Numero struct {
	Label string `json:"label"`
	Valor string `json:"valor"`
}

type Comentarios struct {
	ColetadoEm   string       `json:"coletado_em"`
	Threads      int          `json:"threads"`
	Esperado     int          `json:"esperado"`
	Distribuicao []Faixa      `json:"distribuicao"`
	Top          []Comentario `json:"top"`
	Temas        []Tema       `json:"temas"`
	Veredito     string       `json:"veredito"`
}

type Faixa struct {
	Faixa string `json:"faixa"`
	N     int    `json:"n"`
}

type Comentario struct {
	Autor string `json:"autor"`
	Likes int    `json:"likes"`
	Texto string `json:"texto"`
}

type Tema struct {
	Titulo   string       `json:"titulo"`
	Resumo   string       `json:"resumo"`
	Exemplos []Comentario `json:"exemplos"`
}

// ParseResumo decodes a stored payload.
func ParseResumo(payload string) (Resumo, error) {
	var r Resumo
	if err := json.Unmarshal([]byte(payload), &r); err != nil {
		return Resumo{}, fmt.Errorf("parse resumo: %w", err)
	}
	return r, nil
}

// SearchText builds the lower-case haystack used by the ?q= filter.
func (r Resumo) SearchText() string {
	parts := []string{r.Title, r.Channel, r.Slug, strings.Join(r.Tags, " "), strings.Join(r.TLDR, " ")}
	return strings.ToLower(strings.Join(parts, " "))
}

// DurationLabel renders seconds as "5:52"; empty when unknown.
func (r Resumo) DurationLabel() string {
	if r.DurationS <= 0 {
		return ""
	}
	return fmt.Sprintf("%d:%02d", r.DurationS/60, r.DurationS%60)
}
