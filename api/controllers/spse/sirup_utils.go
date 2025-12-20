package spse

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Massad/gin-boilerplate/db"
	"github.com/Massad/gin-boilerplate/models"
	"github.com/PuerkitoBio/goquery"
)

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
			"kode_rup", "satuan_kerja", "nama_paket", "metode_pemilihan",
			"tanggal_pengumuman", "rencana_pemilihan", "pagu_rup", "kode_satuan_kerja",
			"cara_pengadaan", "jenis_pengadaan", "pdn", "umk", "sumber_dana",
			"kode_rup_lokal", "akhir_pemilihan", "tipe_swakelola",
			"dates", "sumber_dana_sirup", "lokasi_pekerjaan",
		},
		RequiredFields: map[string]string{
			"kode_rup": "", "satuan_kerja": "", "nama_paket": "",
			"metode_pemilihan": "", "tanggal_pengumuman": "", "rencana_pemilihan": "",
			"pagu_rup": "", "kode_satuan_kerja": "", "cara_pengadaan": "",
			"jenis_pengadaan": "", "pdn": "", "umk": "",
			"sumber_dana": "", "kode_rup_lokal": "", "akhir_pemilihan": "",
			"tipe_swakelola": "", "dates": "", "sumber_dana_sirup": "", "lokasi_pekerjaan": "",
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

	// Copy values from map, validating as we go
	for _, fieldName := range fieldMapping.FieldOrder {
		if value, exists := itemMap[fieldName]; exists {
			strValue := fmt.Sprintf("%v", value)

			// Validate field if validator exists
			if validator, exists := fieldMapping.FieldValidators[fieldName]; exists {
				if !validator(strValue) {
					dataset.MappingStatus.InvalidFields++
					continue
				}
			}

			dataset.FieldValues[fieldName] = strValue
			dataset.MappingStatus.MappedFields++
			dataset.MappingStatus.MissingFields--
		}
	}

	// Fill missing fields with defaults
	for fieldName, defaultValue := range fieldMapping.RequiredFields {
		if _, exists := dataset.FieldValues[fieldName]; !exists {
			dataset.FieldValues[fieldName] = defaultValue
		}
	}

	dataset.MappingStatus.SequencePreserved = dataset.MappingStatus.MappedFields == len(fieldMapping.FieldOrder)
	return dataset
}

