package libwebp

import (
	"image/jpeg"
	"io/ioutil"
	"os"
	"testing"

	"github.com/bep/gowebp/libwebp/options"
)

func TestEncode(t *testing.T) {
	t.Run("encode lossy", func(t *testing.T) {
		r, err := os.Open("../test_data/images/source.jpg")
		if err != nil {
			t.Fatal(err)
		}

		img, err := jpeg.Decode(r)
		if err != nil {
			t.Fatal(err)
		}

		if err = Encode(ioutil.Discard, img, options.EncodingOptions{}); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("encode lossless", func(t *testing.T) {
		r, err := os.Open("../test_data/images/source.jpg")
		if err != nil {
			t.Fatal(err)
		}

		img, err := jpeg.Decode(r)
		if err != nil {
			t.Fatal(err)
		}

		if err = Encode(ioutil.Discard, img, options.EncodingOptions{Quality: 75}); err != nil {
			t.Fatal(err)
		}
	})
}

func BenchmarkEncode(b *testing.B) {
	r, err := os.Open("../test_data/images/source.jpg")
	if err != nil {
		b.Fatal(err)
	}

	img, err := jpeg.Decode(r)
	if err != nil {
		b.Fatal(err)
	}

	opts := options.EncodingOptions{Quality: 75}

	for i := 0; i < b.N; i++ {
		if err = Encode(ioutil.Discard, img, opts); err != nil {
			b.Fatal(err)
		}
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
