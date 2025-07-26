export interface User {
  id: string
  email: string
  username: string
  verification_tier: string
  account_status: string
  email_verified: boolean
  created_at: string
}

export interface LoginRequest {
  email: string
}

export interface LoginResponse {
  message: string
  email: string
}

export interface JobPosting {
  id: string
  job_bank_id?: string
  title: string
  employer: string
  location: string
  province?: string
  city?: string
  salary_min?: number
  salary_max?: number
  salary_type?: string
  salary_raw?: string
  posting_date?: string
  url: string
  is_tfw: boolean
  has_lmia: boolean
  description?: string
  scraping_run_id: string
  created_at: string
  updated_at: string
}

export interface JobPostingsResponse {
  jobs: JobPosting[]
  total: number
  limit: number
  offset: number
  has_more: boolean
}

export interface JobPostingFilters {
  search?: string
  employer?: string
  city?: string
  province?: string
  title?: string
  salary_min?: number
  sort_by?: string
  sort_order?: 'asc' | 'desc'
  limit?: number
  offset?: number
  days?: number
}

export interface JobStats {
  total_jobs: number
  total_employers: number
  top_employers: Array<{
    employer: string
    job_count: number
    earliest_posting?: string
    latest_posting?: string
  }>
}