// buildSIRUPInsertFromDataset builds INSERT query for SIRUP data using ordered dataset
func (ctrl SPSEController) buildSIRUPInsertFromDataset(dataset *OrderedDataSet) (string, []interface{}) {
	if dataset == nil || dataset.FieldValues == nil {
		return "", nil
	}

	// Extract fields from ordered dataset
	kodeRUP := dataset.FieldValues["kode_rup"]
	satuanKerja := dataset.FieldValues["satuan_kerja"]
	namaPaket := dataset.FieldValues["nama_paket"]
	metodePemilihan := dataset.FieldValues["metode_pemilihan"]
	tanggalPengumuman := dataset.FieldValues["tanggal_pengumuman"]
	rencanaPemilihan := dataset.FieldValues["rencana_pemilihan"]
	paguRUP := dataset.FieldValues["pagu_rup"]
	kodeSatuanKerja := dataset.FieldValues["kode_satuan_kerja"]
	caraPengadaan := dataset.FieldValues["cara_pengadaan"]
	jenisPengadaan := dataset.FieldValues["jenis_pengadaan"]
	pdn := dataset.FieldValues["pdn"]
	umk := dataset.FieldValues["umk"]
	sumberDana := dataset.FieldValues["sumber_dana"]
	kodeRUPLokal := dataset.FieldValues["kode_rup_lokal"]
	akhirPemilihan := dataset.FieldValues["akhir_pemilihan"]
	tipeSwakelola := dataset.FieldValues["tipe_swakelola"]
	dates := dataset.FieldValues["dates"]
	sumberDanaSirup := dataset.FieldValues["sumber_dana_sirup"]
	lokasiPekerjaan := dataset.FieldValues["lokasi_pekerjaan"]

	query := `INSERT INTO spse_perencanaansirup
		(kode_rup, satuan_kerja, nama_paket, metode_pemilihan, tanggal_pengumuman, rencana_pemilihan, pagu_rup,
		 kode_satuan_kerja, cara_pengadaan, jenis_pengadaan, pdn, umk, sumber_dana, kode_rup_lokal,
		 akhir_pemilihan, tipe_swakelola, dates, sumber_dana_sirup, lokasi_pekerjaan, created_at, last_update, deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, NOW(), $20, NULL)
		ON CONFLICT (kode_rup, nama_paket) DO UPDATE SET
			satuan_kerja = EXCLUDED.satuan_kerja,
			metode_pemilihan = EXCLUDED.metode_pemilihan,
			tanggal_pengumuman = EXCLUDED.tanggal_pengumuman,
			rencana_pemilihan = EXCLUDED.rencana_pemilihan,
			pagu_rup = EXCLUDED.pagu_rup,
			kode_satuan_kerja = EXCLUDED.kode_satuan_kerja,
			cara_pengadaan = EXCLUDED.cara_pengadaan,
			jenis_pengadaan = EXCLUDED.jenis_pengadaan,
			pdn = EXCLUDED.pdn,
			umk = EXCLUDED.umk,
			sumber_dana = EXCLUDED.sumber_dana,
			kode_rup_lokal = EXCLUDED.kode_rup_lokal,
			akhir_pemilihan = EXCLUDED.akhir_pemilihan,
			tipe_swakelola = EXCLUDED.tipe_swakelola,
			dates = EXCLUDED.dates,
			sumber_dana_sirup = EXCLUDED.sumber_dana_sirup,
			lokasi_pekerjaan = EXCLUDED.lokasi_pekerjaan,
			last_update = EXCLUDED.last_update,
			deleted_at = NULL`

	args := []interface{}{kodeRUP, satuanKerja, namaPaket, metodePemilihan, tanggalPengumuman, rencanaPemilihan, paguRUP,
		kodeSatuanKerja, caraPengadaan, jenisPengadaan, pdn, umk, sumberDana, kodeRUPLokal,
		akhirPemilihan, tipeSwakelola, dates, sumberDanaSirup, lokasiPekerjaan, time.Now().Unix()}

	return query, args
}

