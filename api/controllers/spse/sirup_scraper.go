package spse

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Massad/gin-boilerplate/db"
	"github.com/gin-gonic/gin"
)

// SIRUPScrapingResult represents the result of SIRUP scraping operations
type SIRUPScrapingResult struct {
	Success       bool   `json:"success"`
	Message       string `json:"message"`
	RecordsFound  int    `json:"records_found"`
	RecordsStored int    `json:"records_stored"`
	Endpoint      string `json:"endpoint"`
}

// scrapeSIRUPEndpoint scrapes SIRUP data for a specific kodeRUP
func (ctrl SPSEController) scrapeSIRUPEndpoint(kodeRUP string) (map[string]interface{}, error) {
	// Construct SIRUP URL
	sirupURL := fmt.Sprintf("https://sirup.inaproc.id/sirup/rup/detailPaketPenyedia2020?idPaket=%s", kodeRUP)

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Create request
	req, err := http.NewRequest("GET", sirupURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers to mimic browser
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Connection", "keep-alive")

	// Execute request with retry logic
	var resp *http.Response
	maxRetries := 3
	baseDelay := 2 * time.Second // Longer delay for SIRUP

	for attempt := 0; attempt < maxRetries; attempt++ {
		resp, err = client.Do(req)
		if err == nil && resp.StatusCode < 500 {
			break
		}

		if resp != nil {
			resp.Body.Close()
		}

		if attempt < maxRetries-1 {
			delay := time.Duration(attempt+1) * baseDelay
			log.Printf("WARNING: SIRUP request failed (attempt %d/%d), retrying in %v", attempt+1, maxRetries, delay)
			time.Sleep(delay)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to execute request after %d attempts: %v", maxRetries, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	// Read response body
	var body []byte
	var readErr error

	if resp.Header.Get("Content-Encoding") == "gzip" {
		reader, gzipErr := gzip.NewReader(resp.Body)
		if gzipErr != nil {
			return nil, fmt.Errorf("failed to create gzip reader: %v", gzipErr)
		}
		defer reader.Close()
		body, readErr = io.ReadAll(reader)
	} else {
		body, readErr = io.ReadAll(resp.Body)
	}

	if readErr != nil {
		return nil, fmt.Errorf("failed to read response: %v", readErr)
	}

	htmlContent := string(body)

	// Parse HTML content
	spseCtrl := SPSEController{}
	parsedData, err := spseCtrl.parseSIRUPHTML(htmlContent, kodeRUP)
	if err != nil {
		return nil, fmt.Errorf("failed to parse SIRUP HTML: %v", err)
	}

	return parsedData, nil
}

// ScrapeSIRUP godoc
// @Summary Scrape SIRUP data for filtered perencanaan records
// @Schemes
// @Description Scrapes SIRUP data from sirup.inaproc.id for filtered kodeRUP records from spse_perencanaan table
// @Tags SIRUP
// @Accept json
// @Produce json
// @Success 200 {object} SIRUPScrapingResult
// @Failure 500 {object} gin.H
// @Router /spse/sirup/scrape [POST]
func (ctrl SPSEController) ScrapeSIRUP(c *gin.Context) {
	log.Println("Starting filtered SIRUP scraping from perencanaan table...")

	spseCtrl := SPSEController{}

	// Get filtered kodeRUP values from perencanaan table
	kodeRUPs, err := spseCtrl.getFilteredPerencanaanRecords()
	if err != nil {
		log.Printf("Error getting filtered perencanaan records: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
			"message": "Failed to get filtered perencanaan records",
		})
		return
	}

	if len(kodeRUPs) == 0 {
		log.Println("No perencanaan records found to process")
		c.JSON(http.StatusOK, SIRUPScrapingResult{
			Success:       true,
			Message:       "No perencanaan records found to process",
			RecordsFound:  0,
			RecordsStored: 0,
			Endpoint:      "sirup.inaproc.id",
		})
		return
	}

	// Process each kodeRUP
	database := db.GetDB()
	if database == nil {
		log.Printf("ERROR: Database connection is nil!")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Database connection is nil",
			"message": "Database connection failed",
		})
		return
	}

	// Start transaction
	tx, err := database.Db.Begin()
	if err != nil {
		log.Printf("ERROR: Failed to start transaction: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   fmt.Sprintf("Failed to start transaction: %v", err),
			"message": "Transaction start failed",
		})
		return
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

	recordsProcessed := 0
	recordsStored := 0
	recordsFailed := 0
	totalRecords := len(kodeRUPs)

	// Log progress every 100 records or at start
	log.Printf("SIRUP scraping started: %d records to process", totalRecords)

	// Process each kodeRUP with rate limiting
	for i, kodeRUP := range kodeRUPs {
		// Scrape SIRUP data
		sirupData, err := ctrl.scrapeSIRUPEndpoint(kodeRUP)
		if err != nil {
			log.Printf("ERROR: Failed to scrape SIRUP for kodeRUP %s: %v", kodeRUP, err)
			// Store failed record with minimal data
			failedData := map[string]interface{}{
				"kode_rup":      kodeRUP,
				"sirup_scraped": false, // Mark as failed
			}

			// Prepare minimal data using existing function but override sirup_scraped
			enrichedFailedData, _ := spseCtrl.prepareSIRUPDataForInsertion(kodeRUP, failedData)
			enrichedFailedData["sirup_scraped"] = false // Ensure it's marked as failed

			// Add active_year from environment
			activeYear := os.Getenv("SPSE_ACTIVE_YEAR")
			enrichedFailedData["active_year"] = activeYear

			failedDataset := spseCtrl.convertMapToSIRUPDataset(enrichedFailedData)
			failedQuery, failedArgs := spseCtrl.buildSIRUPInsertFromDataset(failedDataset)
			if failedQuery != "" {
				if _, insertErr := tx.Exec(failedQuery, failedArgs...); insertErr != nil {
					log.Printf("ERROR: Failed to insert failed record for kodeRUP %s: %v", kodeRUP, insertErr)
				} else {
					recordsStored++ // Count as stored even though failed
				}
			}
			recordsFailed++
			recordsProcessed++
		} else {
			// Check if we actually got meaningful data
			isSuccessful := false
			if namaPaket, exists := sirupData["nama_paket"].(string); exists && namaPaket != "" {
				isSuccessful = true
			}

			if !isSuccessful {
				// Data extraction failed - treat as failed scraping
				log.Printf("WARNING: Data extraction failed for kodeRUP '%s' (empty nama_paket)", kodeRUP)
				failedData := map[string]interface{}{
					"kode_rup":      kodeRUP,
					"sirup_scraped": false, // Mark as failed due to no data
				}

				enrichedFailedData, _ := spseCtrl.prepareSIRUPDataForInsertion(kodeRUP, failedData)
				enrichedFailedData["sirup_scraped"] = false

				activeYear := os.Getenv("SPSE_ACTIVE_YEAR")
				enrichedFailedData["active_year"] = activeYear

				failedDataset := spseCtrl.convertMapToSIRUPDataset(enrichedFailedData)
				failedQuery, failedArgs := spseCtrl.buildSIRUPInsertFromDataset(failedDataset)
				if failedQuery != "" {
					if _, insertErr := tx.Exec(failedQuery, failedArgs...); insertErr != nil {
						log.Printf("ERROR: Failed to insert failed record for kodeRUP %s: %v", kodeRUP, insertErr)
					} else {
						recordsStored++
					}
				}
				recordsFailed++
				recordsProcessed++
			} else {
				// Successful scraping - prepare and store full data
				enrichedData, err := spseCtrl.prepareSIRUPDataForInsertion(kodeRUP, sirupData)
				if err != nil {
					log.Printf("ERROR: Failed to enrich data for kodeRUP %s: %v", kodeRUP, err)
					recordsFailed++
					recordsProcessed++
					continue
				}

				// Add active_year from environment
				activeYear := os.Getenv("SPSE_ACTIVE_YEAR")
				enrichedData["active_year"] = activeYear

				// Convert to ordered dataset
				orderedDataset := spseCtrl.convertMapToSIRUPDataset(enrichedData)

				// Build and execute insert query
				insertQuery, args := spseCtrl.buildSIRUPInsertFromDataset(orderedDataset)
				if insertQuery == "" {
					log.Printf("ERROR: No insert query generated for kodeRUP %s", kodeRUP)
					recordsFailed++
					recordsProcessed++
					continue
				}

				_, err = tx.Exec(insertQuery, args...)
				if err != nil {
					log.Printf("ERROR: Failed to insert SIRUP record for kodeRUP %s: %v", kodeRUP, err)
					recordsFailed++
					recordsProcessed++
					continue
				}

				recordsStored++
				recordsProcessed++
			}
		}

		// Log progress every 100 records
		if (i+1)%100 == 0 || i == 0 {
			log.Printf("SIRUP scraping progress: %d/%d records processed", i+1, totalRecords)
		}

		// Rate limiting - add delay between requests
		if i < len(kodeRUPs)-1 {
			time.Sleep(1 * time.Second)
		}
	}

	// Commit transaction
	if err == nil {
		err = tx.Commit()
		if err != nil {
			log.Printf("ERROR: Failed to commit transaction: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   fmt.Sprintf("Failed to commit transaction: %v", err),
				"message": "Transaction commit failed",
			})
			return
		}
		log.Printf("Transaction committed successfully")
	}

	log.Printf("SIRUP scraping completed: %d processed, %d stored, %d failed", recordsProcessed, recordsStored, recordsFailed)

	result := SIRUPScrapingResult{
		Success:       true,
		Message:       fmt.Sprintf("Processed %d records: %d successful, %d failed (all stored)", recordsProcessed, recordsStored-recordsFailed, recordsFailed),
		RecordsFound:  recordsProcessed,
		RecordsStored: recordsStored,
		Endpoint:      "sirup.inaproc.id",
	}

	c.JSON(http.StatusOK, result)
}

