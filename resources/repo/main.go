package repo

import (
	"go.mongodb.org/mongo-driver/mongo"
)

var planetsCollection *mongo.Collection
var db *mongo.Client

func Initialize(_db *mongo.Client, databaseName string) {
	db = _db
	planetsCollection = db.Database(databaseName).Collection("planets")
}
