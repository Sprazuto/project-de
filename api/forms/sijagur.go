package forms

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
)

// SijagurForm ...
type SijagurForm struct{}

// RealisasiQueryForm represents the query parameters for realisasi endpoints
type RealisasiQueryForm struct {
	Tahun    string `form:"tahun" json:"tahun" binding:"omitempty,numeric,min=1900,max=2100"`
	Bulan    string `form:"bulan" json:"bulan" binding:"omitempty,numeric,min=1,max=12"`
	Idsatker string `form:"idsatker" json:"idsatker" binding:"omitempty,numeric,min=0"`
}

// GetDefaultValues returns default values for the form fields
func (f RealisasiQueryForm) GetDefaultValues() (tahun, bulan, idsatker string) {
	now := time.Now()
	currentYear := now.Year()
	currentMonth := int(now.Month())

	tahun = strconv.Itoa(currentYear)
	bulan = strconv.Itoa(currentMonth)
	idsatker = "0"

	return
}

// ToInts converts string fields to integers with defaults
func (f RealisasiQueryForm) ToInts() (tahun, bulan, idsatker int, err error) {
	defaultTahun, defaultBulan, defaultIdsatker := f.GetDefaultValues()

	tahunStr := f.Tahun
	if tahunStr == "" {
		tahunStr = defaultTahun
	}
	bulanStr := f.Bulan
	if bulanStr == "" {
		bulanStr = defaultBulan
	}
	idsatkerStr := f.Idsatker
	if idsatkerStr == "" {
		idsatkerStr = defaultIdsatker
	}

	tahun, err = strconv.Atoi(tahunStr)
	if err != nil {
		return 0, 0, 0, err
	}
	bulan, err = strconv.Atoi(bulanStr)
	if err != nil {
		return 0, 0, 0, err
	}
	idsatker, err = strconv.Atoi(idsatkerStr)
	if err != nil {
		return 0, 0, 0, err
	}

	return tahun, bulan, idsatker, nil
}

// Tahun ...
func (f SijagurForm) Tahun(tag string, errMsg ...string) (message string) {
	switch tag {
	case "required":
		if len(errMsg) == 0 {
			return "Please provide a year"
		}
		return errMsg[0]
	case "numeric":
		return "Year must be a valid number"
	case "min", "max":
		return "Year must be between 1900 and 2100"
	default:
		return "Something went wrong, please try again later"
	}
}

// Bulan ...
func (f SijagurForm) Bulan(tag string, errMsg ...string) (message string) {
	switch tag {
	case "required":
		if len(errMsg) == 0 {
			return "Please provide a month"
		}
		return errMsg[0]
	case "numeric":
		return "Month must be a valid number"
	case "min", "max":
		return "Month must be between 1 and 12"
	default:
		return "Something went wrong, please try again later"
	}
}

// Idsatker ...
func (f SijagurForm) Idsatker(tag string, errMsg ...string) (message string) {
	switch tag {
	case "required":
		if len(errMsg) == 0 {
			return "Please provide a satker ID"
		}
		return errMsg[0]
	case "numeric":
		return "Satker ID must be a valid number"
	case "min":
		return "Satker ID must be 0 or greater"
	default:
		return "Something went wrong, please try again later"
	}
}

// ValidateRealisasiQuery ...
func (f SijagurForm) ValidateRealisasiQuery(err error) string {
	switch err.(type) {
	case validator.ValidationErrors:

		if _, ok := err.(*json.UnmarshalTypeError); ok {
			return "Something went wrong, please try again later"
		}

		for _, e := range err.(validator.ValidationErrors) {
			switch e.Field() {
			case "Tahun":
				return f.Tahun(e.Tag())
			case "Bulan":
				return f.Bulan(e.Tag())
			case "Idsatker":
				return f.Idsatker(e.Tag())
			}
		}

	default:
		return "Invalid request"
	}

	return "Something went wrong, please try again later"
}