package stageplot

import "github.com/google/uuid"

type StagePlot struct {
	ID          uuid.UUID   `json:"id" bson:"id"`
	Name        string      `json:"name" bson:"name"`
	IsPublic    bool        `json:"isPublic" bson:"isPublic"`
	Icons       []Icon      `json:"icons" bson:"icons"`
	InputList   InputList   `json:"inputList" bson:"inputList"`
	MonitorList MonitorList `json:"monitorList" bson:"monitorList"`
}

type Icon struct {
	Name     string `json:"name" bson:"monitorList"`
	Filepath string `json:"filepath" bson:"filepath"`
	Position [2]int `json:"position" bson:"position"`
}

type InputList struct {
	Channel int    `json:"channel" bson:"channel"`
	Name    string `json:"name" bson:"name"`
}

type MonitorList struct {
	Channel       int            `json:"channel" bson:"channel"`
	Name          string         `json:"name" bson:"name"`
	MonitorLevels []MonitorLevel `json:"monitorLevels" bson:"monitorLevels"`
}

type MonitorLevel struct {
	InputChannel int `json:"inputChannel" bson:"inputChannel"`
	Level        int `json:"level" bson:"level"`
}
