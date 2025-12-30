package spse

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Massad/gin-boilerplate/db"
	"github.com/Massad/gin-boilerplate/models"
	"github.com/gin-gonic/gin"
)

// GetStatistics godoc
// @Summary Get scraping statistics from database
// @Schemes
// @Description Get statistics of scraped perencanaan data
// @Tags SPSE
// @Accept json
// @Produce json
// @Success 200 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /spse/statistics [GET]
func (ctrl SPSEController) GetStatistics(c *gin.Context) {
	database := db.GetDB()

	var perencanaanCount int
	var persiapanCount int
	var pemilihanCount int
	var hasilPemilihanCount int
	var kontrakCount int
	var serahTerimaCount int

	// Get counts from each table (excluding soft deleted records)
	err := database.QueryRow("SELECT COUNT(*) FROM spse_perencanaan WHERE deleted_at IS NULL").Scan(&perencanaanCount)
	if err != nil {
		perencanaanCount = 0
	}

	err = database.QueryRow("SELECT COUNT(*) FROM spse_persiapan WHERE deleted_at IS NULL").Scan(&persiapanCount)
	if err != nil {
		persiapanCount = 0
	}

	err = database.QueryRow("SELECT COUNT(*) FROM spse_pemilihan WHERE deleted_at IS NULL").Scan(&pemilihanCount)
	if err != nil {
		pemilihanCount = 0
	}

	err = database.QueryRow("SELECT COUNT(*) FROM spse_hasilpemilihan WHERE deleted_at IS NULL").Scan(&hasilPemilihanCount)
	if err != nil {
		hasilPemilihanCount = 0
	}

	err = database.QueryRow("SELECT COUNT(*) FROM spse_kontrak WHERE deleted_at IS NULL").Scan(&kontrakCount)
	if err != nil {
		kontrakCount = 0
	}

	err = database.QueryRow("SELECT COUNT(*) FROM spse_serahterima WHERE deleted_at IS NULL").Scan(&serahTerimaCount)
	if err != nil {
		serahTerimaCount = 0
	}

	total := perencanaanCount + persiapanCount + pemilihanCount + hasilPemilihanCount + kontrakCount + serahTerimaCount

	c.JSON(http.StatusOK, gin.H{
		"statistics": gin.H{
			"perencanaan":     perencanaanCount,
			"persiapan":       persiapanCount,
			"pemilihan":       pemilihanCount,
			"hasil_pemilihan": hasilPemilihanCount,
			"kontrak":         kontrakCount,
			"serah_terima":    serahTerimaCount,
			"total":           total,
		},
		"last_updated": time.Now().Format(time.RFC3339),
	})
}

