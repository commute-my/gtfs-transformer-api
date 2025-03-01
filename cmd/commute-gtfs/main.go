package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/commute-my/gtfs-api/worker"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))

	r := http.NewServeMux()
	// r.HandleFunc("GET /rapidkl/mrt", func(w http.ResponseWriter, r *http.Request) {

	// })

	go func() {
		t := time.NewTicker(5 * time.Second)
		defer t.Stop()

		for {
			select {
			case <-t.C:
				for _, w := range worker.Workers {
					go func() {
						logger.Info("running worker")
						if err := w.Run(); err != nil {
							logger.Error("could not run worker", slog.Any("err", err))
						}
					}()
				}
			}
		}
	}()

	logger.Info("listening and serving http")
	srv := http.Server{
		Addr:    ":8090",
		Handler: r,
	}
	if err := srv.ListenAndServe(); err != nil {
		logger.Error("could not listen and serve http", slog.Any("err", err))
		os.Exit(1)
	}
}
