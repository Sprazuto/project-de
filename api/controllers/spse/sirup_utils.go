package spse

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Massad/gin-boilerplate/db"
	"github.com/PuerkitoBio/goquery"
)

// parseCurrencyToFloat converts currency string like "Rp. 426.600" to float64
func parseCurrencyToFloat(currencyStr string) (float64, error) {
	if currencyStr == "" {
		return 0, nil
	}
	// Remove currency symbols and separators
	cleaned := strings.ReplaceAll(currencyStr, "Rp.", "")
	cleaned = strings.ReplaceAll(cleaned, " ", "")
	cleaned = strings.ReplaceAll(cleaned, ".", "")
	cleaned = strings.ReplaceAll(cleaned, ",", ".")
	return strconv.ParseFloat(cleaned, 64)
}

// parseDateToTime converts date string to *time.Time
func parseDateToTime(dateStr string) *time.Time {
	if dateStr == "" {
		return nil
	}

	// Translate Indonesian month names to English
	translatedStr := translateIndonesianMonths(dateStr)

	// Try common formats
	formats := []string{
		"02/01/2006",             // DD/MM/YYYY
		"2006-01-02",             // YYYY-MM-DD
		"02 January 2006",        // DD Month YYYY
		"02 January 2006 15:04",  // With time English
		"2 January 2006",         // Single digit day English
		"January 2006",           // Month Year English
		"06 November 2025 22:05", // Specific format seen
		"24 December 2025",       // Specific format seen
		"24 December 2025 10:22", // Specific format seen
	}
	for _, format := range formats {
		if t, err := time.Parse(format, translatedStr); err == nil {
			return &t
		}
	}
	log.Printf("Warning: Could not parse date '%s' (translated: '%s')", dateStr, translatedStr)
	return nil
}

// translateIndonesianMonths translates Indonesian month names to English
func translateIndonesianMonths(dateStr string) string {
	translations := map[string]string{
		"Januari":   "January",
		"Februari":  "February",
		"Maret":     "March",
		"April":     "April",
		"Mei":       "May",
		"Juni":      "June",
		"Juli":      "July",
		"Agustus":   "August",
		"September": "September",
		"Oktober":   "October",
		"November":  "November",
		"Desember":  "December",
	}

	result := dateStr
	for id, en := range translations {
		result = strings.ReplaceAll(result, id, en)
	}
	return result
}

// SIRUPFieldMapping defines the precise field mapping structure for SIRUP database columns
type SIRUPFieldMapping struct {
	TableName       string                       `json:"table_name"`
	FieldOrder      []string                     `json:"field_order"`
	RequiredFields  map[string]string            `json:"required_fields"`
	FieldValidators map[string]func(string) bool `json:"-"`
}

// GetSIRUPFieldMapping returns the precise field mapping for SIRUP table
func (ctrl SPSEController) GetSIRUPFieldMapping() *SIRUPFieldMapping {
	return &SIRUPFieldMapping{
		TableName: "spse_perencanaansirup",
		FieldOrder: []string{
			"kode_rup", "nama_paket", "nama_klpd", "satuan_kerja", "tahun_anggaran",
			"total_pagu", "lokasi_pekerjaan", "sumber_dana", "jenis_pengadaan",
			"metode_pemilihan", "pemanfaatan_mulai", "pemanfaatan_akhir",
			"jadwal_kontrak_mulai", "jadwal_kontrak_akhir", "jadwal_pemilihan_mulai",
			"jadwal_pemilihan_akhir", "tanggal_umumkan_paket", "sirup_scraped",
		},
		RequiredFields: map[string]string{
			"kode_rup": "", "nama_paket": "", "nama_klpd": "", "satuan_kerja": "",
			"tahun_anggaran": "", "total_pagu": "", "lokasi_pekerjaan": "",
			"sumber_dana": "", "jenis_pengadaan": "", "metode_pemilihan": "",
			"pemanfaatan_mulai": "", "pemanfaatan_akhir": "", "jadwal_kontrak_mulai": "",
			"jadwal_kontrak_akhir": "", "jadwal_pemilihan_mulai": "",
			"jadwal_pemilihan_akhir": "", "tanggal_umumkan_paket": "", "sirup_scraped": "",
		},
		FieldValidators: map[string]func(string) bool{
			"kode_rup": ctrl.isRupCode,
		},
	}
}

