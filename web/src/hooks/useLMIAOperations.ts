import { useMutation } from '@tanstack/react-query';
import { useApiClient } from './useApiClient';

interface ApiResponse {
  message: string;
  status: string;
}

export function useLMIAOperations() {
  const apiClient = useApiClient();

  const fullUpdateMutation = useMutation({
    mutationFn: async (): Promise<ApiResponse> => {
      const response = await apiClient.post('/lmia/update');
      return response.data as ApiResponse;
    },
  });

  const processorMutation = useMutation({
    mutationFn: async (): Promise<ApiResponse> => {
      const response = await apiClient.post('/lmia/process');
      return response.data as ApiResponse;
    },
  });

  const backfillMutation = useMutation({
    mutationFn: async (): Promise<ApiResponse> => {
      const response = await apiClient.post('/lmia/statistics/backfill');
      return response.data as ApiResponse;
    },
  });

  const geocodingMutation = useMutation({
    mutationFn: async (): Promise<ApiResponse> => {
      const response = await apiClient.post('/admin/lmia/geocode');
      return response.data as ApiResponse;
    },
  });

  return {
    fullUpdate: fullUpdateMutation,
    processor: processorMutation,
    backfill: backfillMutation,
    geocoding: geocodingMutation,
  };
}

export function useScraperOperations() {
  const apiClient = useApiClient();

  const scraperMutation = useMutation({
    mutationFn: async (): Promise<ApiResponse> => {
      const response = await apiClient.post('/admin/scraper/run');
      return response.data as ApiResponse;
    },
  });

  const statisticsMutation = useMutation({
    mutationFn: async (): Promise<ApiResponse> => {
      const response = await apiClient.post('/admin/scraper/statistics');
      return response.data as ApiResponse;
    },
  });

  return {
    scraper: scraperMutation,
    statistics: statisticsMutation,
  };
}