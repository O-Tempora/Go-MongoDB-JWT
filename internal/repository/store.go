package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Store struct {
	db      *mongo.Database
	userRep UserRepository
}

func CreateStore(db *mongo.Database) *Store {
	return &Store{
		db: db,
	}
}
func (s *Store) User() UserRepository {
	if s.userRep != nil {
		return s.userRep
	}
	s.userRep = &UserRep{
		store:      s,
		collection: s.db.Collection("users", nil),
	}
	return s.userRep
}
