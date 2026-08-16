package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"web-app-uji-emisi/config"
	"web-app-uji-emisi/middleware"
	"web-app-uji-emisi/models"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request"})
		return
	}

	var user models.User
	if err := config.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid username or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid username or password"})
		return
	}

	token, err := middleware.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login successful",
		"data": LoginResponse{
			Token: token,
			User:  user,
		},
	})
}

func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Logout successful"})
}

func GetCurrentUser(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var user models.User
	if err := config.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": user})
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func SeedInitialData() {
	var count int64
	config.DB.Model(&models.User{}).Count(&count)
	if count > 0 {
		return
	}

	users := []models.User{
		{
			ID:        uuid.New().String(),
			Nama:      "Superuser Admin",
			Username:  "adminmas",
			Password:  mustHashPassword("adminmas"),
			Role:      "admin",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.New().String(),
			Nama:      "Kang Supervisor",
			Username:  "supervisoraja",
			Password:  mustHashPassword("supervisoraja"),
			Role:      "supervisor",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.New().String(),
			Nama:      "Kroco01",
			Username:  "kroco01",
			Password:  mustHashPassword("kroco01"),
			Role:      "anggota",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	for _, user := range users {
		config.DB.Create(&user)
	}

	seedParameters()
	seedPopupNotes()
}

func mustHashPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash)
}

func seedParameters() {
	parameters := []models.Parameter{
		// Bensin M
		{ID: uuid.New().String(), Kategori: "Bensin M", TahunMin: 0, TahunMax: intPtr(2007), TahunOperator: "<", COMax: float64Ptr(4.0), HCMax: float64Ptr(1000.0), IsActive: true},
		{ID: uuid.New().String(), Kategori: "Bensin M", TahunMin: 2007, TahunMax: intPtr(2018), TahunOperator: "between", COMax: float64Ptr(1.0), HCMax: float64Ptr(150.0), IsActive: true},
		{ID: uuid.New().String(), Kategori: "Bensin M", TahunMin: 2018, TahunMax: nil, TahunOperator: ">", COMax: float64Ptr(0.5), HCMax: float64Ptr(100.0), IsActive: true},
		// Bensin N&O
		{ID: uuid.New().String(), Kategori: "Bensin N&O", TahunMin: 0, TahunMax: intPtr(2007), TahunOperator: "<", COMax: float64Ptr(4.0), HCMax: float64Ptr(1000.0), IsActive: true},
		{ID: uuid.New().String(), Kategori: "Bensin N&O", TahunMin: 2007, TahunMax: intPtr(2018), TahunOperator: "between", COMax: float64Ptr(1.0), HCMax: float64Ptr(150.0), IsActive: true},
		{ID: uuid.New().String(), Kategori: "Bensin N&O", TahunMin: 2018, TahunMax: nil, TahunOperator: ">", COMax: float64Ptr(0.5), HCMax: float64Ptr(100.0), IsActive: true},
		// Solar JBB
		{ID: uuid.New().String(), Kategori: "Solar JBB", TahunMin: 0, TahunMax: intPtr(2010), TahunOperator: "<", OpasitasMax: float64Ptr(65.0), IsActive: true},
		{ID: uuid.New().String(), Kategori: "Solar JBB", TahunMin: 2010, TahunMax: intPtr(2021), TahunOperator: "between", OpasitasMax: float64Ptr(40.0), IsActive: true},
		{ID: uuid.New().String(), Kategori: "Solar JBB", TahunMin: 2021, TahunMax: nil, TahunOperator: ">", OpasitasMax: float64Ptr(30.0), IsActive: true},
		// Solar GVW
		{ID: uuid.New().String(), Kategori: "Solar GVW", TahunMin: 0, TahunMax: intPtr(2010), TahunOperator: "<", OpasitasMax: float64Ptr(65.0), IsActive: true},
		{ID: uuid.New().String(), Kategori: "Solar GVW", TahunMin: 2010, TahunMax: intPtr(2021), TahunOperator: "between", OpasitasMax: float64Ptr(40.0), IsActive: true},
		{ID: uuid.New().String(), Kategori: "Solar GVW", TahunMin: 2021, TahunMax: nil, TahunOperator: ">", OpasitasMax: float64Ptr(35.0), IsActive: true},
	}

	for _, param := range parameters {
		config.DB.Create(&param)
	}
}

func seedPopupNotes() {
	notes := []models.PopupNote{
		{
			ID:             uuid.New().String(),
			JenisKendaraan: "Bensin M, Bensin N&O",
			TahunOperator:  ">",
			TahunValue:     1990,
			ParameterUji:   "O2",
			NilaiOperator:  ">",
			NilaiValue:     20.0,
			Note:           "Indikasi kebocoran exhaust",
			IsActive:       true,
		},
		{
			ID:             uuid.New().String(),
			JenisKendaraan: "Bensin M",
			TahunOperator:  "<",
			TahunValue:     2018,
			ParameterUji:   "HC",
			NilaiOperator:  ">",
			NilaiValue:     4.0,
			Note:           "Indikasi kerusakan: busi mati/lemah, kabel busi putus, koil pengapian rusak, atau kompresi rendah akibat ring piston aus",
			IsActive:       true,
		},
		{
			ID:             uuid.New().String(),
			JenisKendaraan: "Solar JBB",
			TahunOperator:  ">",
			TahunValue:     2021,
			ParameterUji:   "Opasitas",
			NilaiOperator:  ">",
			NilaiValue:     70.0,
			Note:           "Indikasi Boros BBM",
			IsActive:       true,
		},
	}

	for _, note := range notes {
		config.DB.Create(&note)
	}
}

func intPtr(i int) *int {
	return &i
}

func float64Ptr(f float64) *float64 {
	return &f
}
