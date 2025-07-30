import { useState } from "react";
import { JobPosting, User } from "@/types";
import {
  useApproveJob,
  useRejectJob,
  useBulkApproveJobs,
  useBulkRejectJobs,
} from "@/hooks/useAdminJobs";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Checkbox } from "@/components/ui/checkbox";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faExternalLinkAlt,
  faCheck,
  faTimes,
} from "@fortawesome/free-solid-svg-icons";
import { RejectJobDialog } from "./RejectJobDialog";
import { ApprovalConfirmationModal } from "./ApprovalConfirmationModal";
import { Pagination } from "@/components/ui/pagination";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

interface PaginationInfo {
  currentPage: number;
  pageSize: number;
  total: number;
  hasMore: boolean;
}

interface PendingJobsListProps {
  jobs: JobPosting[];
  selectedJobs: string[];
  onJobSelect: (jobId: string, selected: boolean) => void;
  onSelectAll: () => void;
  onClearSelection: () => void;
  user: User;
  pagination: PaginationInfo;
  onPageChange: (page: number) => void;
  onPageSizeChange: (pageSize: number) => void;
}

export function PendingJobsList({
  jobs,
  selectedJobs,
  onJobSelect,
  onSelectAll,
  onClearSelection,
  user,
  pagination,
  onPageChange,
  onPageSizeChange,
}: PendingJobsListProps) {
  const [rejectingJobId, setRejectingJobId] = useState<string | null>(null);
  const [approvalConfirmation, setApprovalConfirmation] = useState<{
    type: 'single' | 'bulk';
    jobIds: string[];
  } | null>(null);

  const approveJobMutation = useApproveJob();
  const rejectJobMutation = useRejectJob();
  const bulkApproveMutation = useBulkApproveJobs();
  const bulkRejectMutation = useBulkRejectJobs();

  const handleApprove = (jobId: string) => {
    setApprovalConfirmation({
      type: 'single',
      jobIds: [jobId],
    });
  };

  const handleConfirmApproval = async () => {
    if (!approvalConfirmation) return;

    try {
      if (approvalConfirmation.type === 'single') {
        await approveJobMutation.mutateAsync({
          jobId: approvalConfirmation.jobIds[0],
          approvedBy: user.email,
        });
      } else {
        await bulkApproveMutation.mutateAsync({
          jobIds: approvalConfirmation.jobIds,
          approvedBy: user.email,
        });
        onClearSelection();
      }
      setApprovalConfirmation(null);
    } catch (error) {
      console.error("Failed to approve job(s):", error);
    }
  };

  const handleReject = async (jobId: string, reason?: string) => {
    try {
      await rejectJobMutation.mutateAsync({
        jobId,
        rejectedBy: user.email,
        reason,
      });
      setRejectingJobId(null);
    } catch (error) {
      console.error("Failed to reject job:", error);
    }
  };

  const handleBulkApprove = () => {
    if (selectedJobs.length === 0) return;

    setApprovalConfirmation({
      type: 'bulk',
      jobIds: selectedJobs,
    });
  };

  const handleBulkReject = async () => {
    if (selectedJobs.length === 0) return;

    try {
      await bulkRejectMutation.mutateAsync({
        jobIds: selectedJobs,
        rejectedBy: user.email,
        reason: "Bulk rejection for clean slate",
      });
      onClearSelection();
    } catch (error) {
      console.error("Failed to bulk reject jobs:", error);
    }
  };

  const formatSalary = (job: JobPosting) => {
    // Prefer structured data over raw string
    if (job.salary_min && job.salary_max) {
      const type = job.salary_type || "hourly";
      if (job.salary_min === job.salary_max) {
        return `$${job.salary_min} ${type}`;
      }
      return `$${job.salary_min} - $${job.salary_max} ${type}`;
    }
    // Fallback to raw salary string if structured data not available
    if (job.salary_raw) return job.salary_raw;
    return "Not specified";
  };

  const formatDate = (dateString?: string) => {
    if (!dateString) return "Not specified";
    return new Date(dateString).toLocaleDateString();
  };

  if (jobs.length === 0) {
    return (
      <Card>
        <CardContent className="pt-6">
          <div className="text-center py-12">
            <FontAwesomeIcon
              icon={faCheck}
              className="mx-auto text-4xl text-green-500"
            />
            <h3 className="mt-4 text-lg font-medium text-gray-900">
              All caught up!
            </h3>
            <p className="mt-2 text-gray-600">
              No jobs are currently pending Reddit approval.
            </p>
          </div>
        </CardContent>
      </Card>
    );
  }

  return (
    <div className="space-y-4">
      {/* Bulk Actions */}
      {selectedJobs.length > 0 && (
        <Card>
          <CardContent>
            <div className="flex items-center justify-between">
              <div className="flex items-center space-x-2">
                <span className="text-sm font-medium text-blue-900">
                  {selectedJobs.length} job
                  {selectedJobs.length !== 1 ? "s" : ""} selected
                </span>
              </div>
              <div className="flex items-center space-x-2">
                <Button variant="outline" size="sm" onClick={onClearSelection}>
                  Clear Selection
                </Button>
                <Button
                  size="sm"
                  variant="outline"
                  onClick={handleBulkReject}
                  disabled={bulkRejectMutation.isPending}
                  className="text-red-600 border-red-300 hover:bg-red-50"
                >
                  {bulkRejectMutation.isPending
                    ? "Rejecting..."
                    : `Reject ${selectedJobs.length} Jobs`}
                </Button>
                <Button
                  size="sm"
                  onClick={handleBulkApprove}
                  disabled={bulkApproveMutation.isPending || approveJobMutation.isPending}
                  className="bg-green-600 hover:bg-green-700"
                >
                  {bulkApproveMutation.isPending
                    ? "Approving..."
                    : `Approve ${selectedJobs.length} Jobs`}
                </Button>
              </div>
            </div>
          </CardContent>
        </Card>
      )}

      {/* Jobs Table */}
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <div>
              <CardTitle className="text-lg">Pending Jobs</CardTitle>
              <p className="text-sm text-gray-600 mt-1">
                {pagination.total} total jobs • Page{" "}
                {pagination.currentPage + 1} of{" "}
                {Math.ceil(pagination.total / pagination.pageSize)} • Showing{" "}
                {jobs.length} jobs
              </p>
            </div>
            <div className="flex items-center space-x-2">
              <Button
                variant="outline"
                size="sm"
                onClick={
                  selectedJobs.length === jobs.length
                    ? onClearSelection
                    : onSelectAll
                }
              >
                {selectedJobs.length === jobs.length
                  ? "Deselect All"
                  : "Select All"}
              </Button>
            </div>
          </div>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead className="w-12">
                  <span className="sr-only">Select</span>
                </TableHead>
                <TableHead>Job Title</TableHead>
                <TableHead>Employer</TableHead>
                <TableHead>Location</TableHead>
                <TableHead>Salary</TableHead>
                <TableHead>Posted</TableHead>
                <TableHead>Status</TableHead>
                <TableHead className="text-right">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {jobs.map((job) => (
                <TableRow key={job.id}>
                  <TableCell>
                    <Checkbox
                      checked={selectedJobs.includes(job.id)}
                      onCheckedChange={(checked) =>
                        onJobSelect(job.id, checked as boolean)
                      }
                    />
                  </TableCell>
                  <TableCell className="font-medium">
                    <div className="flex items-center space-x-2">
                      <span className="max-w-xs truncate">{job.title}</span>
                      <a
                        href={job.url}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="text-blue-600 hover:text-blue-800"
                      >
                        <FontAwesomeIcon icon={faExternalLinkAlt} />
                      </a>
                    </div>
                  </TableCell>
                  <TableCell>{job.employer}</TableCell>
                  <TableCell>{job.location}</TableCell>
                  <TableCell className="text-sm">{formatSalary(job)}</TableCell>
                  <TableCell className="text-sm">
                    {formatDate(job.posting_date)}
                  </TableCell>
                  <TableCell>
                    <Badge
                      variant="secondary"
                      className="bg-yellow-100 text-yellow-800"
                    >
                      {job.reddit_approval_status}
                    </Badge>
                  </TableCell>
                  <TableCell className="text-right">
                    <div className="flex items-center justify-end space-x-2">
                      <Button
                        size="sm"
                        variant="outline"
                        onClick={() => setRejectingJobId(job.id)}
                        disabled={rejectJobMutation.isPending}
                        className="text-red-600 border-red-300 hover:bg-red-50"
                      >
                        <FontAwesomeIcon icon={faTimes} />
                      </Button>
                      <Button
                        size="sm"
                        onClick={() => handleApprove(job.id)}
                        disabled={approveJobMutation.isPending || bulkApproveMutation.isPending}
                        className="bg-green-600 hover:bg-green-700"
                      >
                        <FontAwesomeIcon icon={faCheck} />
                      </Button>
                    </div>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </CardContent>
      </Card>

      {/* Pagination Controls */}
      {pagination.total > 0 && (
        <div className="flex items-center justify-between flex-wrap gap-4">
          <div className="flex items-center space-x-2">
            <span className="text-sm text-gray-600">Jobs per page:</span>
            <Select
              value={pagination.pageSize.toString()}
              onValueChange={(value) => onPageSizeChange(Number(value))}
            >
              <SelectTrigger className="w-20">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="25">25</SelectItem>
                <SelectItem value="50">50</SelectItem>
                <SelectItem value="100">100</SelectItem>
                <SelectItem value="200">200</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div className="flex-1">
            <Pagination
              currentPage={pagination.currentPage + 1} // Convert from 0-based to 1-based
              totalPages={Math.ceil(pagination.total / pagination.pageSize)}
              totalItems={pagination.total}
              itemsPerPage={pagination.pageSize}
              onPageChange={(page) => onPageChange(page - 1)} // Convert back to 0-based
            />
          </div>
        </div>
      )}

      {/* Reject Job Dialog */}
      {rejectingJobId && (
        <RejectJobDialog
          job={jobs.find((j) => j.id === rejectingJobId)!}
          onReject={handleReject}
          onCancel={() => setRejectingJobId(null)}
          isLoading={rejectJobMutation.isPending}
        />
      )}

      {/* Approval Confirmation Modal */}
      {approvalConfirmation && (
        <ApprovalConfirmationModal
          isOpen={true}
          onClose={() => setApprovalConfirmation(null)}
          onConfirm={handleConfirmApproval}
          jobCount={approvalConfirmation.jobIds.length}
          isLoading={approveJobMutation.isPending || bulkApproveMutation.isPending}
        />
      )}
    </div>
  );
}
