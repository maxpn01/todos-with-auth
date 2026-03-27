package handlers

import (
	"net/http"
	"todo-app-with-auth/models"
)

func (h Handler) Signup(w http.ResponseWriter, r *http.Request) (id int64, user models.User) {

}

func (h Handler) Signin(w http.ResponseWriter, r *http.Request) {

}

func (h Handler) Signout(w http.ResponseWriter, r *http.Request) {

}
