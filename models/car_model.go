package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Car struct {
	Id           primitive.ObjectID `json:"id,omitempty"`
	Model        string             `json:"model,omitempty" validate:"required"`
	Manufacturer string             `json:"manufacturer,omitempty" validate:"required"`
	Year         int                `json:"year,omitempty" validate:"required"`
	Kilometres   float64            `json:"kilometres,omitempty" validate:"required"`
}
