package handler

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/azkaainurridho514/bimskrip/database"
	"github.com/azkaainurridho514/bimskrip/model"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *fiber.Ctx) error {
	var req model.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(201).JSON(fiber.Map{
			"status_code": 201,
			"data":        []string{},
			"message":     "Invalid request body",
		})
	}

	if req.DosenPAID != nil {
		var count int64
		err := database.GetDB().Model(&model.User{}).
			Where("dosen_pa_id = ?", req.DosenPAID).
			Count(&count).Error
		if err != nil {
			return c.Status(201).JSON(fiber.Map{"status_code": 400, "data": []string{}, "message": "Failed to check DPA usage"})
		}
		if count >= 15 {
			return c.Status(201).JSON(fiber.Map{"status_code": 400, "data": []string{}, "message": "Dosen PA sudah memiliki 15 mahasiswa"})
		}
	}

	var existingUser model.User
	if err := database.GetDB().Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return c.Status(201).JSON(fiber.Map{"status_code": 400, "data": []string{}, "message": "User already exists"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(201).JSON(fiber.Map{"status_code": 400, "data": []string{}, "message": "Failed to hash password"})
	}

	var photoURL string
	file, err := c.FormFile("photo")
	if err == nil {
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".webp" {
			folder := "./storage/upload/photo"
			if err := os.MkdirAll(folder, os.ModePerm); err == nil {
				dst := fmt.Sprintf("%s/%s", folder, file.Filename)
				if err := c.SaveFile(file, dst); err == nil {
					photoURL = fmt.Sprintf("%s://%s/uploads/photo/%s", c.Protocol(), c.Hostname(), file.Filename)
				}
			}
		}
	}

	roleIDStr := c.FormValue("role_id")
	roleID, _ := strconv.Atoi(roleIDStr) 
	dosenPAIDStr := c.FormValue("dosen_pa_id")
	dosenPAID, _ := strconv.Atoi(dosenPAIDStr) 

	user := model.User{
		Name:      req.Name,
		Email:     req.Email,
		Password:  string(hashedPassword),
		Phone:     req.Phone,
		RoleID:    int(roleID),
		DosenPAID: int(dosenPAID),
		Photo:     photoURL,
	}

	if err := database.GetDB().Create(&user).Error; err != nil {
		return c.Status(201).JSON(fiber.Map{"status_code": 400, "data": []string{}, "message": "Failed to create user"})
	}

	database.GetDB().Preload("Role").Preload("DosenPA.Role").First(&user, user.ID)

	return c.Status(201).JSON(fiber.Map{
		"message":     "User registered successfully",
		"status_code": 200,
		"data":        user,
	})
}

