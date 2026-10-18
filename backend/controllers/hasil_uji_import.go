package controllers

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
	"web-app-uji-emisi/config"
	"web-app-uji-emisi/models"
)

// ---------------------------------------------------------------------------
// Definisi template import hasil uji
// ---------------------------------------------------------------------------

// kolomHasilUjiTemplate adalah urutan kolom template, dipakai juga saat file
// diunggah tanpa baris header.
var kolomHasilUjiTemplate = []string{
	"plat_nomor", "co", "co2", "hc", "o2", "lambda", "opasitas", "tanggal_uji",
}

var kolomUjiXLSX = []string{"A", "B", "C", "D", "E", "F", "G", "H"}

// batasNilai mendokumentasikan rentang dan satuan tiap parameter uji.
type batasNilai struct {
	Kolom  string
	Min    float64
	Maks   float64
	Satuan string
	Contoh string
}

var batasNilaiUji = []batasNilai{
	{"co", 0, 100, "%", "0.85"},
	{"co2", 0, 100, "%", "12.5"},
	{"hc", 0, 10000, "ppm", "120"},
	{"o2", 0, 100, "%", "1.2"},
	{"lambda", 0.1, 10, "", "0.99"},
	{"opasitas", 0, 100, "%", "12.5"},
}

const (
	maksKendaraanDiTemplate = 1000
	maksBarisPencariHeader  = 20
	maksBarisContoh         = 4
	mimeTypeXLSX            = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
)

// aliasKolomUji memetakan berbagai penulisan nama kolom yang lazim dipakai ke
// nama kanonik template.
var aliasKolomUji = map[string][]string{
	"plat_nomor": {"plat_nomor", "platnomor", "no_polisi", "nomor_polisi", "nopol", "kendaraan"},
	"co":         {"co"},
	"co2":        {"co2", "co_2"},
	"hc":         {"hc"},
	"o2":         {"o2", "o_2"},
	"lambda":     {"lambda", "lamda"},
	"opasitas":   {"opasitas", "opasitas_persen", "opacity"},
	"tanggal_uji": {
		"tanggal_uji", "tanggal", "tgl_uji", "tgl", "tanggal_tes",
		"tanggal_pengujian", "tanggal_uji_emisi",
	},
}

var (
	reTanggalISO = regexp.MustCompile(`^(\d{4})-(\d{1,2})-(\d{1,2})$`)
	reTanggalYMD = regexp.MustCompile(`^(\d{4})[/\-.](\d{1,2})[/\-.](\d{1,2})$`)
	reTanggalDMY = regexp.MustCompile(`^(\d{1,2})[/\-.](\d{1,2})[/\-.](\d{4})$`)
)

// ---------------------------------------------------------------------------
// Template CSV
// ---------------------------------------------------------------------------

func DownloadTemplateHasilUjiCSV(c *gin.Context) {
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", `attachment; filename=template_import_hasil_uji.csv`)

	w := csv.NewWriter(c.Writer)
	defer w.Flush()

	for _, line := range panduanTemplateHasilUji() {
		fmt.Fprintf(c.Writer, "# %s\n", line)
	}
	_ = w.Write(kolomHasilUjiTemplate)
	for _, row := range contohTemplateHasilUji() {
		_ = w.Write(row)
	}

	// Daftar kendaraan terdaftar ditulis sebagai baris komentar di akhir supaya
	// pengguna tidak salah ketik plat nomor. Baris '#' otomatis dilewati import.
	if daftar := daftarKendaraanTemplate(); len(daftar) > 0 {
		fmt.Fprintf(c.Writer, "#\n")
		fmt.Fprintf(c.Writer, "# KENDARAAN YANG SAAT INI TERDAFTAR (gunakan plat_nomor di bawah ini):\n")
		for _, row := range daftar {
			fmt.Fprintf(c.Writer, "# %s\n", strings.Join(row, " | "))
		}
	}
}

