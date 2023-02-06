package main

import (
	"UserManager/controllers"
	"UserManager/services"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
)

var (
	server         *gin.Engine
	userservice    services.UserService
	usercontoller  controllers.UserController
	ctx            context.Context
	usercollection *mongo.Collection
	mongoclient    *mongo.Client
	err            error
)

func init() {
	ctx = context.TODO()
	mongoconn := options.Client().ApplyURI("mongodb://localhost:27017")
	mongoclient, err = mongo.Connect(ctx, mongoconn)
	if err != nil {
		log.Fatal(err)
	}
	err = mongoclient.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("mongo connection successful")

	usercollection = mongoclient.Database("userdb").Collection("users")
	userservice = services.NewUserService(usercollection, ctx)
	usercontoller = controllers.NewUserController(userservice)
	server = gin.Default()
}

func main() {
	defer func(mongoclient *mongo.Client, ctx context.Context) {
		err := mongoclient.Disconnect(ctx)
		if err != nil {

		}
	}(mongoclient, ctx)

	basepath := server.Group("/v1")
	usercontoller.RegisterUserRoutes(basepath)

	log.Fatal(server.Run("9090"))
}
