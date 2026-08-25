package domain

import (
	"fmt"
	"time"
)

func Transition(from, to Status) error {
	if !validStatus(from) || !validStatus(to) {
		return ErrInvalidStatus
	}
	if from == to {
		return nil
	}
	allowed := map[Status]map[Status]bool{
		StatusReceived:   {StatusReviewed: true, StatusProcessing: true, StatusCancelled: true},
		StatusReviewed:   {StatusProcessing: true, StatusArchived: true, StatusCancelled: true},
		StatusProcessing: {StatusReviewed: true, StatusArchived: true, StatusCancelled: true},
		StatusArchived:   {}, StatusCancelled: {},
	}
	if !allowed[from][to] {
		return fmt.Errorf("%w: %s -> %s", ErrInvalidTransition, from, to)
	}
	return nil
}

func (r *Record) ApplyTransition(to Status, now time.Time) error {
	if err := Transition(r.Status, to); err != nil {
		return err
	}
	r.Status = to
	r.Touch(now)
	return nil
}

func NextStatusFor(action string) (Status, error) {
	switch action {
	case "review":
		return StatusReviewed, nil
	case "process":
		return StatusProcessing, nil
	case "archive":
		return StatusArchived, nil
	case "cancel":
		return StatusCancelled, nil
	default:
		return "", ErrInvalidTransition
	}
}
