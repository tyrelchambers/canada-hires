package scraper

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/playwright-community/playwright-go"
	"scraper/types"
)

const (
	baseURL = "https://www.jobbank.gc.ca"
	lmiaURL = "https://www.jobbank.gc.ca/jobsearch/jobsearch?fsrc=32"
)

type Scraper struct {
	pw      *playwright.Playwright
	browser playwright.Browser
	page    playwright.Page
	timeout time.Duration
}

func NewScraper() (*Scraper, error) {
	pw, err := playwright.Run()
	if err != nil {
		return nil, fmt.Errorf("could not start playwright: %v", err)
	}
	
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
	})
	if err != nil {
		pw.Stop()
		return nil, fmt.Errorf("could not launch browser: %v", err)
	}
	
	page, err := browser.NewPage()
	if err != nil {
		browser.Close()
		pw.Stop()
		return nil, fmt.Errorf("could not create page: %v", err)
	}
	
	return &Scraper{
		pw:      pw,
		browser: browser,
		page:    page,
		timeout: time.Duration(1000) * time.Millisecond,
	}, nil
}

func (s *Scraper) Close() {
	if s.browser != nil {
		s.browser.Close()
	}
	if s.pw != nil {
		s.pw.Stop()
	}
}

func (s *Scraper) ScrapeLMIAJobs(numberOfPages int) ([]types.JobData, error) {
	fmt.Println("🎯 Navigating directly to LMIA jobs page...")
	
	_, err := s.page.Goto(lmiaURL)
	if err != nil {
		return nil, fmt.Errorf("failed to navigate to LMIA page: %v", err)
	}
	
	// Wait for the page to load
	_, err = s.page.WaitForSelector("#moreresultbutton", playwright.PageWaitForSelectorOptions{
		Timeout: playwright.Float(15000),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to wait for more results button: %v", err)
	}
	
	// Get total results count
	totalResults, err := s.page.Locator("#results-count").TextContent()
	if err != nil {
		return nil, fmt.Errorf("failed to get results count: %v", err)
	}
	
	fmt.Printf("\n📊 Total LMIA jobs to scrape: %s\n", strings.TrimSpace(totalResults))
	
	// Load more pages
	if err := s.loadMorePages(numberOfPages); err != nil {
		return nil, fmt.Errorf("failed to load more pages: %v", err)
	}
	
	// Parse jobs from the page
	jobs, err := s.parseJobs()
	if err != nil {
		return nil, fmt.Errorf("failed to parse jobs: %v", err)
	}
	
	// Clean the data
	cleanedJobs := s.cleanData(jobs)
	
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
		
		// Check if more button exists
		moreButton := s.page.Locator("#moreresultbutton")
		count, err := moreButton.Count()
		if err != nil {
			return fmt.Errorf("failed to check for more button: %v", err)
		}
		
		if count > 0 {
			// Check if button is visible and clickable
			isVisible, err := moreButton.IsVisible()
			if err != nil {
				return fmt.Errorf("failed to check button visibility: %v", err)
			}
			
			if isVisible {
				err = moreButton.Click()
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

func (s *Scraper) parseJobs() ([]types.JobData, error) {
	articles, err := s.page.Locator("article").All()
	if err != nil {
		return nil, fmt.Errorf("failed to get articles: %v", err)
	}
	
	var jobs []types.JobData
	
	for i, article := range articles {
		jobTitle, err := article.Locator(".noctitle").TextContent()
		if err != nil {
			continue
		}
		jobTitle = strings.TrimSpace(jobTitle)
		
		list := article.Locator(".list-unstyled")
		
		business, err := list.Locator(".business").TextContent()
		if err != nil {
			business = ""
		}
		business = strings.TrimSpace(business)
		
		location, err := list.Locator(".location").TextContent()
		if err != nil {
			location = ""
		}
		location = strings.TrimSpace(location)
		
		salary, err := list.Locator(".salary").TextContent()
		if err != nil {
			salary = ""
		}
		salary = strings.TrimSpace(salary)
		
		date, err := list.Locator(".date").TextContent()
		if err != nil {
			date = ""
		}
		date = strings.TrimSpace(date)
		
		href, err := article.Locator("a").GetAttribute("href")
		if err != nil {
			continue
		}
		jobURL := baseURL + href
		
		// Extract job bank ID from URL
		jobBankID := extractJobBankID(jobURL)
		
		if jobTitle != "" && jobURL != "" {
			jobs = append(jobs, types.JobData{
				JobTitle:  jobTitle,
				Business:  business,
				Salary:    salary,
				Location:  location,
				JobURL:    jobURL,
				Date:      date,
				JobBankID: jobBankID,
			})
			
			fmt.Printf("%d job(s) loaded: %s\n", i+1, jobTitle)
		}
	}
	
	fmt.Printf("Scraped %d jobs from the page\n", len(jobs))
	return jobs, nil
}

func extractJobBankID(url string) string {
	re := regexp.MustCompile(`/jobpostingtfw/(\d+)`)
	matches := re.FindStringSubmatch(url)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func (s *Scraper) cleanData(jobs []types.JobData) []types.JobData {
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