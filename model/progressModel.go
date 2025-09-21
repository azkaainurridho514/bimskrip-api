package model

import "time"

type Progress struct {
	ID             int         `json:"id" gorm:"primaryKey"`
	MhsID          int         `json:"mhs_id"`
	NameID         int         `json:"name_id"`
	StatusID       int         `json:"status_id"`
	DosenPAID      int         `json:"dosen_pa_id"`
	Desc      	   string       `json:"desc"`
	Url      	   string       `json:"url"`
	Comment        string       `json:"comment"`
	Student        User         `json:"student" gorm:"foreignKey:MhsID"`
	ProgressName   ProgressName `json:"progress_name" gorm:"foreignKey:NameID"`
	Status         StatusName   `json:"status" gorm:"foreignKey:StatusID"`
	DosenPA        User         `json:"dosen_pa" gorm:"foreignKey:DosenPAID"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}
