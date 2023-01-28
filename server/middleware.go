package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
)

// AuthMiddleware gets the Bearer token, finds the associated user in storage, and places that user on request context
func (s *Server) AuthMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if token == "" {
			s.JSONErr(w, fmt.Errorf("token is required"), http.StatusUnauthorized)
			return
		}
		user, err := s.Storage.LookupToken(r.Context(), token)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				s.JSONErr(w, fmt.Errorf("invalid token"), http.StatusUnauthorized)
				return
			}
			s.JSONErr(w, err, http.StatusInternalServerError)
			return
		}
		h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "user", user)))
	}
}

func (s *Server) CORS(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*") //TODO
		w.Header().Add("Access-Control-Allow-Methods", "*")
		w.Header().Add("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			return
		}
		h.ServeHTTP(w, r)
	}
}
