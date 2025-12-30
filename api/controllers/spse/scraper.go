package spse

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/Massad/gin-boilerplate/db"
	"github.com/gin-gonic/gin"
)

// SPSEController handles SPSE perencanaan data scraping
type SPSEController struct{}

// ScrapingResult represents the result of scraping operations
type ScrapingResult struct {
	Success      bool   `json:"success"`
	Message      string `json:"message"`
	RecordsFound int    `json:"records_found"`
	Endpoint     string `json:"endpoint"`
}

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
	targetURL := os.Getenv("SPSE_BASE_URL") + os.Getenv("SPSE_AUTH_ENDPOINT")

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
		matches = re.FindStringSubmatch(content)
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
	activeInstansi := os.Getenv("SPSE_ACTIVE_INSTANSI")
	if activeInstansi == "" {
		activeInstansi = "D118"
	}
	formData.Set("activeInstansi", activeInstansi)
	activeYear := os.Getenv("SPSE_ACTIVE_YEAR")
	if activeYear == "" {
		activeYear = "2025"
	}
	formData.Set("activeYear", activeYear)

	// Prepare request
	baseURL := os.Getenv("SPSE_BASE_URL")
	if baseURL == "" {
		baseURL = "https://spse.inaproc.id"
	}
	req, err := http.NewRequest("POST", baseURL+endpoint, strings.NewReader(formData.Encode()))
	if err != nil {
		return ScrapingResult{Success: false, Message: fmt.Sprintf("Failed to create request: %v", err)}, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	referer := os.Getenv("SPSE_REFERER")
	if referer == "" {
		referer = "https://spse.inaproc.id/sumedangkab/amel"
	}
	req.Header.Set("Referer", referer)

	// Execute request with retry logic
	var resp *http.Response
	maxRetries := 3
	baseDelay := 1 * time.Second

	for attempt := 0; attempt < maxRetries; attempt++ {
		resp, err = spseClient.Do(req)
		if err == nil && resp.StatusCode < 500 {
			break
		}

		// Close failed response body to avoid resource leak
		if resp != nil {
			resp.Body.Close()
		}

		if attempt < maxRetries-1 {
			delay := time.Duration(attempt+1) * baseDelay
			log.Printf("Request failed (attempt %d/%d), retrying in %v: %v", attempt+1, maxRetries, delay, err)
			time.Sleep(delay)
		}
	}

	if err != nil {
		return ScrapingResult{Success: false, Message: fmt.Sprintf("Failed to execute request after %d attempts: %v", maxRetries, err)}, err
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

	// Collect existing kode_rup values before scraping for soft deletion logic
	existingKodeRUPs := make(map[string]bool)
	fullTableName := fmt.Sprintf("spse_%s", tableName)
	existingQuery := fmt.Sprintf("SELECT kode_rup FROM %s WHERE deleted_at IS NULL", fullTableName)
	rows, err := database.Query(existingQuery)
	if err != nil {
		log.Printf("ERROR: Failed to query existing records: %v", err)
		return ScrapingResult{
			Success:      false,
			Message:      fmt.Sprintf("Failed to query existing records: %v", err),
			RecordsFound: 0,
			Endpoint:     endpoint,
		}, fmt.Errorf("existing records query failed: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var kodeRUP string
		if err := rows.Scan(&kodeRUP); err != nil {
			log.Printf("ERROR: Failed to scan existing kode_rup: %v", err)
			continue
		}
		existingKodeRUPs[kodeRUP] = true
	}

	// Start transaction for robust error handling
	tx, err := database.Db.Begin()
	if err != nil {
		log.Printf("ERROR: Failed to start transaction: %v", err)
		return ScrapingResult{
			Success:      false,
			Message:      fmt.Sprintf("Failed to start transaction: %v", err),
			RecordsFound: 0,
			Endpoint:     endpoint,
		}, fmt.Errorf("transaction start failed: %v", err)
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("ERROR: Failed to rollback transaction: %v", rollbackErr)
			} else {
				log.Printf("Transaction rolled back due to error")
			}
		}
	}()

	recordsStored := 0
	recordsFailed := 0
	currentKodeRUPs := make(map[string]bool) // Track kode_rup values in current API response

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

		// Add active_year to the dataset
		orderedDataset.FieldValues["active_year"] = activeYear

		// Log mapping quality for debugging
		if orderedDataset.MappingStatus.MappedFields < orderedDataset.MappingStatus.TotalFields/2 {
			log.Printf("Warning: Poor field mapping for table %s - only %d/%d fields mapped",
				tableName, orderedDataset.MappingStatus.MappedFields, orderedDataset.MappingStatus.TotalFields)
		}

		// Track kode_rup from current API response
		if kodeRUP, exists := orderedDataset.FieldValues["kode_rup"]; exists && kodeRUP != "" {
			if kodeRUPStr, ok := kodeRUP.(string); ok && kodeRUPStr != "" {
				currentKodeRUPs[kodeRUPStr] = true
			}
		}

		// Store based on table type using ordered dataset
		var insertQuery string
		var args []interface{}
		var satkerQuery string
		var satkerArgs []interface{}
		switch tableName {
		case "perencanaan":
			insertQuery, args = ctrl.buildPerencanaanInsertFromDataset(orderedDataset)
			satkerQuery, satkerArgs = ctrl.buildSatuanKerjaInsertFromDataset(orderedDataset)
		case "persiapan":
			insertQuery, args = ctrl.buildPersiapanInsertFromDataset(orderedDataset)
		case "pemilihan":
			insertQuery, args = ctrl.buildPemilihanInsertFromDataset(orderedDataset)
		case "hasilpemilihan":
			insertQuery, args = ctrl.buildHasilPemilihanInsertFromDataset(orderedDataset)
		case "kontrak":
			insertQuery, args = ctrl.buildKontrakInsertFromDataset(orderedDataset)
		case "serahterima":
			insertQuery, args = ctrl.buildSerahTerimaInsertFromDataset(orderedDataset)
		}

		if insertQuery != "" {
			_, err := tx.Exec(insertQuery, args...)
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

		// Execute satker insert for perencanaan data if kode_satuan_kerja is present
		if satkerQuery != "" && len(satkerArgs) > 0 && satkerArgs[0] != nil && satkerArgs[0] != "" {
			_, err := tx.Exec(satkerQuery, satkerArgs...)
			if err != nil {
				log.Printf("Warning: Error inserting satker record: %v", err)
				// Don't count as failed since satker is secondary data
			}
		}
	}

	// Commit transaction if no errors occurred
	if err == nil {
		err = tx.Commit()
		if err != nil {
			log.Printf("ERROR: Failed to commit transaction: %v", err)
			return ScrapingResult{
				Success:      false,
				Message:      fmt.Sprintf("Failed to commit transaction: %v", err),
				RecordsFound: recordsStored,
				Endpoint:     endpoint,
			}, fmt.Errorf("transaction commit failed: %v", err)
		}
		log.Printf("Transaction committed successfully")
	}

	log.Printf("Stored %d records out of %d fetched (Failed: %d)", recordsStored, len(dataArray), recordsFailed)

	// Soft delete records that are no longer present in the API response
	recordsSoftDeleted := 0
	if len(existingKodeRUPs) > 0 {
		// Find kode_rup values that exist in DB but not in current API response
		var toDelete []string
		for existingKode := range existingKodeRUPs {
			if !currentKodeRUPs[existingKode] {
				toDelete = append(toDelete, existingKode)
			}
		}

		if len(toDelete) > 0 {
			// Build IN clause for soft deletion
			placeholders := make([]string, len(toDelete))
			args := make([]interface{}, len(toDelete))
			for i, kode := range toDelete {
				placeholders[i] = fmt.Sprintf("$%d", i+1)
				args[i] = kode
			}

			softDeleteQuery := fmt.Sprintf(
				"UPDATE %s SET deleted_at = NOW() WHERE kode_rup IN (%s) AND deleted_at IS NULL",
				fullTableName,
				strings.Join(placeholders, ","),
			)

			result, err := database.Exec(softDeleteQuery, args...)
			if err != nil {
				log.Printf("Warning: Failed to soft delete outdated records: %v", err)
			} else {
				rowsAffected, _ := result.RowsAffected()
				recordsSoftDeleted = int(rowsAffected)
				log.Printf("Soft deleted %d records that are no longer in API response", recordsSoftDeleted)
			}
		}
	}

	return ScrapingResult{
		Success:      true,
		Message:      fmt.Sprintf("Successfully stored %d records, soft deleted %d outdated records", recordsStored, recordsSoftDeleted),
		RecordsFound: recordsStored,
		Endpoint:     endpoint,
	}, nil
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

	endpoint := os.Getenv("SPSE_PERENCANAAN_ENDPOINT")

	result, err := ctrl.scrapeEndpoint(endpoint, "perencanaan")
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

	endpoint := os.Getenv("SPSE_PERSIAPAN_ENDPOINT")

	result, err := ctrl.scrapeEndpoint(endpoint, "persiapan")
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

	endpoint := os.Getenv("SPSE_PEMILIHAN_ENDPOINT")

	result, err := ctrl.scrapeEndpoint(endpoint, "pemilihan")
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

	endpoint := os.Getenv("SPSE_HASIL_PEMILIHAN_ENDPOINT")

	result, err := ctrl.scrapeEndpoint(endpoint, "hasilpemilihan")
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

	endpoint := os.Getenv("SPSE_KONTRAK_ENDPOINT")

	result, err := ctrl.scrapeEndpoint(endpoint, "kontrak")
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

	endpoint := os.Getenv("SPSE_SERAH_TERIMA_ENDPOINT")

	result, err := ctrl.scrapeEndpoint(endpoint, "serahterima")
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
		{os.Getenv("SPSE_PERENCANAAN_ENDPOINT"), "perencanaan", "Perencanaan"},
		{os.Getenv("SPSE_PERSIAPAN_ENDPOINT"), "persiapan", "Persiapan"},
		{os.Getenv("SPSE_PEMILIHAN_ENDPOINT"), "pemilihan", "Pemilihan"},
		{os.Getenv("SPSE_HASIL_PEMILIHAN_ENDPOINT"), "hasilpemilihan", "Hasil Pemilihan"},
		{os.Getenv("SPSE_KONTRAK_ENDPOINT"), "kontrak", "Kontrak"},
		{os.Getenv("SPSE_SERAH_TERIMA_ENDPOINT"), "serahterima", "Serah Terima"},
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
