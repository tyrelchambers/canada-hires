import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useApiClient } from "./useApiClient";
import {
  Subreddit,
  SubredditsResponse,
  CreateSubredditRequest,
  UpdateSubredditRequest,
} from "@/types";

// Get all subreddits (admin only)
export const useSubreddits = () => {
  const api = useApiClient();
  return useQuery({
    queryKey: ["subreddits"],
    queryFn: async (): Promise<SubredditsResponse> => {
      const response = await api.get<SubredditsResponse>("/subreddits");
      return response.data;
    },
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
};

// Get only active subreddits (public)
export const useActiveSubreddits = () => {
  const api = useApiClient();
  return useQuery({
    queryKey: ["subreddits", "active"],
    queryFn: async (): Promise<SubredditsResponse> => {
      const response = await api.get<SubredditsResponse>("/subreddits/active");
      return response.data;
    },
    staleTime: 10 * 60 * 1000, // 10 minutes
  });
};

// Create a new subreddit
export const useCreateSubreddit = () => {
  const api = useApiClient();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (data: CreateSubredditRequest): Promise<Subreddit> => {
      const response = await api.post<Subreddit>("/subreddits", data);
      return response.data;
    },
    onSuccess: async () => {
      // Invalidate subreddit queries to refresh the lists
      await queryClient.invalidateQueries({
        queryKey: ["subreddits"],
      });
    },
  });
};

// Update a subreddit
export const useUpdateSubreddit = () => {
  const api = useApiClient();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({
      id,
      data,
    }: {
      id: string;
      data: UpdateSubredditRequest;
    }): Promise<Subreddit> => {
      const response = await api.put<Subreddit>(`/subreddits/${id}`, data);
      return response.data;
    },
    onSuccess: async () => {
      // Invalidate subreddit queries to refresh the lists
      await queryClient.invalidateQueries({
        queryKey: ["subreddits"],
      });
    },
  });
};

// Delete a subreddit
export const useDeleteSubreddit = () => {
  const api = useApiClient();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (id: string): Promise<void> => {
      await api.delete(`/subreddits/${id}`);
    },
    onSuccess: async () => {
      // Invalidate subreddit queries to refresh the lists
      await queryClient.invalidateQueries({
        queryKey: ["subreddits"],
      });
    },
  });
};