// convertMapToSIRUPDataset converts a map to SIRUP ordered dataset format
func (ctrl SPSEController) convertMapToSIRUPDataset(itemMap map[string]interface{}) *OrderedDataSet {
	fieldMapping := ctrl.GetSIRUPFieldMapping()
	dataset := &OrderedDataSet{
		TableName:     fieldMapping.TableName,
		FieldOrder:    fieldMapping.FieldOrder,
		FieldValues:   make(map[string]interface{}),
		OriginalArray: nil,
		MappingStatus: MappingStatus{
			TotalFields:   len(fieldMapping.FieldOrder),
			MappedFields:  0,
			MissingFields: len(fieldMapping.FieldOrder),
			InvalidFields: 0,
		},
	}

	// Fields that should keep their original type (not converted to string)
	keepTypeFields := map[string]bool{
		"total_pagu":            true,
		"tanggal_umumkan_paket": true,
		"sirup_scraped":         true,
	}

	// Copy values from map, validating as we go
	for _, fieldName := range fieldMapping.FieldOrder {
		if value, exists := itemMap[fieldName]; exists {
			// For fields that keep original type, don't convert to string
			if keepTypeFields[fieldName] {
				dataset.FieldValues[fieldName] = value
			} else {
				strValue := fmt.Sprintf("%v", value)

				// Validate field if validator exists
				if validator, exists := fieldMapping.FieldValidators[fieldName]; exists {
					if !validator(strValue) {
						dataset.MappingStatus.InvalidFields++
						continue
					}
				}

				dataset.FieldValues[fieldName] = strValue
			}
			dataset.MappingStatus.MappedFields++
			dataset.MappingStatus.MissingFields--
		}
	}

	// For missing fields, leave them unset (nil) to store null in database

	dataset.MappingStatus.SequencePreserved = dataset.MappingStatus.MappedFields == len(fieldMapping.FieldOrder)
	return dataset
}

// buildSIRUPInsertFromDataset builds INSERT query for SIRUP data using ordered dataset
func (ctrl SPSEController) buildSIRUPInsertFromDataset(dataset *OrderedDataSet) (string, []interface{}) {
	if dataset == nil || dataset.FieldValues == nil {
		return "", nil
	}

	// Extract fields from ordered dataset, keeping original types
	kodeRUP := dataset.FieldValues["kode_rup"]
	namaPaket := dataset.FieldValues["nama_paket"]
	namaKLPD := dataset.FieldValues["nama_klpd"]
	satuanKerja := dataset.FieldValues["satuan_kerja"]
	tahunAnggaran := dataset.FieldValues["tahun_anggaran"]
	totalPagu := dataset.FieldValues["total_pagu"]
	lokasiPekerjaan := dataset.FieldValues["lokasi_pekerjaan"]
	sumberDana := dataset.FieldValues["sumber_dana"]
	jenisPengadaan := dataset.FieldValues["jenis_pengadaan"]
	metodePemilihan := dataset.FieldValues["metode_pemilihan"]
	pemanfaatanMulai := dataset.FieldValues["pemanfaatan_mulai"]
	pemanfaatanAkhir := dataset.FieldValues["pemanfaatan_akhir"]
	jadwalKontrakMulai := dataset.FieldValues["jadwal_kontrak_mulai"]
	jadwalKontrakAkhir := dataset.FieldValues["jadwal_kontrak_akhir"]
	jadwalPemilihanMulai := dataset.FieldValues["jadwal_pemilihan_mulai"]
	jadwalPemilihanAkhir := dataset.FieldValues["jadwal_pemilihan_akhir"]
	tanggalUmumkanPaket := dataset.FieldValues["tanggal_umumkan_paket"]

	query := `INSERT INTO spse_perencanaansirup
		(kode_rup, nama_paket, nama_klpd, satuan_kerja, tahun_anggaran, total_pagu,
		 lokasi_pekerjaan, sumber_dana, jenis_pengadaan, metode_pemilihan,
		 pemanfaatan_mulai, pemanfaatan_akhir, jadwal_kontrak_mulai, jadwal_kontrak_akhir,
		 jadwal_pemilihan_mulai, jadwal_pemilihan_akhir, tanggal_umumkan_paket, sirup_scraped, created_at, last_update, deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, NOW(), $19, NULL)
		ON CONFLICT (kode_rup, nama_paket) DO UPDATE SET
			nama_klpd = EXCLUDED.nama_klpd,
			satuan_kerja = EXCLUDED.satuan_kerja,
			tahun_anggaran = EXCLUDED.tahun_anggaran,
			total_pagu = EXCLUDED.total_pagu,
			lokasi_pekerjaan = EXCLUDED.lokasi_pekerjaan,
			sumber_dana = EXCLUDED.sumber_dana,
			jenis_pengadaan = EXCLUDED.jenis_pengadaan,
			metode_pemilihan = EXCLUDED.metode_pemilihan,
			pemanfaatan_mulai = EXCLUDED.pemanfaatan_mulai,
			pemanfaatan_akhir = EXCLUDED.pemanfaatan_akhir,
			jadwal_kontrak_mulai = EXCLUDED.jadwal_kontrak_mulai,
			jadwal_kontrak_akhir = EXCLUDED.jadwal_kontrak_akhir,
			jadwal_pemilihan_mulai = EXCLUDED.jadwal_pemilihan_mulai,
			jadwal_pemilihan_akhir = EXCLUDED.jadwal_pemilihan_akhir,
			tanggal_umumkan_paket = EXCLUDED.tanggal_umumkan_paket,
			sirup_scraped = EXCLUDED.sirup_scraped,
			last_update = EXCLUDED.last_update,
			deleted_at = NULL`

	sirupScraped := dataset.FieldValues["sirup_scraped"]

	args := []interface{}{kodeRUP, namaPaket, namaKLPD, satuanKerja, tahunAnggaran, totalPagu,
		lokasiPekerjaan, sumberDana, jenisPengadaan, metodePemilihan,
		pemanfaatanMulai, pemanfaatanAkhir, jadwalKontrakMulai, jadwalKontrakAkhir,
		jadwalPemilihanMulai, jadwalPemilihanAkhir, tanggalUmumkanPaket, sirupScraped, time.Now().Unix()}

	return query, args
}

