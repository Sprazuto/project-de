package controllers

import (
	"net/http"
	"strconv"

	"github.com/Massad/gin-boilerplate/forms"
	"github.com/Massad/gin-boilerplate/models"
	"github.com/gin-gonic/gin"
)

// SijagurController ...
type SijagurController struct{}

var sijagurModel = new(models.SijagurData)

// getRealisasiData is a helper function to handle common logic for both bulan and tahun endpoints
func (ctrl SijagurController) getRealisasiData(c *gin.Context, dataType string, getDataFunc func(int, int, int) ([]models.RealisasiData, error)) {
	var queryForm forms.RealisasiQueryForm

	// Bind query parameters
	if err := c.ShouldBindQuery(&queryForm); err != nil {
		sijagurForm := forms.SijagurForm{}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": sijagurForm.ValidateRealisasiQuery(err), "error": err.Error()})
		return
	}

	// Convert to integers with defaults
	tahunInt, bulanInt, idsatkerInt, err := queryForm.ToInts()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invalid query parameters", "error": err.Error()})
		return
	}

	data, err := getDataFunc(tahunInt, bulanInt, idsatkerInt)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Could not get realisasi " + dataType + " data", "error": err.Error()})
		return
	}

	results := []models.SijagurResult{
		{
			Data: data,
			Meta: models.RealisasiMeta{
				Year:      tahunInt,
				Month:     bulanInt,
				MonthName: models.GetMonthName(bulanInt),
				Idsatker:  idsatkerInt,
				Type:      dataType,
			},
		},
	}

	response := models.SijagurResponse{
		Results: results,
	}

	c.JSON(http.StatusOK, response)
}

// GetRealisasiBulan godoc
// @Summary Get Realisasi Bulan data
// @Schemes
// @Description Get Realisasi Bulan data from database
// @Tags Sijagur
// @Accept json
// @Produce json
// @Param tahun query int false "Year"
// @Param bulan query int false "Month"
// @Param idsatker query int false "Satker ID"
// @Success 	 200  {object}  models.RealisasiBulanResponse
// @Failure      400  {object}  gin.H
// @Failure      500  {object}  gin.H
// @Router /realisasi-bulan [GET]
func (ctrl SijagurController) GetRealisasiBulan(c *gin.Context) {
	userID := getUserID(c)
	_ = userID // Mark as used to avoid compiler warning

	ctrl.getRealisasiData(c, "bulan", sijagurModel.GetRealisasiBulanWithParams)
}

// GetRealisasiTahun godoc
// @Summary Get Realisasi Tahun data
// @Schemes
// @Description Get Realisasi Tahun data from database
// @Tags Sijagur
// @Accept json
// @Produce json
// @Param tahun query int false "Year"
// @Param bulan query int false "Month"
// @Param idsatker query int false "Satker ID"
// @Success 	 200  {object}  models.RealisasiTahunResponse
// @Failure      400  {object}  gin.H
// @Failure      500  {object}  gin.H
// @Router /realisasi-tahun [GET]
func (ctrl SijagurController) GetRealisasiTahun(c *gin.Context) {
	userID := getUserID(c)
	_ = userID // Mark as used to avoid compiler warning

	ctrl.getRealisasiData(c, "tahun", sijagurModel.GetRealisasiTahunWithParams)
}

// GetRealisasiPerbulan godoc
// @Summary Get Realisasi Perbulan data (4 categories monthly breakdown)
// @Schemes
// @Description Get Realisasi Perbulan data for Barjas, Fisik, Anggaran, and Kinerja based on tahun and idsatker
// @Tags Sijagur
// @Accept json
// @Produce json
// @Param tahun query int false "Year (default: current year)"
// @Param idsatker query int false "Satker ID (default: 0 for all)"
// @Success 200 {object} models.RealisasiResponse
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /realisasi-perbulan [GET]
func (ctrl SijagurController) GetRealisasiPerbulan(c *gin.Context) {
	userID := getUserID(c)
	_ = userID // Mark as used to avoid compiler warning

	ctrl.getRealisasiData(c, "perbulan", func(year, month, idsatker int) ([]models.RealisasiData, error) {
		return sijagurModel.GetRealisasiPerbulan(year, idsatker)
	})
}

// GetPeringkatKinerja godoc
// @Summary Get Peringkat Kinerja (OPD/SKPD or Kecamatan) with scope
// @Schemes
// @Description Get performance rankings from de_ranking_opd using alias-based scores and jenis_opd-based scope
// @Tags Sijagur
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int false "Month"
// @Param idsatker query int false "Satker ID"
// @Param category query string false "Category filter: all|barjas|fisik|anggaran|kinerja" default(all)
// @Param dimension query string false "Score dimension: kumulatif|capaian|periodik" default(kumulatif)
// @Param scope query string false "Scope: skpd|kecamatan (mapped to jenis_opd)" default(skpd)
// @Param sortBy query string false "Sort by: score_total|score_barjas|score_fisik|score_anggaran|score_kinerja|rank_number"
// @Param sortDir query string false "Sort direction: asc|desc" default(desc)
// @Success 200 {object} models.RankingResponse
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /sijagur/peringkat-kinerja [GET]
func (ctrl SijagurController) GetPeringkatKinerja(c *gin.Context) {
	year, _ := strconv.Atoi(c.DefaultQuery("year", "0"))
	month, _ := strconv.Atoi(c.DefaultQuery("month", "0"))
	idsatker, _ := strconv.Atoi(c.DefaultQuery("idsatker", "0"))
	category := c.DefaultQuery("category", "all")
	dimension := c.DefaultQuery("dimension", "kumulatif")
	scope := c.DefaultQuery("scope", "skpd")
	sortBy := c.DefaultQuery("sortBy", "")
	sortDir := c.DefaultQuery("sortDir", "desc")

	if year <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "year is required",
		})
		return
	}

	resp, err := sijagurModel.GetPeringkatKinerja(
		year,
		month,
		idsatker,
		category,
		dimension,
		scope,
		sortBy,
		sortDir,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "failed to fetch peringkat kinerja",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}
