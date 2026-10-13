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

// Resumo mirrors data/videos/<id>.json (schema v2 — super resumo no estilo
// dhh-rails-world-2026: kicker, h1 editorial, placas de tese, seções numeradas,
// citações com tradução e leitura da audiência).
type Resumo struct {
	YouTubeID   string   `json:"youtube_id"`
	Slug        string   `json:"slug"`
	URL         string   `json:"url"`
	Title       string   `json:"title"`
	Channel     string   `json:"channel"`
	Thumbnail   string   `json:"thumbnail"`
	PublishedAt string   `json:"published_at"`
	DurationS   int      `json:"duration_s"`
	Tags        []string `json:"tags"`

	Kicker string `json:"kicker"`
	H1A    string `json:"h1_a"`
	H1Em   string `json:"h1_em"`
	H1B    string `json:"h1_b"`
	Lede   string `json:"lede"`
	Tese   []Tese `json:"tese"`

	TLDR      []string  `json:"tldr"`
	Secoes    []Secao   `json:"secoes"`
	Numeros   []Numero  `json:"numeros"`
	Citacoes  []Citacao `json:"citacoes"`
	Aplicacao []Item    `json:"aplicacao"`

	Comentarios Comentarios `json:"comentarios"`
}

// Tese é uma placa ("PLACA 01") com lead em negrito e o parágrafo.
type Tese struct {
	Placa string `json:"placa"`
	Lead  string `json:"lead"`
	Texto string `json:"texto"`
}

// Secao é um bloco numerado: "01 · Gancho · 0:00" + título.
type Secao struct {
	N       string   `json:"n"`
	Tag     string   `json:"tag"`
	T       string   `json:"t"`
	Titulo  string   `json:"titulo"`
	Bullets []string `json:"bullets"`
}

type Numero struct {
	Label string `json:"label"`
	Valor string `json:"valor"`
}

type Citacao struct {
	Texto    string `json:"texto"`
	Fonte    string `json:"fonte"`
	Traducao string `json:"traducao"`
}

// Item é um bloco título+texto (aplicação prática e insights).
type Item struct {
	Titulo string `json:"titulo"`
	Texto  string `json:"texto"`
}

type Comentarios struct {
	ColetadoEm        string       `json:"coletado_em"`
	Threads           int          `json:"threads"`
	Esperado          int          `json:"esperado"`
	Distribuicao      []Faixa      `json:"distribuicao"`
	Top               []Comentario `json:"top"`
	Temas             []Tema       `json:"temas"`
	InsightsPositivos []Item       `json:"insights_positivos"`
	InsightsNegativos []Item       `json:"insights_negativos"`
	Consenso          []string     `json:"consenso"`
	Briga             []string     `json:"briga"`
	Veredito          string       `json:"veredito"`
}

type Faixa struct {
	Faixa string `json:"faixa"`
	N     int    `json:"n"`
}

type Comentario struct {
	Autor    string `json:"autor"`
	Likes    int    `json:"likes"`
	Texto    string `json:"texto"`
	Traducao string `json:"traducao"`
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
	parts := []string{r.Title, r.Channel, r.Slug, r.H1A, r.H1Em, r.H1B, r.Lede, strings.Join(r.Tags, " "), strings.Join(r.TLDR, " ")}
	return strings.ToLower(strings.Join(parts, " "))
}

// DurationLabel renders seconds as "5:52"; empty when unknown.
func (r Resumo) DurationLabel() string {
	if r.DurationS <= 0 {
		return ""
	}
	return fmt.Sprintf("%d:%02d", r.DurationS/60, r.DurationS%60)
}
