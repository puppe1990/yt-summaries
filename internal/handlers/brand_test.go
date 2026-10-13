package handlers

import (
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFaviconIsDarkroomMark(t *testing.T) {
	body, err := os.ReadFile(filepath.Join(projectRoot(t), "web", "static", "favicon.svg"))
	if err != nil {
		t.Fatal(err)
	}
	svg := string(body)
	if strings.Contains(svg, "#0b2f3a") {
		t.Fatal("favicon ainda é o placeholder do Cais (harbor #0b2f3a)")
	}
	if !strings.Contains(svg, `viewBox="0 0 32 32"`) {
		t.Error("favicon deve ser o mark compacto 32×32, não o ícone 512")
	}
	for _, color := range []string{"#120e0b", "#efe4cf", "#5fb3a8"} {
		if !strings.Contains(svg, color) {
			t.Errorf("favicon.svg precisa da cor da paleta %s", color)
		}
	}
}

func TestOpenGraphImageExists(t *testing.T) {
	assertPNGSize(t, filepath.Join(projectRoot(t), "web", "static", "og.png"), 1200, 630)
}

func TestBrandIconsHaveExpectedSizes(t *testing.T) {
	root := filepath.Join(projectRoot(t), "web", "static", "icons")
	assertPNGSize(t, filepath.Join(root, "icon.png"), 180, 180)
	assertPNGSize(t, filepath.Join(root, "icon-192.png"), 192, 192)
	assertPNGSize(t, filepath.Join(root, "icon-512.png"), 512, 512)
	assertPNGSize(t, filepath.Join(root, "icon-512-maskable.png"), 512, 512)
}

func TestBrandRastersUseDarkroomPalette(t *testing.T) {
	root := projectRoot(t)
	for _, rel := range []string{
		"web/static/og.png",
		"web/static/icons/icon-512.png",
		"web/static/icons/icon-512-maskable.png",
	} {
		path := filepath.Join(root, filepath.FromSlash(rel))
		if !pngHasApprox(t, path, 0x12, 0x0e, 0x0b, 12) {
			t.Errorf("%s não tem o darkroom #120e0b", rel)
		}
		if !pngHasApprox(t, path, 0xef, 0xe4, 0xcf, 18) {
			t.Errorf("%s não tem o paper #efe4cf", rel)
		}
	}
}

func TestHome_carriesAbsoluteOpenGraphTags(t *testing.T) {
	s := setupTestStore(t)
	code, body := getBody(t, newHome(t, s), "/")
	if code != 200 {
		t.Fatalf("status = %d, want 200", code)
	}
	wants := []string{
		`<meta property="og:site_name" content="YT Resumos" />`,
		`<meta property="og:description" content="`,
		`<meta property="og:image" content="https://cais.example.com/static/og.png" />`,
		`<meta property="og:image:width" content="1200" />`,
		`<meta property="og:image:height" content="630" />`,
		`<meta name="twitter:card" content="summary_large_image" />`,
		`<meta name="twitter:image" content="https://cais.example.com/static/og.png" />`,
		`href="/static/favicon.svg"`,
		`src="/static/favicon.svg"`,
	}
	for _, want := range wants {
		if !strings.Contains(body, want) {
			t.Errorf("missing %s", want)
		}
	}
}

func assertPNGSize(t *testing.T, path string, w, h int) {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("%s: %v", path, err)
	}
	defer func() { _ = f.Close() }()
	cfg, err := png.DecodeConfig(f)
	if err != nil {
		t.Fatalf("%s: decode: %v", path, err)
	}
	if cfg.Width != w || cfg.Height != h {
		t.Errorf("%s = %dx%d, want %dx%d", filepath.Base(path), cfg.Width, cfg.Height, w, h)
	}
}

func pngHasApprox(t *testing.T, path string, r, g, b, slop uint8) bool {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = f.Close() }()
	img, err := png.Decode(f)
	if err != nil {
		t.Fatal(err)
	}
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y += 3 {
		for x := bounds.Min.X; x < bounds.Max.X; x += 3 {
			pr, pg, pb, pa := img.At(x, y).RGBA()
			if pa>>8 < 200 {
				continue
			}
			if near(uint8(pr>>8), r, slop) && near(uint8(pg>>8), g, slop) && near(uint8(pb>>8), b, slop) {
				return true
			}
		}
	}
	return false
}

func near(got, want, slop uint8) bool {
	d := int(got) - int(want)
	if d < 0 {
		d = -d
	}
	return d <= int(slop)
}
