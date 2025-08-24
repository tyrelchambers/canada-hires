import { useQuery } from "@tanstack/react-query";
import { useApiClient } from "./useApiClient";

// Pelias geocoding response interfaces
interface PeliasGeometry {
  type: string;
  coordinates: [number, number]; // [longitude, latitude]
}

interface PeliasProperties {
  id: string;
  gid: string;
  layer: string;
  source: string;
  source_id: string;
  country_code?: string;
  name: string;
  postalcode?: string;
  confidence: number;
  match_type?: string;
  distance?: number;
  accuracy?: string;
  country?: string;
  country_gid?: string;
  country_a?: string;
  region?: string;
  region_gid?: string;
  region_a?: string;
  locality?: string;
  locality_gid?: string;
  label?: string;
}

interface PeliasFeature {
  type: string;
  geometry: PeliasGeometry;
  properties: PeliasProperties;
  bbox?: number[];
}

interface PeliasResponse {
  geocoding: {
    version: string;
    attribution: string;
    query: Record<string, any>;
    warnings?: string[];
    errors?: string[];
    engine: Record<string, any>;
    timestamp: number;
  };
  type: string;
  features: PeliasFeature[];
  bbox?: number[];
}

export function useAddressSearch(
  query: string,
  options: {
    size?: number;
    layers?: string;
    enabled?: boolean;
  } = {}
) {
  const apiClient = useApiClient();
  const { size = 5, layers, enabled = true } = options;

  return useQuery<PeliasResponse>({
    queryKey: ["search", "address", query, size, layers],
    queryFn: async (): Promise<PeliasResponse> => {
      const params: Record<string, any> = {
        text: query,
        size,
      };
      
      if (layers) {
        params.layers = layers;
      }

      const response = await apiClient.get<PeliasResponse>("/search", {
        params,
      });
      return response.data;
    },
    enabled: enabled && query.length >= 3,
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
}

export function useCitySearch(
  query: string,
  options: {
    size?: number;
    enabled?: boolean;
  } = {}
) {
  const { size = 8, enabled = true } = options;
  
  return useAddressSearch(query, {
    size,
    layers: "locality,region",
    enabled,
  });
}

export type { PeliasFeature, PeliasResponse };