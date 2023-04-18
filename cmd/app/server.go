package main

import (
	"encoding/json"
	"net/http"

	"example.com/raft-cache/pkg/cache"

	"github.com/gorilla/mux"
)

type httpServer struct {
	handler http.Handler
	cache   cache.Cache
}

type addCommand struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func newServer(c cache.Cache) *httpServer {
	router := mux.NewRouter()

	s := &httpServer{
		handler: router,
		cache:   c,
	}

	router.HandleFunc("/add", s.addHandler).Methods("POST")
	router.HandleFunc("/get", s.getHandler).Methods("GET")

	return s
}

func (s *httpServer) addHandler(w http.ResponseWriter, r *http.Request) {
	var cmd addCommand

	// Decode the request body into a command.
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.cache.Put(cmd.Key, cmd.Value); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond to the client.
	w.WriteHeader(http.StatusOK)
}

func (s *httpServer) getHandler(w http.ResponseWriter, r *http.Request) {
	// Get the key from the query string.
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "key is required", http.StatusBadRequest)
		return
	}

	// Get the value from the cache.
	value, ok := s.cache.Get(key)
	if !ok {
		http.Error(w, "key not found", http.StatusNotFound)
		return
	}

	// Respond to the client.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"key": key, "value": value})
}
