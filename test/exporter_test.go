package exporter_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/days365/pipeflow/exporter"
)

type MockExecuter struct {
	msgs chan *pubsub.Message
}

func (m MockExecuter) Do(ctx context.Context, msgs []*pubsub.Message) error {
	for _, msg := range msgs {
		m.msgs <- msg
	}
	return nil
}

func TestExporterStart(t *testing.T) {
	c := make(chan *pubsub.Message)
	e := MockExecuter{msgs: c}
	exp := exporter.New(e)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		if err := exp.Start(ctx); err != nil {
			t.Error("exporter Start failed:", err)
			return
		}
	}()

	for _, id := range []string{"1", "2", "3", "4", "5"} {
		exp.In(
			&pubsub.Message{
				Data: []byte("message id " + id),
			},
		)
	}

	var cnt int = 1
	for msg := range c {
		want := fmt.Sprintf("message id %d", cnt)
		got := string(msg.Data)
		if got != want {
			t.Errorf("message id is wrong: want %s, but got %s", want, got)
		}
		cnt++
		if cnt > 5 {
			break
		}
	}
	now := time.Now()

	exporter.ExportInterval = 3 * time.Second
	for _, id := range []string{"1", "2"} {
		exp.In(
			&pubsub.Message{
				Data: []byte("message id " + id),
			},
		)
	}

	<-c
	diff := time.Since(now)
	if diff < 3*time.Second {
		t.Error("message queue flash interval is failed:", diff)
	}

	cancel()
}
