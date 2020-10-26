package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/days365/pipeflow/exporter"
	"github.com/days365/pipeflow/fetcher"
)

func main() {
	ctx := context.Background()

	prefix := os.Getenv("BUCKET_PREFIX")
	bucket := os.Getenv("BUCKET_NAME")
	exe, err := exporter.NewStoreToGCSExecuter(ctx, prefix, bucket)
	if err != nil {
		log.Println("[error] failed to executer", err)
		return
	}

	exp := exporter.New(exe)
	go func() {
		if err := exp.Start(ctx); err != nil {
			log.Println("[error] failed to start exporter:", err)
		}
	}()

	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if _, err = w.Write([]byte("ok")); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		})

		if err = http.ListenAndServe(":8080", nil); err != nil {
			log.Println("[error]", err)
		}
	}()

	// for cloud run healthcheck to run continuously.
	go func() {
		endpoint := os.Getenv("HEALTHCHECK_ENDPOINT")
		ticker := time.NewTicker(time.Second * 10)

		if endpoint != "" {
			for {
				select {
				case <-ticker.C:
					res, err := http.DefaultClient.Get(endpoint)
					if err != nil {
						log.Println("[error] failed to healthcheck", err)
						continue
					}
					res.Body.Close()
				default:
					continue
				}
			}
		}
	}()

	f := fetcher.New(exp)
	sub := os.Getenv("PUBSUB_SUBSCRIPTION")
	projectID := os.Getenv("GCP_PROJECT_ID")
	if err := f.Start(ctx, projectID, sub); err != nil {
		log.Println("[error] failed to start featcher:", err)
	}
}