func panduanTemplateHasilUji() []string {
	baris := []string{
		"TEMPLATE IMPORT DATA HASIL UJI EMISI - APLIKASI UJI EMISI",
		"",
		"CARA PAKAI:",
		"  1. Kendaraan harus SUDAH terdaftar di halaman 'Input Data Kendaraan'.",
		"     Data hasil uji hanya boleh ditautkan ke plat nomor yang sudah ada.",
		"  2. Ganti plat nomor pada baris contoh dengan plat nomor yang terdaftar,",
		"     lalu hapus baris contoh yang tidak dipakai.",
		"  3. Simpan sebagai .csv dengan encoding UTF-8 (Excel: Save As > CSV UTF-8),",
		"     atau pakai template .xlsx yang juga tersedia.",
		"  4. Unggah lewat tombol 'Upload File' pada halaman Data Hasil Uji.",
		"",
		"KOLOM & FORMAT:",
		"  - plat_nomor : WAJIB. Plat nomor kendaraan yang sudah terdaftar di aplikasi.",
		"                 Huruf besar/kecil dan spasi tidak berpengaruh.",
		"  - co         : CO dalam persen (%). Angka desimal, boleh dikosongkan.",
		"  - co2        : CO2 dalam persen (%). Angka desimal, boleh dikosongkan.",
		"  - hc         : HC dalam ppm. Angka desimal, boleh dikosongkan.",
		"  - o2         : O2 dalam persen (%). Angka desimal, boleh dikosongkan.",
		"  - lambda     : Lambda. Angka desimal, boleh dikosongkan.",
		"  - opasitas   : Opasitas dalam persen (%). Angka desimal, boleh dikosongkan.",
		"  - tanggal_uji : TANGGAL pengujian. Boleh dikosongkan (otomatis diisi waktu unggah).",
		"                 Format: 2024-01-15, 2024/01/15, atau 15/01/2024. Tidak boleh di masa depan.",
		"",
		"PENTING:",
		"  - Minimal SATU parameter (co/co2/hc/o2/lambda/opasitas) harus diisi per baris.",
		"  - Kolom 'hasil_uji', 'valid', dan 'catatan' sengaja TIDAK ada di template.",
		"    Sistem menghitungnya otomatis dari menu Parameter, sama seperti input manual.",
		"  - Satu kendaraan hanya boleh punya satu hasil uji. Kendaraan yang sudah",
		"    punya hasil uji akan ditolak.",
		"  - Pemisah kolom bisa koma (,), titik koma (;), atau tab. Otomatis dideteksi.",
		"  - Baris kosong dan baris berawalan '#' otomatis dilewati.",
		"  - Angka boleh memakai pemisah ribuan (1.25) maupun koma (1,25).",
		"  - Baris yang formatnya salah ditolak, nomor baris + alasannya ditampilkan",
		"    pada hasil unggah. Baris lain tetap tersimpan.",
		"",
		"RENTANG NILAI YANG DITERIMA:",
	}

	for _, b := range batasNilaiUji {
		baris = append(baris, fmt.Sprintf("  - %-9s : %g sampai %g %s (contoh: %s)",
			b.Kolom, b.Min, b.Maks, b.Satuan, b.Contoh))
	}

	return baris
}

// contohTemplateHasilUji membuat baris contoh. Bila ada kendaraan terdaftar,
// plat nomor contoh memakai kendaraan nyata sehingga template langsung bisa
// dicoba oleh pengguna.
func contohTemplateHasilUji() [][]string {
	plat := registeredPlats()
	if len(plat) == 0 {
		plat = []string{"B 1234 CDE", "B 1234 CDF", "B 1234 CDG", "B 1234 CDJ"}
	}

	contoh := [][]string{
		{"0.85", "12.5", "120", "1.2", "0.99", ""},
		{"0.40", "11.8", "60", "2.1", "1.01", ""},
		{"", "", "", "", "", "12.5"},
		{"", "", "", "", "", "45.3"},
	}
	perTanggal := []string{"", time.Now().Format("2006-01-02"), time.Now().Format("2006-01-02"), ""}

	baris := make([][]string, 0, len(plat))
	for i, p := range plat {
		if i >= len(contoh) {
			break
		}
		row := make([]string, 0, len(kolomHasilUjiTemplate))
		row = append(row, p)
		row = append(row, contoh[i]...)
		row = append(row, perTanggal[i])
		baris = append(baris, row)
	}
	return baris
}

// registeredPlats mengambil plat kendaraan yang belum punya hasil uji supaya
// baris contoh pada template langsung bisa dipakai.
func registeredPlats() []string {
	var kendaraan []models.Kendaraan
	config.DB.Where("id NOT IN (?)", config.DB.Model(&models.HasilUji{}).Select("kendaraan_id")).
		Order("created_at").Limit(maksBarisContoh).Find(&kendaraan)

	plat := make([]string, 0, len(kendaraan))
	for _, k := range kendaraan {
		plat = append(plat, k.PlatNomor)
	}
	return plat
}

