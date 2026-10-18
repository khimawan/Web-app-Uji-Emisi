package controllers

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"web-app-uji-emisi/config"
	"web-app-uji-emisi/models"
)

// Nilai yang valid untuk kolom kategorikal kendaraan. Template CSV dan form
// webapp harus memakai daftar ini supaya tidak ada data yang tidak konsisten.
var (
	KategoriKendaraanValid = []string{"Dinas", "Umum"}
	JenisKendaraanValid    = []string{"Bensin M", "Bensin N&O", "Solar JBB", "Solar GVW"}
)

// Kolom template CSV, urutan ini yang dipakai saat header tidak dikenali.
var kolomKendaraanTemplate = []string{
	"kategori", "jenis", "plat_nomor", "merek", "tipe", "kapasitas_mesin", "tahun_pembuatan",
}

var (
	rePlatNomor       = regexp.MustCompile(`^[A-Z0-9][A-Z0-9 .\-/]{0,19}$`)
	rePemisahAngka    = regexp.MustCompile(`[.,]`)
	batasMaksErrorRow = 50
)

const maxUploadSize = 5 << 20 // 5 MB

// ---------------------------------------------------------------------------
// Template CSV
// ---------------------------------------------------------------------------

func DownloadCSVTemplate(c *gin.Context) {
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", `attachment; filename=template_kendaraan.csv`)

	w := csv.NewWriter(c.Writer)
	defer w.Flush()

	// Baris panduan: diawali "#", otomatis dilewati oleh import (dan Excel).
	for _, line := range panduanTemplateKendaraan() {
		fmt.Fprintf(c.Writer, "# %s\n", line)
	}

	// Baris header - JANGAN diubah namanya, itu yang dibaca import.
	_ = w.Write(kolomKendaraanTemplate)

	// Baris contoh: satu untuk tiap nilai kategori + jenis yang tersedia.
	// Hapus semua baris contoh sebelum mengunggah data asli.
	for _, row := range contohTemplateKendaraan() {
		_ = w.Write(row)
	}
}

func panduanTemplateKendaraan() []string {
	return []string{
		"TEMPLATE INPUT DATA KENDARAAN - APLIKASI UJI EMISI",
		"",
		"CARA PAKAI:",
		"  1. Unduh template di halaman Input Data Kendaraan: Template CSV atau Template Excel.",
		"  2. Isi data mulai baris paling bawah (baris contoh), 1 baris = 1 kendaraan.",
		"  3. Hapus semua baris contoh sebelum mengunggah.",
		"  4. Simpan kembali dengan format yang sama:",
		"       - Template CSV   -> .csv dengan encoding UTF-8 (Save As > CSV UTF-8).",
		"       - Template Excel -> .xlsx atau .xlsm.",
		"     Catatan: .xls (Excel 97-2003) belum didukung, simpan ulang sebagai .xlsx atau .csv.",
		"  5. Unggah lewat tombol Upload File pada halaman Input Data Kendaraan.",
		"     Nama kolom boleh diacak urutannya dan file hasil salin-tempel dari Excel",
		"     tetap terbaca karena baris header dicari otomatis.",
		"",
		"ATURAN FORMAT:",
		fmt.Sprintf("  - kategori     : salah satu dari [%s]. Huruf besar/kecil bebas.",
			strings.Join(KategoriKendaraanValid, " / ")),
		fmt.Sprintf("  - jenis        : salah satu dari [%s].", strings.Join(JenisKendaraanValid, " / ")),
		"                   Bensin M  = Angkutan Penumpang (motor & mobil <=3.500 cc)",
		"                   Bensin N&O = Angkutan Barang / Gendang (di atas 3.500 cc)",
		"                   Solar JBB  = Solar, berat total kendaraan <= 3,5 ton",
		"                   Solar GVW  = Solar, berat total kendaraan >  3,5 ton",
		"  - plat_nomor   : huruf kapital, tanpa tanda baca. Contoh: B 1234 CDE, B1234CDE, DG 1234 AB",
		"  - merek        : nama pabrikan. Contoh: Toyota, Mitsubishi, Isuzu, Hino",
		"  - tipe         : tipe/series kendaraan. Contoh: Avanza, D-Max, Traga, FM 440",
		"  - kapasitas_mesin : ANGKA BULAT (cc). Tulis 1496, jangan '1496 cc' / '1.496' / '1,5'.",
		"  - tahun_pembuatan : ANGKA BULAT 4 digit. Tulis 2020, jangan '2020.0' / 'th 2020'.",
		"",
		"CATATAN:",
		"  - Pemisah kolom bisa koma (,), titik koma (;), tab, atau pipe (|). Otomatis dideteksi.",
		"  - Baris kosong dan baris berawalan '#' otomatis dilewati.",
		"  - Nilai numerik boleh pakai pemisah ribuan (1.496) atau desimal (1496,00),",
		"    nilainya tetap dibaca sebagai 1496.",
		"  - Baris yang formatnya salah akan ditolak dan nomor baris + alasannya",
		"    ditampilkan pada hasil unggah. Baris lain tetap tersimpan.",
	}
}

