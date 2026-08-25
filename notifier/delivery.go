package notifier

import (
	"context"
	"fmt"
	"time"

	"gesture-nebula-console/domain"
)

type Delivery struct {
	ID       string    `json:"id"`
	RecordID string    `json:"record_id"`
	Outcome  string    `json:"outcome"`
	SentAt   time.Time `json:"sent_at"`
}

func PrepareDelivery(record domain.Record, outcome string, now time.Time) Delivery {
	return Delivery{ID: fmt.Sprintf("delivery-%s", record.ID), RecordID: record.ID, Outcome: outcome, SentAt: now.UTC()}
}

func SendWithRetry(ctx context.Context, sender Sender, msg Message, attempts int) error {
	if attempts < 1 {
		attempts = 1
	}
	var err error
	for i := 0; i < attempts; i++ {
		if e := ctx.Err(); e != nil {
			return e
		}
		err = sender.Send(ctx, msg)
		if err == nil {
			return nil
		}
		select {
		case <-time.After(time.Millisecond):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return err
}
