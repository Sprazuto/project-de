package models

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/Massad/gin-boilerplate/db"
)

// RealisasiData represents raw data from database for frontend processing
type RealisasiData struct {
	Category      string  `json:"category"`      // "barjas", "fisik", "anggaran", "kinerja"
	Progress      float64 `json:"progress"`
	ProgressFormatted string `json:"progress_formatted,omitempty"`
	Capaian       float64 `json:"capaian,omitempty"`
	Items         []RealisasiRawItem `json:"items"`
}

// RealisasiRawItem represents raw item data without formatting
type RealisasiRawItem struct {
	Type      string      `json:"type"`      // "perencanaan", "pemilihan", "pengadaan", "penyerahan", "realisasi", "target"
	Value     interface{} `json:"value"`     // Raw value from database
	Formatted string      `json:"formatted,omitempty"` // Formatted value for display
	Detail    *RealisasiDetail `json:"detail,omitempty"` // Additional details for complex items
}

// RealisasiDetail represents additional details for complex items
type RealisasiDetail struct {
	Selesai   int64 `json:"selesai,omitempty"`
	Target    int64 `json:"target,omitempty"`
	Terlambat int64 `json:"terlambat,omitempty"`
}

// Backward compatibility aliases
type RealisasiBulanCard = RealisasiData
type RealisasiBulanItem = RealisasiRawItem

// RealisasiResponse represents the standard response structure for realisasi data
type RealisasiResponse struct {
	Data []RealisasiData `json:"data"`
	Meta RealisasiMeta   `json:"meta"`
}

// RealisasiMeta represents metadata for the realisasi response
type RealisasiMeta struct {
	Year      int    `json:"year"`
	Month     int    `json:"month"`
	MonthName string `json:"month_name,omitempty"`
	Idsatker  int    `json:"idsatker"`
	Type      string `json:"type"` // "bulan" or "tahun"
}

// RealisasiBulanResponse alias for backward compatibility
type RealisasiBulanResponse = RealisasiResponse

// RealisasiTahunResponse alias for backward compatibility
type RealisasiTahunResponse = RealisasiResponse

// ProgressData represents the progress data from de_ranking_opd
type ProgressData struct {
	CapaianOpd        float64
	CapaianBarjas     float64
	CapaianFisik      float64
	CapaianAnggaran   float64
	CapaianKinerja    float64
	KumulatifOpd      float64
	KumulatifBarjas   float64
	KumulatifFisik    float64
	KumulatifAnggaran float64
	KumulatifKinerja  float64
}

// SijagurModel ...
type SijagurModel struct{}

// getProgressData retrieves progress data from de_ranking_opd for monthly data
func (m SijagurModel) getProgressData(year, month, idsatker int) (ProgressData, error) {
	var progress ProgressData

	query := `
		SELECT
			capaian_opd, capaian_barjas, capaian_fisik, capaian_anggaran, capaian_kinerja,
			kumulatif_opd, kumulatif_barjas, kumulatif_fisik, kumulatif_anggaran, kumulatif_kinerja
		FROM de_ranking_opd
		WHERE tahun = $1 AND bulan = $2 AND idsatker = $3
		LIMIT 1
	`

	err := db.GetDB().QueryRow(query, year, month, idsatker).Scan(
		&progress.CapaianOpd, &progress.CapaianBarjas, &progress.CapaianFisik,
		&progress.CapaianAnggaran, &progress.CapaianKinerja,
		&progress.KumulatifOpd, &progress.KumulatifBarjas, &progress.KumulatifFisik,
		&progress.KumulatifAnggaran, &progress.KumulatifKinerja,
	)

	if err != nil {
		log.Printf("Error querying progress data: %v", err)
		return ProgressData{}, err
	}

	return progress, nil
}

