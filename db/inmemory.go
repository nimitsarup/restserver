package db

import (
	"errors"

	"github.com/nimitsarup/restserver/model"
	"golang.org/x/crypto/bcrypt"
)

var (
	db                       *InMemoryDB
	ErrUserDoesNotExist      = errors.New("user does not exist")
	ErrDatabaseUninitialized = errors.New("database not initialized")
	ErrUserAlreadyExist      = errors.New("user already exist")
)

type InMemoryDB struct {
	users map[string]model.User
}

type UsersInMemoryDB interface {
	AddUser(user model.User) error
	GetUser(id string) (model.User, error)
	GetAllUsers() ([]model.User, error)
}

func New() (m *InMemoryDB) {
	if db != nil {
		return db
	}
	return &InMemoryDB{users: map[string]model.User{}}
}

func (d *InMemoryDB) AddUser(user model.User) error {
	// check unique password (in real DB, add a unique constraint)
	for _, v := range d.users {
		if v.Email == user.Email {
			return ErrUserAlreadyExist
		}
	}
	// store bcrypted password
	user.Password = HashPassword(user.Password)
	d.users[user.Id] = user
	return nil
}

func (d *InMemoryDB) GetUser(id string) (model.User, error) {
	if u, ok := d.users[id]; ok {
		return u, nil
	}
	return model.User{}, ErrUserDoesNotExist
}

func (d *InMemoryDB) GetAllUsers() ([]model.User, error) {
	allUsers := []model.User{}
	for _, v := range d.users {
		allUsers = append(allUsers, v)
	}
	return allUsers, nil
}

func HashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
