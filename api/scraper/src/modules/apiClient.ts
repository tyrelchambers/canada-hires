import axios from 'axios';
import { JobData } from '../types.js';
import { config } from '../config/index.js';

export interface ScrapingRun {
  id: string;
  status: string;
  started_at: string;
  total_pages: number;
  jobs_scraped: number;
  jobs_stored: number;
  last_page_scraped: number;
  created_at: string;
}

export class ApiClient {
  private baseUrl: string;
  private scrapingRunId: string | null = null;

  constructor(baseUrl?: string) {
    this.baseUrl = baseUrl || config.apiUrl;
  }

  async startScrapingSession(): Promise<string> {
    try {
      console.log('🚀 Starting new scraping session...');
      const response = await axios.post<ScrapingRun>(`${this.baseUrl}/api/jobs/scraping-runs`);
      this.scrapingRunId = response.data.id;
      console.log(`✅ Scraping session started: ${this.scrapingRunId}`);
      return this.scrapingRunId;
    } catch (error) {
      console.error('❌ Failed to start scraping session:', error);
      throw error;
    }
  }

  async submitJobs(jobs: JobData[]): Promise<void> {
    if (!this.scrapingRunId) {
      throw new Error('No active scraping session. Call startScrapingSession() first.');
    }

    if (jobs.length === 0) {
      console.log('⚠️  No jobs to submit');
      return;
    }

    // Submit in smaller batches to avoid overwhelming the API
    const batchSize = 500; // Smaller batches for network efficiency
    const totalJobs = jobs.length;
    let totalProcessed = 0;

    console.log(`📤 Submitting ${totalJobs} jobs to API in batches of ${batchSize}...`);

    for (let i = 0; i < totalJobs; i += batchSize) {
      const end = Math.min(i + batchSize, totalJobs);
      const batch = jobs.slice(i, end);

      try {
        const response = await axios.post(
          `${this.baseUrl}/api/jobs/scraping-runs/${this.scrapingRunId}/jobs`,
          batch,
          {
            headers: {
              'Content-Type': 'application/json'
            }
          }
        );
        totalProcessed += response.data.jobs_processed;
        console.log(`✅ Batch ${Math.floor(i/batchSize) + 1}: Submitted ${response.data.jobs_processed} jobs (${totalProcessed}/${totalJobs} total)`);
      } catch (error) {
        console.error(`❌ Failed to submit batch ${Math.floor(i/batchSize) + 1}:`, error);
        throw error;
      }
    }

    console.log(`🎉 Successfully submitted all ${totalProcessed} jobs to API`);
  }

  async completeScrapingSession(totalPages: number, jobsScraped: number, jobsStored: number): Promise<void> {
    if (!this.scrapingRunId) {
      throw new Error('No active scraping session.');
    }

    try {
      console.log('🏁 Completing scraping session...');
      await axios.post(
        `${this.baseUrl}/api/jobs/scraping-runs/${this.scrapingRunId}/complete`,
        {
          total_pages: totalPages,
          jobs_scraped: jobsScraped,
          jobs_stored: jobsStored
        }
      );
      console.log('✅ Scraping session completed successfully');
    } catch (error) {
      console.error('❌ Failed to complete scraping session:', error);
      throw error;
    }
  }

  async testConnection(): Promise<boolean> {
    try {
      console.log('🔗 Testing API connection...');
      await axios.get(`${this.baseUrl}/api/jobs/stats`);
      console.log('✅ API connection successful');
      return true;
    } catch (error) {
      console.error('❌ API connection failed:', error);
      return false;
    }
  }
}