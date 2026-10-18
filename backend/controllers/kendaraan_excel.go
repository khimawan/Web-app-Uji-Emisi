package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"web-app-uji-emisi/config"
	"web-app-uji-emisi/models"
)

// ---------------------------------------------------------------------------
// Template Excel Input Data Kendaraan
//
// Isi template sengaja dibuat sama dengan template CSV (panduan + satu baris
// contoh untuk tiap kombinasi kategori & jenis) supaya pengguna tidak perlu
// belajar dua format yang berbeda.
// ---------------------------------------------------------------------------

var kolomKendaraanXLSX = []string{"A", "B", "C", "D", "E", "F", "G"}

const (
	shPanduanKendaraan      = "Panduan"
	shDataKendaraan         = "Data Kendaraan"
	shReferensiKendaraan    = "Daftar Nilai Valid"
	barisHeaderKendaraan    = 6
	barisAwalDataKendaraan  = barisHeaderKendaraan + 1
	barisAkhirDataKendaraan = 1006
)

// keteranganJenisKendaraan menjelaskan tiap nilai jenis agar pengguna tidak
// salah memilih antara Bensin M, Bensin N&O, Solar JBB, dan Solar GVW.
var keteranganJenisKendaraan = map[string]string{
	"Bensin M":   "Angkutan penumpang (motor & mobil) dengan kapasitas mesin <= 3.500 cc",
	"Bensin N&O": "Angkutan barang / gendang (N&O = Non Operational), kapasitas mesin > 3.500 cc",
	"Solar JBB":  "Solar, berat total kendaraan (JBB = Laden Beban Berat) <= 3,5 ton",
	"Solar GVW":  "Solar, berat total kendaraan (GVW = Gross Vehicle Weight) > 3,5 ton",
}

func DownloadTemplateKendaraanExcel(c *gin.Context) {
	f := excelize.NewFile()
	defer f.Close()

	if _, err := f.NewSheet(shPanduanKendaraan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal membuat template Excel"})
		return
	}
	idxData, err := f.NewSheet(shDataKendaraan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal membuat template Excel"})
		return
	}
	if _, err := f.NewSheet(shReferensiKendaraan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal membuat template Excel"})
		return
	}
	// Hapus sheet kosong bawaan excelize agar template rapi.
	if err := f.DeleteSheet("Sheet1"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal membuat template Excel"})
		return
	}

	gayaJudul, _ := f.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true, Size: 13, Color: "1F2937"}})
	gayaHeader, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"2563EB"}},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
	})
	gayaContoh, _ := f.NewStyle(&excelize.Style{Font: &excelize.Font{Italic: true, Color: "6B7280"}})

	tulisSheetPanduanKendaraan(f, gayaJudul)
	tulisSheetDataKendaraan(f, gayaJudul, gayaHeader, gayaContoh)
	tulisSheetReferensiKendaraan(f, gayaJudul, gayaHeader)

	f.SetActiveSheet(idxData)

	buf, err := f.WriteToBuffer()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal membuat template Excel"})
		return
	}

	c.Header("Content-Disposition", `attachment; filename=template_kendaraan.xlsx`)
	c.Data(http.StatusOK, mimeTypeXLSX, buf.Bytes())
}

func tulisSheetPanduanKendaraan(f *excelize.File, gayaJudul int) {
	f.SetCellValue(shPanduanKendaraan, "A1", "Panduan Input Data Kendaraan - Aplikasi Uji Emisi")
	f.SetCellStyle(shPanduanKendaraan, "A1", "A1", gayaJudul)
	for i, line := range panduanTemplateKendaraan() {
		f.SetCellValue(shPanduanKendaraan, fmt.Sprintf("A%d", i+3), line)
	}
	f.SetColWidth(shPanduanKendaraan, "A", "A", 100)
}

