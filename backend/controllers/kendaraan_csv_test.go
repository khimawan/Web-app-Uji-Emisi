package controllers

import (
	"strings"
	"testing"
)

func TestDeteksiPemisahCSV(t *testing.T) {
	kasus := map[string]rune{
		"kategori;jenis;plat_nomor\nUmum;Bensin M;B 1 A\n":          ';',
		"kategori,jenis,plat_nomor\nUmum,Bensin M,B 1 A\n":          ',',
		"kategori\tjenis\tplat_nomor\nUmum\tBensin M\tB 1 A\n":      '\t',
		"# panduan, dengan koma\n\nkategori,jenis\nUmum,Bensin M\n": ',',
	}
	for masuk, harap := range kasus {
		if dapat := deteksiPemisahCSV(masuk); dapat != harap {
			t.Errorf("deteksiPemisahCSV(%q) = %q, harap %q", masuk, dapat, harap)
		}
	}
}

func TestParseAngkaBulat(t *testing.T) {
	harusBerhasil := map[string]int{
		"1496":      1496,
		"1496 cc":   1496,
		"1496CC":    1496,
		" 1496 ":    1496,
		"1.496":     1496,
		"1,496":     1496,
		"1 496":     1496,
		"1496,00":   1496,
		"1496.00":   1496,
		"2020.0":    2020,
		"12,5":      12,
		"1.234.567": 1234567,
	}
	for masuk, harap := range harusBerhasil {
		dapat, err := parseAngkaBulat(masuk)
		if err != nil || dapat != harap {
			t.Errorf("parseAngkaBulat(%q) = %d, %v; harap %d", masuk, dapat, err, harap)
		}
	}

	harusGagal := []string{"", "abc", "1.2.3.4", "1496ccx", "1,49,6", "-", "th 2020", "1."}
	for _, masuk := range harusGagal {
		if _, err := parseAngkaBulat(masuk); err == nil {
			t.Errorf("parseAngkaBulat(%q) seharusnya gagal", masuk)
		}
	}
}

func TestCanonicalKategoriDanJenis(t *testing.T) {
	Kategori := map[string]string{
		"Umum": "Umum", "umum": "Umum", " DINAS ": "Dinas", "umum  ": "Umum",
	}
	for masuk, harap := range Kategori {
		if dapat, ok := canonicalKategoriKendaraan(masuk); !ok || dapat != harap {
			t.Errorf("canonicalKategoriKendaraan(%q) = %q,%v; harap %q", masuk, dapat, ok, harap)
		}
	}
	if _, ok := canonicalKategoriKendaraan("Mobil"); ok {
		t.Error(`canonicalKategoriKendaraan("Mobil") seharusnya ditolak`)
	}

	jenis := map[string]string{
		"Bensin M":             "Bensin M",
		"bensin m":             "Bensin M",
		"BENSIN_M":             "Bensin M",
		"Bensin N&O":           "Bensin N&O",
		"Bensin N & O":         "Bensin N&O",
		"bensin  n  &  o":      "Bensin N&O",
		"Solar JBB":            "Solar JBB",
		"Solar GVW":            "Solar GVW",
		"Solar JBB (<=3,5 t)":  "Solar JBB",
		"Solar GVW (>3,5 ton)": "Solar GVW",
		"solar  gvw":           "Solar GVW",
	}
	for masuk, harap := range jenis {
		if dapat, ok := canonicalJenisKendaraan(masuk); !ok || dapat != harap {
			t.Errorf("canonicalJenisKendaraan(%q) = %q,%v; harap %q", masuk, dapat, ok, harap)
		}
	}
	for _, masuk := range []string{"Benzin", "Solar", "Bensin MM", ""} {
		if _, ok := canonicalJenisKendaraan(masuk); ok {
			t.Errorf("canonicalJenisKendaraan(%q) seharusnya ditolak", masuk)
		}
	}
}

func TestPetakanKolomKendaraan(t *testing.T) {
	// Urutan kolom acak harus tetap dipetakan lewat nama header.
	p, err := petakanKolomKendaraan([]string{"MEREK", "Plat Nomor", "Kategori", "tahun_pembuatan", "tipe", "kapasitas_mesin", "jenis"})
	if err != nil {
		t.Fatalf("petakanKolomKendaraan: %v", err)
	}
	if p.Merek != 0 || p.PlatNomor != 1 || p.Kategori != 2 || p.TahunPembuatan != 3 ||
		p.Tipe != 4 || p.KapasitasMesin != 5 || p.Jenis != 6 {
		t.Fatalf("pemetaan kolom salah: %+v", p)
	}

	// Header asing dengan jumlah kolom sama -> pakai urutan template.
	if _, err := petakanKolomKendaraan([]string{"a", "b", "c", "d", "e", "f", "g"}); err != nil {
		t.Fatalf("petakanKolomKendaraan fallback: %v", err)
	}

	// Jumlah kolom tidak cocok harus ditolak dengan pesan yang jelas.
	if _, err := petakanKolomKendaraan([]string{"kategori", "jenis"}); err == nil ||
		!strings.Contains(err.Error(), "Download Template") {
		t.Fatalf("petakanKolomKendaraan harus menolak header tidak sesuai, dapat: %v", err)
	}
}

