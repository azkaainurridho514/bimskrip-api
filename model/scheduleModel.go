package model

import "time"

type Schedule struct {
	ID        int      `json:"id" gorm:"primaryKey"`
	Date      string    `json:"date"`
	// Date      time.Time    `json:"date"`
	MhsID     int      `json:"mhs_id"`
	DosenPAID int      `json:"dosen_pa_id"`
	Tempat    string    `json:"tempat"`
	Student   User      `json:"student" gorm:"foreignKey:MhsID"`
	DosenPA   User      `json:"dosen_pa" gorm:"foreignKey:DosenPAID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
