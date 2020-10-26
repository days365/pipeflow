package exporter

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/pubsub"
)

var (
	TmpMsgsLimit   = 5
	ExportInterval = time.Second * 30
)

const (
	tmpMsgsCap = 256
)

type Executer interface {
	Do(context.Context, []*pubsub.Message) error
}

type Exporter struct {
	queue      chan *pubsub.Message
	exportTime time.Time
	exe        Executer
}

func New(exe Executer) Exporter {
	return Exporter{
		queue:      make(chan *pubsub.Message),
		exportTime: time.Now(),
		exe:        exe,
	}
}

func (e Exporter) Start(ctx context.Context) error {
	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()
	tmpMsgs := make([]*pubsub.Message, 0, tmpMsgsCap)
	for {
		select {
		case <-ticker.C:
			select {
			case msg := <-e.queue:
				tmpMsgs = append(tmpMsgs, msg)
				continue
			default:
				break
			}
		case <-ctx.Done():
			log.Println("exit exporter")
			return nil
		}

		if e.canFlash(tmpMsgs) {
			e.exportTime = time.Now()
			if len(tmpMsgs) > 0 {
				if err := e.exe.Do(ctx, tmpMsgs); err != nil {
					log.Println("[error] failed to executer Do:", err)
				}
				tmpMsgs = make([]*pubsub.Message, 0, tmpMsgsCap)
			}
		}
	}
}

func (e Exporter) In(msg *pubsub.Message) {
	e.queue <- msg
}

func (e Exporter) canFlash(msgs []*pubsub.Message) bool {
	diff := time.Since(e.exportTime)
	return len(msgs) >= TmpMsgsLimit || diff > ExportInterval
}
