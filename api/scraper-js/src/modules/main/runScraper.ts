import * as cheerio from "cheerio";
import { Page } from "puppeteer";
import { pushToArray } from "./pushToArray.js";
import { JobData } from "../../types.js";

export const runScraper = async (
  jobArray: JobData[],
  page: Page,
  baseUrl: string,
): Promise<void> => {
  //loads the html from the page and uses cheerio to select the jobs from the page
  const html: string = await page.content();

  const $: cheerio.CheerioAPI = cheerio.load(html);
  const jobList = $("article");

  //loop through the jobList and add the scraped data to the jobArray
  pushToArray(jobArray, jobList, $, baseUrl);

  console.log(`Scraped ${jobArray.length} jobs from the page`);
};
