package app

import (
	"io"
)

type PhotoStore interface {
	Get(photoID string, writer io.Writer) error
}
