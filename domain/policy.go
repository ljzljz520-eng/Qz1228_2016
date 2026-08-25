package domain

import (
	"errors"
	"sort"
	"strings"
	"time"
)

type LabelPolicy struct {
	RequiredPrefix string
	Allowed        map[string]bool
	MaxAge         time.Duration
}

func DefaultLabelPolicy() LabelPolicy {
	return LabelPolicy{Allowed: map[string]bool{"amber": true, "cobalt": true, "cobalt arc": true, "violet": true, "rose": true, "old-label": true, "new-label": true, "persistent": true, "label": true}}
}

func (p LabelPolicy) Check(label string) error {
	clean := NormalizeLabel(label)
	if err := ValidateLabel(clean); err != nil {
		return err
	}
	if p.RequiredPrefix != "" && !strings.HasPrefix(clean, p.RequiredPrefix) {
		return errors.New("label does not match required prefix")
	}
	if len(p.Allowed) > 0 && !p.Allowed[clean] {
		return errors.New("label is not allowed")
	}
	return nil
}

func (p LabelPolicy) Options() []string {
	values := make([]string, 0, len(p.Allowed))
	for key := range p.Allowed {
		values = append(values, key)
	}
	sort.Strings(values)
	return values
}

func (p LabelPolicy) Accepts(label string) bool { return p.Check(label) == nil }

type RolePolicy struct {
	Reviewers map[string]bool
	Operators map[string]bool
	Auditors  map[string]bool
}

func DefaultRolePolicy() RolePolicy {
	return RolePolicy{Reviewers: map[string]bool{"reviewer": true, "operator": true}, Operators: map[string]bool{"operator": true}, Auditors: map[string]bool{"auditor": true, "operator": true}}
}

func (p RolePolicy) CanReview(u User) bool   { return u.Active && p.Reviewers[u.Role] }
func (p RolePolicy) CanRegister(u User) bool { return u.Active && p.Operators[u.Role] }
func (p RolePolicy) CanAudit(u User) bool    { return u.Active && p.Auditors[u.Role] }

func FreshAt(t, now time.Time, maxAge time.Duration) bool {
	if maxAge <= 0 {
		return true
	}
	return !t.IsZero() && now.Sub(t) <= maxAge
}

func StableRecordOrder(records []Record) []Record {
	out := append([]Record(nil), records...)
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].UpdatedAt.Equal(out[j].UpdatedAt) {
			return out[i].ID < out[j].ID
		}
		return out[i].UpdatedAt.Before(out[j].UpdatedAt)
	})
	return out
}
