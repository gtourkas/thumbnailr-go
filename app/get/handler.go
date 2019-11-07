package get

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"thumbnailr/app"
)

type Input struct {
	UserID string
	ThumbnailID string
}

type OutputData struct {
	Base64Contents string
	Format string
	Size string
}

type Handler struct {
	ThumbnailRepo app.ThumbnailRepo
	ThumbnailStore app.ThumbnailStore
}

func (h *Handler) logf(format string, args ...interface{}) {
	log.Printf("get: " + format,args)
}

func (h *Handler) Handle(in Input) (out app.Output) {

	thumbnail := app.Thumbnail{}
	if err := h.ThumbnailRepo.Get(in.ThumbnailID, &thumbnail); err != nil {
		out = app.NewUnexpectedErrorOutput()
		msg := fmt.Sprintf("cannot get thumbnail %s from repo", in.ThumbnailID)
		out.Message = msg
		h.logf("%s / error: %v", msg, err)
		return
	}

	var buff bytes.Buffer
	if err := h.ThumbnailStore.Get(in.UserID, in.ThumbnailID, &buff); err != nil {
		out = app.NewUnexpectedErrorOutput()
		msg := fmt.Sprintf("cannot get thumbnail %s from store", in.ThumbnailID)
		out.Message = msg
		h.logf("%s / error: %v", msg, err)
		return
	}

	out = app.NewSuccessOutput()
	out.Data = OutputData{
		Base64Contents: base64.StdEncoding.EncodeToString(buff.Bytes()),
		Format: string(thumbnail.Format),
		Size: fmt.Sprintf("%dx%d", thumbnail.Size.Width, thumbnail.Size.Length),
	}
	return
}
