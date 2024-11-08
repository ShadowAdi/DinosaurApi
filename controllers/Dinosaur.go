package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"server/models"
	"time"

	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DinosaurController struct {
	session *mongo.Client
}

func NewDinosaurController(s *mongo.Client) *DinosaurController {
	return &DinosaurController{s}
}

func (uc DinosaurController) GetDinosaur(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	collection := uc.session.Database("Dinosaur").Collection("dinos")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var dino models.Dinosaur
	if err := collection.FindOne(ctx, bson.M{"_id": objId}).Decode(&dino); err != nil {
		http.Error(w, "Dinosaur not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dino)
}

func (uc DinosaurController) CreateDinosaur(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var dino models.Dinosaur
	if err := json.NewDecoder(r.Body).Decode(&dino); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	dino.Id = primitive.NewObjectID()
	collection := uc.session.Database("Dinosaur").Collection("dinos")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, dino)
	if err != nil {
		http.Error(w, "Failed to create dinosaur", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dino)
}

func (uc DinosaurController) DeleteDinosaurs(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	collection := uc.session.Database("Dinosaur").Collection("dinos")
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	deleteResult, err := collection.DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil {
		http.Error(w, "Failed to delete dinosaur", http.StatusInternalServerError)
		return
	}

	if deleteResult.DeletedCount == 0 {
		http.Error(w, "Dinosaur not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Deleted Dinosaur with ID: %s\n", id)
}

func (uc DinosaurController) GetDinosaurs(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	collection := uc.session.Database("Dinosaur").Collection("dinos")
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		http.Error(w, "Failed to retrieve dinosaurs", http.StatusInternalServerError)
		return
	}
	defer cur.Close(ctx)

	var dinos []models.Dinosaur
	for cur.Next(ctx) {
		var dino models.Dinosaur
		if err := cur.Decode(&dino); err != nil {
			http.Error(w, "Error decoding dinosaur data", http.StatusInternalServerError)
			return
		}
		dinos = append(dinos, dino)
	}

	if err := cur.Err(); err != nil {
		http.Error(w, "Error reading dinosaurs data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dinos)
}

func (uc DinosaurController) CreateDinosaurs(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var dinos []models.Dinosaur
	if err := json.NewDecoder(r.Body).Decode(&dinos); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	collection := uc.session.Database("Dinosaur").Collection("dinos")
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	for _, dinosaur := range dinos {
		dinosaur.Id = primitive.NewObjectID()
		_, err := collection.InsertOne(ctx, dinosaur)
		if err != nil {
			http.Error(w, "Failed to save dinosaur data", http.StatusInternalServerError)
			log.Println("Insert error:", err)
			return
		}

	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Dinosaurs created successfully"))

}

func (uc DinosaurController) UpdateDinosaur(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var dino models.Dinosaur
	if err := json.NewDecoder(r.Body).Decode(&dino); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	collection := uc.session.Database("Dinosaur").Collection("dinos")
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	update := bson.M{"$set": dino}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	result := collection.FindOneAndUpdate(ctx, bson.M{"_id": objId}, update, opts)

	var updatedDino models.Dinosaur
	if err := result.Decode(&updatedDino); err != nil {
		http.Error(w, "Failed to update dinosaur", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedDino)
}
