package controllers

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"web-app-uji-emisi/config"
	"web-app-uji-emisi/models"
)

type CreateKendaraanRequest struct {
	Kategori       string `json:"kategori" binding:"required"`
	Jenis          string `json:"jenis" binding:"required"`
	PlatNomor      string `json:"plat_nomor" binding:"required"`
	Merek          string `json:"merek" binding:"required"`
	Tipe           string `json:"tipe" binding:"required"`
	KapasitasMesin int    `json:"kapasitas_mesin" binding:"required"`
	TahunPembuatan int    `json:"tahun_pembuatan" binding:"required"`
}

func GetKendaraan(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")
	kategori := c.Query("kategori")
	jenis := c.Query("jenis")

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	query := config.DB.Model(&models.Kendaraan{})

	if search != "" {
		query = query.Where("plat_nomor ILIKE ? OR merek ILIKE ? OR tipe ILIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	if kategori != "" {
		query = query.Where("kategori = ?", kategori)
	}
	if jenis != "" {
		query = query.Where("jenis = ?", jenis)
	}

	var total int64
	query.Count(&total)

	var kendaraan []models.Kendaraan
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&kendaraan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to fetch kendaraan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"items": kendaraan,
			"pagination": gin.H{
				"total":      total,
				"page":       page,
				"limit":      limit,
				"totalPages": (total + int64(limit) - 1) / int64(limit),
			},
		},
	})
}

func GetKendaraanByID(c *gin.Context) {
	id := c.Param("id")

	var kendaraan models.Kendaraan
	if err := config.DB.Where("id = ?", id).First(&kendaraan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Kendaraan not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": kendaraan})
}

func CreateKendaraan(c *gin.Context) {
	var req CreateKendaraanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request"})
		return
	}

	userID, _ := c.Get("user_id")

	kendaraan := models.Kendaraan{
		ID:             uuid.New().String(),
		Kategori:       req.Kategori,
		Jenis:          req.Jenis,
		PlatNomor:      req.PlatNomor,
		Merek:          req.Merek,
		Tipe:           req.Tipe,
		KapasitasMesin: req.KapasitasMesin,
		TahunPembuatan: req.TahunPembuatan,
		CreatedBy:      userID.(string),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := config.DB.Create(&kendaraan).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Failed to create kendaraan. Plat nomor might already exist."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Kendaraan added successfully",
		"data":    kendaraan,
	})
}

func UpdateKendaraan(c *gin.Context) {
	id := c.Param("id")

	var existing models.Kendaraan
	if err := config.DB.Where("id = ?", id).First(&existing).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Kendaraan not found"})
		return
	}

	var req CreateKendaraanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request"})
		return
	}

	existing.Kategori = req.Kategori
	existing.Jenis = req.Jenis
	existing.PlatNomor = req.PlatNomor
	existing.Merek = req.Merek
	existing.Tipe = req.Tipe
	existing.KapasitasMesin = req.KapasitasMesin
	existing.TahunPembuatan = req.TahunPembuatan
	existing.UpdatedAt = time.Now()

	if err := config.DB.Save(&existing).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Failed to update kendaraan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Kendaraan updated successfully",
		"data":    existing,
	})
}

func DeleteKendaraan(c *gin.Context) {
	id := c.Param("id")

	var existing models.Kendaraan
	if err := config.DB.Where("id = ?", id).First(&existing).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Kendaraan not found"})
		return
	}

	if err := config.DB.Delete(&existing).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to delete kendaraan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Kendaraan deleted successfully"})
}

func UploadCSV(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "No file uploaded"})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to open file"})
		return
	}
	defer src.Close()

	reader := csv.NewReader(src)
	records, err := reader.ReadAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid CSV format"})
		return
	}

	if len(records) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "CSV file is empty"})
		return
	}

	userID, _ := c.Get("user_id")
	total := len(records) - 1
	success := 0
	failed := 0
	errors := []gin.H{}

	for i, record := range records[1:] {
		if len(record) < 7 {
			errors = append(errors, gin.H{"row": i + 2, "message": "Invalid number of columns"})
			failed++
			continue
		}

		kapasitasMesin, err := strconv.Atoi(record[5])
		if err != nil {
			errors = append(errors, gin.H{"row": i + 2, "message": "Invalid kapasitas mesin"})
			failed++
			continue
		}

		tahunPembuatan, err := strconv.Atoi(record[6])
		if err != nil {
			errors = append(errors, gin.H{"row": i + 2, "message": "Invalid tahun pembuatan"})
			failed++
			continue
		}

		kendaraan := models.Kendaraan{
			ID:             uuid.New().String(),
			Kategori:       record[0],
			Jenis:          record[1],
			PlatNomor:      record[2],
			Merek:          record[3],
			Tipe:           record[4],
			KapasitasMesin: kapasitasMesin,
			TahunPembuatan: tahunPembuatan,
			CreatedBy:      userID.(string),
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		if err := config.DB.Create(&kendaraan).Error; err != nil {
			errors = append(errors, gin.H{"row": i + 2, "message": fmt.Sprintf("Failed to create: %v", err)})
			failed++
		} else {
			success++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "CSV uploaded successfully",
		"data": gin.H{
			"total":   total,
			"success": success,
			"failed":  failed,
			"errors":  errors,
		},
	})
}

func DownloadCSVTemplate(c *gin.Context) {
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=template_kendaraan.csv")

	c.String(http.StatusOK, "kategori,jenis,plat_nomor,merek,tipe,kapasitas_mesin,tahun_pembuatan\n")
}

func GetAllKendaraan(c *gin.Context) {
	var kendaraan []models.Kendaraan
	if err := config.DB.Order("created_at DESC").Find(&kendaraan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to fetch kendaraan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": kendaraan})
}

func SearchKendaraan(c *gin.Context) {
	search := c.Query("q")

	var kendaraan []models.Kendaraan
	query := config.DB.Model(&models.Kendaraan{})

	if search != "" {
		query = query.Where("plat_nomor ILIKE ? OR merek ILIKE ? OR tipe ILIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Order("created_at DESC").Limit(20).Find(&kendaraan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to search kendaraan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": kendaraan})
}
