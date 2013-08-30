package web

import (
	"fmt"
	"github.com/drone/routes"
	"github.com/politician/goose/worker"
	"net/http"
)

func NewWatchHandler(watcher worker.WatchList) http.Handler {
	const (
		COLLECTION_URL = "/watches"
		WATCH_URL      = "/watches/id:([a-f0-9]+)"
	)

	// Setup routes.
	mux := routes.New()

	// mux.Get(COLLECTION_URL, watches)
	mux.Post(COLLECTION_URL, createFunc(watcher))
	// mux.Del(COLLECTION_URL, deleteAll)
	// mux.Get(WATCH_URL, watch)
	// mux.Del(WATCH_URL, deleteOne)

	return mux
}

func createFunc(watcher worker.WatchList) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Connection", "close")

		watch, err := ParseJsonWatch(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id, err := watcher.Add(watch)
		if err != nil {
			http.Error(w, err.Error(), CRAZY_ERROR)
			return
		}

		// TODO: Provide a real implementation.
		http.Error(w, fmt.Sprintf(`{"id": %d}`, id), http.StatusCreated)
	}
}
