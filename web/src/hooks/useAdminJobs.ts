import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useApiClient } from "./useApiClient";
import {
  JobPostingsResponse,
  RedditApprovalRequest,
  RedditRejectionRequest,
  BulkApprovalRequest,
  BulkRejectionRequest,
  RedditApprovalStats,
  GeneratedRedditPost,
  BulkGenerationResponse,
  BulkGenerationRequest,
  RedditPreview,
} from "@/types";

// Get pending jobs for Reddit approval
export const usePendingJobs = (limit = 50, offset = 0) => {
  const api = useApiClient();
  return useQuery({
    queryKey: ["admin", "jobs", "reddit", "pending", { limit, offset }],
    queryFn: async (): Promise<JobPostingsResponse> => {
      const response = await api.get<JobPostingsResponse>(
        "/admin/jobs/reddit/pending",
        {
          params: { limit, offset },
        },
      );
      return response.data;
    },
    staleTime: 30 * 1000, // 30 seconds
  });
};

// Get posted jobs (approved and posted to Reddit)
export const usePostedJobs = (limit = 50, offset = 0) => {
  const api = useApiClient();
  return useQuery({
    queryKey: ["admin", "jobs", "reddit", "posted", { limit, offset }],
    queryFn: async (): Promise<JobPostingsResponse> => {
      const response = await api.get<JobPostingsResponse>(
        "/admin/jobs/reddit/posted",
        {
          params: { limit, offset },
        },
      );
      return response.data;
    },
    staleTime: 60 * 1000, // 1 minute
  });
};

// Get Reddit approval stats
export const useRedditApprovalStats = () => {
  const api = useApiClient();
  return useQuery<RedditApprovalStats>({
    queryKey: ["admin", "jobs", "reddit", "stats"],
    queryFn: async (): Promise<RedditApprovalStats> => {
      const response = await api.get<RedditApprovalStats>(
        "/admin/jobs/reddit/stats",
      );
      return response.data;
    },
    staleTime: 30 * 1000, // 30 seconds
  });
};

// Approve a job for Reddit posting
export const useApproveJob = () => {
  const api = useApiClient();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({
      jobId,
      approvedBy,
      subredditIds,
    }: {
      jobId: string;
      approvedBy: string;
      subredditIds?: string[];
    }) => {
      const response = await api.post(`/admin/jobs/reddit/approve/${jobId}`, {
        approved_by: approvedBy,
        subreddit_ids: subredditIds,
      } as RedditApprovalRequest);
      return response;
    },
    onSuccess: async () => {
      // Invalidate pending jobs to refresh the list
      await queryClient.invalidateQueries({
        queryKey: ["admin", "jobs", "reddit", "pending"],
      });
      await queryClient.invalidateQueries({
        queryKey: ["admin", "jobs", "reddit", "posted"],
      });
    },
  });
};

// Reject a job for Reddit posting
export const useRejectJob = () => {
  const api = useApiClient();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({
      jobId,
      rejectedBy,
      reason,
    }: {
      jobId: string;
      rejectedBy: string;
      reason?: string;
    }): Promise<RedditRejectionRequest> => {
      return await api
        .post<RedditRejectionRequest>(`/admin/jobs/reddit/reject/${jobId}`, {
          rejected_by: rejectedBy,
          reason,
        })
        .then((res) => res.data);
    },
    onSuccess: async () => {
      // Invalidate pending jobs to refresh the list
      await queryClient.invalidateQueries({
        queryKey: ["admin", "jobs", "reddit", "pending"],
      });
      await queryClient.invalidateQueries({
        queryKey: ["admin", "jobs", "reddit", "posted"],
      });
    },
  });
};

// Bulk approve jobs for Reddit posting
export const useBulkApproveJobs = () => {
  const api = useApiClient();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({
      jobIds,
      approvedBy,
      subredditIds,
    }: {
      jobIds: string[];
      approvedBy: string;
      subredditIds?: string[];
    }) => {
      return api
        .post<BulkApprovalRequest>("/admin/jobs/reddit/bulk-approve", {
          job_ids: jobIds,
          approved_by: approvedBy,
          subreddit_ids: subredditIds,
        })
        .then((res) => res.data);
    },
    onSuccess: async () => {
      // Invalidate pending jobs to refresh the list
      await queryClient.invalidateQueries({
        queryKey: ["admin", "jobs", "reddit", "pending"],
      });
      await queryClient.invalidateQueries({
        queryKey: ["admin", "jobs", "reddit", "posted"],
      });
    },
  });
};

// Bulk reject jobs for Reddit posting
export const useBulkRejectJobs = () => {
  const api = useApiClient();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({
      jobIds,
      rejectedBy,
      reason,
    }: {
      jobIds: string[];
      rejectedBy: string;
      reason?: string;
    }) => {
      return api
        .post<BulkRejectionRequest>("/admin/jobs/reddit/bulk-reject", {
          job_ids: jobIds,
          rejected_by: rejectedBy,
          reason,
        })
        .then((res) => res.data);
    },
    onSuccess: async () => {
      // Invalidate pending jobs to refresh the list
      await queryClient.invalidateQueries({
        queryKey: ["admin", "jobs", "reddit", "pending"],
      });
      await queryClient.invalidateQueries({
        queryKey: ["admin", "jobs", "reddit", "posted"],
      });
    },
  });
};

// Generate Reddit post content for multiple jobs
export const useGenerateRedditPosts = () => {
  const api = useApiClient();

  return useMutation({
    mutationFn: async (jobIds: string[]): Promise<BulkGenerationResponse> => {
      const response = await api.post<BulkGenerationResponse>(
        "/admin/jobs/reddit/generate-content/bulk",
        { job_ids: jobIds } as BulkGenerationRequest,
      );
      return response.data;
    },
  });
};

// Generate Reddit post content for a single job
export const useGenerateRedditPost = () => {
  const api = useApiClient();

  return useMutation({
    mutationFn: async (jobId: string): Promise<GeneratedRedditPost> => {
      const response = await api.post<GeneratedRedditPost>(
        `/admin/jobs/reddit/generate-content/${jobId}`,
      );
      return response.data;
    },
  });
};

// Preview what will be posted to Reddit for a job
export const usePreviewRedditPost = () => {
  const api = useApiClient();

  return useMutation({
    mutationFn: async (jobId: string): Promise<RedditPreview> => {
      const response = await api.get<RedditPreview>(
        `/admin/jobs/reddit/preview/${jobId}`,
      );
      return response.data;
    },
  });
};