// getProgressDataTahun retrieves progress data from de_ranking_opd for yearly data
func (m SijagurModel) getProgressDataTahun(year, month, idsatker int) (ProgressData, error) {
	var progress ProgressData

	query := `
		SELECT
			capaian_opd, capaian_barjas, capaian_fisik, capaian_anggaran, capaian_kinerja,
			kumulatif_opd, kumulatif_barjas, kumulatif_fisik, kumulatif_anggaran, kumulatif_kinerja
		FROM de_ranking_opd
		WHERE tahun = $1 AND bulan = $2 AND idsatker = $3
		LIMIT 1
	`

	err := db.GetDB().QueryRow(query, year, month, idsatker).Scan(
		&progress.CapaianOpd, &progress.CapaianBarjas, &progress.CapaianFisik,
		&progress.CapaianAnggaran, &progress.CapaianKinerja,
		&progress.KumulatifOpd, &progress.KumulatifBarjas, &progress.KumulatifFisik,
		&progress.KumulatifAnggaran, &progress.KumulatifKinerja,
	)

	if err != nil {
		log.Printf("Error querying progress data tahun: %v", err)
		return ProgressData{}, err
	}

	return progress, nil
}


// GetRealisasiBulanWithParams retrieves raw realisasi bulan data for frontend processing
func (m SijagurModel) GetRealisasiBulanWithParams(year, month, idsatker int) ([]RealisasiData, error) {
	// First, get the progress percentages from de_ranking_opd
	progressData, err := m.getProgressData(year, month, idsatker)
	if err != nil {
		log.Printf("Error getting progress data: %v", err)
		return nil, err
	}

	// Use goroutines to fetch data concurrently
	var wg sync.WaitGroup
	var mu sync.Mutex
	data := make([]RealisasiData, 4)
	errors := make([]error, 4)

	// Helper function to fetch data concurrently
	fetchData := func(index int, fetchFunc func() (RealisasiData, error)) {
		defer wg.Done()
		result, err := fetchFunc()
		mu.Lock()
		data[index] = result
		errors[index] = err
		mu.Unlock()
	}

	wg.Add(4)
	go fetchData(0, func() (RealisasiData, error) { return m.getBarjasRawData(year, month, idsatker, progressData.CapaianBarjas) })
	go fetchData(1, func() (RealisasiData, error) { return m.getFisikRawData(year, month, idsatker, progressData.CapaianFisik) })
	go fetchData(2, func() (RealisasiData, error) { return m.getAnggaranRawData(year, month, idsatker, progressData.CapaianAnggaran) })
	go fetchData(3, func() (RealisasiData, error) { return m.getKinerjaRawData(year, month, idsatker, progressData.CapaianKinerja) })

	wg.Wait()

	// Check for errors
	for _, err := range errors {
		if err != nil {
			return nil, err
		}
	}

	return data, nil
}


// GetRealisasiTahunWithParams retrieves raw realisasi tahun data for frontend processing
func (m SijagurModel) GetRealisasiTahunWithParams(year, month, idsatker int) ([]RealisasiData, error) {
	// First, get the progress percentages from de_ranking_opd
	progressData, err := m.getProgressDataTahun(year, month, idsatker)
	if err != nil {
		log.Printf("Error getting progress data tahun: %v", err)
		return nil, err
	}

	// Use goroutines to fetch data concurrently
	var wg sync.WaitGroup
	var mu sync.Mutex
	data := make([]RealisasiData, 4)
	errors := make([]error, 4)

	// Helper function to fetch data concurrently
	fetchData := func(index int, fetchFunc func() (RealisasiData, error)) {
		defer wg.Done()
		result, err := fetchFunc()
		mu.Lock()
		data[index] = result
		errors[index] = err
		mu.Unlock()
	}

	wg.Add(4)
	go fetchData(0, func() (RealisasiData, error) { return m.getBarjasRawDataTahun(year, month, idsatker, progressData.KumulatifBarjas, progressData.CapaianBarjas) })
	go fetchData(1, func() (RealisasiData, error) { return m.getFisikRawDataTahun(year, month, idsatker, progressData.KumulatifFisik, progressData.CapaianFisik) })
	go fetchData(2, func() (RealisasiData, error) { return m.getAnggaranRawDataTahun(year, month, idsatker, progressData.KumulatifAnggaran, progressData.CapaianAnggaran) })
	go fetchData(3, func() (RealisasiData, error) { return m.getKinerjaRawDataTahun(year, month, idsatker, progressData.KumulatifKinerja, progressData.CapaianKinerja) })

	wg.Wait()

	// Check for errors
	for _, err := range errors {
		if err != nil {
			return nil, err
		}
	}

	return data, nil
}

