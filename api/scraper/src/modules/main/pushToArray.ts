import * as cheerio from "cheerio";
import { JobData } from "../../types.js";
import { Element } from "domhandler";
// Extract job bank ID from URL like https://www.jobbank.gc.ca/jobsearch/jobpostingtfw/44739738
const extractJobBankId = (url: string): string | undefined => {
  const match = url.match(/\/jobpostingtfw\/(\d+)/);
  return match ? match[1] : undefined;
};

export const pushToArray = (
  jobArray: JobData[],
  jobList: cheerio.Cheerio<Element>,
  $: cheerio.CheerioAPI,
  baseUrl: string,
): void => {
  //loop through the jobList and add scrape the data to the jobArray
  for (let i = 0; i < jobList.length; i++) {
    const jobTitle: string = $(jobList[i]).find(".noctitle").text().trim();

    const list: cheerio.Cheerio<Element> = $(jobList[i]).find(".list-unstyled");

    const business: string = list.find(".business").text().trim();

    const location: string = list.find(".location").text().trim();

    const salary: string = list.find(".salary").text().trim();

    const date: string = list.find(".date").text().trim();

    const jobUrl: string = baseUrl + $(jobList[i]).find("a").attr("href");

    // Extract job bank ID from URL
    const jobBankId = extractJobBankId(jobUrl);

    if (jobTitle && jobUrl) {
      jobArray.push({
        jobTitle,
        business,
        salary,
        location,
        jobUrl,
        date,
        jobBankId,
      });

      console.log(`${i + 1} job(s) loaded: ${jobTitle}`);
    }
  }
};
