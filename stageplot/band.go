package stageplot

import "github.com/google/uuid"

type Band struct {
	ID         uuid.UUID   `json:"id" bson:"id"`
	Name       string      `json:"name" bson:"name"`
	StagePlots []StagePlot `json:"stagePlots" bson:"stagePlots"`
}
