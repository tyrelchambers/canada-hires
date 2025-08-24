package scraper

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	scraper_types "canada-hires/scraper-types"

	"github.com/chromedp/chromedp"
)

const (
	baseURL             = "https://www.jobbank.gc.ca"
	lmiaURL             = "https://www.jobbank.gc.ca/jobsearch/jobsearch?fsrc=32"
	nonCompliantURL     = "https://www.canada.ca/en/immigration-refugees-citizenship/services/work-canada/employers-non-compliant.html"
)

type Scraper struct {
	ctx     context.Context
	cancel  context.CancelFunc
	timeout time.Duration
}

func NewScraper() (*Scraper, error) {
	// Create chromedp context with options
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
	)

	allocCtx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancel := chromedp.NewContext(allocCtx)

	// Test if chromedp is working
	if err := chromedp.Run(ctx); err != nil {
		cancel()
		return nil, fmt.Errorf("could not start chromedp: %v", err)
	}

	return &Scraper{
		ctx:     ctx,
		cancel:  cancel,
		timeout: time.Duration(1000) * time.Millisecond,
	}, nil
}

func (s *Scraper) Close() {
	if s.cancel != nil {
		s.cancel()
	}
}

func (s *Scraper) ScrapeLMIAJobs(numberOfPages int) ([]scraper_types.JobData, error) {
	fmt.Println("🎯 Navigating directly to LMIA jobs page...")

	err := chromedp.Run(s.ctx,
		chromedp.Navigate(lmiaURL),
		chromedp.WaitVisible("#moreresultbutton", chromedp.ByID),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to navigate to LMIA page: %v", err)
	}

	// Get total results count
	var totalResults string
	err = chromedp.Run(s.ctx,
		chromedp.Text("#results-count", &totalResults),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get results count: %v", err)
	}

	fmt.Printf("\n📊 Total LMIA jobs to scrape: %s\n", strings.TrimSpace(totalResults))

	// Load more pages
	if err := s.loadMorePages(numberOfPages); err != nil {
		return nil, fmt.Errorf("failed to load more pages: %v", err)
	}

	// Parse jobs from the page
	fmt.Println("🔍 Starting to parse jobs from loaded pages...")
	jobs, err := s.parseJobs()
	if err != nil {
		return nil, fmt.Errorf("failed to parse jobs: %v", err)
	}
	fmt.Printf("✅ Finished parsing %d jobs\n", len(jobs))

	// Clean the data
	fmt.Println("🧹 Starting data cleaning...")
	cleanedJobs := s.cleanData(jobs)
	fmt.Printf("✅ Finished cleaning %d jobs\n", len(cleanedJobs))

	return cleanedJobs, nil
}

func (s *Scraper) loadMorePages(numberOfPages int) error {
	scrapeAll := numberOfPages == -1
	i := 0

	for {
		if !scrapeAll && i >= numberOfPages {
			break
		}

		time.Sleep(s.timeout)

		// Check if more button exists and is visible
		var buttonExists bool
		err := chromedp.Run(s.ctx,
			chromedp.Evaluate(`document.querySelector('#moreresultbutton') !== null`, &buttonExists),
		)
		if err != nil {
			return fmt.Errorf("failed to check for more button: %v", err)
		}

		if buttonExists {
			// Check if button is visible
			var isVisible bool
			err := chromedp.Run(s.ctx,
				chromedp.Evaluate(`(function() {
					const btn = document.querySelector('#moreresultbutton');
					return btn && btn.offsetParent !== null;
				})()`, &isVisible),
			)
			if err != nil {
				return fmt.Errorf("failed to check button visibility: %v", err)
			}

			if isVisible {
				err = chromedp.Run(s.ctx,
					chromedp.Click("#moreresultbutton", chromedp.ByID),
					chromedp.Sleep(s.timeout),
				)
				if err != nil {
					return fmt.Errorf("failed to click more button: %v", err)
				}

				i++
				if scrapeAll {
					fmt.Printf("%d 📄(s) loaded (scraping all pages...)\n", i)
				} else {
					fmt.Printf("%d 📄(s) loaded out of %d\n", i, numberOfPages)
				}

				time.Sleep(s.timeout)
			} else {
				fmt.Printf("More button not visible after %d pages 😔\n", i)
				fmt.Println("Finished loading all available pages")
				break
			}
		} else {
			fmt.Printf("No more results after %d pages 😔\n", i)
			fmt.Println("Finished loading all available pages")
			time.Sleep(s.timeout * 7)
			break
		}
	}

	return nil
}

