import { Page } from "puppeteer";

export const numPages = async (
  page: Page,
  numberOfPages: number,
  timeout: number,
): Promise<void> => {
  // If -1 is passed, scrape all available pages
  const scrapeAll = numberOfPages === -1;
  let i = 0;

  while (true) {
    // Check if we should stop based on page limit (unless scraping all)
    if (!scrapeAll && i >= numberOfPages) {
      break;
    }

    await new Promise((resolve) => setTimeout(resolve, timeout));
    const moreButton = await page.$("#moreresultbutton");

    if (moreButton) {
      // @ts-ignore
      await moreButton.evaluate((b) => b.click());
      i++;

      if (scrapeAll) {
        console.log(`${i} 📄(s) loaded (scraping all pages...)`);
      } else {
        console.log(`${i} 📄(s) loaded out of ${numberOfPages}`);
      }

      await new Promise((resolve) => setTimeout(resolve, timeout));
    } else {
      console.log(`No more results after ${i} pages 😔`);
      console.log(`Finished loading all available pages`);
      await new Promise((resolve) => setTimeout(resolve, timeout * 7));
      break;
    }
  }
};
