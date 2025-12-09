package main

import (
	"database/sql"
	"net/http"

	"github.com/f044fs3t5w3f/gophermart/internal/handler"
	"github.com/f044fs3t5w3f/gophermart/internal/logger"
	dbRepo "github.com/f044fs3t5w3f/gophermart/internal/repository/db"
	"github.com/f044fs3t5w3f/gophermart/internal/service"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

func main() {
	config := getConfig()
	db, err := sql.Open("pgx", config.databaseURI)

	logger.Initialize("info")

	if err != nil {
		logger.Log.Fatal("couldn't open db connection", zap.Error(err))
	}

	err = migrateDB(db)

	if err != nil && err != migrate.ErrNoChange {
		logger.Log.Fatal("couldn't migrate database", zap.Error(err))
	}

	repository := dbRepo.NewDBRepository(db)
	service := service.NewService(repository)
	router := handler.GetRouter(service, repository)
	logger.Log.Info("Server has been started", zap.String("addr", config.runAddress))
	err = http.ListenAndServe(config.runAddress, router)
	if err != nil {
		logger.Log.Fatal("couldn't start server", zap.Error(err))
	}
}
