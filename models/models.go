package models

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ShortenedURL struct {
	ID             string             `json:"_id,omitempty" bson:"_id,omitempty"`
	OriginalURL    string             `json:"original_url,omitempty" bson:"original_url,omitempty" validate:"required,url"`
	CreationDate   primitive.DateTime `json:"creation_date,omitempty" bson:"creation_date,omitempty" validate:"required,datetime"`
	ExpirationDate primitive.DateTime `json:"expiration_date,omitempty" bson:"expiration_date,omitempty" validate:"datetime"`
}

type URLStats struct {
	ID          string               `json:"_id,omitempty" bson:"_id,omitempty"`
	AccessTimes []primitive.DateTime `json:"access_times,omitempty" bson:"access_times,omitempty"`
}

func (s ShortenedURL) GenerateShortUrl() string {
	return fmt.Sprintf("/go/%s", s.ID)
}
