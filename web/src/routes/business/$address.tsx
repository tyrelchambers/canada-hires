import { createFileRoute, Link } from "@tanstack/react-router";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Progress } from "@/components/ui/progress";
import { Separator } from "@/components/ui/separator";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faEye,
  faMapMarkerAlt,
  faStar,
  faUsers,
  faShield,
  faExclamationTriangle,
  faChartLine,
  faFlag,
  faShare,
} from "@fortawesome/free-solid-svg-icons";
import { useAddressReports } from "@/hooks/useReports";
import { AuthNav } from "@/components/AuthNav";

function BusinessDetailPage() {
  const { address } = Route.useParams();
  const decodedAddress = decodeURIComponent(address);

  const {
    data: reportsData,
    isLoading: loading,
    error,
  } = useAddressReports(decodedAddress);

  const reports = reportsData?.reports || [];

  const getConfidenceColor = (confidence: number) => {
    if (confidence >= 80) return "text-green-600";
    if (confidence >= 60) return "text-yellow-600";
    return "text-red-600";
  };

  const getTFWRating = (confidenceLevel: number) => {
    if (confidenceLevel >= 8)
      return {
        rating: "red",
        color: "text-white",
        label: "High TFW Usage",
        percentage: 45,
      };
    if (confidenceLevel >= 5)
      return {
        rating: "yellow",
        color: "text-yellow-600",
        label: "Moderate TFW Usage",
        percentage: 25,
      };
    return {
      rating: "green",
      color: "text-green-600",
      label: "Low TFW Usage",
      percentage: 10,
    };
  };

  const getStatusLabel = (status: string) => {
    switch (status) {
      case "approved":
        return "Verified";
      case "pending":
        return "Under Review";
      case "rejected":
        return "Unverified";
      case "flagged":
        return "Flagged";
      default:
        return "Unknown";
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-slate-50">
        <header className="border-b bg-white">
          <div className="container mx-auto px-4 py-4">
            <div className="flex items-center justify-between">
              <Link to="/" className="flex items-center space-x-2">
                <div className="w-8 h-8 bg-red-600 rounded-lg flex items-center justify-center">
                  <FontAwesomeIcon
                    icon={faEye}
                    className="w-5 h-5 text-white"
                  />
                </div>
                <span className="text-xl font-bold text-slate-900">
                  JobWatch Canada
                </span>
              </Link>
              <AuthNav />
            </div>
          </div>
        </header>
        <div className="container mx-auto px-4 py-8">
          <div className="text-center py-12">
            <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-red-600"></div>
            <p className="mt-4 text-gray-600 text-lg">
              Loading business details...
            </p>
          </div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen bg-slate-50">
        <header className="border-b bg-white">
          <div className="container mx-auto px-4 py-4">
            <div className="flex items-center justify-between">
              <Link to="/" className="flex items-center space-x-2">
                <div className="w-8 h-8 bg-red-600 rounded-lg flex items-center justify-center">
                  <FontAwesomeIcon
                    icon={faEye}
                    className="w-5 h-5 text-white"
                  />
                </div>
                <span className="text-xl font-bold text-slate-900">
                  JobWatch Canada
                </span>
              </Link>
              <AuthNav />
            </div>
          </div>
        </header>
        <div className="container mx-auto px-4 py-8">
          <Card className="border-red-200 bg-red-50">
            <CardContent className="p-6">
              <p className="text-red-800">
                Error:{" "}
                {error instanceof Error ? error.message : "An error occurred"}
              </p>
            </CardContent>
          </Card>
        </div>
      </div>
    );
  }

  const businessName =
    reports.length > 0 ? reports[0].business_name : "Unknown Business";
  const averageConfidence =
    reports.length > 0
      ? reports.reduce(
          (sum, report) => sum + (report.confidence_level || 0),
          0,
        ) / reports.length
      : 0;

  const businessRating = getTFWRating(averageConfidence);
  const confidencePercentage = Math.round(averageConfidence * 10);
  const verifiedReports = reports.filter((r) => r.status === "approved").length;

  // Calculate report distribution
  const highTFWReports = reports.filter(
    (r) => (r.confidence_level || 0) >= 8,
  ).length;
  const moderateTFWReports = reports.filter(
    (r) => (r.confidence_level || 0) >= 5 && (r.confidence_level || 0) < 8,
  ).length;
  const lowTFWReports = reports.filter(
    (r) => (r.confidence_level || 0) < 5,
  ).length;

  return (
    <div className="min-h-screen bg-slate-50">
      {/* Header */}
      <AuthNav />{" "}
      <div className="max-w-screen-xl mx-auto px-4 py-8">
        {/* Breadcrumb */}
        <div className="mb-6">
          <nav className="text-sm text-slate-600">
            <Link to="/" className="hover:text-slate-900">
              Home
            </Link>
            <span className="mx-2">/</span>
            <Link to="/directory" className="hover:text-slate-900">
              Directory
            </Link>
            <span className="mx-2">/</span>
            <span className="text-slate-900">{businessName}</span>
          </nav>
        </div>

        <div className="grid lg:grid-cols-3 gap-8">
          {/* Main Content */}
          <div className="lg:col-span-2 space-y-6">
            {/* Business Header */}
            <Card>
              <CardHeader>
                <div className="flex items-start justify-between">
                  <div className="flex-1">
                    <CardTitle className="text-2xl mb-2">
                      {businessName}
                    </CardTitle>
                    <div className="flex items-center space-x-4 text-slate-600 mb-4">
                      <div className="flex items-center space-x-1">
                        <FontAwesomeIcon
                          icon={faMapMarkerAlt}
                          className="w-4 h-4"
                        />
                        <span>{decodedAddress}</span>
                      </div>
                    </div>
                    <div className="flex items-center space-x-4">
                      <Badge variant="secondary">Business</Badge>
                      <Badge>{businessRating.label}</Badge>
                    </div>
                  </div>
                </div>
              </CardHeader>
              <CardContent>
                <p className="text-slate-600">
                  Community-reported business in our directory with TFW usage
                  data.
                </p>
              </CardContent>
            </Card>

            {/* TFW Usage Stats */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center space-x-2">
                  <FontAwesomeIcon icon={faChartLine} className="w-5 h-5" />
                  <span>TFW Usage Analysis</span>
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-6">
                <div className="grid md:grid-cols-3 gap-6">
                  <div className="text-center">
                    <div className="text-3xl font-bold text-slate-900 mb-1">
                      {businessRating.percentage}%
                    </div>
                    <div className="text-sm text-slate-600">
                      Estimated TFW Usage
                    </div>
                  </div>
                  <div className="text-center">
                    <div
                      className={`text-3xl font-bold mb-1 ${getConfidenceColor(confidencePercentage)}`}
                    >
                      {confidencePercentage}%
                    </div>
                    <div className="text-sm text-slate-600">
                      Confidence Score
                    </div>
                  </div>
                  <div className="text-center">
                    <div className="text-3xl font-bold text-slate-900 mb-1">
                      {reports.length}
                    </div>
                    <div className="text-sm text-slate-600">
                      Community Reports
                    </div>
                  </div>
                </div>

                <Separator />

                <div className="space-y-4">
                  <h4 className="font-semibold">Report Distribution</h4>
                  <div className="space-y-3">
                    <div className="flex items-center justify-between">
                      <div className="flex items-center space-x-2">
                        <div className="w-3 h-3 bg-red-500 rounded-full"></div>
                        <span className="text-sm">
                          High TFW (8-10 confidence)
                        </span>
                      </div>
                      <div className="flex items-center space-x-2">
                        <Progress
                          value={
                            reports.length > 0
                              ? (highTFWReports / reports.length) * 100
                              : 0
                          }
                          className="w-24"
                        />
                        <span className="text-sm font-medium">
                          {highTFWReports} reports
                        </span>
                      </div>
                    </div>
                    <div className="flex items-center justify-between">
                      <div className="flex items-center space-x-2">
                        <div className="w-3 h-3 bg-yellow-500 rounded-full"></div>
                        <span className="text-sm">
                          Moderate TFW (5-7 confidence)
                        </span>
                      </div>
                      <div className="flex items-center space-x-2">
                        <Progress
                          value={
                            reports.length > 0
                              ? (moderateTFWReports / reports.length) * 100
                              : 0
                          }
                          className="w-24"
                        />
                        <span className="text-sm font-medium">
                          {moderateTFWReports} reports
                        </span>
                      </div>
                    </div>
                    <div className="flex items-center justify-between">
                      <div className="flex items-center space-x-2">
                        <div className="w-3 h-3 bg-green-500 rounded-full"></div>
                        <span className="text-sm">
                          Low TFW (0-4 confidence)
                        </span>
                      </div>
                      <div className="flex items-center space-x-2">
                        <Progress
                          value={
                            reports.length > 0
                              ? (lowTFWReports / reports.length) * 100
                              : 0
                          }
                          className="w-24"
                        />
                        <span className="text-sm font-medium">
                          {lowTFWReports} reports
                        </span>
                      </div>
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Community Reports */}
            <div>
              <header>
                <h3 className="flex font-bold items-center space-x-2">
                  <FontAwesomeIcon icon={faUsers} className="w-5 h-5" />
                  <span>Community Reports</span>
                </h3>
                <p className="text-muted-foreground">
                  Reports from verified and unverified community members
                </p>
              </header>
              <div className="space-y-4 mt-6">
                {reports.length === 0 ? (
                  <div className="text-center py-8">
                    <p className="text-slate-500">
                      No reports found for this business address.
                    </p>
                  </div>
                ) : (
                  reports.map((report) => {
                    const reportRating = getTFWRating(
                      report.confidence_level || 0,
                    );
                    const isVerified = report.status === "approved";

                    return (
                      <div
                        key={report.id}
                        className="border rounded-lg p-4 space-y-3 bg-white"
                      >
                        <div className="flex items-start justify-between">
                          <div className="flex items-center space-x-3">
                            <Avatar className="w-8 h-8">
                              <AvatarFallback className="text-xs">
                                {report.report_source
                                  ?.substring(0, 2)
                                  .toUpperCase() || "AN"}
                              </AvatarFallback>
                            </Avatar>
                            <div>
                              <div className="flex items-center space-x-2">
                                <span className="text-sm font-medium">
                                  Community Member
                                </span>
                                {isVerified && (
                                  <FontAwesomeIcon
                                    icon={faShield}
                                    className="w-4 h-4 text-green-600"
                                  />
                                )}
                              </div>
                              <div className="text-xs text-slate-500">
                                {report.report_source} •{" "}
                                {getStatusLabel(report.status)}
                              </div>
                            </div>
                          </div>
                          <div className="text-xs text-slate-500">
                            {new Date(report.created_at).toLocaleDateString()}
                          </div>
                        </div>

                        <div className="flex items-center space-x-2">
                          <span className="text-sm">Reported TFW usage:</span>
                          <Badge>{reportRating.label}</Badge>
                        </div>

                        {report.additional_notes && (
                          <p className="text-sm text-slate-600 bg-slate-50 p-3 rounded">
                            "{report.additional_notes}"
                          </p>
                        )}
                      </div>
                    );
                  })
                )}
              </div>
            </div>
          </div>

          {/* Sidebar */}
          <div className="space-y-6">
            {/* Quick Actions */}
            <Card>
              <CardHeader>
                <CardTitle>Take Action</CardTitle>
              </CardHeader>
              <CardContent className="space-y-3">
                <Button className="w-full" asChild>
                  <Link to="/reports/create">Submit Report</Link>
                </Button>
              </CardContent>
            </Card>

            {/* Data Quality */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center space-x-2">
                  <FontAwesomeIcon icon={faStar} className="w-5 h-5" />
                  <span>Data Quality</span>
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="flex items-center justify-between">
                  <span className="text-sm">Confidence Score</span>
                  <span
                    className={`font-bold ${getConfidenceColor(confidencePercentage)}`}
                  >
                    {confidencePercentage}%
                  </span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-sm">Total Reports</span>
                  <span className="font-bold">{reports.length}</span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-sm">Verified Reports</span>
                  <span className="font-bold">{verifiedReports}</span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-sm">Last Updated</span>
                  <span className="text-sm text-slate-600">
                    {reports.length > 0
                      ? new Date(
                          Math.max(
                            ...reports.map((r) =>
                              new Date(r.created_at).getTime(),
                            ),
                          ),
                        ).toLocaleDateString()
                      : "No data"}
                  </span>
                </div>
              </CardContent>
            </Card>

            {/* Disclaimer */}
            <Card className="border-yellow-200 bg-yellow-50">
              <CardHeader>
                <CardTitle className="flex items-center space-x-2 text-yellow-800">
                  <FontAwesomeIcon
                    icon={faExclamationTriangle}
                    className="w-5 h-5"
                  />
                  <span>Disclaimer</span>
                </CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-sm text-yellow-700">
                  This information is based on community reports and may not
                  reflect current hiring practices. Data should be used for
                  informational purposes only.
                </p>
              </CardContent>
            </Card>
          </div>
        </div>
      </div>
    </div>
  );
}

export const Route = createFileRoute("/business/$address")({
  component: BusinessDetailPage,
});
