package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/stinkyfingers/stageplotapi/stageplot"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Mongo struct {
	Client *mongo.Client
	DB     *mongo.Database
}

const (
	databaseName        = "stageplot"
	userCollection      = "user"
	stagePlotCollection = "stageplot"
)

var _ Storage = (*Mongo)(nil)

var (
	ErrAlreadyExist = fmt.Errorf("record already exists")
)

func NewMongo() (*Mongo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil, err
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	return &Mongo{
		Client: client,
		DB:     client.Database(databaseName),
	}, nil
}

func (m *Mongo) Status(ctx context.Context) error {
	return m.Client.Ping(ctx, nil)
}

func (m *Mongo) GetUser(ctx context.Context, id string) (*stageplot.User, error) {
	var u stageplot.User
	err := m.DB.Collection(userCollection).FindOne(ctx, bson.M{"id": id}).Decode(&u)
	return &u, err
}

func (m *Mongo) CreateUser(ctx context.Context, u *stageplot.User) error {
	user, _ := m.GetUser(ctx, u.ID)
	if user != nil {
		return ErrAlreadyExist
	}
	_, err := m.DB.Collection(userCollection).InsertOne(ctx, u)
	return err
}
