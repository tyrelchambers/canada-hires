import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useApiClient } from "./useApiClient";
import {
  Report,
  CreateReportRequest,
  UpdateReportRequest,
  ReportListResponse,
  ReportFilters,
  ModerationRequest,
} from "@/types";

// Query hooks for fetching reports
export const useReports = (filters: ReportFilters = {}) => {
  const api = useApiClient();
  return useQuery({
    queryKey: ["reports", filters],
    queryFn: async (): Promise<ReportListResponse> => {
      const params = new URLSearchParams();
      
      if (filters.limit) params.append("limit", filters.limit.toString());
      if (filters.offset) params.append("offset", filters.offset.toString());

      const response = await api.get<ReportListResponse>(`/reports?${params}`);
      return response.data;
    },
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
};

export const useReport = (id: string) => {
  const api = useApiClient();
  return useQuery({
    queryKey: ["reports", id],
    queryFn: async (): Promise<Report> => {
      const response = await api.get<Report>(`/reports/${id}`);
      return response.data;
    },
    staleTime: 5 * 60 * 1000, // 5 minutes
    enabled: !!id,
  });
};

export const useBusinessReports = (businessName: string, filters: ReportFilters = {}) => {
  const api = useApiClient();
  return useQuery({
    queryKey: ["reports", "business", businessName, filters],
    queryFn: async (): Promise<ReportListResponse> => {
      const params = new URLSearchParams();
      
      if (filters.limit) params.append("limit", filters.limit.toString());
      if (filters.offset) params.append("offset", filters.offset.toString());

      const response = await api.get<ReportListResponse>(`/reports/business/${encodeURIComponent(businessName)}?${params}`);
      return response.data;
    },
    staleTime: 5 * 60 * 1000, // 5 minutes
    enabled: !!businessName,
  });
};

export const useUserReports = (filters: ReportFilters = {}) => {
  const api = useApiClient();
  return useQuery({
    queryKey: ["reports", "user", filters],
    queryFn: async (): Promise<ReportListResponse> => {
      const params = new URLSearchParams();
      
      if (filters.limit) params.append("limit", filters.limit.toString());
      if (filters.offset) params.append("offset", filters.offset.toString());

      const response = await api.get<ReportListResponse>(`/reports/user/me?${params}`);
      return response.data;
    },
    staleTime: 1 * 60 * 1000, // 1 minute
  });
};

export const useReportsByStatus = (status: "pending" | "approved" | "rejected" | "flagged", filters: ReportFilters = {}) => {
  const api = useApiClient();
  return useQuery({
    queryKey: ["reports", "status", status, filters],
    queryFn: async (): Promise<ReportListResponse> => {
      const params = new URLSearchParams();
      
      if (filters.limit) params.append("limit", filters.limit.toString());
      if (filters.offset) params.append("offset", filters.offset.toString());

      const response = await api.get<ReportListResponse>(`/reports/status/${status}?${params}`);
      return response.data;
    },
    staleTime: 30 * 1000, // 30 seconds
    enabled: !!status,
  });
};

// Mutation hooks for creating/updating/deleting reports
export const useCreateReport = () => {
  const api = useApiClient();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (data: CreateReportRequest): Promise<Report> => {
      const response = await api.post<Report>("/reports", data);
      return response.data;
    },
    onSuccess: async () => {
      // Invalidate relevant queries
      await queryClient.invalidateQueries({ queryKey: ["reports"] });
      await queryClient.invalidateQueries({ queryKey: ["reports", "user"] });
    },
  });
};

export const useUpdateReport = () => {
  const api = useApiClient();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ id, data }: { id: string; data: UpdateReportRequest }): Promise<Report> => {
      const response = await api.put<Report>(`/reports/${id}`, data);
      return response.data;
    },
    onSuccess: async (updatedReport) => {
      // Invalidate and update relevant queries
      await queryClient.invalidateQueries({ queryKey: ["reports"] });
      await queryClient.invalidateQueries({ queryKey: ["reports", "user"] });
      await queryClient.invalidateQueries({ queryKey: ["reports", "business"] });
      
      // Update the specific report in cache
      queryClient.setQueryData(["reports", updatedReport.id], updatedReport);
    },
  });
};

export const useDeleteReport = () => {
  const api = useApiClient();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (id: string): Promise<void> => {
      await api.delete(`/reports/${id}`);
    },
    onSuccess: async (_, deletedId) => {
      // Invalidate relevant queries
      await queryClient.invalidateQueries({ queryKey: ["reports"] });
      await queryClient.invalidateQueries({ queryKey: ["reports", "user"] });
      await queryClient.invalidateQueries({ queryKey: ["reports", "business"] });
      
      // Remove the specific report from cache
      queryClient.removeQueries({ queryKey: ["reports", deletedId] });
    },
  });
};

// Admin mutation hooks for moderation
export const useApproveReport = () => {
  const api = useApiClient();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ id, notes }: { id: string; notes?: string }) => {
      const data: ModerationRequest = { notes };
      return api.post(`/reports/${id}/approve`, data);
    },
    onSuccess: async () => {
      // Invalidate admin queries
      await queryClient.invalidateQueries({ queryKey: ["reports", "status"] });
      await queryClient.invalidateQueries({ queryKey: ["reports"] });
    },
  });
};

export const useRejectReport = () => {
  const api = useApiClient();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ id, notes }: { id: string; notes?: string }) => {
      const data: ModerationRequest = { notes };
      return api.post(`/reports/${id}/reject`, data);
    },
    onSuccess: async () => {
      // Invalidate admin queries
      await queryClient.invalidateQueries({ queryKey: ["reports", "status"] });
      await queryClient.invalidateQueries({ queryKey: ["reports"] });
    },
  });
};

export const useFlagReport = () => {
  const api = useApiClient();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ id, notes }: { id: string; notes?: string }) => {
      const data: ModerationRequest = { notes };
      return api.post(`/reports/${id}/flag`, data);
    },
    onSuccess: async () => {
      // Invalidate admin queries
      await queryClient.invalidateQueries({ queryKey: ["reports", "status"] });
      await queryClient.invalidateQueries({ queryKey: ["reports"] });
    },
  });
};