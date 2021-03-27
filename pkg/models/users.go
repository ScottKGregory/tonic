package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	Claims      StandardClaims     `json:"claims"`
	Permissions []string           `json:"permissions"`
}

type StandardClaims struct {
	Subject             string `json:"sub"`
	Name                string `json:"name"`
	GivenName           string `json:"given_name"`
	FamilyName          string `json:"family_name"`
	MiddleName          string `json:"middle_name"`
	Nickname            string `json:"nickname"`
	PreferredUsername   string `json:"preferred_username"`
	Profile             string `json:"profile"`
	Picture             string `json:"picture"`
	Website             string `json:"website"`
	Email               string `json:"email"`
	EmailVerified       bool   `json:"email_verified"`
	Gender              string `json:"gender"`
	Birthdate           string `json:"birthdate"`
	Zoneinfo            string `json:"zoneinfo"`
	Locale              string `json:"locale"`
	PhoneNumber         string `json:"phone_number"`
	PhoneNumberVerified bool   `json:"phone_number_verified"`
	Address             struct {
		Formatted     string `json:"formatted"`
		StreetAddress string `json:"street_address"`
		Locality      string `json:"locality"`
		Region        string `json:"region"`
		PostalCode    string `json:"postal_code"`
		Country       string `json:"country"`
	} `json:"address"`
	UpdatedAt int64 `json:"updated_at"`
}
