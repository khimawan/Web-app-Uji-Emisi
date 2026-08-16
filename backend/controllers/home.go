package controllers

import (
	"net/http"
	"path/filepath"
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

	var workingInstruction []models.HomeContent
	config.DB.Where("content_type = ? AND is_active = ?", "working_instruction", true).
		Order("sort_order").First(&workingInstruction)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"descriptions":        descriptions,
			"images":              images,
			"working_instruction": workingInstruction,
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
 filepath := filepath.Join("uploads", "home", filename)

	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to save file"})
		return
	}

	userID, _ := c.Get("user_id")

	content := models.HomeContent{
		ID:          uuid.New().String(),
		ContentType: "image",
		Title:       &title,
		FilePath:    &filepath,
		FileType:    &file.Header.Get("Content-Type"),
		SortOrder:   0,
		IsActive:    true,
		CreatedBy:   userID.(string),
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

func UploadWorkingInstruction(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "No file uploaded"})
		return
	}

	// Save file
	filename := uuid.New().String() + filepath.Ext(file.Filename)
 filepath := filepath.Join("uploads", "home", filename)

	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to save file"})
		return
	}

	userID, _ := c.Get("user_id")

	content := models.HomeContent{
		ID:          uuid.New().String(),
		ContentType: "working_instruction",
		FilePath:    &filepath,
		FileType:    &file.Header.Get("Content-Type"),
		SortOrder:   0,
		IsActive:    true,
		CreatedBy:   userID.(string),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := config.DB.Create(&content).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to save content"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Working instruction uploaded successfully",
		"data":    content,
	})
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
