package mongodb

import "go.mongodb.org/mongo-driver/mongo"

type MongoResult struct {
	Cursor       *mongo.Cursor
	InsertedID   string
	AffectedRows int64
}
