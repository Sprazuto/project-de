# SPSE Procurement Data Scraper

A minimal, fast web scraper for extracting procurement data from the SPSE (Sistem Pengadaan Secara Elektronik) API endpoints. This scraper is integrated into the existing Gin-based Go backend.

## Overview

This scraper fetches procurement data from three different stages of the procurement process:

1. **Perencanaan (Planning)** - Procurement planning data
2. **Persiapan (Preparation)** - Package preparation data
3. **Pemilihan (Selection)** - Provider selection and contracting data

## Features

- ✅ Automatic authentication token extraction
- ✅ Dynamic database schema creation
- ✅ Multi-endpoint data scraping
- ✅ PostgreSQL integration with Gorp ORM
- ✅ RESTful API endpoints
- ✅ JSON response formatting
- ✅ Error handling and logging
- ✅ Data deduplication via unique constraints

## API Endpoints

### Scraping Endpoints (Public Access)

```
POST /v1/spse/scraper/perencanaan
POST /v1/spse/scraper/persiapan
POST /v1/spse/scraper/pemilihan
POST /v1/spse/scraper/all
```

### Data Retrieval Endpoints (Protected)

```
GET  /v1/spse/statistics
GET  /v1/spse/data/perencanaan
GET  /v1/spse/data/persiapan
GET  /v1/spse/data/pemilihan
```

## Database Schema

The scraper creates three main tables:

### spse_perencanaan

- `id` (Primary Key)
- `kode_rup` - RUP Code
- `satuan_kerja` - Work Unit
- `nama_paket` - Package Name
- `metode_pemilihan` - Selection Method
- `tanggal_pengumuman` - Announcement Date
- `rencana_pemilihan` - Selection Plan
- `pagu_rup` - RUP Budget
- Additional fields for procurement details

### spse_persiapan

- `id` (Primary Key)
- `kode_rup` - RUP Code
- `satuan_kerja` - Work Unit
- `nama_paket` - Package Name
- `metode_pemilihan` - Selection Method
- `tanggal_buat_paket` - Package Creation Date
- `nilai_pagu_rup` - RUP Budget Value
- `nilai_pagu_paket` - Package Budget Value

### spse_pemilihan

- `id` (Primary Key)
- `kode_rup` - RUP Code
- `satuan_kerja` - Work Unit
- `nama_paket` - Package Name
- `metode_pemilihan` - Selection Method
- `rencana_pemilihan` - Selection Plan
- `tanggal_pemilihan` - Selection Date
- `nilai_hps` - HPS Value
- `status_paket` - Package Status
- Additional fields for contract details

## Installation & Setup

### Prerequisites

- Go 1.18+
- PostgreSQL database
- Existing Gin boilerplate project

### 1. Database Setup

Ensure your PostgreSQL database is configured with the following environment variables:

```env
DATABASE_USER=your_username
DATABASE_PASSWORD=your_password
DATABASE_NAME=your_database
DATABASE_HOST=localhost
DATABASE_PORT=5432
```

### 2. Run the Application

```bash
cd api
go run main.go
```

The server will start on the configured PORT (default: 9000).

### 3. Test the Endpoints

```bash
# Test scraping all endpoints
curl -X POST http://localhost:9000/v1/spse/scraper/all

# Get scraping statistics
curl -X GET http://localhost:9000/v1/spse/statistics

# Retrieve perencanaan data with pagination
curl -X GET "http://localhost:9000/v1/spse/data/perencanaan?limit=50&offset=0"
```

## Usage Examples

### Scrape Single Endpoint

```bash
curl -X POST http://localhost:9000/v1/spse/scraper/perencanaan
```

Response:

```json
{
  "success": true,
  "message": "Successfully scraped 150 records from /sumedangkab/amel/dt/detailperencanaan2",
  "records_found": 150,
  "endpoint": "/sumedangkab/amel/dt/detailperencanaan2"
}
```

### Scrape All Endpoints

```bash
curl -X POST http://localhost:9000/v1/spse/scraper/all
```

Response:

