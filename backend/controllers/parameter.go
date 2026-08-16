package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"web-app-uji-emisi/config"
	"web-app-uji-emisi/models"
)

type CreateParameterRequest struct {
	Kategori      string   `json:"kategori" binding:"required"`
	TahunMin      int      `json:"tahun_min" binding:"required"`
	TahunMax      *int     `json:"tahun_max"`
	TahunOperator string   `json:"tahun_operator" binding:"required"`
	COMax         *float64 `json:"co_max"`
	HCMax         *float64 `json:"hc_max"`
	OpasitasMax   *float64 `json:"opasitas_max"`
}

func GetParameters(c *gin.Context) {
	kategori := c.Query("kategori")

	query := config.DB.Model(&models.Parameter{})
	if kategori != "" {
		query = query.Where("kategori = ?", kategori)
	}

	var parameters []models.Parameter
	if err := query.Order("kategori, tahun_min").Find(&parameters).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to fetch parameters"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": parameters})
}

func CreateParameter(c *gin.Context) {
	var req CreateParameterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request"})
		return
	}

	parameter := models.Parameter{
		ID:            uuid.New().String(),
		Kategori:      req.Kategori,
		TahunMin:      req.TahunMin,
		TahunMax:      req.TahunMax,
		TahunOperator: req.TahunOperator,
		COMax:         req.COMax,
		HCMax:         req.HCMax,
		OpasitasMax:   req.OpasitasMax,
		IsActive:      true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := config.DB.Create(&parameter).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to create parameter"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Parameter added successfully",
		"data":    parameter,
	})
}

func UpdateParameter(c *gin.Context) {
	id := c.Param("id")

	var existing models.Parameter
	if err := config.DB.Where("id = ?", id).First(&existing).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Parameter not found"})
		return
	}

	var req CreateParameterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request"})
		return
	}

	existing.Kategori = req.Kategori
	existing.TahunMin = req.TahunMin
	existing.TahunMax = req.TahunMax
	existing.TahunOperator = req.TahunOperator
	existing.COMax = req.COMax
	existing.HCMax = req.HCMax
	existing.OpasitasMax = req.OpasitasMax
	existing.UpdatedAt = time.Now()

	if err := config.DB.Save(&existing).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to update parameter"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Parameter updated successfully",
		"data":    existing,
	})
}

func DeleteParameter(c *gin.Context) {
	id := c.Param("id")

	var existing models.Parameter
	if err := config.DB.Where("id = ?", id).First(&existing).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Parameter not found"})
		return
	}

	if err := config.DB.Delete(&existing).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to delete parameter"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Parameter deleted successfully"})
}