// ScrapeSIRUPSingle godoc
// @Summary Scrape SIRUP data for a specific kodeRUP
// @Schemes
// @Description Scrapes SIRUP data from sirup.inaproc.id for a single kodeRUP
// @Tags SIRUP
// @Accept json
// @Produce json
// @Param kodeRUP query string true "Kode RUP to scrape"
// @Success 200 {object} SIRUPScrapingResult
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /spse/sirup/scrape/single [POST]
func (ctrl SPSEController) ScrapeSIRUPSingle(c *gin.Context) {
	kodeRUP := c.Query("kodeRUP")
	if kodeRUP == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "kodeRUP parameter is required",
			"message": "Missing kodeRUP parameter",
		})
		return
	}

	log.Printf("Starting single SIRUP scraping for kodeRUP: %s", kodeRUP)

	// Scrape SIRUP data
	sirupData, err := ctrl.scrapeSIRUPEndpoint(kodeRUP)
	if err != nil {
		log.Printf("Error scraping SIRUP for kodeRUP %s: %v", kodeRUP, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
			"message": "Failed to scrape SIRUP data",
		})
		return
	}

	spseCtrl := SPSEController{}

	// Prepare SIRUP data for insertion
	enrichedData, err := spseCtrl.prepareSIRUPDataForInsertion(kodeRUP, sirupData)
	if err != nil {
		log.Printf("Error enriching data for kodeRUP %s: %v", kodeRUP, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
			"message": "Failed to enrich perencanaan data",
		})
		return
	}

	// Add active_year from environment
	activeYear := os.Getenv("SPSE_ACTIVE_YEAR")
	enrichedData["active_year"] = activeYear

	// Store in database
	database := db.GetDB()
	if database == nil {
		log.Printf("ERROR: Database connection is nil!")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Database connection is nil",
			"message": "Database connection failed",
		})
		return
	}

	// Convert to ordered dataset
	orderedDataset := spseCtrl.convertMapToSIRUPDataset(enrichedData)

	// Build and execute insert query
	insertQuery, args := spseCtrl.buildSIRUPInsertFromDataset(orderedDataset)
	if insertQuery == "" {
		log.Printf("No insert query generated for kodeRUP %s", kodeRUP)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "No insert query generated",
			"message": "Failed to generate insert query",
		})
		return
	}

	_, err = database.Exec(insertQuery, args...)
	if err != nil {
		log.Printf("Error inserting SIRUP record for kodeRUP %s: %v", kodeRUP, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
			"message": "Failed to store SIRUP data",
		})
		return
	}

	log.Printf("Successfully stored SIRUP data for kodeRUP: %s", kodeRUP)

	result := SIRUPScrapingResult{
		Success:       true,
		Message:       fmt.Sprintf("Successfully stored SIRUP data for kodeRUP %s", kodeRUP),
		RecordsFound:  1,
		RecordsStored: 1,
		Endpoint:      "sirup.inaproc.id",
	}

	c.JSON(http.StatusOK, result)
}
