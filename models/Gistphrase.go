package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Catchphrase struct {
    ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
    GistTitle    string             `json:"gistTitle,omitempty" bson:"gistTitle,omitempty"`
    GistTopic  string             `json:"gistTopic,omitempty" bson:"gistTopic,omitempty"`
    GistContent string             `json:"gistContent,omitempty" bson:"gistContent,omitempty"`
}