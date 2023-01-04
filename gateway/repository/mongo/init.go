package mongodb

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoRepo interface {
}

type mongoRepository struct {
	mongodb *mongo.Database
}

func NewMongoRepository(originalMongo *mongo.Database) MongoRepo {
	return &mongoRepository{
		mongodb: originalMongo,
	}
}
