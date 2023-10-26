package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ShortenedURL struct {
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	OriginalURL    string             `json:"original_url,omitempty" bson:"original_url,omitempty" validate:"required,url"`
	CreationDate   primitive.DateTime `json:"creation_date,omitempty" bson:"creation_date,omitempty" validate:"required,datetime"`
	ExpirationDate primitive.DateTime `json:"expiration_date,omitempty" bson:"expiration_date,omitempty" validate:"required,datetime"`
}
