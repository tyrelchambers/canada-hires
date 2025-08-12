import { useQuery } from "@tanstack/react-query";
import { useApiClient } from "./useApiClient";

export interface RegionData {
  name: string;
  count: number;
}

export interface LMIAStatistics {
  id: string;
  date: string;
  period_type: "daily" | "monthly";
  total_jobs: number;
  unique_employers: number;
  avg_salary_min?: number;
  avg_salary_max?: number;
  top_provinces: RegionData[];
  top_cities: RegionData[];
  created_at: string;
  updated_at: string;
}

export interface TrendsSummary {
  total_jobs_today: number;
  total_jobs_this_month: number;
  total_jobs_last_month: number;
  percentage_change: number;
  top_provinces_today: RegionData[];
  top_cities_today: RegionData[];
  recent_trends: LMIAStatistics[];
}

interface TrendsResponse {
  data: LMIAStatistics[];
  count: number;
}

interface TrendsFilters {
  start_date?: string;
  end_date?: string;
  limit?: number;
}

export const useDailyTrends = (filters?: TrendsFilters) => {
  const apiClient = useApiClient();
  return useQuery({
    queryKey: ["lmia-daily-trends", filters],
    queryFn: async (): Promise<TrendsResponse> => {
      const params = new URLSearchParams();

      if (filters?.start_date) params.append("start_date", filters.start_date);
      if (filters?.end_date) params.append("end_date", filters.end_date);
      if (filters?.limit) params.append("limit", filters.limit.toString());

      const response = await apiClient.get<TrendsResponse>(
        `/lmia/statistics/daily?${params}`,
      );
      return response.data;
    },
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
};

export const useMonthlyTrends = (filters?: TrendsFilters) => {
  const apiClient = useApiClient();

  return useQuery({
    queryKey: ["lmia-monthly-trends", filters],
    queryFn: async (): Promise<TrendsResponse> => {
      const params = new URLSearchParams();

      if (filters?.start_date) params.append("start_date", filters.start_date);
      if (filters?.end_date) params.append("end_date", filters.end_date);
      if (filters?.limit) params.append("limit", filters.limit.toString());

      const response = await apiClient.get<TrendsResponse>(
        `/lmia/statistics/monthly?${params}`,
      );
      return response.data;
    },
    staleTime: 10 * 60 * 1000, // 10 minutes
  });
};

export const useTrendsSummary = () => {
  const apiClient = useApiClient();

  return useQuery({
    queryKey: ["lmia-trends-summary"],
    queryFn: async (): Promise<TrendsSummary> => {
      const response = await apiClient.get<TrendsSummary>(
        "/lmia/statistics/summary",
      );
      return response.data;
    },
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
};

// Utility hook to get trends for a specific time period
export const useTrendsForPeriod = (
  period: "week" | "month" | "quarter" | "year",
  type: "daily" | "monthly" = "daily",
) => {
  const getDateRange = () => {
    const now = new Date();
    const endDate = now.toISOString().split("T")[0]; // Today

    let startDate: Date;
    switch (period) {
      case "week":
        startDate = new Date(now.getTime() - 7 * 24 * 60 * 60 * 1000);
        break;
      case "month":
        startDate = new Date(
          now.getFullYear(),
          now.getMonth() - 1,
          now.getDate(),
        );
        break;
      case "quarter":
        startDate = new Date(
          now.getFullYear(),
          now.getMonth() - 3,
          now.getDate(),
        );
        break;
      case "year":
        startDate = new Date(
          now.getFullYear() - 1,
          now.getMonth(),
          now.getDate(),
        );
        break;
    }

    return {
      start_date: startDate.toISOString().split("T")[0],
      end_date: endDate,
    };
  };

  const dateRange = getDateRange();

  // Always call both hooks but only return the one we need
  // This ensures we don't violate the Rules of Hooks
  const dailyResult = useDailyTrends(type === "daily" ? dateRange : undefined);
  const monthlyResult = useMonthlyTrends(
    type === "monthly" ? dateRange : undefined,
  );

  return type === "daily" ? dailyResult : monthlyResult;
};
