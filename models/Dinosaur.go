package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Dinosaur struct {
	Id                   primitive.ObjectID `bson:"_id,omitempty"`
	Img                  string             `json:"img" bson:"img"`
	DinoName             string             `json:"dinoName" bson:"dinoName"`
	DinoSmallDescription string             `json:"dinoSmallDescription" bson:"dinoSmallDescription"`
	Description          string             `json:"description" bson:"description"`
	Paleontologists      []string           `json:"paleontologists" bson:"paleontologists"`
	LengthInMeter        string             `json:"length_in_meter" bson:"length_in_meter"`
	WeightInKg           string             `json:"weight_in_Kg" bson:"weight_in_Kg"`
	Diet                 string             `json:"diet" bson:"diet"`
	Family               string             `json:"family" bson:"family"`
	MYA                  string             `json:"mya" bson:"mya"`
	Epoch                string             `json:"epoch" bson:"epoch"`
	Age                  string             `json:"age" bson:"age"`
	YearDescribed        string             `json:"year_described" bson:"year_described"`
	YearDiscovered       string             `json:"year_discovered" bson:"year_discovered"`
	DiscoverRegions      []string           `json:"discover_regions" bson:"discover_regions"`
}