// parseSIRUPHTML parses SIRUP HTML response and extracts structured data
func (ctrl SPSEController) parseSIRUPHTML(htmlContent string, kodeRUP string) (map[string]interface{}, error) {

	// Initialize with default empty values
	data := map[string]interface{}{
		"kode_rup":               kodeRUP,
		"nama_paket":             "",
		"nama_klpd":              "",
		"satuan_kerja":           "",
		"tahun_anggaran":         "",
		"total_pagu":             "",
		"lokasi_pekerjaan":       "[]",
		"sumber_dana_sirup":      "[]",
		"jenis_pengadaan":        "",
		"metode_pemilihan":       "",
		"jadwal_kontrak_mulai":   "",
		"jadwal_kontrak_akhir":   "",
		"pemanfaatan_mulai":      "",
		"pemanfaatan_akhir":      "",
		"jadwal_pemilihan_mulai": "",
		"jadwal_pemilihan_akhir": "",
		"tanggal_umumkan_paket":  "",
	}

	// Parse HTML with goquery
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		log.Printf("Failed to parse HTML for kodeRUP %s: %v", kodeRUP, err)
		return data, fmt.Errorf("failed to parse HTML: %v", err)
	}

	var sumberDana []string
	var lokasi []string
	var jenisPengadaan []string

	// Extract data from rows with label-left class
	doc.Find("td.label-left").Each(func(i int, labelCell *goquery.Selection) {
		label := strings.TrimSpace(strings.ToLower(labelCell.Text()))
		valueCell := labelCell.Next()

		// Handle fields present in SIRUP model
		if strings.Contains(label, "nama paket") {
			paketText := strings.TrimSpace(valueCell.Text())
			data["nama_paket"] = paketText
		} else if strings.Contains(label, "nama klpd") {
			klpdText := strings.TrimSpace(valueCell.Text())
			data["nama_klpd"] = klpdText
		} else if strings.Contains(label, "satuan kerja") {
			satkerText := strings.TrimSpace(valueCell.Text())
			data["satuan_kerja"] = satkerText
		} else if strings.Contains(label, "tahun anggaran") {
			data["tahun_anggaran"] = strings.TrimSpace(valueCell.Text())
		} else if strings.Contains(label, "total pagu") || strings.Contains(label, "nilai pagu") {
			paguText := strings.TrimSpace(valueCell.Text())
			data["total_pagu"] = paguText
		} else if strings.Contains(label, "jenis pengadaan") {
			// Extract jenis pengadaan from the table
			var tdIndex int
			valueCell.Find("td").Each(func(i int, td *goquery.Selection) {
				text := strings.TrimSpace(td.Text())
				// Jenis Pengadaan column (assuming 3 columns: No., Jenis Pengadaan, Pagu)
				if tdIndex%3 == 1 && len(text) > 0 && text != "Jenis Pengadaan" {
					jenisPengadaan = append(jenisPengadaan, text)
				}
				tdIndex++
			})
		} else if strings.Contains(label, "metode pemilihan") || strings.Contains(label, "metode pengadaan") {
			metodeText := strings.TrimSpace(valueCell.Text())
			data["metode_pemilihan"] = metodeText
		} else if strings.Contains(label, "tanggal umumkan paket") {
			data["tanggal_umumkan_paket"] = strings.TrimSpace(valueCell.Text())
		} else if strings.Contains(label, "pemanfaatan barang/jasa") {
			// Extract from nested table - find the data row (not header row)
			valueCell.Find("table tr").Each(func(i int, row *goquery.Selection) {
				cells := row.Find("td")
				if cells.Length() == 2 {
					// This is the data row with actual values
					data["pemanfaatan_mulai"] = strings.TrimSpace(cells.First().Text())
					data["pemanfaatan_akhir"] = strings.TrimSpace(cells.Last().Text())
				}
			})
		} else if strings.Contains(label, "jadwal pelaksanaan kontrak") {
			// Extract from nested table - find the data row (not header row)
			valueCell.Find("table tr").Each(func(i int, row *goquery.Selection) {
				cells := row.Find("td")
				if cells.Length() == 2 {
					// This is the data row with actual values
					data["jadwal_kontrak_mulai"] = strings.TrimSpace(cells.First().Text())
					data["jadwal_kontrak_akhir"] = strings.TrimSpace(cells.Last().Text())
				}
			})
		} else if strings.Contains(label, "jadwal pemilihan penyedia") {
			// Extract from nested table - find the data row (not header row)
			valueCell.Find("table tr").Each(func(i int, row *goquery.Selection) {
				cells := row.Find("td")
				if cells.Length() == 2 {
					// This is the data row with actual values
					data["jadwal_pemilihan_mulai"] = strings.TrimSpace(cells.First().Text())
					data["jadwal_pemilihan_akhir"] = strings.TrimSpace(cells.Last().Text())
				}
			})
		} else if strings.Contains(label, "pra dipa") || strings.Contains(label, "sumber dana") || strings.Contains(label, "pembiayaan") {
			// Extract funding sources from the "Sumber Dana" column
			var tdIndex int
			valueCell.Find("td").Each(func(i int, td *goquery.Selection) {
				text := strings.TrimSpace(td.Text())
				// Sumber Dana column
				if tdIndex%6 == 1 && len(text) > 0 && text != "No." && text != "Sumber Dana" && text != "MAK" {
					sumberDana = append(sumberDana, text)
				}
				tdIndex++
			})
		} else if strings.Contains(label, "Detail Lokasi") || strings.Contains(label, "lokasi") {
			// Extract raw location data from column 4 (index 3)
			var tdIndex int
			valueCell.Find("td").Each(func(i int, td *goquery.Selection) {
				text := strings.TrimSpace(td.Text())
				if tdIndex%4 == 3 && len(text) > 0 {
					lokasi = append(lokasi, text)
				}
				tdIndex++
			})
		}
	})

	// Fallback parsing if some fields are empty
	doc.Find("td").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		lowerText := strings.ToLower(text)

		// Fallback for nama_paket
		if data["nama_paket"] == "" && (strings.Contains(lowerText, "paket") || strings.Contains(lowerText, "pekerjaan")) {
			data["nama_paket"] = text
		}

		// Fallback for nama_klpd
		if data["nama_klpd"] == "" && (strings.Contains(lowerText, "klpd") || strings.Contains(lowerText, "kementerian")) {
			data["nama_klpd"] = text
		}

		// Fallback for satuan_kerja
		if data["satuan_kerja"] == "" && (strings.Contains(lowerText, "satuan kerja") || strings.Contains(lowerText, "satker")) {
			data["satuan_kerja"] = text
		}

		// Fallback for tahun
		if data["tahun_anggaran"] == "" && strings.Contains(text, "202") && len(text) <= 4 {
			data["tahun_anggaran"] = text
		}

		// Fallback for total_pagu
		if data["total_pagu"] == "" && (strings.Contains(lowerText, "pagu") || strings.Contains(lowerText, "nilai")) && strings.Contains(text, "Rp") {
			data["total_pagu"] = text
		}

		// Fallback for metode_pemilihan
		if data["metode_pemilihan"] == "" && (strings.Contains(lowerText, "metode") || strings.Contains(lowerText, "pemilihan")) {
			data["metode_pemilihan"] = text
		}

		// Additional funding sources - removed to prevent table content inclusion
		// if strings.Contains(lowerText, "apbd") || strings.Contains(lowerText, "apbn") {
		//     sumberDana = append(sumberDana, text)
		// }

	})

	// Remove duplicates
	uniqueSumberDana := make(map[string]bool)
	var cleanSumberDana []string
	for _, dana := range sumberDana {
		clean := strings.TrimSpace(dana)
		if !uniqueSumberDana[clean] && len(clean) > 0 {
			uniqueSumberDana[clean] = true
			cleanSumberDana = append(cleanSumberDana, clean)
		}
	}

	// Remove duplicates for jenis pengadaan
	uniqueJenisPengadaan := make(map[string]bool)
	var cleanJenisPengadaan []string
	for _, jenis := range jenisPengadaan {
		clean := strings.TrimSpace(jenis)
		if !uniqueJenisPengadaan[clean] && len(clean) > 0 {
			uniqueJenisPengadaan[clean] = true
			cleanJenisPengadaan = append(cleanJenisPengadaan, clean)
		}
	}

	// Remove duplicates for lokasi
	uniqueLokasi := make(map[string]bool)
	var cleanLokasi []string
	for _, loc := range lokasi {
		clean := strings.TrimSpace(loc)
		if !uniqueLokasi[clean] && len(clean) > 0 {
			uniqueLokasi[clean] = true
			cleanLokasi = append(cleanLokasi, clean)
		}
	}

	// Convert to separated string
	sumberDanaStr := strings.Join(cleanSumberDana, ", ")
	jenisPengadaanStr := strings.Join(cleanJenisPengadaan, ", ")
	separator := ", "
	if len(cleanLokasi) > 1 {
		separator = "|"
	}
	lokasiStr := strings.Join(cleanLokasi, separator)

	data["sumber_dana_sirup"] = sumberDanaStr
	data["jenis_pengadaan"] = jenisPengadaanStr
	data["lokasi_pekerjaan"] = lokasiStr

	return data, nil
}

