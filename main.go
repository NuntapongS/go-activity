package main

import (
	"context"
	"fmt"
	"myapp/activity"
	"myapp/login"
	"myapp/mw"

	"myapp/staff"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
		AuthMechanism: "SCRAM-SHA-1",
		AuthSource:    "goecho",
		Username:      "root",
		Password:      "password",
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

	mongodb := client.Database("goecho")

	e := echo.New()
	corsConfig := middleware.CORSConfig{
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}

	e.Use(middleware.CORSWithConfig(corsConfig))

	e.GET("/healths", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	e.POST("/api/login", login.LoginHandler(staff.GetStaffByUsername(mongodb)))

	a := e.Group("/api")
	config := middleware.JWTConfig{
		Claims:     &login.JwtCustomClaims{},
		SigningKey: []byte(viper.GetString("jwt.secret")),
	}

	a.Use(middleware.JWTWithConfig(config))
	onlyAdminRoutes := a.Group("")
	onlyAdminRoutes.Use(mw.OnlyAdmin)

	onlyAdminRoutes.GET("/activity", activity.GetActivityHandler(activity.GetActivity(mongodb)))
	onlyAdminRoutes.GET("/activity/:activityid", activity.GetActivityByActivityIDHandler(activity.GetActivityByActivityID(mongodb)))
	onlyAdminRoutes.POST("/activity", activity.CreateActivityhandler(activity.CreateActivity(mongodb)))
	onlyAdminRoutes.PUT("/activity", activity.UpdateActivityHandler(activity.UpdateActivity(mongodb)))
	onlyAdminRoutes.DELETE("/activity/:activityid", activity.DeleteActivityHandler(activity.DeleteActivity(mongodb)))

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
