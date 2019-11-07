package app

import (
	"github.com/disintegration/imaging"
	"github.com/pkg/errors"
	"io"
)

type ThumbnailCreatorImpl struct {
}

var ErrUnsupportedFormat = errors.New("unsupported format")

func (tc *ThumbnailCreatorImpl) Create(photoReader io.Reader, format Format, size Size, thumbnailWriter io.Writer) error {

	var encodeFormat imaging.Format
	switch format {
	case JPEG :
		encodeFormat = imaging.JPEG
	case PNG:
		encodeFormat = imaging.PNG
	default:
		return ErrUnsupportedFormat
	}

	photo, err:= imaging.Decode(photoReader)
	if err != nil {
		return errors.Wrap(err, "cannot decode the photo image")
	}

	thumbnail := imaging.Thumbnail(photo, int(size.Width), int(size.Length), imaging.NearestNeighbor)

	if err := imaging.Encode(thumbnailWriter, thumbnail, encodeFormat); err != nil {
		return errors.Wrap(err, "cannot encode the thumbnail image")
	}

	return nil
}