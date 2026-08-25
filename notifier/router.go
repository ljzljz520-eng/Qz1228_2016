package notifier

import (
	"context"
	"errors"
	"strings"
	"sync"

	"gesture-nebula-console/domain"
)

type Rule struct {
	Prefix string
	Sender Sender
}
type Router struct {
	mu       sync.RWMutex
	rules    []Rule
	fallback Sender
}

func NewRouter(fallback Sender) *Router { return &Router{fallback: fallback, rules: make([]Rule, 0)} }

func (r *Router) AddRule(prefix string, sender Sender) error {
	if strings.TrimSpace(prefix) == "" || sender == nil {
		return errors.New("prefix and sender are required")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.rules = append(r.rules, Rule{Prefix: prefix, Sender: sender})
	return nil
}

func (r *Router) resolve(label string) Sender {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, rule := range r.rules {
		if strings.HasPrefix(label, rule.Prefix) {
			return rule.Sender
		}
	}
	return r.fallback
}

func (r *Router) Send(ctx context.Context, message Message) error {
	sender := r.resolve(message.Label)
	if sender == nil {
		return errors.New("no sender configured")
	}
	return sender.Send(ctx, message)
}

func (r *Router) SendRecord(ctx context.Context, record domain.Record) error {
	return r.Send(ctx, DeliveryMessage(record))
}

func (r *Router) Rules() []Rule {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return append([]Rule(nil), r.rules...)
}

func RouteByStatus(status domain.Status) string {
	switch status {
	case domain.StatusArchived:
		return "archive"
	case domain.StatusCancelled:
		return "cancel"
	case domain.StatusProcessing:
		return "process"
	default:
		return "active"
	}
}

func DeliveryMessage(record domain.Record) Message {
	return Message{RecordID: record.ID, Label: record.Label, Status: record.Status}
}