// getBarjasRawData retrieves raw barjas data for frontend processing
func (m SijagurModel) getBarjasRawData(year, month, idsatker int, progress float64) (RealisasiData, error) {
	var data RealisasiData

	// Query to get barjas data
	query := `
		SELECT
			c_perencanaan_selesai, c_perencanaan_target, c_perencanaan_terlambat,
			c_pemilihan_selesai, c_pemilihan_target, c_pemilihan_terlambat,
			c_pengadaan_selesai, c_pengadaan_target, c_pengadaan_terlambat,
			c_penyerahan_selesai, c_penyerahan_target, c_penyerahan_terlambat
		FROM de_detail_barjas ddb
		INNER JOIN de_ranking_opd dro ON ddb.id_ranking_opd = dro.id
		WHERE dro.tahun = $1 AND dro.bulan = $2 AND dro.idsatker = $3
	`

	var perencanaanSelesai, perencanaanTarget, perencanaanTerlambat int64
	var pemilihanSelesai, pemilihanTarget, pemilihanTerlambat int64
	var pengadaanSelesai, pengadaanTarget, pengadaanTerlambat int64
	var penyerahanSelesai, penyerahanTarget, penyerahanTerlambat int64

	err := db.GetDB().QueryRow(query, year, month, idsatker).Scan(
		&perencanaanSelesai, &perencanaanTarget, &perencanaanTerlambat,
		&pemilihanSelesai, &pemilihanTarget, &pemilihanTerlambat,
		&pengadaanSelesai, &pengadaanTarget, &pengadaanTerlambat,
		&penyerahanSelesai, &penyerahanTarget, &penyerahanTerlambat,
	)
	if err != nil {
		log.Printf("Error querying barjas data: %v", err)
		return RealisasiData{}, err
	}

	// Build raw items for frontend processing
	items := []RealisasiRawItem{
		{Type: "perencanaan", Value: perencanaanSelesai, Detail: &RealisasiDetail{Selesai: perencanaanSelesai, Target: perencanaanTarget, Terlambat: perencanaanTerlambat}},
		{Type: "pemilihan", Value: pemilihanSelesai, Detail: &RealisasiDetail{Selesai: pemilihanSelesai, Target: pemilihanTarget, Terlambat: pemilihanTerlambat}},
		{Type: "pengadaan", Value: pengadaanSelesai, Detail: &RealisasiDetail{Selesai: pengadaanSelesai, Target: pengadaanTarget, Terlambat: pengadaanTerlambat}},
		{Type: "penyerahan", Value: penyerahanSelesai, Detail: &RealisasiDetail{Selesai: penyerahanSelesai, Target: penyerahanTarget, Terlambat: penyerahanTerlambat}},
	}

	data = RealisasiData{
		Category:          "barjas",
		Progress:           progress,
		ProgressFormatted:  m.formatProgress(progress),
		Items:              items,
	}

	return data, nil
}

// getFisikRawData retrieves raw fisik data for frontend processing
func (m SijagurModel) getFisikRawData(year, month, idsatker int, progress float64) (RealisasiData, error) {
	var data RealisasiData

	// Query to get fisik data
	query := `
		SELECT
			c_fisik_realisasi, c_fisik_target
		FROM de_detail_fisik ddf
		INNER JOIN de_ranking_opd dro ON ddf.id_ranking_opd = dro.id
		WHERE dro.tahun = $1 AND dro.bulan = $2 AND dro.idsatker = $3
	`

	var realisasi, target float64
	err := db.GetDB().QueryRow(query, year, month, idsatker).Scan(&realisasi, &target)
	if err != nil {
		log.Printf("Error querying fisik data: %v", err)
		return RealisasiData{}, err
	}

	// Build raw items for frontend processing
	items := []RealisasiRawItem{
		{Type: "realisasi", Value: realisasi, Formatted: m.formatNumber(realisasi)},
		{Type: "target", Value: target, Formatted: m.formatNumber(target)},
	}

	data = RealisasiData{
		Category:          "fisik",
		Progress:           progress,
		ProgressFormatted:  m.formatProgress(progress),
		Items:              items,
	}

	return data, nil
}

