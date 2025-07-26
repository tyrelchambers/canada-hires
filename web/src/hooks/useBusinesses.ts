import { useQuery } from "@tanstack/react-query";
import { useApiClient } from "./useApiClient";

interface BusinessRating {
  business_id: string;
  current_rating: string;
  confidence_score: number;
  report_count: number;
  avg_tfw_percentage?: number;
  last_updated: string;
  is_disputed: boolean;
}

interface Business {
  id: string;
  name: string;
  address?: string;
  city?: string;
  province?: string;
  postal_code?: string;
  industry_code?: string;
  phone?: string;
  website?: string;
  size_category?: string;
  created_at: string;
  updated_at: string;
  rating?: BusinessRating;
}

interface DirectoryResponse {
  businesses: Business[];
  total: number;
  limit: number;
  offset: number;
}

interface BusinessFilters {
  query?: string;
  city?: string;
  province?: string;
  rating?: string;
  year?: string;
  limit?: number;
  offset?: number;
}

export function useBusinesses(filters: BusinessFilters) {
  const apiClient = useApiClient();

  return useQuery({
    queryKey: ["businesses", filters],
    queryFn: async (): Promise<DirectoryResponse> => {
      const params = new URLSearchParams();
      
      if (filters.query) params.append("query", filters.query);
      if (filters.city) params.append("city", filters.city);
      if (filters.province) params.append("province", filters.province);
      if (filters.rating) params.append("rating", filters.rating);
      if (filters.year) params.append("year", filters.year);
      
      // If no search criteria, default to latest year
      if (
        !filters.query &&
        !filters.city &&
        !filters.province &&
        !filters.rating &&
        !filters.year
      ) {
        params.append("year", new Date().getFullYear().toString());
      }
      
      if (filters.limit) params.append("limit", filters.limit.toString());
      if (filters.offset) params.append("offset", filters.offset.toString());

      const response = await apiClient.get(`/directory?${params}`);
      return response.data;
    },
    staleTime: 1000 * 60 * 5, // 5 minutes
  });
}

export type { Business, BusinessRating, DirectoryResponse, BusinessFilters };