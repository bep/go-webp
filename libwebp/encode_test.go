package libwebp

import (
	"bytes"
	"flag"
	"image"
	"image/draw"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/bep/gowebp/internal/libwebp"

	"github.com/bep/gowebp/libwebp/webpoptions"
	"golang.org/x/image/webp"
)

type testCase struct {
	name      string
	inputFile string
	opts      webpoptions.EncodingOptions
}

var testCases = []testCase{
	{"lossy", "sunset.jpg", webpoptions.EncodingOptions{Quality: 75, EncodingPreset: webpoptions.EncodingPresetPhoto, UseSharpYuv: true}},
	{"lossless", "source.jpg", webpoptions.EncodingOptions{}},
	{"bw", "bw-gopher.png", webpoptions.EncodingOptions{Quality: 75}},
}

func TestEncode(t *testing.T) {
	writeGolden := true
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			r, err := os.Open(filepath.Join("../test_data/images", test.inputFile))
			if err != nil {
				t.Fatal(err)
			}

			img, _, err := image.Decode(r)
			if err != nil {
				t.Fatal(err)
			}

			targetName := strings.TrimSuffix(test.inputFile, filepath.Ext(test.inputFile)) + "-" + test.name + ".webp"

			targetFilename := filepath.Join("../test_data/images/golden", targetName)
			b := &bytes.Buffer{}
			var w io.Writer = b

			if writeGolden {
				f, err := os.Create(targetFilename)
				if err != nil {
					t.Fatal(err)
				}
				w = f
				defer f.Close()
			}

			if err = Encode(w, img, test.opts); err != nil {
				t.Fatal(err)
			}

			if !writeGolden {
				f, err := os.Open(targetFilename)
				if err != nil {
					t.Fatal(err)
				}
				defer f.Close()

				img1 := decodeWebp(t, f)
				img2 := decodeWebp(t, b)

				if !goldenEqual(img1, img2) {
					t.Fatal("images are different")
				}

			}
		})
	}
}

var longrunning = flag.Bool("longrunning", false, "Enable long running tests.")

// go test -v ./libwebp --longrunning
func TestEncodeLongRunning(t *testing.T) {
	if !(*longrunning) {
		t.Skip("Skip long running test...")
	}
	t.Log("Start...")
	for i := 0; i < 60; i++ {
		for _, test := range testCases {
			r, err := os.Open(filepath.Join("../test_data/images", test.inputFile))
			if err != nil {
				t.Fatal(err)
			}

			img, _, err := image.Decode(r)
			if err != nil {
				t.Fatal(err)
			}

			if err = Encode(ioutil.Discard, img, test.opts); err != nil {
				t.Fatal(err)
			}
		}

		time.Sleep(2 * time.Second)
	}

	t.Log("Done...")
}

func BenchmarkEncode(b *testing.B) {
	for _, test := range testCases {
		b.Run(test.name, func(b *testing.B) {
			r, err := os.Open(filepath.Join("../test_data/images", test.inputFile))
			if err != nil {
				b.Fatal(err)
			}

			img, _, err := image.Decode(r)
			if err != nil {
				b.Fatal(err)
			}

			// Encode will convert to NRGBA if needed. Do that here so we
			// don't get those numbers included in the below.
			imgrgba := libwebp.ConvertToNRGBA(img)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if err = Encode(ioutil.Discard, imgrgba, test.opts); err != nil {
					b.Fatal(err)
				}
			}

		})
	}

}

// Just to have something to compare with.
func BenchmarkEncodeJpeg(b *testing.B) {
	r, err := os.Open("../test_data/images/source.jpg")
	if err != nil {
		b.Fatal(err)
	}

	img, err := jpeg.Decode(r)
	if err != nil {
		b.Fatal(err)
	}

	opts := &jpeg.Options{
		Quality: 75,
	}

	for i := 0; i < b.N; i++ {
		if err = jpeg.Encode(ioutil.Discard, img, opts); err != nil {
			b.Fatal(err)
		}
	}
}

// usesFMA indicates whether "fused multiply and add" (FMA) instruction is
// used.  The command "grep FMADD go/test/codegen/floats.go" can help keep
// the FMA-using architecture list updated.
var usesFMA = runtime.GOARCH == "s390x" ||
	runtime.GOARCH == "ppc64" ||
	runtime.GOARCH == "ppc64le" ||
	runtime.GOARCH == "arm64"

func decodeWebp(t *testing.T, r io.Reader) *image.NRGBA {
	img, err := webp.Decode(r)
	if err != nil {
		t.Fatal(err)
	}

	b := img.Bounds()
	m := image.NewNRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(m, m.Bounds(), img, b.Min, draw.Src)

	return m
}

// goldenEqual compares two NRGBA images.
// A small tolerance is allowed on architectures using "fused multiply and add"
// (FMA) instruction to accommodate for floating-point rounding differences
// with control golden images that were generated on amd64 architecture.
// See https://golang.org/ref/spec#Floating_point_operators
// and https://github.com/gohugoio/hugo/issues/6387 for more information.
//
// Borrowed from https://github.com/disintegration/gift/blob/a999ff8d5226e5ab14b64a94fca07c4ac3f357cf/gift_test.go#L598-L625
// Copyright (c) 2014-2019 Grigory Dryapak
// Licensed under the MIT License.
func goldenEqual(img1, img2 *image.NRGBA) bool {
	maxDiff := 0
	if usesFMA {
		maxDiff = 1
	}
	if !img1.Rect.Eq(img2.Rect) {
		return false
	}
	if len(img1.Pix) != len(img2.Pix) {
		return false
	}
	for i := 0; i < len(img1.Pix); i++ {
		diff := int(img1.Pix[i]) - int(img2.Pix[i])
		if diff < 0 {
			diff = -diff
		}
		if diff > maxDiff {
			return false
		}
	}
	return true
}