func contohTemplateKendaraan() [][]string {
	return [][]string{
		{"Umum", "Bensin M", "B 1234 CDE", "Toyota", "Avanza", "1329", "2020"},
		{"Umum", "Bensin M", "B 1234 CDF", "Honda", "Brio", "1199", "2021"},
		{"Umum", "Bensin N&O", "B 1234 CDG", "Isuzu", "Traga", "2477", "2019"},
		{"Dinas", "Bensin N&O", "B 1234 CDH", "Toyota", "Hilux D-Cab", "2755", "2022"},
		{"Umum", "Solar JBB", "B 1234 CDI", "Mitsubishi", "L300", "2477", "2020"},
		{"Umum", "Solar GVW", "B 1234 CDJ", "Hino", "FM 440", "12916", "2018"},
	}
}

// ---------------------------------------------------------------------------
// Import CSV
// ---------------------------------------------------------------------------

// recordCSV menyimpan satu baris CSV beserta nomor baris fisiknya di file.
type recordCSV struct {
	baris []string
	line  int
}

type hasilBarisCSV struct {
	Nomor     int
	Kendaraan models.Kendaraan
	Err       string
}

// UploadKendaraanFile menerima template CSV maupun Excel (.xlsx/.xlsm).
// Ekstensi lain ditolak dengan pesan yang jelas agar pengguna tidak salah pilih file.
func UploadKendaraanFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Tidak ada file yang diunggah"})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal membuka file"})
		return
	}
	defer src.Close()

	pembatas := int64(maxUploadSize)
	if file.Size > pembatas {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Sprintf("Ukuran file maksimal %d MB", pembatas>>20),
		})
		return
	}

	ext := fileExt(file.Filename)
	var records []recordCSV
	switch ext {
	case ".csv", ".txt":
		records, err = bacaRecordsCSV(src)
	case ".xlsx", ".xlsm":
		records, err = bacaRecordsXLSX(src, barisTerlihatHeader)
	case ".xls":
		err = errors.New("Format .xls (Excel 97-2003) belum didukung. Buka di Excel lalu simpan ulang sebagai .xlsx atau .csv")
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Sprintf("Format %s belum didukung. Gunakan .csv, .xlsx, atau .xlsm.", ext),
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
			"message": "File kosong. Gunakan tombol 'Download Template' untuk contoh format yang benar.",
		})
		return
	}

	// Baris header dicari di beberapa baris pertama supaya file hasil salin-tempel
	// dari Excel (yang biasanya punya baris judul) tetap bisa dibaca.
	idxHeader, petakan, err := cariBarisHeaderKendaraan(records)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	barisData := records
	if idxHeader >= 0 {
		barisData = records[idxHeader+1:]
	} else if len(records[0].baris) < len(kolomKendaraanTemplate) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Sprintf("File harus memuat %d kolom: %s. Gunakan tombol 'Download Template'.",
				len(kolomKendaraanTemplate), strings.Join(kolomKendaraanTemplate, ", ")),
		})
		return
	}

	if len(barisData) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "File tidak berisi data. Gunakan tombol 'Download Template' untuk contoh format yang benar.",
		})
		return
	}

	userID, _ := c.Get("user_id")
	createdBy, _ := userID.(string)

	total := len(barisData)
	success := 0
	failed := 0
	errorsList := []gin.H{}
	platDipakai := make(map[string]int, total)

	// Plat nomor yang sudah terdaftar dimuat sekali agar duplikat ditolak dengan
	// pesan yang jelas, bukan error mentah dari database.
	platTerdaftar := map[string]bool{}
	var platRegistered []models.Kendaraan
	config.DB.Select("plat_nomor").Find(&platRegistered)
	for _, k := range platRegistered {
		platTerdaftar[strings.ToUpper(strings.TrimSpace(k.PlatNomor))] = true
	}

	for _, rec := range barisData {
		nomorBaris := rec.line

		hasil := parseBarisKendaraan(rec.baris, petakan, nomorBaris)
		if hasil == nil {
			// Baris kosong sepenuhnya: lewati tanpa dihitung sebagai gagal.
			total--
			continue
		}

		if hasil.Err != "" {
			failed++
			if len(errorsList) < batasMaksErrorRow {
				errorsList = append(errorsList, gin.H{
					"row":     nomorBaris,
					"message": hasil.Err,
				})
			}
			continue
		}

		if barisLama, ada := platDipakai[hasil.Kendaraan.PlatNomor]; ada {
			failed++
			if len(errorsList) < batasMaksErrorRow {
				errorsList = append(errorsList, gin.H{
					"row":     nomorBaris,
					"message": fmt.Sprintf("Plat nomor %s duplikat (sudah dipakai di baris %d)", hasil.Kendaraan.PlatNomor, barisLama),
				})
			}
			continue
		}
		if platTerdaftar[hasil.Kendaraan.PlatNomor] {
			failed++
			if len(errorsList) < batasMaksErrorRow {
				errorsList = append(errorsList, gin.H{
					"row":     nomorBaris,
					"message": fmt.Sprintf("Plat nomor %s sudah terdaftar. Hapus dari file ini atau edit di halaman Input Data Kendaraan.", hasil.Kendaraan.PlatNomor),
				})
			}
			platDipakai[hasil.Kendaraan.PlatNomor] = nomorBaris
			continue
		}
		platDipakai[hasil.Kendaraan.PlatNomor] = nomorBaris
		platTerdaftar[hasil.Kendaraan.PlatNomor] = true

		kendaraan := hasil.Kendaraan
		kendaraan.ID = uuid.New().String()
		kendaraan.CreatedBy = createdBy
		kendaraan.CreatedAt = time.Now()
		kendaraan.UpdatedAt = time.Now()

		if err := config.DB.Create(&kendaraan).Error; err != nil {
			failed++
			if len(errorsList) < batasMaksErrorRow {
				errorsList = append(errorsList, gin.H{
					"row":     nomorBaris,
					"message": fmt.Sprintf("Gagal menyimpan ke database: %v", err),
				})
			}
			continue
		}
		success++
	}

	payload := gin.H{
		"total":   total,
		"success": success,
		"failed":  failed,
		"errors":  errorsList,
	}
	if failed > len(errorsList) {
		payload["errors_truncated"] = failed - len(errorsList)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("Data kendaraan berhasil diunggah: %d tersimpan, %d ditolak", success, failed),
		"data":    payload,
	})
}