func daftarKendaraanTemplate() [][]string {
	var kendaraan []models.Kendaraan
	config.DB.Order("created_at").Limit(maksKendaraanDiTemplate).Find(&kendaraan)

	rows := make([][]string, 0, len(kendaraan))
	for _, k := range kendaraan {
		rows = append(rows, []string{
			k.PlatNomor, k.Kategori, k.Jenis, k.Merek, k.Tipe,
			strconv.Itoa(k.KapasitasMesin), strconv.Itoa(k.TahunPembuatan),
		})
	}
	return rows
}

// ---------------------------------------------------------------------------
// Template Excel (.xlsx)
// ---------------------------------------------------------------------------

const (
	shPanduan      = "Panduan"
	shData         = "Data Hasil Uji"
	shKendaraan    = "Daftar Kendaraan"
	barisHeaderUji = 6
)

func DownloadTemplateHasilUjiExcel(c *gin.Context) {
	f := excelize.NewFile()
	defer f.Close()

	if _, err := f.NewSheet(shPanduan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal membuat template Excel"})
		return
	}
	idxData, err := f.NewSheet(shData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal membuat template Excel"})
		return
	}
	if _, err := f.NewSheet(shKendaraan); err != nil {
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

	// Sheet Panduan
	f.SetCellValue(shPanduan, "A1", "Panduan Import Data Hasil Uji Emisi")
	f.SetCellStyle(shPanduan, "A1", "A1", gayaJudul)
	for i, line := range panduanTemplateHasilUji() {
		f.SetCellValue(shPanduan, fmt.Sprintf("A%d", i+3), line)
	}
	f.SetColWidth(shPanduan, "A", "A", 100)

	// Sheet Data Hasil Uji
	f.SetCellValue(shData, "A1", "Isi data mulai baris 7. Baris contoh (abu-abu) bisa diedit atau dihapus.")
	f.SetCellStyle(shData, "A1", "A1", gayaJudul)
	f.SetCellValue(shData, "A2", "Ganti plat nomor dengan plat nomor yang SUDAH terdaftar, lihat sheet 'Daftar Kendaraan'.")
	f.SetCellStyle(shData, "A2", "A2", gayaJudul)

	for i, nama := range kolomHasilUjiTemplate {
		f.SetCellValue(shData, fmt.Sprintf("%s%d", kolomUjiXLSX[i], barisHeaderUji), nama)
	}
	f.SetCellStyle(shData,
		fmt.Sprintf("A%d", barisHeaderUji),
		fmt.Sprintf("%s%d", kolomUjiXLSX[len(kolomUjiXLSX)-1], barisHeaderUji),
		gayaHeader)

	for i, row := range contohTemplateHasilUji() {
		baris := barisHeaderUji + 1 + i
		for j, nilai := range row {
			f.SetCellValue(shData, fmt.Sprintf("%s%d", kolomUjiXLSX[j], baris), nilai)
		}
		f.SetCellStyle(shData,
			fmt.Sprintf("A%d", baris),
			fmt.Sprintf("%s%d", kolomUjiXLSX[len(kolomUjiXLSX)-1], baris),
			gayaContoh)
	}

	for i, lebar := range []float64{22, 12, 12, 12, 12, 12, 14, 16} {
		f.SetColWidth(shData, kolomUjiXLSX[i], kolomUjiXLSX[i], lebar)
	}
	f.SetRowHeight(shData, barisHeaderUji, 24)
	kunciPanes(f, shData, barisHeaderUji)
	_ = f.AutoFilter(shData, fmt.Sprintf("A%d:%s%d", barisHeaderUji, kolomUjiXLSX[len(kolomUjiXLSX)-1], barisHeaderUji), nil)

	// Sheet Daftar Kendaraan (referensi, tidak ikut diimpor)
	tulisDaftarKendaraan(f, gayaJudul, gayaHeader)

	f.SetActiveSheet(idxData)

	buf, err := f.WriteToBuffer()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal membuat template Excel"})
		return
	}

	c.Header("Content-Disposition", `attachment; filename=template_import_hasil_uji.xlsx`)
	c.Data(http.StatusOK, mimeTypeXLSX, buf.Bytes())
}

