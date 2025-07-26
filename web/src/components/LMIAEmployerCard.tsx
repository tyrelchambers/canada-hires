import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { LMIAEmployer } from "@/types";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faMapMarkerAlt, faBuilding, faUsers, faCalendar } from "@fortawesome/free-solid-svg-icons";

interface LMIAEmployerCardProps {
  employer: LMIAEmployer;
}

export function LMIAEmployerCard({ employer }: LMIAEmployerCardProps) {
  const formatDate = (dateString?: string) => {
    if (!dateString) return null;
    return new Date(dateString).toLocaleDateString();
  };

  const getLocationDisplay = () => {
    const parts = [employer.province_territory].filter(Boolean);
    return parts.join(", ");
  };

  return (
    <Card className="hover:shadow-md transition-shadow">
      <CardHeader className="pb-3">
        <div className="flex justify-between items-start">
          <CardTitle className="text-lg font-semibold text-gray-900 leading-tight">
            {employer.employer}
          </CardTitle>
          {employer.approved_positions && (
            <Badge variant="secondary" className="ml-2 flex items-center gap-1">
              <FontAwesomeIcon icon={faUsers} className="w-3 h-3" />
              {employer.approved_positions}
            </Badge>
          )}
        </div>
        
        {employer.incorporate_status && (
          <p className="text-sm text-gray-600 flex items-center gap-1">
            <FontAwesomeIcon icon={faBuilding} className="w-3 h-3" />
            {employer.incorporate_status}
          </p>
        )}
      </CardHeader>
      
      <CardContent className="pt-0">
        <div className="space-y-2">
          {getLocationDisplay() && (
            <div className="flex items-center gap-1 text-sm text-gray-600">
              <FontAwesomeIcon icon={faMapMarkerAlt} className="w-3 h-3" />
              {getLocationDisplay()}
            </div>
          )}
          
          {employer.address && (
            <p className="text-sm text-gray-600 truncate" title={employer.address}>
              {employer.address}
            </p>
          )}
          
          {employer.occupation && (
            <div className="flex flex-wrap gap-1">
              <Badge variant="outline" className="text-xs">
                {employer.occupation}
              </Badge>
            </div>
          )}
          
          {employer.program_stream && (
            <Badge variant="outline" className="text-xs">
              {employer.program_stream}
            </Badge>
          )}
          
          <div className="flex justify-between items-center pt-2 text-xs text-gray-500">
            {employer.approved_lmias && (
              <div className="flex items-center gap-1">
                <FontAwesomeIcon icon={faCalendar} className="w-3 h-3" />
                LMIAs: {employer.approved_lmias}
              </div>
            )}
            <div>
              Added: {formatDate(employer.created_at)}
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}