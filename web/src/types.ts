export interface User {
  id: string;
  email: string;
  username: string;
  verification_tier: string;
  account_status: string;
  email_verified: boolean;
  created_at: string;
}

export interface LoginRequest {
  email: string;
}

export interface LoginResponse {
  message: string;
  email: string;
}

export interface LMIAEmployer {
  id: string;
  resource_id: string;
  province_territory?: string;
  program_stream?: string;
  employer: string;
  address?: string;
  occupation?: string;
  incorporate_status?: string;
  approved_lmias?: number;
  approved_positions?: number;
  created_at: string;
  updated_at: string;
  year: number;
  quarter: string;
}

export interface LMIAResource {
  id: string;
  resource_id: string;
  name: string;
  quarter: string;
  year: number;
  url: string;
  format: string;
  language: string;
  size_bytes?: number;
  last_modified?: string;
  date_published?: string;
  downloaded_at?: string;
  processed_at?: string;
  created_at: string;
  updated_at: string;
}

export interface CronJob {
  id: string;
  job_name: string;
  status: string;
  started_at: string;
  completed_at?: string;
  error_message?: string;
  resources_processed: number;
  records_processed: number;
  created_at: string;
}

export interface LMIASearchResponse {
  employers: LMIAEmployer[];
  count: number;
  query: string;
  limit: number;
}

export interface LMIALocationResponse {
  employers: LMIAEmployer[];
  count: number;
  city: string;
  province: string;
  limit: number;
}

export interface LMIAResourcesResponse {
  resources: LMIAResource[];
  count: number;
}

export interface LMIAStatsResponse {
  total_resources: number;
  processed_resources: number;
  last_update?: string;
  last_update_status: string;
  total_records_processed: number;
  total_records: number;
  distinct_employers: number;
  year_range: {
    min_year: number;
    max_year: number;
  };
}

export interface JobPosting {
  id: string;
  job_bank_id?: string;
  title: string;
  employer: string;
  location: string;
  province?: string;
  city?: string;
  salary_min?: number;
  salary_max?: number;
  salary_type?: string;
  salary_raw?: string;
  posting_date?: string;
  url: string;
  is_tfw: boolean;
  has_lmia: boolean;
  description?: string;
  scraping_run_id: string;
  created_at: string;
  updated_at: string;
}

export interface JobPostingsResponse {
  jobs: JobPosting[];
  total: number;
  limit: number;
  offset: number;
  has_more: boolean;
}

export interface JobPostingFilters {
  search?: string;
  employer?: string;
  city?: string;
  province?: string;
  title?: string;
  salary_min?: number;
  sort_by?: string;
  sort_order?: 'asc' | 'desc';
  limit?: number;
  offset?: number;
  days?: number;
}

export interface JobStats {
  total_jobs: number;
  total_employers: number;
  top_employers: Array<{
    employer: string;
    job_count: number;
    earliest_posting?: string;
    latest_posting?: string;
  }>;
}
