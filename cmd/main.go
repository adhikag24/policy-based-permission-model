package main

import (
	"fmt"
	"log/slog"

	"github.com/adhikag24/policy-based-permission-model/domain/policies"
	"github.com/adhikag24/policy-based-permission-model/http"
	handlerspolicies "github.com/adhikag24/policy-based-permission-model/http/handlers/policies"
	"github.com/adhikag24/policy-based-permission-model/infrastructure/mysql"
	mysqlpolicies "github.com/adhikag24/policy-based-permission-model/infrastructure/mysql/policies"
	"github.com/adhikag24/policy-based-permission-model/utils"
	"github.com/labstack/echo/v5"
)

func main() {
	utils.LoadEnv()

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
	var (
		mysqlUsername = utils.EnvKey("MYSQL_USERNAME").GetValue()
		mysqlPassword = utils.EnvKey("MYSQL_PASSWORD").GetValue()
		mysqlHost     = utils.EnvKey("MYSQL_HOST").GetValue()
		mysqlDatabase = utils.EnvKey("MYSQL_DATABASE").GetValue()
		mysqlPort     = utils.EnvKey("MYSQL_PORT").GetValue()
	)

	if mysqlHost == "" || mysqlUsername == "" || mysqlPassword == "" || mysqlDatabase == "" {
		panic("missing required MySQL environment variables")
	}

	if mysqlPort == "" {
		mysqlPort = "3306"
	}

	return &Config{
		MySQL: mysql.MySQLConfig{
			Username: mysqlUsername,
			Password: mysqlPassword,
			Host:     mysqlHost,
			Port:     mysqlPort,
			DBName:   mysqlDatabase,
		},
	}
}
