package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository interface {
	IsPresent(guid string) (bool, error)
	UpdateRefresh(guid primitive.ObjectID, refresh string) error
}

type UserRep struct {
	store      *Store
	collection *mongo.Collection
}

func (r *UserRep) IsPresent(guid string) (bool, error) {
	id, err := primitive.ObjectIDFromHex(guid)
	if err != nil {
		return false, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	count, err := r.collection.CountDocuments(ctx, bson.D{{Key: "_id", Value: /*fmt.Sprintf("ObjectId(%s)", guid)*/ id}}, nil)
	if err != nil {
		return false, err
	} else if count == 0 {
		return false, nil
	}
	return true, nil
}
func (r *UserRep) UpdateRefresh(guid primitive.ObjectID, refresh string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := r.collection.UpdateByID(ctx, guid, bson.D{{Key: "$set", Value: bson.D{{Key: "refreshtoken", Value: refresh}}}}, options.Update().SetUpsert(true))
	return err
}
