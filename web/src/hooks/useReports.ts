import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useApiClient } from "./useApiClient";

interface Report {
  id: string;
  user_id: string;
  business_name: string;
  business_address: string;
  report_source: string;
  confidence_level?: number;
  additional_notes?: string;
  status: "pending" | "approved" | "rejected" | "flagged";
  moderated_by?: string;
  moderation_notes?: string;
  ip_address?: string;
  created_at: string;
  updated_at: string;
}

interface PaginationInfo {
  limit: number;
  offset: number;
  total?: number;
}

interface ReportsResponse {
  reports: Report[];
  pagination: PaginationInfo;
}

interface ReportFilters {
  query?: string;
  city?: string;
  province?: string;
  status?: string;
  year?: string;
  limit?: number;
  offset?: number;
}

interface CreateReportRequest {
  business_name: string;
  business_address: string;
  report_source: string;
  confidence_level?: number;
  additional_notes?: string;
}

export function useReports(filters: ReportFilters) {
  const apiClient = useApiClient();

  return useQuery({
    queryKey: ["reports", filters],
    queryFn: async (): Promise<ReportsResponse> => {
      const params = new URLSearchParams();

      if (filters.query) params.append("query", filters.query);
      if (filters.city) params.append("city", filters.city);
      if (filters.province) params.append("province", filters.province);
      if (filters.status) params.append("status", filters.status);
      if (filters.year) params.append("year", filters.year);

      if (filters.limit) params.append("limit", filters.limit.toString());
      if (filters.offset) params.append("offset", filters.offset.toString());

      const response = apiClient
        .get<ReportsResponse>(`/reports?${params}`)
        .then((res) => res.data);
      return response;
    },
    staleTime: 1000 * 60 * 5, // 5 minutes
  });
}

export function useCreateReport() {
  const apiClient = useApiClient();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (data: CreateReportRequest): Promise<Report> => {
      const response = apiClient
        .post<Report>("/reports", data)
        .then((res) => res.data);
      return response;
    },
    onSuccess: async () => {
      // Invalidate reports queries to refetch data
      await queryClient.invalidateQueries({ queryKey: ["reports"] });
    },
  });
}

export type { Report, ReportsResponse, ReportFilters, CreateReportRequest };