// bacaRecordsCSV mendeteksi pemisah kolom, membuang BOM, serta melewati baris
// kosong dan baris komentar (awalan '#').
func bacaRecordsCSV(src io.Reader) ([]recordCSV, error) {
	head := make([]byte, 4096)
	n, err := src.Read(head)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("Gagal membaca file: %v", err)
	}
	head = head[:n]
	head = bytes.TrimPrefix(head, []byte{0xEF, 0xBB, 0xBF})

	r := csv.NewReader(io.MultiReader(bytes.NewReader(head), src))
	r.Comment = '#'
	r.FieldsPerRecord = -1
	r.LazyQuotes = true
	r.TrimLeadingSpace = true
	r.ReuseRecord = false
	r.Comma = deteksiPemisahCSV(string(head))

	records := []recordCSV{}
	line := 0
	for {
		rec, err := r.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("Format CSV tidak valid di sekitar baris %d: %v", line+1, err)
		}
		line, _ = r.FieldPos(0)
		records = append(records, recordCSV{baris: rec, line: line})
	}
	return records, nil
}

// deteksiPemisahCSV memilih koma, titik koma, atau tab sesuai baris pertama
// yang berisi data - umum terjadi saat file disimpan dari Excel locale lokal.
func deteksiPemisahCSV(head string) rune {
	kandidat := []rune{',', ';', '\t', '|'}
	terbaik := ','
	terbanyak := 0

	for _, k := range kandidat {
		if n := strings.Count(head, string(k)); n > terbanyak {
			terbanyak, terbaik = n, k
		}
	}
	return terbaik
}

