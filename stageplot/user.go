package stageplot

type User struct {
	ID    string `json:"id" bson:"id"`
	Email string `json:"email" bson:"email"`
	Bands []Band `json:"bands" bson:"bands"`
}
