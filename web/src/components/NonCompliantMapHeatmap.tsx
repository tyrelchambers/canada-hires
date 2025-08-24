import { useState, useCallback, useMemo } from "react";
import { MapContainer, TileLayer, Marker, Popup, Circle } from "react-leaflet";
import { LatLngExpression } from "leaflet";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faSpinner,
  faExclamationTriangle,
  faInfoCircle,
} from "@fortawesome/free-solid-svg-icons";
import {
  useNonCompliantLocations,
  useNonCompliantByPostalCode,
  useNonCompliantByCoordinates,
} from "@/hooks/useNonCompliant";
import { MapSearch } from "./MapSearch";
import { ReasonCodeTooltip } from "./ReasonCodeTooltip";
import { ReasonCodeModal } from "./ReasonCodeModal";
import type { NonCompliantPostalCodeLocation } from "@/types";
import "leaflet/dist/leaflet.css";

export function NonCompliantMapHeatmap() {
  const [selectedLocation, setSelectedLocation] =
    useState<NonCompliantPostalCodeLocation | null>(null);
  const [selectedPostalCode, setSelectedPostalCode] = useState<string>("");
  const [selectedCoordinates, setSelectedCoordinates] = useState<{
    lat: number | null;
    lng: number | null;
  }>({ lat: null, lng: null });
  const [modalState, setModalState] = useState<{
    isOpen: boolean;
    reasonCodes: string[];
    businessName: string;
  }>({
    isOpen: false,
    reasonCodes: [],
    businessName: "",
  });

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

  // Fetch employers for selected coordinates
  const {
    data: coordinateEmployersData,
    isLoading: coordinateEmployersLoading,
    error: coordinateEmployersError,
  } = useNonCompliantByCoordinates(
    selectedCoordinates.lat,
    selectedCoordinates.lng,
    100,
  );

  // Handle marker click
  const handleMarkerClick = useCallback(
    (location: NonCompliantPostalCodeLocation) => {
      setSelectedLocation(location);

      // Check if this is a coordinate-based location
      if (location.postal_code.startsWith("COORD_")) {
        // Extract coordinates from the COORD_lat_lng format
        const coordParts = location.postal_code
          .replace("COORD_", "")
          .split("_");
        if (coordParts.length === 2) {
          const lat = parseFloat(coordParts[0]);
          const lng = parseFloat(coordParts[1]);
          setSelectedCoordinates({ lat, lng });
          setSelectedPostalCode(""); // Clear postal code query
        }
      } else {
        // Regular postal code
        setSelectedPostalCode(location.postal_code);
        setSelectedCoordinates({ lat: null, lng: null }); // Clear coordinate query
      }
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

  // Handle opening the reason code modal
  const handleShowReasonDetails = useCallback(
    (reasonCodes: string[], businessName: string) => {
      setModalState({
        isOpen: true,
        reasonCodes,
        businessName,
      });
    },
    [],
  );

  // Handle closing the modal
  const handleCloseModal = useCallback(() => {
    setModalState((prev) => ({
      ...prev,
      isOpen: false,
    }));
  }, []);

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
    if (!data?.locations || !Array.isArray(data.locations)) return [];

    return data.locations.map((location) => {
      // Calculate circle radius based on total penalty amount (min 200m, max 800m)
      const circleRadius = Math.min(
        Math.max(200, (location.total_penalty_amount || 0) / 1000 + 100),
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
                <div>
                  {location.employer_count || 0} non-compliant employers
                </div>
                <div>{location.violation_count || 0} violations</div>
                <div className="text-orange-700 font-medium">
                  {formatCurrency(location.total_penalty_amount || 0)} in
                  penalties
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
    <div className=" h-full flex flex-col lg:flex-row">
      {/* Sidebar - Desktop only */}
      <div className="hidden lg:flex w-[400px] bg-gray-50 border-r border-gray-200 flex-col overflow-hidden">
        <div className="p-4 border-b border-gray-200 bg-white">
          <h1 className="text-xl font-bold text-gray-900">
            Non-Compliant Employers Map
          </h1>
          <p className="text-sm text-gray-600 mt-1">
            Explore labor violations across Canada since 2016.
          </p>
        </div>

        <div className="flex-1 overflow-y-auto p-4 space-y-4">
          {selectedPostalCode ||
          (selectedCoordinates.lat !== null &&
            selectedCoordinates.lng !== null) ? (
            /* Employer Detail View */
            <>
              {/* Back Button */}
              <div className="mb-4">
                <button
                  onClick={() => {
                    setSelectedPostalCode("");
                    setSelectedCoordinates({ lat: null, lng: null });
                    setSelectedLocation(null);
                  }}
                  className="flex items-center gap-2 text-gray-600 hover:text-gray-900 text-sm font-medium"
                >
                  ← Back to Map Overview
                </button>
              </div>

              {/* Location Header */}
              <div className="bg-white rounded-lg shadow p-4">
                <h2 className="text-xl font-bold text-gray-900 mb-2 flex items-center gap-2">
                  <FontAwesomeIcon
                    icon={faExclamationTriangle}
                    className="text-orange-600"
                  />
                  {selectedPostalCode ||
                    (selectedCoordinates.lat !== null &&
                    selectedCoordinates.lng !== null
                      ? `Coordinates: ${selectedCoordinates.lat?.toFixed(6)}, ${selectedCoordinates.lng?.toFixed(6)}`
                      : "Unknown Location")}
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
                {(() => {
                  // Determine which data source to use based on whether we have coordinates or postal code
                  const isCoordinateBased =
                    selectedCoordinates.lat !== null &&
                    selectedCoordinates.lng !== null;
                  const currentEmployersData = isCoordinateBased
                    ? coordinateEmployersData
                    : employersData;
                  const currentLoading = isCoordinateBased
                    ? coordinateEmployersLoading
                    : employersLoading;
                  const currentError = isCoordinateBased
                    ? coordinateEmployersError
                    : employersError;

                  if (currentLoading) {
                    return (
                      <div className="flex items-center gap-2 text-gray-500 text-sm">
                        <FontAwesomeIcon
                          icon={faSpinner}
                          className="animate-spin"
                        />
                        <span>Loading employers...</span>
                      </div>
                    );
                  }

                  if (currentError) {
                    return (
                      <div className="text-red-600 text-sm">
                        Error loading employers: {currentError.message}
                      </div>
                    );
                  }

                  if (
                    currentEmployersData &&
                    Array.isArray(currentEmployersData.employers) &&
                    currentEmployersData.employers.length > 0
                  ) {
                    return (
                      <div className="space-y-3">
                        {currentEmployersData.employers.map((employer) => (
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
                              Array.isArray(employer.reason_codes) &&
                              employer.reason_codes.length > 0 && (
                                <div className="text-sm text-gray-700 mb-2 space-y-1">
                                  <div className="flex items-center justify-between">
                                    <strong>Violation Codes:</strong>
                                    <button
                                      onClick={() =>
                                        handleShowReasonDetails(
                                          employer.reason_codes,
                                          employer.business_operating_name,
                                        )
                                      }
                                      className="text-xs text-blue-600 hover:text-blue-800 font-medium flex items-center gap-1"
                                      title="View detailed violation descriptions"
                                    >
                                      <FontAwesomeIcon
                                        icon={faInfoCircle}
                                        className="w-3 h-3"
                                      />
                                      Details
                                    </button>
                                  </div>
                                  <div className="space-x-1">
                                    {employer.reason_codes.map(
                                      (code, index) => (
                                        <span key={code}>
                                          <ReasonCodeTooltip reasonCode={code}>
                                            <span className="inline-block bg-orange-100 text-orange-800 text-xs px-2 py-1 rounded-full font-medium">
                                              {code}
                                            </span>
                                          </ReasonCodeTooltip>
                                          {index <
                                            employer.reason_codes.length - 1 &&
                                            " "}
                                        </span>
                                      ),
                                    )}
                                  </div>
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
                        {currentEmployersData.count &&
                          currentEmployersData.employers &&
                          currentEmployersData.count >
                            currentEmployersData.employers.length && (
                            <div className="text-sm text-gray-500 text-center py-3 border-t">
                              Showing {currentEmployersData.employers.length} of{" "}
                              {currentEmployersData.count} employers
                            </div>
                          )}
                      </div>
                    );
                  }

                  return (
                    <div className="text-gray-500 text-sm py-8 text-center">
                      No non-compliant employers found for this location.
                    </div>
                  );
                })()}
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
                        {data.locations && Array.isArray(data.locations)
                          ? data.locations.reduce(
                              (sum, loc) => sum + (loc.employer_count || 0),
                              0,
                            )
                          : 0}
                      </span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-gray-600">Total Violations:</span>
                      <span className="font-medium text-orange-600">
                        {data.locations && Array.isArray(data.locations)
                          ? data.locations.reduce(
                              (sum, loc) => sum + (loc.violation_count || 0),
                              0,
                            )
                          : 0}
                      </span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-gray-600">Total Penalties:</span>
                      <span className="font-medium text-orange-700">
                        {formatCurrency(
                          data.locations && Array.isArray(data.locations)
                            ? data.locations.reduce(
                                (sum, loc) =>
                                  sum + (loc.total_penalty_amount || 0),
                                0,
                              )
                            : 0,
                        )}
                      </span>
                    </div>
                  </div>
                ) : null}
              </div>
            </>
          )}
        </div>
        <div className="px-4 pb-2 text-xs text-gray-400">
          Data sourced from Canadian government databases. Location accuracy may vary.
        </div>
      </div>

      {/* Map */}
      <div className=" h-[300px] md:h-[400px] lg:h-full lg:flex-1 z-0">
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
      </div>

      {/* Mobile Bottom Section - Mobile only */}
      <div className="lg:hidden bg-gray-50 border-t border-gray-200">
        <div className="p-4 bg-white border-b border-gray-200">
          <h1 className="text-lg font-bold text-gray-900">
            Non-Compliant Employers Map
          </h1>
          <p className="text-sm text-gray-600 mt-1">
            Explore labor violations across Canada
          </p>
        </div>

        <div className="p-4 space-y-4">
          {selectedPostalCode ||
          (selectedCoordinates.lat !== null &&
            selectedCoordinates.lng !== null) ? (
            /* Employer Detail View */
            <>
              {/* Back Button */}
              <div className="mb-4">
                <button
                  onClick={() => {
                    setSelectedPostalCode("");
                    setSelectedCoordinates({ lat: null, lng: null });
                    setSelectedLocation(null);
                  }}
                  className="flex items-center gap-2 text-gray-600 hover:text-gray-900 text-sm font-medium"
                >
                  ← Back to Map Overview
                </button>
              </div>

              {/* Location Header */}
              <div className="bg-white rounded-lg shadow p-4">
                <h2 className="text-xl font-bold text-gray-900 mb-2 flex items-center gap-2">
                  <FontAwesomeIcon
                    icon={faExclamationTriangle}
                    className="text-orange-600"
                  />
                  {selectedPostalCode ||
                    (selectedCoordinates.lat !== null &&
                    selectedCoordinates.lng !== null
                      ? `Coordinates: ${selectedCoordinates.lat?.toFixed(6)}, ${selectedCoordinates.lng?.toFixed(6)}`
                      : "Unknown Location")}
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
                {(() => {
                  // Determine which data source to use based on whether we have coordinates or postal code
                  const isCoordinateBased =
                    selectedCoordinates.lat !== null &&
                    selectedCoordinates.lng !== null;
                  const currentEmployersData = isCoordinateBased
                    ? coordinateEmployersData
                    : employersData;
                  const currentLoading = isCoordinateBased
                    ? coordinateEmployersLoading
                    : employersLoading;
                  const currentError = isCoordinateBased
                    ? coordinateEmployersError
                    : employersError;

                  if (currentLoading) {
                    return (
                      <div className="flex items-center gap-2 text-gray-500 text-sm">
                        <FontAwesomeIcon
                          icon={faSpinner}
                          className="animate-spin"
                        />
                        <span>Loading employers...</span>
                      </div>
                    );
                  }

                  if (currentError) {
                    return (
                      <div className="text-red-600 text-sm">
                        Error loading employers: {currentError.message}
                      </div>
                    );
                  }

                  if (
                    currentEmployersData &&
                    Array.isArray(currentEmployersData.employers) &&
                    currentEmployersData.employers.length > 0
                  ) {
                    return (
                      <div className="space-y-3">
                        {currentEmployersData.employers.map((employer) => (
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
                              Array.isArray(employer.reason_codes) &&
                              employer.reason_codes.length > 0 && (
                                <div className="text-sm text-gray-700 mb-2 space-y-1">
                                  <div className="flex items-center justify-between">
                                    <strong>Violation Codes:</strong>
                                    <button
                                      onClick={() =>
                                        handleShowReasonDetails(
                                          employer.reason_codes,
                                          employer.business_operating_name,
                                        )
                                      }
                                      className="text-xs text-blue-600 hover:text-blue-800 font-medium flex items-center gap-1"
                                      title="View detailed violation descriptions"
                                    >
                                      <FontAwesomeIcon
                                        icon={faInfoCircle}
                                        className="w-3 h-3"
                                      />
                                      Details
                                    </button>
                                  </div>
                                  <div className="space-x-1">
                                    {employer.reason_codes.map(
                                      (code, index) => (
                                        <span key={code}>
                                          <ReasonCodeTooltip reasonCode={code}>
                                            <span className="inline-block bg-orange-100 text-orange-800 text-xs px-2 py-1 rounded-full font-medium">
                                              {code}
                                            </span>
                                          </ReasonCodeTooltip>
                                          {index <
                                            employer.reason_codes.length - 1 &&
                                            " "}
                                        </span>
                                      ),
                                    )}
                                  </div>
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
                        {currentEmployersData.count &&
                          currentEmployersData.employers &&
                          currentEmployersData.count >
                            currentEmployersData.employers.length && (
                            <div className="text-sm text-gray-500 text-center py-3 border-t">
                              Showing {currentEmployersData.employers.length} of{" "}
                              {currentEmployersData.count} employers
                            </div>
                          )}
                      </div>
                    );
                  }

                  return (
                    <div className="text-gray-500 text-sm py-8 text-center">
                      No non-compliant employers found for this location.
                    </div>
                  );
                })()}
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
                        {data.locations && Array.isArray(data.locations)
                          ? data.locations.reduce(
                              (sum, loc) => sum + (loc.employer_count || 0),
                              0,
                            )
                          : 0}
                      </span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-gray-600">Total Violations:</span>
                      <span className="font-medium text-orange-600">
                        {data.locations && Array.isArray(data.locations)
                          ? data.locations.reduce(
                              (sum, loc) => sum + (loc.violation_count || 0),
                              0,
                            )
                          : 0}
                      </span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-gray-600">Total Penalties:</span>
                      <span className="font-medium text-orange-700">
                        {formatCurrency(
                          data.locations && Array.isArray(data.locations)
                            ? data.locations.reduce(
                                (sum, loc) =>
                                  sum + (loc.total_penalty_amount || 0),
                                0,
                              )
                            : 0,
                        )}
                      </span>
                    </div>
                  </div>
                ) : null}
              </div>
            </>
          )}
        </div>
        <div className="px-4 pb-2 text-xs text-gray-400">
          Data sourced from Canadian government databases. Location accuracy may vary.
        </div>
      </div>

      {/* Reason Code Modal */}
      <ReasonCodeModal
        reasonCodes={modalState.reasonCodes}
        businessName={modalState.businessName}
        isOpen={modalState.isOpen}
        onClose={handleCloseModal}
      />
    </div>
  );
}