// barisTerlihatHeader menandai apakah baris pertama berisi nama kolom.
func barisTerlihatHeader(row []string) bool {
	ditemukan := 0
	for _, v := range row {
		switch namaKolomNormal(v) {
		case "kategori", "jenis", "plat_nomor", "merek", "tipe", "kapasitas_mesin", "tahun_pembuatan":
			ditemukan++
		}
	}
	return ditemukan >= 2
}

type indeksKolom struct {
	Kategori, Jenis, PlatNomor, Merek, Tipe int
	KapasitasMesin, TahunPembuatan          int
}

// cariBarisHeaderKendaraan mencari baris header pada baris-baris pertama file.
// Mengembalikan indeks baris header (-1 bila file tidak punya header) beserta
// pemetaan kolomnya.
func cariBarisHeaderKendaraan(records []recordCSV) (int, indeksKolom, error) {
	posisi := indeksKolom{0, 1, 2, 3, 4, 5, 6}
	for i := 0; i < len(records) && i < maksBarisPencariHeader; i++ {
		if !barisTerlihatHeader(records[i].baris) {
			continue
		}
		petakan, err := petakanKolomKendaraan(records[i].baris)
		if err != nil {
			return -1, posisi, err
		}
		return i, petakan, nil
	}
	return -1, posisi, nil
}

// aliasKolomKendaraan memetakan nama kolom kanon ke nama lain yang lazim
// dipakai pengguna, termasuk versi ber-satuan seperti "Kapasitas (cc)".
var aliasKolomKendaraan = map[string][]string{
	"kategori": {"kategori", "kategori_kendaraan"},
	"jenis":    {"jenis", "jenis_kendaraan"},
	"plat_nomor": {
		"plat_nomor", "platnomor", "no_polisi", "nomor_polisi", "nopol",
		"kendaraan", "no_kendaraan",
	},
	"merek":           {"merek", "pabrikan", "make", "merek_kendaraan"},
	"tipe":            {"tipe", "type", "series", "model", "tipe_kendaraan"},
	"kapasitas_mesin": {"kapasitas_mesin", "kapasitas", "kapasitas_mesin_cc", "mesin", "cc"},
	"tahun_pembuatan": {"tahun_pembuatan", "tahun", "thn", "th", "tahun_kei"},
}

// kanonKolomKendaraan mengenali nama kolom: alias persis lebih dulu, lalu nama
// dengan satuan dalam kurung dibuang ("Kapasitas (cc)" -> "kapasitas_mesin").
func kanonKolomKendaraan(raw string) string {
	return kanonKolomKendaraanNormal(namaKolomNormal(raw), 0)
}

func kanonKolomKendaraanNormal(nama string, kedalaman int) string {
	if kedalaman > 2 {
		return nama
	}
	for kanon, alias := range aliasKolomKendaraan {
		for _, a := range alias {
			if nama == a {
				return kanon
			}
		}
	}
	if i := strings.Index(nama, "("); i >= 0 {
		potong := strings.Trim(strings.SplitN(nama, "(", 2)[0], "_")
		if potong != "" {
			return kanonKolomKendaraanNormal(potong, kedalaman+1)
		}
	}
	return nama
}

// petakanKolomKendaraan mencocokkan nama kolom pada baris header. Bila semua
// nama dikenali, urutan kolom bebas; bila tidak, diasumsikan mengikuti template.
func petakanKolomKendaraan(header []string) (indeksKolom, error) {
	normalisasi := make([]string, len(header))
	for i, h := range header {
		normalisasi[i] = kanonKolomKendaraan(h)
	}

	cari := func(target string) (int, bool) {
		for i, h := range normalisasi {
			if h == target {
				return i, true
			}
		}
		return 0, false
	}

	var p indeksKolom
	petakan := map[string]*int{
		"kategori":        &p.Kategori,
		"jenis":           &p.Jenis,
		"plat_nomor":      &p.PlatNomor,
		"merek":           &p.Merek,
		"tipe":            &p.Tipe,
		"kapasitas_mesin": &p.KapasitasMesin,
		"tahun_pembuatan": &p.TahunPembuatan,
	}

	lengkap := true
	for nama, ptr := range petakan {
		if idx, ada := cari(nama); ada {
			*ptr = idx
		} else {
			lengkap = false
		}
	}
	if lengkap {
		return p, nil
	}

	if len(header) == len(kolomKendaraanTemplate) {
		return indeksKolom{0, 1, 2, 3, 4, 5, 6}, nil
	}

	return indeksKolom{}, fmt.Errorf(
		"Kolom file tidak sesuai. File harus memuat kolom: %s. Gunakan tombol 'Download Template'.",
		strings.Join(kolomKendaraanTemplate, ", "))
}

