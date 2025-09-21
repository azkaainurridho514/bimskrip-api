package model

type StatusName struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
}