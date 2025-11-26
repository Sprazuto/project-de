package controllers

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Massad/gin-boilerplate/db"

	"github.com/gin-gonic/gin"
)

// PerencanaanData represents planning stage data (matches SPSEPerencanaan model)
type PerencanaanData struct {
	ID                int64     `json:"id"`
	KodeRUP           string    `json:"kode_rup"`
	SatuanKerja       string    `json:"satuan_kerja"`
	NamaPaket         string    `json:"nama_paket"`
	MetodePemilihan   string    `json:"metode_pemilihan"`
	TanggalPengumuman string    `json:"tanggal_pengumuman"`
	RencanaPemilihan  string    `json:"rencana_pemilihan"`
	PaguRUP           string    `json:"pagu_rup"`
	KodeSatuanKerja   string    `json:"kode_satuan_kerja"`
	CaraPengadaan     string    `json:"cara_pengadaan"`
	JenisPengadaan    string    `json:"jenis_pengadaan"`
	PDN               string    `json:"pdn"`
	UMK               string    `json:"umk"`
	SumberDana        string    `json:"sumber_dana"`
	KodeRUPLokal      string    `json:"kode_rup_lokal"`
	AkhirPemilihan    string    `json:"akhir_pemilihan"`
	TipeSwakelola     string    `json:"tipe_swakelola"`
	CreatedAt         time.Time `json:"created_at"`
	LastUpdate        int64     `json:"last_update"`
}

// PersiapanData represents preparation stage data (matches SPSEPersiapan model)
type PersiapanData struct {
	ID               int64     `json:"id"`
	KodeRUP          string    `json:"kode_rup"`
	SatuanKerja      string    `json:"satuan_kerja"`
	NamaPaket        string    `json:"nama_paket"`
	MetodePemilihan  string    `json:"metode_pemilihan"`
	TanggalBuatPaket string    `json:"tanggal_buat_paket"`
	NilaiPaguRUP     string    `json:"nilai_pagu_rup"`
	NilaiPaguPaket   string    `json:"nilai_pagu_paket"`
	KodeSatuanKerja  string    `json:"kode_satuan_kerja"`
	CaraPengadaan    string    `json:"cara_pengadaan"`
	JenisPengadaan   string    `json:"jenis_pengadaan"`
	PDN              string    `json:"pdn"`
	UMK              string    `json:"umk"`
	SumberDana       string    `json:"sumber_dana"`
	KodeRUPLokal     string    `json:"kode_rup_lokal"`
	MetodePengadaan  string    `json:"metode_pengadaan"`
	TipeSwakelola    string    `json:"tipe_swakelola"`
	CreatedAt        time.Time `json:"created_at"`
	LastUpdate       int64     `json:"last_update"`
}

// PemilihanData represents selection/contract stage data (matches SPSEPemilihan model)
type PemilihanData struct {
	ID               int64     `json:"id"`
	KodeRUP          string    `json:"kode_rup"`
	SatuanKerja      string    `json:"satuan_kerja"`
	NamaPaket        string    `json:"nama_paket"`
	MetodePemilihan  string    `json:"metode_pemilihan"`
	RencanaPemilihan string    `json:"rencana_pemilihan"`
	TanggalPemilihan string    `json:"tanggal_pemilihan"`
	NilaiHPS         string    `json:"nilai_hps"`
	StatusPaket      string    `json:"status_paket"`
	KodeSatuanKerja  string    `json:"kode_satuan_kerja"`
	CaraPengadaan    string    `json:"cara_pengadaan"`
	JenisPengadaan   string    `json:"jenis_pengadaan"`
	PDN              string    `json:"pdn"`
	UMK              string    `json:"umk"`
	SumberDana       string    `json:"sumber_dana"`
	KodeRUPLokal     string    `json:"kode_rup_lokal"`
	MetodePengadaan  string    `json:"metode_pengadaan"`
	PaguRUP          string    `json:"pagu_rup"`
	TipeSwakelola    string    `json:"tipe_swakelola"`
	AkhirPemilihan   string    `json:"akhir_pemilihan"`
	CreatedAt        time.Time `json:"created_at"`
	LastUpdate       int64     `json:"last_update"`
}

// ScrapingResult represents the result of scraping operations
type ScrapingResult struct {
	Success      bool   `json:"success"`
	Message      string `json:"message"`
	RecordsFound int    `json:"records_found"`
	Endpoint     string `json:"endpoint"`
}

// SPSEController handles SPSE perencanaan data scraping
type SPSEController struct{}

// Global cookie jar for maintaining session
var spseJar, _ = cookiejar.New(nil)

var spseClient = &http.Client{
	Transport: &http.Transport{
		MaxIdleConns:          10,
		IdleConnTimeout:       30 * time.Second,
		DisableCompression:    false,
		DisableKeepAlives:     false,
		ResponseHeaderTimeout: 30 * time.Second,
	},
	Timeout: 30 * time.Second,
	Jar:     spseJar,
}

