package app

import (
	"io"
	"time"
)

type Size struct {
	Width  uint
	Length uint
}

type Format string

const (
	JPEG Format = "JPEG"
	PNG  Format = "PNG"
)

type ThumbnailStore interface {
	Put(userID string, id string, reader io.Reader) error
	Get(userID string, id string, writer io.Writer) error
}

type ThumbnailCreator interface {
	Create(src io.Reader, format Format, size Size, dest io.Writer) (err error)
}

type Thumbnail struct {
	ID          string
	UserID      string
	PhotoID     string
	Size        Size
	Format      Format
	RequestedAt time.Time
	IsCreated   bool
	CreatedAt   *time.Time
}

type ThumbnailRepo interface {
	Get(id string) (*Thumbnail, error)
	Save(thumbnail *Thumbnail) error
}