```json
{
  "results": {
    "Perencanaan": {
      "success": true,
      "message": "Successfully scraped 150 records",
      "records_found": 150,
      "endpoint": "/sumedangkab/amel/dt/detailperencanaan2"
    },
    "Persiapan": {
      "success": true,
      "message": "Successfully scraped 120 records",
      "records_found": 120,
      "endpoint": "/sumedangkab/amel/dt/detailpersiapan2"
    },
    "Pemilihan": {
      "success": true,
      "message": "Successfully scraped 80 records",
      "records_found": 80,
      "endpoint": "/sumedangkab/amel/dt/detailpemilihan2"
    }
  },
  "total_success": 3,
  "total_endpoints": 3,
  "message": "Scraping completed: 3/3 endpoints successful"
}
```

### Get Data with Pagination

```bash
curl -X GET "http://localhost:9000/v1/spse/data/perencanaan?limit=10&offset=0" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

Response:

```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "kode_rup": "123456789",
      "satuan_kerja": "Dinas Pekerjaan Umum",
      "nama_paket": "Pembangunan Jalan Raya",
      "metode_pemilihan": "e-Tendering",
      "tanggal_pengumuman": "2025-01-15",
      "rencana_pemilihan": "2025-02-01",
      "pagu_rup": "1,000,000,000",
      "created_at": "2025-01-20T10:30:00Z"
    }
  ],
  "pagination": {
    "limit": 10,
    "offset": 0,
    "count": 10
  }
}
```

## Integration with Frontend

The scraper provides RESTful endpoints that can be consumed by the Vue.js frontend:

```javascript
// Example Vue.js service
class SPSEService {
  async scrapeAll() {
    const response = await fetch("/v1/spse/scraper/all", {
      method: "POST",
    });
    return response.json();
  }

  async getPerencanaan(limit = 100, offset = 0) {
    const response = await fetch(
      `/v1/spse/data/perencanaan?limit=${limit}&offset=${offset}`,
      {
        headers: {
          Authorization: `Bearer ${localStorage.getItem("token")}`,
        },
      }
    );
    return response.json();
  }

  async getStatistics() {
    const response = await fetch("/v1/spse/statistics", {
      headers: {
        Authorization: `Bearer ${localStorage.getItem("token")}`,
      },
    });
    return response.json();
  }
}
```

## Error Handling

The scraper includes comprehensive error handling:

- **Authentication Token Issues**: Automatic token extraction and regeneration
- **Network Timeouts**: Configurable HTTP client timeouts
- **Database Errors**: Transaction rollback and error logging
- **Data Validation**: Input sanitization and format validation

## Logging

The scraper logs all operations to:

- Console output for development
- `spse_scraper.log` file for persistent logging

Log levels:

- **INFO**: Successful operations
- **WARNING**: Non-critical issues (e.g., missing optional data)
- **ERROR**: Critical failures that require attention

## Performance Considerations

- **Connection Pooling**: Reuses HTTP connections for efficiency
- **Batch Processing**: Processes multiple records in single transactions
- **Pagination Support**: Large datasets retrieved in manageable chunks
- **Indexing**: Unique constraints on key fields for fast lookups

## Security Features

- **CSRF Protection**: Extracts and uses authenticity tokens
- **Request Rate Limiting**: Built into the Gin framework
- **Input Validation**: Sanitizes all input data
- **SQL Injection Prevention**: Uses parameterized queries

## Troubleshooting

### Common Issues

1. **"Authenticity token not found"**

   - Check if the SPSE website structure has changed
   - Verify network connectivity to spse.inaproc.id

2. **"Failed to create tables"**

   - Ensure PostgreSQL is running and accessible
   - Check database user permissions
   - Verify database connection configuration

3. **"Invalid response format"**
   - The API may have changed response structure
   - Check logs for detailed error information

### Debug Mode

Enable debug logging by setting:

```go
log.SetLevel(log.DebugLevel)
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is part of the SPSE procurement system integration and follows the same licensing terms.

## Support

For issues and questions:

- Check the logs for detailed error messages
- Review the API endpoint documentation
- Test individual endpoints using curl commands
- Verify database connectivity and permissions
