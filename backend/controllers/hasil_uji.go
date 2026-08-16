package controllers

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"web-app-uji-emisi/config"
	"web-app-uji-emisi/models"
)

type CreateHasilUjiRequest struct {
	KendaraanID string   `json:"kendaraan_id" binding:"required"`
	CO          *float64 `json:"co"`
	CO2         *float64 `json:"co2"`
	HC          *float64 `json:"hc"`
	O2          *float64 `json:"o2"`
	Lambda      *float64 `json:"lambda"`
	Opasitas    *float64 `json:"opasitas"`
}

func GetHasilUji(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")
	hasilUji := c.Query("hasil_uji")
	kategori := c.Query("kategori")
	tanggalMulai := c.Query("tanggal_mulai")
	tanggalAkhir := c.Query("tanggal_akhir")

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	type HasilUjiWithKendaraan struct {
		models.HasilUji
		models.Kendaraan
	}

	query := config.DB.Table("hasil_uji").
		Select("hasil_uji.*, kendaraan.plat_nomor, kendaraan.merek, kendaraan.tipe, kendaraan.kategori, kendaraan.jenis, kendaraan.kapasitas_mesin, kendaraan.tahun_pembuatan").
		Joins("JOIN kendaraan ON hasil_uji.kendaraan_id = kendaraan.id").
		Where("hasil_uji.deleted_at IS NULL")

	if search != "" {
		query = query.Where("kendaraan.plat_nomor ILIKE ? OR kendaraan.merek ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if hasilUji != "" {
		query = query.Where("hasil_uji.hasil_uji = ?", hasilUji)
	}
	if kategori != "" {
		query = query.Where("kendaraan.kategori = ?", kategori)
	}
	if tanggalMulai != "" {
		query = query.Where("hasil_uji.tested_at >= ?", tanggalMulai)
	}
	if tanggalAkhir != "" {
		query = query.Where("hasil_uji.tested_at <= ?", tanggalAkhir+" 23:59:59")
	}

	var total int64
	query.Count(&total)

	var results []map[string]interface{}
	if err := query.Offset(offset).Limit(limit).Order("hasil_uji.tested_at DESC").Find(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to fetch hasil uji"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"items": results,
			"pagination": gin.H{
				"total":      total,
				"page":       page,
				"limit":      limit,
				"totalPages": (total + int64(limit) - 1) / int64(limit),
			},
		},
	})
}

func CreateHasilUji(c *gin.Context) {
	var req CreateHasilUjiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request"})
		return
	}

	var kendaraan models.Kendaraan
	if err := config.DB.Where("id = ?", req.KendaraanID).First(&kendaraan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Kendaraan not found"})
		return
	}

	userID, _ := c.Get("user_id")

	// Determine hasil uji based on parameters
	hasilUji, valid, catatan := evaluateEmission(kendaraan, req.CO, req.CO2, req.HC, req.O2, req.Opasitas)

	now := time.Now()
	hasilUjiRecord := models.HasilUji{
		ID:          uuid.New().String(),
		KendaraanID: req.KendaraanID,
		CO:          req.CO,
		CO2:         req.CO2,
		HC:          req.HC,
		O2:          req.O2,
		Lambda:      req.Lambda,
		Opasitas:    req.Opasitas,
		HasilUji:    hasilUji,
		Valid:       valid,
		Catatan:     catatan,
		TestedBy:    userID.(string),
		TestedAt:    now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := config.DB.Create(&hasilUjiRecord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to create hasil uji"})
		return
	}

	// Get pop-up notes
	var popupNotes []models.PopupNote
	config.DB.Where("is_active = ?", true).Find(&popupNotes)

	var triggeredNotes []string
	for _, note := range popupNotes {
		if checkPopupNote(kendaraan, note, req.CO, req.CO2, req.HC, req.O2, req.Opasitas) {
			triggeredNotes = append(triggeredNotes, note.Note)
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Hasil uji added successfully",
		"data": gin.H{
			"hasil_uji":    hasilUjiRecord,
			"popup_notes":  triggeredNotes,
		},
	})
}

