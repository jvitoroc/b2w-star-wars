package resources

import (
	"github.com/gorilla/mux"
	"github.com/jvitoroc/b2w-star-wars/resources/handlers"
	"github.com/jvitoroc/b2w-star-wars/resources/repo"
	"go.mongodb.org/mongo-driver/mongo"
)

func Initialize(r *mux.Router, db *mongo.Client) {
	repo.Initialize(db)
	handlers.Initialize(r)
}
