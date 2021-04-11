package libwebp

import (
	"image"
	"io"

	"github.com/bep/gowebp/libwebp/webpoptions"

	"github.com/bep/gowebp/internal/libwebp"
)

// Encode encodes src as Webp into w using the options in o.
func Encode(w io.Writer, src image.Image, o webpoptions.EncodingOptions) error {
	return libwebp.Encode(w, src, o)
}
