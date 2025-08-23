export interface User {
  id: string;
  email: string;
  username: string;
  role: 'user' | 'admin';
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

export interface JobSubredditPost {
  id: string;
  job_posting_id: string;
  subreddit_id: string;
  reddit_post_id?: string;
  reddit_post_url?: string;
  posted_at: string;
  created_at: string;
  subreddit_name?: string; // Joined data when queried with subreddit info
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
  reddit_posted: boolean;
  reddit_approval_status: 'pending' | 'approved' | 'rejected';
  reddit_approved_by?: string;
  reddit_approved_at?: string;
  reddit_rejection_reason?: string;
  description?: string;
  scraping_run_id: string;
  created_at: string;
  updated_at: string;
  subreddit_posts?: JobSubredditPost[]; // Joined data when queried with posting info
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

// Admin interfaces for Reddit approval
export interface RedditApprovalRequest {
  approved_by: string;
  subreddit_ids?: string[];
}

export interface RedditRejectionRequest {
  rejected_by: string;
  reason?: string;
}

export interface BulkApprovalRequest {
  job_ids: string[];
  approved_by: string;
  subreddit_ids?: string[];
}

export interface BulkRejectionRequest {
  job_ids: string[];
  rejected_by: string;
  reason?: string;
}


export interface Subreddit {
  id: string;
  name: string;
  is_active: boolean;
  post_count: number;
  last_posted_at?: string;
  created_at: string;
  updated_at: string;
}

export interface CreateSubredditRequest {
  name: string;
  is_active?: boolean;
}

export interface UpdateSubredditRequest {
  is_active?: boolean;
}

export interface SubredditsResponse {
  subreddits: Subreddit[];
}

// Report Types
export interface Report {
  id: string;
  user_id: string;
  business_name: string;
  business_address: string;
  report_source: 'employment' | 'observation' | 'public_record';
  confidence_level?: number; // Deprecated: use tfw_ratio
  tfw_ratio?: 'few' | 'many' | 'most' | 'all';
  additional_notes?: string;
  created_at: string;
  updated_at: string;
}

// Request Types
export interface CreateReportRequest {
  business_name: string;
  business_address: string;
  report_source: 'employment' | 'observation' | 'public_record';
  confidence_level?: number; // Deprecated: use tfw_ratio
  tfw_ratio?: 'few' | 'many' | 'most' | 'all';
  additional_notes?: string;
}

export interface UpdateReportRequest extends CreateReportRequest {}

export interface ReportListResponse {
  reports: Report[];
  pagination: {
    limit: number;
    offset: number;
    total?: number;
  };
}

export interface ReportFilters {
  limit?: number;
  offset?: number;
}


export interface ReportsByAddress {
  business_name: string;
  business_address: string;
  report_count: number;
  confidence_level: number; // Deprecated: use tfw_ratio_* fields
  tfw_ratio_few: number;
  tfw_ratio_many: number;
  tfw_ratio_most: number;
  tfw_ratio_all: number;
  latest_report: string;
}

export interface ReportsByAddressResponse {
  data: ReportsByAddress[];
  limit: number;
  offset: number;
  count: number;
}

export interface RedditApprovalStats {
  pending_count: number;
  approved_count: number;
  rejected_count: number;
}

export interface LMIAEmployersByResourceResponse {
  employers: LMIAEmployer[];
  count: number;
}

export interface LMIAEmployerGeoLocation {
  id: string;
  employer: string;
  address?: string;
  province_territory?: string;
  approved_lmias?: number;
  approved_positions?: number;
  quarter: string;
  year: number;
  latitude?: number;
  longitude?: number;
  total_lmias: number;
}

export interface LMIAGeolocationResponse {
  employers: LMIAEmployerGeoLocation[];
  count: number;
  year: number;
  quarter?: string;
  limit: number;
}

export interface PostalCodeBusiness {
  employer: string;
  occupation: string;
  approved_lmias: number;
  approved_positions: number;
}

export interface PostalCodeLocation {
  postal_code: string;
  latitude: number;
  longitude: number;
  businesses: PostalCodeBusiness[];
  total_lmias: number;
  business_count: number;
}

export interface PostalCodeLocationsResponse {
  locations: PostalCodeLocation[];
  count: number;
  year: number;
  quarter?: string;
  limit: number;
}

export interface LMIAEmployersByPostalCodeResponse {
  employers: LMIAEmployer[];
  count: number;
  postal_code: string;
  year: number;
  quarter?: string;
  limit: number;
}

// Gemini content generation types
export interface GeneratedRedditPost {
  job_id: string;
  content: string;
  error?: string;
}

export interface BulkGenerationResponse {
  posts: GeneratedRedditPost[];
}

export interface BulkGenerationRequest {
  job_ids: string[];
}

// Reddit preview types
export interface RedditPreview {
  job_id: string;
  title: string;
  body: string;
  content_type: 'ai' | 'template';
  error?: string;
}
