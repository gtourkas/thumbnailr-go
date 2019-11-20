package request_creation

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
	"thumbnailr/app"
	"thumbnailr/app/create"
	"time"
)

type Input struct {
	UserID  string
	PhotoID string
	Format  string
	Width   int
	Length  int
}

type Handler struct {
	QuotaRepo     app.QuotaRepo
	ThumbnailRepo app.ThumbnailRepo
	CommandIssuer app.CommandIssuer
}

func NewHandler(quotaRepo app.QuotaRepo, thumbnailRepo app.ThumbnailRepo, commandIssuer app.CommandIssuer) *Handler {
	return &Handler{QuotaRepo: quotaRepo, ThumbnailRepo: thumbnailRepo, CommandIssuer: commandIssuer}
}

type OutputData struct {
	ThumbnailID string
}

var (
	ErrBadSize           = errors.New("BAD_SIZE")
	ErrUnsupportedFormat = errors.New("UNSUPPORTED_FORMAT")
	ErrQuotaReached      = errors.New("QUOTA_REACHED")
)

func (h *Handler) logf(format string, args ...interface{}) {
	log.Printf("creation request handler: "+format, args)
}

func (h *Handler) Handle(in Input) (out app.Output) {

	// sanity checks
	if in.Width <= 0 || in.Length <= 0 {
		out = app.NewErrorOutput(ErrBadSize)
		return
	}

	if in.Format != string(app.PNG) &&
		in.Format != string(app.JPEG) {
		out = app.NewErrorOutput(ErrUnsupportedFormat)
		return
	}

	// check the user quota
	var err error
	var state *app.QuotaState
	state, err = h.QuotaRepo.Get(in.UserID)
	if err != nil {
		out = app.NewUnexpectedErrorOutput()
		msg := fmt.Sprintf("cannot get user quota for user %s", in.UserID)
		out.Message = msg
		h.logf("%s / error: %v", msg, err)
		return
	}

	// create quota if the user doesn't have any
	var quota app.Quota
	if state == nil {
		quota = app.NewQuota(in.UserID, time.Now().UTC(), 100)
	} else {
		quota = app.NewQuotaFromState(state)
	}

	// if quota is reached do not proceed
	if quota.IsReached() {
		out = app.NewErrorOutput(ErrQuotaReached)
		return
	}

	thumbnailID := uuid.New().String()
	now := time.Now().UTC()
	if err := h.ThumbnailRepo.Save(&app.Thumbnail{
		ID:          thumbnailID,
		UserID:      in.UserID,
		PhotoID:     in.PhotoID,
		Size:        app.Size{Length: uint(in.Length), Width: uint(in.Width)},
		Format:      app.Format(in.Format),
		IsCreated:   false,
		RequestedAt: now,
	}); err != nil {
		out = app.NewUnexpectedErrorOutput()
		msg := "cannot add thumbnail"
		out.Message = msg
		h.logf("%s / error: %v", msg, err)
		return
	}

	// send the command for thumbnail creation
	if err := h.CommandIssuer.Send(create.Input{
		UserID:      in.UserID,
		ThumbnailID: thumbnailID,
		PhotoID:     in.PhotoID,
		Format:      app.Format(in.Format),
		Size:        app.Size{Length: uint(in.Length), Width: uint(in.Width)},
	}); err != nil {
		out = app.NewUnexpectedErrorOutput()
		msg := "cannot issue command for thumbnail creation"
		out.Message = msg
		h.logf("%s / error: %v", msg, err)
		return
	}

	// add to the user quota
	quota.Add(time.Now().UTC())

	s := quota.GetState()
	state = &s
	if err := h.QuotaRepo.Save(state); err != nil {
		out = app.NewUnexpectedErrorOutput()
		msg := fmt.Sprintf("cannot update user quota for user %s", in.UserID)
		out.Message = msg
		h.logf("%s / error: %v", msg, err)
		return
	}

	out = app.NewSuccessOutput()
	out.Data = OutputData{ThumbnailID: thumbnailID}
	return
}
