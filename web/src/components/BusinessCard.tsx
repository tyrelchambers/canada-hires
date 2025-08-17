import { Badge } from "@/components/ui/badge";
import { CardContent } from "@/components/ui/card";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faMapMarkerAlt,
  faStar,
  faStarHalfAlt,
  faCalendar,
} from "@fortawesome/free-solid-svg-icons";
import { faStar as faStarEmpty } from "@fortawesome/free-regular-svg-icons";
import { Link } from "@tanstack/react-router";

interface Report {
  id: string;
  business_name: string;
  business_address: string;
  additional_notes?: string;
  report_source: string;
  created_at: string;
  status: string;
  confidence_level?: number;
}

interface GroupedBusiness {
  business_name: string;
  business_address: string;
  report_count: number;
  confidence_level: number;
  latest_report: string;
}

interface BusinessCardProps {
  report?: Report;
  business?: GroupedBusiness;
}

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
      <FontAwesomeIcon key={i} icon={faStar} className="text-yellow-400" />,
    );
  }

  if (hasHalfStar) {
    stars.push(
      <FontAwesomeIcon
        key="half"
        icon={faStarHalfAlt}
        className="text-yellow-400"
      />,
    );
  }

  const emptyStars = 5 - Math.ceil(rating);
  for (let i = 0; i < emptyStars; i++) {
    stars.push(
      <FontAwesomeIcon
        key={`empty-${i}`}
        icon={faStarEmpty}
        className="text-gray-300"
      />,
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

export function BusinessCard({ report, business }: BusinessCardProps) {
  // Handle both individual reports and grouped business data
  const isGrouped = !!business;
  const businessName = isGrouped ? business.business_name : report?.business_name;
  const businessAddress = isGrouped ? business.business_address : report?.business_address;
  const confidenceLevel = isGrouped ? business.confidence_level : (report?.confidence_level || 0);
  const reportCount = isGrouped ? business.report_count : 1;
  const latestDate = isGrouped ? business.latest_report : report?.created_at;
  
  const tfwRating = getTFWRating(confidenceLevel);

  const cardContent = (
    <div className="hover:shadow-lg transition-shadow cursor-pointer bg-white rounded-md border border-border">
      <CardContent className="p-6">
        <div className="flex flex-col md:flex-row gap-4">
          {/* Business Info */}
          <div className="flex-1">
            <div className="flex flex-col md:flex-row md:justify-between md:items-start mb-2">
              <div>
                <h3 className="text-xl font-bold text-gray-900 mb-1">
                  {businessName}
                </h3>
                <div className="flex items-center gap-2 mb-2">
                  <div className="flex items-center">
                    {renderStars(tfwRating.rating)}
                  </div>
                  <span className={`font-medium ${tfwRating.color}`}>
                    {tfwRating.label}
                  </span>
                </div>
              </div>

              {/* Show report count for grouped data, status for individual reports */}
              {isGrouped ? (
                <Badge className="bg-blue-100 text-blue-800 border-blue-200" variant="outline">
                  {reportCount} {reportCount === 1 ? 'report' : 'reports'}
                </Badge>
              ) : (
                <Badge
                  className={getStatusColor(report?.status || '')}
                  variant="outline"
                >
                  {getStatusLabel(report?.status || '')}
                </Badge>
              )}
            </div>

            <div className="flex items-center text-gray-600 mb-3">
              <FontAwesomeIcon icon={faMapMarkerAlt} className="mr-2" />
              <span>{businessAddress}</span>
            </div>

            {/* Show latest report date for grouped data */}
            {isGrouped && latestDate && (
              <div className="flex items-center text-gray-500 text-sm">
                <FontAwesomeIcon icon={faCalendar} className="mr-2" />
                <span>Latest report: {new Date(latestDate).toLocaleDateString()}</span>
              </div>
            )}
          </div>
        </div>
      </CardContent>
    </div>
  );

  // Wrap grouped business cards in Link, leave individual report cards as-is
  if (isGrouped && businessAddress) {
    return (
      <Link
        to="/business/$address"
        params={{ address: encodeURIComponent(businessAddress) }}
        className="block"
      >
        {cardContent}
      </Link>
    );
  }

  return cardContent;
}
