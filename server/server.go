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
	mux.HandlerFunc("GET", "/user/:id", s.AuthMiddleware(s.GetUser))
	mux.HandlerFunc("POST", "/user", s.AuthMiddleware(s.CreateUser))
	mux.HandlerFunc("PUT", "/user", s.AuthMiddleware(s.UpdateUser))
	mux.HandlerFunc("DELETE", "/user/:id", s.AuthMiddleware(s.DeleteUser))
	mux.HandlerFunc("POST", "/user/login", s.LoginUser)

	mux.HandlerFunc("GET", "/plot/:plotId", s.AuthMiddleware(s.GetStagePlot))
	mux.HandlerFunc("POST", "/plot", s.AuthMiddleware(s.CreateStagePlot))
	mux.HandlerFunc("PUT", "/plot", s.AuthMiddleware(s.UpdateStagePlot))
	mux.HandlerFunc("DELETE", "/plot/:plotId", s.AuthMiddleware(s.DeleteStagePlot))
	mux.HandlerFunc("GET", "/plots", s.AuthMiddleware(s.ListStagePlots))

	mux.HandlerFunc("GET", "/status", s.Status)
	return http.ListenAndServe(s.Port, s.CORS(cors.Default().Handler(mux)))
}

func (s *Server) JSON(w http.ResponseWriter, data interface{}, dataErr error) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Type", "application/json")
	if dataErr != nil {
		s.JSONErr(w, dataErr, http.StatusInternalServerError)
		return
	}
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		s.JSONErr(w, err, http.StatusInternalServerError)
	}
}

func (s *Server) JSONErr(w http.ResponseWriter, dataErr error, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(map[string]interface{}{
		"error": dataErr.Error(),
		"code":  code,
	})
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}
