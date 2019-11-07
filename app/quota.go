package app

import (
	"errors"
	"time"
)

type QuotaState struct {
	UserID string
	Month uint
	Year uint
	Current uint
	Limit uint
	Reached bool
}

type Quota struct {
	userID string
	month uint
	year uint
	current uint
	limit uint
	reached bool
}

func(q *Quota) IsReached() bool {
	return q.reached
}

var (
	ErrQuotaReached error = errors.New("quota is reached")
	ErrCannotAddToPastDate error = errors.New("cannot add to quota of a past date")
)

func(q *Quota) Add(now time.Time) error {

	nowMonth := uint(now.Month())
	nowYear := uint(now.Year())

	if nowMonth < q.month && nowYear <= q.year {
		return ErrCannotAddToPastDate
	}

	// reset quota if on a later date
	if nowMonth > q.month && nowYear >= q.year {
		q.month = nowMonth
		q.year = nowYear
		q.current = 0
		q.reached = false
	}

	if q.reached {
		return ErrQuotaReached
	}

	q.current ++

	if q.limit == q.current {
		q.reached = true
	}

	return nil
}

func NewQuota(userID string, now time.Time, limit uint) Quota {
	return Quota{
		userID: userID,
		month:   uint(now.Month()),
		year:    uint(now.Year()),
		current: 0,
		limit:   limit,
		reached: false,
	}
}

func NewQuotaFromState(state *QuotaState) Quota {
	return Quota{
		userID: state.UserID,
		month: state.Month,
		year: state.Year,
		current: state.Current,
		limit: state.Limit,
		reached: state.Reached,
	}
}

func(q *Quota) GetState() QuotaState {
	return QuotaState{
		UserID: q.userID,
		Month: q.month,
		Year: q.year,
		Current: q.current,
		Limit: q.limit,
		Reached: q.reached,
	}
}

type QuotaRepo interface {
	Get(userID string, out *QuotaState) error
	Save(quota *QuotaState) error
}
