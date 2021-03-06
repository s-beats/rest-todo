package server

import (
	"context"
	"net/http"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/s-beats/rest-todo/domain/repository"
	"github.com/s-beats/rest-todo/infra/rdb"
	"github.com/s-beats/rest-todo/usecase"
)

func Start() error {
	address := os.Getenv("HOST")
	port := os.Getenv("PORT")
	return http.ListenAndServe(address+":"+port, nil)
}

func UsecaseMiddlewarefunc(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db, err := rdb.NewDB()
		if err != nil {
			log.Fatal().Err(err)
		}

		taskRepo := repository.NewTask(db)
		userRepo := repository.NewUser(db)
		taskUsecase := usecase.NewTask(taskRepo, userRepo)
		userUsecase := usecase.NewUser(userRepo)

		reqWithContext := func(r *http.Request) func(key string, val any) *http.Request {
			return func(key string, val any) *http.Request {
				return r.WithContext(context.WithValue(r.Context(), key, val))
			}
		}(r)
		r = reqWithContext("taskUsecase", taskUsecase)
		r = reqWithContext("userUsecase", userUsecase)

		f(w, r)
		log.Printf("[%v] ", r.Method)
	}
}
