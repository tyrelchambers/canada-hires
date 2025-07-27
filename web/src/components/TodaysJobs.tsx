import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faMapMarkerAlt,
  faBuilding,
  faFlag,
  faUsers,
  faExternalLinkAlt,
  faInfoCircle,
} from "@fortawesome/free-solid-svg-icons";
import { useJobPostings } from "@/hooks/useJobPostings";
import { Button } from "@/components/ui/button";
import { useState, useEffect } from "react";

export function TodaysJobs() {
  const [showingYesterday, setShowingYesterday] = useState(false);

  // First try to get today's jobs
  const { data: todayData, isLoading: todayLoading } = useJobPostings({
    limit: 10,
    days: 1,
    sort_by: "posting_date",
    sort_order: "desc",
  });

  // Fallback to yesterday's jobs if needed
  const { data: yesterdayData, isLoading: yesterdayLoading } = useJobPostings({
    limit: 10,
    days: 2, // Last 2 days to capture yesterday's jobs
    sort_by: "posting_date",
    sort_order: "desc",
  });

  // Determine which data to show
  const hasTodayJobs = todayData && todayData?.jobs?.length > 0;
  const jobData = hasTodayJobs ? todayData : yesterdayData;
  const isLoading = todayLoading || (!hasTodayJobs && yesterdayLoading);

  // Update state when we determine we're showing yesterday's jobs
  useEffect(() => {
    if (
      !todayLoading &&
      !hasTodayJobs &&
      yesterdayData &&
      yesterdayData?.jobs?.length > 0
    ) {
      setShowingYesterday(true);
    } else if (hasTodayJobs) {
      setShowingYesterday(false);
    }
  }, [todayLoading, hasTodayJobs, yesterdayData]);

  const formatSalary = (min?: number, max?: number, type?: string) => {
    if (!min && !max) return null;

    const formatAmount = (amount: number) => {
      return new Intl.NumberFormat("en-CA", {
        style: "currency",
        currency: "CAD",
        minimumFractionDigits: 0,
        maximumFractionDigits: 0,
      }).format(amount);
    };

    let salaryText = "";
    if (min && max && min !== max) {
      salaryText = `${formatAmount(min)} - ${formatAmount(max)}`;
    } else if (min) {
      salaryText = formatAmount(min);
    } else if (max) {
      salaryText = formatAmount(max);
    }

    if (type && salaryText) {
      salaryText += ` ${type}`;
    }

    return salaryText;
  };

  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Today's Job Postings</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex items-center justify-center py-8">
            <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-gray-900"></div>
            <span className="ml-2">Loading...</span>
          </div>
        </CardContent>
      </Card>
    );
  }

  if (!jobData?.jobs?.length) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Today's Job Postings</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-gray-600 text-center py-8">
            No job postings found for today.
          </p>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>
          {showingYesterday ? "Recent Job Postings" : "Today's Job Postings"}
        </CardTitle>
        {showingYesterday && (
          <div className="flex items-center gap-2 text-sm text-amber-700 bg-amber-50 p-3 rounded-lg border border-amber-200">
            <FontAwesomeIcon icon={faInfoCircle} className="w-4 h-4" />
            <span>
              No job postings found for today. Showing recent postings from
              yesterday.
            </span>
          </div>
        )}
      </CardHeader>
      <CardContent>
        <div className="bg-white rounded-lg border">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Job Title</TableHead>
                <TableHead>Employer</TableHead>
                <TableHead>Location</TableHead>
                <TableHead>Salary</TableHead>
                <TableHead>Status</TableHead>
                <TableHead></TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {jobData.jobs.map((job) => (
                <TableRow key={job.id}>
                  <TableCell className="font-medium">
                    <div className="max-w-xs truncate" title={job.title}>
                      {job.title}
                    </div>
                  </TableCell>
                  <TableCell>
                    <div className="flex items-center max-w-xs">
                      <FontAwesomeIcon
                        icon={faBuilding}
                        className="mr-2 text-gray-400 w-3 h-3"
                      />
                      <span className="truncate" title={job.employer}>
                        {job.employer}
                      </span>
                    </div>
                  </TableCell>
                  <TableCell>
                    <div className="flex items-center">
                      <FontAwesomeIcon
                        icon={faMapMarkerAlt}
                        className="mr-2 text-gray-400 w-3 h-3"
                      />
                      <span
                        className="text-sm truncate max-w-32"
                        title={job.location}
                      >
                        {job.location}
                      </span>
                    </div>
                  </TableCell>
                  <TableCell>
                    {formatSalary(
                      job.salary_min,
                      job.salary_max,
                      job.salary_type,
                    ) ? (
                      <div className="flex items-center">
                        <span className="text-sm">
                          {formatSalary(
                            job.salary_min,
                            job.salary_max,
                            job.salary_type,
                          )}
                        </span>
                      </div>
                    ) : (
                      <span className="text-gray-400">-</span>
                    )}
                  </TableCell>
                  <TableCell>
                    <div className="flex items-center gap-1">
                      <Badge
                        variant="outline"
                        className="bg-yellow-50 text-yellow-700 border-yellow-200"
                      >
                        <FontAwesomeIcon
                          icon={faFlag}
                          className="w-3 h-3 mr-1"
                        />
                        LMIA
                      </Badge>
                      {job.is_tfw && (
                        <Badge
                          variant="outline"
                          className="bg-blue-50 text-blue-700 border-blue-200"
                        >
                          <FontAwesomeIcon
                            icon={faUsers}
                            className="w-3 h-3 mr-1"
                          />
                          TFW
                        </Badge>
                      )}
                    </div>
                  </TableCell>
                  <TableCell>
                    <Button variant="outline" size="sm" asChild>
                      <a
                        href={job.url}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="flex items-center gap-1"
                      >
                        <FontAwesomeIcon
                          icon={faExternalLinkAlt}
                          className="w-3 h-3"
                        />
                        Apply
                      </a>
                    </Button>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </div>
      </CardContent>
    </Card>
  );
}
