package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	IsPresent(guid string) (bool, error)
	UpdateRefresh(guid string, refresh string) error
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
func (r *UserRep) UpdateRefresh(guid string, refresh string) error {
	id, err := primitive.ObjectIDFromHex(guid)
	if err != nil {
		return err
	}
	hashToken, err := bcrypt.GenerateFromPassword([]byte(refresh), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err = r.collection.UpdateByID(ctx, id, bson.D{{Key: "$set", Value: bson.D{{Key: "refreshtoken", Value: string(hashToken)}}}}, options.Update().SetUpsert(true))
	return err
}
