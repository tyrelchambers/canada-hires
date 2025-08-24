import { useState, useCallback, useMemo } from "react";
import { Marker, Popup, Circle } from "react-leaflet";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faMapMarkerAlt } from "@fortawesome/free-solid-svg-icons";
import {
  useLMIAPostalCodeLocations,
  useLMIAEmployersByPostalCode,
} from "@/hooks/useLMIA";
import { MapSearch } from "./MapSearch";
import { QuarterSelector } from "./QuarterSelector";
import { MapSidebar } from "./shared/MapSidebar";
import { MapMobileSection } from "./shared/MapMobileSection";
import { MapContainer } from "./shared/MapContainer";
import { MapBackButton } from "./shared/MapBackButton";
import { MapLoadingOverlay } from "./shared/MapLoadingOverlay";
import { MapLoadingSpinner } from "./shared/MapLoadingSpinner";
import type { PostalCodeLocation } from "@/types";
import "leaflet/dist/leaflet.css";

export function LMIAMapHeatmap() {
  const currentYear = new Date().getFullYear();
  const [year, setYear] = useState(currentYear);
  const [quarter, setQuarter] = useState<string | undefined>(undefined);
  const [selectedLocation, setSelectedLocation] =
    useState<PostalCodeLocation | null>(null);
  const [selectedPostalCode, setSelectedPostalCode] = useState<string>("");

  const [mapRef, setMapRef] = useState<L.Map | null>(null);

  // Fetch LMIA postal code data
  const { data, isLoading, error } = useLMIAPostalCodeLocations(
    year,
    quarter,
    1000,
  );

  // Fetch businesses for selected postal code
  const {
    data: businessesData,
    isLoading: businessesLoading,
    error: businessesError,
  } = useLMIAEmployersByPostalCode(selectedPostalCode, year, quarter, 100);

  // Handle marker click
  const handleMarkerClick = useCallback((location: PostalCodeLocation) => {
    setSelectedLocation(location);
    setSelectedPostalCode(location.postal_code);
  }, []);

  // Handle location search
  const handleLocationSelect = useCallback(
    (location: { lng: number; lat: number; name: string }) => {
      if (mapRef) {
        mapRef.setView([location.lat, location.lng], 11);
      }
    },
    [mapRef],
  );

  // Handle quarter change - convert empty string to undefined
  const handleQuarterChange = useCallback((newQuarter?: string) => {
    setQuarter(newQuarter === "" ? undefined : newQuarter);
  }, []);

  // Handle back button
  const handleBack = useCallback(() => {
    setSelectedPostalCode("");
    setSelectedLocation(null);
  }, []);

  // Memoize markers and circles for better performance
  const markersAndCircles = useMemo(() => {
    if (!data?.locations) return [];

    return data.locations.map((location) => {
      // Calculate circle radius based on total LMIAs (min 50m, max 500m)
      const circleRadius = Math.min(
        Math.max(300, location.total_lmias * 10),
        700,
      );

      return (
        <div key={location.postal_code}>
          {/* Circle overlay showing postal code coverage area */}
          <Circle
            center={[location.latitude, location.longitude]}
            radius={circleRadius}
            pathOptions={{
              color: "#dc2626",
              fillColor: "#dc2626",
              fillOpacity: 0.3,
              weight: 2,
              opacity: 0.3,
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
                <div className="font-semibold">{location.postal_code}</div>
                <div>{location.business_count} businesses</div>
                <div>{location.total_lmias} LMIAs</div>
              </div>
            </Popup>
          </Marker>
        </div>
      );
    });
  }, [data?.locations, handleMarkerClick]);

  const sidebarContent = selectedPostalCode ? (
    /* Business Detail View */
    <>
      <MapBackButton onBack={handleBack} />

      {/* Postal Code Header */}
      <div className="bg-white rounded-lg shadow p-4">
        <h2 className="text-xl font-bold text-gray-900 mb-2">
          {selectedPostalCode}
        </h2>
        <div className="grid grid-cols-1 gap-3 text-sm">
          <div className="flex justify-between">
            <span className="text-gray-600">Businesses:</span>
            <span className="font-semibold text-lg">
              {selectedLocation?.business_count || 0}
            </span>
          </div>
          <div className="flex justify-between">
            <span className="text-gray-600">Total LMIAs:</span>
            <span className="font-semibold text-lg text-red-600">
              {selectedLocation?.total_lmias || 0}
            </span>
          </div>
        </div>
      </div>

      {/* Business List */}
      <div className="bg-white rounded-lg shadow p-4">
        <h3 className="font-semibold mb-3">Businesses with LMIAs</h3>
        {businessesLoading ? (
          <MapLoadingSpinner text="Loading businesses..." size="sm" />
        ) : businessesError ? (
          <div className="text-red-600 text-sm">
            Error loading businesses: {businessesError.message}
          </div>
        ) : businessesData && businessesData.employers.length > 0 ? (
          <div className="space-y-3">
            {businessesData.employers.map((employer) => (
              <div
                key={employer.id}
                className="bg-gray-50 rounded p-3 border-l-4 border-red-500"
              >
                <div className="font-medium text-gray-900 mb-2">
                  {employer.employer}
                </div>
                <div className="text-gray-600 mb-2 text-sm">
                  {employer.address || "Address not available"}
                </div>
                <div className="text-sm text-gray-700 mb-2">
                  <strong>Occupation:</strong>{" "}
                  {employer.occupation || "Not specified"}
                </div>
                <div className="flex items-center justify-between text-sm">
                  <span className="text-red-600 font-semibold">
                    {employer.approved_lmias || 0} LMIAs Approved
                  </span>
                  <span className="text-gray-600">
                    {employer.approved_positions || 0} Positions
                  </span>
                </div>
              </div>
            ))}
            {businessesData.count > businessesData.employers.length && (
              <div className="text-sm text-gray-500 text-center py-3 border-t">
                Showing {businessesData.employers.length} of{" "}
                {businessesData.count} businesses
              </div>
            )}
          </div>
        ) : (
          <div className="text-gray-500 text-sm py-8 text-center">
            No businesses found for this postal code.
          </div>
        )}
      </div>
    </>
  ) : (
    /* Default Map Overview */
    <>
      {/* Quarter Selector */}
      <QuarterSelector
        year={year}
        quarter={quarter}
        onYearChange={setYear}
        onQuarterChange={handleQuarterChange}
      />

      {/* City Search */}
      <div className="bg-white rounded-lg shadow p-4">
        <h3 className="font-semibold mb-3">Search Location</h3>
        <MapSearch onLocationSelect={handleLocationSelect} />
      </div>

      {/* Stats */}
      <div className="bg-white rounded-lg shadow p-4">
        <h3 className="font-semibold mb-3">Current Data</h3>
        {isLoading ? (
          <MapLoadingSpinner />
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
              <span className="text-gray-600">Businesses:</span>
              <span className="font-medium">
                {data.locations?.reduce(
                  (sum, loc) => sum + loc.business_count,
                  0,
                ) || 0}
              </span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-600">Period:</span>
              <span className="font-medium">
                {quarter || "All quarters"} {year}
              </span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-600">Total LMIAs:</span>
              <span className="font-medium text-red-600">
                {data.locations?.reduce(
                  (sum, loc) => sum + loc.total_lmias,
                  0,
                ) || 0}
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
            <div className="w-6 h-6 bg-red-600 rounded-full flex items-center justify-center relative">
              <FontAwesomeIcon
                icon={faMapMarkerAlt}
                className="text-white h-2 w-2"
              />
              <span className="absolute -top-1 -right-1 bg-yellow-500 text-black text-xs rounded-full w-3 h-3 flex items-center justify-center font-bold">
                2
              </span>
            </div>
            <span>Postal Code Area</span>
          </div>

          <div className="flex items-center gap-2">
            <div className="w-6 h-6 border-2 border-red-600 rounded-full opacity-30 bg-red-600"></div>
            <span>Coverage Radius</span>
          </div>

          <div className="text-xs text-gray-500 space-y-1">
            <p>• Markers show postal code locations with LMIA businesses</p>
            <p>• Number badges indicate business count in that area</p>
            <p>• Larger markers = more businesses</p>
            <p>• Circles show approximate postal code coverage</p>
            <p>• Click markers to see all businesses in that area</p>
          </div>
        </div>
      </div>
    </>
  );

  return (
    <div className="h-full flex flex-col lg:flex-row">
      <MapSidebar
        title="LMIA Heatmap"
        description="Explore LMIA approvals across Canada"
        footer="Data sourced from Canadian government databases. Location accuracy may vary."
      >
        {sidebarContent}
      </MapSidebar>

      {/* Map */}
      <div className=" h-[300px] md:h-[400px] lg:h-full lg:flex-1 z-0">
        <MapContainer mapRef={setMapRef}>
          {/* Postal Code Markers and Circles */}
          {markersAndCircles}
        </MapContainer>

        {/* Loading Overlay */}
        <MapLoadingOverlay
          isLoading={isLoading}
          loadingText="Loading LMIA data..."
        />
      </div>

      {/* Mobile Bottom Section */}
      <MapMobileSection
        title="LMIA Heatmap"
        description="Explore LMIA approvals across Canada"
        footer="Data sourced from Canadian government databases. Location accuracy may vary."
      >
        {sidebarContent}
      </MapMobileSection>
    </div>
  );
}
