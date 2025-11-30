package spse

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// SPSEFieldMapping defines the precise field mapping structure for database columns
type SPSEFieldMapping struct {
	TableName       string                       `json:"table_name"`
	FieldOrder      []string                     `json:"field_order"`
	RequiredFields  map[string]string            `json:"required_fields"`
	FieldValidators map[string]func(string) bool `json:"-"`
}

// OrderedDataSet represents an ordered dataset with proper field mapping
type OrderedDataSet struct {
	TableName     string                 `json:"table_name"`
	FieldOrder    []string               `json:"field_order"`
	FieldValues   map[string]interface{} `json:"field_values"`
	OriginalArray []interface{}          `json:"original_array"`
	MappingStatus MappingStatus          `json:"mapping_status"`
}

// MappingStatus tracks the quality of the field mapping
type MappingStatus struct {
	TotalFields       int  `json:"total_fields"`
	MappedFields      int  `json:"mapped_fields"`
	MissingFields     int  `json:"missing_fields"`
	InvalidFields     int  `json:"invalid_fields"`
	SequencePreserved bool `json:"sequence_preserved"`
}

// GetFieldMapping returns the precise field mapping for a specific table
func (ctrl SPSEController) GetFieldMapping(tableName string) *SPSEFieldMapping {
	mappings := map[string]*SPSEFieldMapping{
		"perencanaan": {
			TableName: "spse_perencanaan",
			FieldOrder: []string{
				"kode_rup", "satuan_kerja", "nama_paket", "metode_pemilihan",
				"tanggal_pengumuman", "rencana_pemilihan", "pagu_rup", "kode_satuan_kerja",
				"cara_pengadaan", "jenis_pengadaan", "pdn", "umk", "sumber_dana",
				"kode_rup_lokal", "akhir_pemilihan", "tipe_swakelola",
			},
			RequiredFields: map[string]string{
				"kode_rup": "", "satuan_kerja": "", "nama_paket": "",
				"metode_pemilihan": "", "tanggal_pengumuman": "", "rencana_pemilihan": "",
				"pagu_rup": "", "kode_satuan_kerja": "", "cara_pengadaan": "",
				"jenis_pengadaan": "", "pdn": "", "umk": "",
				"sumber_dana": "", "kode_rup_lokal": "", "akhir_pemilihan": "",
				"tipe_swakelola": "",
			},
			FieldValidators: map[string]func(string) bool{
				"kode_rup":           ctrl.isRupCode,
				"tanggal_pengumuman": ctrl.isDate,
				"pagu_rup":           ctrl.isCurrency,
			},
		},
		"persiapan": {
			TableName: "spse_persiapan",
			FieldOrder: []string{
				"kode_rup", "satuan_kerja", "nama_paket", "metode_pemilihan",
				"tanggal_buat_paket", "nilai_pagu_rup", "nilai_pagu_paket", "kode_satuan_kerja",
				"cara_pengadaan", "jenis_pengadaan", "pdn", "umk", "sumber_dana",
				"kode_rup_lokal", "metode_pengadaan", "tipe_swakelola",
			},
			RequiredFields: map[string]string{
				"kode_rup": "", "satuan_kerja": "", "nama_paket": "",
				"metode_pemilihan": "", "tanggal_buat_paket": "", "nilai_pagu_rup": "",
				"nilai_pagu_paket": "", "kode_satuan_kerja": "", "cara_pengadaan": "",
				"jenis_pengadaan": "", "pdn": "", "umk": "",
				"sumber_dana": "", "kode_rup_lokal": "", "metode_pengadaan": "",
				"tipe_swakelola": "",
			},
			FieldValidators: map[string]func(string) bool{
				"kode_rup":           ctrl.isRupCode,
				"tanggal_buat_paket": ctrl.isDate,
				"nilai_pagu_rup":     ctrl.isCurrency,
				"nilai_pagu_paket":   ctrl.isCurrency,
			},
		},
		"pemilihan": {
			TableName: "spse_pemilihan",
			FieldOrder: []string{
				"kode_rup", "satuan_kerja", "nama_paket", "metode_pemilihan",
				"rencana_pemilihan", "tanggal_pemilihan", "nilai_hps", "status_paket",
				"kode_satuan_kerja", "cara_pengadaan", "jenis_pengadaan", "pdn", "umk",
				"sumber_dana", "kode_rup_lokal", "metode_pengadaan", "pagu_rup",
				"tipe_swakelola", "akhir_pemilihan",
			},
			RequiredFields: map[string]string{
				"kode_rup": "", "satuan_kerja": "", "nama_paket": "",
				"metode_pemilihan": "", "rencana_pemilihan": "", "tanggal_pemilihan": "",
				"nilai_hps": "", "status_paket": "", "kode_satuan_kerja": "",
				"cara_pengadaan": "", "jenis_pengadaan": "", "pdn": "",
				"umk": "", "sumber_dana": "", "kode_rup_lokal": "",
				"metode_pengadaan": "", "pagu_rup": "", "tipe_swakelola": "",
				"akhir_pemilihan": "",
			},
			FieldValidators: map[string]func(string) bool{
				"kode_rup":          ctrl.isRupCode,
				"tanggal_pemilihan": ctrl.isDate,
				"nilai_hps":         ctrl.isCurrency,
				"pagu_rup":          ctrl.isCurrency,
			},
		},
		"hasilpemilihan": {
			TableName: "spse_hasilpemilihan",
			FieldOrder: []string{
				"kode_rup", "satuan_kerja", "nama_paket", "metode_pemilihan",
				"tanggal_hasil_pemilihan", "nilai_hasil_pemilihan", "status_paket",
				"kode_satuan_kerja", "cara_pengadaan", "jenis_pengadaan", "pdn", "umk",
				"sumber_dana", "kode_rup_lokal", "metode_pengadaan", "pagu_rup",
				"tipe_swakelola",
			},
			RequiredFields: map[string]string{
				"kode_rup": "", "satuan_kerja": "", "nama_paket": "",
				"metode_pemilihan": "", "tanggal_hasil_pemilihan": "", "nilai_hasil_pemilihan": "",
				"status_paket": "", "kode_satuan_kerja": "", "cara_pengadaan": "",
				"jenis_pengadaan": "", "pdn": "", "umk": "",
				"sumber_dana": "", "kode_rup_lokal": "", "metode_pengadaan": "",
				"pagu_rup": "", "tipe_swakelola": "",
			},
			FieldValidators: map[string]func(string) bool{
				"kode_rup":                ctrl.isRupCode,
				"tanggal_hasil_pemilihan": ctrl.isDate,
				"nilai_hasil_pemilihan":   ctrl.isCurrency,
				"pagu_rup":                ctrl.isCurrency,
			},
		},
		"kontrak": {
			TableName: "spse_kontrak",
			FieldOrder: []string{
				"kode_rup", "satuan_kerja", "nama_paket", "metode_pemilihan",
				"tanggal_kontrak", "nilai_kontrak", "status_paket", "mulai_kontrak",
				"nilai_bap", "selesai_kontrak", "kode_satuan_kerja", "cara_pengadaan",
				"jenis_pengadaan", "pdn", "umk", "sumber_dana", "kode_rup_lokal",
				"metode_pengadaan", "tipe_swakelola",
			},
			RequiredFields: map[string]string{
				"kode_rup": "", "satuan_kerja": "", "nama_paket": "",
				"metode_pemilihan": "", "tanggal_kontrak": "", "nilai_kontrak": "",
				"status_paket": "", "mulai_kontrak": "", "nilai_bap": "",
				"selesai_kontrak": "", "kode_satuan_kerja": "", "cara_pengadaan": "",
				"jenis_pengadaan": "", "pdn": "", "umk": "",
				"sumber_dana": "", "kode_rup_lokal": "", "metode_pengadaan": "",
				"tipe_swakelola": "",
			},
			FieldValidators: map[string]func(string) bool{
				"kode_rup":        ctrl.isRupCode,
				"tanggal_kontrak": ctrl.isDate,
				"nilai_kontrak":   ctrl.isCurrency,
				"nilai_bap":       ctrl.isCurrency,
			},
		},
		"serahterima": {
			TableName: "spse_serahterima",
			FieldOrder: []string{
				"kode_rup", "satuan_kerja", "nama_paket", "metode_pemilihan",
				"tanggal_serah_terima", "nilai_bap", "status_paket", "kode_satuan_kerja",
				"cara_pengadaan", "jenis_pengadaan", "pdn", "umk", "sumber_dana",
				"kode_rup_lokal", "metode_pengadaan", "tipe_swakelola",
			},
			RequiredFields: map[string]string{
				"kode_rup": "", "satuan_kerja": "", "nama_paket": "",
				"metode_pemilihan": "", "tanggal_serah_terima": "", "nilai_bap": "",
				"status_paket": "", "kode_satuan_kerja": "", "cara_pengadaan": "",
				"jenis_pengadaan": "", "pdn": "", "umk": "",
				"sumber_dana": "", "kode_rup_lokal": "", "metode_pengadaan": "",
				"tipe_swakelola": "",
			},
			FieldValidators: map[string]func(string) bool{
				"kode_rup":             ctrl.isRupCode,
				"tanggal_serah_terima": ctrl.isDate,
				"nilai_bap":            ctrl.isCurrency,
			},
		},
	}

	return mappings[tableName]
}

