package notifier

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"gesture-nebula-console/domain"
	"gesture-nebula-console/store"
)

type Message struct {
	RecordID string
	Label    string
	Status   domain.Status
}
type Sender interface {
	Send(context.Context, Message) error
}

type MemorySender struct {
	mu       sync.Mutex
	Messages []Message
	Fail     bool
}

func (m *MemorySender) Send(_ context.Context, msg Message) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.Fail {
		return errors.New("sender failure")
	}
	m.Messages = append(m.Messages, msg)
	return nil
}
func (m *MemorySender) Snapshot() []Message {
	m.mu.Lock()
	defer m.mu.Unlock()
	return append([]Message(nil), m.Messages...)
}

type Notifier struct {
	Store  *store.Store
	Sender Sender
	now    func() time.Time
	seq    uint64
	mu     sync.Mutex
}

func New(s *store.Store, sender Sender) *Notifier {
	return &Notifier{Store: s, Sender: sender, now: func() time.Time { return time.Now().UTC().Round(time.Microsecond) }}
}

func (n *Notifier) nextID() string {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.seq++
	return fmt.Sprintf("notify-%06d", n.seq)
}

func (n *Notifier) Dispatch(ctx context.Context, record domain.Record) error {
	if n.Store == nil || n.Sender == nil {
		return errors.New("notifier is not configured")
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	err := n.Sender.Send(ctx, Message{RecordID: record.ID, Label: record.Label, Status: record.Status})
	outcome := "delivered"
	if err != nil {
		outcome = "failed"
	}
	auditErr := n.Store.SaveAudit(domain.Audit{ID: n.nextID(), RecordID: record.ID, Action: "notify", Outcome: outcome, CreatedAt: n.now()})
	if err != nil {
		return fmt.Errorf("send notification: %w", err)
	}
	return auditErr
}
