import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faMapMarkerAlt,
  faBuilding,
  faUsers,
  faTimes,
  faFileText,
  faSpinner,
} from "@fortawesome/free-solid-svg-icons";
import { Button } from "@/components/ui/button";
import { useLMIAEmployersByPostalCode } from "@/hooks/useLMIA";
import type { PostalCodeLocation } from "@/types";

interface PostalCodePopoverProps {
  location: PostalCodeLocation;
  year: number;
  quarter?: string;
  onClose: () => void;
}

export function PostalCodePopover({ location, year, quarter, onClose }: PostalCodePopoverProps) {
  // Fetch businesses for the selected postal code
  const { 
    data: businessesData, 
    isLoading: businessesLoading, 
    error: businessesError 
  } = useLMIAEmployersByPostalCode(location.postal_code, year, quarter, 100);
  
  console.log("PostalCodePopover rendering:", { 
    location, 
    businessesLoading, 
    businessesError, 
    businessesData 
  });
  
  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" onClick={onClose}>
      <div className="bg-white rounded-lg shadow-lg border-0 overflow-hidden max-w-md mx-4 min-h-96" onClick={(e) => e.stopPropagation()}>
        {/* Header */}
        <div className="bg-red-600 text-white p-3 relative">
          <button
            className="absolute top-1 right-1 h-6 w-6 text-white hover:bg-red-500 rounded bg-red-700"
            onClick={onClose}
          >
            ×
          </button>
          
          <div className="flex items-start gap-2 pr-8">
            <div>📍</div>
            <div>
              <h3 className="font-semibold text-sm leading-tight">
                Postal Code: {location.postal_code}
              </h3>
              <div className="text-red-100 text-xs mt-1">
                {location.business_count} {location.business_count === 1 ? 'Business' : 'Businesses'} with LMIAs
              </div>
            </div>
          </div>
        </div>

        {/* Content */}
        <div className="p-4 space-y-4 bg-white text-black">
          <div className="text-black">
            <h3 className="font-bold text-lg mb-2">Debug Info</h3>
            <p>Postal Code: {location.postal_code}</p>
            <p>Business Count: {location.business_count}</p>
            <p>Total LMIAs: {location.total_lmias}</p>
            <p>Loading: {businessesLoading ? 'Yes' : 'No'}</p>
            <p>Error: {businessesError ? businessesError.message : 'None'}</p>
            <p>Data: {businessesData ? `${businessesData.employers.length} employers` : 'None'}</p>
          </div>
          
          {businessesLoading && (
            <div className="text-blue-600">Loading businesses...</div>
          )}
          
          {businessesError && (
            <div className="text-red-600">Error: {businessesError.message}</div>
          )}
          
          {businessesData && businessesData.employers.length > 0 && (
            <div>
              <h4 className="font-semibold mb-2">Businesses:</h4>
              <div className="space-y-2">
                {businessesData.employers.slice(0, 3).map((business) => (
                  <div key={business.id} className="bg-gray-100 p-2 rounded text-sm">
                    <div className="font-medium">{business.employer}</div>
                    <div className="text-gray-600">{business.occupation || 'N/A'}</div>
                    <div className="text-red-600">{business.approved_lmias || 0} LMIAs</div>
                  </div>
                ))}
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}