func tulisDaftarKendaraan(f *excelize.File, gayaJudul, gayaHeader int) {
	kolom := []string{"plat_nomor", "kategori", "jenis", "merek", "tipe", "kapasitas_mesin", "tahun_pembuatan", "sudah_ada_hasil_uji"}

	var kendaraan []models.Kendaraan
	config.DB.Order("created_at").Limit(maksKendaraanDiTemplate).Find(&kendaraan)

	adaUji := map[string]bool{}
	var hasilUji []models.HasilUji
	config.DB.Select("kendaraan_id").Find(&hasilUji)
	for _, h := range hasilUji {
		adaUji[h.KendaraanID] = true
	}

	f.SetCellValue(shKendaraan, "A1", fmt.Sprintf(
		"Kendaraan yang saat ini terdaftar (%d). Salin plat_nomor dari sini ke sheet 'Data Hasil Uji'.", len(kendaraan)))
	f.SetCellStyle(shKendaraan, "A1", "A1", gayaJudul)

	for i, nama := range kolom {
		f.SetCellValue(shKendaraan, fmt.Sprintf("%s3", kolomUjiXLSX[i]), nama)
		f.SetColWidth(shKendaraan, kolomUjiXLSX[i], kolomUjiXLSX[i], 18)
	}
	f.SetCellStyle(shKendaraan, "A3", fmt.Sprintf("%s3", kolomUjiXLSX[len(kolom)-1]), gayaHeader)

	for i, k := range kendaraan {
		ada := "Tidak"
		if adaUji[k.ID] {
			ada = "Ya"
		}
		values := []interface{}{k.PlatNomor, k.Kategori, k.Jenis, k.Merek, k.Tipe,
			k.KapasitasMesin, k.TahunPembuatan, ada}
		for j, v := range values {
			f.SetCellValue(shKendaraan, fmt.Sprintf("%s%d", kolomUjiXLSX[j], 4+i), v)
		}
	}

	kunciPanes(f, shKendaraan, 3)
	_ = f.AutoFilter(shKendaraan, fmt.Sprintf("A3:%s3", kolomUjiXLSX[len(kolom)-1]), nil)
}

func kunciPanes(f *excelize.File, sheet string, baris int) {
	atas := baris + 1
	f.SetPanes(sheet, &excelize.Panes{
		Freeze: true, Split: false, XSplit: 0, YSplit: baris,
		TopLeftCell: fmt.Sprintf("A%d", atas), ActivePane: "bottomLeft",
		Selection: []excelize.Selection{{SQRef: fmt.Sprintf("A%d", atas), ActiveCell: fmt.Sprintf("A%d", atas), Pane: "bottomLeft"}},
	})
}

// ---------------------------------------------------------------------------
// Import
// ---------------------------------------------------------------------------

// indeksKolomUji menyimpan posisi setiap kolom; nilai -1 berarti kolom tidak ada
// pada file.
type indeksKolomUji struct {
	PlatNomor, CO, CO2, HC, O2, Lambda, Opasitas, TanggalUji int
}

type hasilBarisUji struct {
	Nomor                   int
	PlatNomor               string
	CO, CO2, HC, O2, Lambda *float64
	Opasitas                *float64
	Tanggal                 *time.Time
	Err                     string
}

func UploadHasilUjiFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Tidak ada file yang diunggah"})
		return
	}

	if int64(file.Size) > int64(maxUploadSize) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Sprintf("Ukuran file maksimal %d MB", int64(maxUploadSize)>>20),
		})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal membuka file"})
		return
	}
	defer src.Close()

	var records []recordCSV
	switch strings.ToLower(fileExt(file.Filename)) {
	case ".csv", ".txt":
		records, err = bacaRecordsCSV(src)
	case ".xlsx", ".xlsm":
		records, err = bacaRecordsXLSX(src, barisTerlihatHeaderUji)
	case ".xls":
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format .xls (Excel 97-2003) belum didukung. Buka di Excel lalu simpan ulang sebagai .xlsx atau .csv.",
		})
		return
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format file tidak didukung. Gunakan .csv atau .xlsx - template bisa diunduh di halaman ini.",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	if len(records) < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "File kosong. Gunakan tombol 'Download Template' untuk melihat contoh format yang benar.",
		})
		return
	}

	barisHeader, petakan := petakanKolomHasilUji(records)
	if barisHeader < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Sprintf(
				"Kolom file tidak dikenali. File harus memuat kolom: %s. Gunakan tombol 'Download Template'.",
				strings.Join(kolomHasilUjiTemplate, ", ")),
		})
		return
	}

	barisData := records[barisHeader+1:]
	if len(barisData) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "File tidak berisi data. Gunakan tombol 'Download Template' untuk melihat contoh format yang benar.",
		})
		return
	}

	// Muat data referensi sekali saja agar tidak query per baris.
	kendaraanByPlat := map[string]models.Kendaraan{}
	var semuaKendaraan []models.Kendaraan
	config.DB.Find(&semuaKendaraan)
	for _, k := range semuaKendaraan {
		kendaraanByPlat[strings.ToUpper(strings.TrimSpace(k.PlatNomor))] = k
	}

	hasilUjiAktif := map[string]bool{}
	adaHasilUjiTerhapus := map[string]bool{}
	var semuaHasilUji []models.HasilUji
	config.DB.Unscoped().Select("kendaraan_id", "deleted_at").Find(&semuaHasilUji)
	for _, h := range semuaHasilUji {
		if h.DeletedAt.Valid {
			adaHasilUjiTerhapus[h.KendaraanID] = true
		} else {
			hasilUjiAktif[h.KendaraanID] = true
		}
	}

	userID, _ := c.Get("user_id")
	ditestor, _ := userID.(string)

	total := len(barisData)
	success, failed, lulus, tidakLulus := 0, 0, 0, 0
	errorsList := []gin.H{}
	platDipakai := make(map[string]int, total)

	for _, rec := range barisData {
		hasil := parseBarisHasilUji(rec.baris, petakan, rec.line)
		if hasil == nil {
			total--
			continue
		}

		tolak := func(pesan string) {
			failed++
			if len(errorsList) < batasMaksErrorRow {
				errorsList = append(errorsList, gin.H{"row": hasil.Nomor, "message": pesan})
			}
		}

		if hasil.Err != "" {
			tolak(hasil.Err)
			continue
		}

		kunci := strings.ToUpper(hasil.PlatNomor)
		if barisLama, duplikat := platDipakai[kunci]; duplikat {
			tolak(fmt.Sprintf("Plat nomor %s duplikat di dalam file (sudah dipakai di baris %d)",
				hasil.PlatNomor, barisLama))
			continue
		}
		platDipakai[kunci] = hasil.Nomor

		kendaraan, terdaftar := kendaraanByPlat[kunci]
		if !terdaftar {
			tolak(fmt.Sprintf("Kendaraan dengan plat nomor %s tidak terdaftar. Daftarkan dulu di halaman 'Input Data Kendaraan'.",
				hasil.PlatNomor))
			continue
		}
		if hasilUjiAktif[kendaraan.ID] {
			tolak(fmt.Sprintf("Kendaraan %s sudah memiliki data hasil uji", hasil.PlatNomor))
			continue
		}

		hasilUji, valid, catatan := evaluateEmission(kendaraan, hasil.CO, hasil.CO2, hasil.HC, hasil.O2, hasil.Opasitas)

		waktuUji := time.Now()
		if hasil.Tanggal != nil {
			waktuUji = *hasil.Tanggal
		}

		// Hasil uji lama yang sudah dihapus (soft delete) tetap menempati unique
		// index kendaraan_id, jadi harus dibersihkan dulu agar bisa diunggah ulang.
		if adaHasilUjiTerhapus[kendaraan.ID] {
			config.DB.Unscoped().
				Where("kendaraan_id = ? AND deleted_at IS NOT NULL", kendaraan.ID).
				Delete(&models.HasilUji{})
		}

		record := models.HasilUji{
			ID:          uuid.New().String(),
			KendaraanID: kendaraan.ID,
			CO:          hasil.CO,
			CO2:         hasil.CO2,
			HC:          hasil.HC,
			O2:          hasil.O2,
			Lambda:      hasil.Lambda,
			Opasitas:    hasil.Opasitas,
			HasilUji:    hasilUji,
			Valid:       valid,
			Catatan:     catatan,
			TestedBy:    ditestor,
			TestedAt:    waktuUji,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		if err := config.DB.Create(&record).Error; err != nil {
			tolak(fmt.Sprintf("Gagal menyimpan ke database: %v", err))
			continue
		}

		hasilUjiAktif[kendaraan.ID] = true
		if valid {
			lulus++
		} else {
			tidakLulus++
		}
		success++
	}

	payload := gin.H{
		"total":       total,
		"success":     success,
		"failed":      failed,
		"lulus":       lulus,
		"tidak_lulus": tidakLulus,
		"errors":      errorsList,
	}
	if failed > len(errorsList) {
		payload["errors_truncated"] = failed - len(errorsList)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("Import selesai: %d tersimpan, %d ditolak", success, failed),
		"data":    payload,
	})
}