func UpdateHasilUji(c *gin.Context) {
	id := c.Param("id")

	var existing models.HasilUji
	if err := config.DB.Where("id = ?", id).First(&existing).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Hasil uji not found"})
		return
	}

	var req CreateHasilUjiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request"})
		return
	}

	var kendaraan models.Kendaraan
	if err := config.DB.Where("id = ?", req.KendaraanID).First(&kendaraan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Kendaraan not found"})
		return
	}

	hasilUji, valid, catatan := evaluateEmission(kendaraan, req.CO, req.CO2, req.HC, req.O2, req.Opasitas)

	existing.CO = req.CO
	existing.CO2 = req.CO2
	existing.HC = req.HC
	existing.O2 = req.O2
	existing.Lambda = req.Lambda
	existing.Opasitas = req.Opasitas
	existing.HasilUji = hasilUji
	existing.Valid = valid
	existing.Catatan = catatan
	existing.UpdatedAt = time.Now()

	if err := config.DB.Save(&existing).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to update hasil uji"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Hasil uji updated successfully",
		"data":    existing,
	})
}

func DeleteHasilUji(c *gin.Context) {
	id := c.Param("id")

	var existing models.HasilUji
	if err := config.DB.Where("id = ?", id).First(&existing).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Hasil uji not found"})
		return
	}

	if err := config.DB.Delete(&existing).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to delete hasil uji"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Hasil uji deleted successfully"})
}

func GetStatistics(c *gin.Context) {
	var totalKendaraan int64
	var totalDinas int64
	var totalUmum int64
	var totalLulus int64
	var totalTidakLulus int64
	var totalBensin int64
	var totalSolar int64

	config.DB.Model(&models.Kendaraan{}).Count(&totalKendaraan)
	config.DB.Model(&models.Kendaraan{}).Where("kategori = ?", "Dinas").Count(&totalDinas)
	config.DB.Model(&models.Kendaraan{}).Where("kategori = ?", "Umum").Count(&totalUmum)
	config.DB.Model(&models.HasilUji{}).Where("hasil_uji = ?", "Lulus").Count(&totalLulus)
	config.DB.Model(&models.HasilUji{}).Where("hasil_uji = ?", "Tidak Lulus").Count(&totalTidakLulus)
	config.DB.Model(&models.Kendaraan{}).Where("jenis LIKE ?", "%Bensin%").Count(&totalBensin)
	config.DB.Model(&models.Kendaraan{}).Where("jenis LIKE ?", "%Solar%").Count(&totalSolar)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"total_kendaraan":   totalKendaraan,
			"total_dinas":       totalDinas,
			"total_umum":        totalUmum,
			"total_lulus":       totalLulus,
			"total_tidak_lulus": totalTidakLulus,
			"total_bensin":      totalBensin,
			"total_solar":       totalSolar,
		},
	})
}

func ExportCSV(c *gin.Context) {
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=hasil_uji_emisi.csv")

	var results []map[string]interface{}
	config.DB.Table("hasil_uji").
		Select("hasil_uji.*, kendaraan.plat_nomor, kendaraan.merek, kendaraan.tipe, kendaraan.kategori, kendaraan.jenis, kendaraan.kapasitas_mesin, kendaraan.tahun_pembuatan").
		Joins("JOIN kendaraan ON hasil_uji.kendaraan_id = kendaraan.id").
		Where("hasil_uji.deleted_at IS NULL").
		Order("hasil_uji.tested_at DESC").
		Find(&results)

	w := csv.NewWriter(c.Writer)
	defer w.Flush()

	// Header
	w.Write([]string{
		"No", "Waktu Input", "Kategori Kendaraan", "Jenis Kendaraan", "Plat Nomor",
		"Merek", "Tipe", "Kapasitas Mesin", "Tahun Pembuatan",
		"CO", "CO2", "HC", "O2", "Lambda", "Opasitas",
		"Hasil Uji", "Valid", "Catatan",
	})

	// Data
	for i, row := range results {
		w.Write([]string{
			strconv.Itoa(i + 1),
			fmt.Sprintf("%v", row["tested_at"]),
			fmt.Sprintf("%v", row["kategori"]),
			fmt.Sprintf("%v", row["jenis"]),
			fmt.Sprintf("%v", row["plat_nomor"]),
			fmt.Sprintf("%v", row["merek"]),
			fmt.Sprintf("%v", row["tipe"]),
			fmt.Sprintf("%v", row["kapasitas_mesin"]),
			fmt.Sprintf("%v", row["tahun_pembuatan"]),
			fmt.Sprintf("%v", row["co"]),
			fmt.Sprintf("%v", row["co2"]),
			fmt.Sprintf("%v", row["hc"]),
			fmt.Sprintf("%v", row["o2"]),
			fmt.Sprintf("%v", row["lambda"]),
			fmt.Sprintf("%v", row["opasitas"]),
			fmt.Sprintf("%v", row["hasil_uji"]),
			fmt.Sprintf("%v", row["valid"]),
			fmt.Sprintf("%v", row["catatan"]),
		})
	}
}

