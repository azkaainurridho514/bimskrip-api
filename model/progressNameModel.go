package model

type ProgressName struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
}