// getAuthenticityToken retrieves the authenticity token from the main page
func (ctrl SPSEController) getAuthenticityToken() (string, error) {
	targetURL := "https://spse.inaproc.id/sumedangkab/amel/perencanaan/detail"

	log.Printf("Fetching authenticity token from: %s", targetURL)

	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")

	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	resp, err := spseClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to get page: %v", err)
	}
	defer resp.Body.Close()

	var body []byte
	var readErr error

	// Handle gzip compressed response
	if resp.Header.Get("Content-Encoding") == "gzip" {
		reader, gzipErr := gzip.NewReader(resp.Body)
		if gzipErr != nil {
			return "", fmt.Errorf("failed to create gzip reader: %v", gzipErr)
		}
		defer reader.Close()
		body, readErr = io.ReadAll(reader)
	} else {
		body, readErr = io.ReadAll(resp.Body)
	}

	if readErr != nil {
		return "", fmt.Errorf("failed to read response body: %v", readErr)
	}

	content := string(body)

	// Try to find token using the exact pattern we found
	re := regexp.MustCompile(`name=["']authenticityToken["'][^>]*value=["']([^"']+)["']`)
	matches := re.FindStringSubmatch(content)
	if len(matches) > 1 {
		token := matches[1]
		return token, nil
	}

	// Also try the JavaScript pattern we found
	re = regexp.MustCompile(`d\.authenticityToken\s*=\s*['"]([^'"]+)['"]`)
	matches = re.FindStringSubmatch(content)
	if len(matches) > 1 {
		token := matches[1]
		return token, nil
	}

	// Try alternative patterns
	patterns := []string{
		`<input[^>]*name=["']authenticityToken["'][^>]*value=["']([^"']+)["']`,
		`<input[^>]*value=["']([^"']+)["'][^>]*name=["']authenticityToken["']`,
		`authenticityToken["']?\s*[:=]\s*["']([^"']+)["']`,
		`csrf[_-]?token["']?\s*[:=]\s*["']([^"']+)["']`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(content)
		if len(matches) > 1 {
			token := matches[1]
			return token, nil
		}
	}

	return "", fmt.Errorf("authenticity token not found")
}

