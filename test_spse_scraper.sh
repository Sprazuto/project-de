#!/bin/bash
# SPSE Scraper Test Script
# This script demonstrates how to test all the SPSE scraper endpoints

BASE_URL="http://localhost:9000/v1/spse"
echo "🧪 SPSE Procurement Data Scraper Test Suite"
echo "============================================"
echo ""

# Test 1: Health check
echo "1. Testing server health..."
if curl -s "${BASE_URL%/*}/health" > /dev/null 2>&1; then
    echo "✅ Server is running"
else
    echo "❌ Server is not running. Please start the server first."
    exit 1
fi
echo ""

# Test 2: Get initial statistics
echo "2. Getting initial statistics..."
curl -s -X GET "${BASE_URL}/statistics" | jq '.' || curl -s -X GET "${BASE_URL}/statistics"
echo ""

# Test 3: Test individual scraping endpoints
echo "3. Testing Perencanaan scraping..."
echo "POST ${BASE_URL}/scraper/perencanaan"
curl -s -X POST "${BASE_URL}/scraper/perencanaan" | jq '.' || curl -s -X POST "${BASE_URL}/scraper/perencanaan"
echo ""

echo "4. Testing Persiapan scraping..."
echo "POST ${BASE_URL}/scraper/persiapan"
curl -s -X POST "${BASE_URL}/scraper/persiapan" | jq '.' || curl -s -X POST "${BASE_URL}/scraper/persiapan"
echo ""

echo "5. Testing Pemilihan scraping..."
echo "POST ${BASE_URL}/scraper/pemilihan"
curl -s -X POST "${BASE_URL}/scraper/pemilihan" | jq '.' || curl -s -X POST "${BASE_URL}/scraper/pemilihan"
echo ""

echo "6. Testing Hasil Pemilihan scraping..."
echo "POST ${BASE_URL}/scraper/hasilpemilihan"
curl -s -X POST "${BASE_URL}/scraper/hasilpemilihan" | jq '.' || curl -s -X POST "${BASE_URL}/scraper/hasilpemilihan"
echo ""

echo "7. Testing Kontrak scraping..."
echo "POST ${BASE_URL}/scraper/kontrak"
curl -s -X POST "${BASE_URL}/scraper/kontrak" | jq '.' || curl -s -X POST "${BASE_URL}/scraper/kontrak"
echo ""

echo "8. Testing Serah Terima scraping..."
echo "POST ${BASE_URL}/scraper/serahterima"
curl -s -X POST "${BASE_URL}/scraper/serahterima" | jq '.' || curl -s -X POST "${BASE_URL}/scraper/serahterima"
echo ""

# Test 9: Test comprehensive scraping
echo "9. Testing comprehensive scraping (all endpoints)..."
echo "POST ${BASE_URL}/scraper/all"
curl -s -X POST "${BASE_URL}/scraper/all" | jq '.' || curl -s -X POST "${BASE_URL}/scraper/all"
echo ""

# Test 10: Get updated statistics
echo "10. Getting updated statistics after scraping..."
curl -s -X GET "${BASE_URL}/statistics" | jq '.' || curl -s -X GET "${BASE_URL}/statistics"
echo ""

# Test 11: Test data retrieval with pagination
echo "11. Testing data retrieval with pagination..."
echo "GET ${BASE_URL}/data/perencanaan?limit=5&offset=0"
curl -s -X GET "${BASE_URL}/data/perencanaan?limit=5&offset=0" | jq '.' || curl -s -X GET "${BASE_URL}/data/perencanaan?limit=5&offset=0"
echo ""

echo "12. Testing Persiapan data retrieval..."
echo "GET ${BASE_URL}/data/persiapan?limit=5&offset=0"
curl -s -X GET "${BASE_URL}/data/persiapan?limit=5&offset=0" | jq '.' || curl -s -X GET "${BASE_URL}/data/persiapan?limit=5&offset=0"
echo ""

echo "13. Testing Pemilihan data retrieval..."
echo "GET ${BASE_URL}/data/pemilihan?limit=5&offset=0"
curl -s -X GET "${BASE_URL}/data/pemilihan?limit=5&offset=0" | jq '.' || curl -s -X GET "${BASE_URL}/data/pemilihan?limit=5&offset=0"
echo ""

echo "14. Testing Hasil Pemilihan data retrieval..."
echo "GET ${BASE_URL}/data/hasilpemilihan?limit=5&offset=0"
curl -s -X GET "${BASE_URL}/data/hasilpemilihan?limit=5&offset=0" | jq '.' || curl -s -X GET "${BASE_URL}/data/hasilpemilihan?limit=5&offset=0"
echo ""

echo "15. Testing Kontrak data retrieval..."
echo "GET ${BASE_URL}/data/kontrak?limit=5&offset=0"
curl -s -X GET "${BASE_URL}/data/kontrak?limit=5&offset=0" | jq '.' || curl -s -X GET "${BASE_URL}/data/kontrak?limit=5&offset=0"
echo ""

echo "16. Testing Serah Terima data retrieval..."
echo "GET ${BASE_URL}/data/serahterima?limit=5&offset=0"
curl -s -X GET "${BASE_URL}/data/serahterima?limit=5&offset=0" | jq '.' || curl -s -X GET "${BASE_URL}/data/serahterima?limit=5&offset=0"
echo ""

echo "✅ Test suite completed!"
echo "📊 Check the logs for detailed scraping information"
echo "🔍 All endpoints are accessible at: ${BASE_URL}"