package controllers

import (
	"bytes"
	"testing"
	"time"

	"github.com/xuri/excelize/v2"
)

func TestParseAngkaDesimal(t *testing.T) {
	harusNilai := map[string]float64{
		"0.85":       0.85,
		"0,85":       0.85,
		"12.5":       12.5,
		"12,5":       12.5,
		"12.5%":      12.5,
		" 12,5 ppm ": 12.5,
		"120":        120,
		"0":          0,
		"1.234,5":    1234.5,
		"1,234.5":    1234.5,
		"1.234.567":  1234567,
		"1,234,567":  1234567,
		"0.850":      0.85,
		"1.234":      1.234,
		"+5.25":      5.25,
	}

	for mentah, urous := range harusNilai {
		hasil, err := parseAngkaDesimal(mentah)
		if err != nil {
			t.Fatalf("parseAngkaDesimal(%q) error: %v", mentah, err)
		}
		if hasil != urous {
			t.Errorf("parseAngkaDesimal(%q) = %v, ingin %v", mentah, hasil, urous)
		}
	}

	harusGagal := []string{"", "   ", "abc", "1.2.3.4", "0,85,2", "12 ppm %", "1..2", "-", "NaN"}
	for _, mentah := range harusGagal {
		if _, err := parseAngkaDesimal(mentah); err == nil {
			t.Errorf("parseAngkaDesimal(%q) seharusnya gagal", mentah)
		}
	}
}

func TestParseTanggalUji(t *testing.T) {
	harus, err := parseTanggalUji("2024-01-15")
	if err != nil || harus.Year() != 2024 || harus.Month() != time.January || harus.Day() != 15 {
		t.Fatalf("ISO: %v (err %v)", harus, err)
	}

	harus, err = parseTanggalUji("15/01/2024")
	if err != nil || harus.Year() != 2024 || harus.Month() != time.January || harus.Day() != 15 {
		t.Fatalf("DMY: %v (err %v)", harus, err)
	}

	harus, err = parseTanggalUji("15-01-2024")
	if err != nil || harus.Day() != 15 || harus.Month() != time.January {
		t.Fatalf("DMY dash: %v (err %v)", harus, err)
	}

	harus, err = parseTanggalUji("2024-01-15 08:30:00")
	if err != nil || harus.Day() != 15 {
		t.Fatalf("dengan waktu: %v (err %v)", harus, err)
	}

	// Nomor seri Excel untuk 2024-01-15.
	harus, err = parseTanggalUji("45306")
	if err != nil || harus.Year() != 2024 || harus.Month() != time.January || harus.Day() != 15 {
		t.Fatalf("serial excel: %v (err %v)", harus, err)
	}

	if kosong, err := parseTanggalUji(""); err != nil || !kosong.IsZero() {
		t.Fatalf("kosong harus menghasilkan waktu zero, dapat %v (err %v)", kosong, err)
	}

	for _, salah := range []string{"abc", "2024-13-01", "32/01/2024", "31/02/2024", "1", "2024"} {
		if _, err := parseTanggalUji(salah); err == nil {
			t.Errorf("parseTanggalUji(%q) seharusnya gagal", salah)
		}
	}
}

func TestPetakanKolomHasilUji(t *testing.T) {
	records := []recordCSV{
		{baris: []string{"LAPORAN UJI"}, line: 1},
		{baris: []string{"plat_nomor", "co2", "co", "opasitas", "tanggal_uji"}, line: 2},
		{baris: []string{"B 1 AA", "1", "2", "3", "2024-01-15"}, line: 3},
	}

	baris, p := petakanKolomHasilUji(records)
	if baris != 1 {
		t.Fatalf("baris header = %d, ingin 1", baris)
	}
	if p.PlatNomor != 0 || p.CO != 2 || p.CO2 != 1 || p.Opasitas != 3 || p.TanggalUji != 4 {
		t.Fatalf("pemetaan kolom salah: %+v", p)
	}
	if p.HC != -1 || p.O2 != -1 || p.Lambda != -1 {
		t.Fatalf("kolom yang tidak ada harus -1: %+v", p)
	}

	// Alias penamaan umum harus dikenali.
	if _, ok := kanonKolomUji("no_polisi"); !ok {
		t.Error("alias no_polisi tidak dikenali")
	}
	if kanon, _ := kanonKolomUji("Plat Nomor"); kanon != "plat_nomor" {
		t.Errorf("'Plat Nomor' -> %q, ingin plat_nomor", kanon)
	}

	// Tanpa kolom plat_nomor harus ditolak.
	if _, p := petakanKolomHasilUji([]recordCSV{{baris: []string{"co", "co2"}, line: 1}}); p.PlatNomor != -1 {
		t.Error("header tanpa plat_nomor seharusnya tidak dipetakan")
	}
}

