package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/stinkyfingers/stageplotapi/stageplot"
	"github.com/stinkyfingers/stageplotapi/storage"
	"go.mongodb.org/mongo-driver/mongo"
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

func (s *Server) LoginUser(w http.ResponseWriter, r *http.Request) {
	var u stageplot.User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		s.JSONErr(w, err, http.StatusInternalServerError)
		return
	}
	now := time.Now()
	user, err := s.Storage.GetUser(r.Context(), u.ID)
	if err != nil && err != mongo.ErrNoDocuments {
		s.JSONErr(w, err, http.StatusInternalServerError)
		return
	}
	if err != mongo.ErrNoDocuments {
		user.LastLogin = &now
		user.Token = u.Token
		if err = s.Storage.UpdateUser(r.Context(), user); err != nil {
			s.JSONErr(w, err, http.StatusInternalServerError)
			return
		}
	} else {
		user = stageplot.NewUser(u.ID, u.Name, u.Email, u.Token)
		user.LastLogin = &now
		if err = s.Storage.CreateUser(r.Context(), user); err != nil {
			s.JSONErr(w, err, http.StatusInternalServerError)
			return
		}
	}
	s.JSON(w, user, nil)
}

// TODO unneeded? LoginUser does everything here
func (s *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	var u stageplot.User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		s.JSONErr(w, err, http.StatusInternalServerError)
		return
	}
	currentUser, _ := s.Storage.GetUser(r.Context(), u.ID)
	if currentUser != nil {
		s.JSONErr(w, storage.ErrAlreadyExist, http.StatusConflict)
		return
	}
	user := stageplot.NewUser(u.ID, u.Name, u.Email, u.Token)
	err = s.Storage.CreateUser(r.Context(), user)
	if err != nil {
		s.JSONErr(w, err, http.StatusInternalServerError)
		return
	}
	s.JSON(w, user, nil)
}

func (s *Server) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var u stageplot.User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		s.JSONErr(w, err, http.StatusInternalServerError)
		return
	}
	err = s.Storage.UpdateUser(r.Context(), &u)
	if err != nil {
		s.JSONErr(w, err, http.StatusInternalServerError)
		return
	}
	s.JSON(w, u, nil)
}

func (s *Server) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := httprouter.ParamsFromContext(r.Context()).ByName("id")
	u := stageplot.User{
		ID: id,
	}
	err := s.Storage.DeleteUser(r.Context(), &u)
	if err != nil {
		s.JSONErr(w, err, http.StatusInternalServerError)
		return
	}
	s.JSON(w, u, nil)
}

func (s *Server) GetStagePlot(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*stageplot.User)
	if !ok {
		s.JSONErr(w, fmt.Errorf("invalid token"), http.StatusUnauthorized)
		return
	}
	id := httprouter.ParamsFromContext(r.Context()).ByName("plotId")
	uid, err := uuid.Parse(id)
	if err != nil {
		s.JSONErr(w, err, http.StatusInternalServerError)
		return
	}
	found := false
	for _, plotId := range user.StagePlotIDs {
		if plotId == uid {
			found = true
			break
		}
	}
	if !found {
		s.JSONErr(w, fmt.Errorf("user cannot access stage plot"), http.StatusUnauthorized)
		return
	}
	sp := &stageplot.StagePlot{ID: uid}
	err = s.Storage.Get(r.Context(), sp)
	if err != nil {
		s.JSONErr(w, err, http.StatusInternalServerError)
		return
	}
	s.JSON(w, sp, nil)
}

func (s *Server) CreateStagePlot(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*stageplot.User)
	if !ok {
		s.JSONErr(w, fmt.Errorf("invalid token"), http.StatusUnauthorized)
		return
	}
	var sp stageplot.StagePlot
	err := json.NewDecoder(r.Body).Decode(&sp)
	if err != nil {
		s.JSONErr(w, err, http.StatusInternalServerError)
		return
	}
	sp.ID = uuid.New()
	err = s.Storage.Create(r.Context(), &sp)
	if err != nil {
		s.JSONErr(w, err, http.StatusInternalServerError)
		return
	}
	user.StagePlotIDs = append(user.StagePlotIDs, sp.ID)
	err = s.Storage.UpdateUser(r.Context(), user)
	if err != nil {
		s.JSONErr(w, err, http.StatusInternalServerError)
		return
	}
	s.JSON(w, sp, nil)
}

func (s *Server) UpdateStagePlot(w http.ResponseWriter, r *http.Request) {
	var sp stageplot.StagePlot
	err := json.NewDecoder(r.Body).Decode(&sp)
	if err != nil {
		s.JSONErr(w, err, http.StatusInternalServerError)
		return
	}

	err = s.Storage.Replace(r.Context(), &sp)
	if err != nil {
		s.JSONErr(w, err, http.StatusInternalServerError)
		return
	}
	s.JSON(w, sp, nil)
}

func (s *Server) DeleteStagePlot(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*stageplot.User)
	if !ok {
		s.JSONErr(w, fmt.Errorf("invalid token"), http.StatusUnauthorized)
		return
	}

	id := httprouter.ParamsFromContext(r.Context()).ByName("id")
	uid, err := uuid.Parse(id)
	if err != nil {
		s.JSONErr(w, err, http.StatusInternalServerError)
		return
	}
	sp := &stageplot.StagePlot{ID: uid}

	err = s.Storage.Delete(r.Context(), sp)
	if err != nil {
		s.JSONErr(w, err, http.StatusInternalServerError)
		return
	}
	for i, plotId := range user.StagePlotIDs {
		if plotId == uid {
			user.StagePlotIDs = append(user.StagePlotIDs[:i], user.StagePlotIDs[i+1:]...)
			break
		}
	}
	if err = s.Storage.UpdateUser(r.Context(), user); err != nil {
		s.JSONErr(w, err, http.StatusInternalServerError)
		return
	}
	s.JSON(w, sp, nil)
}

func (s *Server) ListStagePlots(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*stageplot.User)
	if !ok {
		s.JSONErr(w, fmt.Errorf("invalid token"), http.StatusUnauthorized)
		return
	}

	plots, err := s.Storage.List(r.Context(), user)
	if err != nil {
		s.JSONErr(w, err, http.StatusInternalServerError)
		return
	}
	s.JSON(w, plots, nil)
}
