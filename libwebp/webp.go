package libwebp

import (
	"image"
	"io"

	"github.com/bep/go-webp/internal/libwebp"
)

func Encode(w io.Writer, src image.Image) error {
	// TODO1 opts
	options, err := libwebp.NewLossyEncoderOptions(libwebp.PresetDefault, float32(75))
	if err != nil {
		return err
	}
	if enc, err := libwebp.NewEncoder(src, options); err != nil {
		return err
	} else {
		return enc.Encode(w)
	}
}
