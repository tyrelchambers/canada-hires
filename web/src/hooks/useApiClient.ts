import { useMemo } from "react";
import axios, { AxiosError, AxiosInstance } from "axios";
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
    client.interceptors.request.use((config) => {
      return config;
    });

    // Add response interceptor to handle 401 errors
    client.interceptors.response.use(
      (response) => response,
      (error: AxiosError) => {
        if (error.response?.status === 401) {
          // Redirect to login on 401 Unauthorized
          window.location.href = "/auth/login";
        }
        return Promise.reject(error);
      }
    );

    return client;
  }, []);

  return apiClient;
}