// getAnggaranRawData retrieves raw anggaran data for frontend processing
func (m SijagurModel) getAnggaranRawData(year, month, idsatker int, progress float64) (RealisasiData, error) {
	var data RealisasiData

	// Query to get anggaran data
	query := `
		SELECT
			c_anggaran_realisasi, c_anggaran_target
		FROM de_detail_anggaran dda
		INNER JOIN de_ranking_opd dro ON dda.id_ranking_opd = dro.id
		WHERE dro.tahun = $1 AND dro.bulan = $2 AND dro.idsatker = $3
	`

	var realisasi, target float64
	err := db.GetDB().QueryRow(query, year, month, idsatker).Scan(&realisasi, &target)
	if err != nil {
		log.Printf("Error querying anggaran data: %v", err)
		return RealisasiData{}, err
	}

	// Build raw items for frontend processing
	items := []RealisasiRawItem{
		{Type: "realisasi", Value: realisasi, Formatted: m.formatCurrency(realisasi)},
		{Type: "target", Value: target, Formatted: m.formatCurrency(target)},
	}

	data = RealisasiData{
		Category:          "anggaran",
		Progress:           progress,
		ProgressFormatted:  m.formatProgress(progress),
		Items:              items,
	}

	return data, nil
}

// getKinerjaRawData retrieves raw kinerja data for frontend processing
func (m SijagurModel) getKinerjaRawData(year, month, idsatker int, progress float64) (RealisasiData, error) {
	var data RealisasiData

	// Query to get kinerja data
	query := `
		SELECT
			c_kinerja_realisasi, c_kinerja_target
		FROM de_detail_kinerja ddk
		INNER JOIN de_ranking_opd dro ON ddk.id_ranking_opd = dro.id
		WHERE dro.tahun = $1 AND dro.bulan = $2 AND dro.idsatker = $3
	`

	var realisasi, target float64
	err := db.GetDB().QueryRow(query, year, month, idsatker).Scan(&realisasi, &target)
	if err != nil {
		log.Printf("Error querying kinerja data: %v", err)
		return RealisasiData{}, err
	}

	// Build raw items for frontend processing
	items := []RealisasiRawItem{
		{Type: "realisasi", Value: realisasi, Formatted: m.formatNumber(realisasi)},
		{Type: "target", Value: target, Formatted: m.formatNumber(target)},
	}

	data = RealisasiData{
		Category:          "kinerja",
		Progress:           progress,
		ProgressFormatted:  m.formatProgress(progress),
		Items:              items,
	}

	return data, nil
}

// getBarjasRawDataTahun retrieves raw barjas tahun data for frontend processing
func (m SijagurModel) getBarjasRawDataTahun(year, month, idsatker int, progress float64, capaian float64) (RealisasiData, error) {
	var data RealisasiData

	// Query to get barjas data for the year up to the specified month
	query := `
		SELECT
			k_barjas_realisasi, k_barjas_target
		FROM de_detail_barjas ddb
		INNER JOIN de_ranking_opd dro ON ddb.id_ranking_opd = dro.id
		WHERE dro.tahun = $1 AND dro.bulan = $2 AND dro.idsatker = $3
	`

	var totalRealisasi, totalTarget float64
	err := db.GetDB().QueryRow(query, year, month, idsatker).Scan(&totalRealisasi, &totalTarget)
	if err != nil {
		log.Printf("Error querying barjas tahun data: %v", err)
		return RealisasiData{}, err
	}

	// Build raw items for frontend processing
	items := []RealisasiRawItem{
		{Type: "realisasi", Value: totalRealisasi, Formatted: m.formatNumber(totalRealisasi)},
		{Type: "target", Value: totalTarget, Formatted: m.formatNumber(totalTarget)},
	}

	data = RealisasiData{
		Category:          "barjas",
		Progress:           progress,
		ProgressFormatted:  m.formatProgress(progress),
		Capaian:            capaian,
		Items:              items,
	}

	return data, nil
}

