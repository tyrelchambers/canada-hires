# Job Bank Scraper (Go Version)

A simplified Go-based web scraper for Canadian Job Bank LMIA (Labour Market Impact Assessment) jobs.

## Features

- Scrapes LMIA jobs from jobbank.gc.ca
- Browser automation using chromedp (Chrome DevTools Protocol)
- HTML parsing with goquery
- Data cleaning and province name normalization
- Statistical reporting (job counts by location and employer)
- Outputs structured data ready for database insertion

## Setup

1. Install dependencies:
```bash
go mod tidy
```

2. Build the scraper:
```bash
go build -o scraper main.go
```

3. Run the scraper:
```bash
./scraper
```

## Architecture

- `main.go` - Entry point with statistical reporting
- `types/` - Data structures (JobData)
- `scraper/` - Browser automation and scraping logic

## Key Differences from TypeScript Version

- Uses chromedp instead of Puppeteer for browser automation
- Uses goquery instead of Cheerio for HTML parsing
- Simplified - no API client or config management
- Focused on data extraction with helpful statistics
- Ready for direct database integration within the main API

## Usage

The scraper:
- Scrapes all LMIA jobs (no filtering)
- Processes all available pages (-1 means unlimited)
- Provides helpful statistics about scraped data
- Returns clean, structured job data ready for database storage