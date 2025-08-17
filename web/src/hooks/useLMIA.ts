import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useApiClient } from "./useApiClient";
import {
  CronJob,
  LMIASearchResponse,
  LMIALocationResponse,
  LMIAResourcesResponse,
  LMIAStatsResponse,
  LMIAEmployersByResourceResponse,
} from "@/types";

interface LMIASearchParams {
  q: string;
  year?: string;
  limit?: number;
}

interface LMIALocationParams {
  city?: string;
  province?: string;
  year?: string;
  limit?: number;
}

export function useLMIASearch(
  query: string,
  year?: string,
  limit: number = 0, // Default to no limit (fetch all records)
) {
  const apiClient = useApiClient();

  return useQuery<LMIASearchResponse>({
    queryKey: ["lmia", "search", query, year, limit],
    queryFn: async (): Promise<LMIASearchResponse> => {
      // If no query but year is provided, use wildcard search for that year
      if ((!query || query.trim().length < 2) && year) {
        const params: LMIASearchParams = { q: "*", year };
        if (limit > 0) params.limit = limit;
        const response = await apiClient.get<LMIASearchResponse>(
          `/lmia/employers/search`,
          {
            params,
          },
        );
        return response.data;
      }

      // If no query and no year, return empty result
      if (!query || query.trim().length < 2) {
        return { employers: [], count: 0, query: "", limit: 0 };
      }

      const params: LMIASearchParams = { q: query };
      if (year) params.year = year;
      if (limit > 0) params.limit = limit;

      const response = await apiClient.get<LMIASearchResponse>(
        `/lmia/employers/search`,
        {
          params,
        },
      );
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

  return useQuery<LMIALocationResponse>({
    queryKey: ["lmia", "location", city, province, year, limit],
    queryFn: async (): Promise<LMIALocationResponse> => {
      // Show latest year by default if no search criteria - use search endpoint with wildcard query
      const shouldShowDefault = !city && !province && !year;
      if (shouldShowDefault) {
        const currentYear = new Date().getFullYear();
        const params: LMIASearchParams = {
          q: "*",
          year: currentYear.toString(),
        };
        if (limit > 0) params.limit = limit;
        const response = await apiClient.get<LMIALocationResponse>(
          `/lmia/employers/search`,
          {
            params,
          },
        );
        return response.data;
      }

      // If only year is provided (no city/province), use search endpoint
      if (!city && !province && year) {
        const params: LMIASearchParams = { q: "*", year };
        if (limit > 0) params.limit = limit;
        const response = await apiClient.get<LMIALocationResponse>(
          `/lmia/employers/search`,
          {
            params,
          },
        );
        return response.data;
      }

      // If we have city or province, use location endpoint
      if (city || province) {
        const params: LMIALocationParams = {};
        if (city) params.city = city;
        if (province) params.province = province;
        if (year) params.year = year;
        if (limit > 0) params.limit = limit;

        const response = await apiClient.get<LMIALocationResponse>(
          `/lmia/employers/location`,
          {
            params,
          },
        );
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

  return useQuery<LMIAResourcesResponse>({
    queryKey: ["lmia", "resources"],
    queryFn: async (): Promise<LMIAResourcesResponse> => {
      const response = await apiClient.get<LMIAResourcesResponse>(
        "/lmia/resources",
      );
      return response.data;
    },
  });
}

export function useLMIAStats() {
  const apiClient = useApiClient();

  return useQuery<LMIAStatsResponse>({
    queryKey: ["lmia", "stats"],
    queryFn: async (): Promise<LMIAStatsResponse> => {
      const response = await apiClient.get<LMIAStatsResponse>("/lmia/stats");
      return response.data;
    },
  });
}

export function useLMIAStatus() {
  const apiClient = useApiClient();

  return useQuery<CronJob>({
    queryKey: ["lmia", "status"],
    queryFn: async (): Promise<CronJob> => {
      const response = await apiClient.get<CronJob>("/lmia/status");
      return response.data;
    },
  });
}

export function useTriggerLMIAUpdate() {
  const apiClient = useApiClient();
  const queryClient = useQueryClient();

  return useMutation<void, Error, void>({
    mutationFn: async () => {
      await apiClient.post("/lmia/update");
    },
    onSuccess: () => {
      // Invalidate related queries to refresh data
      void queryClient.invalidateQueries({ queryKey: ["lmia", "stats"] });
      void queryClient.invalidateQueries({ queryKey: ["lmia", "status"] });
      void queryClient.invalidateQueries({ queryKey: ["lmia", "resources"] });
    },
  });
}

export function useLMIAEmployersByResource(resourceId: string) {
  const apiClient = useApiClient();

  return useQuery<LMIAEmployersByResourceResponse>({
    queryKey: ["lmia", "employers", "resource", resourceId],
    queryFn: async (): Promise<LMIAEmployersByResourceResponse> => {
      const response = await apiClient.get<LMIAEmployersByResourceResponse>(
        `/lmia/employers/resource/${resourceId}`,
      );
      return response.data;
    },
    enabled: !!resourceId,
  });
}
