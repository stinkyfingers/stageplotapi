package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"

	"github.com/stinkyfingers/stageplotapi/storage"
)

// TODO limit API to UI

type Server struct {
	Port    string
	Storage storage.Storage
}

func NewServer(port string, store storage.Storage) *Server {
	return &Server{
		Port:    fmt.Sprintf(":%s", port),
		Storage: store,
	}
}

func (s *Server) Run() error {
	mux := httprouter.New()
	mux.HandlerFunc("GET", "/user/:id", s.GetUser)
	mux.HandlerFunc("POST", "/user", s.CreateUser)
	mux.HandlerFunc("GET", "/status", s.Status)
	return http.ListenAndServe(s.Port, cors.Default().Handler(mux))
}

func (s *Server) JSON(w http.ResponseWriter, data interface{}, dataErr error) {
	if dataErr != nil {
		s.JSONErr(w, dataErr, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		s.JSONErr(w, err, http.StatusInternalServerError)
	}
}

func (s *Server) JSONErr(w http.ResponseWriter, dataErr error, code int) {
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(map[string]interface{}{
		"error": dataErr.Error(),
		"code":  code,
	})
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}