func evaluateEmission(kendaraan models.Kendaraan, co, co2, hc, o2, opasitas *float64) (string, bool, *string) {
	var parameters []models.Parameter
	config.DB.Where("kategori = ? AND is_active = ?", kendaraan.Jenis, true).Find(&parameters)

	for _, param := range parameters {
		match := false
		switch param.TahunOperator {
		case "<":
			match = kendaraan.TahunPembuatan < param.TahunMin
		case ">":
			match = kendaraan.TahunPembuatan > param.TahunMin
		case "between":
			match = kendaraan.TahunPembuatan >= param.TahunMin && kendaraan.TahunPembuatan <= *param.TahunMax
		}

		if match {
			valid := true
			var catatanParts []string

			// Check Bensin parameters
			if param.COMax != nil && co != nil {
				if *co > *param.COMax {
					valid = false
					catatanParts = append(catatanParts, fmt.Sprintf("CO %.2f%% melebihi batas %.2f%%", *co, *param.COMax))
				}
			}
			if param.HCMax != nil && hc != nil {
				if *hc > *param.HCMax {
					valid = false
					catatanParts = append(catatanParts, fmt.Sprintf("HC %.0f ppm melebihi batas %.0f ppm", *hc, *param.HCMax))
				}
			}

			// Check Solar parameters
			if param.OpasitasMax != nil && opasitas != nil {
				if *opasitas > *param.OpasitasMax {
					valid = false
					catatanParts = append(catatanParts, fmt.Sprintf("Opasitas %.1f%% melebihi batas %.1f%%", *opasitas, *param.OpasitasMax))
				}
			}

			if valid {
				return "Lulus", true, nil
			} else {
				catatan := ""
				if len(catatanParts) > 0 {
					catatan = catatanParts[0]
					for _, part := range catatanParts[1:] {
						catatan += "; " + part
					}
				}
				return "Tidak Lulus", false, &catatan
			}
		}
	}

	return "Tidak Lulus", false, stringPtr("Parameter tidak ditemukan")
}

func checkPopupNote(kendaraan models.Kendaraan, note models.PopupNote, co, co2, hc, o2, opasitas *float64) bool {
	// Check if jenis kendaraan matches
	jenisList := splitAndTrim(note.JenisKendaraan)
	jenisMatch := false
	for _, j := range jenisList {
		if j == kendaraan.Jenis {
			jenisMatch = true
			break
		}
	}
	if !jenisMatch {
		return false
	}

	// Check tahun
	tahunMatch := false
	switch note.TahunOperator {
	case ">":
		tahunMatch = kendaraan.TahunPembuatan > note.TahunValue
	case "<":
		tahunMatch = kendaraan.TahunPembuatan < note.TahunValue
	case ">=":
		tahunMatch = kendaraan.TahunPembuatan >= note.TahunValue
	case "<=":
		tahunMatch = kendaraan.TahunPembuatan <= note.TahunValue
	case "=":
		tahunMatch = kendaraan.TahunPembuatan == note.TahunValue
	}
	if !tahunMatch {
		return false
	}

	// Check parameter value
	var nilai *float64
	switch note.ParameterUji {
	case "CO":
		nilai = co
	case "CO2":
		nilai = co2
	case "HC":
		nilai = hc
	case "O2":
		nilai = o2
	case "Opasitas":
		nilai = opasitas
	}

	if nilai == nil {
		return false
	}

	switch note.NilaiOperator {
	case ">":
		return *nilai > note.NilaiValue
	case "<":
		return *nilai < note.NilaiValue
	case ">=":
		return *nilai >= note.NilaiValue
	case "<=":
		return *nilai <= note.NilaiValue
	case "=":
		return *nilai == note.NilaiValue
	}

	return false
}

func splitAndTrim(s string) []string {
	var result []string
	for _, item := range splitString(s, ",") {
		trimmed := trimSpace(item)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func splitString(s, delimiter string) []string {
	var result []string
	start := 0
	for i := 0; i <= len(s)-len(delimiter); i++ {
		if s[i:i+len(delimiter)] == delimiter {
			result = append(result, s[start:i])
			start = i + len(delimiter)
		}
	}
	result = append(result, s[start:])
	return result
}

func trimSpace(s string) string {
	start := 0
	end := len(s)
	for start < end && s[start] == ' ' {
		start++
	}
	for end > start && s[end-1] == ' ' {
		end--
	}
	return s[start:end]
}

func stringPtr(s string) *string {
	return &s
}
