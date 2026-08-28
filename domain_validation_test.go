package nebula

import (
	"gesture-nebula-console/domain"
	"testing"
)

func TestDomainValidation(t *testing.T) {
	if err := (domain.Record{ID: "x", Gesture: 14, Label: "x", Status: domain.StatusReceived, Version: 1, OwnerID: "u"}).Validate(); err != domain.ErrInvalidGesture {
		t.Fatalf("gesture err %v", err)
	}
	if err := domain.ValidateLabel(" "); err != domain.ErrEmptyLabel {
		t.Fatalf("label err %v", err)
	}
	if err := domain.ValidateUser(domain.User{ID: "u", Name: "x", Role: "r", Active: false}); err == nil {
		t.Fatal("expected inactive user error")
	}
}

func TestStateTransitions(t *testing.T) {
	if err := domain.Transition(domain.StatusReceived, domain.StatusReviewed); err != nil {
		t.Fatal(err)
	}
	if err := domain.Transition(domain.StatusArchived, domain.StatusProcessing); err == nil {
		t.Fatal("expected terminal transition error")
	}
}
