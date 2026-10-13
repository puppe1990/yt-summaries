#!/usr/bin/env bash
# Rasteriza ícones PWA e o card Open Graph a partir dos vetores em
# web/static/icons/*.svg e design/og.html.
# Requer: librsvg (rsvg-convert) e playwright-cli.
set -euo pipefail

root="$(cd "$(dirname "$0")/.." && pwd)"
cd "$root"

mkdir -p web/static/icons

rsvg-convert -w 180 -h 180 web/static/icons/icon.svg -o web/static/icons/icon.png
rsvg-convert -w 192 -h 192 web/static/icons/icon.svg -o web/static/icons/icon-192.png
rsvg-convert -w 512 -h 512 web/static/icons/icon.svg -o web/static/icons/icon-512.png
rsvg-convert -w 512 -h 512 web/static/icons/icon-maskable.svg -o web/static/icons/icon-512-maskable.png

echo "ícones em web/static/icons/ (favicon.svg é o mark 32×32, não copiar o 512)"

# OG card: HTML + Fraunces via Playwright (rsvg não carrega Google Fonts).
# Screenshots do Playwright ficam em .playwright-cli (regra do workspace).
shot_dir="/Users/matheuspuppe/Desktop/Projetos/.playwright-cli"
mkdir -p "$shot_dir"

port=8793
python3 -m http.server "$port" --bind 127.0.0.1 --directory "$root" >/tmp/yt-resumos-og-http.log 2>&1 &
http_pid=$!
cleanup() {
  kill "$http_pid" 2>/dev/null || true
  playwright-cli close >/dev/null 2>&1 || true
}
trap cleanup EXIT
sleep 0.4

playwright-cli close >/dev/null 2>&1 || true
playwright-cli open "http://127.0.0.1:${port}/design/og.html"
playwright-cli resize 1200 630
playwright-cli run-code "async page => { await page.evaluate(() => document.fonts.ready); }"
(
  cd "$shot_dir"
  playwright-cli screenshot --filename=yt-resumos-og.png
)
cp "$shot_dir/yt-resumos-og.png" web/static/og.png
python3 - <<'PY'
from PIL import Image
im = Image.open("web/static/og.png").convert("RGB")
im.quantize(colors=64, method=Image.Quantize.MEDIANCUT).save("web/static/og.png", optimize=True)
PY

echo "og.png 1200x630 em web/static/og.png"
