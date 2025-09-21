package model

import "time"
type User struct {
	ID         int   `json:"id" gorm:"primaryKey"`
	RoleID     int   `json:"role_id"`
	DosenPAID  int  `json:"dosen_pa_id,omitempty"`
	Name       string `json:"name"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Photo   	string `json:"photo"`
	Role       Role   `json:"role" gorm:"foreignKey:RoleID"`
	DosenPA    *User  `json:"dosen_pa,omitempty" gorm:"foreignKey:DosenPAID"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
