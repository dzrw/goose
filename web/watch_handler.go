package web

import (
	"fmt"
	"github.com/drone/routes"
	"github.com/politician/goose/worker"
	"net/http"
	"strconv"
)

func NewWatchHandler(ws worker.WatchService) http.Handler {
	const (
		COLLECTION_URL = "/watches"
		WATCH_URL      = "/watches/:id([a-f0-9]+)"
	)

	// Setup routes.
	mux := routes.New()

	// mux.Get(COLLECTION_URL, watches)
	mux.Post(COLLECTION_URL, closeConnectionAdvice(createFunc(ws)))
	mux.Del(COLLECTION_URL, closeConnectionAdvice(deleteAllFunc(ws)))
	// mux.Get(WATCH_URL, watch)
	mux.Del(WATCH_URL, closeConnectionAdvice(deleteOneFunc(ws)))

	return mux
}

func closeConnectionAdvice(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Connection", "close")
		next(w, req)
	}
}

func createFunc(ws worker.WatchService) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		watch, err := ParseJsonWatch(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id, tag, err := ws.Add(watch)
		if err != nil {
			http.Error(w, err.Error(), GOOSE_ERROR)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, `{"id": %d, "tag": "%s"}`, id, tag)
	}
}

func deleteAllFunc(ws worker.WatchService) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		err := ws.Clear()
		if err != nil {
			http.Error(w, err.Error(), GOOSE_ERROR)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func deleteOneFunc(ws worker.WatchService) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		params := req.URL.Query()

		id, err := strconv.Atoi(params.Get(":id"))
		if err != nil {
			http.Error(w, err.Error(), GOOSE_ERROR)
			return
		}

		_, err = ws.Remove(id)
		if err != nil {
			http.Error(w, err.Error(), GOOSE_ERROR)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
