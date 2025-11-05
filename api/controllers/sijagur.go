package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Massad/gin-boilerplate/models"

	"github.com/gin-gonic/gin"
)

// SijagurController ...
type SijagurController struct{}

var sijagurModel = new(models.SijagurModel)

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
// @Failure      500  {object}  gin.H
// @Router /realisasi-bulan [GET]
func (ctrl SijagurController) GetRealisasiBulan(c *gin.Context) {
	now := time.Now()
	currentYear := now.Year()
	currentMonth := int(now.Month())

	tahun := c.DefaultQuery("tahun", strconv.Itoa(currentYear))
	bulan := c.DefaultQuery("bulan", strconv.Itoa(currentMonth))
	idsatker := c.DefaultQuery("idsatker", "0")

	tahunInt, _ := strconv.Atoi(tahun)
	bulanInt, _ := strconv.Atoi(bulan)
	idsatkerInt, _ := strconv.Atoi(idsatker)

	data, err := sijagurModel.GetRealisasiBulanWithParams(tahunInt, bulanInt, idsatkerInt)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Could not get realisasi bulan data", "error": err.Error()})
		return
	}

	response := models.RealisasiResponse{
		Data: data,
		Meta: models.RealisasiMeta{
			Year:      tahunInt,
			Month:     bulanInt,
			MonthName: sijagurModel.GetMonthName(bulanInt),
			Idsatker:  idsatkerInt,
			Type:      "bulan",
		},
	}

	c.JSON(http.StatusOK, response)
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
// @Failure      500  {object}  gin.H
// @Router /realisasi-tahun [GET]
func (ctrl SijagurController) GetRealisasiTahun(c *gin.Context) {
	now := time.Now()
	currentYear := now.Year()
	currentMonth := int(now.Month())

	tahun := c.DefaultQuery("tahun", strconv.Itoa(currentYear))
	bulan := c.DefaultQuery("bulan", strconv.Itoa(currentMonth))
	idsatker := c.DefaultQuery("idsatker", "0")

	tahunInt, _ := strconv.Atoi(tahun)
	bulanInt, _ := strconv.Atoi(bulan)
	idsatkerInt, _ := strconv.Atoi(idsatker)

	data, err := sijagurModel.GetRealisasiTahunWithParams(tahunInt, bulanInt, idsatkerInt)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Could not get realisasi tahun data", "error": err.Error()})
		return
	}

	response := models.RealisasiResponse{
		Data: data,
		Meta: models.RealisasiMeta{
			Year:      tahunInt,
			Month:     bulanInt,
			MonthName: sijagurModel.GetMonthName(bulanInt),
			Idsatker:  idsatkerInt,
			Type:      "tahun",
		},
	}

	c.JSON(http.StatusOK, response)
}