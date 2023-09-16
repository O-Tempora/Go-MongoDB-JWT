package service

import (
	"errors"
	"gomongojwt/internal/repository"
	"gomongojwt/internal/util"

	"go.mongodb.org/mongo-driver/mongo"
)

type Service interface {
	RefreshTokens(oldAccess, oldRefresh string) (newAccess, newRefresh string, err error)
	AuthorizeUser(guid string) (access, refresh string, err error)
}

type ServiceInstance struct {
	store *repository.Store
	db    *mongo.Database
}

func InitService(store *repository.Store, db *mongo.Database) *ServiceInstance {
	serv := &ServiceInstance{
		store: store,
		db:    db,
	}
	return serv
}
func (s *ServiceInstance) Store() *repository.Store {
	return s.store
}
func (s *ServiceInstance) DB() *mongo.Database {
	return s.db
}

func (s *ServiceInstance) RefreshTokens(oldAccess, oldRefresh string) (newAccess, newRefresh string, err error) {
	guid, err := util.ValidateJWT(oldAccess)
	if err != nil {
		return "", "", err
	}
	same, err := s.store.User().CompareRefreshAndHash(oldRefresh, guid.User)
	if err != nil {
		return "", "", err
	} else if !same {
		return "", "", errors.New("refresh tokens don't match")
	}
	newAccess, newRefresh, err = util.GetTokenPair(guid.User)
	if err != nil {
		return "", "", err
	}
	if err = s.store.User().UpdateRefresh(guid.User, newRefresh); err != nil {
		return "", "", err
	}
	return newAccess, newRefresh, nil
}

func (s *ServiceInstance) AuthorizeUser(guid string) (access, refresh string, err error) {
	access, refresh, err = util.GetTokenPair(guid)
	if err != nil {
		return "", "", err
	}
	err = s.store.User().UpdateRefresh(guid, refresh)
	if err != nil {
		return "", "", err
	}
	return access, refresh, err
}