func (s *Scraper) parseJobs() ([]scraper_types.JobData, error) {
	// Get all job articles
	var articlesHTML string
	err := chromedp.Run(s.ctx,
		chromedp.Evaluate(`(function() {
			const articles = document.querySelectorAll('article');
			let result = [];
			for (let i = 0; i < articles.length; i++) {
				const article = articles[i];
				const titleEl = article.querySelector('.noctitle');
				const businessEl = article.querySelector('.list-unstyled .business');
				const locationEl = article.querySelector('.list-unstyled .location');
				const salaryEl = article.querySelector('.list-unstyled .salary');
				const dateEl = article.querySelector('.list-unstyled .date');
				const linkEl = article.querySelector('.resultJobItem');

				result.push({
					title: titleEl ? titleEl.textContent.trim() : '',
					business: businessEl ? businessEl.textContent.trim() : '',
					location: locationEl ? locationEl.textContent.trim() : '',
					salary: salaryEl ? salaryEl.textContent.trim() : '',
					date: dateEl ? dateEl.textContent.trim() : '',
					href: linkEl ? linkEl.getAttribute('href') : ''
				});
			}
			return JSON.stringify(result);
		})()`, &articlesHTML),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get articles: %v", err)
	}

	// Parse the JSON response
	var articles []struct {
		Title    string `json:"title"`
		Business string `json:"business"`
		Location string `json:"location"`
		Salary   string `json:"salary"`
		Date     string `json:"date"`
		Href     string `json:"href"`
	}

	if err := json.Unmarshal([]byte(articlesHTML), &articles); err != nil {
		return nil, fmt.Errorf("failed to parse articles JSON: %v", err)
	}

	fmt.Printf("🔍 Found %d article elements on page\n", len(articles))

	var jobs []scraper_types.JobData

	for i, article := range articles {
		jobTitle := strings.TrimSpace(article.Title)
		business := strings.TrimSpace(article.Business)
		location := strings.TrimSpace(article.Location)
		salary := strings.TrimSpace(article.Salary)
		date := strings.TrimSpace(article.Date)
		href := article.Href

		if href == "" {
			continue
		}
		rawJobURL := baseURL + href

		// Clean URL and extract job bank ID
		jobURL, jobBankID := cleanJobURL(rawJobURL)

		if jobTitle != "" && jobURL != "" {
			jobs = append(jobs, scraper_types.JobData{
				JobTitle:  jobTitle,
				Business:  business,
				Salary:    salary,
				Location:  location,
				JobURL:    jobURL,
				Date:      date,
				JobBankID: jobBankID,
			})

		} else {
			if i < 5 {
				fmt.Printf("Skipping article %d: title='%s', url='%s'\n", i+1, jobTitle, jobURL)
			}
		}
	}

	fmt.Printf("Scraped %d jobs from the page\n", len(jobs))
	return jobs, nil
}

