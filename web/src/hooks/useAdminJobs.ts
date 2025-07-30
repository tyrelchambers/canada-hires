import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useApiClient } from "./useApiClient";
import {
  JobPostingsResponse,
  RedditApprovalRequest,
  RedditRejectionRequest,
  BulkApprovalRequest,
  BulkRejectionRequest,
  RedditApprovalStats,
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

// Approve a job for Reddit posting
export const useApproveJob = () => {
  const api = useApiClient();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({
      jobId,
      approvedBy,
    }: {
      jobId: string;
      approvedBy: string;
    }) => {
      const response = await api.post(`/admin/jobs/reddit/approve/${jobId}`, {
        approved_by: approvedBy,
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
      await queryClient.invalidateQueries({
        queryKey: ["admin", "jobs", "reddit", "stats"],
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
      await queryClient.invalidateQueries({
        queryKey: ["admin", "jobs", "reddit", "stats"],
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
    }: {
      jobIds: string[];
      approvedBy: string;
    }) => {
      return api
        .post<BulkApprovalRequest>("/admin/jobs/reddit/bulk-approve", {
          job_ids: jobIds,
          approved_by: approvedBy,
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
      await queryClient.invalidateQueries({
        queryKey: ["admin", "jobs", "reddit", "stats"],
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
      await queryClient.invalidateQueries({
        queryKey: ["admin", "jobs", "reddit", "stats"],
      });
    },
  });
};

// Get Reddit approval statistics
export const useRedditApprovalStats = () => {
  const api = useApiClient();
  return useQuery({
    queryKey: ["admin", "jobs", "reddit", "stats"],
    queryFn: async (): Promise<RedditApprovalStats> => {
      const response = await api.get<RedditApprovalStats>(
        "/admin/jobs/reddit/stats",
      );
      return response.data;
    },
    staleTime: 60 * 1000, // 1 minute
  });
};
