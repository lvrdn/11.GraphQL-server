package user

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"shop/pkg/session"

	"github.com/99designs/gqlgen/graphql"
)

type Response map[string]interface{}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		St: NewUserStorage(),
		Sm: session.NewSessionStorage(),
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
		log.Printf("read body error: [%s], path: [%s]\n", err.Error(), r.URL.Path)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	dataFromBody := make(map[string]*User)
	err = json.Unmarshal(body, &dataFromBody)
	if err != nil {
		log.Printf("unmarshal body error: [%s], path: [%s]\n", err.Error(), r.URL.Path)
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
		log.Printf("add user error: [%s], path: [%s]\n", err.Error(), r.URL.Path)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sessionToken, err := uh.Sm.CreateSessionToken(newUserID)
	if err != nil {
		log.Printf("create session error: [%s], path: [%s]\n", err.Error(), r.URL.Path)
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
		log.Printf("marshal response error: [%s], path: [%s]\n", err.Error(), r.URL.Path)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(dataResponse)
}