// namaKolomNormal menormalkan nama kolom agar pencocokan tidak bergantung pada
// huruf besar, spasi, tanda hubung, atau titik (mis. "No. Polisi" -> "no_polisi").
func namaKolomNormal(s string) string {
	s = strings.ToLower(normalisasiTeks(s))
	for _, ganti := range []string{" ", "-", ".", ",", ":"} {
		s = strings.ReplaceAll(s, ganti, "_")
	}
	for strings.Contains(s, "__") {
		s = strings.ReplaceAll(s, "__", "_")
	}
	return strings.Trim(s, "_")
}

// parseBarisKendaraan memvalidasi & menormalkan satu baris. Mengembalikan nil
// bila baris kosong dan tidak perlu diproses.
func parseBarisKendaraan(record []string, p indeksKolom, nomor int) *hasilBarisCSV {
	ambil := func(idx int) string {
		if idx < 0 || idx >= len(record) {
			return ""
		}
		return record[idx]
	}

	if semuaKosong(record) {
		return nil
	}

	ada := func(v string) bool { return normalisasiTeks(v) != "" }

	// --- kategori ---
	var kategori string
	if !ada(ambil(p.Kategori)) {
		return barisGagal(nomor, "kategori wajib diisi")
	}
	var ok bool
	if kategori, ok = canonicalKategoriKendaraan(ambil(p.Kategori)); !ok {
		return barisGagal(nomor, fmt.Sprintf("kategori %q tidak dikenal, pilih salah satu: %s",
			bersihkanTampilan(ambil(p.Kategori)), strings.Join(KategoriKendaraanValid, " / ")))
	}

	// --- jenis ---
	if !ada(ambil(p.Jenis)) {
		return barisGagal(nomor, "jenis wajib diisi")
	}
	var jenis string
	if jenis, ok = canonicalJenisKendaraan(ambil(p.Jenis)); !ok {
		return barisGagal(nomor, fmt.Sprintf("jenis %q tidak dikenal, pilih salah satu: %s",
			bersihkanTampilan(ambil(p.Jenis)), strings.Join(JenisKendaraanValid, " / ")))
	}

	// --- plat nomor ---
	plat := strings.ToUpper(normalisasiTeks(ambil(p.PlatNomor)))
	if plat == "" {
		return barisGagal(nomor, "plat_nomor wajib diisi")
	}
	if !rePlatNomor.MatchString(plat) {
		return barisGagal(nomor, fmt.Sprintf("plat_nomor %q tidak valid (maks. 20 karakter, huruf/angka/ spasi/tanda hubung)",
			bersihkanTampilan(ambil(p.PlatNomor))))
	}

	// --- merek ---
	merek := normalisasiTeks(ambil(p.Merek))
	if merek == "" {
		return barisGagal(nomor, "merek wajib diisi")
	}
	if len([]rune(merek)) > 100 {
		return barisGagal(nomor, "merek maksimal 100 karakter")
	}

	// --- tipe ---
	tipe := normalisasiTeks(ambil(p.Tipe))
	if tipe == "" {
		return barisGagal(nomor, "tipe wajib diisi")
	}
	if len([]rune(tipe)) > 100 {
		return barisGagal(nomor, "tipe maksimal 100 karakter")
	}

	// --- kapasitas mesin ---
	kapasitas, err := parseAngkaBulat(ambil(p.KapasitasMesin))
	if err != nil {
		return barisGagal(nomor, fmt.Sprintf("kapasitas_mesin %q bukan angka bulat yang valid (satuan cc, contoh: 1496)",
			bersihkanTampilan(ambil(p.KapasitasMesin))))
	}
	if kapasitas < 1 || kapasitas > 100000 {
		return barisGagal(nomor, fmt.Sprintf("kapasitas_mesin %d di luar rentang yang diperbolehkan (1 - 100.000 cc)", kapasitas))
	}

	// --- tahun pembuatan ---
	tahun, err := parseAngkaBulat(ambil(p.TahunPembuatan))
	if err != nil {
		return barisGagal(nomor, fmt.Sprintf("tahun_pembuatan %q bukan angka bulat yang valid (contoh: 2020)",
			bersihkanTampilan(ambil(p.TahunPembuatan))))
	}
	batasTahun := time.Now().Year() + 2
	if tahun < 1900 || tahun > batasTahun {
		return barisGagal(nomor, fmt.Sprintf("tahun_pembuatan %d di luar rentang yang diperbolehkan (1900 - %d)", tahun, batasTahun))
	}

	return &hasilBarisCSV{
		Nomor: nomor,
		Kendaraan: models.Kendaraan{
			Kategori:       kategori,
			Jenis:          jenis,
			PlatNomor:      plat,
			Merek:          merek,
			Tipe:           tipe,
			KapasitasMesin: kapasitas,
			TahunPembuatan: tahun,
		},
	}
}

