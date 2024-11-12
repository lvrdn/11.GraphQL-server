package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/99designs/gqlgen/graphql"
)

type Response map[string]interface{}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		St: NewUserStorage(),
		Sm: NewSessionStorage(),
	}
}

type UserHandler struct {
	St UserStorage
	Sm SessionManager
}

type User struct {
	ID       int
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

type UserStorage interface {
	Add(*User) (int, error)
}
type SessionManager interface {
	CreateSessionToken(int) (string, error)
	AuthMiddleWare(http.Handler) http.Handler
	CheckSession(context.Context, interface{}, graphql.Resolver) (interface{}, error)
	GetUserIdFromCtx(context.Context) (int, error)
}

func (uh *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("error with read r.body", r.URL.Path, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	dataFromBody := make(map[string]*User)
	err = json.Unmarshal(body, &dataFromBody)
	if err != nil {
		fmt.Println("error with unmarshal json from r.body", r.URL.Path, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	newUser, ok := dataFromBody["user"]
	if !ok {
		http.Error(w, "bad request, no user data", http.StatusBadRequest)
		return
	}

	newUserID, err := uh.St.Add(newUser)
	if err != nil {
		fmt.Println("error with add new user to storage", r.URL.Path, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sessionToken, err := uh.Sm.CreateSessionToken(newUserID)
	if err != nil {
		fmt.Println("error with getting session token for new user", r.URL.Path, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := Response{
		"body": Response{
			"token": sessionToken,
		},
	}

	dataResponse, err := json.Marshal(response)
	if err != nil {
		fmt.Println("error with marshal data response", r.URL.Path, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(dataResponse)
}
