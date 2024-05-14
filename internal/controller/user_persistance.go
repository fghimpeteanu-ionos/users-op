package controller

import (
	userv1 "github.com/fghimpeteanu-ionos/user-op/api/v1"
	"github.com/google/uuid"
)

type PersistUser interface {
	Persist(user *userv1.User) (string, error)
}

type ReadUser interface {
	Read(uuid string) (*userv1.User, error)
}

type UserPersistence interface {
	PersistUser
	ReadUser
}

type UserPersistenceImpl struct{}

func (up *UserPersistenceImpl) Persist(user *userv1.User) (string, error) {
	return uuid.New().String(), nil
}

func (up *UserPersistenceImpl) Read(uuid string) (*userv1.User, error) {
	return &userv1.User{}, nil
}
