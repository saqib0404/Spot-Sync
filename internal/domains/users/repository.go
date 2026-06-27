package users

import (
	"errors"

	"gorm.io/gorm"
)

var ErrorUserAlreadyExists = errors.New("User with this email already exists")
var ErrorInvalidCredentials = errors.New("Invalid email or password")

type Repository interface {
	CreateUser(user *User) error
	GetUserByEmail(email string) (*User, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) CreateUser(user *User) error {
	result := r.db.Create(user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return ErrorUserAlreadyExists
		}
		return result.Error
	}
	return nil
}

func (r *repository) GetUserByEmail(email string) (*User, error) {
	user := &User{}
	result := r.db.Where(&User{Email: email}).First(user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("User not found")
		}
		return nil, result.Error
	}
	return user, nil
}
