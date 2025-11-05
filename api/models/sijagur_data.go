package models

import (
	"log"

	"github.com/Massad/gin-boilerplate/db"
)

// getProgressData retrieves progress data from de_ranking_opd
func (m SijagurModel) getProgressData(year, month, idsatker int) (ProgressData, error) {
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
func (m SijagurModel) getBarjasRawData(year, month, idsatker int, progress float64) (RealisasiData, error) {
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
		Progress:           progress,
		ProgressFormatted:  formatter.FormatProgress(progress),
		Items:              items,
	}, nil
}

// getFisikRawData retrieves raw fisik data for frontend processing
func (m SijagurModel) getFisikRawData(year, month, idsatker int, progress float64) (RealisasiData, error) {
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
		Progress:           progress,
		ProgressFormatted:  formatter.FormatProgress(progress),
		Items: []RealisasiRawItem{
			{Type: "realisasi", Value: realisasi, Formatted: formatter.FormatNumber(realisasi)},
			{Type: "target", Value: target, Formatted: formatter.FormatNumber(target)},
		},
	}, nil
}

// getAnggaranRawData retrieves raw anggaran data for frontend processing
func (m SijagurModel) getAnggaranRawData(year, month, idsatker int, progress float64) (RealisasiData, error) {
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
		Progress:           progress,
		ProgressFormatted:  formatter.FormatProgress(progress),
		Items: []RealisasiRawItem{
			{Type: "realisasi", Value: realisasi, Formatted: formatter.FormatCurrency(realisasi)},
			{Type: "target", Value: target, Formatted: formatter.FormatCurrency(target)},
		},
	}, nil
}

// getKinerjaRawData retrieves raw kinerja data for frontend processing
func (m SijagurModel) getKinerjaRawData(year, month, idsatker int, progress float64) (RealisasiData, error) {
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
		Progress:           progress,
		ProgressFormatted:  formatter.FormatProgress(progress),
		Items: []RealisasiRawItem{
			{Type: "realisasi", Value: realisasi, Formatted: formatter.FormatNumber(realisasi)},
			{Type: "target", Value: target, Formatted: formatter.FormatNumber(target)},
		},
	}, nil
}

// getBarjasRawDataTahun retrieves raw barjas tahun data for frontend processing
func (m SijagurModel) getBarjasRawDataTahun(year, month, idsatker int, progress float64, capaian float64) (RealisasiData, error) {
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
		Progress:           progress,
		ProgressFormatted:  formatter.FormatProgress(progress),
		Capaian:            capaian,
		Items: []RealisasiRawItem{
			{Type: "realisasi", Value: totalRealisasi, Formatted: formatter.FormatNumber(totalRealisasi)},
			{Type: "target", Value: totalTarget, Formatted: formatter.FormatNumber(totalTarget)},
		},
	}, nil
}

func (m SijagurModel) getFisikRawDataTahun(year, month, idsatker int, progress float64, capaian float64) (RealisasiData, error) {
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
		Progress:           progress,
		ProgressFormatted:  formatter.FormatProgress(progress),
		Capaian:            capaian,
		Items: []RealisasiRawItem{
			{Type: "realisasi", Value: realisasi, Formatted: formatter.FormatNumber(realisasi)},
			{Type: "target", Value: target, Formatted: formatter.FormatNumber(target)},
		},
	}, nil
}

func (m SijagurModel) getAnggaranRawDataTahun(year, month, idsatker int, progress float64, capaian float64) (RealisasiData, error) {
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
		Progress:           progress,
		ProgressFormatted:  formatter.FormatProgress(progress),
		Capaian:            capaian,
		Items: []RealisasiRawItem{
			{Type: "realisasi", Value: realisasi, Formatted: formatter.FormatCurrency(realisasi)},
			{Type: "target", Value: target, Formatted: formatter.FormatCurrency(target)},
		},
	}, nil
}

func (m SijagurModel) getKinerjaRawDataTahun(year, month, idsatker int, progress float64, capaian float64) (RealisasiData, error) {
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
		Progress:           progress,
		ProgressFormatted:  formatter.FormatProgress(progress),
		Capaian:            capaian,
		Items: []RealisasiRawItem{
			{Type: "realisasi", Value: realisasi, Formatted: formatter.FormatNumber(realisasi)},
			{Type: "target", Value: target, Formatted: formatter.FormatNumber(target)},
		},
	}, nil
}
