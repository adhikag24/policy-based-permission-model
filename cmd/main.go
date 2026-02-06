package main

import (
	"fmt"
	"log/slog"

	"github.com/adhikag24/policy-based-permission-model/domain/policies"
	"github.com/adhikag24/policy-based-permission-model/http"
	handlerspolicies "github.com/adhikag24/policy-based-permission-model/http/handlers/policies"
	"github.com/adhikag24/policy-based-permission-model/infrastructure/mysql"
	mysqlpolicies "github.com/adhikag24/policy-based-permission-model/infrastructure/mysql/policies"
	"github.com/labstack/echo/v5"
)

func main() {
	e := echo.New()

	config := initializeConfig()

	db, err := mysql.Connect(config.MySQL)
	if err != nil {
		fmt.Println("err", err)
		panic("failed to connect database")
	}

	policiesRepository := mysqlpolicies.NewRepository(db)
	policiesService := policies.NewService(policiesRepository)
	policiesHandler := handlerspolicies.NewHandler(policiesService)

	http.RegisterRoutes(e, &http.Handlers{
		Policies: policiesHandler,
	})

	slog.Info("starting server on :8080")
	e.Start(":8080")
}

type Config struct {
	MySQL mysql.MySQLConfig
}

func initializeConfig() *Config {
	return &Config{
		MySQL: mysql.MySQLConfig{
			Username: "root",
			Password: "admin",
			Host:     "127.0.0.1",
			Port:     3306,
			DBName:   "policy-based-permission-model",
		},
	}
}
