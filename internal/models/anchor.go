package models

import (
	"strings"
	"unicode"
)

// Anchor turns a heading label into a URL fragment: "01 · Anúncio" → "01-anuncio".
func Anchor(label string) string {
	var b strings.Builder
	dash := true
	for _, r := range strings.ToLower(strings.TrimSpace(label)) {
		r = foldLatin(r)
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
			dash = false
			continue
		}
		if !dash {
			b.WriteByte('-')
			dash = true
		}
	}
	return strings.Trim(b.String(), "-")
}

func (s Secao) Anchor() string {
	label := strings.TrimSpace(s.N)
	tag := strings.TrimSpace(s.Tag)
	if label != "" && tag != "" {
		return Anchor(label + " " + tag)
	}
	if tag != "" {
		return Anchor(tag)
	}
	return Anchor(label)
}

func foldLatin(r rune) rune {
	switch r {
	case 'á', 'à', 'â', 'ã', 'ä':
		return 'a'
	case 'é', 'è', 'ê', 'ë':
		return 'e'
	case 'í', 'ì', 'î', 'ï':
		return 'i'
	case 'ó', 'ò', 'ô', 'õ', 'ö':
		return 'o'
	case 'ú', 'ù', 'û', 'ü':
		return 'u'
	case 'ç':
		return 'c'
	case 'ñ':
		return 'n'
	default:
		return r
	}
}
