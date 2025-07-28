package scraper

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	scraper_types "canada-hires/scraper-types"

	"github.com/chromedp/chromedp"
)

const (
	baseURL = "https://www.jobbank.gc.ca"
	lmiaURL = "https://www.jobbank.gc.ca/jobsearch/jobsearch?fsrc=32"
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

			if len(jobs) <= 5 || len(jobs)%100 == 0 {
				fmt.Printf("Job %d: %s at %s\n", len(jobs), jobTitle, business)
			}
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
		cleanURL = fmt.Sprintf("%s/jobpostingtfw/%s", baseURL, jobBankID)
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
