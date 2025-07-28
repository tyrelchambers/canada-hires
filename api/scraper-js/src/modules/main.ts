//Puppeteer allows you to control headless Chrome or Chromium over the DevTools Protocol.
import puppeteer, { Browser, Page } from "puppeteer";
import { displayMessage } from "./main/displayMessage.js";
import { numPages } from "./main/numPages.js";
import { runScraper } from "./main/runScraper.js";
import { cleanData } from "./main/cleanData.js";
import fs from "fs";
import { JobData } from "../types.js";
import { ApiClient } from "./apiClient.js";

const baseUrl: string = "https://www.jobbank.gc.ca";
const lmiaUrl: string = "https://www.jobbank.gc.ca/jobsearch/jobsearch?fsrc=32";
//controls timeouts to avoid being blocked by the website
const timeout: number = Math.floor(Math.random() * 1000);

export const main = async (
  jobTitle: string, 
  province: string, 
  numberOfPages: number, 
  saveToJSON: boolean,
  saveToAPI: boolean = true
): Promise<void> => {
  // Initialize API client if saving to API
  let apiClient: ApiClient | null = null;
  if (saveToAPI) {
    apiClient = new ApiClient();
    const isConnected = await apiClient.testConnection();
    if (!isConnected) {
      console.log('⚠️  API connection failed, continuing with JSON-only save');
      saveToAPI = false;
    }
  }

  const browser: Browser = await puppeteer.launch({ headless: true });
  const page: Page = await browser.newPage();

  console.log("🎯 Navigating directly to LMIA jobs page...");
  await page.goto(lmiaUrl, { waitUntil: 'networkidle2' });
  
  // Wait for the page to load - no need to search since we're already on LMIA jobs
  await page.waitForSelector("#moreresultbutton", { timeout: 15000 });

  // Get and display total LMIA results count
  const totalResults = await page.$eval("#results-count", (el) => el.textContent?.trim() || "0");
  console.log(`\n📊 Total LMIA jobs to scrape: ${totalResults}`);

  // Start API scraping session if enabled
  if (saveToAPI && apiClient) {
    await apiClient.startScrapingSession();
  }

  //clicks on the more results button to load more results
  await numPages(page, numberOfPages, timeout);
  
  //scrapes the jobs from the page and adds them to the jobArray
  let jobArray: JobData[] = [];
  await runScraper(jobArray, page, baseUrl);

  //cleans the data
  jobArray = cleanData(jobArray);

  // Log results to console
  console.log("\n=== SCRAPED JOBS ===");
  console.log(`Total jobs found: ${jobArray.length}`);
  jobArray.forEach((job: JobData, index: number) => {
    console.log(`\nJob ${index + 1}:`);
    console.log(`  Title: ${job.jobTitle}`);
    console.log(`  Business: ${job.business}`);
    console.log(`  Location: ${job.location}`);
    console.log(`  Salary: ${job.salary}`);
    console.log(`  Date: ${job.date}`);
    console.log(`  URL: ${job.jobUrl}`);
  });

  // Save to API if enabled
  if (saveToAPI && apiClient) {
    try {
      await apiClient.submitJobs(jobArray);
      await apiClient.completeScrapingSession(numberOfPages, jobArray.length, jobArray.length);
    } catch (error) {
      console.error('❌ Failed to save to API:', error);
    }
  }

  //optional save the jobArray to JSON file
  if (saveToJSON) {
    const timestamp: string = new Date().toISOString().replace(/[:.]/g, '-');
    const filename: string = `jobs_${timestamp}.json`;
    fs.writeFileSync(filename, JSON.stringify(jobArray, null, 2));
    console.log(`\nResults saved to ${filename}`);
  }

  browser.close();
  console.log("Browser closed 👋");
  process.exit();
};