package users

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `gorm:"not null" json:"name"`
	Email    string `gorm:"unique;not null" json:"email"`
	Password string `json:"password" gorm:"not null"`
	Role     string `gorm:"type:varchar(20);default:'driver';not null" json:"role"`
}

func (u *User) hashPassword(pass string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return nil
}

func (u *User) checkPassword(pass string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pass))
}