// parseSIRUPHTML parses SIRUP HTML response and extracts structured data
func (ctrl SPSEController) parseSIRUPHTML(htmlContent string, kodeRUP string) (map[string]interface{}, error) {
	log.Printf("Parsing SIRUP HTML for kodeRUP: %s", kodeRUP)

	// Initialize with default empty values
	data := map[string]interface{}{
		"kode_rup":           kodeRUP,
		"dates":              "[]", // JSON array of dates
		"sumber_dana_sirup":  "[]", // JSON array of funding sources
		"lokasi_pekerjaan":   "[]", // JSON array of work locations
	}

	// Parse HTML with goquery
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		log.Printf("Failed to parse HTML for kodeRUP %s: %v", kodeRUP, err)
		return data, fmt.Errorf("failed to parse HTML: %v", err)
	}

	var dates []string
	var sumberDana []string
	var lokasi []string

	// Extract dates from specific sections
	doc.Find("table").Each(func(tableIndex int, table *goquery.Selection) {
		table.Find("tr").Each(func(rowIndex int, row *goquery.Selection) {
			cells := row.Find("td")
			if cells.Length() >= 2 {
				label := strings.TrimSpace(cells.First().Text())
				value := strings.TrimSpace(cells.Last().Text())

				// Look for date-related fields
				if strings.Contains(strings.ToLower(label), "tanggal") ||
				   strings.Contains(strings.ToLower(label), "jadwal") ||
				   strings.Contains(strings.ToLower(label), "waktu") {
					if len(value) > 0 && (strings.Contains(value, "/") || strings.Contains(value, "-") ||
					   strings.Contains(value, "Januari") || strings.Contains(value, "Februari") ||
					   strings.Contains(value, "Maret") || strings.Contains(value, "April") ||
					   strings.Contains(value, "Mei") || strings.Contains(value, "Juni") ||
					   strings.Contains(value, "Juli") || strings.Contains(value, "Agustus") ||
					   strings.Contains(value, "September") || strings.Contains(value, "Oktober") ||
					   strings.Contains(value, "November") || strings.Contains(value, "Desember")) {
						dates = append(dates, value)
					}
				}
			}
		})
	})

	// Extract funding sources - look for cells containing "APBD"
	doc.Find("td").Each(func(i int, cell *goquery.Selection) {
		text := strings.TrimSpace(cell.Text())
		if text == "APBD" {
			sumberDana = append(sumberDana, text)
		}
	})

	// Extract work locations - look for specific location patterns
	doc.Find("td").Each(func(i int, cell *goquery.Selection) {
		text := strings.TrimSpace(cell.Text())
		// Extract clean location data
		if text == "Sumedang (Kab.)" || (strings.Contains(text, "Badan Keuangan") && !strings.Contains(text, "No.")) {
			// Clean up the text
			cleanText := strings.ReplaceAll(text, "\n", " ")
			cleanText = strings.ReplaceAll(cleanText, "\t", " ")
			// Remove extra whitespace
			for strings.Contains(cleanText, "  ") {
				cleanText = strings.ReplaceAll(cleanText, "  ", " ")
			}
			if len(cleanText) > 0 && len(cleanText) < 100 { // Reasonable length limit
				lokasi = append(lokasi, strings.TrimSpace(cleanText))
			}
		}
	})

	// If no structured data found, try fallback parsing
	if len(dates) == 0 && len(sumberDana) == 0 && len(lokasi) == 0 {
		log.Printf("No structured data found, using fallback parsing for kodeRUP %s", kodeRUP)

		// Fallback: look for any text containing dates
		doc.Find("td").Each(func(i int, s *goquery.Selection) {
			text := strings.TrimSpace(s.Text())
			if len(text) > 5 && (strings.Contains(text, "/") || strings.Contains(text, "-")) {
				if strings.Contains(text, "202") || strings.Contains(text, "202") {
					dates = append(dates, text)
				}
			}
		})

		// Fallback: look for funding sources
		doc.Find("td").Each(func(i int, s *goquery.Selection) {
			text := strings.TrimSpace(s.Text())
			lowerText := strings.ToLower(text)
			if strings.Contains(lowerText, "apbd") || strings.Contains(lowerText, "apbn") ||
			   (strings.Contains(lowerText, "dana") && len(text) > 3) {
				sumberDana = append(sumberDana, text)
			}
		})

		// Fallback: look for locations
		doc.Find("td").Each(func(i int, s *goquery.Selection) {
			text := strings.TrimSpace(s.Text())
			lowerText := strings.ToLower(text)
			if strings.Contains(lowerText, "kabupaten") || strings.Contains(lowerText, "kota") ||
			   strings.Contains(lowerText, "provinsi") || strings.Contains(lowerText, "sumedang") {
				lokasi = append(lokasi, text)
			}
		})
	}

	// Remove duplicates and clean data
	uniqueDates := make(map[string]bool)
	var cleanDates []string
	for _, date := range dates {
		clean := strings.TrimSpace(date)
		if !uniqueDates[clean] && len(clean) > 0 {
			uniqueDates[clean] = true
			cleanDates = append(cleanDates, clean)
		}
	}

	uniqueSumberDana := make(map[string]bool)
	var cleanSumberDana []string
	for _, dana := range sumberDana {
		clean := strings.TrimSpace(dana)
		if !uniqueSumberDana[clean] && len(clean) > 0 {
			uniqueSumberDana[clean] = true
			cleanSumberDana = append(cleanSumberDana, clean)
		}
	}

	uniqueLokasi := make(map[string]bool)
	var cleanLokasi []string
	for _, loc := range lokasi {
		clean := strings.TrimSpace(loc)
		if !uniqueLokasi[clean] && len(clean) > 0 {
			uniqueLokasi[clean] = true
			cleanLokasi = append(cleanLokasi, clean)
		}
	}

	// Convert to JSON strings
	datesJSON, _ := json.Marshal(cleanDates)
	sumberDanaJSON, _ := json.Marshal(cleanSumberDana)
	lokasiJSON, _ := json.Marshal(cleanLokasi)

	data["dates"] = string(datesJSON)
	data["sumber_dana_sirup"] = string(sumberDanaJSON)
	data["lokasi_pekerjaan"] = string(lokasiJSON)

	log.Printf("Extracted for kodeRUP %s: %d dates, %d funding sources, %d locations",
		kodeRUP, len(cleanDates), len(cleanSumberDana), len(cleanLokasi))

	return data, nil
}

