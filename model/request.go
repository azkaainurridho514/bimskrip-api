package model

type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6"`
	Name      string `json:"name" validate:"required"`
	Phone     string `json:"phone" validate:"required"`
	RoleID    int    `json:"role_id" validate:"required"`
	DosenPAID *int   `json:"dosen_pa_id,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
	Data       User   `json:"data"`
}

type ProgressRequest struct {
	MhsID     int    `json:"mhs_id" validate:"required"`
	NameID    int    `json:"name_id" validate:"required"`
	DosenPAID int    `json:"dosen_pa_id" validate:"required"`
	Desc      string `json:"desc" validate:"required"`
	Url       string `json:"url" validate:"required"`
}

type ProgressUpdateRequest struct {
	NameID int    `json:"name_id" validate:"required"`
	Url    string `json:"url" validate:"required"`
}

type ProgressStatusUpdate struct {
	StatusID  int `json:"status_id" validate:"required"`
	DosenPAID int `json:"dosen_pa_id" validate:"required"`
}

type ScheduleRequest struct {
	Date      string `json:"date" validate:"required"`
	MhsID     int    `json:"mhs_id" validate:"required"`
	DosenPAID int    `json:"dosen_pa_id" validate:"required"`
	Tempat    string `json:"tempat" validate:"required"`
}

type GetProgressRequest struct {
	UserID int `json:"user_id" validate:"required"`
	RoleID int `json:"role_id" validate:"required"`
}

type GetScheduleRequest struct {
	UserID int `json:"user_id" validate:"required"`
	RoleID int `json:"role_id" validate:"required"`
}