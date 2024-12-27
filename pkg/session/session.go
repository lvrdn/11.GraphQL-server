package session

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
)

type ctxKey int

const sessionCtxKey ctxKey = 1

var errNoUserID error = fmt.Errorf("no user with this session key")

type Session struct {
	UserID int
	Key    string
}

type memoryStorage struct {
	Sessions []*Session
	mu       *sync.RWMutex
}

func NewSessionStorage() *memoryStorage {
	return &memoryStorage{
		Sessions: make([]*Session, 0),
		mu:       &sync.RWMutex{},
	}
}

func (ms *memoryStorage) GetUserIdFromCtx(ctx context.Context) (int, error) {
	session, ok := ctx.Value(sessionCtxKey).(*Session)
	if !ok {
		return 0, fmt.Errorf("user not authorized")
	}

	return session.UserID, nil
}

func (ms *memoryStorage) CheckSession(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	session, ok := ctx.Value(sessionCtxKey).(*Session)
	if !ok {
		return nil, fmt.Errorf("user not authorized")
	}

	ms.mu.RLock()
	defer ms.mu.RUnlock()
	for _, sessionFromStorage := range ms.Sessions {
		if *sessionFromStorage == *session {
			return next(ctx)
		}
	}
	return nil, fmt.Errorf("user not authorized")
}

func (ms *memoryStorage) CreateSessionToken(userID int) (string, error) {
	key := uuid.New().String()
	ms.mu.Lock()
	ms.Sessions = append(ms.Sessions,
		&Session{
			UserID: userID,
			Key:    key,
		})
	ms.mu.Unlock()
	return key, nil
}

func (ms *memoryStorage) GetUserID(key string) (int, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	for _, sessionFromStorage := range ms.Sessions {
		if sessionFromStorage.Key == key {
			return sessionFromStorage.UserID, nil
		}
	}
	return 0, errNoUserID
}

func (ms *memoryStorage) AuthMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		keyFromReq := r.Header.Get("Authorization")

		if keyFromReq != "" {
			keyFromReq = strings.TrimPrefix(keyFromReq, "Token ")
			userID, err := ms.GetUserID(keyFromReq)
			if err != nil {
				if err == errNoUserID {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				log.Printf("get user id error in middleware: [%s]", err.Error())
				return
			}
			session := &Session{
				Key:    keyFromReq,
				UserID: userID,
			}
			ctx := context.WithValue(r.Context(), sessionCtxKey, session)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		next.ServeHTTP(w, r)
	})
}
