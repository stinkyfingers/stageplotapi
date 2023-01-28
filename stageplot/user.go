package stageplot

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           string      `json:"id" bson:"id"`
	Name         string      `json:"name" bson:"name"`
	Email        string      `json:"email" bson:"email"`
	Token        string      `json:"token" bson:"token"`
	LastLogin    *time.Time  `json:"lastLogin" bson:"lastLogin"`
	StagePlotIDs []uuid.UUID `json:"stagePlotIds" bson:"stagePlotIds"`
}

func NewUser(id, name, email, token string) *User {
	return &User{
		ID:           id,
		Name:         name,
		Email:        email,
		Token:        token,
		StagePlotIDs: []uuid.UUID{},
	}
}