func cleanJobURL(url string) (cleanURL string, jobBankID string) {
	// Extract job bank ID from URL (handles both /jobpostingtfw/ and other formats)
	re := regexp.MustCompile(`/jobpostingtfw/(\d+)`)
	matches := re.FindStringSubmatch(url)
	if len(matches) > 1 {
		jobBankID = matches[1]
		// Clean URL by removing everything after the job ID
		cleanURL = fmt.Sprintf("%s/jobsearch/jobpostingtfw/%s", baseURL, jobBankID)
		return cleanURL, jobBankID
	}

	// If no job ID found in expected format, try other patterns
	re = regexp.MustCompile(`/(\d+)`)
	matches = re.FindStringSubmatch(url)
	if len(matches) > 1 {
		jobBankID = matches[len(matches)-1] // Get the last number found
		// For other URL formats, still try to clean session IDs and query params
		re = regexp.MustCompile(`(.*?)(;jsessionid=.*|\?.*)`)
		cleanMatches := re.FindStringSubmatch(url)
		if len(cleanMatches) > 1 {
			cleanURL = cleanMatches[1]
		} else {
			cleanURL = url
		}
		return cleanURL, jobBankID
	}

	// Fallback: just clean session IDs and query params
	re = regexp.MustCompile(`(.*?)(;jsessionid=.*|\?.*)`)
	matches = re.FindStringSubmatch(url)
	if len(matches) > 1 {
		cleanURL = matches[1]
	} else {
		cleanURL = url
	}
	return cleanURL, ""
}

// Deprecated: use cleanJobURL instead
func extractJobBankID(url string) string {
	_, jobBankID := cleanJobURL(url)
	return jobBankID
}

func (s *Scraper) cleanData(jobs []scraper_types.JobData) []scraper_types.JobData {
	for i := range jobs {
		jobs[i].JobTitle = removeTabsAndNewLines(jobs[i].JobTitle)
		jobs[i].Business = removeTabsAndNewLines(jobs[i].Business)
		jobs[i].Salary = removeTabsAndNewLines(jobs[i].Salary)
		jobs[i].Location = removeTabsAndNewLines(jobs[i].Location)
		jobs[i].Date = removeTabsAndNewLines(jobs[i].Date)
	}
	return jobs
}

func removeTabsAndNewLines(str string) string {
	// Remove tabs, newlines, and various labels
	re := regexp.MustCompile(`(\t|\n|Location)`)
	str = re.ReplaceAllString(str, "")

	// Remove leading whitespace
	re = regexp.MustCompile(`(^\s*)`)
	str = re.ReplaceAllString(str, "")

	// Remove salary label and negotiation text
	re = regexp.MustCompile(`(Salary:|to be negotiated)`)
	str = re.ReplaceAllString(str, "")

	// Remove parentheses
	re = regexp.MustCompile(`(\(|\))`)
	str = re.ReplaceAllString(str, "")

	// Replace province abbreviations with full names
	provinceMap := map[string]string{
		"BC": "British Columbia",
		"ON": "Ontario",
		"QC": "Quebec",
		"SK": "Saskatchewan",
		"AB": "Alberta",
		"MB": "Manitoba",
		"NB": "New Brunswick",
		"NL": "Newfoundland and Labrador",
		"NS": "Nova Scotia",
		"PE": "Prince Edward Island",
		"NT": "Northwest Territories",
		"NU": "Nunavut",
		"YT": "Yukon",
	}

	for abbrev, full := range provinceMap {
		re = regexp.MustCompile(`\b` + abbrev + `\b`)
		str = re.ReplaceAllString(str, full)
	}

	return strings.TrimSpace(str)
}

