package model

type Topics struct {
	Name     string     `bson:"name"`
	Thoughts []Thoughts `bson:"thoughts"`
}
