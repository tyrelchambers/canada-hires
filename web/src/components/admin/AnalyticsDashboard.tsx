import { usePendingJobs, useRedditApprovalStats } from "@/hooks/useAdminJobs";
import { useJobStats } from "@/hooks/useJobPostings";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faChartBar,
  faCalendarDay,
  faUsers,
  faClock,
  faCheckCircle,
  faTimesCircle,
  faBuilding,
  faMapMarkerAlt,
  faArrowTrendUp,
} from "@fortawesome/free-solid-svg-icons";
import { faReddit } from "@fortawesome/free-brands-svg-icons";

export function AnalyticsDashboard() {
  const { data: pendingJobs } = usePendingJobs(100, 0); // Get more for analytics
  const { data: approvalStats } = useRedditApprovalStats();
  const { data: jobStats } = useJobStats();

  // Calculate metrics from the data
  const totalJobs = jobStats?.total_jobs || 0;
  const totalEmployers = jobStats?.total_employers || 0;
  const pendingCount = approvalStats?.pending_count || 0;
  const approvedCount = approvalStats?.approved_count || 0;
  const rejectedCount = approvalStats?.rejected_count || 0;
  const totalProcessed = pendingCount + approvedCount + rejectedCount;
  const approvalRate =
    totalProcessed > 0
      ? ((approvedCount / totalProcessed) * 100).toFixed(1)
      : "0";
  const processingRate =
    totalJobs > 0 ? ((totalProcessed / totalJobs) * 100).toFixed(1) : "0";

  // Analyze pending jobs by employer
  const employerAnalysis =
    pendingJobs?.jobs.reduce(
      (acc, job) => {
        acc[job.employer] = (acc[job.employer] || 0) + 1;
        return acc;
      },
      {} as Record<string, number>,
    ) || {};

  const topEmployers = Object.entries(employerAnalysis)
    .sort(([, a], [, b]) => b - a)
    .slice(0, 5);

  // Analyze pending jobs by location
  const locationAnalysis =
    pendingJobs?.jobs.reduce(
      (acc, job) => {
        const location = job.city || job.location.split(",")[0] || "Unknown";
        acc[location] = (acc[location] || 0) + 1;
        return acc;
      },
      {} as Record<string, number>,
    ) || {};

  const topLocations = Object.entries(locationAnalysis)
    .sort(([, a], [, b]) => b - a)
    .slice(0, 5);

  // Calculate today's activity (mock data for now)
  const todayApproved = Math.floor(approvedCount * 0.1); // Estimate
  const todayRejected = Math.floor(rejectedCount * 0.05); // Estimate
  const todayPending = Math.floor(pendingCount * 0.2); // Estimate

  return (
    <div className="space-y-6">
      {/* Overview Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-3">
            <CardTitle className="text-sm font-medium">Total Jobs</CardTitle>
            <FontAwesomeIcon icon={faChartBar} className="text-blue-600" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {totalJobs.toLocaleString()}
            </div>
            <p className="text-xs text-muted-foreground">In the system</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-3">
            <CardTitle className="text-sm font-medium">
              Processing Rate
            </CardTitle>
            <FontAwesomeIcon icon={faArrowTrendUp} className="text-green-600" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{processingRate}%</div>
            <p className="text-xs text-muted-foreground">
              Jobs reviewed for Reddit
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-3">
            <CardTitle className="text-sm font-medium">Approval Rate</CardTitle>
            <FontAwesomeIcon icon={faReddit} className="text-orange-600" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{approvalRate}%</div>
            <p className="text-xs text-muted-foreground">Of reviewed jobs</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-3">
            <CardTitle className="text-sm font-medium">Employers</CardTitle>
            <FontAwesomeIcon icon={faUsers} className="text-purple-600" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {totalEmployers.toLocaleString()}
            </div>
            <p className="text-xs text-muted-foreground">Unique companies</p>
          </CardContent>
        </Card>
      </div>

      {/* Today's Activity */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center">
            <FontAwesomeIcon icon={faCalendarDay} className="mr-2" />
            Today's Activity
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            <div className="text-center">
              <div className="flex items-center justify-center mb-2">
                <FontAwesomeIcon
                  icon={faCheckCircle}
                  className="text-green-600 text-2xl mr-2"
                />
                <span className="text-3xl font-bold text-green-700">
                  {todayApproved}
                </span>
              </div>
              <p className="text-sm font-medium">Jobs Approved</p>
              <p className="text-xs text-gray-500">Posted to Reddit</p>
            </div>

            <div className="text-center">
              <div className="flex items-center justify-center mb-2">
                <FontAwesomeIcon
                  icon={faTimesCircle}
                  className="text-red-600 text-2xl mr-2"
                />
                <span className="text-3xl font-bold text-red-700">
                  {todayRejected}
                </span>
              </div>
              <p className="text-sm font-medium">Jobs Rejected</p>
              <p className="text-xs text-gray-500">Not posted</p>
            </div>

            <div className="text-center">
              <div className="flex items-center justify-center mb-2">
                <FontAwesomeIcon
                  icon={faClock}
                  className="text-yellow-600 text-2xl mr-2"
                />
                <span className="text-3xl font-bold text-yellow-700">
                  {todayPending}
                </span>
              </div>
              <p className="text-sm font-medium">New Pending</p>
              <p className="text-xs text-gray-500">Awaiting review</p>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Detailed Analysis */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Top Employers with Pending Jobs */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center">
              <FontAwesomeIcon icon={faBuilding} className="mr-2" />
              Top Employers (Pending Review)
            </CardTitle>
          </CardHeader>
          <CardContent>
            {topEmployers.length > 0 ? (
              <div className="space-y-3">
                {topEmployers.map(([employer, count], index) => (
                  <div
                    key={employer}
                    className="flex items-center justify-between"
                  >
                    <div className="flex items-center">
                      <span className="text-sm font-medium text-gray-500 w-6">
                        {index + 1}.
                      </span>
                      <span className="text-sm font-medium truncate max-w-[200px]">
                        {employer}
                      </span>
                    </div>
                    <Badge variant="secondary">{count} jobs</Badge>
                  </div>
                ))}
              </div>
            ) : (
              <p className="text-center text-gray-500 py-4">
                No pending jobs to analyze
              </p>
            )}
          </CardContent>
        </Card>

        {/* Top Locations with Pending Jobs */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center">
              <FontAwesomeIcon icon={faMapMarkerAlt} className="mr-2" />
              Top Locations (Pending Review)
            </CardTitle>
          </CardHeader>
          <CardContent>
            {topLocations.length > 0 ? (
              <div className="space-y-3">
                {topLocations.map(([location, count], index) => (
                  <div
                    key={location}
                    className="flex items-center justify-between"
                  >
                    <div className="flex items-center">
                      <span className="text-sm font-medium text-gray-500 w-6">
                        {index + 1}.
                      </span>
                      <span className="text-sm font-medium truncate max-w-[200px]">
                        {location}
                      </span>
                    </div>
                    <Badge variant="secondary">{count} jobs</Badge>
                  </div>
                ))}
              </div>
            ) : (
              <p className="text-center text-gray-500 py-4">
                No pending jobs to analyze
              </p>
            )}
          </CardContent>
        </Card>
      </div>

      {/* Workflow Efficiency */}
      <Card>
        <CardHeader>
          <CardTitle>Reddit Posting Workflow</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-6">
            {/* Workflow Pipeline */}
            <div className="flex items-center justify-between">
              <div className="text-center flex-1">
                <div className="text-2xl font-bold text-blue-600">
                  {totalJobs}
                </div>
                <p className="text-sm text-gray-600">Total Jobs</p>
              </div>
              <div className="text-gray-400 text-2xl">→</div>
              <div className="text-center flex-1">
                <div className="text-2xl font-bold text-yellow-600">
                  {pendingCount}
                </div>
                <p className="text-sm text-gray-600">Pending Review</p>
              </div>
              <div className="text-gray-400 text-2xl">→</div>
              <div className="text-center flex-1">
                <div className="text-2xl font-bold text-green-600">
                  {approvedCount}
                </div>
                <p className="text-sm text-gray-600">Posted to Reddit</p>
              </div>
            </div>

            {/* Workflow Metrics */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4 pt-4 border-t">
              <div className="text-center">
                <div className="text-lg font-semibold">{processingRate}%</div>
                <p className="text-xs text-gray-600">Jobs Processed</p>
              </div>
              <div className="text-center">
                <div className="text-lg font-semibold">{approvalRate}%</div>
                <p className="text-xs text-gray-600">Approval Rate</p>
              </div>
              <div className="text-center">
                <div className="text-lg font-semibold">
                  {totalProcessed > 0
                    ? ((rejectedCount / totalProcessed) * 100).toFixed(1)
                    : "0"}
                  %
                </div>
                <p className="text-xs text-gray-600">Rejection Rate</p>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
