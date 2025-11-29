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

// HasilPemilihanData represents election results stage data (matches SPSEHasilPemilihan model)
type HasilPemilihanData struct {
	ID                    int64     `json:"id"`
	KodeRUP               string    `json:"kode_rup"`
	SatuanKerja           string    `json:"satuan_kerja"`
	NamaPaket             string    `json:"nama_paket"`
	MetodePemilihan       string    `json:"metode_pemilihan"`
	TanggalHasilPemilihan string    `json:"tanggal_hasil_pemilihan"`
	NilaiHasilPemilihan   string    `json:"nilai_hasil_pemilihan"`
	StatusPaket           string    `json:"status_paket"`
	KodeSatuanKerja       string    `json:"kode_satuan_kerja"`
	CaraPengadaan         string    `json:"cara_pengadaan"`
	JenisPengadaan        string    `json:"jenis_pengadaan"`
	PDN                   string    `json:"pdn"`
	UMK                   string    `json:"umk"`
	SumberDana            string    `json:"sumber_dana"`
	KodeRUPLokal          string    `json:"kode_rup_lokal"`
	MetodePengadaan       string    `json:"metode_pengadaan"`
	PaguRUP               string    `json:"pagu_rup"`
	TipeSwakelola         string    `json:"tipe_swakelola"`
	CreatedAt             time.Time `json:"created_at"`
	LastUpdate            int64     `json:"last_update"`
}

// KontrakData represents contract stage data (matches SPSEKontrak model)
type KontrakData struct {
	ID              int64     `json:"id"`
	KodeRUP         string    `json:"kode_rup"`
	SatuanKerja     string    `json:"satuan_kerja"`
	NamaPaket       string    `json:"nama_paket"`
	MetodePemilihan string    `json:"metode_pemilihan"`
	TanggalKontrak  string    `json:"tanggal_kontrak"`
	NilaiKontrak    string    `json:"nilai_kontrak"`
	StatusPaket     string    `json:"status_paket"`
	MulaiKontrak    string    `json:"mulai_kontrak"`
	NilaiBAP        string    `json:"nilai_bap"`
	SelesaiKontrak  string    `json:"selesai_kontrak"`
	KodeSatuanKerja string    `json:"kode_satuan_kerja"`
	CaraPengadaan   string    `json:"cara_pengadaan"`
	JenisPengadaan  string    `json:"jenis_pengadaan"`
	PDN             string    `json:"pdn"`
	UMK             string    `json:"umk"`
	SumberDana      string    `json:"sumber_dana"`
	KodeRUPLokal    string    `json:"kode_rup_lokal"`
	MetodePengadaan string    `json:"metode_pengadaan"`
	TipeSwakelola   string    `json:"tipe_swakelola"`
	CreatedAt       time.Time `json:"created_at"`
	LastUpdate      int64     `json:"last_update"`
}