// ---------------------------------------------------------------------------
// Pembacaan file
// ---------------------------------------------------------------------------

// bacaRecordsXLSX mengubah sheet yang memuat data menjadi []recordCSV. Nilai sel
// dikembalikan sebagai string, nomor baris mengikuti nomor baris di Excel.
// Predicate `adalahHeader` menentukan sheet mana yang dipilih sehingga fungsi ini
// bisa dipakai ulang oleh import kendaraan.
func bacaRecordsXLSX(src io.Reader, adalahHeader func([]string) bool) ([]recordCSV, error) {
	f, err := excelize.OpenReader(src)
	if err != nil {
		return nil, fmt.Errorf("File Excel tidak dapat dibaca: %v", err)
	}
	defer f.Close()

	namaSheet := pilihSheetData(f, adalahHeader)
	if namaSheet == "" {
		return nil, fmt.Errorf("File Excel tidak memiliki sheet")
	}

	rows, err := f.GetRows(namaSheet)
	if err != nil {
		return nil, fmt.Errorf("Gagal membaca sheet %q: %v", namaSheet, err)
	}

	records := make([]recordCSV, 0, len(rows))
	for i, row := range rows {
		records = append(records, recordCSV{baris: row, line: i + 1})
	}
	return records, nil
}

// pilihSheetData mencari sheet yang paling mungkin berisi data, bukan panduan.
func pilihSheetData(f *excelize.File, adalahHeader func([]string) bool) string {
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return ""
	}
	for _, name := range sheets {
		rows, err := f.GetRows(name)
		if err != nil {
			continue
		}
		for _, row := range rows {
			if adalahHeader(row) {
				return name
			}
		}
	}
	return sheets[0]
}

// barisTerlihatHeaderUji menandai baris sebagai header bila memuat minimal dua
// nama kolom yang dikenal.
func barisTerlihatHeaderUji(row []string) bool {
	ditemukan := 0
	for _, v := range row {
		if _, ok := kanonKolomUji(v); ok {
			ditemukan++
		}
	}
	return ditemukan >= 2
}

// kanonKolomUji mengembalikan nama kanonik template dari sebuah header kolom.
// Input dinormalkan lebih dulu sehingga "Plat Nomor" dan "plat_nomor" dianggap sama.
func kanonKolomUji(raw string) (string, bool) {
	nama := namaKolomNormal(raw)
	for kanon, alias := range aliasKolomUji {
		for _, a := range alias {
			if nama == a || namaKolomNormal(a) == nama {
				return kanon, true
			}
		}
	}
	// Header yang membawa satuan dalam kurung, mis. "CO (%)" atau "HC (ppm)".
	if strings.Contains(nama, "(") {
		if kanon, ok := kanonKolomUji(strings.SplitN(nama, "(", 2)[0]); ok {
			return kanon, true
		}
	}
	return "", false
}

// petakanKolomHasilUji mencari baris header lalu memetakan indeks setiap kolom.
// Mengembalikan -1 bila tidak ada header yang cocok.
func petakanKolomHasilUji(records []recordCSV) (int, indeksKolomUji) {
	batas := len(records)
	if batas > maksBarisPencariHeader {
		batas = maksBarisPencariHeader
	}

	for i := 0; i < batas; i++ {
		if !barisTerlihatHeaderUji(records[i].baris) {
			continue
		}

		posisi := map[string]int{}
		for j, v := range records[i].baris {
			if kanon, ok := kanonKolomUji(v); ok {
				posisi[kanon] = j
			}
		}

		petakan := indeksKolomUji{-1, -1, -1, -1, -1, -1, -1, -1}
		peta := map[string]*int{
			"plat_nomor":  &petakan.PlatNomor,
			"co":          &petakan.CO,
			"co2":         &petakan.CO2,
			"hc":          &petakan.HC,
			"o2":          &petakan.O2,
			"lambda":      &petakan.Lambda,
			"opasitas":    &petakan.Opasitas,
			"tanggal_uji": &petakan.TanggalUji,
		}
		for kolom, ptr := range peta {
			if idx, ok := posisi[kolom]; ok {
				*ptr = idx
			}
		}

		adaParameter := petakan.CO >= 0 || petakan.CO2 >= 0 || petakan.HC >= 0 ||
			petakan.O2 >= 0 || petakan.Lambda >= 0 || petakan.Opasitas >= 0
		if petakan.PlatNomor >= 0 && adaParameter {
			return i, petakan
		}
	}
	return -1, indeksKolomUji{-1, -1, -1, -1, -1, -1, -1, -1}
}

