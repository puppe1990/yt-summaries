package data

import "embed"

// Videos are the resumo payloads imported into SQLite at boot
// (payloads written by the yt-resumos-zap skill).
//
//go:embed videos/*.json
var Videos embed.FS