// ScrapeReasonDescriptions scrapes all reason code descriptions from the non-compliant page
func (s *Scraper) ScrapeReasonDescriptions() (map[string]string, error) {
	fmt.Println("🎯 Scraping reason descriptions from non-compliant page...")
	
	err := chromedp.Run(s.ctx,
		chromedp.Navigate(nonCompliantURL),
		chromedp.WaitVisible("ol li", chromedp.ByQuery),
		chromedp.Sleep(2*time.Second), // Wait for page to fully load
	)
	if err != nil {
		return nil, fmt.Errorf("failed to navigate to non-compliant page: %v", err)
	}

	// Extract reason descriptions using JavaScript
	var reasonsJSON string
	err = chromedp.Run(s.ctx,
		chromedp.Evaluate(`(function() {
			const reasonsMap = {};
			
			// Look for list items with ids like "list1", "list2", etc.
			for (let i = 1; i <= 50; i++) {
				const element = document.getElementById('list' + i);
				if (element) {
					const description = element.textContent.trim();
					if (description) {
						reasonsMap[i.toString()] = description;
					}
				}
			}
			
			// If no elements with list IDs found, try to extract from ordered list
			if (Object.keys(reasonsMap).length === 0) {
				const listItems = document.querySelectorAll('ol li');
				for (let i = 0; i < listItems.length; i++) {
					const item = listItems[i];
					const description = item.textContent.trim();
					if (description) {
						// Use 1-based indexing to match reason codes
						reasonsMap[(i + 1).toString()] = description;
					}
				}
			}
			
			// If still no results, try to find any numbered list pattern
			if (Object.keys(reasonsMap).length === 0) {
				const allLists = document.querySelectorAll('ol, ul');
				for (let list of allLists) {
					const items = list.querySelectorAll('li');
					if (items.length > 10) { // Likely the main reasons list
						for (let i = 0; i < items.length; i++) {
							const item = items[i];
							const description = item.textContent.trim();
							if (description) {
								reasonsMap[(i + 1).toString()] = description;
							}
						}
						break;
					}
				}
			}
			
			return JSON.stringify(reasonsMap);
		})()`, &reasonsJSON),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to extract reason descriptions: %v", err)
	}

	// Parse JSON data
	var reasonsMap map[string]string
	err = json.Unmarshal([]byte(reasonsJSON), &reasonsMap)
	if err != nil {
		return nil, fmt.Errorf("failed to parse reasons JSON: %v", err)
	}

	fmt.Printf("✅ Extracted %d reason descriptions\n", len(reasonsMap))
	for code, desc := range reasonsMap {
		if len(desc) > 100 {
			fmt.Printf("   Reason %s: %s...\n", code, desc[:100])
		} else {
			fmt.Printf("   Reason %s: %s\n", code, desc)
		}
	}

	return reasonsMap, nil
}

// ScrapeNonCompliantEmployersWithReasons scrapes both reason descriptions and employer data
func (s *Scraper) ScrapeNonCompliantEmployersWithReasons() ([]scraper_types.NonCompliantEmployerData, map[string]string, error) {
	// First, scrape the reason descriptions
	reasonDescriptions, err := s.ScrapeReasonDescriptions()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to scrape reason descriptions: %v", err)
	}

	// Then scrape the employer data
	employers, err := s.ScrapeNonCompliantEmployers()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to scrape employers: %v", err)
	}

	return employers, reasonDescriptions, nil
}

