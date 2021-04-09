package libwebp

/*
#include <stdlib.h>
#include <encode.h>

static uint8_t* encodeNRGBA(WebPConfig* config, const uint8_t* rgba, int width, int height, int stride, size_t* output_size) {
	WebPPicture pic;
	WebPMemoryWriter wrt;
	int ok;
	if (!WebPPictureInit(&pic)) {
		return NULL;
	}
	pic.use_argb = 1;
	pic.width = width;
	pic.height = height;
	pic.writer = WebPMemoryWrite;
	pic.custom_ptr = &wrt;
	WebPMemoryWriterInit(&wrt);
	ok = WebPPictureImportRGBA(&pic, rgba, stride) && WebPEncode(config, &pic);
	WebPPictureFree(&pic);
	if (!ok) {
		WebPMemoryWriterClear(&wrt);
		return NULL;
	}
	*output_size = wrt.size;
	return wrt.mem;
}

*/
import "C"
import (
	"errors"
	"image"
	"image/draw"
	"io"
	"unsafe"

	"github.com/bep/gowebp/libwebp/options"
)

type (
	Encoder struct {
		config *C.WebPConfig
		img    *image.NRGBA
	}
)

// Encode encodes src into w considering the options in o.
//
// TODO(bep) ColorSpace
// TODO(bep) Can we handle *image.YCbCr without conversion?
// TODO(bep) Grayscale
func Encode(w io.Writer, src image.Image, o options.EncodingOptions) error {
	config, err := encodingOptionsToCConfig(o)
	if err != nil {
		return err
	}

	var (
		bounds = src.Bounds()
		stride int
		rgba   *C.uint8_t
	)

	switch v := src.(type) {
	case *image.RGBA:
		rgba = (*C.uint8_t)(&v.Pix[0])
		stride = v.Stride
	case *image.NRGBA:
		rgba = (*C.uint8_t)(&v.Pix[0])
		stride = v.Stride
	default:
		img := ConvertToNRGBA(src)
		rgba = (*C.uint8_t)(&img.Pix[0])
		stride = img.Stride
	}

	var size C.size_t
	output := C.encodeNRGBA(
		config,
		rgba,
		C.int(bounds.Max.X),
		C.int(bounds.Max.Y),
		C.int(stride),
		&size,
	)

	if output == nil || size == 0 {
		return errors.New("cannot encode webppicture")
	}
	defer C.free(unsafe.Pointer(output))

	_, err = w.Write(((*[1 << 30]byte)(unsafe.Pointer(output)))[0:int(size):int(size)])

	return err
}

func ConvertToNRGBA(src image.Image) *image.NRGBA {
	dst := image.NewNRGBA(src.Bounds())
	draw.Draw(dst, dst.Bounds(), src, src.Bounds().Min, draw.Src)

	return dst
}

func encodingOptionsToCConfig(o options.EncodingOptions) (*C.WebPConfig, error) {
	cfg := &C.WebPConfig{}
	quality := C.float(o.Quality)

	if C.WebPConfigPreset(cfg, C.WebPPreset(o.EncodingPreset), quality) == 0 {
		return nil, errors.New("failed to init encoder config")
	}

	cfg.quality = quality
	if quality == 0 {
		cfg.lossless = C.int(1)
	}

	if C.WebPValidateConfig(cfg) == 0 {
		return nil, errors.New("failed to validate config")
	}

	return cfg, nil

}
