package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

//UserSessionInfo ...
type UserSessionInfo struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

//JSONRaw ...
type JSONRaw json.RawMessage

//Value ...
func (j JSONRaw) Value() (driver.Value, error) {
	byteArr := []byte(j)
	return driver.Value(byteArr), nil
}

//Scan ...
func (j *JSONRaw) Scan(src interface{}) error {
	asBytes, ok := src.([]byte)
	if !ok {
		return error(errors.New("Scan source was not []bytes"))
	}
	err := json.Unmarshal(asBytes, &j)
	if err != nil {
		return error(errors.New("Scan could not unmarshal to []string"))
	}
	return nil
}

//MarshalJSON ...
func (j *JSONRaw) MarshalJSON() ([]byte, error) {
	return *j, nil
}

//UnmarshalJSON ...
func (j *JSONRaw) UnmarshalJSON(data []byte) error {
	if j == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}
	*j = append((*j)[0:0], data...)
	return nil
}

//DataList ....
type DataList struct {
	Data JSONRaw `db:"data" json:"data"`
	Meta JSONRaw `db:"meta" json:"meta"`
}

// Formatter provides utility functions for formatting data
type Formatter struct{}

// FormatCurrency formats a float64 value as Indonesian Rupiah currency
func (f Formatter) FormatCurrency(amount float64) string {
	s := fmt.Sprintf("%.0f", amount)
	if s == "0" {
		return "Rp0,00"
	}
	// Add dots every 3 digits from the right
	var result []rune
	for i, r := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			result = append(result, '.')
		}
		result = append(result, r)
	}
	return "Rp" + string(result) + ",00"
}

// FormatNumber formats a float64 value as a number string, removing trailing zeros
func (f Formatter) FormatNumber(value float64) string {
	// Remove trailing zeros and decimal point if necessary
	s := fmt.Sprintf("%.2f", value)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	return s
}

// FormatProgress formats a float64 progress value, truncating to 2 decimal places
func (f Formatter) FormatProgress(value float64) string {
	// For progress, round down to 2 decimal places without rounding up
	// Multiply by 100, floor, then divide by 100
	truncated := float64(int(value*100)) / 100
	s := fmt.Sprintf("%.2f", truncated)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	return s
}

var monthNames = []string{
	"", "Januari", "Februari", "Maret", "April", "Mei", "Juni",
	"Juli", "Agustus", "September", "Oktober", "November", "Desember",
}

func GetMonthName(month int) string {
	if month >= 1 && month <= 12 {
		return monthNames[month]
	}
	return "Unknown"
}
