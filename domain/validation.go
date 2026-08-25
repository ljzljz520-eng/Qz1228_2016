package domain

import (
	"errors"
	"strings"
)

func ValidateUser(u User) error {
	if strings.TrimSpace(u.ID) == "" {
		return errors.New("user id is required")
	}
	if strings.TrimSpace(u.Name) == "" {
		return errors.New("user name is required")
	}
	if strings.TrimSpace(u.Role) == "" {
		return errors.New("user role is required")
	}
	if !u.Active {
		return errors.New("user is inactive")
	}
	return nil
}

func ValidateLabel(label string) error {
	clean := strings.TrimSpace(label)
	if clean == "" {
		return ErrEmptyLabel
	}
	if len([]rune(clean)) > 80 {
		return errors.New("label exceeds 80 characters")
	}
	return nil
}

func NormalizeLabel(label string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(label)), " ")
}

func ValidateVersion(expected, actual int) error {
	if expected <= 0 || actual <= 0 {
		return ErrStaleVersion
	}
	if expected != actual {
		return ErrStaleVersion
	}
	return nil
}

func IsTerminal(status Status) bool { return status == StatusArchived || status == StatusCancelled }

func CanEdit(status Status) bool {
	return status == StatusReceived || status == StatusReviewed || status == StatusProcessing
}
