import { useQuery } from "@tanstack/react-query";
import { useApiClient } from "./useApiClient";
import { JobPostingsResponse, JobPostingFilters, JobStats } from "@/types";

export function useJobPostings(filters: JobPostingFilters = {}) {
  const apiClient = useApiClient();

  return useQuery({
    queryKey: ["jobPostings", filters],
    queryFn: async (): Promise<JobPostingsResponse> => {
      const params = new URLSearchParams();

      // Add filters to params
      if (filters.search) params.append("search", filters.search);
      if (filters.employer) params.append("employer", filters.employer);
      if (filters.city) params.append("city", filters.city);
      if (filters.province) params.append("province", filters.province);
      if (filters.title) params.append("title", filters.title);
      if (filters.salary_min) params.append("salary_min", filters.salary_min.toString());
      if (filters.sort_by) params.append("sort_by", filters.sort_by);
      if (filters.sort_order) params.append("sort_order", filters.sort_order);
      if (filters.limit) params.append("limit", filters.limit.toString());
      if (filters.offset) params.append("offset", filters.offset.toString());
      if (filters.days !== undefined) params.append("days", filters.days.toString());

      const response = await apiClient.get(`/jobs?${params.toString()}`);
      return response.data;
    },
    // Keep data fresh for 5 minutes
    staleTime: 5 * 60 * 1000,
  });
}

export function useJobStats() {
  const apiClient = useApiClient();

  return useQuery({
    queryKey: ["jobStats"],
    queryFn: async (): Promise<JobStats> => {
      const response = await apiClient.get("/jobs/stats");
      return response.data;
    },
    // Keep stats fresh for 10 minutes since they don't change often
    staleTime: 10 * 60 * 1000,
  });
}