// barisCSV menyusun baris dengan urutan kolom template.
func barisCSV(nilai ...string) []string {
	baris := make([]string, len(kolomKendaraanTemplate))
	for i := range baris {
		baris[i] = "-"
	}
	copy(baris, nilai)
	return baris
}

func TestParseBarisKendaraan(t *testing.T) {
	p := indeksKolom{0, 1, 2, 3, 4, 5, 6}
	// Setiap baris contoh template harus valid.
	for _, contoh := range contohTemplateKendaraan() {
		hasil := parseBarisKendaraan(contoh, p, 2)
		if hasil == nil || hasil.Err != "" {
			t.Fatalf("baris contoh %v ditolak: %v", contoh, hasil)
		}
		if hasil.Kendaraan.Kategori != contoh[0] || hasil.Kendaraan.Jenis != contoh[1] {
			t.Errorf("kategori/jenis tidak dinormalkan: %+v", hasil.Kendaraan)
		}
		if hasil.Kendaraan.KapasitasMesin <= 0 || hasil.Kendaraan.TahunPembuatan <= 0 {
			t.Errorf("angka tidak ter-parse: %+v", hasil.Kendaraan)
		}
	}

	// Baris kosong diabaikan.
	if parseBarisKendaraan([]string{"", "  ", ""}, p, 2) != nil {
		t.Error("baris kosong seharusnya dilewati")
	}

	gagal := []struct {
		baris     []string
		harusKata string
	}{
		{barisCSV("", "Bensin M", "B 1 A", "Toyota", "Avanza", "1496", "2020"), "kategori"},
		{barisCSV("Umum", "Benzin", "B 1 A", "Toyota", "Avanza", "1496", "2020"), "tidak dikenal"},
		{barisCSV("Mobil", "Bensin M", "B 1 A", "Toyota", "Avanza", "1496", "2020"), "tidak dikenal"},
		{barisCSV("Umum", "Bensin M", "", "Toyota", "Avanza", "1496", "2020"), "plat_nomor"},
		{barisCSV("Umum", "Bensin M", "B 1 A", "", "Avanza", "1496", "2020"), "merek"},
		{barisCSV("Umum", "Bensin M", "B 1 A", "Toyota", "", "1496", "2020"), "tipe"},
		{barisCSV("Umum", "Bensin M", "B 1 A", "Toyota", "Avanza", "abc", "2020"), "kapasitas_mesin"},
		{barisCSV("Umum", "Bensin M", "B 1 A", "Toyota", "Avanza", "1496", "th 2020"), "tahun_pembuatan"},
		{barisCSV("Umum", "Bensin M", "B 1 A", "Toyota", "Avanza", "1496", "1800"), "tahun_pembuatan"},
		{barisCSV("Umum", "Bensin M", "B 1 A", "Toyota", "Avanza", "0", "2020"), "kapasitas_mesin"},
	}
	for _, tc := range gagal {
		hasil := parseBarisKendaraan(tc.baris, p, 5)
		if hasil == nil || !strings.Contains(hasil.Err, tc.harusKata) {
			t.Errorf("baris %v: err=%v, harus memuat %q", tc.baris, hasil, tc.harusKata)
		}
	}

	// Normalisasi angka & kapitalisasi plat.
	hasil := parseBarisKendaraan([]string{"umum", "bensin m", "b 1234 cde", " Toyota ", " Avanza 1.0 ", "1.496", "2020.0"}, p, 2)
	if hasil == nil || hasil.Err != "" {
		t.Fatalf("baris ternormalisasi ditolak: %v", hasil)
	}
	if hasil.Kendaraan.PlatNomor != "B 1234 CDE" || hasil.Kendaraan.Merek != "Toyota" ||
		hasil.Kendaraan.Tipe != "Avanza 1.0" || hasil.Kendaraan.KapasitasMesin != 1496 ||
		hasil.Kendaraan.TahunPembuatan != 2020 {
		t.Errorf("normalisasi tidak sesuai: %+v", hasil.Kendaraan)
	}
}
