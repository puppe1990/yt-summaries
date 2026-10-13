package models

import "testing"

func TestAnchor_KebabFromPortugueseLabel(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"Planos", "planos"},
		{"Em 30 segundos", "em-30-segundos"},
		{"01 · Anúncio", "01-anuncio"},
		{"Citações", "citacoes"},
		{"A plateia", "a-plateia"},
		{"Levar para casa", "levar-para-casa"},
		{"  Foo   Bar  ", "foo-bar"},
		{"", ""},
		{"---", ""},
	}
	for _, tc := range cases {
		if got := Anchor(tc.in); got != tc.want {
			t.Errorf("Anchor(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestSecaoAnchor_UsesNumberAndTag(t *testing.T) {
	s := Secao{N: "01", Tag: "Anúncio", T: "0:00", Titulo: "Uma thread longa"}
	if got := s.Anchor(); got != "01-anuncio" {
		t.Errorf("Secao.Anchor() = %q, want %q", got, "01-anuncio")
	}
}
