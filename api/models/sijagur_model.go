package models

import (
	"log"
	"sync"
)

// RealisasiMonthlyItem represents monthly data for perbulan endpoint
type RealisasiMonthlyItem struct {
	Month              string  `json:"month"`
	Value              float64 `json:"value"`
	ValueFormatted     string  `json:"value_formatted"`
	Realisasi          float64 `json:"realisasi,omitempty"`
	Target             float64 `json:"target,omitempty"`
	RealisasiFormatted string  `json:"realisasi_formatted,omitempty"`
	TargetFormatted    string  `json:"target_formatted,omitempty"`
}

// RealisasiData represents raw data from database for frontend processing
type RealisasiData struct {
	Category          string                 `json:"category"` // "barjas", "fisik", "anggaran", "kinerja"
	Progress          float64                `json:"progress"`
	ProgressFormatted string                 `json:"progress_formatted,omitempty"`
	Capaian           float64                `json:"capaian,omitempty"`
	Items             []RealisasiRawItem     `json:"items,omitempty"`
	Monthly           []RealisasiMonthlyItem `json:"monthly,omitempty"`       // For perbulan data
	CurrentMonth      *RealisasiMonthlyItem  `json:"current_month,omitempty"` // For perbulan data
}

// RealisasiRawItem represents raw item data without formatting
type RealisasiRawItem struct {
	Type      string           `json:"type"`                // "perencanaan", "pemilihan", "pengadaan", "penyerahan", "realisasi", "target"
	Value     interface{}      `json:"value"`               // Raw value from database
	Formatted string           `json:"formatted,omitempty"` // Formatted value for display
	Detail    *RealisasiDetail `json:"detail,omitempty"`    // Additional details for complex items
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

// SijagurData ...
type SijagurData struct{}

// GetRealisasiBulanWithParams retrieves raw realisasi bulan data for frontend processing
func (m SijagurData) GetRealisasiBulanWithParams(year, month, idsatker int) ([]RealisasiData, error) {
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
	go fetchData(0, func() (RealisasiData, error) {
		return m.getBarjasRawData(year, month, idsatker, progressData.CapaianBarjas)
	})
	go fetchData(1, func() (RealisasiData, error) {
		return m.getFisikRawData(year, month, idsatker, progressData.CapaianFisik)
	})
	go fetchData(2, func() (RealisasiData, error) {
		return m.getAnggaranRawData(year, month, idsatker, progressData.CapaianAnggaran)
	})
	go fetchData(3, func() (RealisasiData, error) {
		return m.getKinerjaRawData(year, month, idsatker, progressData.CapaianKinerja)
	})

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
func (m SijagurData) GetRealisasiTahunWithParams(year, month, idsatker int) ([]RealisasiData, error) {
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
	go fetchData(0, func() (RealisasiData, error) {
		return m.getBarjasRawDataTahun(year, month, idsatker, progressData.KumulatifBarjas, progressData.CapaianBarjas)
	})
	go fetchData(1, func() (RealisasiData, error) {
		return m.getFisikRawDataTahun(year, month, idsatker, progressData.KumulatifFisik, progressData.CapaianFisik)
	})
	go fetchData(2, func() (RealisasiData, error) {
		return m.getAnggaranRawDataTahun(year, month, idsatker, progressData.KumulatifAnggaran, progressData.CapaianAnggaran)
	})
	go fetchData(3, func() (RealisasiData, error) {
		return m.getKinerjaRawDataTahun(year, month, idsatker, progressData.KumulatifKinerja, progressData.CapaianKinerja)
	})

	wg.Wait()

	// Check for errors
	for _, err := range errors {
		if err != nil {
			return nil, err
		}
	}

	return data, nil
}

// ProcessRealisasiPerbulan processes pre-fetched monthly progress data into RealisasiData format
func (m SijagurData) ProcessRealisasiPerbulan(progressByMonth []struct {
	Month             int
	PeriodikBarjas    float64
	PeriodikFisik     float64
	PeriodikAnggaran  float64
	PeriodikKinerja   float64
	RealisasiBarjas   float64
	TargetBarjas      float64
	RealisasiFisik    float64
	TargetFisik       float64
	RealisasiAnggaran float64
	TargetAnggaran    float64
	RealisasiKinerja  float64
	TargetKinerja     float64
}) []RealisasiData {
	// Month labels in Indonesian, 1-based index.
	monthLabels := []string{"Januari", "Februari", "Maret", "April", "Mei", "Juni", "Juli", "Agustus", "September", "Oktober", "November", "Desember"}

	// Helper to build a category series from extractor, matching legacy-style formatting via Formatter.
	buildCategory := func(category string, extract func(p struct {
		Month             int
		PeriodikBarjas    float64
		PeriodikFisik     float64
		PeriodikAnggaran  float64
		PeriodikKinerja   float64
		RealisasiBarjas   float64
		TargetBarjas      float64
		RealisasiFisik    float64
		TargetFisik       float64
		RealisasiAnggaran float64
		TargetAnggaran    float64
		RealisasiKinerja  float64
		TargetKinerja     float64
	}) float64, extractRealisasi func(p struct {
		Month             int
		PeriodikBarjas    float64
		PeriodikFisik     float64
		PeriodikAnggaran  float64
		PeriodikKinerja   float64
		RealisasiBarjas   float64
		TargetBarjas      float64
		RealisasiFisik    float64
		TargetFisik       float64
		RealisasiAnggaran float64
		TargetAnggaran    float64
		RealisasiKinerja  float64
		TargetKinerja     float64
	}) float64, extractTarget func(p struct {
		Month             int
		PeriodikBarjas    float64
		PeriodikFisik     float64
		PeriodikAnggaran  float64
		PeriodikKinerja   float64
		RealisasiBarjas   float64
		TargetBarjas      float64
		RealisasiFisik    float64
		TargetFisik       float64
		RealisasiAnggaran float64
		TargetAnggaran    float64
		RealisasiKinerja  float64
		TargetKinerja     float64
	}) float64) RealisasiData {
		var monthly []RealisasiMonthlyItem
		var latestMonthIdx = -1

		formatter := Formatter{}

		// Initialize all 12 months with 0.
		for i := 0; i < 12; i++ {
			monthly = append(monthly, RealisasiMonthlyItem{
				Month:          monthLabels[i],
				Value:          0,
				ValueFormatted: formatter.FormatProgress(0),
				Realisasi:      0,
				Target:         100,
			})
		}

		// Fill from DB results.
		for _, p := range progressByMonth {
			value := extract(p)
			idx := p.Month - 1
			if idx >= 0 && idx < len(monthly) {
				monthly[idx].Value = value
				monthly[idx].ValueFormatted = formatter.FormatProgress(value)
				// Get realisasi and target from the appropriate fields based on category
				switch category {
				case "barjas":
					monthly[idx].Realisasi = p.RealisasiBarjas
					monthly[idx].Target = p.TargetBarjas
					monthly[idx].RealisasiFormatted = formatter.FormatNumber(p.RealisasiBarjas)
					monthly[idx].TargetFormatted = formatter.FormatNumber(p.TargetBarjas)
				case "fisik":
					monthly[idx].Realisasi = p.RealisasiFisik
					monthly[idx].Target = p.TargetFisik
					monthly[idx].RealisasiFormatted = formatter.FormatNumber(p.RealisasiFisik)
					monthly[idx].TargetFormatted = formatter.FormatNumber(p.TargetFisik)
				case "anggaran":
					monthly[idx].Realisasi = p.RealisasiAnggaran
					monthly[idx].Target = p.TargetAnggaran
					monthly[idx].RealisasiFormatted = formatter.FormatCurrency(p.RealisasiAnggaran)
					monthly[idx].TargetFormatted = formatter.FormatCurrency(p.TargetAnggaran)
				case "kinerja":
					monthly[idx].Realisasi = p.RealisasiKinerja
					monthly[idx].Target = p.TargetKinerja
					monthly[idx].RealisasiFormatted = formatter.FormatNumber(p.RealisasiKinerja)
					monthly[idx].TargetFormatted = formatter.FormatNumber(p.TargetKinerja)
				}
				if value > 0 {
					latestMonthIdx = idx
				}
			}
		}

		// Determine current_month:
		// - Latest month with value > 0
		// - If none, fallback to Desember with 0
		var currentMonth *RealisasiMonthlyItem
		if latestMonthIdx >= 0 {
			currentMonth = &RealisasiMonthlyItem{
				Month:              monthLabels[latestMonthIdx],
				Value:              monthly[latestMonthIdx].Value,
				ValueFormatted:     monthly[latestMonthIdx].ValueFormatted,
				Realisasi:          monthly[latestMonthIdx].Realisasi,
				Target:             monthly[latestMonthIdx].Target,
				RealisasiFormatted: monthly[latestMonthIdx].RealisasiFormatted,
				TargetFormatted:    monthly[latestMonthIdx].TargetFormatted,
			}
		} else {
			currentMonth = &RealisasiMonthlyItem{
				Month:          monthLabels[11], // "Desember"
				Value:          0,
				ValueFormatted: formatter.FormatProgress(0),
				Realisasi:      0,
				Target:         100,
			}
		}

		return RealisasiData{
			Category:     category,
			Progress:     0, // Not used for perbulan
			Monthly:      monthly,
			CurrentMonth: currentMonth,
		}
	}

	// Build categories using periodik_* columns as monthly values, with legacy-style formatted strings.
	data := []RealisasiData{
		buildCategory("barjas", func(p struct {
			Month             int
			PeriodikBarjas    float64
			PeriodikFisik     float64
			PeriodikAnggaran  float64
			PeriodikKinerja   float64
			RealisasiBarjas   float64
			TargetBarjas      float64
			RealisasiFisik    float64
			TargetFisik       float64
			RealisasiAnggaran float64
			TargetAnggaran    float64
			RealisasiKinerja  float64
			TargetKinerja     float64
		}) float64 {
			return p.PeriodikBarjas
		}, func(p struct {
			Month             int
			PeriodikBarjas    float64
			PeriodikFisik     float64
			PeriodikAnggaran  float64
			PeriodikKinerja   float64
			RealisasiBarjas   float64
			TargetBarjas      float64
			RealisasiFisik    float64
			TargetFisik       float64
			RealisasiAnggaran float64
			TargetAnggaran    float64
			RealisasiKinerja  float64
			TargetKinerja     float64
		}) float64 {
			return p.RealisasiBarjas
		}, func(p struct {
			Month             int
			PeriodikBarjas    float64
			PeriodikFisik     float64
			PeriodikAnggaran  float64
			PeriodikKinerja   float64
			RealisasiBarjas   float64
			TargetBarjas      float64
			RealisasiFisik    float64
			TargetFisik       float64
			RealisasiAnggaran float64
			TargetAnggaran    float64
			RealisasiKinerja  float64
			TargetKinerja     float64
		}) float64 {
			return p.TargetBarjas
		}),
		buildCategory("fisik", func(p struct {
			Month             int
			PeriodikBarjas    float64
			PeriodikFisik     float64
			PeriodikAnggaran  float64
			PeriodikKinerja   float64
			RealisasiBarjas   float64
			TargetBarjas      float64
			RealisasiFisik    float64
			TargetFisik       float64
			RealisasiAnggaran float64
			TargetAnggaran    float64
			RealisasiKinerja  float64
			TargetKinerja     float64
		}) float64 {
			return p.PeriodikFisik
		}, func(p struct {
			Month             int
			PeriodikBarjas    float64
			PeriodikFisik     float64
			PeriodikAnggaran  float64
			PeriodikKinerja   float64
			RealisasiBarjas   float64
			TargetBarjas      float64
			RealisasiFisik    float64
			TargetFisik       float64
			RealisasiAnggaran float64
			TargetAnggaran    float64
			RealisasiKinerja  float64
			TargetKinerja     float64
		}) float64 {
			return p.RealisasiFisik
		}, func(p struct {
			Month             int
			PeriodikBarjas    float64
			PeriodikFisik     float64
			PeriodikAnggaran  float64
			PeriodikKinerja   float64
			RealisasiBarjas   float64
			TargetBarjas      float64
			RealisasiFisik    float64
			TargetFisik       float64
			RealisasiAnggaran float64
			TargetAnggaran    float64
			RealisasiKinerja  float64
			TargetKinerja     float64
		}) float64 {
			return p.TargetFisik
		}),
		buildCategory("anggaran", func(p struct {
			Month             int
			PeriodikBarjas    float64
			PeriodikFisik     float64
			PeriodikAnggaran  float64
			PeriodikKinerja   float64
			RealisasiBarjas   float64
			TargetBarjas      float64
			RealisasiFisik    float64
			TargetFisik       float64
			RealisasiAnggaran float64
			TargetAnggaran    float64
			RealisasiKinerja  float64
			TargetKinerja     float64
		}) float64 {
			return p.PeriodikAnggaran
		}, func(p struct {
			Month             int
			PeriodikBarjas    float64
			PeriodikFisik     float64
			PeriodikAnggaran  float64
			PeriodikKinerja   float64
			RealisasiBarjas   float64
			TargetBarjas      float64
			RealisasiFisik    float64
			TargetFisik       float64
			RealisasiAnggaran float64
			TargetAnggaran    float64
			RealisasiKinerja  float64
			TargetKinerja     float64
		}) float64 {
			return p.RealisasiAnggaran
		}, func(p struct {
			Month             int
			PeriodikBarjas    float64
			PeriodikFisik     float64
			PeriodikAnggaran  float64
			PeriodikKinerja   float64
			RealisasiBarjas   float64
			TargetBarjas      float64
			RealisasiFisik    float64
			TargetFisik       float64
			RealisasiAnggaran float64
			TargetAnggaran    float64
			RealisasiKinerja  float64
			TargetKinerja     float64
		}) float64 {
			return p.TargetAnggaran
		}),
		buildCategory("kinerja", func(p struct {
			Month             int
			PeriodikBarjas    float64
			PeriodikFisik     float64
			PeriodikAnggaran  float64
			PeriodikKinerja   float64
			RealisasiBarjas   float64
			TargetBarjas      float64
			RealisasiFisik    float64
			TargetFisik       float64
			RealisasiAnggaran float64
			TargetAnggaran    float64
			RealisasiKinerja  float64
			TargetKinerja     float64
		}) float64 {
			return p.PeriodikKinerja
		}, func(p struct {
			Month             int
			PeriodikBarjas    float64
			PeriodikFisik     float64
			PeriodikAnggaran  float64
			PeriodikKinerja   float64
			RealisasiBarjas   float64
			TargetBarjas      float64
			RealisasiFisik    float64
			TargetFisik       float64
			RealisasiAnggaran float64
			TargetAnggaran    float64
			RealisasiKinerja  float64
			TargetKinerja     float64
		}) float64 {
			return p.RealisasiKinerja
		}, func(p struct {
			Month             int
			PeriodikBarjas    float64
			PeriodikFisik     float64
			PeriodikAnggaran  float64
			PeriodikKinerja   float64
			RealisasiBarjas   float64
			TargetBarjas      float64
			RealisasiFisik    float64
			TargetFisik       float64
			RealisasiAnggaran float64
			TargetAnggaran    float64
			RealisasiKinerja  float64
			TargetKinerja     float64
		}) float64 {
			return p.TargetKinerja
		}),
	}

	return data
}

// GetRealisasiPerbulan retrieves raw realisasi perbulan data for frontend processing
func (m SijagurData) GetRealisasiPerbulan(year, idsatker int) ([]RealisasiData, error) {
	// Fetch raw data from database
	rawData, err := FetchRealisasiPerbulanData(year, idsatker)
	if err != nil {
		return nil, err
	}

	// Process the raw data into RealisasiData format
	return m.ProcessRealisasiPerbulan(rawData), nil
}
