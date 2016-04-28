package main

type Category struct {
	ID string `bson:"_id" json:"_id"`
	Name string `bson:"name" json:"name"`
	Pivot string `bson:"pivot" json:"pivot"`
}
