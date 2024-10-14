package main

import (
	"context"
	"flag"
	"log"

	"github.com/amradel55/hotel-reservation/api"
	"github.com/amradel55/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config = fiber.Config{
	ErrorHandler: func(ctx *fiber.Ctx, err error) error {
		return ctx.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {

	listenAddr := flag.String("listenAddr", "localhost:5000", "The listen address of the API server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	// handler initialization
	var (
		hs           = db.NewMongoHotelStore(client, db.DBNAME)
		rs           = db.NewMongoRoomStore(client, hs, db.DBNAME)
		userHandler  = api.NewUserHandler(db.NewMongoUserStore(client, db.DBNAME))
		hotelHandler = api.NewHotelHandler(hs, rs)
		app          = fiber.New(config)
		apiv1        = app.Group("/api/v1")
	)

	// users handlers
	apiv1.Put("user/:id", userHandler.HandlePutUser)
	apiv1.Delete("user/:id", userHandler.HandleDeleteUser)
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)

	//hotel handlers
	apiv1.Get("/hotel", hotelHandler.HandleGetHotels)

	app.Listen(*listenAddr)
}
