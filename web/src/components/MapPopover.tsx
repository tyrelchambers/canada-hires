import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faBuilding,
  faMapMarkerAlt,
  faCalendarAlt,
  faUsers,
  faTimes,
} from "@fortawesome/free-solid-svg-icons";
import { Button } from "@/components/ui/button";
import type { LMIAEmployerGeoLocation } from "@/types";

interface MapPopoverProps {
  employer: LMIAEmployerGeoLocation;
  onClose: () => void;
}

export function MapPopover({ employer, onClose }: MapPopoverProps) {
  if (!employer.latitude || !employer.longitude) {
    return null;
  }

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" onClick={onClose}>
      <div className="bg-white rounded-lg shadow-lg border-0 overflow-hidden max-w-sm mx-4" onClick={(e) => e.stopPropagation()}>
        {/* Header */}
        <div className="bg-red-600 text-white p-3 relative">
          <Button
            variant="ghost"
            size="sm"
            className="absolute top-1 right-1 h-6 w-6 p-0 text-white hover:bg-red-500 hover:text-white"
            onClick={onClose}
          >
            <FontAwesomeIcon icon={faTimes} className="h-3 w-3" />
          </Button>
          
          <div className="flex items-start gap-2 pr-8">
            <FontAwesomeIcon icon={faBuilding} className="mt-1 flex-shrink-0" />
            <div>
              <h3 className="font-semibold text-sm leading-tight">
                {employer.employer}
              </h3>
              <div className="flex items-center gap-1 text-red-100 text-xs mt-1">
                <FontAwesomeIcon icon={faMapMarkerAlt} className="h-3 w-3" />
                <span className="truncate">
                  {employer.address || "Address not available"}
                </span>
              </div>
            </div>
          </div>
        </div>

        {/* Content */}
        <div className="p-3 space-y-3">
          {/* Location Info */}
          <div className="flex items-center justify-between text-sm">
            <span className="text-gray-600">Province:</span>
            <span className="font-medium">
              {employer.province_territory || "N/A"}
            </span>
          </div>

          {/* Period Info */}
          <div className="flex items-center justify-between text-sm">
            <span className="text-gray-600 flex items-center gap-1">
              <FontAwesomeIcon icon={faCalendarAlt} className="h-3 w-3" />
              Period:
            </span>
            <span className="font-medium">
              {employer.quarter} {employer.year}
            </span>
          </div>

          {/* LMIA Stats */}
          <div className="border-t pt-2 space-y-2">
            <div className="flex items-center justify-between text-sm">
              <span className="text-gray-600">Current Period LMIAs:</span>
              <span className="font-semibold text-red-600">
                {employer.approved_lmias || 0}
              </span>
            </div>
            
            <div className="flex items-center justify-between text-sm">
              <span className="text-gray-600">Current Period Positions:</span>
              <span className="font-semibold text-red-600">
                {employer.approved_positions || 0}
              </span>
            </div>

            <div className="flex items-center justify-between text-sm border-t pt-2">
              <span className="text-gray-600 flex items-center gap-1">
                <FontAwesomeIcon icon={faUsers} className="h-3 w-3" />
                Total LMIAs (All Data):
              </span>
              <span className="font-bold text-red-700 text-base">
                {employer.total_lmias}
              </span>
            </div>
          </div>

          {/* Footer */}
          <div className="text-xs text-gray-500 border-t pt-2">
            <p>
              This employer has received {employer.total_lmias} LMIA approvals across all available data periods.
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}