// parseBarisHasilUji memvalidasi dan menormalkan satu baris. Mengembalikan nil
// bila baris kosong sehingga tidak perlu diproses.
func parseBarisHasilUji(record []string, p indeksKolomUji, nomor int) *hasilBarisUji {
	if semuaKosong(record) {
		return nil
	}

	ambil := func(idx int) string {
		if idx < 0 || idx >= len(record) {
			return ""
		}
		return normalisasiTeks(record[idx])
	}

	plat := strings.ToUpper(ambil(p.PlatNomor))
	if plat == "" {
		return &hasilBarisUji{Nomor: nomor, Err: "plat_nomor wajib diisi"}
	}
	if !rePlatNomor.MatchString(plat) {
		return &hasilBarisUji{Nomor: nomor, Err: fmt.Sprintf(
			"plat_nomor %q tidak valid (maks. 20 karakter, huruf/angka/spasi/tanda hubung)",
			bersihkanTampilan(ambil(p.PlatNomor)))}
	}

	hasil := &hasilBarisUji{Nomor: nomor, PlatNomor: plat}

	kolomUji := []struct {
		idx    int
		kolom  string
		batas  batasNilai
		tujuan **float64
	}{
		{p.CO, "co", batasNilaiUji[0], &hasil.CO},
		{p.CO2, "co2", batasNilaiUji[1], &hasil.CO2},
		{p.HC, "hc", batasNilaiUji[2], &hasil.HC},
		{p.O2, "o2", batasNilaiUji[3], &hasil.O2},
		{p.Lambda, "lambda", batasNilaiUji[4], &hasil.Lambda},
		{p.Opasitas, "opasitas", batasNilaiUji[5], &hasil.Opasitas},
	}

	for _, spec := range kolomUji {
		mentah := ambil(spec.idx)
		if mentah == "" {
			continue
		}
		nilai, err := parseAngkaDesimal(mentah)
		if err != nil {
			return &hasilBarisUji{Nomor: nomor, Err: fmt.Sprintf("%s %q bukan angka yang valid (contoh: %s)",
				spec.kolom, bersihkanTampilan(mentah), spec.batas.Contoh)}
		}
		if nilai < spec.batas.Min || nilai > spec.batas.Maks {
			return &hasilBarisUji{Nomor: nomor, Err: fmt.Sprintf("%s %s di luar rentang yang diperbolehkan (%g - %g %s)",
				spec.kolom, strconv.FormatFloat(nilai, 'f', -1, 64),
				spec.batas.Min, spec.batas.Maks, spec.batas.Satuan)}
		}
		v := nilai
		*spec.tujuan = &v
	}

	if hasil.CO == nil && hasil.CO2 == nil && hasil.HC == nil &&
		hasil.O2 == nil && hasil.Lambda == nil && hasil.Opasitas == nil {
		return &hasilBarisUji{Nomor: nomor,
			Err: "minimal satu parameter uji (co/co2/hc/o2/lambda/opasitas) harus diisi"}
	}

	if mentah := ambil(p.TanggalUji); mentah != "" {
		tgl, err := parseTanggalUji(mentah)
		if err != nil {
			return &hasilBarisUji{Nomor: nomor, Err: fmt.Sprintf(
				"tanggal_uji %q tidak valid (format 2024-01-15 atau 15/01/2024)",
				bersihkanTampilan(mentah))}
		}
		if tgl.After(time.Now().AddDate(0, 0, 1)) {
			return &hasilBarisUji{Nomor: nomor, Err: fmt.Sprintf(
				"tanggal_uji %s tidak boleh berada di masa depan", tgl.Format("2006-01-02"))}
		}
		hasil.Tanggal = &tgl
	}

	return hasil
}

// ---------------------------------------------------------------------------
// Parser angka & tanggal
// ---------------------------------------------------------------------------

