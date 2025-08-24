import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useApiClient } from "./useApiClient";
import {
  NonCompliantReason,
  NonCompliantLocationResponse,
  NonCompliantEmployersByPostalCodeResponse,
} from "@/types";

export interface NonCompliantReasonsResponse {
  reasons: NonCompliantReason[];
  count: number;
}

export function useNonCompliantReasons() {
  const apiClient = useApiClient();

  return useQuery<NonCompliantReasonsResponse>({
    queryKey: ["non-compliant", "reasons"],
    queryFn: async (): Promise<NonCompliantReasonsResponse> => {
      const response = await apiClient.get<NonCompliantReasonsResponse>(
        "/non-compliant/reasons",
      );
      return response.data;
    },
  });
}

export function useNonCompliantLocations(limit: number = 1000) {
  const apiClient = useApiClient();

  return useQuery<NonCompliantLocationResponse>({
    queryKey: ["non-compliant", "locations", limit],
    queryFn: async (): Promise<NonCompliantLocationResponse> => {
      const response = await apiClient.get<NonCompliantLocationResponse>(
        "/non-compliant/locations",
        {
          params: { limit },
        },
      );
      return response.data;
    },
  });
}

export function useNonCompliantByPostalCode(
  postalCode: string,
  limit: number = 100,
) {
  const apiClient = useApiClient();

  return useQuery<NonCompliantEmployersByPostalCodeResponse>({
    queryKey: ["non-compliant", "postal-code", postalCode, limit],
    queryFn: async (): Promise<NonCompliantEmployersByPostalCodeResponse> => {
      const response =
        await apiClient.get<NonCompliantEmployersByPostalCodeResponse>(
          `/non-compliant/employers/postal-code/${postalCode}`,
          {
            params: { limit },
          },
        );
      return response.data;
    },
    enabled: !!postalCode,
  });
}

export function useNonCompliantByCoordinates(
  lat: number | null,
  lng: number | null,
  limit: number = 100,
) {
  const apiClient = useApiClient();

  return useQuery<NonCompliantEmployersByPostalCodeResponse>({
    queryKey: ["non-compliant", "coordinates", lat, lng, limit],
    queryFn: async (): Promise<NonCompliantEmployersByPostalCodeResponse> => {
      if (lat === null || lng === null) {
        throw new Error("Both lat and lng coordinates are required");
      }
      const response =
        await apiClient.get<NonCompliantEmployersByPostalCodeResponse>(
          `/non-compliant/employers/coordinates/${lat}/${lng}`,
          {
            params: { limit },
          },
        );
      return response.data;
    },
    enabled: lat !== null && lng !== null,
  });
}

export function useTriggerNonCompliantGeocoding() {
  const apiClient = useApiClient();
  const queryClient = useQueryClient();

  return useMutation<{ message: string }, Error, void>({
    mutationFn: async () => {
      const response = await apiClient.post<{ message: string }>(
        "/admin/non-compliant/geocode",
      );
      return response.data;
    },
    onSuccess: () => {
      // Invalidate related queries to refresh data
      void queryClient.invalidateQueries({
        queryKey: ["non-compliant", "locations"],
      });
      void queryClient.invalidateQueries({
        queryKey: ["non-compliant", "employers"],
      });
    },
  });
}
