import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useApiClient } from "@/hooks/useApiClient";

export interface Boycott {
  id: string;
  business_name: string;
  business_address: string;
  created_at: string;
  is_boycotting: boolean;
}

export interface BoycottStats {
  business_name: string;
  business_address: string;
  boycott_count: number;
}

export interface BoycottStatsDetail {
  business_name: string;
  business_address: string;
  boycott_count: number;
  is_boycotted_by_user: boolean;
}

interface ToggleBoycottRequest {
  business_name: string;
  business_address: string;
}

interface BoycottListResponse {
  data: Boycott[];
  limit: number;
  offset: number;
  count: number;
}

// Toggle boycott status
export const useToggleBoycott = () => {
  const queryClient = useQueryClient();
  const api = useApiClient();

  return useMutation({
    mutationFn: async (data: ToggleBoycottRequest): Promise<Boycott> => {
      return await api
        .post<Boycott>("/boycotts/toggle", data)
        .then((res) => res.data);
    },
    onSuccess: async () => {
      // Invalidate and refetch related queries
      await queryClient.invalidateQueries({ queryKey: ["boycotts"] });
      await queryClient.invalidateQueries({ queryKey: ["boycott-stats"] });
      await queryClient.invalidateQueries({ queryKey: ["top-boycotted"] });
    },
  });
};

// Get user's boycotts
export const useUserBoycotts = (limit = 50, offset = 0) => {
  const api = useApiClient();

  return useQuery({
    queryKey: ["boycotts", "user", limit, offset],
    queryFn: async () => {
      const response = await api.get<BoycottListResponse>(
        `/boycotts/my?limit=${limit}&offset=${offset}`,
      );
      return response.data;
    },
  });
};

// Get top boycotted businesses
export const useTopBoycotted = (limit = 3) => {
  const api = useApiClient();

  return useQuery({
    queryKey: ["top-boycotted", limit],
    queryFn: async () => {
      const response = await api.get<BoycottStats[]>(
        `/boycotts/top?limit=${limit}`,
      );
      return response.data;
    },
  });
};

// Get boycott stats for a specific business
export const useBoycottStats = (
  businessName: string,
  businessAddress: string,
) => {
  const api = useApiClient();

  return useQuery({
    queryKey: ["boycott-stats", businessName, businessAddress],
    queryFn: async () => {
      const params = new URLSearchParams({
        business_name: businessName,
        business_address: businessAddress,
      });
      const response = await api.get<BoycottStatsDetail>(
        `/boycotts/stats?${params.toString()}`,
      );
      return response.data;
    },
    enabled: !!businessName, // Only run the query if businessName is provided
  });
};
