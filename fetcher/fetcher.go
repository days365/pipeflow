package fetcher

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/days365/pipeflow/exporter"
)

type Fetcher struct {
	exp exporter.Exporter
}

func New(exp exporter.Exporter) Fetcher {
	return Fetcher{
		exp: exp,
	}
}

func (f Fetcher) Start(ctx context.Context, projectID, sub string) error {
	c, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return err
	}
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	sig := make(chan os.Signal, 1)

	signal.Notify(sig,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer signal.Stop(sig)

	subsc := c.Subscription(sub)
	for {
		select {
		case <-ticker.C:
			err = subsc.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
				f.exp.In(m)
			})
			if err != nil {
				log.Println("[error] subsc.Receive failed:", err)
				return err
			}

		case s := <-sig:
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				log.Println("stop pipeflow")
				return nil
			}
		case <-ctx.Done():
			log.Println("stop pipeflow")
			return nil
		}
	}
}
