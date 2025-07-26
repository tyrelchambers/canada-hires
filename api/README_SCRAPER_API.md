# Job Scraper API Integration

This document explains how to integrate your job scraper with the Canada Hires API.

## Database Migration

First, run the API server to apply the new database migration:

```bash
cd api
go run main.go
```

The migration `011_update_job_postings_for_scraper` will automatically run and update the database schema.

## API Endpoints

### 1. Start a Scraping Session

**POST** `/api/jobs/scraping-runs`

Creates a new scraping session and returns a `scraping_run_id` to use for this batch.

```bash
curl -X POST http://localhost:8000/api/jobs/scraping-runs
```

Response:
```json
{
  "id": "uuid-string",
  "status": "running",
  "started_at": "2025-07-26T...",
  "total_pages": 0,
  "jobs_scraped": 0,
  "jobs_stored": 0,
  "last_page_scraped": 0,
  "created_at": "2025-07-26T..."
}
```

### 2. Submit Job Data

**POST** `/api/jobs/scraping-runs/{scraping_run_id}/jobs`

Submit your scraped job data using the structure from your scraper:

```bash
curl -X POST http://localhost:8000/api/jobs/scraping-runs/{scraping_run_id}/jobs \
  -H "Content-Type: application/json" \
  -d '[
    {
      "jobTitle": "cook",
      "business": "Le Bistro Montebello", 
      "salary": "$20.00 to $25.00 hourly",
      "location": "Montebello Quebec",
      "jobUrl": "https://www.jobbank.gc.ca/jobsearch/jobpostingtfw/44736629",
      "date": "July 25, 2025"
    }
  ]'
```

Response:
```json
{
  "message": "Jobs successfully stored",
  "jobs_processed": 1,
  "scraping_run_id": "uuid-string"
}
```

### 3. Complete Scraping Session

**POST** `/api/jobs/scraping-runs/{scraping_run_id}/complete`

Mark the scraping session as completed:

```bash
curl -X POST http://localhost:8000/api/jobs/scraping-runs/{scraping_run_id}/complete \
  -H "Content-Type: application/json" \
  -d '{
    "total_pages": 10,
    "jobs_scraped": 100,
    "jobs_stored": 95
  }'
```

## Query Endpoints

### Get Job Postings

**GET** `/api/jobs/`

Query parameters:
- `employer` - Filter by employer name
- `city` - Filter by city
- `province` - Filter by province  
- `limit` - Number of results (default: 20)

```bash
curl "http://localhost:8000/api/jobs/?employer=bistro&limit=10"
```

### Get Job Statistics

**GET** `/api/jobs/stats`

```bash
curl http://localhost:8000/api/jobs/stats
```

Response:
```json
{
  "total_jobs": 1500,
  "total_employers": 450,
  "top_employers": [
    {
      "employer": "Restaurant Group Inc",
      "job_count": 25,
      "earliest_posting": "2025-07-01T...",
      "latest_posting": "2025-07-26T..."
    }
  ]
}
```

## Data Processing

The API automatically processes your scraper data:

1. **Salary Parsing**: Converts "$20.00 to $25.00 hourly" into structured min/max salary and type
2. **Location Parsing**: Splits "Montebello Quebec" into city and province fields
3. **Date Parsing**: Converts "July 25, 2025" into a proper timestamp
4. **URL Deduplication**: Uses the job URL as the unique identifier to prevent duplicates

## Integration Example

Here's a simple Node.js example of how to integrate your scraper:

```javascript
const axios = require('axios');

async function submitJobsToAPI(jobsData) {
  // 1. Start scraping session
  const sessionResponse = await axios.post('http://localhost:8000/api/jobs/scraping-runs');
  const scrapingRunId = sessionResponse.data.id;
  
  // 2. Submit job data
  await axios.post(`http://localhost:8000/api/jobs/scraping-runs/${scrapingRunId}/jobs`, jobsData);
  
  // 3. Complete session  
  await axios.post(`http://localhost:8000/api/jobs/scraping-runs/${scrapingRunId}/complete`, {
    total_pages: 1,
    jobs_scraped: jobsData.length,
    jobs_stored: jobsData.length
  });
  
  console.log(`Successfully submitted ${jobsData.length} jobs`);
}

// Your existing scraper data
const scrapedJobs = [
  {
    jobTitle: "cook",
    business: "Le Bistro Montebello",
    salary: "$20.00 to $25.00 hourly",
    location: "Montebello Quebec", 
    jobUrl: "https://www.jobbank.gc.ca/jobsearch/jobpostingtfw/44736629",
    date: "July 25, 2025"
  }
  // ... more jobs
];

submitJobsToAPI(scrapedJobs);
```

## Benefits

- **Automatic Data Processing**: No need to manually parse salaries, locations, or dates
- **Deduplication**: Prevents duplicate job postings based on URL
- **Structured Storage**: Jobs are stored in a searchable, normalized format
- **API Access**: Other parts of your application can easily query the job data
- **Statistics**: Get insights about employers and job trends