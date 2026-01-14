package main

import (
	"base-service/internal/config"
	"base-service/internal/controller"
	"base-service/internal/repository"
	"base-service/internal/service"
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Load()

	db, err := sql.Open("postgres", cfg.PostgresDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	router := mux.NewRouter()
	userController.RegisterRoutes(router)

	log.Printf("Server started on port %d", cfg.HTTPPort)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(cfg.HTTPPort), router))
}