func Login(c *fiber.Ctx) error {
	email := c.Query("email")
	password := c.Query("password")

	if email == "" || password == "" {
		return c.Status(201).JSON(fiber.Map{
			"status_code": 400,
			"data":        []string{},
			"message":     "Email and password are required",
		})
	}

	var user model.User
	if err := database.GetDB().
		Preload("Role").
		Preload("DosenPA.Role").
		Where("email = ?", email).
		First(&user).Error; err != nil {
		return c.Status(201).JSON(fiber.Map{
			"status_code": 400,
			"data":        []string{},
			"message":     "Invalid credentials",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return c.Status(201).JSON(fiber.Map{
			"status_code": 400,
			"data":        []string{},
			"message":     "Invalid credentials",
		})
	}

	return c.Status(201).JSON(model.LoginResponse{
		Message:    "Login successful",
		Data:       user,
		StatusCode: 200,
	})
}

func UpdateUserProfile(c *fiber.Ctx) error {
	var req struct {
		ID    int    `json:"id" validate:"required"`
		Name  string `json:"name"`
		Email string `json:"email"`
		Phone string `json:"phone"`
		Photo string `json:"photo"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(201).JSON(fiber.Map{
			"status_code": 400,
			"message":     "Invalid request body",
			"data": []string{},
		})
	}

	var user model.User
	if err := database.GetDB().First(&user, req.ID).Error; err != nil {
		return  c.Status(201).JSON(fiber.Map{
			"status_code": 404,
			"data": []string{},
			"message":     "User not found",
		})
	}
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Photo != "" {
		var photoURL string
		file, err := c.FormFile("photo")
		if err == nil {
			ext := strings.ToLower(filepath.Ext(file.Filename))
			if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".webp" {
				folder := "./storage/upload/photo"
				if err := os.MkdirAll(folder, os.ModePerm); err == nil {
				dst := fmt.Sprintf("%s/%s", folder, file.Filename)
				if err := c.SaveFile(file, dst); err == nil {
					photoURL = fmt.Sprintf("%s://%s/uploads/photo/%s", c.Protocol(), c.Hostname(), file.Filename)
				}
				}
			}
		}
		user.Photo = photoURL
	}
	user.UpdatedAt = time.Now()
	if err := database.GetDB().Save(&user).Error; err != nil {
		return  c.Status(201).JSON(fiber.Map{
			"status_code": 500,
			"data": []string{},
			"message":     "Failed to update user profile",
		})
	}

	return  c.Status(201).JSON(fiber.Map{
		"status_code": 200,
		"message":     "Profile updated successfully",
		"data":        user,
	})
}

func GetProgresses(c *fiber.Ctx) error {
	userID := c.Query("user_id")
	roleID := c.Query("role_id")

	if userID == "" || roleID == "" {
		return c.Status(201).Status(201).JSON(fiber.Map{"status_code": 400,
			"data": []string{},"message": "Missing user_id or role_id in query"})
	}

	var progresses []model.Progress
	query := database.GetDB().
		Preload("Student.Role").
		Preload("ProgressName").
		Preload("Status").
		Preload("DosenPA.Role")

	switch roleID {
	case "1":
		query = query.Where("mhs_id = ?", userID)
	default:
		query = query.Where("dosen_pa_id = ?", userID)
	}

	if err := query.Find(&progresses).Error; err != nil {
		return c.Status(201).Status(201).JSON(fiber.Map{"status_code": 200,
			"data": []string{},"message": "Failed to fetch progresses"})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Progresses retrieved successfully",
		"status_code": 200,
		"data":     progresses,
	})
}

func CreateProgress(c *fiber.Ctx) error {
	var req model.ProgressRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(201).JSON(fiber.Map{"status_code": 400,
			"data": []string{},"message": "Invalid request body"})
	}

	mhsID := c.FormValue("mhs_id")
	nameID := c.FormValue("name_id")
	dosenPaID := c.FormValue("dosen_pa_id")
	desc := c.FormValue("desc")
	mhsIDInt, _ := strconv.Atoi(mhsID)
	nameIDInt, _ := strconv.Atoi(nameID)
	dosenPaIDInt, _ := strconv.Atoi(dosenPaID)

	var student model.User
	if err := database.GetDB().Where("id = ? AND role_id = 1", mhsIDInt).First(&student).Error; err != nil {
		return c.Status(201).JSON(fiber.Map{"status_code": 400,
			"data": []string{},"message": "Student not found"})
	}

	var dosen model.User
	if err := database.GetDB().Where("id = ? AND role_id = 2", dosenPaIDInt).First(&dosen).Error; err != nil {
		return c.Status(201).JSON(fiber.Map{"status_code": 400,
			"data": []string{},"message": "Supervisor not found"})
	}
	var fileURL string
	file, err := c.FormFile("file")
	if err == nil {
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext == ".pdf" || ext == ".doc" || ext == ".docx" {
			folder := "./storage/upload/files"
			if err := os.MkdirAll(folder, os.ModePerm); err == nil {
				dst := fmt.Sprintf("%s/%s", folder, file.Filename)
				if err := c.SaveFile(file, dst); err == nil {
					fileURL = fmt.Sprintf("%s://%s/uploads/files/%s", c.Protocol(), c.Hostname(), file.Filename)
					
				}
			}
		}
	}

	progress := model.Progress{
		MhsID:     mhsIDInt,
		NameID:    nameIDInt,
		StatusID:  1, 
		DosenPAID: dosenPaIDInt,
		Desc:      desc,
		Url:       fileURL,
	}

	if err := database.GetDB().Create(&progress).Error; err != nil {
		return c.Status(201).JSON(fiber.Map{"status_code": 400,
			"data": []string{},"message": "Failed to create progress"})
	}

	database.GetDB().Preload("Student.Role").Preload("ProgressName").Preload("Status").Preload("DosenPA.Role").First(&progress, progress.ID)

	return c.Status(201).JSON(fiber.Map{
		"message": "Progress created successfully",
		"status_code": 200,
		"data":     progress,
	})
}

func DeleteProgress(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(201).JSON(fiber.Map{"status_code": 400,
			"data": []string{},"message": "Invalid ID"})
	}

	var progress model.Progress
	if err := database.GetDB().First(&progress, id).Error; err != nil {
		return c.Status(201).JSON(fiber.Map{"status_code": 400,
			"data": []string{},"message": "Progress not found"})
	}

	if err := database.GetDB().Delete(&progress, id).Error; err != nil {
		return c.Status(201).JSON(fiber.Map{"status_code": 400,
			"data": []string{},"message": "Failed to delete progress"})
	}

	return c.Status(201).JSON(fiber.Map{"status_code": 200,
		"data": []string{},"message": "Progress deleted successfully"})
}

func UpdateProgressStatus(c *fiber.Ctx) error {
	var req struct {
		ID      int `json:"id" validate:"required"`
		StatusID int `json:"status_id" validate:"required"`
		Comment string `json:"comment"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(201).JSON(fiber.Map{
			"status_code": 400,
			"data":        []string{},
			"message":     "Invalid request body",
		})
	}

	var progress model.Progress
	if err := database.GetDB().First(&progress, req.ID).Error; err != nil {
		return c.Status(201).JSON(fiber.Map{
			"status_code": 404,
			"data":        []string{},
			"message":     "Progress not found",
		})
	}

	progress.StatusID = req.StatusID
	progress.Comment = req.Comment
	if err := database.GetDB().Save(&progress).Error; err != nil {
		return c.Status(201).JSON(fiber.Map{
			"status_code": 500,
			"data":        []string{},
			"message":     "Failed to update progress status",
		})
	}

	database.GetDB().
		Preload("Student.Role").
		Preload("ProgressName").
		Preload("Status").
		Preload("DosenPA.Role").
		First(&progress, progress.ID)

	return c.Status(201).JSON(fiber.Map{
		"status_code": 200,
		"message":     "Progress status updated successfully",
		"data":        progress,
	})
}

func GetTodaySchedules(c *fiber.Ctx) error {
	userID := c.Query("user_id")
	roleID := c.Query("role_id")

	if userID == "" || roleID == "" {
		return c.Status(201).JSON(fiber.Map{
			"status_code": 400,
			"data":        []string{},
			"message":     "Missing query parameters: user_id and role_id are required",
		})
	}

	today := time.Now().Format("2006-01-02") 
	fmt.Println(today)

	var schedules []model.Schedule
	query := database.GetDB().Preload("Student.Role").Preload("DosenPA.Role")

	if roleID == "1" {
		query = query.Where("mhs_id = ?", userID)
	} else if roleID == "2" {
		query = query.Where("dosen_pa_id = ?", userID)
	} else {
		return c.Status(201).JSON(fiber.Map{
			"status_code": 400,
			"data":        []string{},
			"message":     "Invalid role_id: must be 1 (Mahasiswa) or 2 (Dosen)",
		})
	}

	query = query.Where("date = ?", today)

	if err := query.Find(&schedules).Error; err != nil {
		return c.Status(201).JSON(fiber.Map{
			"status_code": 500,
			"data":        []string{},
			"message":     "Failed to fetch today's schedules",
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message":     "Today's schedules retrieved successfully",
		"status_code": 200,
		"data":        schedules,
	})
}

func GetSchedules(c *fiber.Ctx) error {
	userID := c.Query("user_id")
	roleID := c.Query("role_id")

	if userID == "" || roleID == "" {
		return c.Status(201).JSON(fiber.Map{
			"status_code": 400,
			"data":        []string{},
			"message":     "Missing query parameters: user_id and role_id are required",
		})
	}

	var schedules []model.Schedule
	query := database.GetDB().Preload("Student.Role").Preload("DosenPA.Role")

	if roleID == "1" {
		query = query.Where("mhs_id = ?", userID)
	} else {
		query = query.Where("dosen_pa_id = ?", userID)
	}

	if err := query.Find(&schedules).Error; err != nil {
		return c.Status(201).JSON(fiber.Map{
			"status_code": 500,
			"data":        []string{},
			"message":     "Failed to fetch schedules",
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message":     "Schedules retrieved successfully",
		"status_code": 200,
		"data":        schedules,
	})
}

func CreateSchedule(c *fiber.Ctx) error {
	var req model.ScheduleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(201).JSON(fiber.Map{"status_code": 400,
			"data": []string{},"message": "Invalid request body"})
	}

	var student model.User
	if err := database.GetDB().Where("id = ? AND role_id = 1", req.MhsID).First(&student).Error; err != nil {
		return c.Status(201).JSON(fiber.Map{"status_code": 400,
			"data": []string{},"message": "Student not found"})
	}

	var dosen model.User
	if err := database.GetDB().Where("id = ? AND role_id = 2", req.DosenPAID).First(&dosen).Error; err != nil {
		return c.Status(201).JSON(fiber.Map{"status_code": 400,
			"data": []string{},"message": "Supervisor not found"})
	}

	schedule := model.Schedule{
		Date:      req.Date,
		MhsID:     req.MhsID,
		DosenPAID: req.DosenPAID,
		Tempat:    req.Tempat,
	}

	if err := database.GetDB().Create(&schedule).Error; err != nil {
		return c.Status(201).JSON(fiber.Map{"status_code": 400,
			"data": []string{},"message": "Failed to create schedule"})
	}

	database.GetDB().Preload("Student.Role").Preload("DosenPA.Role").First(&schedule, schedule.ID)

	return c.Status(201).JSON(fiber.Map{
		"message": "Schedule created successfully",
		"status_code": 200,
		"data":    schedule,
	})
}

func GetProgressNames(c *fiber.Ctx) error {
	var progressNames []model.ProgressName
	if err := database.GetDB().Find(&progressNames).Error; err != nil {
		return c.Status(201).JSON(fiber.Map{"status_code": 400,
			"data": []string{},"message": "Failed to fetch progress names"})
	}
	return c.Status(201).JSON(fiber.Map{
		"message": "Progress names retrieved successfully",
		"status_code": 200,
		"data":    progressNames,
	})
}

func GetStatusNames(c *fiber.Ctx) error {
	var statusNames []model.StatusName
	if err := database.GetDB().Find(&statusNames).Error; err != nil {
		return c.Status(201).JSON(fiber.Map{"status_code": 400,
			"data": []string{},"message": "Failed to fetch status names"})
	}
	return c.Status(201).JSON(fiber.Map{
		"message": "Status names retrieved successfully",
		"status_code": 200,
		"data":     statusNames,
	})
}

func GetUsersByDosenPA(c *fiber.Ctx) error {
	dosenPaID := c.Query("dosen_pa_id")
	if dosenPaID == "" {
		return c.Status(201).JSON(fiber.Map{
			"status_code": 400,
			"data":        []string{},
			"message":     "Missing dosen_pa_id query parameter",
		})
	}

	var users []model.User
	if err := database.GetDB().
		Preload("Role").
		Where("dosen_pa_id = ?", dosenPaID).
		Find(&users).Error; err != nil {
		return c.Status(201).JSON(fiber.Map{
			"status_code": 500,
			"data":        []string{},
			"message":     "Failed to fetch users",
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"status_code": 200,
		"message":     "Users retrieved successfully",
		"data":        users,
	})
}

func GetMahasiswa(c *fiber.Ctx) error {
	var mahasiswa []model.User
	if err := database.GetDB().Preload("Role").Where("role_id = 1").Find(&mahasiswa).Error; err != nil {
		return c.Status(201).JSON(fiber.Map{"status_code": 400,
			"data": []string{},"message": "Failed to fetch mahasiswa"})
	}
	return c.Status(201).JSON(fiber.Map{
		"message": "Mahasiswa retrieved successfully",
		"status_code": 200,
		"data":    mahasiswa,
	})
}

func GetDosen(c *fiber.Ctx) error {
	var dosen []model.User
	if err := database.GetDB().
		Preload("Role").
		Where("role_id = 2").
		Find(&dosen).Error; err != nil {
		return c.Status(201).JSON(fiber.Map{
			"status_code": 400,
			"data":        []string{},
			"message":     "Failed to fetch dosen",
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message":     "Dosen retrieved successfully",
		"status_code": 200,
		"data":        dosen,
	})
}

func GetDashboardSummary(c *fiber.Ctx) error {
		var totalMahasiswa int64
		var totalDosen int64
		var totalProgress int64
		var progressTerima int64
		var progressTolak int64
		var progressPending int64

		database.GetDB().Model(&model.User{}).Where("role_id = ?", 1).Count(&totalMahasiswa)
		database.GetDB().Model(&model.User{}).Where("role_id = ?", 2).Count(&totalDosen)
		database.GetDB().Model(&model.Progress{}).Count(&totalProgress)
		database.GetDB().Model(&model.Progress{}).Where("status_id = ?", 1).Count(&progressPending)
		database.GetDB().Model(&model.Progress{}).Where("status_id = ?", 2).Count(&progressTerima)
		database.GetDB().Model(&model.Progress{}).Where("status_id NOT IN ?", []int{1, 2}).Count(&progressTolak)

		return c.Status(201).JSON(fiber.Map{
			"message": "Dashboard summary retrieved successfully",
			"status_code": 200,
			"data": fiber.Map{
				"users": fiber.Map{
					"total_mahasiswa": totalMahasiswa,
					"total_dosen":     totalDosen,
				},
				"progress": fiber.Map{
					"total":   totalProgress,
					"terima":  progressTerima,
					"tolak":   progressTolak,
					"pending": progressPending,
				},
			},
		})
	
}