package check_creation

import (
	"fmt"
	"log"
	"thumbnailr/app"
)

type Input struct {
	UserID string
	ThumbnailID string
}

type OutputData struct {
	IsCreated bool
}

type Handler struct {
	ThumbnailRepo app.ThumbnailRepo
}

func (h *Handler) logf(format string, args ...interface{}) {
	log.Printf("check creation: " + format,args)
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

	out = app.NewSuccessOutput()
	out.Data = OutputData{ IsCreated: thumbnail.IsCreated }
	return
}