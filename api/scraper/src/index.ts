import { main } from "./modules/main.js";

console.log("🇨🇦🍁 Job Bank Scraper - Starting...");
console.log("Scraping ALL LMIA jobs (fsrc=32)");

// Fixed parameters for simple execution
const jobTitle: string = ""; // Empty = all job titles
const province: string = ""; // Empty = all provinces
const numberOfPages: number = -1; // -1 means scrape all available pages
const saveToJSON: boolean = false; // Save to JSON file
const saveToAPI: boolean = true; // Save to API database

// Run the scraper immediately
main(jobTitle, province, numberOfPages, saveToJSON, saveToAPI);
