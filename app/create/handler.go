package create

import (
	"bytes"
	"fmt"
	"log"
	"thumbnailr/app"
	"time"
)

type Input struct {
	UserID string
	ThumbnailID string
	PhotoID string
	Format app.Format
	Size app.Size
}

type Handler struct {
	PhotoStore       app.PhotoStore
	ThumbnailStore app.ThumbnailStore
	ThumbnailRepo app.ThumbnailRepo

	thumbnailCreator app.ThumbnailCreator
}

func NewHandler(photoStore app.PhotoStore, thumbnailStore app.ThumbnailStore, thumbnailRepo app.ThumbnailRepo) *Handler {
	h := &Handler{PhotoStore: photoStore,
		ThumbnailStore: thumbnailStore,
		ThumbnailRepo: thumbnailRepo}

	h.thumbnailCreator = &app.ThumbnailCreatorImpl{}

	return h
}


func (h *Handler) logf(format string, args ...interface{}) {
	log.Printf("create handler: " + format,args)
}

func (h *Handler) Handle(in Input) (out app.Output)  {

	// check if the thumbnail is already created; if so end
	thumbnail := app.Thumbnail{}
	if err := h.ThumbnailRepo.Get(in.ThumbnailID, &thumbnail); err != nil {
		out = app.NewUnexpectedErrorOutput()
		msg := fmt.Sprintf("cannot get thumbnail %s from repo", in.PhotoID)
		out.Message = msg
		h.logf("%s / error: %v", msg, err)
		return
	}

	if thumbnail.IsCreated {
		out = app.NewSuccessOutput()
		return
	}

	// fetch the photo from the store
	var buffPhoto bytes.Buffer
	if err := h.PhotoStore.Get(in.PhotoID, &buffPhoto); err != nil {
		out = app.NewUnexpectedErrorOutput()
		msg := fmt.Sprintf("cannot load photo %s from store", in.PhotoID)
		out.Message = msg
		h.logf("%s / error: %v", msg, err)
		return
	}

	// create the thumbnail of the photo
	var buffThumb bytes.Buffer
	if err := h.thumbnailCreator.Create(&buffPhoto,
		in.Format,
		in.Size,
		&buffThumb); err != nil {
		out = app.NewUnexpectedErrorOutput()
		msg := fmt.Sprintf("cannot create image for thumbnail %s", in.ThumbnailID)
		out.Message = msg
		h.logf("%s / error: %v", msg, err)
		return
	}

	// put it to the user store
	if err := h.ThumbnailStore.Put(in.UserID, in.ThumbnailID, &buffThumb); err != nil {
		out = app.NewUnexpectedErrorOutput()
		msg := fmt.Sprintf("cannot save thumbnail %s for user %s to store", in.ThumbnailID, in.UserID)
		out.Message = msg
		h.logf("%s / error: %v", msg, err)
		return
	}

	// update the thumbnail creation job
	now := time.Now().UTC()
	thumbnail.IsCreated = true
	thumbnail.CreatedAt = &now

	if err := h.ThumbnailRepo.Save(&thumbnail); err != nil {
		out = app.NewUnexpectedErrorOutput()
		msg := fmt.Sprintf("cannot update thumbnail %s for user %s to repo", in.ThumbnailID, in.UserID)
		out.Message = msg
		h.logf("%s / error: %v", msg, err)
		return
	}

	out = app.NewSuccessOutput()
	return

}