func barisGagal(nomor int, pesan string) *hasilBarisCSV {
	return &hasilBarisCSV{Nomor: nomor, Err: pesan}
}

// ---------------------------------------------------------------------------
// Normalisasi & validasi nilai
// ---------------------------------------------------------------------------

func normalisasiTeks(s string) string {
	s = strings.ReplaceAll(s, "\u00a0", " ")
	return strings.Join(strings.Fields(s), " ")
}

func semuaKosong(record []string) bool {
	for _, v := range record {
		if normalisasiTeks(v) != "" {
			return false
		}
	}
	return true
}

func bersihkanTampilan(s string) string {
	const maks = 40
	s = normalisasiTeks(s)
	if len([]rune(s)) > maks {
		return string([]rune(s)[:maks]) + "..."
	}
	return s
}

func kunciKategori(s string) string {
	k := strings.ToUpper(normalisasiTeks(s))
	k = strings.ReplaceAll(k, "_", " ")
	return strings.Join(strings.Fields(k), " ")
}

func canonicalKategoriKendaraan(v string) (string, bool) {
	kunci := kunciKategori(v)
	for _, k := range KategoriKendaraanValid {
		if kunci == strings.ToUpper(k) {
			return k, true
		}
	}
	return "", false
}

// canonicalJenisKendaraan menerima nilai persis maupun label panjang dari UI
// (mis. "Solar JBB (<=3,5 ton)") dan mengembalikannya ke nilai kanonik.
func canonicalJenisKendaraan(v string) (string, bool) {
	kunci := kunciKategori(v)
	kunci = strings.ReplaceAll(kunci, " & ", "&")
	kunci = strings.ReplaceAll(kunci, " AND ", "&")

	for _, j := range JenisKendaraanValid {
		kanon := strings.ToUpper(j)
		if kunci == kanon {
			return j, true
		}
	}
	for _, j := range JenisKendaraanValid {
		kanon := strings.ToUpper(j)
		if len(kunci) > len(kanon) && strings.HasPrefix(kunci, kanon) {
			lanjutan := kunci[len(kanon):]
			if !hurufAngka(rune(lanjutan[0])) {
				return j, true
			}
		}
	}
	return "", false
}

func hurufAngka(r rune) bool {
	return (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9')
}

func hanyaAngka(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// parseAngkaBulat menerima 1496, "1496 cc", 1.496, 1.234.567, 1496,00, dan
// 2020.0. Kelompok 3 digit di belakang separator dibaca sebagai pemisah ribuan,
// sedangkan kelompok 1-2 digit dibaca sebagai desimal dan diabaikan, karena
// kedua penulisan sama lazim dipakai di ekspor Excel/spreadsheet.
func parseAngkaBulat(raw string) (int, error) {
	s := strings.ToUpper(normalisasiTeks(raw))
	s = strings.TrimSpace(strings.TrimSuffix(s, "CC"))
	s = strings.ReplaceAll(s, " ", "")
	if s == "" {
		return 0, strconv.ErrSyntax
	}

	kelompok := rePemisahAngka.Split(s, -1)
	if len(kelompok) == 0 {
		return 0, strconv.ErrSyntax
	}
	for _, k := range kelompok {
		if !hanyaAngka(k) {
			return 0, strconv.ErrSyntax
		}
	}

	terakhir := kelompok[len(kelompok)-1]
	switch {
	case len(kelompok) == 1:
		s = kelompok[0]
	case len(terakhir) == 3:
		// Pemisah ribuan: 1.496 / 1.234.567 -> 1496 / 1234567
		s = strings.Join(kelompok, "")
	case len(kelompok) == 2 && len(terakhir) <= 2:
		// Desimal: 1496,0 / 2020.0 -> 1496 / 2020
		s = kelompok[0]
	default:
		return 0, strconv.ErrSyntax
	}

	return strconv.Atoi(s)
}
