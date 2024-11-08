package main

import (
	"context"
	"log"
	"net/http"
	"server/controllers"
	"time"

	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	r := httprouter.New()
	client, err := getSession()
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	uc := controllers.NewDinosaurController(client)
	r.POST("/dinosaur", uc.CreateDinosaur)
	r.GET("/dinosaurs", uc.GetDinosaurs)
	r.DELETE("/dinosaur/:id", uc.DeleteDinosaurs)
	r.GET("/dinosaur/:id", uc.GetDinosaur)
	r.PUT("/dinosaur/:id", uc.UpdateDinosaur)
	r.POST("/dinosaursAll", uc.CreateDinosaurs)

	http.ListenAndServe("localhost:9001", r)
}

func getSession() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI("mongodb+srv://AdityaShukla45:aRsHGqtwtYxeUTQm@cluster0.nuful.mongodb.net/Dinosaur?retryWrites=true&w=majority")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	if err = client.Ping(ctx, nil); err != nil {
		return nil, err
	}
	log.Println("Connected to MongoDB!")
	return client, nil
}
