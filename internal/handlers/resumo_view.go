package handlers

import (
	"github.com/puppe1990/yt-resumos/internal/models"
)

// resumoView is the payload shape for partials/resumo.html.
type resumoView struct {
	Slug        string
	Title       string
	Channel     string
	URL         string
	Thumbnail   string
	PublishedAt string
	Duration    string
	Tags        []string
	TLDR        []string
	Secoes      []models.Secao
	Numeros     []models.Numero
	Citacoes    []string
	Comentarios models.Comentarios
}

func buildResumoView(v models.Video) (resumoView, error) {
	r, err := models.ParseResumo(v.Payload)
	if err != nil {
		return resumoView{}, err
	}
	return resumoView{
		Slug:        v.Slug,
		Title:       r.Title,
		Channel:     r.Channel,
		URL:         r.URL,
		Thumbnail:   r.Thumbnail,
		PublishedAt: r.PublishedAt,
		Duration:    r.DurationLabel(),
		Tags:        r.Tags,
		TLDR:        r.TLDR,
		Secoes:      r.Secoes,
		Numeros:     r.Numeros,
		Citacoes:    r.Citacoes,
		Comentarios: r.Comentarios,
	}, nil
}

// tabView is one video link on the home tab bar. The first tab is always the
// "próximo" (newest) — project invariant.
type tabView struct {
	Slug    string
	Title   string
	Channel string
	Active  bool
	Proximo bool
}
