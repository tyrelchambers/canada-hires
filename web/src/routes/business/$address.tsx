import { createFileRoute, Link } from "@tanstack/react-router";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
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
} from "@fortawesome/free-solid-svg-icons";
import { useAddressReports } from "@/hooks/useReports";
import { AuthNav } from "@/components/AuthNav";
import { BoycottButton } from "@/components/BoycottButton";
import { useBoycottStats } from "@/hooks/useBoycotts";

function BusinessDetailPage() {
  const { address } = Route.useParams();
  const decodedAddress = decodeURIComponent(address);

  const {
    data: reportsData,
    isLoading: loading,
    error,
  } = useAddressReports(decodedAddress);

  const reports = reportsData?.reports || [];
  const businessName = reports.length > 0 ? reports[0].business_name : "Unknown Business";
  
  const { data: boycottStats } = useBoycottStats(
    businessName,
    decodedAddress
  );

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

  // Calculate TFW ratio distribution
  const tfwRatioDistribution = {
    few: reports.filter(r => r.tfw_ratio === 'few').length,
    many: reports.filter(r => r.tfw_ratio === 'many').length,
    most: reports.filter(r => r.tfw_ratio === 'most').length,
    all: reports.filter(r => r.tfw_ratio === 'all').length,
  };

  // Get most common ratio
  const mostCommonRatio = Object.entries(tfwRatioDistribution)
    .reduce((a, b) => a[1] > b[1] ? a : b)[0] as 'few' | 'many' | 'most' | 'all';

  const getTFWRatingFromRatio = (ratio: string) => {
    switch (ratio) {
      case 'all':
        return { rating: "red", color: "text-white", label: "High TFW Usage", percentage: 80 };
      case 'most':
        return { rating: "red", color: "text-white", label: "High TFW Usage", percentage: 65 };
      case 'many':
        return { rating: "yellow", color: "text-yellow-600", label: "Moderate TFW Usage", percentage: 40 };
      case 'few':
      default:
        return { rating: "green", color: "text-green-600", label: "Low TFW Usage", percentage: 15 };
    }
  };

  const businessRating = getTFWRatingFromRatio(mostCommonRatio);
  const verifiedReports = reports.length;

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
            <Link to="/reports" className="hover:text-slate-900">
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
                <div className="grid grid-cols-4 gap-6">
                  <div className="text-center">
                    <div className="text-3xl font-bold text-slate-900 mb-1 capitalize">
                      {mostCommonRatio}
                    </div>
                    <div className="text-sm text-slate-600">
                      Most Common Ratio
                    </div>
                  </div>
                  <div className="text-center">
                    <div className="text-3xl font-bold text-slate-900 mb-1">
                      {businessRating.percentage}%
                    </div>
                    <div className="text-sm text-slate-600">
                      Estimated TFW Usage
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
                  <div className="text-center">
                    <div className="text-3xl font-bold text-red-600 mb-1">
                      {boycottStats?.boycott_count || 0}
                    </div>
                    <div className="text-sm text-slate-600">
                      People Boycotting
                    </div>
                  </div>
                </div>

                <Separator />

                <div className="space-y-4">
                  <h4 className="font-semibold">TFW Ratio Distribution</h4>
                  <div className="space-y-3">
                    <div className="flex items-center justify-between">
                      <div className="flex items-center space-x-2">
                        <div className="w-3 h-3 bg-red-600 rounded-full"></div>
                        <span className="text-sm">All TFW workers</span>
                      </div>
                      <div className="flex items-center space-x-2">
                        <Progress
                          value={
                            reports.length > 0
                              ? (tfwRatioDistribution.all / reports.length) * 100
                              : 0
                          }
                          className="w-24"
                        />
                        <span className="text-sm font-medium">
                          {tfwRatioDistribution.all} reports
                        </span>
                      </div>
                    </div>
                    <div className="flex items-center justify-between">
                      <div className="flex items-center space-x-2">
                        <div className="w-3 h-3 bg-red-400 rounded-full"></div>
                        <span className="text-sm">Most TFW workers</span>
                      </div>
                      <div className="flex items-center space-x-2">
                        <Progress
                          value={
                            reports.length > 0
                              ? (tfwRatioDistribution.most / reports.length) * 100
                              : 0
                          }
                          className="w-24"
                        />
                        <span className="text-sm font-medium">
                          {tfwRatioDistribution.most} reports
                        </span>
                      </div>
                    </div>
                    <div className="flex items-center justify-between">
                      <div className="flex items-center space-x-2">
                        <div className="w-3 h-3 bg-yellow-500 rounded-full"></div>
                        <span className="text-sm">Many TFW workers</span>
                      </div>
                      <div className="flex items-center space-x-2">
                        <Progress
                          value={
                            reports.length > 0
                              ? (tfwRatioDistribution.many / reports.length) * 100
                              : 0
                          }
                          className="w-24"
                        />
                        <span className="text-sm font-medium">
                          {tfwRatioDistribution.many} reports
                        </span>
                      </div>
                    </div>
                    <div className="flex items-center justify-between">
                      <div className="flex items-center space-x-2">
                        <div className="w-3 h-3 bg-green-500 rounded-full"></div>
                        <span className="text-sm">Few TFW workers</span>
                      </div>
                      <div className="flex items-center space-x-2">
                        <Progress
                          value={
                            reports.length > 0
                              ? (tfwRatioDistribution.few / reports.length) * 100
                              : 0
                          }
                          className="w-24"
                        />
                        <span className="text-sm font-medium">
                          {tfwRatioDistribution.few} reports
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
                    const reportRatio = report.tfw_ratio || 'few';
                    const reportRating = getTFWRatingFromRatio(reportRatio);
                    const isVerified = true; // All reports are now automatically accepted

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
                                Verified
                              </div>
                            </div>
                          </div>
                          <div className="text-xs text-slate-500">
                            {new Date(report.created_at).toLocaleDateString()}
                          </div>
                        </div>

                        <div className="flex items-center space-x-2">
                          <span className="text-sm">TFW ratio observed:</span>
                          <Badge className="capitalize">{reportRatio} TFW workers</Badge>
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
                  <Link 
                    to="/reports/create" 
                    search={{ 
                      businessName: businessName, 
                      businessAddress: decodedAddress 
                    }}
                  >
                    Submit Report
                  </Link>
                </Button>
                <BoycottButton 
                  businessName={businessName}
                  businessAddress={decodedAddress}
                  className="w-full"
                />
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
                  <span className="text-sm">Most Common Ratio</span>
                  <span className="font-bold capitalize text-slate-900">
                    {mostCommonRatio}
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
                  <span className="text-sm">People Boycotting</span>
                  <span className="font-bold text-red-600">
                    {boycottStats?.boycott_count || 0}
                  </span>
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
