package controllers

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/xuri/excelize/v2"
)

func TestCariBarisHeaderKendaraan(t *testing.T) {
	header := func() []string {
		return []string{"kategori", "jenis", "plat_nomor", "merek", "tipe", "kapasitas_mesin", "tahun_pembuatan"}
	}

	t.Run("header di baris pertama", func(t *testing.T) {
		records := []recordCSV{
			{baris: header(), line: 1},
			{baris: []string{"Umum", "Bensin M", "B 1 AA", "Toyota", "Avanza", "1329", "2020"}, line: 2},
		}
		idx, p, err := cariBarisHeaderKendaraan(records)
		if err != nil {
			t.Fatal(err)
		}
		if idx != 0 {
			t.Fatalf("indeks header = %d, ingin 0", idx)
		}
		if p.PlatNomor != 2 || p.TahunPembuatan != 6 {
			t.Fatalf("pemetaan kolom salah: %+v", p)
		}
	})

	t.Run("header setelah baris judul ala Excel", func(t *testing.T) {
		records := []recordCSV{
			{baris: []string{"Laporan kendaraan 2026"}, line: 1},
			{baris: nil, line: 2},
			{baris: []string{"Dataclidesomething"}, line: 3},
			{baris: []string{"kategori", "jenis", "plat_nomor", "merek", "tipe", "kapasitas_mesin", "tahun_pembuatan"}, line: 4},
			{baris: []string{"Umum", "Solar JBB", "B 1 AB", "Hino", "FM 440", "12916", "2018"}, line: 5},
		}
		idx, p, err := cariBarisHeaderKendaraan(records)
		if err != nil {
			t.Fatal(err)
		}
		if idx != 3 {
			t.Fatalf("indeks header = %d, ingin 3", idx)
		}
		if p.Kategori != 0 || p.Jenis != 1 {
			t.Fatalf("pemetaan kolom salah: %+v", p)
		}
	})

	t.Run("tanpa header memakai urutan template", func(t *testing.T) {
		records := []recordCSV{
			{baris: []string{"Umum", "Bensin N&O", "B 1 AC", "Isuzu", "Traga", "2477", "2019"}, line: 1},
		}
		idx, p, err := cariBarisHeaderKendaraan(records)
		if err != nil {
			t.Fatal(err)
		}
		if idx != -1 {
			t.Fatalf("indeks header = %d, ingin -1 (tanpa header)", idx)
		}
		if p.Kategori != 0 || p.Jenis != 1 || p.PlatNomor != 2 || p.TahunPembuatan != 6 {
			t.Fatalf("pemetaan template salah: %+v", p)
		}
	})

	t.Run("header di luar batas pencarian", func(t *testing.T) {
		records := make([]recordCSV, 0, 25)
		for i := 0; i < 25; i++ {
			records = append(records, recordCSV{baris: []string{"catatan ke-" + string(rune('a'+i))}, line: i + 1})
		}
		records = append(records, recordCSV{baris: header(), line: 26})
		idx, _, err := cariBarisHeaderKendaraan(records)
		if err != nil {
			t.Fatal(err)
		}
		if idx != -1 {
			t.Fatalf("indeks header = %d, ingin -1 (di luar 20 baris pertama)", idx)
		}
	})
}

func TestBacaRecordsXLSXKendaraan(t *testing.T) {
	// Meniru struktur template Excel kendaraan: sheet panduan, sheet data dengan
	// header di baris 6, dan baris contoh.
	f := excelize.NewFile()
	defer f.Close()

	if _, err := f.NewSheet(shPanduanKendaraan); err != nil {
		t.Fatal(err)
	}
	if _, err := f.NewSheet(shDataKendaraan); err != nil {
		t.Fatal(err)
	}
	f.SetCellValue(shPanduanKendaraan, "A1", "Panduan Input Data Kendaraan")
	f.SetCellValue(shDataKendaraan, "A1", "Isi data mulai baris 7.")
	for i, nama := range kolomKendaraanTemplate {
		f.SetCellValue(shDataKendaraan, selKendaraan(kolomKendaraanXLSX[i], barisHeaderKendaraan), nama)
	}
	for i, row := range contohTemplateKendaraan() {
		for j, nilai := range row {
			f.SetCellValue(shDataKendaraan, selKendaraan(kolomKendaraanXLSX[j], barisAwalDataKendaraan+i), nilai)
		}
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		t.Fatal(err)
	}

	records, err := bacaRecordsXLSX(bytes.NewReader(buf.Bytes()), barisTerlihatHeader)
	if err != nil {
		t.Fatalf("bacaRecordsXLSX: %v", err)
	}

	idx, p, err := cariBarisHeaderKendaraan(records)
	if err != nil {
		t.Fatal(err)
	}
	if idx != barisHeaderKendaraan-1 {
		t.Fatalf("indeks header = %d, ingin %d", idx, barisHeaderKendaraan-1)
	}
	if p.Kategori != 0 || p.PlatNomor != 2 || p.KapasitasMesin != 5 {
		t.Fatalf("pemetaan kolom salah: %+v", p)
	}

	// Baris pertama setelah header harus terbaca sebagai data kendaraan.
	barisData := records[idx+1:]
	if len(barisData) != len(contohTemplateKendaraan()) {
		t.Fatalf("jumlah baris data = %d, ingin %d", len(barisData), len(contohTemplateKendaraan()))
	}
	hasil := parseBarisKendaraan(barisData[0].baris, p, barisData[0].line)
	if hasil == nil || hasil.Err != "" {
		t.Fatalf("baris contoh template ditolak: %+v", hasil)
	}
	if hasil.Kendaraan.PlatNomor != "B 1234 CDE" || hasil.Kendaraan.KapasitasMesin != 1329 {
		t.Fatalf("nilai baris template salah: %+v", hasil.Kendaraan)
	}
}

func TestKanonKolomKendaraan(t *testing.T) {
	harus := map[string]string{
		"plat_nomor":      "plat_nomor",
		"Plat Nomor":      "plat_nomor",
		"Plat  Nomor":     "plat_nomor",
		"No. Polisi":      "plat_nomor",
		"Kapasitas (cc)":  "kapasitas_mesin",
		"KAPASITAS MESIN": "kapasitas_mesin",
		"Tahun Pembuatan": "tahun_pembuatan",
		"Tahun":           "tahun_pembuatan",
		"Merek":           "merek",
		"Type":            "tipe",
		"tipe_kendaraan":  "tipe",
		" kategori ":      "kategori",
		"kolom_asing":     "kolom_asing",
	}
	for raw, want := range harus {
		if got := kanonKolomKendaraan(raw); got != want {
			t.Errorf("kanonKolomKendaraan(%q) = %q, ingin %q", raw, got, want)
		}
	}
}

func TestKeteranganJenisKendaraanLengkap(t *testing.T) {
	for _, jenis := range JenisKendaraanValid {
		if keteranganJenisKendaraan[jenis] == "" {
			t.Errorf("keterangan untuk jenis %q kosong", jenis)
		}
	}
}

func selKendaraan(kolom string, baris int) string {
	return fmt.Sprintf("%s%d", kolom, baris)
}
