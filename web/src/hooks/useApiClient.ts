import { useMemo } from "react";
import axios, { AxiosInstance } from "axios";
import { API_BASE_URL } from "@/constants";

export function useApiClient(): AxiosInstance {
  const apiClient = useMemo(() => {
    const client = axios.create({
      baseURL: API_BASE_URL,
      headers: {
        "Content-Type": "application/json",
      },
      withCredentials: true,
    });

    // Add Clerk token to every request
    client.interceptors.request.use(async (config) => {
      return config;
    });

    return client;
  }, []);

  return apiClient;
}