// convertMapToOrderedDataset converts a map to ordered dataset format
func (ctrl SPSEController) convertMapToOrderedDataset(tableName string, itemMap map[string]interface{}) *OrderedDataSet {
	fieldMapping := ctrl.GetFieldMapping(tableName)
	if fieldMapping == nil {
		return &OrderedDataSet{
			TableName:     tableName,
			FieldOrder:    []string{},
			FieldValues:   itemMap,
			OriginalArray: nil,
			MappingStatus: MappingStatus{TotalFields: 0, MappedFields: 0, MissingFields: 0, InvalidFields: 0},
		}
	}

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

// parseStringData parses string data that might be delimited or formatted
func (ctrl SPSEController) parseStringData(data string) map[string]interface{} {
	itemMap := make(map[string]interface{})

	// Try to split by common delimiters and extract meaningful data
	// This is a basic implementation - you may need to adjust based on actual data format

	// Look for patterns that might indicate RUP codes (numeric)
	re := regexp.MustCompile(`(\d{8,})`)
	matches := re.FindStringSubmatch(data)
	if len(matches) > 0 {
		itemMap["kode_rup"] = matches[1]
	}

	// Try to extract currency amounts
	currencyRe := regexp.MustCompile(`Rp\.?\s*([\d,.]+)`)
	currencyMatches := currencyRe.FindStringSubmatch(data)
	if len(currencyMatches) > 1 {
		itemMap["pagu_rup"] = "Rp. " + currencyMatches[1]
	}

	// Extract dates (basic pattern for Indonesian date format)
	dateRe := regexp.MustCompile(`(\d{1,2}\s+(?:Januari|Februari|Maret|April|Mei|Juni|Juli|Agustus|September|Oktober|November|Desember)\s+\d{4})`)
	dateMatches := dateRe.FindStringSubmatch(data)
	if len(dateMatches) > 0 {
		itemMap["tanggal_pengumuman"] = dateMatches[1]
	}

	// If we couldn't parse specific fields, put the whole string in nama_paket as fallback
	if len(itemMap) == 0 {
		itemMap["nama_paket"] = data
	}

	return itemMap
}

// mapArrayToOrderedDataset converts unkeyed array to ordered dataset with precise field mapping
func (ctrl SPSEController) mapArrayToOrderedDataset(tableName string, arr []interface{}) *OrderedDataSet {
	fieldMapping := ctrl.GetFieldMapping(tableName)
	if fieldMapping == nil {
		return &OrderedDataSet{
			TableName:     tableName,
			FieldOrder:    []string{},
			FieldValues:   make(map[string]interface{}),
			OriginalArray: arr,
			MappingStatus: MappingStatus{TotalFields: 0, MappedFields: 0, MissingFields: 0, InvalidFields: 0},
		}
	}

	// Convert array to strings for processing
	arrStr := make([]string, len(arr))
	for i, v := range arr {
		if v != nil {
			arrStr[i] = fmt.Sprintf("%v", v)
		} else {
			arrStr[i] = ""
		}
	}

	// Create ordered dataset with precise field mapping
	dataset := &OrderedDataSet{
		TableName:     fieldMapping.TableName,
		FieldOrder:    fieldMapping.FieldOrder,
		FieldValues:   make(map[string]interface{}),
		OriginalArray: arr,
		MappingStatus: MappingStatus{
			TotalFields:   len(fieldMapping.FieldOrder),
			MappedFields:  0,
			MissingFields: len(fieldMapping.FieldOrder),
			InvalidFields: 0,
		},
	}

	// Map array elements to database fields based on precise field order
	for i, fieldName := range fieldMapping.FieldOrder {
		if i < len(arrStr) && arrStr[i] != "" {
			value := strings.TrimSpace(arrStr[i])

			// Validate field if validator exists
			if validator, exists := fieldMapping.FieldValidators[fieldName]; exists {
				if !validator(value) {
					dataset.MappingStatus.InvalidFields++
					continue // Skip invalid values
				}
			}

			dataset.FieldValues[fieldName] = value
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

	// Check if sequence is preserved (all expected fields mapped)
	dataset.MappingStatus.SequencePreserved = dataset.MappingStatus.MappedFields == len(fieldMapping.FieldOrder)

	return dataset
}

// Helper functions for pattern recognition

// isRupCode checks if a string looks like a RUP code
func (ctrl SPSEController) isRupCode(value string) bool {
	// RUP codes are typically numeric and 8+ digits
	_, err := strconv.ParseInt(value, 10, 64)
	return err == nil && len(value) >= 8
}

// isDate checks if a string looks like a date
func (ctrl SPSEController) isDate(value string) bool {
	// Check for common date patterns
	datePatterns := []string{
		`\d{1,2}/\d{1,2}/\d{4}`, // DD/MM/YYYY
		`\d{4}-\d{1,2}-\d{1,2}`, // YYYY-MM-DD
		`\d{1,2}-\d{1,2}-\d{4}`, // DD-MM-YYYY
	}

	for _, pattern := range datePatterns {
		matched, _ := regexp.MatchString(pattern, value)
		if matched {
			return true
		}
	}

	// Check for Indonesian date formats
	months := []string{"januari", "februari", "maret", "april", "mei", "juni",
		"juli", "agustus", "september", "oktober", "november", "desember"}
	lowerValue := strings.ToLower(value)
	for _, month := range months {
		if strings.Contains(lowerValue, month) {
			return true
		}
	}

	return false
}

// isCurrency checks if a string looks like a currency amount
func (ctrl SPSEController) isCurrency(value string) bool {
	// Check for currency patterns
	currencyPatterns := []string{
		`Rp\.?\s*\d+`,   // Rp 1000
		`\d+(\.\d{3})*`, // 1000 or 1.000.000
		`\d+,\d{3}`,     // 1,000
	}

	for _, pattern := range currencyPatterns {
		matched, _ := regexp.MatchString(pattern, value)
		if matched {
			return true
		}
	}

	return false
}

// isSelectionMethod checks if a string looks like a selection method
func (ctrl SPSEController) isSelectionMethod(value string) bool {
	methods := []string{"e-tendering", "e-purchasing", "e-catalog", "swakelola", "pengadaan langsung"}
	lowerValue := strings.ToLower(value)
	for _, method := range methods {
		if strings.Contains(lowerValue, method) {
			return true
		}
	}
	return false
}

// buildPerencanaanInsertFromDataset builds INSERT query for perencanaan data using ordered dataset
func (ctrl SPSEController) buildPerencanaanInsertFromDataset(dataset *OrderedDataSet) (string, []interface{}) {
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

	query := `INSERT INTO spse_perencanaan
		(kode_rup, satuan_kerja, nama_paket, metode_pemilihan, tanggal_pengumuman, rencana_pemilihan, pagu_rup,
		 kode_satuan_kerja, cara_pengadaan, jenis_pengadaan, pdn, umk, sumber_dana, kode_rup_lokal,
		 akhir_pemilihan, tipe_swakelola, created_at, last_update, deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, NOW(), $17, NULL)
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
			last_update = EXCLUDED.last_update,
			deleted_at = NULL`

	args := []interface{}{kodeRUP, satuanKerja, namaPaket, metodePemilihan, tanggalPengumuman, rencanaPemilihan, paguRUP,
		kodeSatuanKerja, caraPengadaan, jenisPengadaan, pdn, umk, sumberDana, kodeRUPLokal,
		akhirPemilihan, tipeSwakelola, time.Now().Unix()}

	return query, args
}

// buildPersiapanInsertFromDataset builds INSERT query for persiapan data using ordered dataset
func (ctrl SPSEController) buildPersiapanInsertFromDataset(dataset *OrderedDataSet) (string, []interface{}) {
	if dataset == nil || dataset.FieldValues == nil {
		return "", nil
	}

	// Extract fields from ordered dataset
	kodeRUP := dataset.FieldValues["kode_rup"]
	satuanKerja := dataset.FieldValues["satuan_kerja"]
	namaPaket := dataset.FieldValues["nama_paket"]
	metodePemilihan := dataset.FieldValues["metode_pemilihan"]
	tanggalBuatPaket := dataset.FieldValues["tanggal_buat_paket"]
	nilaiPaguRUP := dataset.FieldValues["nilai_pagu_rup"]
	nilaiPaguPaket := dataset.FieldValues["nilai_pagu_paket"]
	kodeSatuanKerja := dataset.FieldValues["kode_satuan_kerja"]
	caraPengadaan := dataset.FieldValues["cara_pengadaan"]
	jenisPengadaan := dataset.FieldValues["jenis_pengadaan"]
	pdn := dataset.FieldValues["pdn"]
	umk := dataset.FieldValues["umk"]
	sumberDana := dataset.FieldValues["sumber_dana"]
	kodeRUPLokal := dataset.FieldValues["kode_rup_lokal"]
	metodePengadaan := dataset.FieldValues["metode_pengadaan"]
	tipeSwakelola := dataset.FieldValues["tipe_swakelola"]

	query := `INSERT INTO spse_persiapan
		(kode_rup, satuan_kerja, nama_paket, metode_pemilihan, tanggal_buat_paket, nilai_pagu_rup, nilai_pagu_paket,
		 kode_satuan_kerja, cara_pengadaan, jenis_pengadaan, pdn, umk, sumber_dana, kode_rup_lokal,
		 metode_pengadaan, tipe_swakelola, created_at, last_update, deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, NOW(), $17, NULL)
		ON CONFLICT (kode_rup, nama_paket) DO UPDATE SET
			satuan_kerja = EXCLUDED.satuan_kerja,
			metode_pemilihan = EXCLUDED.metode_pemilihan,
			tanggal_buat_paket = EXCLUDED.tanggal_buat_paket,
			nilai_pagu_rup = EXCLUDED.nilai_pagu_rup,
			nilai_pagu_paket = EXCLUDED.nilai_pagu_paket,
			kode_satuan_kerja = EXCLUDED.kode_satuan_kerja,
			cara_pengadaan = EXCLUDED.cara_pengadaan,
			jenis_pengadaan = EXCLUDED.jenis_pengadaan,
			pdn = EXCLUDED.pdn,
			umk = EXCLUDED.umk,
			sumber_dana = EXCLUDED.sumber_dana,
			kode_rup_lokal = EXCLUDED.kode_rup_lokal,
			metode_pengadaan = EXCLUDED.metode_pengadaan,
			tipe_swakelola = EXCLUDED.tipe_swakelola,
			last_update = EXCLUDED.last_update,
			deleted_at = NULL`

	args := []interface{}{kodeRUP, satuanKerja, namaPaket, metodePemilihan, tanggalBuatPaket, nilaiPaguRUP, nilaiPaguPaket,
		kodeSatuanKerja, caraPengadaan, jenisPengadaan, pdn, umk, sumberDana, kodeRUPLokal,
		metodePengadaan, tipeSwakelola, time.Now().Unix()}

	return query, args
}

// buildPemilihanInsertFromDataset builds INSERT query for pemilihan data using ordered dataset
func (ctrl SPSEController) buildPemilihanInsertFromDataset(dataset *OrderedDataSet) (string, []interface{}) {
	if dataset == nil || dataset.FieldValues == nil {
		return "", nil
	}

	// Extract fields from ordered dataset
	kodeRUP := dataset.FieldValues["kode_rup"]
	satuanKerja := dataset.FieldValues["satuan_kerja"]
	namaPaket := dataset.FieldValues["nama_paket"]
	metodePemilihan := dataset.FieldValues["metode_pemilihan"]
	rencanaPemilihan := dataset.FieldValues["rencana_pemilihan"]
	tanggalPemilihan := dataset.FieldValues["tanggal_pemilihan"]
	nilaiHPS := dataset.FieldValues["nilai_hps"]
	statusPaket := dataset.FieldValues["status_paket"]
	kodeSatuanKerja := dataset.FieldValues["kode_satuan_kerja"]
	caraPengadaan := dataset.FieldValues["cara_pengadaan"]
	jenisPengadaan := dataset.FieldValues["jenis_pengadaan"]
	pdn := dataset.FieldValues["pdn"]
	umk := dataset.FieldValues["umk"]
	sumberDana := dataset.FieldValues["sumber_dana"]
	kodeRUPLokal := dataset.FieldValues["kode_rup_lokal"]
	metodePengadaan := dataset.FieldValues["metode_pengadaan"]
	paguRUP := dataset.FieldValues["pagu_rup"]
	tipeSwakelola := dataset.FieldValues["tipe_swakelola"]
	akhirPemilihan := dataset.FieldValues["akhir_pemilihan"]

	query := `INSERT INTO spse_pemilihan
		(kode_rup, satuan_kerja, nama_paket, metode_pemilihan, rencana_pemilihan, tanggal_pemilihan, nilai_hps, status_paket,
		 kode_satuan_kerja, cara_pengadaan, jenis_pengadaan, pdn, umk, sumber_dana, kode_rup_lokal,
		 metode_pengadaan, pagu_rup, tipe_swakelola, akhir_pemilihan, created_at, last_update, deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, NOW(), $20, NULL)
		ON CONFLICT (kode_rup, nama_paket) DO UPDATE SET
			satuan_kerja = EXCLUDED.satuan_kerja,
			metode_pemilihan = EXCLUDED.metode_pemilihan,
			rencana_pemilihan = EXCLUDED.rencana_pemilihan,
			tanggal_pemilihan = EXCLUDED.tanggal_pemilihan,
			nilai_hps = EXCLUDED.nilai_hps,
			status_paket = EXCLUDED.status_paket,
			kode_satuan_kerja = EXCLUDED.kode_satuan_kerja,
			cara_pengadaan = EXCLUDED.cara_pengadaan,
			jenis_pengadaan = EXCLUDED.jenis_pengadaan,
			pdn = EXCLUDED.pdn,
			umk = EXCLUDED.umk,
			sumber_dana = EXCLUDED.sumber_dana,
			kode_rup_lokal = EXCLUDED.kode_rup_lokal,
			metode_pengadaan = EXCLUDED.metode_pengadaan,
			pagu_rup = EXCLUDED.pagu_rup,
			tipe_swakelola = EXCLUDED.tipe_swakelola,
			akhir_pemilihan = EXCLUDED.akhir_pemilihan,
			last_update = EXCLUDED.last_update,
			deleted_at = NULL`

	args := []interface{}{kodeRUP, satuanKerja, namaPaket, metodePemilihan, rencanaPemilihan, tanggalPemilihan, nilaiHPS, statusPaket,
		kodeSatuanKerja, caraPengadaan, jenisPengadaan, pdn, umk, sumberDana, kodeRUPLokal,
		metodePengadaan, paguRUP, tipeSwakelola, akhirPemilihan, time.Now().Unix()}

	return query, args
}

// buildHasilPemilihanInsertFromDataset builds INSERT query for hasilpemilihan data using ordered dataset
func (ctrl SPSEController) buildHasilPemilihanInsertFromDataset(dataset *OrderedDataSet) (string, []interface{}) {
	if dataset == nil || dataset.FieldValues == nil {
		return "", nil
	}

	// Extract fields from ordered dataset
	kodeRUP := dataset.FieldValues["kode_rup"]
	satuanKerja := dataset.FieldValues["satuan_kerja"]
	namaPaket := dataset.FieldValues["nama_paket"]
	metodePemilihan := dataset.FieldValues["metode_pemilihan"]
	tanggalHasilPemilihan := dataset.FieldValues["tanggal_hasil_pemilihan"]
	nilaiHasilPemilihan := dataset.FieldValues["nilai_hasil_pemilihan"]
	statusPaket := dataset.FieldValues["status_paket"]
	kodeSatuanKerja := dataset.FieldValues["kode_satuan_kerja"]
	caraPengadaan := dataset.FieldValues["cara_pengadaan"]
	jenisPengadaan := dataset.FieldValues["jenis_pengadaan"]
	pdn := dataset.FieldValues["pdn"]
	umk := dataset.FieldValues["umk"]
	sumberDana := dataset.FieldValues["sumber_dana"]
	kodeRUPLokal := dataset.FieldValues["kode_rup_lokal"]
	metodePengadaan := dataset.FieldValues["metode_pengadaan"]
	paguRUP := dataset.FieldValues["pagu_rup"]
	tipeSwakelola := dataset.FieldValues["tipe_swakelola"]

	query := `INSERT INTO spse_hasilpemilihan
		(kode_rup, satuan_kerja, nama_paket, metode_pemilihan, tanggal_hasil_pemilihan, nilai_hasil_pemilihan, status_paket,
		 kode_satuan_kerja, cara_pengadaan, jenis_pengadaan, pdn, umk, sumber_dana, kode_rup_lokal,
		 metode_pengadaan, pagu_rup, tipe_swakelola, created_at, last_update, deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, NOW(), $18, NULL)
		ON CONFLICT (kode_rup, nama_paket) DO UPDATE SET
			satuan_kerja = EXCLUDED.satuan_kerja,
			metode_pemilihan = EXCLUDED.metode_pemilihan,
			tanggal_hasil_pemilihan = EXCLUDED.tanggal_hasil_pemilihan,
			nilai_hasil_pemilihan = EXCLUDED.nilai_hasil_pemilihan,
			status_paket = EXCLUDED.status_paket,
			kode_satuan_kerja = EXCLUDED.kode_satuan_kerja,
			cara_pengadaan = EXCLUDED.cara_pengadaan,
			jenis_pengadaan = EXCLUDED.jenis_pengadaan,
			pdn = EXCLUDED.pdn,
			umk = EXCLUDED.umk,
			sumber_dana = EXCLUDED.sumber_dana,
			kode_rup_lokal = EXCLUDED.kode_rup_lokal,
			metode_pengadaan = EXCLUDED.metode_pengadaan,
			pagu_rup = EXCLUDED.pagu_rup,
			tipe_swakelola = EXCLUDED.tipe_swakelola,
			last_update = EXCLUDED.last_update,
			deleted_at = NULL`

	args := []interface{}{kodeRUP, satuanKerja, namaPaket, metodePemilihan, tanggalHasilPemilihan, nilaiHasilPemilihan, statusPaket,
		kodeSatuanKerja, caraPengadaan, jenisPengadaan, pdn, umk, sumberDana, kodeRUPLokal,
		metodePengadaan, paguRUP, tipeSwakelola, time.Now().Unix()}

	return query, args
}

// buildKontrakInsertFromDataset builds INSERT query for kontrak data using ordered dataset
func (ctrl SPSEController) buildKontrakInsertFromDataset(dataset *OrderedDataSet) (string, []interface{}) {
	if dataset == nil || dataset.FieldValues == nil {
		return "", nil
	}

	// Extract fields from ordered dataset
	kodeRUP := dataset.FieldValues["kode_rup"]
	satuanKerja := dataset.FieldValues["satuan_kerja"]
	namaPaket := dataset.FieldValues["nama_paket"]
	metodePemilihan := dataset.FieldValues["metode_pemilihan"]
	tanggalKontrak := dataset.FieldValues["tanggal_kontrak"]
	nilaiKontrak := dataset.FieldValues["nilai_kontrak"]
	statusPaket := dataset.FieldValues["status_paket"]
	mulaiKontrak := dataset.FieldValues["mulai_kontrak"]
	nilaiBAP := dataset.FieldValues["nilai_bap"]
	selesaiKontrak := dataset.FieldValues["selesai_kontrak"]
	kodeSatuanKerja := dataset.FieldValues["kode_satuan_kerja"]
	caraPengadaan := dataset.FieldValues["cara_pengadaan"]
	jenisPengadaan := dataset.FieldValues["jenis_pengadaan"]
	pdn := dataset.FieldValues["pdn"]
	umk := dataset.FieldValues["umk"]
	sumberDana := dataset.FieldValues["sumber_dana"]
	kodeRUPLokal := dataset.FieldValues["kode_rup_lokal"]
	metodePengadaan := dataset.FieldValues["metode_pengadaan"]
	tipeSwakelola := dataset.FieldValues["tipe_swakelola"]

	query := `INSERT INTO spse_kontrak
		(kode_rup, satuan_kerja, nama_paket, metode_pemilihan, tanggal_kontrak, nilai_kontrak, status_paket,
		 mulai_kontrak, nilai_bap, selesai_kontrak, kode_satuan_kerja, cara_pengadaan, jenis_pengadaan,
		 pdn, umk, sumber_dana, kode_rup_lokal, metode_pengadaan, tipe_swakelola, created_at, last_update, deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, NOW(), $20, NULL)
		ON CONFLICT (kode_rup, nama_paket) DO UPDATE SET
			satuan_kerja = EXCLUDED.satuan_kerja,
			metode_pemilihan = EXCLUDED.metode_pemilihan,
			tanggal_kontrak = EXCLUDED.tanggal_kontrak,
			nilai_kontrak = EXCLUDED.nilai_kontrak,
			status_paket = EXCLUDED.status_paket,
			mulai_kontrak = EXCLUDED.mulai_kontrak,
			nilai_bap = EXCLUDED.nilai_bap,
			selesai_kontrak = EXCLUDED.selesai_kontrak,
			kode_satuan_kerja = EXCLUDED.kode_satuan_kerja,
			cara_pengadaan = EXCLUDED.cara_pengadaan,
			jenis_pengadaan = EXCLUDED.jenis_pengadaan,
			pdn = EXCLUDED.pdn,
			umk = EXCLUDED.umk,
			sumber_dana = EXCLUDED.sumber_dana,
			kode_rup_lokal = EXCLUDED.kode_rup_lokal,
			metode_pengadaan = EXCLUDED.metode_pengadaan,
			tipe_swakelola = EXCLUDED.tipe_swakelola,
			last_update = EXCLUDED.last_update,
			deleted_at = NULL`

	args := []interface{}{kodeRUP, satuanKerja, namaPaket, metodePemilihan, tanggalKontrak, nilaiKontrak, statusPaket,
		mulaiKontrak, nilaiBAP, selesaiKontrak, kodeSatuanKerja, caraPengadaan, jenisPengadaan,
		pdn, umk, sumberDana, kodeRUPLokal, metodePengadaan, tipeSwakelola, time.Now().Unix()}

	return query, args
}

// buildSerahTerimaInsertFromDataset builds INSERT query for serahterima data using ordered dataset
func (ctrl SPSEController) buildSerahTerimaInsertFromDataset(dataset *OrderedDataSet) (string, []interface{}) {
	if dataset == nil || dataset.FieldValues == nil {
		return "", nil
	}

	// Extract fields from ordered dataset
	kodeRUP := dataset.FieldValues["kode_rup"]
	satuanKerja := dataset.FieldValues["satuan_kerja"]
	namaPaket := dataset.FieldValues["nama_paket"]
	metodePemilihan := dataset.FieldValues["metode_pemilihan"]
	tanggalSerahTerima := dataset.FieldValues["tanggal_serah_terima"]
	nilaiBAP := dataset.FieldValues["nilai_bap"]
	statusPaket := dataset.FieldValues["status_paket"]
	kodeSatuanKerja := dataset.FieldValues["kode_satuan_kerja"]
	caraPengadaan := dataset.FieldValues["cara_pengadaan"]
	jenisPengadaan := dataset.FieldValues["jenis_pengadaan"]
	pdn := dataset.FieldValues["pdn"]
	umk := dataset.FieldValues["umk"]
	sumberDana := dataset.FieldValues["sumber_dana"]
	kodeRUPLokal := dataset.FieldValues["kode_rup_lokal"]
	metodePengadaan := dataset.FieldValues["metode_pengadaan"]
	tipeSwakelola := dataset.FieldValues["tipe_swakelola"]

	query := `INSERT INTO spse_serahterima
		(kode_rup, satuan_kerja, nama_paket, metode_pemilihan, tanggal_serah_terima, nilai_bap, status_paket,
		 kode_satuan_kerja, cara_pengadaan, jenis_pengadaan, pdn, umk, sumber_dana, kode_rup_lokal,
		 metode_pengadaan, tipe_swakelola, created_at, last_update, deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, NOW(), $17, NULL)
		ON CONFLICT (kode_rup, nama_paket) DO UPDATE SET
			satuan_kerja = EXCLUDED.satuan_kerja,
			metode_pemilihan = EXCLUDED.metode_pemilihan,
			tanggal_serah_terima = EXCLUDED.tanggal_serah_terima,
			nilai_bap = EXCLUDED.nilai_bap,
			status_paket = EXCLUDED.status_paket,
			kode_satuan_kerja = EXCLUDED.kode_satuan_kerja,
			cara_pengadaan = EXCLUDED.cara_pengadaan,
			jenis_pengadaan = EXCLUDED.jenis_pengadaan,
			pdn = EXCLUDED.pdn,
			umk = EXCLUDED.umk,
			sumber_dana = EXCLUDED.sumber_dana,
			kode_rup_lokal = EXCLUDED.kode_rup_lokal,
			metode_pengadaan = EXCLUDED.metode_pengadaan,
			tipe_swakelola = EXCLUDED.tipe_swakelola,
			last_update = EXCLUDED.last_update,
			deleted_at = NULL`

	args := []interface{}{kodeRUP, satuanKerja, namaPaket, metodePemilihan, tanggalSerahTerima, nilaiBAP, statusPaket,
		kodeSatuanKerja, caraPengadaan, jenisPengadaan, pdn, umk, sumberDana, kodeRUPLokal,
		metodePengadaan, tipeSwakelola, time.Now().Unix()}

	return query, args
}