// parseAngkaDesimal menerima 0.85, "0,85", "12.5%", "12,5 ppm", dan 1.234,5.
// Kelompok tepat tiga digit dibaca sebagai pemisah ribuan, sisanya sebagai desimal.
func parseAngkaDesimal(raw string) (float64, error) {
	s := strings.ToUpper(normalisasiTeks(raw))
	s = strings.TrimSpace(s)
	s = strings.TrimSuffix(s, "PPM")
	s = strings.TrimSuffix(s, "%")
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "+")

	negatif := strings.HasPrefix(s, "-")
	if negatif {
		s = s[1:]
	}
	s = strings.ReplaceAll(s, " ", "")
	if s == "" {
		return 0, strconv.ErrSyntax
	}

	kelompok := rePemisahAngka.Split(s, -1)
	for _, k := range kelompok {
		if !hanyaAngka(k) {
			return 0, strconv.ErrSyntax
		}
	}

	switch {
	case len(kelompok) == 1:
		// tanpa separator sama sekali
	case len(kelompok) >= 3 && selainPertamaTigaDigit(kelompok):
		// seluruh separator adalah pemisah ribuan: 1.234.567 / 1,234,567
		kelompok = []string{strings.Join(kelompok, "")}
	default:
		// separator terakhir adalah desimal, sebelumnya pemisah ribuan.
		// "1.234,5" -> 1234.5, "0,85" -> 0.85, "0.850" -> 0.850.
		depan, belakang := kelompok[:len(kelompok)-1], kelompok[len(kelompok)-1]
		if len(belakang) > 3 {
			return 0, strconv.ErrSyntax
		}
		for i, g := range depan {
			if i > 0 && len(g) != 3 {
				return 0, strconv.ErrSyntax
			}
		}
		kelompok = []string{strings.Join(depan, "") + "." + belakang}
	}

	nilai, err := strconv.ParseFloat(kelompok[0], 64)
	if err != nil {
		return 0, strconv.ErrSyntax
	}
	if negatif {
		nilai = -nilai
	}
	return nilai, nil
}

// selainPertamaTigaDigit memeriksa pola pengelompokan ribuan: grup pertama
// boleh 1-3 digit, grup berikutnya tepat 3 digit.
func selainPertamaTigaDigit(kelompok []string) bool {
	if len(kelompok) < 2 {
		return false
	}
	for i, g := range kelompok {
		if i == 0 {
			if len(g) == 0 || len(g) > 3 {
				return false
			}
			continue
		}
		if len(g) != 3 {
			return false
		}
	}
	return true
}

// parseTanggalUji menerima 2024-01-15, 2024/01/15, 15/01/2024, dan nomor seri Excel.
// String kosong menghasilkan waktu zero yang berarti "otomatis".
func parseTanggalUji(raw string) (time.Time, error) {
	s := strings.TrimSpace(normalisasiTeks(raw))
	if s == "" {
		return time.Time{}, nil
	}
	if idx := strings.IndexAny(s, " T"); idx > 0 {
		s = s[:idx]
	}

	if hanyaAngka(s) {
		if n, err := strconv.ParseFloat(s, 64); err == nil && n > 20000 && n < 60000 {
			if tgl, err := excelize.ExcelDateToTime(n, false); err == nil {
				return time.Date(tgl.Year(), tgl.Month(), tgl.Day(), 0, 0, 0, 0, time.Local), nil
			}
		}
		return time.Time{}, strconv.ErrSyntax
	}

	if m := reTanggalISO.FindStringSubmatch(s); m != nil {
		return waktuDariParts(m[1], m[2], m[3])
	}
	if m := reTanggalYMD.FindStringSubmatch(s); m != nil {
		return waktuDariParts(m[1], m[2], m[3])
	}
	if m := reTanggalDMY.FindStringSubmatch(s); m != nil {
		return waktuDariParts(m[3], m[2], m[1])
	}
	return time.Time{}, strconv.ErrSyntax
}

func waktuDariParts(tahun, bulan, hari string) (time.Time, error) {
	y, err := strconv.Atoi(tahun)
	if err != nil {
		return time.Time{}, strconv.ErrSyntax
	}
	bl, err := strconv.Atoi(bulan)
	if err != nil {
		return time.Time{}, strconv.ErrSyntax
	}
	h, err := strconv.Atoi(hari)
	if err != nil || bl < 1 || bl > 12 || h < 1 || h > 31 || y < 1900 || y > 2100 {
		return time.Time{}, strconv.ErrSyntax
	}

	tgl := time.Date(y, time.Month(bl), h, 0, 0, 0, 0, time.Local)
	if tgl.Year() != y || tgl.Month() != time.Month(bl) || tgl.Day() != h {
		return time.Time{}, strconv.ErrSyntax
	}
	return tgl, nil
}

func fileExt(nama string) string {
	if idx := strings.LastIndex(nama, "."); idx >= 0 {
		return nama[idx:]
	}
	return ""
}
