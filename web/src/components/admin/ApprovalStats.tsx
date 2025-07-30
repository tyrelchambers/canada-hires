import { RedditApprovalStats } from "@/types";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faClock,
  faCheckCircle,
  faTimesCircle,
  faChartLine,
} from "@fortawesome/free-solid-svg-icons";

interface ApprovalStatsProps {
  stats?: RedditApprovalStats;
  isLoading: boolean;
}

export function ApprovalStats({ stats, isLoading }: ApprovalStatsProps) {
  if (isLoading) {
    return (
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {Array.from({ length: 3 }).map((_, i) => (
          <Card key={i}>
            <CardContent className="pt-6">
              <div className="animate-pulse">
                <div className="h-4 bg-gray-200 rounded w-3/4 mb-2"></div>
                <div className="h-8 bg-gray-200 rounded w-1/2 mb-1"></div>
                <div className="h-3 bg-gray-200 rounded w-full"></div>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>
    );
  }

  if (!stats) {
    return (
      <Card>
        <CardContent className="pt-6">
          <div className="text-center text-gray-600">
            <p>Unable to load statistics</p>
          </div>
        </CardContent>
      </Card>
    );
  }

  const total =
    stats.pending_count + stats.approved_count + stats.rejected_count;
  const approvalRate =
    total > 0 ? ((stats.approved_count / total) * 100).toFixed(1) : "0";

  return (
    <div className="space-y-6">
      {/* Overview Stats */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-3">
            <CardTitle className="text-sm font-medium">Pending</CardTitle>
            <FontAwesomeIcon icon={faClock} className="text-yellow-600" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-yellow-700">
              {stats.pending_count}
            </div>
            <p className="text-xs text-muted-foreground">Awaiting review</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-3">
            <CardTitle className="text-sm font-medium">Approved</CardTitle>
            <FontAwesomeIcon icon={faCheckCircle} className="text-green-600" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-green-700">
              {stats.approved_count}
            </div>
            <p className="text-xs text-muted-foreground">Posted to Reddit</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-3">
            <CardTitle className="text-sm font-medium">Rejected</CardTitle>
            <FontAwesomeIcon icon={faTimesCircle} className="text-red-600" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-red-700">
              {stats.rejected_count}
            </div>
            <p className="text-xs text-muted-foreground">Not posted</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-3">
            <CardTitle className="text-sm font-medium">Approval Rate</CardTitle>
            <FontAwesomeIcon icon={faChartLine} className="text-blue-600" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-blue-700">
              {approvalRate}%
            </div>
            <p className="text-xs text-muted-foreground">Of reviewed jobs</p>
          </CardContent>
        </Card>
      </div>

      {/* Detailed Breakdown */}
      <Card>
        <CardHeader>
          <CardTitle>Review Summary</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            <div className="flex items-center justify-between">
              <span className="text-sm font-medium">Total Jobs Processed</span>
              <Badge variant="secondary">{total}</Badge>
            </div>

            <div className="space-y-2">
              <div className="flex items-center justify-between text-sm">
                <span className="flex items-center">
                  <div className="w-3 h-3 bg-green-500 rounded-full mr-2"></div>
                  Approved
                </span>
                <span>
                  {stats.approved_count} (
                  {total > 0
                    ? ((stats.approved_count / total) * 100).toFixed(1)
                    : 0}
                  %)
                </span>
              </div>

              <div className="flex items-center justify-between text-sm">
                <span className="flex items-center">
                  <div className="w-3 h-3 bg-red-500 rounded-full mr-2"></div>
                  Rejected
                </span>
                <span>
                  {stats.rejected_count} (
                  {total > 0
                    ? ((stats.rejected_count / total) * 100).toFixed(1)
                    : 0}
                  %)
                </span>
              </div>

              <div className="flex items-center justify-between text-sm">
                <span className="flex items-center">
                  <div className="w-3 h-3 bg-yellow-500 rounded-full mr-2"></div>
                  Pending
                </span>
                <span>
                  {stats.pending_count} (
                  {total > 0
                    ? ((stats.pending_count / total) * 100).toFixed(1)
                    : 0}
                  %)
                </span>
              </div>
            </div>

            {/* Progress Bar */}
            <div className="w-full bg-gray-200 rounded-full h-2 mt-4">
              <div className="flex h-2 rounded-full overflow-hidden">
                <div
                  className="bg-green-500"
                  style={{
                    width:
                      total > 0
                        ? `${(stats.approved_count / total) * 100}%`
                        : "0%",
                  }}
                ></div>
                <div
                  className="bg-red-500"
                  style={{
                    width:
                      total > 0
                        ? `${(stats.rejected_count / total) * 100}%`
                        : "0%",
                  }}
                ></div>
                <div
                  className="bg-yellow-500"
                  style={{
                    width:
                      total > 0
                        ? `${(stats.pending_count / total) * 100}%`
                        : "0%",
                  }}
                ></div>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Action Items */}
      {stats.pending_count > 0 && (
        <Card className="border-yellow-200 bg-yellow-50">
          <CardContent className="pt-6">
            <div className="flex items-center">
              <FontAwesomeIcon
                icon={faClock}
                className="text-yellow-600 mr-2"
              />
              <div>
                <p className="font-medium text-yellow-900">
                  {stats.pending_count} job
                  {stats.pending_count !== 1 ? "s" : ""} awaiting review
                </p>
                <p className="text-sm text-yellow-700">
                  Switch to the "Pending Jobs" tab to review and approve jobs
                  for Reddit posting.
                </p>
              </div>
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  );
}