func TestParseBarisHasilUji(t *testing.T) {
	p := indeksKolomUji{PlatNomor: 0, CO: 1, CO2: 2, HC: 3, O2: 4, Lambda: 5, Opasitas: 6, TanggalUji: 7}

	baris := []string{"b 1234 cde", "0,85", "12.5", "120", "1.2", "0.99", "", "2024-01-15"}
	hasil := parseBarisHasilUji(baris, p, 10)
	if hasil == nil || hasil.Err != "" {
		t.Fatalf("baris valid ditolak: %+v", hasil)
	}
	if hasil.PlatNomor != "B 1234 CDE" {
		t.Errorf("plat nomor = %q, ingin B 1234 CDE (huruf besar)", hasil.PlatNomor)
	}
	if hasil.CO == nil || *hasil.CO != 0.85 {
		t.Errorf("co = %v, ingin 0.85", hasil.CO)
	}
	if hasil.Opasitas != nil {
		t.Errorf("opasitas kosong seharusnya nil, dapat %v", *hasil.Opasitas)
	}
	if hasil.Tanggal == nil || hasil.Tanggal.Year() != 2024 {
		t.Errorf("tanggal = %v, ingin 2024-01-15", hasil.Tanggal)
	}

	if parseBarisHasilUji([]string{"", "", "", "", "", "", "", ""}, p, 11) != nil {
		t.Error("baris kosong harus menghasilkan nil")
	}

	gagal := map[int]string{
		12: "plat kosong",
		13: "nilai di luar rentang",
		14: "tanpa parameter uji",
		15: "angka tidak valid",
		16: "tanggal di masa depan",
	}
	kasus := map[int][]string{
		12: {"", "0.85", "", "", "", "", "", ""},
		13: {"B 1 AA", "250", "", "", "", "", "", ""},
		14: {"B 1 AA", "", "", "", "", "", "", ""},
		15: {"B 1 AA", "abc", "", "", "", "", "", ""},
		16: {"B 1 AA", "0.85", "", "", "", "", "", time.Now().AddDate(0, 0, 30).Format("2006-01-02")},
	}
	for nomor, baris := range kasus {
		if h := parseBarisHasilUji(baris, p, nomor); h == nil || h.Err == "" {
			t.Errorf("baris %d (%s) seharusnya ditolak", nomor, gagal[nomor])
		}
	}

	// Kolom tanggal tidak ada pada file.
	tanpaTanggal := indeksKolomUji{PlatNomor: 0, CO: 1, CO2: 2, HC: -1, O2: -1, Lambda: -1, Opasitas: -1, TanggalUji: -1}
	if h := parseBarisHasilUji([]string{"B 1 AA", "0.85", "", "", "", "", "", ""}, tanpaTanggal, 20); h == nil || h.Err != "" {
		t.Errorf("file tanpa kolom tanggal harus tetap bisa: %+v", h)
	}
}

func TestBacaRecordsXLSX(t *testing.T) {
	f := excelize.NewFile()
	defer f.Close()

	if _, err := f.NewSheet("Data Hasil Uji"); err != nil {
		t.Fatal(err)
	}
	f.SetCellValue("Data Hasil Uji", "A1", "Judul laporan yang tidak boleh terbaca")
	f.SetCellValue("Data Hasil Uji", "A3", "plat_nomor")
	f.SetCellValue("Data Hasil Uji", "B3", "co")
	f.SetCellValue("Data Hasil Uji", "C3", "opasitas")
	f.SetCellValue("Data Hasil Uji", "A4", "B 1234 CDE")
	f.SetCellValue("Data Hasil Uji", "B4", 0.85)
	f.SetCellValue("Data Hasil Uji", "C4", 12.5)

	buf, err := f.WriteToBuffer()
	if err != nil {
		t.Fatal(err)
	}

	records, err := bacaRecordsXLSX(bytes.NewReader(buf.Bytes()), barisTerlihatHeaderUji)
	if err != nil {
		t.Fatalf("bacaRecordsXLSX: %v", err)
	}
	if len(records) == 0 {
		t.Fatal("tidak ada record terbaca")
	}

	// Sheet pertama (Sheet1) tidak punya header, jadi harus memilih "Data Hasil Uji".
	baris, p := petakanKolomHasilUji(records)
	if baris != 2 {
		t.Fatalf("baris header = %d, ingin 2", baris)
	}
	if p.PlatNomor != 0 || p.CO != 1 || p.Opasitas != 2 {
		t.Fatalf("pemetaan kolom xlsx salah: %+v", p)
	}

	hasil := parseBarisHasilUji(records[3].baris, p, 4)
	if hasil == nil || hasil.Err != "" {
		t.Fatalf("baris data xlsx ditolak: %+v", hasil)
	}
	if hasil.CO == nil || *hasil.CO != 0.85 {
		t.Errorf("co dari xlsx = %v, ingin 0.85", hasil.CO)
	}
	if hasil.Opasitas == nil || *hasil.Opasitas != 12.5 {
		t.Errorf("opasitas dari xlsx = %v, ingin 12.5", hasil.Opasitas)
	}
}