// GetPerencanaan godoc
// @Summary Get Perencanaan data from database
// @Schemes
// @Description Retrieve stored perencanaan perencanaan data with pagination
// @Tags SPSE
// @Accept json
// @Produce json
// @Param limit query int false "Limit results (default: 100)"
// @Param offset query int false "Offset for pagination (default: 0)"
// @Success 200 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /spse/data/perencanaan [GET]
func (ctrl SPSEController) GetPerencanaan(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	database := db.GetDB()

	rows, err := database.Query(`
		SELECT id, kode_rup, satuan_kerja, nama_paket, metode_pemilihan,
			   tanggal_pengumuman, rencana_pemilihan, pagu_rup, kode_satuan_kerja,
			   cara_pengadaan, jenis_pengadaan, pdn, umk, sumber_dana, kode_rup_lokal,
			   akhir_pemilihan, tipe_swakelola, active_year, created_at, last_update
		FROM spse_perencanaan
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

	var data []models.SPSEPerencanaan
	for rows.Next() {
		var item models.SPSEPerencanaan
		err := rows.Scan(&item.ID, &item.KodeRUP, &item.SatuanKerja, &item.NamaPaket,
			&item.MetodePemilihan, &item.TanggalPengumuman, &item.RencanaPemilihan,
			&item.PaguRUP, &item.KodeSatuanKerja, &item.CaraPengadaan, &item.JenisPengadaan,
			&item.PDN, &item.UMK, &item.SumberDana, &item.KodeRUPLokal, &item.AkhirPemilihan,
			&item.TipeSwakelola, &item.ActiveYear, &item.CreatedAt, &item.LastUpdate)
		if err != nil {
			continue
		}
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

// GetKontrak godoc
// @Summary Get Kontrak data from database
// @Schemes
// @Description Retrieve stored kontrak contract data with pagination
// @Tags SPSE
// @Accept json
// @Produce json
// @Param limit query int false "Limit results (default: 100)"
// @Param offset query int false "Offset for pagination (default: 0)"
// @Success 200 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /spse/data/kontrak [GET]
func (ctrl SPSEController) GetKontrak(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	database := db.GetDB()

	rows, err := database.Query(`
		SELECT id, kode_rup, satuan_kerja, nama_paket, metode_pemilihan,
			   tanggal_kontrak, nilai_kontrak, status_paket, mulai_kontrak,
			   nilai_bap, selesai_kontrak, kode_satuan_kerja, cara_pengadaan,
			   jenis_pengadaan, pdn, umk, sumber_dana, kode_rup_lokal,
			   metode_pengadaan, tipe_swakelola, active_year, created_at, last_update
		FROM spse_kontrak
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

	var data []models.SPSEKontrak
	for rows.Next() {
		var item models.SPSEKontrak
		err := rows.Scan(&item.ID, &item.KodeRUP, &item.SatuanKerja, &item.NamaPaket,
			&item.MetodePemilihan, &item.TanggalKontrak, &item.NilaiKontrak, &item.StatusPaket,
			&item.MulaiKontrak, &item.NilaiBAP, &item.SelesaiKontrak, &item.KodeSatuanKerja,
			&item.CaraPengadaan, &item.JenisPengadaan, &item.PDN, &item.UMK, &item.SumberDana,
			&item.KodeRUPLokal, &item.MetodePengadaan, &item.TipeSwakelola, &item.ActiveYear, &item.CreatedAt, &item.LastUpdate)
		if err != nil {
			continue
		}
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

// GetSerahTerima godoc
// @Summary Get Serah Terima data from database
// @Schemes
// @Description Retrieve stored serah terima handover data with pagination
// @Tags SPSE
// @Accept json
// @Produce json
// @Param limit query int false "Limit results (default: 100)"
// @Param offset query int false "Offset for pagination (default: 0)"
// @Success 200 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /spse/data/serahterima [GET]
func (ctrl SPSEController) GetSerahTerima(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	database := db.GetDB()

	rows, err := database.Query(`
		SELECT id, kode_rup, satuan_kerja, nama_paket, metode_pemilihan,
			   tanggal_serah_terima, nilai_bap, status_paket, kode_satuan_kerja,
			   cara_pengadaan, jenis_pengadaan, pdn, umk, sumber_dana, kode_rup_lokal,
			   metode_pengadaan, tipe_swakelola, active_year, created_at, last_update
		FROM spse_serahterima
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

	var data []models.SPSESerahTerima
	for rows.Next() {
		var item models.SPSESerahTerima
		err := rows.Scan(&item.ID, &item.KodeRUP, &item.SatuanKerja, &item.NamaPaket,
			&item.MetodePemilihan, &item.TanggalSerahTerima, &item.NilaiBAP, &item.StatusPaket,
			&item.KodeSatuanKerja, &item.CaraPengadaan, &item.JenisPengadaan, &item.PDN,
			&item.UMK, &item.SumberDana, &item.KodeRUPLokal, &item.MetodePengadaan,
			&item.TipeSwakelola, &item.ActiveYear, &item.CreatedAt, &item.LastUpdate)
		if err != nil {
			continue
		}
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

// GetPersiapan godoc
// @Summary Get Persiapan data from database
// @Schemes
// @Description Retrieve stored persiapan perencanaan data with pagination
// @Tags SPSE
// @Accept json
// @Produce json
// @Param limit query int false "Limit results (default: 100)"
// @Param offset query int false "Offset for pagination (default: 0)"
// @Success 200 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /spse/data/persiapan [GET]
func (ctrl SPSEController) GetPersiapan(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	database := db.GetDB()

	rows, err := database.Query(`
		SELECT id, kode_rup, satuan_kerja, nama_paket, metode_pemilihan,
			   tanggal_buat_paket, nilai_pagu_rup, nilai_pagu_paket, kode_satuan_kerja,
			   cara_pengadaan, jenis_pengadaan, pdn, umk, sumber_dana, kode_rup_lokal,
			   metode_pengadaan, tipe_swakelola, active_year, created_at, last_update
		FROM spse_persiapan
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

	var data []models.SPSEPersiapan
	for rows.Next() {
		var item models.SPSEPersiapan
		err := rows.Scan(&item.ID, &item.KodeRUP, &item.SatuanKerja, &item.NamaPaket,
			&item.MetodePemilihan, &item.TanggalBuatPaket, &item.NilaiPaguRUP,
			&item.NilaiPaguPaket, &item.KodeSatuanKerja, &item.CaraPengadaan, &item.JenisPengadaan,
			&item.PDN, &item.UMK, &item.SumberDana, &item.KodeRUPLokal, &item.MetodePengadaan,
			&item.TipeSwakelola, &item.ActiveYear, &item.CreatedAt, &item.LastUpdate)
		if err != nil {
			continue
		}
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

// GetPemilihan godoc
// @Summary Get Pemilihan data from database
// @Schemes
// @Description Retrieve stored pemilihan perencanaan data with pagination
// @Tags SPSE
// @Accept json
// @Produce json
// @Param limit query int false "Limit results (default: 100)"
// @Param offset query int false "Offset for pagination (default: 0)"
// @Success 200 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /spse/data/pemilihan [GET]
func (ctrl SPSEController) GetPemilihan(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	database := db.GetDB()

	rows, err := database.Query(`
		SELECT id, kode_rup, satuan_kerja, nama_paket, metode_pemilihan,
			   rencana_pemilihan, tanggal_pemilihan, nilai_hps, status_paket,
			   kode_satuan_kerja, cara_pengadaan, jenis_pengadaan, pdn, umk,
			   sumber_dana, kode_rup_lokal, metode_pengadaan, pagu_rup,
			   tipe_swakelola, akhir_pemilihan, active_year, created_at, last_update
		FROM spse_pemilihan
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

	var data []models.SPSEPemilihan
	for rows.Next() {
		var item models.SPSEPemilihan
		err := rows.Scan(&item.ID, &item.KodeRUP, &item.SatuanKerja, &item.NamaPaket,
			&item.MetodePemilihan, &item.RencanaPemilihan, &item.TanggalPemilihan,
			&item.NilaiHPS, &item.StatusPaket, &item.KodeSatuanKerja, &item.CaraPengadaan,
			&item.JenisPengadaan, &item.PDN, &item.UMK, &item.SumberDana, &item.KodeRUPLokal,
			&item.MetodePengadaan, &item.PaguRUP, &item.TipeSwakelola, &item.AkhirPemilihan,
			&item.ActiveYear, &item.CreatedAt, &item.LastUpdate)
		if err != nil {
			continue
		}
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

// GetHasilPemilihan godoc
// @Summary Get Hasil Pemilihan data from database
// @Schemes
// @Description Retrieve stored hasil pemilihan election results data with pagination
// @Tags SPSE
// @Accept json
// @Produce json
// @Param limit query int false "Limit results (default: 100)"
// @Param offset query int false "Offset for pagination (default: 0)"
// @Success 200 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /spse/data/hasilpemilihan [GET]
func (ctrl SPSEController) GetHasilPemilihan(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	database := db.GetDB()

	rows, err := database.Query(`
		SELECT id, kode_rup, satuan_kerja, nama_paket, metode_pemilihan,
			   tanggal_hasil_pemilihan, nilai_hasil_pemilihan, status_paket,
			   kode_satuan_kerja, cara_pengadaan, jenis_pengadaan, pdn, umk,
			   sumber_dana, kode_rup_lokal, metode_pengadaan, pagu_rup,
			   tipe_swakelola, active_year, created_at, last_update
		FROM spse_hasilpemilihan
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

	var data []models.SPSEHasilPemilihan
	for rows.Next() {
		var item models.SPSEHasilPemilihan
		err := rows.Scan(&item.ID, &item.KodeRUP, &item.SatuanKerja, &item.NamaPaket,
			&item.MetodePemilihan, &item.TanggalHasilPemilihan, &item.NilaiHasilPemilihan,
			&item.StatusPaket, &item.KodeSatuanKerja, &item.CaraPengadaan, &item.JenisPengadaan,
			&item.PDN, &item.UMK, &item.SumberDana, &item.KodeRUPLokal, &item.MetodePengadaan,
			&item.PaguRUP, &item.TipeSwakelola, &item.ActiveYear, &item.CreatedAt, &item.LastUpdate)
		if err != nil {
			continue
		}
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
