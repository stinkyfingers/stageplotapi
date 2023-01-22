package server

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/stinkyfingers/stageplotapi/stageplot"
)

func (s *Server) Status(w http.ResponseWriter, r *http.Request) {
	msg := "ok"
	if err := s.Storage.Status(r.Context()); err != nil {
		msg = err.Error()
	}

	status := struct {
		Message string `json:"message"`
	}{
		Message: msg,
	}
	s.JSON(w, status, nil)
}

func (s *Server) GetUser(w http.ResponseWriter, r *http.Request) {
	id := httprouter.ParamsFromContext(r.Context()).ByName("id")
	user, err := s.Storage.GetUser(r.Context(), id)
	if err != nil {
		s.JSONErr(w, err, http.StatusInternalServerError)
		return
	}
	s.JSON(w, user, nil)
}

func (s *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user stageplot.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		s.JSONErr(w, err, http.StatusInternalServerError)
		return
	}

	err = s.Storage.CreateUser(r.Context(), &user)
	if err != nil {
		s.JSONErr(w, err, http.StatusInternalServerError)
		return
	}
	s.JSON(w, user, nil)
}
