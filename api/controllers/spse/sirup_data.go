package spse

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Massad/gin-boilerplate/db"
	"github.com/Massad/gin-boilerplate/models"
	"github.com/gin-gonic/gin"
)

// GetSIRUPStatistics godoc
// @Summary Get SIRUP scraping statistics from database
// @Schemes
// @Description Get statistics of scraped SIRUP perencanaan data
// @Tags SIRUP
// @Accept json
// @Produce json
// @Success 200 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /spse/sirup/statistics [GET]
func (ctrl SPSEController) GetSIRUPStatistics(c *gin.Context) {
	database := db.GetDB()

	var sirupCount int

	// Get count from SIRUP table (excluding soft deleted records)
	err := database.QueryRow("SELECT COUNT(*) FROM spse_perencanaansirup WHERE deleted_at IS NULL").Scan(&sirupCount)
	if err != nil {
		sirupCount = 0
	}

	c.JSON(http.StatusOK, gin.H{
		"statistics": gin.H{
			"sirup_perencanaan": sirupCount,
		},
		"last_updated": "2025-12-19T15:00:00Z", // This should be dynamic
	})
}

// GetSIRUPPerencanaan godoc
// @Summary Get SIRUP Perencanaan data from database
// @Schemes
// @Description Retrieve stored SIRUP perencanaan data with pagination
// @Tags SIRUP
// @Accept json
// @Produce json
// @Param limit query int false "Limit results (default: 100)"
// @Param offset query int false "Offset for pagination (default: 0)"
// @Success 200 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /spse/sirup/data/perencanaan [GET]
func (ctrl SPSEController) GetSIRUPPerencanaan(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	database := db.GetDB()

	rows, err := database.Query(`
		SELECT id, kode_rup, nama_paket, nama_klpd, satuan_kerja, tahun_anggaran,
			   total_pagu, lokasi_pekerjaan, sumber_dana, jenis_pengadaan,
			   metode_pemilihan, pemanfaatan_mulai, pemanfaatan_akhir,
			   jadwal_kontrak_mulai, jadwal_kontrak_akhir, jadwal_pemilihan_mulai,
			   jadwal_pemilihan_akhir, tanggal_umumkan_paket, sirup_scraped, active_year, created_at, last_update
		FROM spse_perencanaansirup
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to retrieve data",
		})
		return
	}
	defer rows.Close()

	var data []models.SPSESirupPerencanaan
	for rows.Next() {
		var item models.SPSESirupPerencanaan
		var kodeRUPPtr, namaPaket, namaKLPD, satuanKerja, tahunAnggaran, lokasiPekerjaan, sumberDana, jenisPengadaan, metodePemilihan, pemanfaatanMulai, pemanfaatanAkhir, jadwalKontrakMulai, jadwalKontrakAkhir, jadwalPemilihanMulai, jadwalPemilihanAkhir, activeYear *string
		var totalPagu *float64
		var sirupScraped *bool
		err := rows.Scan(&item.ID, &kodeRUPPtr, &namaPaket, &namaKLPD,
			&satuanKerja, &tahunAnggaran, &totalPagu, &lokasiPekerjaan,
			&sumberDana, &jenisPengadaan, &metodePemilihan,
			&pemanfaatanMulai, &pemanfaatanAkhir, &jadwalKontrakMulai,
			&jadwalKontrakAkhir, &jadwalPemilihanMulai, &jadwalPemilihanAkhir,
			&item.TanggalUmumkanPaket, &sirupScraped, &activeYear, &item.CreatedAt, &item.LastUpdate)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		item.KodeRUP = kodeRUPPtr
		item.NamaPaket = namaPaket
		item.NamaKLPD = namaKLPD
		item.SatuanKerja = satuanKerja
		item.TahunAnggaran = tahunAnggaran
		item.TotalPagu = totalPagu
		item.LokasiPekerjaan = lokasiPekerjaan
		item.SumberDana = sumberDana
		item.JenisPengadaan = jenisPengadaan
		item.MetodePemilihan = metodePemilihan
		item.PemanfaatanMulai = pemanfaatanMulai
		item.PemanfaatanAkhir = pemanfaatanAkhir
		item.JadwalKontrakMulai = jadwalKontrakMulai
		item.JadwalKontrakAkhir = jadwalKontrakAkhir
		item.JadwalPemilihanMulai = jadwalPemilihanMulai
		item.JadwalPemilihanAkhir = jadwalPemilihanAkhir
		item.SirupScraped = sirupScraped
		item.ActiveYear = activeYear
		data = append(data, item)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
		"pagination": gin.H{
			"limit":  limit,
			"offset": offset,
			"count":  len(data),
		},
	})
}

// GetSIRUPByKodeRUP godoc
// @Summary Get SIRUP data for a specific kodeRUP
// @Schemes
// @Description Retrieve SIRUP data for a specific kodeRUP
// @Tags SIRUP
// @Accept json
// @Produce json
// @Param kodeRUP path string true "Kode RUP"
// @Success 200 {object} models.SPSESirupPerencanaan
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /spse/sirup/data/{kodeRUP} [GET]
func (ctrl SPSEController) GetSIRUPByKodeRUP(c *gin.Context) {
	kodeRUP := c.Param("kodeRUP")

	database := db.GetDB()

	var item models.SPSESirupPerencanaan
	var kodeRUPPtr, namaPaket, namaKLPD, satuanKerja, tahunAnggaran, lokasiPekerjaan, sumberDana, jenisPengadaan, metodePemilihan, pemanfaatanMulai, pemanfaatanAkhir, jadwalKontrakMulai, jadwalKontrakAkhir, jadwalPemilihanMulai, jadwalPemilihanAkhir, activeYear *string
	var totalPagu *float64
	var sirupScraped *bool
	err := database.QueryRow(`
		SELECT id, kode_rup, nama_paket, nama_klpd, satuan_kerja, tahun_anggaran,
			   total_pagu, lokasi_pekerjaan, sumber_dana, jenis_pengadaan,
			   metode_pemilihan, pemanfaatan_mulai, pemanfaatan_akhir,
			   jadwal_kontrak_mulai, jadwal_kontrak_akhir, jadwal_pemilihan_mulai,
			   jadwal_pemilihan_akhir, tanggal_umumkan_paket, sirup_scraped, active_year, created_at, last_update
		FROM spse_perencanaansirup
		WHERE kode_rup = $1 AND deleted_at IS NULL
		LIMIT 1
	`, kodeRUP).Scan(
		&item.ID, &kodeRUPPtr, &namaPaket, &namaKLPD,
		&satuanKerja, &tahunAnggaran, &totalPagu, &lokasiPekerjaan,
		&sumberDana, &jenisPengadaan, &metodePemilihan,
		&pemanfaatanMulai, &pemanfaatanAkhir, &jadwalKontrakMulai,
		&jadwalKontrakAkhir, &jadwalPemilihanMulai, &jadwalPemilihanAkhir,
		&item.TanggalUmumkanPaket, &sirupScraped, &activeYear, &item.CreatedAt, &item.LastUpdate)
	if err == nil {
		item.KodeRUP = kodeRUPPtr
		item.NamaPaket = namaPaket
		item.NamaKLPD = namaKLPD
		item.SatuanKerja = satuanKerja
		item.TahunAnggaran = tahunAnggaran
		item.TotalPagu = totalPagu
		item.LokasiPekerjaan = lokasiPekerjaan
		item.SumberDana = sumberDana
		item.JenisPengadaan = jenisPengadaan
		item.MetodePemilihan = metodePemilihan
		item.PemanfaatanMulai = pemanfaatanMulai
		item.PemanfaatanAkhir = pemanfaatanAkhir
		item.JadwalKontrakMulai = jadwalKontrakMulai
		item.JadwalKontrakAkhir = jadwalKontrakAkhir
		item.JadwalPemilihanMulai = jadwalPemilihanMulai
		item.JadwalPemilihanAkhir = jadwalPemilihanAkhir
		item.SirupScraped = sirupScraped
		item.ActiveYear = activeYear
	}

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "SIRUP data not found",
			"message": "No SIRUP data found for the specified kodeRUP",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    item,
	})
}

// GetSIRUPComparison godoc
// @Summary Compare perencanaan and SIRUP data
// @Schemes
// @Description Get comparison between basic perencanaan data and enriched SIRUP data
// @Tags SIRUP
// @Accept json
// @Produce json
// @Param limit query int false "Limit results (default: 50)"
// @Param offset query int false "Offset for pagination (default: 0)"
// @Success 200 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /spse/sirup/comparison [GET]
func (ctrl SPSEController) GetSIRUPComparison(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	database := db.GetDB()

	// Query that joins perencanaan and sirup tables
	rows, err := database.Query(`
		SELECT
			p.kode_rup,
			p.nama_paket,
			p.satuan_kerja,
			p.pagu_rup,
			p.metode_pemilihan,
			s.dates,
			s.sumber_dana_sirup,
			s.lokasi_pekerjaan,
			p.created_at as perencanaan_created,
			s.created_at as sirup_created
		FROM spse_perencanaan p
		LEFT JOIN spse_perencanaansirup s ON p.kode_rup = s.kode_rup AND s.deleted_at IS NULL
		WHERE p.deleted_at IS NULL
		ORDER BY p.created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to retrieve comparison data",
		})
		return
	}
	defer rows.Close()

	type ComparisonItem struct {
		KodeRUP            string `json:"kode_rup"`
		NamaPaket          string `json:"nama_paket"`
		SatuanKerja        string `json:"satuan_kerja"`
		PaguRUP            string `json:"pagu_rup"`
		MetodePemilihan    string `json:"metode_pemilihan"`
		Dates              string `json:"dates"`
		SumberDanaSirup    string `json:"sumber_dana_sirup"`
		LokasiPekerjaan    string `json:"lokasi_pekerjaan"`
		PerencanaanCreated string `json:"perencanaan_created"`
		SirupCreated       string `json:"sirup_created"`
		HasSirupData       bool   `json:"has_sirup_data"`
	}

	var data []ComparisonItem
	for rows.Next() {
		var item ComparisonItem
		var dates, sumberDanaSirup, lokasiPekerjaan *string
		var sirupCreated *string

		err := rows.Scan(
			&item.KodeRUP,
			&item.NamaPaket,
			&item.SatuanKerja,
			&item.PaguRUP,
			&item.MetodePemilihan,
			&dates,
			&sumberDanaSirup,
			&lokasiPekerjaan,
			&item.PerencanaanCreated,
			&sirupCreated,
		)
		if err != nil {
			continue
		}

		if dates != nil {
			item.Dates = *dates
		}
		if sumberDanaSirup != nil {
			item.SumberDanaSirup = *sumberDanaSirup
		}
		if lokasiPekerjaan != nil {
			item.LokasiPekerjaan = *lokasiPekerjaan
		}
		if sirupCreated != nil {
			item.SirupCreated = *sirupCreated
		}

		item.HasSirupData = (dates != nil && *dates != "" && *dates != "[]")

		data = append(data, item)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
		"pagination": gin.H{
			"limit":  limit,
			"offset": offset,
			"count":  len(data),
		},
	})
}
