package main

import (
	"context"
	"database/sql"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/f044fs3t5w3f/gophermart/internal/accrual"
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	accrualService := accrual.NewAccuralService(ctx, repository, logger.Log, config.accrualSystemAddress)
	accrualService.LoadOld()
	service := service.NewService(repository, accrualService)
	router := handler.GetRouter(service, repository)
	srv := &http.Server{
		Addr:    config.runAddress,
		Handler: router,
		BaseContext: func(l net.Listener) context.Context {
			return ctx
		},
	}

	go func() {
		logger.Log.Info("Starting http server", zap.String("addr", config.runAddress))
		err = srv.ListenAndServe()

		if err != nil {
			logger.Log.Fatal("couldn't start server", zap.Error(err))
		}

	}()

	sig := <-signals
	logger.Log.Info("shutting down", zap.String("signal", sig.String()))
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	srv.Shutdown(shutdownCtx)
	accrualService.Wait()
}
