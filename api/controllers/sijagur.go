package controllers

import (
	"net/http"

	"github.com/Massad/gin-boilerplate/forms"
	"github.com/Massad/gin-boilerplate/models"

	"github.com/gin-gonic/gin"
)

// SijagurController ...
type SijagurController struct{}

var sijagurModel = new(models.SijagurModel)

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

	response := models.RealisasiResponse{
		Data: data,
		Meta: models.RealisasiMeta{
			Year:      tahunInt,
			Month:     bulanInt,
			MonthName: models.GetMonthName(bulanInt),
			Idsatker:  idsatkerInt,
			Type:      dataType,
		},
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
	ctrl.getRealisasiData(c, "tahun", sijagurModel.GetRealisasiTahunWithParams)
}