// getExistingPerencanaanRecords retrieves all kodeRUP values from spse_perencanaan table
func (ctrl SPSEController) getExistingPerencanaanRecords() ([]string, error) {
	database := db.GetDB()
	if database == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	rows, err := database.Query("SELECT kode_rup FROM spse_perencanaan WHERE deleted_at IS NULL ORDER BY kode_rup")
	if err != nil {
		return nil, fmt.Errorf("failed to query perencanaan records: %v", err)
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

// enrichPerencanaanWithSIRUP combines perencanaan data with SIRUP data
func (ctrl SPSEController) enrichPerencanaanWithSIRUP(kodeRUP string, sirupData map[string]interface{}) (map[string]interface{}, error) {
	database := db.GetDB()
	if database == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	// Get base perencanaan data
	var perencanaan models.SPSEPerencanaan
	err := database.QueryRow(`
		SELECT id, kode_rup, satuan_kerja, nama_paket, metode_pemilihan,
			   tanggal_pengumuman, rencana_pemilihan, pagu_rup, kode_satuan_kerja,
			   cara_pengadaan, jenis_pengadaan, pdn, umk, sumber_dana, kode_rup_lokal,
			   akhir_pemilihan, tipe_swakelola, created_at, last_update
		FROM spse_perencanaan
		WHERE kode_rup = $1 AND deleted_at IS NULL
		LIMIT 1
	`, kodeRUP).Scan(
		&perencanaan.ID, &perencanaan.KodeRUP, &perencanaan.SatuanKerja, &perencanaan.NamaPaket,
		&perencanaan.MetodePemilihan, &perencanaan.TanggalPengumuman, &perencanaan.RencanaPemilihan,
		&perencanaan.PaguRUP, &perencanaan.KodeSatuanKerja, &perencanaan.CaraPengadaan,
		&perencanaan.JenisPengadaan, &perencanaan.PDN, &perencanaan.UMK, &perencanaan.SumberDana,
		&perencanaan.KodeRUPLokal, &perencanaan.AkhirPemilihan, &perencanaan.TipeSwakelola,
		&perencanaan.CreatedAt, &perencanaan.LastUpdate,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get perencanaan data for kodeRUP %s: %v", kodeRUP, err)
	}

	// Combine data
	enrichedData := map[string]interface{}{
		"kode_rup":           perencanaan.KodeRUP,
		"satuan_kerja":       perencanaan.SatuanKerja,
		"nama_paket":         perencanaan.NamaPaket,
		"metode_pemilihan":   perencanaan.MetodePemilihan,
		"tanggal_pengumuman": perencanaan.TanggalPengumuman,
		"rencana_pemilihan":  perencanaan.RencanaPemilihan,
		"pagu_rup":           perencanaan.PaguRUP,
		"kode_satuan_kerja":  perencanaan.KodeSatuanKerja,
		"cara_pengadaan":     perencanaan.CaraPengadaan,
		"jenis_pengadaan":    perencanaan.JenisPengadaan,
		"pdn":                perencanaan.PDN,
		"umk":                perencanaan.UMK,
		"sumber_dana":        perencanaan.SumberDana,
		"kode_rup_lokal":     perencanaan.KodeRUPLokal,
		"akhir_pemilihan":    perencanaan.AkhirPemilihan,
		"tipe_swakelola":     perencanaan.TipeSwakelola,
	}

	// Add SIRUP-specific fields
	if dates, exists := sirupData["dates"]; exists {
		enrichedData["dates"] = dates
	} else {
		enrichedData["dates"] = "[]"
	}

	if sumberDanaSirup, exists := sirupData["sumber_dana_sirup"]; exists {
		enrichedData["sumber_dana_sirup"] = sumberDanaSirup
	} else {
		enrichedData["sumber_dana_sirup"] = "[]"
	}

	if lokasiPekerjaan, exists := sirupData["lokasi_pekerjaan"]; exists {
		enrichedData["lokasi_pekerjaan"] = lokasiPekerjaan
	} else {
		enrichedData["lokasi_pekerjaan"] = "[]"
	}

	return enrichedData, nil
}
