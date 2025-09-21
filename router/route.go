package router

import (
	"github.com/azkaainurridho514/bimskrip/handler"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	app.Static("/uploads", "./storage/upload")
	auth := api.Group("/auth")
	auth.Post("/register", handler.Register)
	auth.Get("/login", handler.Login)

	progress := api.Group("/progress")
	progress.Get("/list", handler.GetProgresses)       
	progress.Post("/create", handler.CreateProgress)    
	progress.Delete("/delete/:id", handler.DeleteProgress) 
	progress.Put("/status", handler.UpdateProgressStatus) 

	schedule := api.Group("/schedule")
	schedule.Get("/list", handler.GetSchedules)    
	schedule.Get("/today", handler.GetTodaySchedules)   
	schedule.Post("/create", handler.CreateSchedule) 

	api.Get("/progress-names", handler.GetProgressNames)
	api.Get("/status-names", handler.GetStatusNames)
	api.Get("/users", handler.GetUsersByDosenPA) 
	api.Get("/dosen", handler.GetDosen) 
	api.Get("/mahasiswa", handler.GetMahasiswa) 
	api.Put("/user/update", handler.UpdateUserProfile) 
	api.Get("/dashboard", handler.GetDashboardSummary) 
}