// getExistingSIRUPRecords retrieves all kodeRUP values from spse_perencanaansirup table
func (ctrl SPSEController) getExistingSIRUPRecords() ([]string, error) {
	database := db.GetDB()
	if database == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	rows, err := database.Query("SELECT kode_rup FROM spse_perencanaansirup WHERE deleted_at IS NULL ORDER BY kode_rup")
	if err != nil {
		return nil, fmt.Errorf("failed to query SIRUP records: %v", err)
	}
	defer rows.Close()

	var kodeRUPs []string
	for rows.Next() {
		var kodeRUP string
		if err := rows.Scan(&kodeRUP); err != nil {
			log.Printf("Error scanning kode_rup: %v", err)
			continue
		}
		kodeRUPs = append(kodeRUPs, kodeRUP)
	}

	return kodeRUPs, nil
}

// prepareSIRUPDataForInsertion prepares SIRUP data for database insertion, setting null for empty fields
func (ctrl SPSEController) prepareSIRUPDataForInsertion(kodeRUP string, sirupData map[string]interface{}) (map[string]interface{}, error) {
	// Initialize with appropriate zero values for all fields
	enrichedData := map[string]interface{}{
		"kode_rup":               kodeRUP,
		"nama_paket":             "",         // string
		"nama_klpd":              nil,        // *string
		"satuan_kerja":           nil,        // *string
		"tahun_anggaran":         "",         // string
		"total_pagu":             float64(0), // float64
		"lokasi_pekerjaan":       "",         // string
		"sumber_dana":            "",         // string
		"jenis_pengadaan":        "",         // string
		"metode_pemilihan":       nil,        // *string
		"pemanfaatan_mulai":      "",         // string
		"pemanfaatan_akhir":      "",         // string
		"jadwal_kontrak_mulai":   "",         // string
		"jadwal_kontrak_akhir":   "",         // string
		"jadwal_pemilihan_mulai": "",         // string
		"jadwal_pemilihan_akhir": "",         // string
		"tanggal_umumkan_paket":  nil,        // *time.Time
		"sirup_scraped":          true,       // bool
	}

	// Override with SIRUP data where available and not empty
	if namaPaket, exists := sirupData["nama_paket"].(string); exists && namaPaket != "" {
		enrichedData["nama_paket"] = namaPaket
	}

	if namaKLPD, exists := sirupData["nama_klpd"].(string); exists && namaKLPD != "" {
		enrichedData["nama_klpd"] = namaKLPD
	}

	if satuanKerja, exists := sirupData["satuan_kerja"].(string); exists && satuanKerja != "" {
		enrichedData["satuan_kerja"] = satuanKerja
	}

	if tahunAnggaran, exists := sirupData["tahun_anggaran"].(string); exists && tahunAnggaran != "" {
		enrichedData["tahun_anggaran"] = tahunAnggaran
	}

	if totalPaguStr, exists := sirupData["total_pagu"].(string); exists && totalPaguStr != "" {
		if totalPagu, err := parseCurrencyToFloat(totalPaguStr); err == nil {
			enrichedData["total_pagu"] = totalPagu
		}
	}

	if lokasiPekerjaan, exists := sirupData["lokasi_pekerjaan"].(string); exists && lokasiPekerjaan != "" && lokasiPekerjaan != "[]" {
		enrichedData["lokasi_pekerjaan"] = lokasiPekerjaan
	}

	if sumberDanaSirup, exists := sirupData["sumber_dana_sirup"].(string); exists && sumberDanaSirup != "" {
		enrichedData["sumber_dana"] = sumberDanaSirup
	}

	if jenisPengadaan, exists := sirupData["jenis_pengadaan"].(string); exists && jenisPengadaan != "" {
		enrichedData["jenis_pengadaan"] = jenisPengadaan
	}

	if metodePemilihan, exists := sirupData["metode_pemilihan"].(string); exists && metodePemilihan != "" {
		enrichedData["metode_pemilihan"] = metodePemilihan
	}

	if jadwalKontrakMulai, exists := sirupData["jadwal_kontrak_mulai"].(string); exists && jadwalKontrakMulai != "" {
		if t := parseDateToTime(jadwalKontrakMulai); t != nil {
			enrichedData["jadwal_kontrak_mulai"] = t.Format("2006-01-02")
		} else {
			enrichedData["jadwal_kontrak_mulai"] = jadwalKontrakMulai
		}
	}

	if jadwalKontrakAkhir, exists := sirupData["jadwal_kontrak_akhir"].(string); exists && jadwalKontrakAkhir != "" {
		if t := parseDateToTime(jadwalKontrakAkhir); t != nil {
			// For "akhir" dates that are month-only, set to last day of month
			if strings.Contains(jadwalKontrakAkhir, t.Format("January")) || strings.Contains(jadwalKontrakAkhir, t.Format("Januari")) ||
				len(strings.Fields(jadwalKontrakAkhir)) == 2 { // Month Year format
				year, month, _ := t.Date()
				lastDay := time.Date(year, month+1, 1, 0, 0, 0, 0, t.Location()).AddDate(0, 0, -1)
				enrichedData["jadwal_kontrak_akhir"] = lastDay.Format("2006-01-02")
			} else {
				enrichedData["jadwal_kontrak_akhir"] = t.Format("2006-01-02")
			}
		} else {
			enrichedData["jadwal_kontrak_akhir"] = jadwalKontrakAkhir
		}
	}

	if pemanfaatanMulai, exists := sirupData["pemanfaatan_mulai"].(string); exists && pemanfaatanMulai != "" {
		if t := parseDateToTime(pemanfaatanMulai); t != nil {
			enrichedData["pemanfaatan_mulai"] = t.Format("2006-01-02")
		} else {
			enrichedData["pemanfaatan_mulai"] = pemanfaatanMulai
		}
	}

	if pemanfaatanAkhir, exists := sirupData["pemanfaatan_akhir"].(string); exists && pemanfaatanAkhir != "" {
		if t := parseDateToTime(pemanfaatanAkhir); t != nil {
			// For "akhir" dates that are month-year only (2 fields), set to last day of month
			if len(strings.Fields(pemanfaatanAkhir)) == 2 { // Month Year format like "November 2025"
				year, month, _ := t.Date()
				lastDay := time.Date(year, month+1, 1, 0, 0, 0, 0, t.Location()).AddDate(0, 0, -1)
				enrichedData["pemanfaatan_akhir"] = lastDay.Format("2006-01-02")
			} else {
				enrichedData["pemanfaatan_akhir"] = t.Format("2006-01-02")
			}
		} else {
			enrichedData["pemanfaatan_akhir"] = pemanfaatanAkhir
		}
	}

	if jadwalPemilihanMulaiStr, exists := sirupData["jadwal_pemilihan_mulai"].(string); exists && jadwalPemilihanMulaiStr != "" {
		if t := parseDateToTime(jadwalPemilihanMulaiStr); t != nil {
			enrichedData["jadwal_pemilihan_mulai"] = t.Format("2006-01-02")
		} else {
			enrichedData["jadwal_pemilihan_mulai"] = jadwalPemilihanMulaiStr
		}
	}

	if jadwalPemilihanAkhirStr, exists := sirupData["jadwal_pemilihan_akhir"].(string); exists && jadwalPemilihanAkhirStr != "" {
		if t := parseDateToTime(jadwalPemilihanAkhirStr); t != nil {
			// For "akhir" dates that are month-only, set to last day of month
			if strings.Contains(jadwalPemilihanAkhirStr, t.Format("January")) || strings.Contains(jadwalPemilihanAkhirStr, t.Format("Januari")) ||
				len(strings.Fields(jadwalPemilihanAkhirStr)) == 2 { // Month Year format
				year, month, _ := t.Date()
				lastDay := time.Date(year, month+1, 1, 0, 0, 0, 0, t.Location()).AddDate(0, 0, -1)
				enrichedData["jadwal_pemilihan_akhir"] = lastDay.Format("2006-01-02")
			} else {
				enrichedData["jadwal_pemilihan_akhir"] = t.Format("2006-01-02")
			}
		} else {
			enrichedData["jadwal_pemilihan_akhir"] = jadwalPemilihanAkhirStr
		}
	}

	if tanggalUmumkanPaketStr, exists := sirupData["tanggal_umumkan_paket"].(string); exists && tanggalUmumkanPaketStr != "" {
		// Parse tanggal_umumkan_paket in Jakarta timezone (+7)
		jakartaLoc, _ := time.LoadLocation("Asia/Jakarta")
		translatedStr := translateIndonesianMonths(tanggalUmumkanPaketStr)
		formats := []string{
			"02 January 2006 15:04",  // DD Month YYYY HH:MM
			"06 November 2025 22:05", // Specific format seen
			"24 December 2025 10:22", // Specific format seen
		}
		for _, format := range formats {
			if t, err := time.ParseInLocation(format, translatedStr, jakartaLoc); err == nil {
				enrichedData["tanggal_umumkan_paket"] = &t
				break
			}
		}
		if enrichedData["tanggal_umumkan_paket"] == nil {
			log.Printf("Warning: Could not parse tanggal_umumkan_paket '%s' (translated: '%s') in Jakarta timezone", tanggalUmumkanPaketStr, translatedStr)
			if t := parseDateToTime(tanggalUmumkanPaketStr); t != nil {
				enrichedData["tanggal_umumkan_paket"] = t // Fallback
			}
		}
	}

	return enrichedData, nil
}
