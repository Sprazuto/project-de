package models

import (
	"log"
	"sync"
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
	progressData, err := m.getProgressData(year, month, idsatker)
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
