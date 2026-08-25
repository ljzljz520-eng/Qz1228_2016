package nebula

import (
	"context"
	"gesture-nebula-console/domain"
	"gesture-nebula-console/notifier"
	"testing"
	"time"
)

func TestNotifierDispatch(t *testing.T) {
	s, _ := openFixture(t)
	sender := &notifier.MemorySender{}
	n := notifier.New(s, sender)
	r := domain.NewRecord("notify", "u1", "label", time.Now())
	if err := s.SaveRecord(r); err != nil {
		t.Fatal(err)
	}
	if err := n.Dispatch(context.Background(), r); err != nil {
		t.Fatal(err)
	}
	if len(sender.Snapshot()) != 1 {
		t.Fatal("message not sent")
	}
}
