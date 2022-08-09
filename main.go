package main

import (
	"context"
	"myapp/user"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
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

	mongodb := client.Database("gouser")

	e := echo.New()

	e.GET("/user", user.GetUserHandler(user.GetUser(mongodb)))
	e.GET("/user/:UserID", user.GetUserByIDHandler(user.GetUserByID(mongodb)))
	e.POST("/user", user.CreateUserHandler(user.CreateUser(mongodb)))

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
