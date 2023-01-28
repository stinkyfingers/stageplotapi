package storage

import (
	"context"

	"github.com/stinkyfingers/stageplotapi/stageplot"
)

type Storage interface {
	Status(ctx context.Context) error
	GetUser(ctx context.Context, id string) (*stageplot.User, error)
	CreateUser(ctx context.Context, u *stageplot.User) error
	UpdateUser(ctx context.Context, u *stageplot.User) error
	DeleteUser(ctx context.Context, u *stageplot.User) error
	LookupToken(ctx context.Context, token string) (*stageplot.User, error)
	Get(ctx context.Context, obj stageplot.IDer) error
	Create(ctx context.Context, obj stageplot.IDer) error
	Replace(ctx context.Context, obj stageplot.IDer) error
	Delete(ctx context.Context, obj stageplot.IDer) error
	List(ctx context.Context, u *stageplot.User) ([]stageplot.StagePlot, error)
}
