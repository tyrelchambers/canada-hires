import { JobPosting } from "@/types";
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
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faExternalLinkAlt, faCheck } from "@fortawesome/free-solid-svg-icons";
import { Pagination } from "@/components/ui/pagination";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { faReddit } from "@fortawesome/free-brands-svg-icons";

interface PaginationInfo {
  currentPage: number;
  pageSize: number;
  total: number;
  hasMore: boolean;
}

interface PostedJobsListProps {
  jobs: JobPosting[];
  pagination: PaginationInfo;
  onPageChange: (page: number) => void;
  onPageSizeChange: (pageSize: number) => void;
  isLoading?: boolean;
}

export function PostedJobsList({
  jobs,
  pagination,
  onPageChange,
  onPageSizeChange,
  isLoading = false,
}: PostedJobsListProps) {
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

  const formatDateTime = (dateString?: string) => {
    if (!dateString) return "Not specified";
    return new Date(dateString).toLocaleString();
  };

  if (isLoading) {
    return (
      <Card>
        <CardContent className="pt-6">
          <div className="text-center py-12">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>
            <p className="mt-2 text-gray-600">Loading posted jobs...</p>
          </div>
        </CardContent>
      </Card>
    );
  }

  if (jobs.length === 0) {
    return (
      <Card>
        <CardContent className="pt-6">
          <div className="text-center py-12">
            <FontAwesomeIcon
              icon={faReddit}
              className="mx-auto text-4xl text-orange-500"
            />
            <h3 className="mt-4 text-lg font-medium text-gray-900">
              No posted jobs yet
            </h3>
            <p className="mt-2 text-gray-600">
              Jobs will appear here after they've been approved and posted to
              Reddit.
            </p>
          </div>
        </CardContent>
      </Card>
    );
  }

  return (
    <div className="space-y-4">
      {/* Jobs Table */}
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <div>
              <CardTitle className="text-lg flex items-center space-x-2">
                <FontAwesomeIcon icon={faReddit} className="text-orange-500" />
                <span>Posted Jobs</span>
              </CardTitle>
              <p className="text-sm text-gray-600 mt-1">
                {pagination.total} total posted jobs • Page{" "}
                {pagination.currentPage + 1} of{" "}
                {Math.ceil(pagination.total / pagination.pageSize)} • Showing{" "}
                {jobs.length} jobs
              </p>
            </div>
          </div>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Job Title</TableHead>
                <TableHead>Employer</TableHead>
                <TableHead>Location</TableHead>
                <TableHead>Salary</TableHead>
                <TableHead>Posted Date</TableHead>
                <TableHead>Approved</TableHead>
                <TableHead>Posted To</TableHead>
                <TableHead>Status</TableHead>
                <TableHead className="text-right">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {jobs.map((job) => (
                <TableRow key={job.id}>
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
                  <TableCell className="text-sm">
                    <div className="space-y-1">
                      <div className="flex items-center space-x-1">
                        <FontAwesomeIcon
                          icon={faCheck}
                          className="text-green-600 text-xs"
                        />
                        <span className="text-xs">
                          {formatDateTime(job.reddit_approved_at)}
                        </span>
                      </div>
                      {job.reddit_approved_by && (
                        <div className="text-xs text-gray-500">
                          by {job.reddit_approved_by}
                        </div>
                      )}
                    </div>
                  </TableCell>
                  <TableCell>
                    <div className="space-y-1">
                      {job.subreddit_posts && job.subreddit_posts.length > 0 ? (
                        job.subreddit_posts.map((post) => (
                          <div key={post.id} className="flex items-center space-x-2">
                            {post.reddit_post_url ? (
                              <a
                                href={post.reddit_post_url}
                                target="_blank"
                                rel="noopener noreferrer"
                                className="flex items-center space-x-1 text-orange-600 hover:text-orange-800"
                              >
                                <FontAwesomeIcon icon={faReddit} className="text-xs" />
                                <span className="text-xs font-mono">
                                  r/{post.subreddit_name || 'unknown'}
                                </span>
                                <FontAwesomeIcon icon={faExternalLinkAlt} className="text-xs" />
                              </a>
                            ) : (
                              <div className="flex items-center space-x-1 text-gray-600">
                                <FontAwesomeIcon icon={faReddit} className="text-xs" />
                                <span className="text-xs font-mono">
                                  r/{post.subreddit_name || 'unknown'}
                                </span>
                              </div>
                            )}
                            <span className="text-xs text-gray-500">
                              {formatDateTime(post.posted_at)}
                            </span>
                          </div>
                        ))
                      ) : (
                        <span className="text-xs text-gray-500">Not posted yet</span>
                      )}
                    </div>
                  </TableCell>
                  <TableCell>
                    <div className="flex flex-col space-y-1">
                      <Badge
                        variant="default"
                        className="bg-green-100 text-green-800 w-fit"
                      >
                        {job.reddit_approval_status}
                      </Badge>
                      {job.reddit_posted && (
                        <Badge
                          variant="secondary"
                          className="bg-orange-100 text-orange-800 w-fit"
                        >
                          <FontAwesomeIcon icon={faReddit} className="mr-1" />
                          Posted
                        </Badge>
                      )}
                    </div>
                  </TableCell>
                  <TableCell className="text-right">
                    <div className="flex items-center justify-end space-x-2">
                      <a
                        href={job.url}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="inline-flex"
                      >
                        <Button
                          size="sm"
                          variant="outline"
                          className="text-blue-600 border-blue-300 hover:bg-blue-50"
                        >
                          <FontAwesomeIcon icon={faExternalLinkAlt} />
                        </Button>
                      </a>
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
    </div>
  );
}
