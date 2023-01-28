package storage

import (
	"context"
	"fmt"
	"os"
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

	url := "mongodb://localhost:27017"
	if val := os.Getenv("MONGO_URL"); val != "" {
		url = val
	}
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
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
	_, err := m.DB.Collection(userCollection).InsertOne(ctx, u)
	return err
}

func (m *Mongo) UpdateUser(ctx context.Context, u *stageplot.User) error {
	_, err := m.DB.Collection(userCollection).UpdateOne(ctx, bson.M{
		"id": u.ID,
	}, bson.M{
		"$set": bson.M{
			"email":        u.Email,
			"name":         u.Name,
			"stagePlotIds": u.StagePlotIDs,
			"lastLogin":    u.LastLogin,
			"token":        u.Token,
		},
	})
	return err
}
func (m *Mongo) DeleteUser(ctx context.Context, u *stageplot.User) error {
	_, err := m.DB.Collection(userCollection).DeleteOne(ctx, bson.M{"id": u.ID})
	return err
}

func (m *Mongo) LookupToken(ctx context.Context, token string) (*stageplot.User, error) {
	var user stageplot.User
	err := m.DB.Collection(userCollection).FindOne(ctx, bson.M{"token": token}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (m *Mongo) Get(ctx context.Context, obj stageplot.IDer) error {
	err := m.DB.Collection(stagePlotCollection).FindOne(ctx, bson.M{"id": obj.GetID()}).Decode(obj)
	return err
}

func (m *Mongo) Create(ctx context.Context, obj stageplot.IDer) error {
	_, err := m.DB.Collection(stagePlotCollection).InsertOne(ctx, obj)
	return err
}

func (m *Mongo) Replace(ctx context.Context, obj stageplot.IDer) error {
	_, err := m.DB.Collection(stagePlotCollection).ReplaceOne(ctx, bson.M{
		"id": obj.GetID(),
	}, obj)
	return err
}
func (m *Mongo) Delete(ctx context.Context, obj stageplot.IDer) error {
	_, err := m.DB.Collection(stagePlotCollection).DeleteOne(ctx, bson.M{"id": obj.GetID()})
	return err
}

func (m *Mongo) List(ctx context.Context, u *stageplot.User) ([]stageplot.StagePlot, error) {
	var user stageplot.User
	err := m.DB.Collection(userCollection).FindOne(ctx, bson.M{"id": u.ID}).Decode(&user)
	if err != nil {
		return nil, err
	}
	if len(user.StagePlotIDs) == 0 {
		return nil, nil
	}

	cursor, err := m.DB.Collection(stagePlotCollection).Find(ctx, bson.M{"id": bson.M{"$in": user.StagePlotIDs}})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	var stagePlots []stageplot.StagePlot
	for cursor.Next(ctx) {
		var stagePlot stageplot.StagePlot
		err = cursor.Decode(&stagePlot)
		if err != nil {
			return nil, err
		}
		stagePlots = append(stagePlots, stagePlot)
	}
	return stagePlots, err
}