// SerahTerimaData represents handover stage data (matches SPSESerahTerima model)
type SerahTerimaData struct {
	ID                 int64     `json:"id"`
	KodeRUP            string    `json:"kode_rup"`
	SatuanKerja        string    `json:"satuan_kerja"`
	NamaPaket          string    `json:"nama_paket"`
	MetodePemilihan    string    `json:"metode_pemilihan"`
	TanggalSerahTerima string    `json:"tanggal_serah_terima"`
	NilaiBAP           string    `json:"nilai_bap"`
	StatusPaket        string    `json:"status_paket"`
	KodeSatuanKerja    string    `json:"kode_satuan_kerja"`
	CaraPengadaan      string    `json:"cara_pengadaan"`
	JenisPengadaan     string    `json:"jenis_pengadaan"`
	PDN                string    `json:"pdn"`
	UMK                string    `json:"umk"`
	SumberDana         string    `json:"sumber_dana"`
	KodeRUPLokal       string    `json:"kode_rup_lokal"`
	MetodePengadaan    string    `json:"metode_pengadaan"`
	TipeSwakelola      string    `json:"tipe_swakelola"`
	CreatedAt          time.Time `json:"created_at"`
	LastUpdate         int64     `json:"last_update"`
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
		var orderedDataset *OrderedDataSet

		// Handle different data formats with ordered dataset approach
		switch v := item.(type) {
		case map[string]interface{}:
			// Convert map to ordered dataset
			orderedDataset = ctrl.convertMapToOrderedDataset(tableName, v)
		case []interface{}:
			// Array format - convert to ordered dataset with precise field mapping
			orderedDataset = ctrl.mapArrayToOrderedDataset(tableName, v)
		case string:
			// String format - try to parse as delimited data
			itemMap := ctrl.parseStringData(v)
			orderedDataset = ctrl.convertMapToOrderedDataset(tableName, itemMap)
		default:
			recordsFailed++
			continue
		}

		// Log mapping quality for debugging
		if orderedDataset.MappingStatus.MappedFields < orderedDataset.MappingStatus.TotalFields/2 {
			log.Printf("Warning: Poor field mapping for table %s - only %d/%d fields mapped",
				tableName, orderedDataset.MappingStatus.MappedFields, orderedDataset.MappingStatus.TotalFields)
		}

		// Store based on table type using ordered dataset
		var insertQuery string
		switch tableName {
		case "perencanaan":
			insertQuery = ctrl.buildPerencanaanInsertFromDataset(orderedDataset)
		case "persiapan":
			insertQuery = ctrl.buildPersiapanInsertFromDataset(orderedDataset)
		case "pemilihan":
			insertQuery = ctrl.buildPemilihanInsertFromDataset(orderedDataset)
		case "hasilpemilihan":
			insertQuery = ctrl.buildHasilPemilihanInsertFromDataset(orderedDataset)
		case "kontrak":
			insertQuery = ctrl.buildKontrakInsertFromDataset(orderedDataset)
		case "serahterima":
			insertQuery = ctrl.buildSerahTerimaInsertFromDataset(orderedDataset)
		}

		if insertQuery != "" {
			_, err := database.Exec(insertQuery)
			if err != nil {
				recordsFailed++
				log.Printf("Error inserting record: %v", err)
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

// buildPerencanaanInsertFromDataset builds INSERT query for perencanaan data using ordered dataset
func (ctrl SPSEController) buildPerencanaanInsertFromDataset(dataset *OrderedDataSet) string {
	if dataset == nil || dataset.FieldValues == nil {
		return ""
	}

	// Extract fields with proper escaping from ordered dataset
	kodeRUP := ctrl.escapeString(dataset.FieldValues["kode_rup"])
	satuanKerja := ctrl.escapeString(dataset.FieldValues["satuan_kerja"])
	namaPaket := ctrl.escapeString(dataset.FieldValues["nama_paket"])
	metodePemilihan := ctrl.escapeString(dataset.FieldValues["metode_pemilihan"])
	tanggalPengumuman := ctrl.escapeString(dataset.FieldValues["tanggal_pengumuman"])
	rencanaPemilihan := ctrl.escapeString(dataset.FieldValues["rencana_pemilihan"])
	paguRUP := ctrl.escapeString(dataset.FieldValues["pagu_rup"])
	kodeSatuanKerja := ctrl.escapeString(dataset.FieldValues["kode_satuan_kerja"])
	caraPengadaan := ctrl.escapeString(dataset.FieldValues["cara_pengadaan"])
	jenisPengadaan := ctrl.escapeString(dataset.FieldValues["jenis_pengadaan"])
	pdn := ctrl.escapeString(dataset.FieldValues["pdn"])
	umk := ctrl.escapeString(dataset.FieldValues["umk"])
	sumberDana := ctrl.escapeString(dataset.FieldValues["sumber_dana"])
	kodeRUPLokal := ctrl.escapeString(dataset.FieldValues["kode_rup_lokal"])
	akhirPemilihan := ctrl.escapeString(dataset.FieldValues["akhir_pemilihan"])
	tipeSwakelola := ctrl.escapeString(dataset.FieldValues["tipe_swakelola"])

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

// buildPersiapanInsertFromDataset builds INSERT query for persiapan data using ordered dataset
func (ctrl SPSEController) buildPersiapanInsertFromDataset(dataset *OrderedDataSet) string {
	if dataset == nil || dataset.FieldValues == nil {
		return ""
	}

	// Extract fields with proper escaping from ordered dataset
	kodeRUP := ctrl.escapeString(dataset.FieldValues["kode_rup"])
	satuanKerja := ctrl.escapeString(dataset.FieldValues["satuan_kerja"])
	namaPaket := ctrl.escapeString(dataset.FieldValues["nama_paket"])
	metodePemilihan := ctrl.escapeString(dataset.FieldValues["metode_pemilihan"])
	tanggalBuatPaket := ctrl.escapeString(dataset.FieldValues["tanggal_buat_paket"])
	nilaiPaguRUP := ctrl.escapeString(dataset.FieldValues["nilai_pagu_rup"])
	nilaiPaguPaket := ctrl.escapeString(dataset.FieldValues["nilai_pagu_paket"])
	kodeSatuanKerja := ctrl.escapeString(dataset.FieldValues["kode_satuan_kerja"])
	caraPengadaan := ctrl.escapeString(dataset.FieldValues["cara_pengadaan"])
	jenisPengadaan := ctrl.escapeString(dataset.FieldValues["jenis_pengadaan"])
	pdn := ctrl.escapeString(dataset.FieldValues["pdn"])
	umk := ctrl.escapeString(dataset.FieldValues["umk"])
	sumberDana := ctrl.escapeString(dataset.FieldValues["sumber_dana"])
	kodeRUPLokal := ctrl.escapeString(dataset.FieldValues["kode_rup_lokal"])
	metodePengadaan := ctrl.escapeString(dataset.FieldValues["metode_pengadaan"])
	tipeSwakelola := ctrl.escapeString(dataset.FieldValues["tipe_swakelola"])

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

// buildPemilihanInsertFromDataset builds INSERT query for pemilihan data using ordered dataset
func (ctrl SPSEController) buildPemilihanInsertFromDataset(dataset *OrderedDataSet) string {
	if dataset == nil || dataset.FieldValues == nil {
		return ""
	}

	// Extract fields with proper escaping from ordered dataset
	kodeRUP := ctrl.escapeString(dataset.FieldValues["kode_rup"])
	satuanKerja := ctrl.escapeString(dataset.FieldValues["satuan_kerja"])
	namaPaket := ctrl.escapeString(dataset.FieldValues["nama_paket"])
	metodePemilihan := ctrl.escapeString(dataset.FieldValues["metode_pemilihan"])
	rencanaPemilihan := ctrl.escapeString(dataset.FieldValues["rencana_pemilihan"])
	tanggalPemilihan := ctrl.escapeString(dataset.FieldValues["tanggal_pemilihan"])
	nilaiHPS := ctrl.escapeString(dataset.FieldValues["nilai_hps"])
	statusPaket := ctrl.escapeString(dataset.FieldValues["status_paket"])
	kodeSatuanKerja := ctrl.escapeString(dataset.FieldValues["kode_satuan_kerja"])
	caraPengadaan := ctrl.escapeString(dataset.FieldValues["cara_pengadaan"])
	jenisPengadaan := ctrl.escapeString(dataset.FieldValues["jenis_pengadaan"])
	pdn := ctrl.escapeString(dataset.FieldValues["pdn"])
	umk := ctrl.escapeString(dataset.FieldValues["umk"])
	sumberDana := ctrl.escapeString(dataset.FieldValues["sumber_dana"])
	kodeRUPLokal := ctrl.escapeString(dataset.FieldValues["kode_rup_lokal"])
	metodePengadaan := ctrl.escapeString(dataset.FieldValues["metode_pengadaan"])
	paguRUP := ctrl.escapeString(dataset.FieldValues["pagu_rup"])
	tipeSwakelola := ctrl.escapeString(dataset.FieldValues["tipe_swakelola"])
	akhirPemilihan := ctrl.escapeString(dataset.FieldValues["akhir_pemilihan"])

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

// buildHasilPemilihanInsertFromDataset builds INSERT query for hasilpemilihan data using ordered dataset
func (ctrl SPSEController) buildHasilPemilihanInsertFromDataset(dataset *OrderedDataSet) string {
	if dataset == nil || dataset.FieldValues == nil {
		return ""
	}

	// Extract fields with proper escaping from ordered dataset
	kodeRUP := ctrl.escapeString(dataset.FieldValues["kode_rup"])
	satuanKerja := ctrl.escapeString(dataset.FieldValues["satuan_kerja"])
	namaPaket := ctrl.escapeString(dataset.FieldValues["nama_paket"])
	metodePemilihan := ctrl.escapeString(dataset.FieldValues["metode_pemilihan"])
	tanggalHasilPemilihan := ctrl.escapeString(dataset.FieldValues["tanggal_hasil_pemilihan"])
	nilaiHasilPemilihan := ctrl.escapeString(dataset.FieldValues["nilai_hasil_pemilihan"])
	statusPaket := ctrl.escapeString(dataset.FieldValues["status_paket"])
	kodeSatuanKerja := ctrl.escapeString(dataset.FieldValues["kode_satuan_kerja"])
	caraPengadaan := ctrl.escapeString(dataset.FieldValues["cara_pengadaan"])
	jenisPengadaan := ctrl.escapeString(dataset.FieldValues["jenis_pengadaan"])
	pdn := ctrl.escapeString(dataset.FieldValues["pdn"])
	umk := ctrl.escapeString(dataset.FieldValues["umk"])
	sumberDana := ctrl.escapeString(dataset.FieldValues["sumber_dana"])
	kodeRUPLokal := ctrl.escapeString(dataset.FieldValues["kode_rup_lokal"])
	metodePengadaan := ctrl.escapeString(dataset.FieldValues["metode_pengadaan"])
	paguRUP := ctrl.escapeString(dataset.FieldValues["pagu_rup"])
	tipeSwakelola := ctrl.escapeString(dataset.FieldValues["tipe_swakelola"])

	return fmt.Sprintf(`INSERT INTO spse_hasilpemilihan
		(kode_rup, satuan_kerja, nama_paket, metode_pemilihan, tanggal_hasil_pemilihan, nilai_hasil_pemilihan, status_paket,
		 kode_satuan_kerja, cara_pengadaan, jenis_pengadaan, pdn, umk, sumber_dana, kode_rup_lokal,
		 metode_pengadaan, pagu_rup, tipe_swakelola, created_at, last_update)
		VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', NOW(), %d)
		ON CONFLICT (kode_rup, nama_paket) DO NOTHING`,
		kodeRUP, satuanKerja, namaPaket, metodePemilihan, tanggalHasilPemilihan, nilaiHasilPemilihan, statusPaket,
		kodeSatuanKerja, caraPengadaan, jenisPengadaan, pdn, umk, sumberDana, kodeRUPLokal,
		metodePengadaan, paguRUP, tipeSwakelola, time.Now().Unix())
}

// buildKontrakInsertFromDataset builds INSERT query for kontrak data using ordered dataset
func (ctrl SPSEController) buildKontrakInsertFromDataset(dataset *OrderedDataSet) string {
	if dataset == nil || dataset.FieldValues == nil {
		return ""
	}

	// Extract fields with proper escaping from ordered dataset
	kodeRUP := ctrl.escapeString(dataset.FieldValues["kode_rup"])
	satuanKerja := ctrl.escapeString(dataset.FieldValues["satuan_kerja"])
	namaPaket := ctrl.escapeString(dataset.FieldValues["nama_paket"])
	metodePemilihan := ctrl.escapeString(dataset.FieldValues["metode_pemilihan"])
	tanggalKontrak := ctrl.escapeString(dataset.FieldValues["tanggal_kontrak"])
	nilaiKontrak := ctrl.escapeString(dataset.FieldValues["nilai_kontrak"])
	statusPaket := ctrl.escapeString(dataset.FieldValues["status_paket"])
	mulaiKontrak := ctrl.escapeString(dataset.FieldValues["mulai_kontrak"])
	nilaiBAP := ctrl.escapeString(dataset.FieldValues["nilai_bap"])
	selesaiKontrak := ctrl.escapeString(dataset.FieldValues["selesai_kontrak"])
	kodeSatuanKerja := ctrl.escapeString(dataset.FieldValues["kode_satuan_kerja"])
	caraPengadaan := ctrl.escapeString(dataset.FieldValues["cara_pengadaan"])
	jenisPengadaan := ctrl.escapeString(dataset.FieldValues["jenis_pengadaan"])
	pdn := ctrl.escapeString(dataset.FieldValues["pdn"])
	umk := ctrl.escapeString(dataset.FieldValues["umk"])
	sumberDana := ctrl.escapeString(dataset.FieldValues["sumber_dana"])
	kodeRUPLokal := ctrl.escapeString(dataset.FieldValues["kode_rup_lokal"])
	metodePengadaan := ctrl.escapeString(dataset.FieldValues["metode_pengadaan"])
	tipeSwakelola := ctrl.escapeString(dataset.FieldValues["tipe_swakelola"])

	return fmt.Sprintf(`INSERT INTO spse_kontrak
		(kode_rup, satuan_kerja, nama_paket, metode_pemilihan, tanggal_kontrak, nilai_kontrak, status_paket,
		 mulai_kontrak, nilai_bap, selesai_kontrak, kode_satuan_kerja, cara_pengadaan, jenis_pengadaan,
		 pdn, umk, sumber_dana, kode_rup_lokal, metode_pengadaan, tipe_swakelola, created_at, last_update)
		VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', NOW(), %d)
		ON CONFLICT (kode_rup, nama_paket) DO NOTHING`,
		kodeRUP, satuanKerja, namaPaket, metodePemilihan, tanggalKontrak, nilaiKontrak, statusPaket,
		mulaiKontrak, nilaiBAP, selesaiKontrak, kodeSatuanKerja, caraPengadaan, jenisPengadaan,
		pdn, umk, sumberDana, kodeRUPLokal, metodePengadaan, tipeSwakelola, time.Now().Unix())
}

// buildSerahTerimaInsertFromDataset builds INSERT query for serahterima data using ordered dataset
func (ctrl SPSEController) buildSerahTerimaInsertFromDataset(dataset *OrderedDataSet) string {
	if dataset == nil || dataset.FieldValues == nil {
		return ""
	}

	// Extract fields with proper escaping from ordered dataset
	kodeRUP := ctrl.escapeString(dataset.FieldValues["kode_rup"])
	satuanKerja := ctrl.escapeString(dataset.FieldValues["satuan_kerja"])
	namaPaket := ctrl.escapeString(dataset.FieldValues["nama_paket"])
	metodePemilihan := ctrl.escapeString(dataset.FieldValues["metode_pemilihan"])
	tanggalSerahTerima := ctrl.escapeString(dataset.FieldValues["tanggal_serah_terima"])
	nilaiBAP := ctrl.escapeString(dataset.FieldValues["nilai_bap"])
	statusPaket := ctrl.escapeString(dataset.FieldValues["status_paket"])
	kodeSatuanKerja := ctrl.escapeString(dataset.FieldValues["kode_satuan_kerja"])
	caraPengadaan := ctrl.escapeString(dataset.FieldValues["cara_pengadaan"])
	jenisPengadaan := ctrl.escapeString(dataset.FieldValues["jenis_pengadaan"])
	pdn := ctrl.escapeString(dataset.FieldValues["pdn"])
	umk := ctrl.escapeString(dataset.FieldValues["umk"])
	sumberDana := ctrl.escapeString(dataset.FieldValues["sumber_dana"])
	kodeRUPLokal := ctrl.escapeString(dataset.FieldValues["kode_rup_lokal"])
	metodePengadaan := ctrl.escapeString(dataset.FieldValues["metode_pengadaan"])
	tipeSwakelola := ctrl.escapeString(dataset.FieldValues["tipe_swakelola"])

	return fmt.Sprintf(`INSERT INTO spse_serahterima
		(kode_rup, satuan_kerja, nama_paket, metode_pemilihan, tanggal_serah_terima, nilai_bap, status_paket,
		 kode_satuan_kerja, cara_pengadaan, jenis_pengadaan, pdn, umk, sumber_dana, kode_rup_lokal,
		 metode_pengadaan, tipe_swakelola, created_at, last_update)
		VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', NOW(), %d)
		ON CONFLICT (kode_rup, nama_paket) DO NOTHING`,
		kodeRUP, satuanKerja, namaPaket, metodePemilihan, tanggalSerahTerima, nilaiBAP, statusPaket,
		kodeSatuanKerja, caraPengadaan, jenisPengadaan, pdn, umk, sumberDana, kodeRUPLokal,
		metodePengadaan, tipeSwakelola, time.Now().Unix())
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

// SPSEFieldMapping defines the precise field mapping structure for database columns
type SPSEFieldMapping struct {
	TableName       string                       `json:"table_name"`
	FieldOrder      []string                     `json:"field_order"`
	RequiredFields  map[string]string            `json:"required_fields"`
	FieldValidators map[string]func(string) bool `json:"-"`
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

// ScrapeHasilPemilihan godoc
// @Summary Scrape Hasil Pemilihan (Election Results) data from SPSE API
// @Schemes
// @Description Scrapes election results data from SPSE API endpoint
// @Tags SPSE
// @Accept json
// @Produce json
// @Success 200 {object} ScrapingResult
// @Failure 500 {object} gin.H
// @Router /spse/scraper/hasilpemilihan [POST]
func (ctrl SPSEController) ScrapeHasilPemilihan(c *gin.Context) {
	log.Println("Starting Hasil Pemilihan scraping...")

	result, err := ctrl.scrapeEndpoint("/sumedangkab/amel/dt/detailhasilpemilihan2", "hasilpemilihan")
	if err != nil {
		log.Printf("Error scraping Hasil Pemilihan: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
			"message": "Failed to scrape Hasil Pemilihan data",
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ScrapeKontrak godoc
// @Summary Scrape Kontrak (Contract) data from SPSE API
// @Schemes
// @Description Scrapes kontrak contract data from SPSE API endpoint
// @Tags SPSE
// @Accept json
// @Produce json
// @Success 200 {object} ScrapingResult
// @Failure 500 {object} gin.H
// @Router /spse/scraper/kontrak [POST]
func (ctrl SPSEController) ScrapeKontrak(c *gin.Context) {
	log.Println("Starting Kontrak scraping...")

	result, err := ctrl.scrapeEndpoint("/sumedangkab/amel/dt/detailkontrak2", "kontrak")
	if err != nil {
		log.Printf("Error scraping Kontrak: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
			"message": "Failed to scrape Kontrak data",
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ScrapeSerahTerima godoc
// @Summary Scrape Serah Terima (Handover) data from SPSE API
// @Schemes
// @Description Scrapes serah terima handover data from SPSE API endpoint
// @Tags SPSE
// @Accept json
// @Produce json
// @Success 200 {object} ScrapingResult
// @Failure 500 {object} gin.H
// @Router /spse/scraper/serahterima [POST]
func (ctrl SPSEController) ScrapeSerahTerima(c *gin.Context) {
	log.Println("Starting Serah Terima scraping...")

	result, err := ctrl.scrapeEndpoint("/sumedangkab/amel/dt/detailserahterima2", "serahterima")
	if err != nil {
		log.Printf("Error scraping Serah Terima: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
			"message": "Failed to scrape Serah Terima data",
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
		{"/sumedangkab/amel/dt/detailhasilpemilihan2", "hasilpemilihan", "Hasil Pemilihan"},
		{"/sumedangkab/amel/dt/detailkontrak2", "kontrak", "Kontrak"},
		{"/sumedangkab/amel/dt/detailserahterima2", "serahterima", "Serah Terima"},
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
	var hasilPemilihanCount int
	var kontrakCount int
	var serahTerimaCount int

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

	err = database.QueryRow("SELECT COUNT(*) FROM spse_hasilpemilihan").Scan(&hasilPemilihanCount)
	if err != nil {
		log.Printf("Error getting hasil pemilihan count: %v", err)
		hasilPemilihanCount = 0
	}

	err = database.QueryRow("SELECT COUNT(*) FROM spse_kontrak").Scan(&kontrakCount)
	if err != nil {
		log.Printf("Error getting kontrak count: %v", err)
		kontrakCount = 0
	}

	err = database.QueryRow("SELECT COUNT(*) FROM spse_serahterima").Scan(&serahTerimaCount)
	if err != nil {
		log.Printf("Error getting serah terima count: %v", err)
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
			   metode_pengadaan, tipe_swakelola, created_at, last_update
		FROM spse_kontrak
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)

	if err != nil {
		log.Printf("Error querying kontrak data: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to retrieve data",
		})
		return
	}
	defer rows.Close()

	var data []KontrakData
	for rows.Next() {
		var item KontrakData
		err := rows.Scan(&item.ID, &item.KodeRUP, &item.SatuanKerja, &item.NamaPaket,
			&item.MetodePemilihan, &item.TanggalKontrak, &item.NilaiKontrak, &item.StatusPaket,
			&item.MulaiKontrak, &item.NilaiBAP, &item.SelesaiKontrak, &item.KodeSatuanKerja,
			&item.CaraPengadaan, &item.JenisPengadaan, &item.PDN, &item.UMK, &item.SumberDana,
			&item.KodeRUPLokal, &item.MetodePengadaan, &item.TipeSwakelola, &item.CreatedAt, &item.LastUpdate)
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
			   metode_pengadaan, tipe_swakelola, created_at, last_update
		FROM spse_serahterima
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)

	if err != nil {
		log.Printf("Error querying serah terima data: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to retrieve data",
		})
		return
	}
	defer rows.Close()

	var data []SerahTerimaData
	for rows.Next() {
		var item SerahTerimaData
		err := rows.Scan(&item.ID, &item.KodeRUP, &item.SatuanKerja, &item.NamaPaket,
			&item.MetodePemilihan, &item.TanggalSerahTerima, &item.NilaiBAP, &item.StatusPaket,
			&item.KodeSatuanKerja, &item.CaraPengadaan, &item.JenisPengadaan, &item.PDN,
			&item.UMK, &item.SumberDana, &item.KodeRUPLokal, &item.MetodePengadaan,
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
			   tipe_swakelola, created_at, last_update
		FROM spse_hasilpemilihan
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)

	if err != nil {
		log.Printf("Error querying hasil pemilihan data: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to retrieve data",
		})
		return
	}
	defer rows.Close()

	var data []HasilPemilihanData
	for rows.Next() {
		var item HasilPemilihanData
		err := rows.Scan(&item.ID, &item.KodeRUP, &item.SatuanKerja, &item.NamaPaket,
			&item.MetodePemilihan, &item.TanggalHasilPemilihan, &item.NilaiHasilPemilihan,
			&item.StatusPaket, &item.KodeSatuanKerja, &item.CaraPengadaan, &item.JenisPengadaan,
			&item.PDN, &item.UMK, &item.SumberDana, &item.KodeRUPLokal, &item.MetodePengadaan,
			&item.PaguRUP, &item.TipeSwakelola, &item.CreatedAt, &item.LastUpdate)
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
