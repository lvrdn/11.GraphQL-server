package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/99designs/gqlgen/graphql"
)

type ctxKey int

const sessionCtxKey ctxKey = 1

type SessionData struct {
	UserID  int
	Session string
}

type SessionStorage struct {
	Sessions []*SessionData
}

func NewSessionStorage() *SessionStorage {
	return &SessionStorage{
		Sessions: make([]*SessionData, 0),
	}
}

func (ss *SessionStorage) GetUserIdFromCtx(ctx context.Context) (int, error) {
	sessionData, ok := ctx.Value(sessionCtxKey).(*SessionData)
	if !ok {
		return 0, fmt.Errorf("User not authorized")
	}

	return sessionData.UserID, nil
}

func (ss *SessionStorage) CheckSession(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	sessionData, ok := ctx.Value(sessionCtxKey).(*SessionData)
	if !ok {
		return nil, fmt.Errorf("User not authorized")
	}

	for _, sessionDataFromStorage := range ss.Sessions {

		if *sessionDataFromStorage == *sessionData {
			return next(ctx)
		}
	}
	return nil, fmt.Errorf("User not authorized")
}

func (ss *SessionStorage) CreateSessionToken(userID int) (string, error) {
	session := "some session"
	ss.Sessions = append(ss.Sessions,
		&SessionData{
			UserID:  userID,
			Session: session,
		})
	return session, nil
}

func (ss *SessionStorage) GetUserID(session string) (int, error) {
	for _, sessionFromStorage := range ss.Sessions {
		if sessionFromStorage.Session == session {
			return sessionFromStorage.UserID, nil
		}
	}
	return 0, fmt.Errorf("no user with this session")
}

func (ss *SessionStorage) AuthMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionFromReq := r.Header.Get("Authorization")

		if sessionFromReq != "" {
			sessionFromReq = strings.TrimPrefix(sessionFromReq, "Token ")
			userID, err := ss.GetUserID(sessionFromReq)
			if err != nil {
				fmt.Println("error getting user id with session - auth middle ware:", err)
				return
			}
			session := &SessionData{
				Session: sessionFromReq,
				UserID:  userID,
			}
			ctx := context.WithValue(r.Context(), sessionCtxKey, session)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		next.ServeHTTP(w, r)
	})
}