// ScrapeNonCompliantEmployers scrapes the non-compliant employers page with pagination
func (s *Scraper) ScrapeNonCompliantEmployers() ([]scraper_types.NonCompliantEmployerData, error) {
	fmt.Println("🎯 Navigating to non-compliant employers page...")

	err := chromedp.Run(s.ctx,
		chromedp.Navigate(nonCompliantURL),
		chromedp.WaitVisible("table tbody tr", chromedp.ByQuery),
		chromedp.Sleep(2*time.Second), // Wait for page to fully load
	)
	if err != nil {
		return nil, fmt.Errorf("failed to navigate to non-compliant page: %v", err)
	}

	// Try to set page length to 100 entries per page (optional - don't fail if it doesn't work)
	fmt.Println("🔧 Attempting to set page length to 100 entries...")
	err = chromedp.Run(s.ctx,
		chromedp.Evaluate(`(function() {
			// Try multiple possible selectors for the page length dropdown
			const selectors = [
				'select[name="wb-auto-4_length"]',
				'select[name*="_length"]',
				'select.dataTables_length',
				'.dataTables_length select'
			];
			
			for (let selector of selectors) {
				const select = document.querySelector(selector);
				if (select) {
					// Try to find and select the 100 option
					const options = select.querySelectorAll('option');
					for (let option of options) {
						if (option.value === '100' || option.textContent.includes('100')) {
							select.value = option.value;
							select.dispatchEvent(new Event('change', { bubbles: true }));
							console.log('Set page length to 100 using selector:', selector);
							return true;
						}
					}
				}
			}
			console.log('Could not find page length selector, continuing with default');
			return false;
		})()`, nil),
		chromedp.Sleep(3*time.Second), // Wait for potential page reload
	)
	if err != nil {
		fmt.Printf("⚠️  Could not set page length (continuing anyway): %v\n", err)
		// Don't return error - continue with default page size
	}

	var allEmployers []scraper_types.NonCompliantEmployerData
	pageNumber := 1

	for {
		fmt.Printf("📄 Scraping page %d...\n", pageNumber)
		
		// Wait for table to be ready
		err = chromedp.Run(s.ctx,
			chromedp.WaitVisible("table tbody tr", chromedp.ByQuery),
			chromedp.Sleep(1*time.Second),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to wait for table on page %d: %v", pageNumber, err)
		}

		// Scrape current page
		employers, err := s.parseNonCompliantEmployersPage()
		if err != nil {
			return nil, fmt.Errorf("failed to parse page %d: %v", pageNumber, err)
		}

		allEmployers = append(allEmployers, employers...)
		fmt.Printf("✅ Scraped %d employers from page %d (total: %d)\n", len(employers), pageNumber, len(allEmployers))

		// Check if next button exists and is enabled
		var nextButtonDisabled bool
		err = chromedp.Run(s.ctx,
			chromedp.Evaluate(`(function() {
				// Try multiple selectors for the next button
				const selectors = [
					'#wb-auto-4_next',
					'a[aria-label="Next"]',
					'.paginate_button.next',
					'[id*="_next"]'
				];
				
				for (let selector of selectors) {
					const nextButton = document.querySelector(selector);
					if (nextButton) {
						const isDisabled = nextButton.classList.contains('disabled') || 
						                  nextButton.parentElement.classList.contains('disabled') ||
						                  nextButton.hasAttribute('disabled') ||
						                  nextButton.getAttribute('aria-disabled') === 'true';
						if (!isDisabled) {
							return false; // Found enabled next button
						}
					}
				}
				
				return true; // No enabled next button found
			})()`, &nextButtonDisabled),
		)
		if err != nil {
			fmt.Printf("⚠️  Could not check next button status: %v\n", err)
			break
		}

		if nextButtonDisabled {
			fmt.Println("🏁 Reached last page or no next button found, stopping pagination")
			break
		}

		// Click next button using multiple selectors
		fmt.Printf("➡️  Navigating to page %d...\n", pageNumber+1)
		var clickSuccess bool
		err = chromedp.Run(s.ctx,
			chromedp.Evaluate(`(function() {
				// Try multiple selectors for the next button
				const selectors = [
					'#wb-auto-4_next',
					'a[aria-label="Next"]',
					'.paginate_button.next:not(.disabled)',
					'[id*="_next"]:not(.disabled)'
				];
				
				for (let selector of selectors) {
					const nextButton = document.querySelector(selector);
					if (nextButton && !nextButton.classList.contains('disabled')) {
						nextButton.click();
						return true;
					}
				}
				
				return false;
			})()`, &clickSuccess),
			chromedp.Sleep(2*time.Second), // Wait for page to load
		)
		if err != nil || !clickSuccess {
			fmt.Printf("⚠️  Could not click next button: %v, success: %v\n", err, clickSuccess)
			break
		}

		pageNumber++
	}

	fmt.Printf("🎉 Scraping completed! Total employers scraped: %d\n", len(allEmployers))
	return allEmployers, nil
}

