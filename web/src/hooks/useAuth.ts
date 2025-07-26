import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useApiClient } from "./useApiClient";

// React Query hooks
export const useCurrentUser = () => {
  const api = useApiClient();
  return useQuery({
    queryKey: ["profile"],
    queryFn: async () => {
      const response = await api.get("/user/profile");
      return response.data;
    },
    retry: false,
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
};

export const useSendLoginLink = () => {
  const api = useApiClient();
  return useMutation({
    mutationFn: (email: string) => api.post("/auth/send-login-link", { email }),
  });
};

export const useLogout = () => {
  const queryClient = useQueryClient();
  const api = useApiClient();
  return useMutation({
    mutationFn: () => api.post("/auth/logout"),
    onSuccess: () => {
      // Clear all queries when user logs out
      queryClient.clear();
    },
  });
};
