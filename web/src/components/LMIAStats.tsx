import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  useLMIAStats,
  useLMIAStatus,
  useTriggerLMIAUpdate,
} from "@/hooks/useLMIA";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faDatabase,
  faSync,
  faClock,
  faCheckCircle,
  faTimesCircle,
  faExclamationCircle,
  faDownload,
} from "@fortawesome/free-solid-svg-icons";

export function LMIAStats() {
  const { data: stats, isLoading: statsLoading } = useLMIAStats();
  const { data: status, isLoading: statusLoading } = useLMIAStatus();
  const { mutate: triggerUpdate, isPending: isUpdating } =
    useTriggerLMIAUpdate();

  const getStatusIcon = (statusStr?: string) => {
    switch (statusStr) {
      case "completed":
        return (
          <FontAwesomeIcon
            icon={faCheckCircle}
            className="w-4 h-4 text-green-500"
          />
        );
      case "running":
        return (
          <FontAwesomeIcon
            icon={faSync}
            className="w-4 h-4 text-blue-500 animate-spin"
          />
        );
      case "failed":
        return (
          <FontAwesomeIcon
            icon={faTimesCircle}
            className="w-4 h-4 text-red-500"
          />
        );
      default:
        return (
          <FontAwesomeIcon
            icon={faExclamationCircle}
            className="w-4 h-4 text-gray-500"
          />
        );
    }
  };

  const getStatusColor = (statusStr?: string) => {
    switch (statusStr) {
      case "completed":
        return "bg-green-100 text-green-800";
      case "running":
        return "bg-blue-100 text-blue-800";
      case "failed":
        return "bg-red-100 text-red-800";
      default:
        return "bg-gray-100 text-gray-800";
    }
  };

  const formatDate = (dateString?: string) => {
    if (!dateString) return "Never";
    return new Date(dateString).toLocaleString();
  };

  if (statsLoading || statusLoading) {
    return (
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        {[...Array(4)].map((_, i) => (
          <Card key={i}>
            <CardContent className="p-6">
              <div className="animate-pulse">
                <div className="h-4 bg-gray-200 rounded w-3/4 mb-2"></div>
                <div className="h-6 bg-gray-200 rounded w-1/2"></div>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Stats Cards */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Total Resources
            </CardTitle>
            <FontAwesomeIcon
              icon={faDatabase}
              className="h-4 w-4 text-muted-foreground"
            />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {stats?.total_resources || 0}
            </div>
            <p className="text-xs text-muted-foreground">
              LMIA data files available
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Processed</CardTitle>
            <FontAwesomeIcon
              icon={faCheckCircle}
              className="h-4 w-4 text-muted-foreground"
            />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {stats?.processed_resources || 0}
            </div>
            <p className="text-xs text-muted-foreground">
              Files successfully processed
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Records</CardTitle>
            <FontAwesomeIcon
              icon={faDownload}
              className="h-4 w-4 text-muted-foreground"
            />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {stats?.total_records?.toLocaleString() || 0}
            </div>
            <p className="text-xs text-muted-foreground">
              Employer records in database
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Last Update</CardTitle>
            <FontAwesomeIcon
              icon={faClock}
              className="h-4 w-4 text-muted-foreground"
            />
          </CardHeader>
          <CardContent>
            <div className="text-sm font-medium">
              {formatDate(stats?.last_update)}
            </div>
            <div className="flex items-center gap-1 mt-1">
              {getStatusIcon(stats?.last_update_status)}
              <Badge
                variant="secondary"
                className={`text-xs ${getStatusColor(stats?.last_update_status)}`}
              >
                {stats?.last_update_status || "unknown"}
              </Badge>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