// parseNonCompliantEmployersPage parses the current page of the non-compliant employers table
func (s *Scraper) parseNonCompliantEmployersPage() ([]scraper_types.NonCompliantEmployerData, error) {
	// Extract table data using JavaScript
	var tableDataJSON string
	err := chromedp.Run(s.ctx,
		chromedp.Evaluate(`(function() {
			// Function to clean and parse addresses
			function cleanAddress(rawAddress) {
				if (!rawAddress) return rawAddress;
				
				// Canadian provinces/territories (both English and French)
				const provinces = [
					'Alberta', 'AB', 'Colombie-Britannique', 'British Columbia', 'BC',
					'Manitoba', 'MB', 'Nouveau-Brunswick', 'New Brunswick', 'NB',
					'Terre-Neuve-et-Labrador', 'Newfoundland and Labrador', 'NL',
					'Territoires du Nord-Ouest', 'Northwest Territories', 'NT',
					'Nouvelle-Écosse', 'Nova Scotia', 'NS', 'Nunavut', 'NU',
					'Ontario', 'ON', 'Île-du-Prince-Édouard', 'Prince Edward Island', 'PE',
					'Québec', 'Quebec', 'QC', 'Saskatchewan', 'SK', 'Yukon', 'YT'
				];
				
				// Create regex pattern for provinces
				const provincePattern = '(' + provinces.join('|') + ')';
				const regex = new RegExp(provincePattern + '\\s*$', 'i');
				
				// Check if address ends with a province
				const provinceMatch = rawAddress.match(regex);
				if (!provinceMatch) return rawAddress; // No province found, return as-is
				
				const province = provinceMatch[1];
				const addressWithoutProvince = rawAddress.replace(regex, '').trim();
				
				// Look for common patterns where city runs into street address
				// Pattern 1: Street name directly followed by city name (no comma)
				// Look for patterns like "rue SomethingCityName" or "Street NameCityName"
				const streetPatterns = [
					/^(.+(?:rue|street|ave|avenue|blvd|boulevard|rd|road|dr|drive|place|pl|way)\s+\w+)([A-Z][a-z]+(?:\s+[A-Z][a-z]+)*)\s*,?\s*$/i,
					/^(.+\s+)([A-Z][a-z]+(?:\s+[A-Z][a-z]+)*)\s*,?\s*$/
				];
				
				for (let pattern of streetPatterns) {
					const match = addressWithoutProvince.match(pattern);
					if (match && match[2]) {
						const streetAddress = match[1].trim();
						const cityName = match[2].trim();
						
						// Additional validation: city name should be reasonable length (2+ chars, not all caps)
						if (cityName.length > 1 && !cityName.match(/^[A-Z]+$/)) {
							return streetAddress + ', ' + cityName + ', ' + province;
						}
					}
				}
				
				// Fallback: if we have a province but couldn't parse the city, just add commas where needed
				if (addressWithoutProvince.includes(',')) {
					return rawAddress; // Already has commas, probably formatted correctly
				} else {
					// Try to add comma before last word before province (assuming it's the city)
					const parts = addressWithoutProvince.trim().split(/\s+/);
					if (parts.length > 1) {
						const cityPart = parts[parts.length - 1];
						const streetPart = parts.slice(0, -1).join(' ');
						return streetPart + ', ' + cityPart + ', ' + province;
					}
				}
				
				return rawAddress; // Return original if can't parse
			}
			
			const rows = document.querySelectorAll('table tbody tr');
			let data = [];
			
			for (let row of rows) {
				const cells = row.querySelectorAll('td');
				if (cells.length >= 7) {
					// Extract reason codes from the Reason(s) column
					const reasonCell = cells[3];
					const reasonLinks = reasonCell.querySelectorAll('a');
					const reasonCodes = Array.from(reasonLinks).map(link => {
						// Extract reason code from href attribute like "#list6" -> "6"
						const href = link.getAttribute('href');
						if (href && href.startsWith('#list')) {
							return href.replace('#list', '');
						}
						// Fallback to text content if href doesn't match expected format
						return link.textContent.trim();
					}).filter(code => code && code !== '');
					
					// Parse penalty amount
					const penaltyText = cells[5].textContent.trim();
					const penaltyMatch = penaltyText.match(/\$?([\d,]+)/);
					const penaltyAmount = penaltyMatch ? parseInt(penaltyMatch[1].replace(/,/g, '')) : 0;
					
					data.push({
						businessOperatingName: cells[0].textContent.trim(),
						businessLegalName: cells[1].textContent.trim(),
						address: cleanAddress(cells[2].textContent.trim()),
						reasonCodes: reasonCodes,
						dateOfFinalDecision: cells[4].textContent.trim(),
						penaltyAmount: penaltyAmount,
						penaltyCurrency: 'CAD',
						status: cells[6].textContent.trim()
					});
				}
			}
			
			return JSON.stringify(data);
		})()`, &tableDataJSON),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to extract table data: %v", err)
	}

	// Parse JSON data
	var employers []scraper_types.NonCompliantEmployerData
	err = json.Unmarshal([]byte(tableDataJSON), &employers)
	if err != nil {
		return nil, fmt.Errorf("failed to parse table JSON: %v", err)
	}

	// Clean and validate the data
	var cleanedEmployers []scraper_types.NonCompliantEmployerData
	for _, employer := range employers {
		// Skip empty rows
		if employer.BusinessOperatingName == "" {
			continue
		}

		// Clean date format (expecting YYYY-MM-DD)
		if employer.DateOfFinalDecision != "" {
			employer.DateOfFinalDecision = s.cleanDate(employer.DateOfFinalDecision)
		}

		// Ensure reason codes are not empty
		var validReasons []string
		for _, reason := range employer.ReasonCodes {
			if reason != "" {
				validReasons = append(validReasons, reason)
			}
		}
		employer.ReasonCodes = validReasons

		cleanedEmployers = append(cleanedEmployers, employer)
	}

	return cleanedEmployers, nil
}

