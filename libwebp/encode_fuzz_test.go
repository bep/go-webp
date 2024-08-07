package libwebp

import (
	"bytes"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/bep/gowebp/libwebp/webpoptions"
)

func FuzzEncodePNG(f *testing.F) {
	names := []string{"bw-gopher.png", "fuzzy-cirlcle.png"}
	opts := webpoptions.EncodingOptions{Quality: 75}

	for _, name := range names {
		b, err := os.ReadFile(filepath.Join("..", "test_data", "images", name))
		if err != nil {
			f.Fatal(err)
		}
		f.Add(b)
	}

	f.Fuzz(func(t *testing.T, data []byte) {
		img, err := png.Decode(bytes.NewReader(data))
		if err != nil {
			if img != nil {
				t.Fatalf("img != nil, but err: %s", err)
			}
			return
		}
		err = Encode(io.Discard, img, opts)
		if err != nil {
			t.Fatal(err)
		}
	})
}

func FuzzEncodeJPG(f *testing.F) {
	names := []string{"source.jpg", "sunset.jpg"}
	opts := webpoptions.EncodingOptions{Quality: 75}

	for _, name := range names {
		b, err := os.ReadFile(filepath.Join("..", "test_data", "images", name))
		if err != nil {
			f.Fatal(err)
		}
		f.Add(b)
	}

	f.Fuzz(func(t *testing.T, data []byte) {
		img, err := jpeg.Decode(bytes.NewReader(data))
		if err != nil {
			if img != nil {
				t.Fatalf("img != nil, but err: %s", err)
			}
			return
		}
		err = Encode(io.Discard, img, opts)
		if err != nil {
			t.Fatal(err)
		}
	})
}
