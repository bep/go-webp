package libwebp

import (
	"image"
	"io"

	"github.com/bep/gowebp/libwebp/options"

	"github.com/bep/gowebp/internal/libwebp"
)

// Encode encodes src as Webp into w using the options in o.
func Encode(w io.Writer, src image.Image, o options.EncodingOptions) error {
	return libwebp.Encode(w, src, o)
}
