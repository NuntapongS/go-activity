package main

import (
	"context"
	"fmt"
	"myapp/activity"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
)

func main() {
	initConfig()

	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()

	credential := options.Credential{
		Username: "root",
		Password: "password",
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017").SetAuth(credential))
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		panic(err)
	}

	mongodb := client.Database("golang101")

	e := echo.New()

	e.GET("/healths", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	e.GET("/activity", activity.GetActivityHandler(activity.GetActivity(mongodb)))
	e.GET("/activity/:activityid", activity.GetActivityByActivityIDHandler(activity.GetActivityByActivityID(mongodb)))
	e.POST("/activity", activity.CreateActivityhandler(activity.CreateActivity(mongodb)))

	go func() {
		if err := e.Start(":8000"); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
func initConfig() {
	viper.SetDefault("app.port", "8000")
	viper.SetDefault("mongo.uri", "mongodb://localhost:27017")
	viper.SetDefault("mongo.db", "goecho")
	viper.SetDefault("mongo.user", "root")
	viper.SetDefault("mongo.pass", "password")

	viper.SetDefault("jwt.secret", "secret")

	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")      // optionally look for config in the working directory
	err := viper.ReadInConfig()   // Find and read the config file
	if err != nil {
		zap.L().Warn(fmt.Sprintf("Fatal error config file: %s \n", err)) // Handle errors reading the config file
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}
