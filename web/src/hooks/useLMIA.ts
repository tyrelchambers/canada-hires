import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useApiClient } from "./useApiClient";
import {
  CronJob,
  LMIASearchResponse,
  LMIALocationResponse,
  LMIAResourcesResponse,
  LMIAStatsResponse,
} from "@/types";

export function useLMIASearch(
  query: string,
  year?: string,
  limit: number = 0, // Default to no limit (fetch all records)
) {
  const apiClient = useApiClient();

  return useQuery({
    queryKey: ["lmia", "search", query, year, limit],
    queryFn: async (): Promise<LMIASearchResponse> => {
      // If no query but year is provided, use wildcard search for that year
      if ((!query || query.trim().length < 2) && year) {
        const params: any = { q: "*", year };
        if (limit > 0) params.limit = limit;
        const response = await apiClient.get(`/lmia/employers/search`, {
          params,
        });
        return response.data;
      }

      // If no query and no year, return empty result
      if (!query || query.trim().length < 2) {
        return { employers: [], count: 0, query: "", limit: 0 };
      }

      const params: any = { q: query };
      if (year) params.year = year;
      if (limit > 0) params.limit = limit;

      const response = await apiClient.get(`/lmia/employers/search`, {
        params,
      });
      return response.data;
    },
    enabled:
      (!!query && query.trim().length >= 2) ||
      (!!year && (!query || query.trim().length < 2)),
  });
}

export function useLMIALocation(
  city: string,
  province: string,
  year?: string,
  limit: number = 0, // Default to no limit (fetch all records)
) {
  const apiClient = useApiClient();

  return useQuery({
    queryKey: ["lmia", "location", city, province, year, limit],
    queryFn: async (): Promise<LMIALocationResponse> => {
      // Show latest year by default if no search criteria - use search endpoint with wildcard query
      const shouldShowDefault = !city && !province && !year;
      if (shouldShowDefault) {
        const currentYear = new Date().getFullYear();
        const params: any = {
          q: "*",
          year: currentYear.toString(),
        };
        if (limit > 0) params.limit = limit;
        const response = await apiClient.get(`/lmia/employers/search`, {
          params,
        });
        return response.data;
      }

      // If only year is provided (no city/province), use search endpoint
      if (!city && !province && year) {
        const params: any = { q: "*", year };
        if (limit > 0) params.limit = limit;
        const response = await apiClient.get(`/lmia/employers/search`, {
          params,
        });
        return response.data;
      }

      // If we have city or province, use location endpoint
      if (city || province) {
        const params: any = {};
        if (city) params.city = city;
        if (province) params.province = province;
        if (year) params.year = year;
        if (limit > 0) params.limit = limit;

        const response = await apiClient.get(`/lmia/employers/location`, {
          params,
        });
        return response.data;
      }

      // Fallback case - should not reach here but return empty result
      return { employers: [], count: 0, city: "", province: "", limit: 0 };
    },
    enabled: !!(city || province || year) || (!city && !province && !year), // Enable for default case too
  });
}

export function useLMIAResources() {
  const apiClient = useApiClient();

  return useQuery({
    queryKey: ["lmia", "resources"],
    queryFn: async (): Promise<LMIAResourcesResponse> => {
      const response = await apiClient.get("/lmia/resources");
      return response.data;
    },
  });
}

export function useLMIAStats() {
  const apiClient = useApiClient();

  return useQuery({
    queryKey: ["lmia", "stats"],
    queryFn: async (): Promise<LMIAStatsResponse> => {
      const response = await apiClient.get("/lmia/stats");
      return response.data;
    },
  });
}

export function useLMIAStatus() {
  const apiClient = useApiClient();

  return useQuery({
    queryKey: ["lmia", "status"],
    queryFn: async (): Promise<CronJob> => {
      const response = await apiClient.get("/lmia/status");
      return response.data;
    },
  });
}

export function useTriggerLMIAUpdate() {
  const apiClient = useApiClient();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async () => {
      const response = await apiClient.post("/lmia/update");
      return response.data;
    },
    onSuccess: () => {
      // Invalidate related queries to refresh data
      queryClient.invalidateQueries({ queryKey: ["lmia", "stats"] });
      queryClient.invalidateQueries({ queryKey: ["lmia", "status"] });
      queryClient.invalidateQueries({ queryKey: ["lmia", "resources"] });
    },
  });
}

export function useLMIAEmployersByResource(resourceId: string) {
  const apiClient = useApiClient();

  return useQuery({
    queryKey: ["lmia", "employers", "resource", resourceId],
    queryFn: async () => {
      const response = await apiClient.get(
        `/lmia/employers/resource/${resourceId}`,
      );
      return response.data;
    },
    enabled: !!resourceId,
  });
}
