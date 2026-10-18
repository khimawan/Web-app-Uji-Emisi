package controllers

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"web-app-uji-emisi/config"
	"web-app-uji-emisi/models"
)

type UpdateDescriptionRequest struct {
	Title       string `json:"title"`
	Description string `json:"description" binding:"required"`
}

func GetHomeContent(c *gin.Context) {
	var descriptions []models.HomeContent
	config.DB.Where("content_type = ? AND is_active = ?", "description", true).
		Order("sort_order").Find(&descriptions)

	var images []models.HomeContent
	config.DB.Where("content_type = ? AND is_active = ?", "image", true).
		Order("sort_order").Find(&images)

	// Normalisasi path file lama (tanpa garis depan) agar tetap bisa diakses.
	for i := range images {
		images[i].FilePath = normalizeFilePath(images[i].FilePath)
	}

	// Working instruction aktif hanya satu; dikirim sebagai objek (bukan array)
	// agar konsisten dengan sisi frontend.
	var workingInstruction models.HomeContent
	var workingInstructionData *models.HomeContent
	if err := config.DB.Where("content_type = ? AND is_active = ?", "working_instruction", true).
		Order("sort_order").Order("created_at DESC").
		First(&workingInstruction).Error; err == nil {
		workingInstruction.FilePath = normalizeFilePath(workingInstruction.FilePath)
		workingInstructionData = &workingInstruction
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"descriptions":        descriptions,
			"images":              images,
			"working_instruction": workingInstructionData,
		},
	})
}

func UpdateDescription(c *gin.Context) {
	id := c.Param("id")

	var existing models.HomeContent
	if err := config.DB.Where("id = ?", id).First(&existing).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Content not found"})
		return
	}

	var req UpdateDescriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request"})
		return
	}

	existing.Title = &req.Title
	existing.Description = &req.Description
	existing.UpdatedAt = time.Now()

	if err := config.DB.Save(&existing).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to update description"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Description updated successfully",
		"data":    existing,
	})
}

func UploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "No file uploaded"})
		return
	}

	title := c.PostForm("title")

	// Save file
	filename := uuid.New().String() + filepath.Ext(file.Filename)
	destPath := filepath.Join("uploads", "home", filename)

	if err := c.SaveUploadedFile(file, destPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal menyimpan file"})
		return
	}

	userID, _ := c.Get("user_id")
	createdBy, _ := userID.(string)
	fileType := file.Header.Get("Content-Type")

	content := models.HomeContent{
		ID:          uuid.New().String(),
		ContentType: "image",
		Title:       &title,
		FilePath:    ptr(publicFilePath(destPath)),
		FileType:    &fileType,
		SortOrder:   0,
		IsActive:    true,
		CreatedBy:   createdBy,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := config.DB.Create(&content).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to save content"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Image uploaded successfully",
		"data":    content,
	})
}

func DeleteImage(c *gin.Context) {
	id := c.Param("id")

	var existing models.HomeContent
	if err := config.DB.Where("id = ?", id).First(&existing).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Image not found"})
		return
	}

	if err := config.DB.Delete(&existing).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to delete image"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Image deleted successfully"})
}

const maksUkuranWorkingInstruction = 10 << 20 // 10 MB

func UploadWorkingInstruction(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Tidak ada file yang diunggah"})
		return
	}

	if file.Size > maksUkuranWorkingInstruction {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Sprintf("Ukuran file maksimal %d MB", maksUkuranWorkingInstruction>>20),
		})
		return
	}

	if strings.ToLower(filepath.Ext(file.Filename)) != ".pdf" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Format file harus PDF (.pdf)"})
		return
	}

	// Validasi magic bytes: dokumen PDF selalu diawali "%PDF-".
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal membuka file"})
		return
	}
	defer src.Close()

	header := make([]byte, 5)
	if _, err := io.ReadFull(src, header); err != nil || string(header) != "%PDF-" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "File bukan PDF yang valid"})
		return
	}

	destPath := filepath.Join("uploads", "home", uuid.New().String()+".pdf")
	if err := c.SaveUploadedFile(file, destPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal menyimpan file"})
		return
	}

	// Hanya satu working instruction yang aktif, jadi file lama dinonaktifkan
	// ketika yang baru diunggah.
	if err := config.DB.Model(&models.HomeContent{}).
		Where("content_type = ? AND is_active = ?", "working_instruction", true).
		Update("is_active", false).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal mengganti working instruction lama"})
		return
	}

	userID, _ := c.Get("user_id")
	createdBy, _ := userID.(string)
	title := filepath.Base(file.Filename)

	content := models.HomeContent{
		ID:          uuid.New().String(),
		ContentType: "working_instruction",
		Title:       &title,
		FilePath:    ptr(publicFilePath(destPath)),
		FileType:    ptr("application/pdf"),
		SortOrder:   0,
		IsActive:    true,
		CreatedBy:   createdBy,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := config.DB.Create(&content).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal menyimpan content"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Working instruction berhasil diunggah",
		"data":    content,
	})
}

func ptr(v string) *string { return &v }

// publicFilePath mengubah path di disk menjadi URL absolut pada origin yang sama
// (mis. "uploads/home/x.pdf" -> "/uploads/home/x.pdf") supaya file dapat diakses
// baik dari frontend SPA maupun langsung dari backend.
func publicFilePath(path string) string {
	return "/" + strings.TrimPrefix(filepath.ToSlash(path), "/")
}

func normalizeFilePath(p *string) *string {
	if p == nil || *p == "" {
		return p
	}
	return ptr(publicFilePath(*p))
}

func DeleteWorkingInstruction(c *gin.Context) {
	id := c.Param("id")

	var existing models.HomeContent
	if err := config.DB.Where("id = ?", id).First(&existing).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Working instruction not found"})
		return
	}

	if err := config.DB.Delete(&existing).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to delete working instruction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Working instruction deleted successfully"})
}
