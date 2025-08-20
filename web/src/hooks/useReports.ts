import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useApiClient } from "./useApiClient";
import { ReportsByAddress, ReportsByAddressResponse } from "@/types";

interface Report {
  id: string;
  user_id: string;
  business_name: string;
  business_address: string;
  report_source: string;
  confidence_level?: number; // Deprecated: use tfw_ratio
  tfw_ratio?: 'few' | 'many' | 'most' | 'all';
  additional_notes?: string;
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
  year?: string;
  limit?: number;
  offset?: number;
}

interface CreateReportRequest {
  business_name: string;
  business_address: string;
  report_source: string;
  confidence_level?: number; // Deprecated: use tfw_ratio
  tfw_ratio?: 'few' | 'many' | 'most' | 'all';
  additional_notes?: string;
}

interface UpdateReportRequest {
  business_name: string;
  business_address: string;
  report_source: string;
  confidence_level?: number; // Deprecated: use tfw_ratio
  tfw_ratio?: 'few' | 'many' | 'most' | 'all';
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
    onSuccess: async (data, variables) => {
      // Invalidate general reports queries to refetch data
      await queryClient.invalidateQueries({ queryKey: ["reports"] });
      await queryClient.invalidateQueries({
        queryKey: ["reportsGroupedByAddress"],
      });
      
      // Invalidate the specific business address query if we have an address
      if (variables.business_address) {
        await queryClient.invalidateQueries({
          queryKey: ["addressReports", variables.business_address],
        });
      }
    },
  });
}

export function useReportsGroupedByAddress(filters: ReportFilters) {
  const apiClient = useApiClient();

  return useQuery({
    queryKey: ["reportsGroupedByAddress", filters],
    queryFn: async (): Promise<ReportsByAddressResponse> => {
      const params = new URLSearchParams();

      if (filters.query) params.append("query", filters.query);
      if (filters.city) params.append("city", filters.city);
      if (filters.province) params.append("province", filters.province);
      if (filters.year) params.append("year", filters.year);

      if (filters.limit) params.append("limit", filters.limit.toString());
      if (filters.offset) params.append("offset", filters.offset.toString());

      const response = await apiClient.get<ReportsByAddressResponse>(
        `/reports/grouped-by-address?${params}`,
      );
      return response.data;
    },
    staleTime: 1000 * 60 * 5, // 5 minutes
  });
}

export function useAddressReports(address: string) {
  const apiClient = useApiClient();

  return useQuery({
    queryKey: ["addressReports", address],
    queryFn: async (): Promise<ReportsResponse> => {
      const params = new URLSearchParams();
      params.append("address", address);

      const response = await apiClient.get<ReportsResponse>(
        `/reports/address?${params}`,
      );
      return response.data;
    },
    enabled: !!address, // Only run query if address is provided
    staleTime: 1000 * 60 * 5, // 5 minutes
  });
}


// Admin hooks for report management
export function useUpdateReport() {
  const apiClient = useApiClient();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ id, data }: { id: string; data: UpdateReportRequest }): Promise<Report> => {
      const response = await apiClient.put<Report>(`/reports/${id}`, data);
      return response.data;
    },
    onSuccess: async () => {
      // Invalidate all report queries to refetch data
      await queryClient.invalidateQueries({ queryKey: ["reports"] });
      await queryClient.invalidateQueries({ queryKey: ["reportsGroupedByAddress"] });
      await queryClient.invalidateQueries({ queryKey: ["addressReports"] });
    },
  });
}

export function useDeleteReport() {
  const apiClient = useApiClient();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (id: string): Promise<void> => {
      await apiClient.delete(`/reports/${id}`);
    },
    onSuccess: async () => {
      // Invalidate all report queries to refetch data
      await queryClient.invalidateQueries({ queryKey: ["reports"] });
      await queryClient.invalidateQueries({ queryKey: ["reportsGroupedByAddress"] });
      await queryClient.invalidateQueries({ queryKey: ["addressReports"] });
    },
  });
}

export function useReportStats() {
  const apiClient = useApiClient();

  return useQuery({
    queryKey: ["reports", "stats"],
    queryFn: async (): Promise<{ total_reports: number }> => {
      // Get reports with minimal data to get total count
      const response = await apiClient.get<ReportsResponse>("/reports?limit=1");
      return {
        total_reports: response.data.reports.length || 0,
      };
    },
    staleTime: 1000 * 60 * 5, // 5 minutes
  });
}

export type {
  Report,
  ReportsResponse,
  ReportFilters,
  CreateReportRequest,
  UpdateReportRequest,
  ReportsByAddress,
  ReportsByAddressResponse,
};