// Helper function to get minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// scrapeEndpoint scrapes a specific endpoint
func (ctrl SPSEController) scrapeEndpoint(endpoint string, tableName string) (ScrapingResult, error) {
	// Get authenticity token
	authToken, err := ctrl.getAuthenticityToken()
	if err != nil {
		return ScrapingResult{Success: false, Message: fmt.Sprintf("Failed to get auth token: %v", err)}, err
	}

	// Prepare request data
	formData := url.Values{}
	formData.Set("authenticityToken", authToken)
	formData.Set("activeSatker", "")
	formData.Set("activeInstansi", "D118")
	formData.Set("activeYear", "2025")

	// Prepare request
	req, err := http.NewRequest("POST", "https://spse.inaproc.id"+endpoint, strings.NewReader(formData.Encode()))
	if err != nil {
		return ScrapingResult{Success: false, Message: fmt.Sprintf("Failed to create request: %v", err)}, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Referer", "https://spse.inaproc.id/sumedangkab/amel")

	// Execute request
	resp, err := spseClient.Do(req)
	if err != nil {
		return ScrapingResult{Success: false, Message: fmt.Sprintf("Failed to execute request: %v", err)}, err
	}
	defer resp.Body.Close()

	var body []byte
	var readErr error

	// Handle gzip compressed response for API calls too
	if resp.Header.Get("Content-Encoding") == "gzip" {
		reader, gzipErr := gzip.NewReader(resp.Body)
		if gzipErr != nil {
			return ScrapingResult{Success: false, Message: fmt.Sprintf("Failed to create gzip reader: %v", gzipErr)}, gzipErr
		}
		defer reader.Close()
		body, readErr = io.ReadAll(reader)
	} else {
		body, readErr = io.ReadAll(resp.Body)
	}

	if readErr != nil {
		return ScrapingResult{Success: false, Message: fmt.Sprintf("Failed to read response: %v", readErr)}, readErr
	}

	// Try to parse as JSON first
	var responseData map[string]interface{}
	if err := json.Unmarshal(body, &responseData); err != nil {
		return ScrapingResult{
			Success:      false,
			Message:      fmt.Sprintf("Failed to parse response as JSON: %v", err),
			RecordsFound: 0,
			Endpoint:     endpoint,
		}, fmt.Errorf("failed to parse JSON response")
	}

	// Check if we have data
	dataValue, hasData := responseData["data"]
	if !hasData {
		return ScrapingResult{
			Success:      true,
			Message:      "No data field found in response",
			RecordsFound: 0,
			Endpoint:     endpoint,
		}, nil
	}

	var dataArray []interface{}
	var ok bool

	// Try different possible data structures
	switch v := dataValue.(type) {
	case []interface{}:
		dataArray = v
		ok = true
	case []map[string]interface{}:
		// Convert []map[string]interface{} to []interface{}
		for _, item := range v {
			dataArray = append(dataArray, item)
		}
		ok = true
	case map[string]interface{}:
		// The data might be wrapped in another object
		for _, val := range v {
			if arr, isArray := val.([]interface{}); isArray {
				dataArray = arr
				ok = true
				break
			}
		}
		if !ok {
			// If no array found, treat the map as a single record
			dataArray = []interface{}{v}
			ok = true
		}
	default:
		ok = false
	}

	if !ok {
		return ScrapingResult{
			Success:      false,
			Message:      fmt.Sprintf("Data field has unexpected format: %T", dataValue),
			RecordsFound: 0,
			Endpoint:     endpoint,
		}, fmt.Errorf("unexpected data format")
	}

	// Store data in database
	database := db.GetDB()
	if database == nil {
		log.Printf("ERROR: Database connection is nil!")
		return ScrapingResult{
			Success:      false,
			Message:      "Database connection is nil",
			RecordsFound: 0,
			Endpoint:     endpoint,
		}, fmt.Errorf("database connection is nil")
	}

	// Test database connection
	err = database.Db.Ping()
	if err != nil {
		log.Printf("ERROR: Cannot ping database: %v", err)
		return ScrapingResult{
			Success:      false,
			Message:      fmt.Sprintf("Database ping failed: %v", err),
			RecordsFound: 0,
			Endpoint:     endpoint,
		}, fmt.Errorf("database ping failed: %v", err)
	}

	// Check if table exists
	var tableExists int
	tableCheckQuery := fmt.Sprintf("SELECT 1 FROM information_schema.tables WHERE table_name = 'spse_%s'", tableName)
	err = database.QueryRow(tableCheckQuery).Scan(&tableExists)
	if err != nil {
		log.Printf("ERROR: Cannot check table existence for spse_%s: %v", tableName, err)
		return ScrapingResult{
			Success:      false,
			Message:      fmt.Sprintf("Table spse_%s does not exist or cannot be accessed: %v", tableName, err),
			RecordsFound: 0,
			Endpoint:     endpoint,
		}, fmt.Errorf("table check failed")
	}

	if tableExists != 1 {
		log.Printf("ERROR: Table spse_%s does not exist in database", tableName)
		return ScrapingResult{
			Success:      false,
			Message:      fmt.Sprintf("Table spse_%s does not exist", tableName),
			RecordsFound: 0,
			Endpoint:     endpoint,
		}, fmt.Errorf("table does not exist")
	}

	log.Printf("Database connection OK, table spse_%s exists", tableName)
	recordsStored := 0
	recordsFailed := 0

	for _, item := range dataArray {
		var itemMap map[string]interface{}

		// Handle different data formats
		switch v := item.(type) {
		case map[string]interface{}:
			// Already a map, use directly
			itemMap = v
		case []interface{}:
			// Array format - convert to map based on table-specific field mapping
			itemMap = ctrl.mapArrayToFields(tableName, v)
		case string:
			// String format - try to parse as delimited data
			itemMap = ctrl.parseStringData(v)
		default:
			recordsFailed++
			continue
		}

		// Store based on table type
		var insertQuery string
		switch tableName {
		case "perencanaan":
			insertQuery = ctrl.buildPerencanaanInsert(itemMap)
		case "persiapan":
			insertQuery = ctrl.buildPersiapanInsert(itemMap)
		case "pemilihan":
			insertQuery = ctrl.buildPemilihanInsert(itemMap)
		}

		if insertQuery != "" {
			_, err := database.Exec(insertQuery)
			if err != nil {
				recordsFailed++
				// Continue with next record instead of failing completely
			} else {
				recordsStored++
			}
		} else {
			recordsFailed++
		}
	}

	log.Printf("Stored %d records out of %d fetched (Failed: %d)", recordsStored, len(dataArray), recordsFailed)

	return ScrapingResult{
		Success:      true,
		Message:      fmt.Sprintf("Successfully stored %d records", recordsStored),
		RecordsFound: recordsStored,
		Endpoint:     endpoint,
	}, nil
}

// buildPerencanaanInsert builds INSERT query for perencanaan data
func (ctrl SPSEController) buildPerencanaanInsert(data map[string]interface{}) string {
	// Extract fields with proper escaping
	kodeRUP := ctrl.escapeString(data["kode_rup"])
	satuanKerja := ctrl.escapeString(data["satuan_kerja"])
	namaPaket := ctrl.escapeString(data["nama_paket"])
	metodePemilihan := ctrl.escapeString(data["metode_pemilihan"])
	tanggalPengumuman := ctrl.escapeString(data["tanggal_pengumuman"])
	rencanaPemilihan := ctrl.escapeString(data["rencana_pemilihan"])
	paguRUP := ctrl.escapeString(data["pagu_rup"])
	kodeSatuanKerja := ctrl.escapeString(data["kode_satuan_kerja"])
	caraPengadaan := ctrl.escapeString(data["cara_pengadaan"])
	jenisPengadaan := ctrl.escapeString(data["jenis_pengadaan"])
	pdn := ctrl.escapeString(data["pdn"])
	umk := ctrl.escapeString(data["umk"])
	sumberDana := ctrl.escapeString(data["sumber_dana"])
	kodeRUPLokal := ctrl.escapeString(data["kode_rup_lokal"])
	akhirPemilihan := ctrl.escapeString(data["akhir_pemilihan"])
	tipeSwakelola := ctrl.escapeString(data["tipe_swakelola"])

	return fmt.Sprintf(`INSERT INTO spse_perencanaan
		(kode_rup, satuan_kerja, nama_paket, metode_pemilihan, tanggal_pengumuman, rencana_pemilihan, pagu_rup,
		 kode_satuan_kerja, cara_pengadaan, jenis_pengadaan, pdn, umk, sumber_dana, kode_rup_lokal,
		 akhir_pemilihan, tipe_swakelola, created_at, last_update)
		VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', NOW(), %d)
		ON CONFLICT (kode_rup, nama_paket) DO NOTHING`,
		kodeRUP, satuanKerja, namaPaket, metodePemilihan, tanggalPengumuman, rencanaPemilihan, paguRUP,
		kodeSatuanKerja, caraPengadaan, jenisPengadaan, pdn, umk, sumberDana, kodeRUPLokal,
		akhirPemilihan, tipeSwakelola, time.Now().Unix())
}

// buildPersiapanInsert builds INSERT query for persiapan data
func (ctrl SPSEController) buildPersiapanInsert(data map[string]interface{}) string {
	// Extract fields with proper escaping
	kodeRUP := ctrl.escapeString(data["kode_rup"])
	satuanKerja := ctrl.escapeString(data["satuan_kerja"])
	namaPaket := ctrl.escapeString(data["nama_paket"])
	metodePemilihan := ctrl.escapeString(data["metode_pemilihan"])
	tanggalBuatPaket := ctrl.escapeString(data["tanggal_buat_paket"])
	nilaiPaguRUP := ctrl.escapeString(data["nilai_pagu_rup"])
	nilaiPaguPaket := ctrl.escapeString(data["nilai_pagu_paket"])
	kodeSatuanKerja := ctrl.escapeString(data["kode_satuan_kerja"])
	caraPengadaan := ctrl.escapeString(data["cara_pengadaan"])
	jenisPengadaan := ctrl.escapeString(data["jenis_pengadaan"])
	pdn := ctrl.escapeString(data["pdn"])
	umk := ctrl.escapeString(data["umk"])
	sumberDana := ctrl.escapeString(data["sumber_dana"])
	kodeRUPLokal := ctrl.escapeString(data["kode_rup_lokal"])
	metodePengadaan := ctrl.escapeString(data["metode_pengadaan"])
	tipeSwakelola := ctrl.escapeString(data["tipe_swakelola"])

	return fmt.Sprintf(`INSERT INTO spse_persiapan
		(kode_rup, satuan_kerja, nama_paket, metode_pemilihan, tanggal_buat_paket, nilai_pagu_rup, nilai_pagu_paket,
		 kode_satuan_kerja, cara_pengadaan, jenis_pengadaan, pdn, umk, sumber_dana, kode_rup_lokal,
		 metode_pengadaan, tipe_swakelola, created_at, last_update)
		VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', NOW(), %d)
		ON CONFLICT (kode_rup, nama_paket) DO NOTHING`,
		kodeRUP, satuanKerja, namaPaket, metodePemilihan, tanggalBuatPaket, nilaiPaguRUP, nilaiPaguPaket,
		kodeSatuanKerja, caraPengadaan, jenisPengadaan, pdn, umk, sumberDana, kodeRUPLokal,
		metodePengadaan, tipeSwakelola, time.Now().Unix())
}

// buildPemilihanInsert builds INSERT query for pemilihan data
func (ctrl SPSEController) buildPemilihanInsert(data map[string]interface{}) string {
	// Extract fields with proper escaping
	kodeRUP := ctrl.escapeString(data["kode_rup"])
	satuanKerja := ctrl.escapeString(data["satuan_kerja"])
	namaPaket := ctrl.escapeString(data["nama_paket"])
	metodePemilihan := ctrl.escapeString(data["metode_pemilihan"])
	rencanaPemilihan := ctrl.escapeString(data["rencana_pemilihan"])
	tanggalPemilihan := ctrl.escapeString(data["tanggal_pemilihan"])
	nilaiHPS := ctrl.escapeString(data["nilai_hps"])
	statusPaket := ctrl.escapeString(data["status_paket"])
	kodeSatuanKerja := ctrl.escapeString(data["kode_satuan_kerja"])
	caraPengadaan := ctrl.escapeString(data["cara_pengadaan"])
	jenisPengadaan := ctrl.escapeString(data["jenis_pengadaan"])
	pdn := ctrl.escapeString(data["pdn"])
	umk := ctrl.escapeString(data["umk"])
	sumberDana := ctrl.escapeString(data["sumber_dana"])
	kodeRUPLokal := ctrl.escapeString(data["kode_rup_lokal"])
	metodePengadaan := ctrl.escapeString(data["metode_pengadaan"])
	paguRUP := ctrl.escapeString(data["pagu_rup"])
	tipeSwakelola := ctrl.escapeString(data["tipe_swakelola"])
	akhirPemilihan := ctrl.escapeString(data["akhir_pemilihan"])

	return fmt.Sprintf(`INSERT INTO spse_pemilihan
		(kode_rup, satuan_kerja, nama_paket, metode_pemilihan, rencana_pemilihan, tanggal_pemilihan, nilai_hps, status_paket,
		 kode_satuan_kerja, cara_pengadaan, jenis_pengadaan, pdn, umk, sumber_dana, kode_rup_lokal,
		 metode_pengadaan, pagu_rup, tipe_swakelola, akhir_pemilihan, created_at, last_update)
		VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', NOW(), %d)
		ON CONFLICT (kode_rup, nama_paket) DO NOTHING`,
		kodeRUP, satuanKerja, namaPaket, metodePemilihan, rencanaPemilihan, tanggalPemilihan, nilaiHPS, statusPaket,
		kodeSatuanKerja, caraPengadaan, jenisPengadaan, pdn, umk, sumberDana, kodeRUPLokal,
		metodePengadaan, paguRUP, tipeSwakelola, akhirPemilihan, time.Now().Unix())
}

// escapeString escapes single quotes for SQL
func (ctrl SPSEController) escapeString(value interface{}) string {
	if value == nil {
		return ""
	}
	str := fmt.Sprintf("%v", value)
	// Escape single quotes by doubling them
	str = strings.ReplaceAll(str, "'", "''")
	return str
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

// mapArrayToFields maps array data to field names using smart field detection and validation
func (ctrl SPSEController) mapArrayToFields(tableName string, arr []interface{}) map[string]interface{} {
	itemMap := make(map[string]interface{})

	// Convert array to strings for easier processing
	arrStr := make([]string, len(arr))
	for i, v := range arr {
		if v != nil {
			arrStr[i] = fmt.Sprintf("%v", v)
		} else {
			arrStr[i] = ""
		}
	}

	// Apply table-specific smart mapping logic
	switch tableName {
	case "perencanaan":
		itemMap = ctrl.mapPerencanaanFields(arrStr)
	case "persiapan":
		itemMap = ctrl.mapPersiapanFields(arrStr)
	case "pemilihan":
		itemMap = ctrl.mapPemilihanFields(arrStr)
	default:
		// Fallback: try to map first few elements to common fields
		if len(arrStr) >= 1 {
			itemMap["kode_rup"] = arrStr[0]
		}
		if len(arrStr) >= 2 {
			itemMap["nama_paket"] = arrStr[1]
		}
		if len(arrStr) >= 3 {
			itemMap["satuan_kerja"] = arrStr[2]
		}
	}

	// Validate and fill missing fields with appropriate defaults
	itemMap = ctrl.validateAndFillFields(tableName, itemMap)

	return itemMap
}

// mapPerencanaanFields maps fields for perencanaan table with smart detection
func (ctrl SPSEController) mapPerencanaanFields(arr []string) map[string]interface{} {
	itemMap := make(map[string]interface{})

	// Perencanaan specific mapping - more flexible approach
	// Try to identify fields by content patterns rather than fixed positions

	for i, val := range arr {
		trimmedVal := strings.TrimSpace(val)
		if trimmedVal == "" {
			continue
		}

		// Pattern-based field detection
		switch {
		case i == 0 || ctrl.isRupCode(trimmedVal):
			if itemMap["kode_rup"] == "" {
				itemMap["kode_rup"] = trimmedVal
			}
		case i == 1 || ctrl.isSatuanKerja(trimmedVal):
			if itemMap["satuan_kerja"] == "" {
				itemMap["satuan_kerja"] = trimmedVal
			}
		case i == 2 || ctrl.isPackageName(trimmedVal):
			if itemMap["nama_paket"] == "" {
				itemMap["nama_paket"] = trimmedVal
			}
		case ctrl.isDate(trimmedVal) && itemMap["tanggal_pengumuman"] == "":
			itemMap["tanggal_pengumuman"] = trimmedVal
		case ctrl.isCurrency(trimmedVal) && itemMap["pagu_rup"] == "":
			itemMap["pagu_rup"] = trimmedVal
		case ctrl.isSelectionMethod(trimmedVal) && itemMap["metode_pemilihan"] == "":
			itemMap["metode_pemilihan"] = trimmedVal
		default:
			// For unidentified fields, assign to next available slot
			ctrl.assignToNextAvailableField(itemMap, trimmedVal, "perencanaan")
		}
	}

	// Fill remaining fields from array positions (backup method)
	ctrl.assignRemainingPerencanaanFields(arr, itemMap)

	return itemMap
}

// mapPersiapanFields maps fields for persiapan table with direct positional mapping
func (ctrl SPSEController) mapPersiapanFields(arr []string) map[string]interface{} {
	itemMap := make(map[string]interface{})

	// Direct positional mapping based on the observed array structure
	// Array structure: [0]=kode_rup, [1]=satuan_kerja, [2]=nama_paket, [3]=metode_pemilihan, etc.
	if len(arr) >= 1 && arr[0] != "" {
		itemMap["kode_rup"] = strings.TrimSpace(arr[0])
	}
	if len(arr) >= 2 && arr[1] != "" {
		itemMap["satuan_kerja"] = strings.TrimSpace(arr[1])
	}
	if len(arr) >= 3 && arr[2] != "" {
		itemMap["nama_paket"] = strings.TrimSpace(arr[2])
	}
	if len(arr) >= 4 && arr[3] != "" {
		itemMap["metode_pemilihan"] = strings.TrimSpace(arr[3])
	}
	if len(arr) >= 5 && arr[4] != "" {
		itemMap["rencana_pemilihan"] = strings.TrimSpace(arr[4])
	}
	if len(arr) >= 6 && arr[5] != "" {
		itemMap["tanggal_pemilihan"] = strings.TrimSpace(arr[5])
	}
	if len(arr) >= 7 && arr[6] != "" {
		itemMap["nilai_hps"] = strings.TrimSpace(arr[6])
	}
	if len(arr) >= 8 && arr[7] != "" {
		itemMap["status_paket"] = strings.TrimSpace(arr[7])
	}
	if len(arr) >= 9 && arr[8] != "" {
		itemMap["kode_satuan_kerja"] = strings.TrimSpace(arr[8])
	}
	if len(arr) >= 10 && arr[9] != "" {
		itemMap["jenis_pengadaan"] = strings.TrimSpace(arr[9])
	}
	if len(arr) >= 11 && arr[10] != "" {
		itemMap["cara_pengadaan"] = strings.TrimSpace(arr[10])
	}
	if len(arr) >= 12 && arr[11] != "" {
		itemMap["pdn"] = strings.TrimSpace(arr[11])
	}
	if len(arr) >= 13 && arr[12] != "" {
		itemMap["umk"] = strings.TrimSpace(arr[12])
	}
	if len(arr) >= 14 && arr[13] != "" {
		itemMap["sumber_dana"] = strings.TrimSpace(arr[13])
	}
	if len(arr) >= 15 && arr[14] != "" {
		itemMap["kode_rup_lokal"] = strings.TrimSpace(arr[14])
	}
	if len(arr) >= 16 && arr[15] != "" {
		itemMap["tipe_swakelola"] = strings.TrimSpace(arr[15])
	}
	if len(arr) >= 17 && arr[16] != "" {
		itemMap["pagu_rup"] = strings.TrimSpace(arr[16])
	}
	if len(arr) >= 19 && arr[18] != "" {
		itemMap["akhir_pemilihan"] = strings.TrimSpace(arr[18])
	}

	// Fill any missing fields using the backup positional method
	ctrl.assignRemainingPersiapanFields(arr, itemMap)

	return itemMap
}

// mapPemilihanFields maps fields for pemilihan table with direct positional mapping
func (ctrl SPSEController) mapPemilihanFields(arr []string) map[string]interface{} {
	itemMap := make(map[string]interface{})

	// Direct positional mapping based on the observed array structure
	// Array structure: [0]=kode_rup, [1]=satuan_kerja, [2]=nama_paket, [3]=metode_pemilihan, [4]=rencana_pemilihan, [5]=tanggal_pemilihan, [6]=nilai_hps, [7]=status_paket, etc.
	if len(arr) >= 1 && arr[0] != "" {
		itemMap["kode_rup"] = strings.TrimSpace(arr[0])
	}
	if len(arr) >= 2 && arr[1] != "" {
		itemMap["satuan_kerja"] = strings.TrimSpace(arr[1])
	}
	if len(arr) >= 3 && arr[2] != "" {
		itemMap["nama_paket"] = strings.TrimSpace(arr[2])
	}
	if len(arr) >= 4 && arr[3] != "" {
		itemMap["metode_pemilihan"] = strings.TrimSpace(arr[3])
	}
	if len(arr) >= 5 && arr[4] != "" {
		itemMap["rencana_pemilihan"] = strings.TrimSpace(arr[4])
	}
	if len(arr) >= 6 && arr[5] != "" {
		itemMap["tanggal_pemilihan"] = strings.TrimSpace(arr[5])
	}
	if len(arr) >= 7 && arr[6] != "" {
		itemMap["nilai_hps"] = strings.TrimSpace(arr[6])
	}
	if len(arr) >= 8 && arr[7] != "" {
		itemMap["status_paket"] = strings.TrimSpace(arr[7])
	}
	if len(arr) >= 9 && arr[8] != "" {
		itemMap["kode_satuan_kerja"] = strings.TrimSpace(arr[8])
	}
	if len(arr) >= 10 && arr[9] != "" {
		itemMap["jenis_pengadaan"] = strings.TrimSpace(arr[9])
	}
	if len(arr) >= 11 && arr[10] != "" {
		itemMap["cara_pengadaan"] = strings.TrimSpace(arr[10])
	}
	if len(arr) >= 12 && arr[11] != "" {
		itemMap["pdn"] = strings.TrimSpace(arr[11])
	}
	if len(arr) >= 13 && arr[12] != "" {
		itemMap["umk"] = strings.TrimSpace(arr[12])
	}
	if len(arr) >= 14 && arr[13] != "" {
		itemMap["sumber_dana"] = strings.TrimSpace(arr[13])
	}
	if len(arr) >= 15 && arr[14] != "" {
		itemMap["kode_rup_lokal"] = strings.TrimSpace(arr[14])
	}
	if len(arr) >= 16 && arr[15] != "" {
		itemMap["tipe_swakelola"] = strings.TrimSpace(arr[15])
	}
	if len(arr) >= 17 && arr[16] != "" {
		itemMap["pagu_rup"] = strings.TrimSpace(arr[16])
	}
	if len(arr) >= 19 && arr[18] != "" {
		itemMap["akhir_pemilihan"] = strings.TrimSpace(arr[18])
	}

	// Fill any missing fields using the backup positional method
	ctrl.assignRemainingPemilihanFields(arr, itemMap)

	return itemMap
}

// assignRemainingPerencanaanFields fills remaining fields using positional mapping as backup
func (ctrl SPSEController) assignRemainingPerencanaanFields(arr []string, itemMap map[string]interface{}) {
	// Standard Perencanaan field order: kode_rup, satuan_kerja, nama_paket, metode_pemilihan, tanggal_pengumuman, rencana_pemilihan, pagu_rup, kode_satuan_kerja, cara_pengadaan, jenis_pengadaan, pdn, umk, sumber_dana, kode_rup_lokal, akhir_pemilihan, tipe_swakelola
	fields := []string{"kode_rup", "satuan_kerja", "nama_paket", "metode_pemilihan", "tanggal_pengumuman", "rencana_pemilihan", "pagu_rup", "kode_satuan_kerja", "cara_pengadaan", "jenis_pengadaan", "pdn", "umk", "sumber_dana", "kode_rup_lokal", "akhir_pemilihan", "tipe_swakelola"}

	for i, field := range fields {
		if itemMap[field] == "" && i < len(arr) {
			itemMap[field] = strings.TrimSpace(arr[i])
		}
	}
}

// assignRemainingPersiapanFields fills remaining fields using positional mapping as backup
func (ctrl SPSEController) assignRemainingPersiapanFields(arr []string, itemMap map[string]interface{}) {
	// Standard Persiapan field order: kode_rup, satuan_kerja, nama_paket, metode_pemilihan, tanggal_buat_paket, nilai_pagu_rup, nilai_pagu_paket, kode_satuan_kerja, cara_pengadaan, jenis_pengadaan, pdn, umk, sumber_dana, kode_rup_lokal, metode_pengadaan, tipe_swakelola
	fields := []string{"kode_rup", "satuan_kerja", "nama_paket", "metode_pemilihan", "tanggal_buat_paket", "nilai_pagu_rup", "nilai_pagu_paket", "kode_satuan_kerja", "cara_pengadaan", "jenis_pengadaan", "pdn", "umk", "sumber_dana", "kode_rup_lokal", "metode_pengadaan", "tipe_swakelola"}

	for i, field := range fields {
		// Check for nil or empty string
		if itemMap[field] == nil || itemMap[field] == "" {
			if i < len(arr) && arr[i] != "" {
				itemMap[field] = strings.TrimSpace(arr[i])
			}
		}
	}
}

// assignRemainingPemilihanFields fills remaining fields using positional mapping as backup
func (ctrl SPSEController) assignRemainingPemilihanFields(arr []string, itemMap map[string]interface{}) {
	// Standard Pemilihan field order: kode_rup, satuan_kerja, nama_paket, metode_pemilihan, rencana_pemilihan, tanggal_pemilihan, nilai_hps, status_paket, kode_satuan_kerja, cara_pengadaan, jenis_pengadaan, pdn, umk, sumber_dana, kode_rup_lokal, metode_pengadaan, pagu_rup, tipe_swakelola, akhir_pemilihan
	fields := []string{"kode_rup", "satuan_kerja", "nama_paket", "metode_pemilihan", "rencana_pemilihan", "tanggal_pemilihan", "nilai_hps", "status_paket", "kode_satuan_kerja", "cara_pengadaan", "jenis_pengadaan", "pdn", "umk", "sumber_dana", "kode_rup_lokal", "metode_pengadaan", "pagu_rup", "tipe_swakelola", "akhir_pemilihan"}

	for i, field := range fields {
		// Check for nil or empty string
		if itemMap[field] == nil || itemMap[field] == "" {
			if i < len(arr) && arr[i] != "" {
				itemMap[field] = strings.TrimSpace(arr[i])
			}
		}
	}
}

// validateAndFillFields ensures all required database fields are populated with appropriate defaults
func (ctrl SPSEController) validateAndFillFields(tableName string, itemMap map[string]interface{}) map[string]interface{} {
	// Define required fields for each table with their default values
	requiredFields := make(map[string]string)

	switch tableName {
	case "perencanaan":
		requiredFields = map[string]string{
			"kode_rup": "UNKNOWN", "satuan_kerja": "N/A", "nama_paket": "N/A",
			"metode_pemilihan": "N/A", "tanggal_pengumuman": "N/A", "rencana_pemilihan": "N/A",
			"pagu_rup": "N/A", "kode_satuan_kerja": "N/A", "cara_pengadaan": "N/A",
			"jenis_pengadaan": "N/A", "pdn": "N/A", "umk": "N/A",
			"sumber_dana": "N/A", "kode_rup_lokal": "N/A", "akhir_pemilihan": "N/A",
			"tipe_swakelola": "N/A",
		}
	case "persiapan":
		requiredFields = map[string]string{
			"kode_rup": "UNKNOWN", "satuan_kerja": "N/A", "nama_paket": "N/A",
			"metode_pemilihan": "N/A", "tanggal_buat_paket": "N/A", "nilai_pagu_rup": "N/A",
			"nilai_pagu_paket": "N/A", "kode_satuan_kerja": "N/A", "cara_pengadaan": "N/A",
			"jenis_pengadaan": "N/A", "pdn": "N/A", "umk": "N/A",
			"sumber_dana": "N/A", "kode_rup_lokal": "N/A", "metode_pengadaan": "N/A",
			"tipe_swakelola": "N/A",
		}
	case "pemilihan":
		requiredFields = map[string]string{
			"kode_rup": "UNKNOWN", "satuan_kerja": "N/A", "nama_paket": "N/A",
			"metode_pemilihan": "N/A", "rencana_pemilihan": "N/A", "tanggal_pemilihan": "N/A",
			"nilai_hps": "N/A", "status_paket": "N/A", "kode_satuan_kerja": "N/A",
			"cara_pengadaan": "N/A", "jenis_pengadaan": "N/A", "pdn": "N/A",
			"umk": "N/A", "sumber_dana": "N/A", "kode_rup_lokal": "N/A",
			"metode_pengadaan": "N/A", "pagu_rup": "N/A", "tipe_swakelola": "N/A",
			"akhir_pemilihan": "N/A",
		}
	}

	// Fill missing fields with defaults
	for field, defaultValue := range requiredFields {
		if itemMap[field] == "" || itemMap[field] == nil {
			itemMap[field] = defaultValue
		}
	}

	return itemMap
}

// Helper functions for pattern recognition

// isRupCode checks if a string looks like a RUP code
func (ctrl SPSEController) isRupCode(value string) bool {
	// RUP codes are typically numeric and 8+ digits
	_, err := strconv.ParseInt(value, 10, 64)
	return err == nil && len(value) >= 8
}

// isSatuanKerja checks if a string looks like a work unit name
func (ctrl SPSEController) isSatuanKerja(value string) bool {
	// Typically contains department/institution keywords
	keywords := []string{"dinas", "kecamatan", "kelurahan", "rsud", "rs", "kantor", "badan", "dinas"}
	lowerValue := strings.ToLower(value)
	for _, keyword := range keywords {
		if strings.Contains(lowerValue, keyword) {
			return true
		}
	}
	return false
}

// isPackageName checks if a string looks like a package name
func (ctrl SPSEController) isPackageName(value string) bool {
	// Package names typically contain procurement-related keywords
	keywords := []string{"pembangunan", "pengadaan", "jasa", "konsultasi", "konstruksi", "maintenance"}
	lowerValue := strings.ToLower(value)
	for _, keyword := range keywords {
		if strings.Contains(lowerValue, keyword) {
			return true
		}
	}
	// Also check if it's a reasonable length for a package name
	return len(value) > 10 && len(value) < 200
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

// assignToNextAvailableField assigns a value to the next unfilled field for a specific table
func (ctrl SPSEController) assignToNextAvailableField(itemMap map[string]interface{}, value, tableName string) {
	var fields []string

	switch tableName {
	case "perencanaan":
		fields = []string{"rencana_pemilihan", "kode_satuan_kerja", "cara_pengadaan", "jenis_pengadaan", "pdn", "umk", "sumber_dana", "kode_rup_lokal", "akhir_pemilihan", "tipe_swakelola"}
	case "persiapan":
		fields = []string{"tanggal_buat_paket", "nilai_pagu_paket", "kode_satuan_kerja", "cara_pengadaan", "jenis_pengadaan", "pdn", "umk", "sumber_dana", "kode_rup_lokal", "metode_pengadaan", "tipe_swakelola"}
	case "pemilihan":
		fields = []string{"rencana_pemilihan", "status_paket", "kode_satuan_kerja", "cara_pengadaan", "jenis_pengadaan", "pdn", "umk", "sumber_dana", "kode_rup_lokal", "metode_pengadaan", "pagu_rup", "tipe_swakelola", "akhir_pemilihan"}
	}

	// Find the first empty field and assign the value
	for _, field := range fields {
		if itemMap[field] == "" {
			itemMap[field] = value
			break
		}
	}
}

// ScrapePerencanaan godoc
// @Summary Scrape Perencanaan (Planning) data from SPSE API
// @Schemes
// @Description Scrapes perencanaan planning data from SPSE API endpoint
// @Tags SPSE
// @Accept json
// @Produce json
// @Success 200 {object} ScrapingResult
// @Failure 500 {object} gin.H
// @Router /spse/scraper/perencanaan [POST]
func (ctrl SPSEController) ScrapePerencanaan(c *gin.Context) {
	log.Println("Starting Perencanaan scraping...")

	result, err := ctrl.scrapeEndpoint("/sumedangkab/amel/dt/detailperencanaan2", "perencanaan")
	if err != nil {
		log.Printf("Error scraping Perencanaan: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
			"message": "Failed to scrape Perencanaan data",
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ScrapePersiapan godoc
// @Summary Scrape Persiapan (Preparation) data from SPSE API
// @Schemes
// @Description Scrapes perencanaan preparation data from SPSE API endpoint
// @Tags SPSE
// @Accept json
// @Produce json
// @Success 200 {object} ScrapingResult
// @Failure 500 {object} gin.H
// @Router /spse/scraper/persiapan [POST]
func (ctrl SPSEController) ScrapePersiapan(c *gin.Context) {
	log.Println("Starting Persiapan scraping...")

	result, err := ctrl.scrapeEndpoint("/sumedangkab/amel/dt/detailpersiapan2", "persiapan")
	if err != nil {
		log.Printf("Error scraping Persiapan: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
			"message": "Failed to scrape Persiapan data",
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ScrapePemilihan godoc
// @Summary Scrape Pemilihan (Selection) data from SPSE API
// @Schemes
// @Description Scrapes perencanaan selection/contract data from SPSE API endpoint
// @Tags SPSE
// @Accept json
// @Produce json
// @Success 200 {object} ScrapingResult
// @Failure 500 {object} gin.H
// @Router /spse/scraper/pemilihan [POST]
func (ctrl SPSEController) ScrapePemilihan(c *gin.Context) {
	log.Println("Starting Pemilihan scraping...")

	result, err := ctrl.scrapeEndpoint("/sumedangkab/amel/dt/detailpemilihan2", "pemilihan")
	if err != nil {
		log.Printf("Error scraping Pemilihan: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
			"message": "Failed to scrape Pemilihan data",
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ScrapeAll godoc
// @Summary Scrape all SPSE perencanaan data endpoints
// @Schemes
// @Description Scrapes data from all three SPSE perencanaan endpoints simultaneously
// @Tags SPSE
// @Accept json
// @Produce json
// @Success 200 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /spse/scraper/all [POST]
func (ctrl SPSEController) ScrapeAll(c *gin.Context) {
	log.Println("Starting comprehensive SPSE scraping...")

	endpoints := []struct {
		path  string
		table string
		name  string
	}{
		{"/sumedangkab/amel/dt/detailperencanaan2", "perencanaan", "Perencanaan"},
		{"/sumedangkab/amel/dt/detailpersiapan2", "persiapan", "Persiapan"},
		{"/sumedangkab/amel/dt/detailpemilihan2", "pemilihan", "Pemilihan"},
	}

	results := make(map[string]ScrapingResult)
	successCount := 0

	for _, endpoint := range endpoints {
		result, err := ctrl.scrapeEndpoint(endpoint.path, endpoint.table)
		if err != nil {
			log.Printf("Error scraping %s: %v", endpoint.name, err)
			results[endpoint.name] = ScrapingResult{
				Success: false,
				Message: err.Error(),
			}
		} else {
			results[endpoint.name] = result
			if result.Success {
				successCount++
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"results":         results,
		"total_success":   successCount,
		"total_endpoints": len(endpoints),
		"message":         fmt.Sprintf("Scraping completed: %d/%d endpoints successful", successCount, len(endpoints)),
	})
}

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

	// Get counts from each table
	err := database.QueryRow("SELECT COUNT(*) FROM spse_perencanaan").Scan(&perencanaanCount)
	if err != nil {
		log.Printf("Error getting perencanaan count: %v", err)
		perencanaanCount = 0
	}

	err = database.QueryRow("SELECT COUNT(*) FROM spse_persiapan").Scan(&persiapanCount)
	if err != nil {
		log.Printf("Error getting persiapan count: %v", err)
		persiapanCount = 0
	}

	err = database.QueryRow("SELECT COUNT(*) FROM spse_pemilihan").Scan(&pemilihanCount)
	if err != nil {
		log.Printf("Error getting pemilihan count: %v", err)
		pemilihanCount = 0
	}

	total := perencanaanCount + persiapanCount + pemilihanCount

	c.JSON(http.StatusOK, gin.H{
		"statistics": gin.H{
			"perencanaan": perencanaanCount,
			"persiapan":   persiapanCount,
			"pemilihan":   pemilihanCount,
			"total":       total,
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
			   akhir_pemilihan, tipe_swakelola, created_at, last_update
		FROM spse_perencanaan
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)

	if err != nil {
		log.Printf("Error querying perencanaan data: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to retrieve data",
		})
		return
	}
	defer rows.Close()

	var data []PerencanaanData
	for rows.Next() {
		var item PerencanaanData
		err := rows.Scan(&item.ID, &item.KodeRUP, &item.SatuanKerja, &item.NamaPaket,
			&item.MetodePemilihan, &item.TanggalPengumuman, &item.RencanaPemilihan,
			&item.PaguRUP, &item.KodeSatuanKerja, &item.CaraPengadaan, &item.JenisPengadaan,
			&item.PDN, &item.UMK, &item.SumberDana, &item.KodeRUPLokal, &item.AkhirPemilihan,
			&item.TipeSwakelola, &item.CreatedAt, &item.LastUpdate)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
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
			   metode_pengadaan, tipe_swakelola, created_at, last_update
		FROM spse_persiapan
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)

	if err != nil {
		log.Printf("Error querying persiapan data: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to retrieve data",
		})
		return
	}
	defer rows.Close()

	var data []PersiapanData
	for rows.Next() {
		var item PersiapanData
		err := rows.Scan(&item.ID, &item.KodeRUP, &item.SatuanKerja, &item.NamaPaket,
			&item.MetodePemilihan, &item.TanggalBuatPaket, &item.NilaiPaguRUP,
			&item.NilaiPaguPaket, &item.KodeSatuanKerja, &item.CaraPengadaan, &item.JenisPengadaan,
			&item.PDN, &item.UMK, &item.SumberDana, &item.KodeRUPLokal, &item.MetodePengadaan,
			&item.TipeSwakelola, &item.CreatedAt, &item.LastUpdate)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
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
			   tipe_swakelola, akhir_pemilihan, created_at, last_update
		FROM spse_pemilihan
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)

	if err != nil {
		log.Printf("Error querying pemilihan data: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to retrieve data",
		})
		return
	}
	defer rows.Close()

	var data []PemilihanData
	for rows.Next() {
		var item PemilihanData
		err := rows.Scan(&item.ID, &item.KodeRUP, &item.SatuanKerja, &item.NamaPaket,
			&item.MetodePemilihan, &item.RencanaPemilihan, &item.TanggalPemilihan,
			&item.NilaiHPS, &item.StatusPaket, &item.KodeSatuanKerja, &item.CaraPengadaan,
			&item.JenisPengadaan, &item.PDN, &item.UMK, &item.SumberDana, &item.KodeRUPLokal,
			&item.MetodePengadaan, &item.PaguRUP, &item.TipeSwakelola, &item.AkhirPemilihan,
			&item.CreatedAt, &item.LastUpdate)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
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
