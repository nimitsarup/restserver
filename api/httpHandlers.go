package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nimitsarup/restserver/handlers"
)

type API struct {
	Handlers handlers.HandlersInterface
}

func (a *API) AddUser(w http.ResponseWriter, r *http.Request) {
	log.Println("invoking AddUser")
	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) <= 0 {
		log.Printf("error/invalid request body %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	status, err := a.Handlers.AddUser(string(body))
	if err != nil {
		log.Printf("caught error %v", err)
	}

	w.WriteHeader(status)
}

func (a *API) GetUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	log.Printf("invoking GetUser with [%s]", id)
	status, user := a.Handlers.GetUser(id)

	w.WriteHeader(status)
	err := writeResponse(w, user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (a *API) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	log.Println("invoking GetAllUsers")
	status, users := a.Handlers.GetAllUsers()
	w.WriteHeader(status)
	err := writeResponse(w, users)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func writeResponse(w http.ResponseWriter, resp interface{}) error {
	b, err := json.Marshal(resp)
	if err != nil {
		log.Printf("caught marshal error %v", err)
		return err
	}
	if _, err = w.Write(b); err != nil {
		log.Printf("caught write error %v", err)
		return err
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return nil
}
