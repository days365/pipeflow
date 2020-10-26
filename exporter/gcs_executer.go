package exporter

import (
	"context"
	"fmt"
	"log"
	"path"
	"time"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
)

type StoreToGCSExecuter struct {
	c      *storage.Client
	prefix string
	bucket string
}

func NewStoreToGCSExecuter(ctx context.Context, prefix, bucket string) (Executer, error) {
	c, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return StoreToGCSExecuter{
		c:      c,
		prefix: prefix,
		bucket: bucket,
	}, nil
}

func (e StoreToGCSExecuter) Do(ctx context.Context, msgs []*pubsub.Message) error {
	savePath := path.Join(e.prefix, fmt.Sprintf("pipeflow-%s.json", genTimestamp()))
	w := e.c.Bucket(e.bucket).Object(savePath).NewWriter(ctx)
	w.ContentType = "application/json; charset=utf8"
	defer func() {
		if err := w.Close(); err != nil {
			log.Println("[error] failed to Close", err)
		}
	}()

	for _, msg := range msgs {
		if _, err := w.Write(msg.Data); err != nil {
			log.Println("[error] failed to write message data:", err)
			continue
		}
		msg.Ack()
	}

	return nil
}

func genTimestamp() string {
	n := time.Now()
	return fmt.Sprintf("%d-%d-%d-%d-%d-%d", n.Year(), int(n.Month()), n.Day(), n.Hour(), n.Minute(), n.Second())
}
