package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"web-app-uji-emisi/config"
	"web-app-uji-emisi/models"
)

type CreatePopupNoteRequest struct {
	JenisKendaraan string  `json:"jenis_kendaraan" binding:"required"`
	TahunOperator  string  `json:"tahun_operator" binding:"required"`
	TahunValue     int     `json:"tahun_value" binding:"required"`
	ParameterUji   string  `json:"parameter_uji" binding:"required"`
	NilaiOperator  string  `json:"nilai_operator" binding:"required"`
	NilaiValue     float64 `json:"nilai_value" binding:"required"`
	Note           string  `json:"note" binding:"required"`
}

func GetPopupNotes(c *gin.Context) {
	var notes []models.PopupNote
	if err := config.DB.Order("created_at DESC").Find(&notes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to fetch popup notes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": notes})
}

func CreatePopupNote(c *gin.Context) {
	var req CreatePopupNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request"})
		return
	}

	note := models.PopupNote{
		ID:             uuid.New().String(),
		JenisKendaraan: req.JenisKendaraan,
		TahunOperator:  req.TahunOperator,
		TahunValue:     req.TahunValue,
		ParameterUji:   req.ParameterUji,
		NilaiOperator:  req.NilaiOperator,
		NilaiValue:     req.NilaiValue,
		Note:           req.Note,
		IsActive:       true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := config.DB.Create(&note).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to create popup note"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Popup note added successfully",
		"data":    note,
	})
}

func UpdatePopupNote(c *gin.Context) {
	id := c.Param("id")

	var existing models.PopupNote
	if err := config.DB.Where("id = ?", id).First(&existing).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Popup note not found"})
		return
	}

	var req CreatePopupNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request"})
		return
	}

	existing.JenisKendaraan = req.JenisKendaraan
	existing.TahunOperator = req.TahunOperator
	existing.TahunValue = req.TahunValue
	existing.ParameterUji = req.ParameterUji
	existing.NilaiOperator = req.NilaiOperator
	existing.NilaiValue = req.NilaiValue
	existing.Note = req.Note
	existing.UpdatedAt = time.Now()

	if err := config.DB.Save(&existing).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to update popup note"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Popup note updated successfully",
		"data":    existing,
	})
}

func DeletePopupNote(c *gin.Context) {
	id := c.Param("id")

	var existing models.PopupNote
	if err := config.DB.Where("id = ?", id).First(&existing).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Popup note not found"})
		return
	}

	if err := config.DB.Delete(&existing).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to delete popup note"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Popup note deleted successfully"})
}
