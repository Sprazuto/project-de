package models

import (
	"log"

	"github.com/Massad/gin-boilerplate/db"
)

// getProgressData retrieves progress data from de_ranking_opd
func (m SijagurData) getProgressData(year, month, idsatker int) (ProgressData, error) {
	query := `
		SELECT
			capaian_opd, capaian_barjas, capaian_fisik, capaian_anggaran, capaian_kinerja,
			kumulatif_opd, kumulatif_barjas, kumulatif_fisik, kumulatif_anggaran, kumulatif_kinerja
		FROM de_ranking_opd
		WHERE tahun = $1 AND bulan = $2 AND idsatker = $3
		LIMIT 1
	`

	var progress ProgressData
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

// getBarjasRawData retrieves raw barjas data for frontend processing
func (m SijagurData) getBarjasRawData(year, month, idsatker int, progress float64) (RealisasiData, error) {
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

	formatter := Formatter{}
	// Build raw items for frontend processing
	items := []RealisasiRawItem{
		{Type: "perencanaan", Value: perencanaanSelesai, Detail: &RealisasiDetail{Selesai: perencanaanSelesai, Target: perencanaanTarget, Terlambat: perencanaanTerlambat}},
		{Type: "pemilihan", Value: pemilihanSelesai, Detail: &RealisasiDetail{Selesai: pemilihanSelesai, Target: pemilihanTarget, Terlambat: pemilihanTerlambat}},
		{Type: "pengadaan", Value: pengadaanSelesai, Detail: &RealisasiDetail{Selesai: pengadaanSelesai, Target: pengadaanTarget, Terlambat: pengadaanTerlambat}},
		{Type: "penyerahan", Value: penyerahanSelesai, Detail: &RealisasiDetail{Selesai: penyerahanSelesai, Target: penyerahanTarget, Terlambat: penyerahanTerlambat}},
	}

	return RealisasiData{
		Category:          "barjas",
		Progress:          progress,
		ProgressFormatted: formatter.FormatProgress(progress),
		Items:             items,
	}, nil
}

// getFisikRawData retrieves raw fisik data for frontend processing
func (m SijagurData) getFisikRawData(year, month, idsatker int, progress float64) (RealisasiData, error) {
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

	formatter := Formatter{}
	return RealisasiData{
		Category:          "fisik",
		Progress:          progress,
		ProgressFormatted: formatter.FormatProgress(progress),
		Items: []RealisasiRawItem{
			{Type: "realisasi", Value: realisasi, Formatted: formatter.FormatNumber(realisasi)},
			{Type: "target", Value: target, Formatted: formatter.FormatNumber(target)},
		},
	}, nil
}

// getAnggaranRawData retrieves raw anggaran data for frontend processing
func (m SijagurData) getAnggaranRawData(year, month, idsatker int, progress float64) (RealisasiData, error) {
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

	formatter := Formatter{}
	return RealisasiData{
		Category:          "anggaran",
		Progress:          progress,
		ProgressFormatted: formatter.FormatProgress(progress),
		Items: []RealisasiRawItem{
			{Type: "realisasi", Value: realisasi, Formatted: formatter.FormatCurrency(realisasi)},
			{Type: "target", Value: target, Formatted: formatter.FormatCurrency(target)},
		},
	}, nil
}

// getKinerjaRawData retrieves raw kinerja data for frontend processing
func (m SijagurData) getKinerjaRawData(year, month, idsatker int, progress float64) (RealisasiData, error) {
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

	formatter := Formatter{}
	return RealisasiData{
		Category:          "kinerja",
		Progress:          progress,
		ProgressFormatted: formatter.FormatProgress(progress),
		Items: []RealisasiRawItem{
			{Type: "realisasi", Value: realisasi, Formatted: formatter.FormatNumber(realisasi)},
			{Type: "target", Value: target, Formatted: formatter.FormatNumber(target)},
		},
	}, nil
}

// getBarjasRawDataTahun retrieves raw barjas tahun data for frontend processing
func (m SijagurData) getBarjasRawDataTahun(year, month, idsatker int, progress float64, capaian float64) (RealisasiData, error) {
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

	formatter := Formatter{}
	return RealisasiData{
		Category:          "barjas",
		Progress:          progress,
		ProgressFormatted: formatter.FormatProgress(progress),
		Capaian:           capaian,
		Items: []RealisasiRawItem{
			{Type: "realisasi", Value: totalRealisasi, Formatted: formatter.FormatNumber(totalRealisasi)},
			{Type: "target", Value: totalTarget, Formatted: formatter.FormatNumber(totalTarget)},
		},
	}, nil
}

func (m SijagurData) getFisikRawDataTahun(year, month, idsatker int, progress float64, capaian float64) (RealisasiData, error) {
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

	formatter := Formatter{}
	return RealisasiData{
		Category:          "fisik",
		Progress:          progress,
		ProgressFormatted: formatter.FormatProgress(progress),
		Capaian:           capaian,
		Items: []RealisasiRawItem{
			{Type: "realisasi", Value: realisasi, Formatted: formatter.FormatNumber(realisasi)},
			{Type: "target", Value: target, Formatted: formatter.FormatNumber(target)},
		},
	}, nil
}

func (m SijagurData) getAnggaranRawDataTahun(year, month, idsatker int, progress float64, capaian float64) (RealisasiData, error) {
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

	formatter := Formatter{}
	return RealisasiData{
		Category:          "anggaran",
		Progress:          progress,
		ProgressFormatted: formatter.FormatProgress(progress),
		Capaian:           capaian,
		Items: []RealisasiRawItem{
			{Type: "realisasi", Value: realisasi, Formatted: formatter.FormatCurrency(realisasi)},
			{Type: "target", Value: target, Formatted: formatter.FormatCurrency(target)},
		},
	}, nil
}

func (m SijagurData) getKinerjaRawDataTahun(year, month, idsatker int, progress float64, capaian float64) (RealisasiData, error) {
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

	formatter := Formatter{}
	return RealisasiData{
		Category:          "kinerja",
		Progress:          progress,
		ProgressFormatted: formatter.FormatProgress(progress),
		Capaian:           capaian,
		Items: []RealisasiRawItem{
			{Type: "realisasi", Value: realisasi, Formatted: formatter.FormatNumber(realisasi)},
			{Type: "target", Value: target, Formatted: formatter.FormatNumber(target)},
		},
	}, nil
}

// FetchRealisasiPerbulanData fetches raw monthly progress data for perbulan processing
func FetchRealisasiPerbulanData(year, idsatker int) ([]struct {
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
}, error) {
	// Query all 12 months for the given tahun & idsatker from de_ranking_opd and detail tables.
	rows, err := db.GetDB().Query(`
		SELECT
			dro.bulan,
			dro.periodik_barjas,
			dro.periodik_fisik,
			dro.periodik_anggaran,
			dro.periodik_kinerja,
			COALESCE(ddb.p_barjas_realisasi, 0) as realisasi_barjas,
			COALESCE(ddb.p_barjas_target, 0) as target_barjas,
			COALESCE(ddf.p_fisik_realisasi, 0) as realisasi_fisik,
			COALESCE(ddf.p_fisik_target, 0) as target_fisik,
			COALESCE(dda.p_anggaran_realisasi, 0) as realisasi_anggaran,
			COALESCE(dda.p_anggaran_target, 0) as target_anggaran,
			COALESCE(ddk.p_kinerja_realisasi, 0) as realisasi_kinerja,
			COALESCE(ddk.p_kinerja_target, 0) as target_kinerja
		FROM de_ranking_opd dro
		LEFT JOIN de_detail_barjas ddb ON ddb.id_ranking_opd = dro.id
		LEFT JOIN de_detail_fisik ddf ON ddf.id_ranking_opd = dro.id
		LEFT JOIN de_detail_anggaran dda ON dda.id_ranking_opd = dro.id
		LEFT JOIN de_detail_kinerja ddk ON ddk.id_ranking_opd = dro.id
		WHERE dro.tahun = $1 AND dro.idsatker = $2
		ORDER BY dro.bulan ASC
	`, year, idsatker)
	if err != nil {
		log.Printf("FetchRealisasiPerbulanData: error querying de_ranking_opd: %v", err)
		return nil, err
	}
	defer rows.Close()

	var progressByMonth []struct {
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
	}
	for rows.Next() {
		var p struct {
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
		}
		if scanErr := rows.Scan(&p.Month, &p.PeriodikBarjas, &p.PeriodikFisik, &p.PeriodikAnggaran, &p.PeriodikKinerja, &p.RealisasiBarjas, &p.TargetBarjas, &p.RealisasiFisik, &p.TargetFisik, &p.RealisasiAnggaran, &p.TargetAnggaran, &p.RealisasiKinerja, &p.TargetKinerja); scanErr != nil {
			log.Printf("FetchRealisasiPerbulanData: scan error: %v", scanErr)
			return nil, scanErr
		}
		// Only months 1-12 are relevant; ignore invalid months defensively.
		if p.Month >= 1 && p.Month <= 12 {
			progressByMonth = append(progressByMonth, p)
		}
	}

	return progressByMonth, nil
}
