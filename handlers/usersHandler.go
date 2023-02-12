package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/nimitsarup/restserver/db"
	"github.com/nimitsarup/restserver/model"
	"github.com/valyala/fastjson"
)

type Handlers struct {
	DB db.UsersInMemoryDB
}

//go:generate moq -out mock/HandlersInterface.go -pkg mock . HandlersInterface
type HandlersInterface interface {
	AddUser(body string) (int, error)
	GetUser(id string) (int, *model.User)
	GetAllUsers() (int, []model.User)
}

// returns http-status + error (if any)
func (h *Handlers) AddUser(body string) (int, error) {
	err := fastjson.Validate(body)
	if err != nil {
		log.Printf("invalid message body, err [%v]", err)
		return http.StatusBadRequest, err
	}
	user := model.CreateUserReq{}
	if err := json.Unmarshal([]byte(body), &user); err != nil {
		log.Printf("invalid message body, err [%v]", err)
		return http.StatusBadRequest, err
	}
	usr := model.ToUser(&user)
	usr.Id = uuid.NewString()

	err = h.DB.AddUser(*usr)
	if errors.Is(err, db.ErrUserAlreadyExist) {
		return http.StatusConflict, nil
	} else if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusCreated, nil
}

func (h *Handlers) GetUser(id string) (int, *model.User) {
	usr, err := h.DB.GetUser(id)
	if errors.Is(err, db.ErrUserDoesNotExist) {
		return http.StatusNotFound, nil
	} else if err != nil {
		return http.StatusInternalServerError, nil
	}
	return http.StatusOK, &usr
}

func (h *Handlers) GetAllUsers() (int, []model.User) {
	users, err := h.DB.GetAllUsers()
	if err != nil {
		return http.StatusInternalServerError, users
	}
	return http.StatusOK, users
}
