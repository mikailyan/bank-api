package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/therealadik/bank-api/internal/config"
	"github.com/therealadik/bank-api/internal/db"
	"github.com/therealadik/bank-api/internal/handler"
	"github.com/therealadik/bank-api/internal/middleware"
	"github.com/therealadik/bank-api/internal/repository"
	"github.com/therealadik/bank-api/internal/service"
)

func runMigrations(dsn string) {
	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		logrus.Fatalf("Ошибка миграций : %v", err)
	}

	err = m.Up()

	switch {
	case errors.Is(err, migrate.ErrNoChange):
		logrus.Info("Миграции не требуются, схема в актуальном состоянии")
		return

	case err != nil:
		logrus.Fatalf("Ошибка при применении миграций: %v", err)
	}

	logrus.Info("Миграции успешно применены")
}

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)

	ctx := context.Background()
	dbCfg := config.LoadDB()
	jwtCfg := config.LoadJWT()
	cryptoCfg := config.LoadCrypto()

	dsn := db.BuildDSN(dbCfg)
	runMigrations(dsn)

	pool, err := db.New(ctx, dbCfg)
	if err != nil {
		logger.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer pool.Close()
	logger.Info("Подключение к БД успешно установлено")

	userRepo := repository.NewUserRepository(pool)
	accountRepo := repository.NewAccountRepository(pool)
	transactionRepo := repository.NewTransactionRepository(pool)
	cardRepo := repository.NewCardRepository(pool)

	authService := service.NewAuthService(userRepo, jwtCfg)
	accountService := service.NewAccountService(accountRepo, transactionRepo)
	cardService := service.NewCardService(cardRepo, pool, cryptoCfg.HMACKey)

	authHandler := handler.NewAuthHandler(authService, logger)
	accountHandler := handler.NewAccountHandler(accountService, logger)
	cardHandler := handler.NewCardHandler(cardService, logger)

	jwtMiddleware := middleware.NewJWTMiddleware(authService, logger)

	r := mux.NewRouter().PathPrefix("/api").Subrouter()

	r.HandleFunc("/register", authHandler.Register).Methods(http.MethodPost)
	r.HandleFunc("/login", authHandler.Login).Methods(http.MethodPost)

	apiRouter := r.PathPrefix("").Subrouter()
	apiRouter.Use(jwtMiddleware.Middleware)

	apiRouter.HandleFunc("/accounts", accountHandler.CreateAccount).Methods(http.MethodPost)
	apiRouter.HandleFunc("/accounts", accountHandler.GetAccounts).Methods(http.MethodGet)
	apiRouter.HandleFunc("/accounts/{id}/balance", accountHandler.UpdateBalance).Methods(http.MethodPatch)
	apiRouter.HandleFunc("/accounts/{id}/transactions", accountHandler.GetTransactions).Methods(http.MethodGet)
	apiRouter.HandleFunc("/transfer", accountHandler.Transfer).Methods(http.MethodPost)

	apiRouter.HandleFunc("/cards", cardHandler.CreateCard).Methods(http.MethodPost)
	apiRouter.HandleFunc("/cards", cardHandler.GetCards).Methods(http.MethodGet)
	apiRouter.HandleFunc("/cards/{id}", cardHandler.GetCardDetails).Methods(http.MethodGet)
	apiRouter.HandleFunc("/payments", cardHandler.ProcessPayment).Methods(http.MethodPost)

	// Настройка сервера
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", "8080"),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Запуск сервера в горутине
	go func() {
		logger.Infof("Сервер запущен на порту %s", "8080")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalf("Ошибка запуска сервера: %v", err)
		}
	}()

	// Канал для сигналов завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Завершение работы сервера...")

	// Ожидание завершения текущих запросов
	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxShutdown); err != nil {
		logger.Fatalf("Ошибка при остановке сервера: %v", err)
	}
	logger.Info("Сервер успешно остановлен")
}
