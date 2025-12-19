package models

type EntityB struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	MaxVarcharEx string `json:"max_varchar_ex"`
	BoolEx       bool   `json:"bool_ex"`
	IntEx        int    `json:"int_ex"`
}