func tulisSheetDataKendaraan(f *excelize.File, gayaJudul, gayaHeader, gayaContoh int) {
	f.SetCellValue(shDataKendaraan, "A1", fmt.Sprintf("Isi data mulai baris %d. Baris contoh (abu-abu) bisa diedit atau dihapus.", barisAwalDataKendaraan))
	f.SetCellStyle(shDataKendaraan, "A1", "A1", gayaJudul)
	f.SetCellValue(shDataKendaraan, "A2", "Jangan ubah nama kolom pada baris "+fmt.Sprint(barisHeaderKendaraan)+". Gunakan dropdown pada kolom Kategori dan Jenis.")
	f.SetCellStyle(shDataKendaraan, "A2", "A2", gayaJudul)

	for i, nama := range kolomKendaraanTemplate {
		f.SetCellValue(shDataKendaraan, fmt.Sprintf("%s%d", kolomKendaraanXLSX[i], barisHeaderKendaraan), nama)
	}
	last := fmt.Sprintf("%s%d", kolomKendaraanXLSX[len(kolomKendaraanXLSX)-1], barisHeaderKendaraan)
	f.SetCellStyle(shDataKendaraan, fmt.Sprintf("A%d", barisHeaderKendaraan), last, gayaHeader)

	for i, row := range contohTemplateKendaraan() {
		baris := barisAwalDataKendaraan + i
		for j, nilai := range row {
			f.SetCellValue(shDataKendaraan, fmt.Sprintf("%s%d", kolomKendaraanXLSX[j], baris), nilai)
		}
		f.SetCellStyle(shDataKendaraan,
			fmt.Sprintf("A%d", baris),
			fmt.Sprintf("%s%d", kolomKendaraanXLSX[len(kolomKendaraanXLSX)-1], baris),
			gayaContoh)
	}

	for i, lebar := range []float64{14, 16, 20, 18, 20, 18, 18} {
		f.SetColWidth(shDataKendaraan, kolomKendaraanXLSX[i], kolomKendaraanXLSX[i], lebar)
	}
	f.SetRowHeight(shDataKendaraan, barisHeaderKendaraan, 24)
	kunciPanes(f, shDataKendaraan, barisHeaderKendaraan)
	_ = f.AutoFilter(shDataKendaraan, fmt.Sprintf("A%d:%s", barisHeaderKendaraan, last), nil)

	// Dropdown nilai valid: mencegah salah ketik kategori/jenis sejak awal.
	rentang := func(kolom string) string {
		return fmt.Sprintf("%s%d:%s%d", kolom, barisAwalDataKendaraan, kolom, barisAkhirDataKendaraan)
	}
	dvKategori := excelize.NewDataValidation(true)
	dvKategori.Sqref = rentang("A")
	if err := dvKategori.SetDropList(KategoriKendaraanValid); err == nil {
		_ = f.AddDataValidation(shDataKendaraan, dvKategori)
	}
	dvJenis := excelize.NewDataValidation(true)
	dvJenis.Sqref = rentang("B")
	if err := dvJenis.SetDropList(JenisKendaraanValid); err == nil {
		_ = f.AddDataValidation(shDataKendaraan, dvJenis)
	}
}

func tulisSheetReferensiKendaraan(f *excelize.File, gayaJudul, gayaHeader int) {
	// Bagian 1: daftar nilai kategori & jenis yang valid beserta keterangannya.
	f.SetCellValue(shReferensiKendaraan, "A1", "Nilai yang WAJIB dipakai pada kolom kategori dan jenis")
	f.SetCellStyle(shReferensiKendaraan, "A1", "A1", gayaJudul)

	refHeader := []string{"kolom", "nilai", "keterangan"}
	for i, nama := range refHeader {
		f.SetCellValue(shReferensiKendaraan, fmt.Sprintf("%s3", kolomKendaraanXLSX[i]), nama)
	}
	f.SetCellStyle(shReferensiKendaraan, "A3", "C3", gayaHeader)

	baris := 4
	for _, k := range KategoriKendaraanValid {
		f.SetCellValue(shReferensiKendaraan, fmt.Sprintf("A%d", baris), "kategori")
		f.SetCellValue(shReferensiKendaraan, fmt.Sprintf("B%d", baris), k)
		f.SetCellValue(shReferensiKendaraan, fmt.Sprintf("C%d", baris), "Kategori kendaraan")
		baris++
	}
	for _, j := range JenisKendaraanValid {
		f.SetCellValue(shReferensiKendaraan, fmt.Sprintf("A%d", baris), "jenis")
		f.SetCellValue(shReferensiKendaraan, fmt.Sprintf("B%d", baris), j)
		f.SetCellValue(shReferensiKendaraan, fmt.Sprintf("C%d", baris), keteranganJenisKendaraan[j])
		baris++
	}

	// Bagian 2: plat nomor yang sudah terdaftar, supaya tidak terduplikasi.
	baris++
	baris++
	f.SetCellValue(shReferensiKendaraan, fmt.Sprintf("A%d", baris),
		"Plat nomor yang SUDAH terdaftar (jangan diunggah ulang)")
	f.SetCellStyle(shReferensiKendaraan, fmt.Sprintf("A%d", baris), fmt.Sprintf("A%d", baris), gayaJudul)
	baris++
	f.SetCellValue(shReferensiKendaraan, fmt.Sprintf("A%d", baris), "plat_nomor")
	f.SetCellStyle(shReferensiKendaraan, fmt.Sprintf("A%d", baris), fmt.Sprintf("A%d", baris), gayaHeader)
	baris++

	var kendaraan []models.Kendaraan
	config.DB.Order("created_at").Limit(maksKendaraanDiTemplate).Find(&kendaraan)
	for _, k := range kendaraan {
		f.SetCellValue(shReferensiKendaraan, fmt.Sprintf("A%d", baris), k.PlatNomor)
		baris++
	}
	if len(kendaraan) == 0 {
		f.SetCellValue(shReferensiKendaraan, fmt.Sprintf("A%d", baris), "(belum ada kendaraan terdaftar)")
	}

	f.SetColWidth(shReferensiKendaraan, "A", "A", 34)
	f.SetColWidth(shReferensiKendaraan, "B", "B", 20)
	f.SetColWidth(shReferensiKendaraan, "C", "C", 70)
}
