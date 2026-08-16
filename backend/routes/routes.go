package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"web-app-uji-emisi/config"
	"web-app-uji-emisi/controllers"
	"web-app-uji-emisi/middleware"
)

func SetupRoutes() *gin.Engine {
	r := gin.Default()

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Serve static files
	r.Static("/uploads", "./uploads")

	// API routes
	api := r.Group("/api")
	{
		// Auth routes (public)
		auth := api.Group("/auth")
		{
			auth.POST("/login", controllers.Login)
			auth.POST("/logout", middleware.AuthMiddleware(), controllers.Logout)
			auth.GET("/me", middleware.AuthMiddleware(), controllers.GetCurrentUser)
		}

		// Home content (public for GET, protected for PUT/POST/DELETE)
		home := api.Group("/home")
		{
			home.GET("", controllers.GetHomeContent)
			home.PUT("/description/:id", middleware.AuthMiddleware(), middleware.AdminOrSupervisorMiddleware(), controllers.UpdateDescription)
			home.POST("/image", middleware.AuthMiddleware(), middleware.AdminOrSupervisorMiddleware(), controllers.UploadImage)
			home.DELETE("/image/:id", middleware.AuthMiddleware(), middleware.AdminOrSupervisorMiddleware(), controllers.DeleteImage)
			home.POST("/working-instruction", middleware.AuthMiddleware(), middleware.AdminOrSupervisorMiddleware(), controllers.UploadWorkingInstruction)
			home.DELETE("/working-instruction/:id", middleware.AuthMiddleware(), middleware.AdminOrSupervisorMiddleware(), controllers.DeleteWorkingInstruction)
		}

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// Kendaraan routes
			kendaraan := protected.Group("/kendaraan")
			{
				kendaraan.GET("", controllers.GetKendaraan)
				kendaraan.GET("/all", controllers.GetAllKendaraan)
				kendaraan.GET("/search", controllers.SearchKendaraan)
				kendaraan.GET("/:id", controllers.GetKendaraanByID)
				kendaraan.POST("", controllers.CreateKendaraan)
				kendaraan.PUT("/:id", controllers.UpdateKendaraan)
				kendaraan.DELETE("/:id", controllers.DeleteKendaraan)
				kendaraan.POST("/upload", controllers.UploadCSV)
				kendaraan.GET("/template/csv", controllers.DownloadCSVTemplate)
			}

			// Hasil uji routes
			hasilUji := protected.Group("/hasil-uji")
			{
				hasilUji.GET("", controllers.GetHasilUji)
				hasilUji.GET("/statistics", controllers.GetStatistics)
				hasilUji.POST("", controllers.CreateHasilUji)
				hasilUji.PUT("/:id", controllers.UpdateHasilUji)
				hasilUji.DELETE("/:id", controllers.DeleteHasilUji)
				hasilUji.GET("/export/csv", controllers.ExportCSV)
			}

			// Parameters routes (admin/supervisor only)
			parameters := protected.Group("/parameters")
			parameters.Use(middleware.AdminOrSupervisorMiddleware())
			{
				parameters.GET("", controllers.GetParameters)
				parameters.POST("", controllers.CreateParameter)
				parameters.PUT("/:id", controllers.UpdateParameter)
				parameters.DELETE("/:id", controllers.DeleteParameter)
			}

			// Popup notes routes (admin/supervisor only)
			popupNotes := protected.Group("/popup-notes")
			popupNotes.Use(middleware.AdminOrSupervisorMiddleware())
			{
				popupNotes.GET("", controllers.GetPopupNotes)
				popupNotes.POST("", controllers.CreatePopupNote)
				popupNotes.PUT("/:id", controllers.UpdatePopupNote)
				popupNotes.DELETE("/:id", controllers.DeletePopupNote)
			}

			// Users routes (admin/supervisor only)
			users := protected.Group("/users")
			users.Use(middleware.AdminOrSupervisorMiddleware())
			{
				users.GET("", controllers.GetUsers)
				users.POST("", controllers.CreateUser)
				users.PUT("/:id", controllers.UpdateUser)
				users.DELETE("/:id", controllers.DeleteUser)
			}
		}
	}

	return r
}
