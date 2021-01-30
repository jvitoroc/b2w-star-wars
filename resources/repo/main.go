package repo

import (
	"go.mongodb.org/mongo-driver/mongo"
)

var planetsCollection *mongo.Collection
var db *mongo.Client

func Initialize(_db *mongo.Client) {
	db = _db
	planetsCollection = db.Database("b2w").Collection("planets")
}