func (m SijagurModel) getFisikRawDataTahun(year, month, idsatker int, progress float64, capaian float64) (RealisasiData, error) {
	var data RealisasiData

	query := `
		SELECT
			k_fisik_realisasi, k_fisik_target
		FROM de_detail_fisik ddf
		INNER JOIN de_ranking_opd dro ON ddf.id_ranking_opd = dro.id
		WHERE dro.tahun = $1 AND dro.bulan = $2 AND dro.idsatker = $3
	`

	var realisasi, target float64
	err := db.GetDB().QueryRow(query, year, month, idsatker).Scan(&realisasi, &target)
	if err != nil {
		log.Printf("Error querying fisik tahun data: %v", err)
		return RealisasiData{}, err
	}

	// Build raw items for frontend processing
	items := []RealisasiRawItem{
		{Type: "realisasi", Value: realisasi, Formatted: m.formatNumber(realisasi)},
		{Type: "target", Value: target, Formatted: m.formatNumber(target)},
	}

	data = RealisasiData{
		Category:          "fisik",
		Progress:           progress,
		ProgressFormatted:  m.formatProgress(progress),
		Capaian:            capaian,
		Items:              items,
	}

	return data, nil
}

func (m SijagurModel) getAnggaranRawDataTahun(year, month, idsatker int, progress float64, capaian float64) (RealisasiData, error) {
	var data RealisasiData

	query := `
		SELECT
			k_anggaran_realisasi, k_anggaran_target
		FROM de_detail_anggaran dda
		INNER JOIN de_ranking_opd dro ON dda.id_ranking_opd = dro.id
		WHERE dro.tahun = $1 AND dro.bulan = $2 AND dro.idsatker = $3
	`

	var realisasi, target float64
	err := db.GetDB().QueryRow(query, year, month, idsatker).Scan(&realisasi, &target)
	if err != nil {
		log.Printf("Error querying anggaran tahun data: %v", err)
		return RealisasiData{}, err
	}

	// Build raw items for frontend processing
	items := []RealisasiRawItem{
		{Type: "realisasi", Value: realisasi, Formatted: m.formatCurrency(realisasi)},
		{Type: "target", Value: target, Formatted: m.formatCurrency(target)},
	}

	data = RealisasiData{
		Category:          "anggaran",
		Progress:           progress,
		ProgressFormatted:  m.formatProgress(progress),
		Capaian:            capaian,
		Items:              items,
	}

	return data, nil
}

func (m SijagurModel) getKinerjaRawDataTahun(year, month, idsatker int, progress float64, capaian float64) (RealisasiData, error) {
	var data RealisasiData

	query := `
		SELECT
			k_kinerja_realisasi, k_kinerja_target
		FROM de_detail_kinerja ddk
		INNER JOIN de_ranking_opd dro ON ddk.id_ranking_opd = dro.id
		WHERE dro.tahun = $1 AND dro.bulan = $2 AND dro.idsatker = $3
	`

	var realisasi, target float64
	err := db.GetDB().QueryRow(query, year, month, idsatker).Scan(&realisasi, &target)
	if err != nil {
		log.Printf("Error querying kinerja tahun data: %v", err)
		return RealisasiData{}, err
	}

	// Build raw items for frontend processing
	items := []RealisasiRawItem{
		{Type: "realisasi", Value: realisasi, Formatted: m.formatNumber(realisasi)},
		{Type: "target", Value: target, Formatted: m.formatNumber(target)},
	}

	data = RealisasiData{
		Category:          "kinerja",
		Progress:           progress,
		ProgressFormatted:  m.formatProgress(progress),
		Capaian:            capaian,
		Items:              items,
	}

	return data, nil
}

func (m SijagurModel) formatCurrency(amount float64) string {
	s := fmt.Sprintf("%.0f", amount)
	if s == "0" {
		return "Rp0,00"
	}
	// Add dots every 3 digits from the right
	var result []rune
	for i, r := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			result = append(result, '.')
		}
		result = append(result, r)
	}
	return "Rp" + string(result) + ",00"
}

func (m SijagurModel) formatNumber(value float64) string {
	// Remove trailing zeros and decimal point if necessary
	s := fmt.Sprintf("%.2f", value)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	return s
}

func (m SijagurModel) formatProgress(value float64) string {
	// For progress, round down to 2 decimal places without rounding up
	// Multiply by 100, floor, then divide by 100
	truncated := float64(int(value*100)) / 100
	s := fmt.Sprintf("%.2f", truncated)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	return s
}

var monthNames = []string{
	"", "Januari", "Februari", "Maret", "April", "Mei", "Juni",
	"Juli", "Agustus", "September", "Oktober", "November", "Desember",
}

func (m SijagurModel) GetMonthName(month int) string {
	if month >= 1 && month <= 12 {
		return monthNames[month]
	}
	return "Unknown"
}
