package repo

import (
	"errors"
	"fmt"

	"github.com/jvitoroc/b2w-star-wars/resources/common"
	"github.com/jvitoroc/b2w-star-wars/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Planet struct {
	ObjectID        primitive.ObjectID `json:"id" bson:"_id"`
	Name            string             `json:"name"`
	Climate         string             `json:"climate"`
	Terrain         string             `json:"terrain"`
	FilmsAppearedIn int                `json:"filmsAppearedIn"`
}

func CreatePlanet(name string, climate string, terrain string, filmsAppearedIn int) (primitive.ObjectID, *common.Error) {
	ctx, cancel := utils.WithTimeout(5)
	defer cancel()
	res, err := planetsCollection.InsertOne(ctx, bson.D{
		{Key: "name", Value: name},
		{Key: "climate", Value: climate},
		{Key: "terrain", Value: terrain},
		{Key: "filmsAppearedIn", Value: filmsAppearedIn},
	})

	if err != nil {
		return primitive.ObjectID{}, common.CreateGenericInternalError(err)
	} else {
		return res.InsertedID.(primitive.ObjectID), nil
	}
}

func GetPlanetByID(id primitive.ObjectID) (*Planet, *common.Error) {
	ctx, cancel := utils.WithTimeout(5)
	defer cancel()
	planet := Planet{}
	filter := bson.D{{Key: "_id", Value: id}}
	err := planetsCollection.FindOne(ctx, filter).Decode(&planet)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, common.CreateNotFoundError(fmt.Sprintf("Planet not found under given id (%s).", id))
		}
		return nil, common.CreateGenericInternalError(err)
	}

	return &planet, nil
}

func GetMatchedPlanets(criteria map[string]string) ([]*Planet, *common.Error) {
	ctx, cancel := utils.WithTimeout(5)
	defer cancel()
	planets := make([]*Planet, 0)
	filter := bson.D{}

	if criteria != nil {
		for k, v := range criteria {
			if v != "" {
				filter = append(filter, bson.E{Key: k, Value: primitive.Regex{Pattern: v, Options: "i"}})
			}
		}
	}

	cur, err := planetsCollection.Find(ctx, filter)
	if err != nil {
		return nil, common.CreateGenericInternalError(err)
	}

	for cur.Next(ctx) {
		var planet Planet
		err := cur.Decode(&planet)

		if err != nil {
			return nil, common.CreateGenericInternalError(err)
		}

		planets = append(planets, &planet)
	}

	return planets, nil
}

func GetAllPlanets() ([]*Planet, *common.Error) {
	return GetMatchedPlanets(nil)
}

func DeletePlanet(id primitive.ObjectID) *common.Error {
	ctx, cancel := utils.WithTimeout(5)
	defer cancel()
	filter := bson.D{{Key: "_id", Value: id}}
	res, err := planetsCollection.DeleteOne(ctx, filter)

	if res.DeletedCount == 0 && err == nil {
		return common.CreateNotFoundError(fmt.Sprintf("Planet not found under given id (%s).", id))
	} else if err != nil {
		return common.CreateGenericInternalError(err)
	}

	return nil
}
