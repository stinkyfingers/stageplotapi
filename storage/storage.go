package storage

import (
	"context"

	"github.com/stinkyfingers/stageplotapi/stageplot"
)

type Storage interface {
	Status(ctx context.Context) error
	GetUser(ctx context.Context, id string) (*stageplot.User, error)
	CreateUser(ctx context.Context, u *stageplot.User) error
	//UpdateUser()
	//DeleteUser()
	//GetStagePlot(id uuid.UUID)
	//CreateStagePlot()
	//UpdateStagePlot()
	//DeleteStagePlot()
	//CreateBand()
	//UpdateBand()
	//DeleteStageBand()
}
