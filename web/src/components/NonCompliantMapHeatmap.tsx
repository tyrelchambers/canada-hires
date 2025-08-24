import { useState, useCallback, useMemo } from "react";
import { MapContainer, TileLayer, Marker, Popup, Circle } from "react-leaflet";
import { LatLngExpression } from "leaflet";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faSpinner,
  faExclamationTriangle,
} from "@fortawesome/free-solid-svg-icons";
import {
  useNonCompliantLocations,
  useNonCompliantByPostalCode,
} from "@/hooks/useNonCompliant";
import { MapSearch } from "./MapSearch";
import type { NonCompliantPostalCodeLocation } from "@/types";
import "leaflet/dist/leaflet.css";

export function NonCompliantMapHeatmap() {
  const [selectedLocation, setSelectedLocation] =
    useState<NonCompliantPostalCodeLocation | null>(null);
  const [selectedPostalCode, setSelectedPostalCode] = useState<string>("");

  // Map center for Canada
  const mapCenter: LatLngExpression = [61.0666922, -95.712891];
  const mapZoom = 4;
  const [mapRef, setMapRef] = useState<L.Map | null>(null);

  // Fetch non-compliant location data
  const { data, isLoading, error } = useNonCompliantLocations(2000);

  // Fetch employers for selected postal code
  const {
    data: employersData,
    isLoading: employersLoading,
    error: employersError,
  } = useNonCompliantByPostalCode(selectedPostalCode, 100);

  // Handle marker click
  const handleMarkerClick = useCallback(
    (location: NonCompliantPostalCodeLocation) => {
      setSelectedLocation(location);
      setSelectedPostalCode(location.postal_code);
    },
    [],
  );

  // Handle location search
  const handleLocationSelect = useCallback(
    (location: { lng: number; lat: number; name: string }) => {
      if (mapRef) {
        mapRef.setView([location.lat, location.lng], 11);
      }
    },
    [mapRef],
  );

  // Format currency
  const formatCurrency = (amount: number, currency: string = "CAD") => {
    return new Intl.NumberFormat("en-CA", {
      style: "currency",
      currency,
      minimumFractionDigits: 0,
      maximumFractionDigits: 0,
    }).format(amount);
  };

  // Format date
  const formatDate = (dateString?: string) => {
    if (!dateString) return "Not available";
    return new Date(dateString).toLocaleDateString("en-CA", {
      year: "numeric",
      month: "short",
      day: "numeric",
    });
  };

  // Memoize markers and circles for better performance
  const markersAndCircles = useMemo(() => {
    if (!data?.locations) return [];

    return data.locations.map((location) => {
      // Calculate circle radius based on total penalty amount (min 200m, max 800m)
      const circleRadius = Math.min(
        Math.max(200, location.total_penalty_amount / 1000 + 100),
        800,
      );

      return (
        <div key={location.postal_code}>
          {/* Circle overlay showing postal code coverage area */}
          <Circle
            center={[location.latitude, location.longitude]}
            radius={circleRadius}
            pathOptions={{
              color: "#ea580c",
              fillColor: "#ea580c",
              fillOpacity: 0.3,
              weight: 2,
              opacity: 0.4,
            }}
          />

          {/* Postal code marker */}
          <Marker
            position={[location.latitude, location.longitude]}
            eventHandlers={{
              click: () => handleMarkerClick(location),
            }}
          >
            <Popup>
              <div className="text-sm">
                <div className="font-semibold flex items-center gap-2">
                  <FontAwesomeIcon
                    icon={faExclamationTriangle}
                    className="text-orange-600"
                  />
                  {location.postal_code}
                </div>
                <div>{location.employer_count} non-compliant employers</div>
                <div>{location.violation_count} violations</div>
                <div className="text-orange-700 font-medium">
                  {formatCurrency(location.total_penalty_amount)} in penalties
                </div>
                {location.most_recent_violation && (
                  <div className="text-xs text-gray-600 mt-1">
                    Latest: {formatDate(location.most_recent_violation)}
                  </div>
                )}
              </div>
            </Popup>
          </Marker>
        </div>
      );
    });
  }, [data?.locations, handleMarkerClick]);

  return (
    <div className="h-screen flex">
      {/* Sidebar */}
      <div className="w-80 bg-gray-50 border-r border-gray-200 flex flex-col overflow-hidden">
        <div className="p-4 border-b border-gray-200 bg-white">
          <h1 className="text-xl font-bold text-gray-900">
            Non-Compliant Employers Map
          </h1>
          <p className="text-sm text-gray-600 mt-1">
            Explore labor violations across Canada
          </p>
        </div>

        <div className="flex-1 overflow-y-auto p-4 space-y-4">
          {selectedPostalCode ? (
            /* Employer Detail View */
            <>
              {/* Back Button */}
              <div className="mb-4">
                <button
                  onClick={() => {
                    setSelectedPostalCode("");
                    setSelectedLocation(null);
                  }}
                  className="flex items-center gap-2 text-gray-600 hover:text-gray-900 text-sm font-medium"
                >
                  ← Back to Map Overview
                </button>
              </div>

              {/* Postal Code Header */}
              <div className="bg-white rounded-lg shadow p-4">
                <h2 className="text-xl font-bold text-gray-900 mb-2 flex items-center gap-2">
                  <FontAwesomeIcon
                    icon={faExclamationTriangle}
                    className="text-orange-600"
                  />
                  {selectedPostalCode}
                </h2>
                <div className="grid grid-cols-1 gap-3 text-sm">
                  <div className="flex justify-between">
                    <span className="text-gray-600">
                      Non-Compliant Employers:
                    </span>
                    <span className="font-semibold text-lg">
                      {selectedLocation?.employer_count || 0}
                    </span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-gray-600">Total Violations:</span>
                    <span className="font-semibold text-lg text-orange-600">
                      {selectedLocation?.violation_count || 0}
                    </span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-gray-600">Total Penalties:</span>
                    <span className="font-semibold text-lg text-orange-700">
                      {selectedLocation
                        ? formatCurrency(selectedLocation.total_penalty_amount)
                        : "$0"}
                    </span>
                  </div>
                  {selectedLocation?.most_recent_violation && (
                    <div className="flex justify-between">
                      <span className="text-gray-600">Most Recent:</span>
                      <span className="font-medium">
                        {formatDate(selectedLocation.most_recent_violation)}
                      </span>
                    </div>
                  )}
                </div>
              </div>

              {/* Employer List */}
              <div className="bg-white rounded-lg shadow p-4">
                <h3 className="font-semibold mb-3">Non-Compliant Employers</h3>
                {employersLoading ? (
                  <div className="flex items-center gap-2 text-gray-500 text-sm">
                    <FontAwesomeIcon
                      icon={faSpinner}
                      className="animate-spin"
                    />
                    <span>Loading employers...</span>
                  </div>
                ) : employersError ? (
                  <div className="text-red-600 text-sm">
                    Error loading employers: {employersError.message}
                  </div>
                ) : employersData && employersData.employers.length > 0 ? (
                  <div className="space-y-3">
                    {employersData.employers.map((employer) => (
                      <div
                        key={employer.id}
                        className="bg-gray-50 rounded p-3 border-l-4 border-orange-500"
                      >
                        <div className="font-medium text-gray-900 mb-2">
                          {employer.business_operating_name}
                        </div>
                        {employer.business_legal_name &&
                          employer.business_legal_name !==
                            employer.business_operating_name && (
                            <div className="text-gray-600 mb-1 text-sm">
                              Legal Name: {employer.business_legal_name}
                            </div>
                          )}
                        <div className="text-gray-600 mb-2 text-sm">
                          {employer.address || "Address not available"}
                        </div>
                        {employer.reason_codes &&
                          employer.reason_codes.length > 0 && (
                            <div className="text-sm text-gray-700 mb-2">
                              <strong>Violation Codes:</strong>{" "}
                              {employer.reason_codes.join(", ")}
                            </div>
                          )}
                        <div className="flex items-center justify-between text-sm">
                          {employer.penalty_amount && (
                            <span className="text-orange-700 font-semibold">
                              {formatCurrency(
                                employer.penalty_amount,
                                employer.penalty_currency,
                              )}
                            </span>
                          )}
                          {employer.date_of_final_decision && (
                            <span className="text-gray-600">
                              Decision:{" "}
                              {formatDate(employer.date_of_final_decision)}
                            </span>
                          )}
                        </div>
                        {employer.status && (
                          <div className="text-xs text-gray-500 mt-1">
                            Status: {employer.status}
                          </div>
                        )}
                      </div>
                    ))}
                    {employersData.count > employersData.employers.length && (
                      <div className="text-sm text-gray-500 text-center py-3 border-t">
                        Showing {employersData.employers.length} of{" "}
                        {employersData.count} employers
                      </div>
                    )}
                  </div>
                ) : (
                  <div className="text-gray-500 text-sm py-8 text-center">
                    No non-compliant employers found for this postal code.
                  </div>
                )}
              </div>
            </>
          ) : (
            /* Default Map Overview */
            <>
              {/* City Search */}
              <div className="bg-white rounded-lg shadow p-4">
                <h3 className="font-semibold mb-3">Search Location</h3>
                <MapSearch onLocationSelect={handleLocationSelect} />
              </div>

              {/* Stats */}
              <div className="bg-white rounded-lg shadow p-4">
                <h3 className="font-semibold mb-3">Current Data</h3>
                {isLoading ? (
                  <div className="flex items-center gap-2 text-gray-500">
                    <FontAwesomeIcon
                      icon={faSpinner}
                      className="animate-spin"
                    />
                    <span>Loading...</span>
                  </div>
                ) : error ? (
                  <div className="text-red-600 text-sm">
                    Error loading data: {error.message}
                  </div>
                ) : data ? (
                  <div className="space-y-2 text-sm">
                    <div className="flex justify-between">
                      <span className="text-gray-600">Postal Codes:</span>
                      <span className="font-medium">{data.count}</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-gray-600">
                        Non-Compliant Employers:
                      </span>
                      <span className="font-medium">
                        {data.locations?.reduce(
                          (sum, loc) => sum + loc.employer_count,
                          0,
                        ) || 0}
                      </span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-gray-600">Total Violations:</span>
                      <span className="font-medium text-orange-600">
                        {data.locations?.reduce(
                          (sum, loc) => sum + loc.violation_count,
                          0,
                        ) || 0}
                      </span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-gray-600">Total Penalties:</span>
                      <span className="font-medium text-orange-700">
                        {formatCurrency(
                          data.locations?.reduce(
                            (sum, loc) => sum + loc.total_penalty_amount,
                            0,
                          ) || 0,
                        )}
                      </span>
                    </div>
                  </div>
                ) : null}
              </div>

              {/* Legend */}
              <div className="bg-white rounded-lg shadow p-4">
                <h3 className="font-semibold mb-3">Legend</h3>
                <div className="space-y-3 text-sm">
                  <div className="flex items-center gap-2">
                    <div className="w-6 h-6 bg-orange-600 rounded-full flex items-center justify-center">
                      <FontAwesomeIcon
                        icon={faExclamationTriangle}
                        className="text-white h-2 w-2"
                      />
                    </div>
                    <span>Non-Compliant Employer Location</span>
                  </div>

                  <div className="flex items-center gap-2">
                    <div className="w-6 h-6 border-2 border-orange-600 rounded-full opacity-30 bg-orange-600"></div>
                    <span>Violation Impact Radius</span>
                  </div>

                  <div className="text-xs text-gray-500 space-y-1">
                    <p>
                      • Markers show postal code locations with non-compliant
                      employers
                    </p>
                    <p>• Circle size indicates penalty amount severity</p>
                    <p>• Orange theme represents labor violations</p>
                    <p>• Click markers to see all violations in that area</p>
                    <p>• Data sourced from official government records</p>
                  </div>
                </div>
              </div>
            </>
          )}
        </div>
      </div>

      {/* Map */}
      <div className="flex-1 relative">
        <MapContainer
          center={mapCenter}
          zoom={mapZoom}
          style={{ height: "100%", width: "100%" }}
          ref={setMapRef}
        >
          <TileLayer
            attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
            url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
          />

          {/* Postal Code Markers and Circles */}
          {markersAndCircles}
        </MapContainer>

        {/* Loading Overlay */}
        {isLoading && (
          <div className="absolute inset-0 bg-black bg-opacity-20 flex items-center justify-center">
            <div className="bg-white rounded-lg p-4 shadow-lg flex items-center gap-3">
              <FontAwesomeIcon
                icon={faSpinner}
                className="animate-spin text-orange-600"
              />
              <span>Loading non-compliant employer data...</span>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
