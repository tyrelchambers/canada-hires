import { createFileRoute, Link } from "@tanstack/react-router";
import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faArrowLeft,
  faMapMarkerAlt,
  faCalendar,
  faUser,
  faFileText,
  faStar,
  faStarHalfAlt,
} from "@fortawesome/free-solid-svg-icons";
import { faStar as faStarEmpty } from "@fortawesome/free-regular-svg-icons";
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

  const getStatusColor = (status: string) => {
    switch (status) {
      case "approved":
        return "bg-green-100 text-green-800 border-green-200";
      case "pending":
        return "bg-yellow-100 text-yellow-800 border-yellow-200";
      case "rejected":
        return "bg-red-100 text-red-800 border-red-200";
      case "flagged":
        return "bg-orange-100 text-orange-800 border-orange-200";
      default:
        return "bg-gray-100 text-gray-800 border-gray-200";
    }
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

  const renderStars = (rating: number) => {
    const stars = [];
    const fullStars = Math.floor(rating);
    const hasHalfStar = rating % 1 !== 0;

    for (let i = 0; i < fullStars; i++) {
      stars.push(
        <FontAwesomeIcon key={i} icon={faStar} className="text-yellow-400" />
      );
    }

    if (hasHalfStar) {
      stars.push(
        <FontAwesomeIcon
          key="half"
          icon={faStarHalfAlt}
          className="text-yellow-400"
        />
      );
    }

    const emptyStars = 5 - Math.ceil(rating);
    for (let i = 0; i < emptyStars; i++) {
      stars.push(
        <FontAwesomeIcon
          key={`empty-${i}`}
          icon={faStarEmpty}
          className="text-gray-300"
        />
      );
    }

    return stars;
  };

  const getTFWRating = (confidenceLevel: number) => {
    if (confidenceLevel >= 8)
      return { rating: 2, color: "text-red-600", label: "High TFW Usage" };
    if (confidenceLevel >= 5)
      return { rating: 3, color: "text-yellow-600", label: "Moderate TFW Usage" };
    return { rating: 4, color: "text-green-600", label: "Low TFW Usage" };
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50">
        <AuthNav />
        <div className="container mx-auto px-4 py-8">
          <div className="text-center py-12">
            <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-red-600"></div>
            <p className="mt-4 text-gray-600 text-lg">Loading business details...</p>
          </div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen bg-gray-50">
        <AuthNav />
        <div className="container mx-auto px-4 py-8">
          <Card className="border-red-200 bg-red-50">
            <CardContent className="p-6">
              <p className="text-red-800">
                Error: {error instanceof Error ? error.message : "An error occurred"}
              </p>
            </CardContent>
          </Card>
        </div>
      </div>
    );
  }

  const businessName = reports.length > 0 ? reports[0].business_name : "Unknown Business";
  const averageConfidence = reports.length > 0 
    ? reports.reduce((sum, report) => sum + (report.confidence_level || 0), 0) / reports.length
    : 0;

  return (
    <div className="min-h-screen bg-gray-50">
      <AuthNav />
      
      <div className="container mx-auto px-4 py-8">
        {/* Back to Directory */}
        <div className="mb-6">
          <Link to="/directory">
            <Button variant="outline" className="mb-4">
              <FontAwesomeIcon icon={faArrowLeft} className="mr-2" />
              Back to Directory
            </Button>
          </Link>
        </div>

        {/* Business Header */}
        <Card className="mb-8">
          <CardContent className="p-8">
            <div className="flex flex-col lg:flex-row lg:justify-between lg:items-start gap-6">
              <div className="flex-1">
                <h1 className="text-3xl font-bold text-gray-900 mb-2">
                  {businessName}
                </h1>
                
                <div className="flex items-center text-gray-600 mb-4">
                  <FontAwesomeIcon icon={faMapMarkerAlt} className="mr-2" />
                  <span className="text-lg">{decodedAddress}</span>
                </div>

                {reports.length > 0 && (
                  <div className="flex items-center gap-3">
                    <div className="flex items-center">
                      {renderStars(getTFWRating(averageConfidence).rating)}
                    </div>
                    <span className={`font-medium text-lg ${getTFWRating(averageConfidence).color}`}>
                      {getTFWRating(averageConfidence).label}
                    </span>
                  </div>
                )}
              </div>

              <div className="text-right">
                <div className="text-3xl font-bold text-blue-600 mb-1">
                  {reports.length}
                </div>
                <div className="text-gray-600">
                  {reports.length === 1 ? 'Report' : 'Reports'}
                </div>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Individual Reports */}
        <div className="space-y-6">
          <h2 className="text-2xl font-bold text-gray-900">All Reports</h2>
          
          {reports.length === 0 ? (
            <Card className="text-center py-12">
              <CardContent>
                <p className="text-gray-500 text-lg">
                  No reports found for this business address.
                </p>
              </CardContent>
            </Card>
          ) : (
            reports.map((report) => {
              const tfwRating = getTFWRating(report.confidence_level || 0);
              
              return (
                <Card key={report.id}>
                  <CardContent className="p-6">
                    <div className="flex flex-col lg:flex-row lg:justify-between lg:items-start gap-4">
                      <div className="flex-1">
                        <div className="flex items-center gap-3 mb-3">
                          <div className="flex items-center">
                            {renderStars(tfwRating.rating)}
                          </div>
                          <span className={`font-medium ${tfwRating.color}`}>
                            {tfwRating.label}
                          </span>
                        </div>

                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
                          <div className="flex items-center text-gray-600">
                            <FontAwesomeIcon icon={faCalendar} className="mr-2" />
                            <span>Reported: {new Date(report.created_at).toLocaleDateString()}</span>
                          </div>
                          <div className="flex items-center text-gray-600">
                            <FontAwesomeIcon icon={faUser} className="mr-2" />
                            <span>Source: {report.report_source}</span>
                          </div>
                        </div>

                        {report.additional_notes && (
                          <div className="mt-4">
                            <div className="flex items-start text-gray-600">
                              <FontAwesomeIcon icon={faFileText} className="mr-2 mt-1" />
                              <div>
                                <span className="font-medium block mb-1">Additional Notes:</span>
                                <p className="text-gray-700">{report.additional_notes}</p>
                              </div>
                            </div>
                          </div>
                        )}
                      </div>

                      <Badge
                        className={getStatusColor(report.status)}
                        variant="outline"
                      >
                        {getStatusLabel(report.status)}
                      </Badge>
                    </div>
                  </CardContent>
                </Card>
              );
            })
          )}
        </div>
      </div>
    </div>
  );
}

export const Route = createFileRoute("/business/$address")({
  component: BusinessDetailPage,
});