// cleanDate attempts to parse and normalize various date formats to YYYY-MM-DD
func (s *Scraper) cleanDate(dateStr string) string {
	dateStr = strings.TrimSpace(dateStr)
	if dateStr == "" {
		return ""
	}

	// Try different date formats that might appear on the page
	formats := []string{
		"2006-01-02", // YYYY-MM-DD (already correct)
		"2006/01/02", // YYYY/MM/DD
		"01/02/2006", // MM/DD/YYYY
		"02/01/2006", // DD/MM/YYYY
		"January 2, 2006",
		"Jan 2, 2006",
		"2 January 2006",
		"2 Jan 2006",
	}

	for _, format := range formats {
		if parsedTime, err := time.Parse(format, dateStr); err == nil {
			return parsedTime.Format("2006-01-02")
		}
	}

	// Handle special case of invalid dates like "2019-02-29"
	if matched := regexp.MustCompile(`(\d{4})-(\d{2})-(\d{2})`).FindStringSubmatch(dateStr); matched != nil {
		year := matched[1]
		month := matched[2]
		day := matched[3]
		
		// Try to fix common invalid dates
		fixedDateStr := s.fixInvalidDate(year, month, day)
		if fixedDateStr != dateStr {
			fmt.Printf("🔧 Fixed invalid date %s -> %s\n", dateStr, fixedDateStr)
			return fixedDateStr
		}
	}

	// If no format matches, return empty string instead of invalid date
	fmt.Printf("⚠️  Could not parse date: %s, skipping\n", dateStr)
	return ""
}

// fixInvalidDate attempts to fix common invalid dates like Feb 29 in non-leap years
func (s *Scraper) fixInvalidDate(year, month, day string) string {
	yearInt, _ := strconv.Atoi(year)
	monthInt, _ := strconv.Atoi(month)
	dayInt, _ := strconv.Atoi(day)
	
	// Fix February 29 in non-leap years
	if monthInt == 2 && dayInt == 29 && !isLeapYear(yearInt) {
		return fmt.Sprintf("%s-02-28", year) // Change to Feb 28
	}
	
	// Fix other impossible dates
	daysInMonth := []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	if isLeapYear(yearInt) {
		daysInMonth[1] = 29 // February in leap year
	}
	
	if monthInt >= 1 && monthInt <= 12 && dayInt > daysInMonth[monthInt-1] {
		// Day exceeds maximum for the month, set to last day of month
		return fmt.Sprintf("%04d-%02d-%02d", yearInt, monthInt, daysInMonth[monthInt-1])
	}
	
	// Return original if no fix needed
	return fmt.Sprintf("%s-%s-%s", year, month, day)
}

// isLeapYear checks if a year is a leap year
func isLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}
