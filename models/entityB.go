package models

type EntityB struct {
	Id            int    `json:"id"`
	OwnerName     string `json:"owner_name"`
	OwnerPassword string `json:"owner_password"`
	IsPremium     bool   `json:"is_premium"`
	CreationYear  int    `json:"creation_year"`
}