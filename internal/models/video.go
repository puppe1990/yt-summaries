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
	SourceID   string `json:"source_id,omitempty"`
	Source     string `json:"source,omitempty"`
	SourceType string `json:"source_type,omitempty"`
	// Legacy field names (pre-2026-10 payloads): youtube_id/channel.
	YouTubeID string `json:"youtube_id,omitempty"`
	Channel   string `json:"channel,omitempty"`

	Slug        string   `json:"slug"`
	URL         string   `json:"url"`
	Title       string   `json:"title"`
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
	Planos    Planos    `json:"planos,omitempty"`

	Comentarios Comentarios `json:"comentarios"`
}

// Planos is an optional comparison table (subscription tiers, feature matrix).
type Planos struct {
	Titulo   string        `json:"titulo"`
	Legenda  string        `json:"legenda"`
	Colunas  []string      `json:"colunas"`
	Precos   []string      `json:"precos"`
	Destaque string        `json:"destaque,omitempty"`
	Grupos   []PlanosGrupo `json:"grupos"`
}

type PlanosGrupo struct {
	Nome   string        `json:"nome"`
	Linhas []PlanosLinha `json:"linhas"`
}

type PlanosLinha struct {
	Recurso string   `json:"recurso"`
	Celulas []string `json:"celulas"`
}

// ColSpan is feature column plus one cell per plan.
func (p Planos) ColSpan() int {
	return len(p.Colunas) + 1
}

// HotCol is true when column i is the highlighted plan.
func (p Planos) HotCol(i int) bool {
	if i < 0 || i >= len(p.Colunas) {
		return false
	}
	return p.Destaque != "" && p.Colunas[i] == p.Destaque
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

// Key returns the unique source id: source_id, falling back to the legacy
// youtube_id so old payloads keep working.
func (r Resumo) Key() string {
	if r.SourceID != "" {
		return r.SourceID
	}
	return r.YouTubeID
}

// SourceName returns the source label (channel, profile, author, site):
// source, falling back to the legacy channel field.
func (r Resumo) SourceName() string {
	if r.Source != "" {
		return r.Source
	}
	return r.Channel
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
	parts := []string{r.Title, r.SourceName(), r.SourceType, r.Slug, r.H1A, r.H1Em, r.H1B, r.Lede, strings.Join(r.Tags, " "), strings.Join(r.TLDR, " ")}
	return strings.ToLower(strings.Join(parts, " "))
}

// DurationLabel renders seconds as "5:52"; empty when unknown.
func (r Resumo) DurationLabel() string {
	if r.DurationS <= 0 {
		return ""
	}
	return fmt.Sprintf("%d:%02d", r.DurationS/60, r.DurationS%60)
}
