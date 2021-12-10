package repositories

import (
	"errors"
	"github.com/go-pg/pg/v10"
	"market/models"
)

var NotFoundErr = errors.New("not found")

// Users - products repository interface
type Users interface {
	CreateUser(phone string) (models.User, error)
	GetUserByID(id int) (models.User, error)
	GetUserByPhone(phone string) (models.User, error)
}

type users struct {
	db *pg.DB
}

// NewUsers - new users repository
func NewUsers(db *pg.DB) Users {
	return &users{
		db: db,
	}
}

// CreateUser - create new user in database
func (u *users) CreateUser(phone string) (models.User, error) {
	user := models.User{
		Phone: phone,
	}

	_, err := u.db.Model(&user).Returning("*").Insert()

	return user, err
}

// GetUserByID - get user by id from database
func (u *users) GetUserByID(id int) (models.User, error) {
	var user models.User
	err := u.db.Model(&user).Where("id = ?", id).First()

	return user, err
}

// GetUserByPhone - get user by phone from database
func (u *users) GetUserByPhone(phone string) (models.User, error) {
	var user models.User
	err := u.db.Model(&user).Where("phone = ?", phone).First()
	if err == pg.ErrNoRows {
		return user, NotFoundErr
	}

	return user, nil
}
