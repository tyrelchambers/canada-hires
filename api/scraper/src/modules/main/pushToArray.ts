import * as cheerio from "cheerio";
import { JobData } from "../../types.js";

export const pushToArray = (
  jobArray: JobData[], 
  jobList: cheerio.Cheerio<cheerio.Element>, 
  $: cheerio.CheerioAPI, 
  baseUrl: string
): void => {
  //loop through the jobList and add scrape the data to the jobArray
  for (let i = 0; i < jobList.length; i++) {
    const jobTitle: string = $(jobList[i])
      .find(".noctitle")
      .text()
      .trim();

    const list: cheerio.Cheerio<cheerio.Element> = $(jobList[i])
      .find(".list-unstyled");

    const business: string = list
      .find(".business")
      .text()
      .trim();

    const location: string = list
      .find(".location")
      .text()
      .trim();

    const salary: string = list
      .find(".salary")
      .text()
      .trim();

    const date: string = list
      .find(".date")
      .text()
      .trim();

    const jobUrl: string = baseUrl + $(jobList[i]).find("a").attr("href");

    if (jobTitle && jobUrl) {
      jobArray.push({
        jobTitle,
        business,
        salary,
        location,
        jobUrl,
        date,
      });

      console.log(`${i + 1} job(s) loaded: ${jobTitle}`);
    }
  }
};