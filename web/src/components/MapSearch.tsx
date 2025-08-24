import { useState, useEffect, useRef } from "react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faTimes, faMapMarkerAlt } from "@fortawesome/free-solid-svg-icons";
// Pelias geocoding response interfaces
interface PeliasGeometry {
  type: string;
  coordinates: [number, number]; // [longitude, latitude]
}

interface PeliasProperties {
  id: string;
  gid: string;
  layer: string;
  source: string;
  source_id: string;
  country_code?: string;
  name: string;
  postalcode?: string;
  confidence: number;
  match_type?: string;
  distance?: number;
  accuracy?: string;
  country?: string;
  country_gid?: string;
  country_a?: string;
  region?: string;
  region_gid?: string;
  region_a?: string;
  locality?: string;
  locality_gid?: string;
  label?: string;
}

interface PeliasFeature {
  type: string;
  geometry: PeliasGeometry;
  properties: PeliasProperties;
  bbox?: number[];
}

interface PeliasResponse {
  geocoding: {
    version: string;
    attribution: string;
    query: Record<string, any>;
    warnings?: string[];
    errors?: string[];
    engine: Record<string, any>;
    timestamp: number;
  };
  type: string;
  features: PeliasFeature[];
  bbox?: number[];
}

interface MapSearchProps {
  onLocationSelect: (location: {
    lng: number;
    lat: number;
    name: string;
  }) => void;
  placeholder?: string;
}

export function MapSearch({
  onLocationSelect,
  placeholder = "Search for a city...",
}: MapSearchProps) {
  const [query, setQuery] = useState("");
  const [suggestions, setSuggestions] = useState<PeliasFeature[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [showDropdown, setShowDropdown] = useState(false);
  const [selectedIndex, setSelectedIndex] = useState(-1);

  const inputRef = useRef<HTMLInputElement>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const timeoutRef = useRef<NodeJS.Timeout | null>(null);

  // Get Pelias server URL from environment
  const peliasServerURL =
    import.meta.env.VITE_PELIAS_SERVER_URL || "http://homeserver:4000";

  // Debounced search function
  const searchCities = async (searchQuery: string) => {
    if (searchQuery.length < 2) {
      setSuggestions([]);
      setShowDropdown(false);
      return;
    }

    setIsLoading(true);

    try {
      const response = await fetch(
        `${peliasServerURL}/v1/search?text=${encodeURIComponent(searchQuery)}&size=8&layers=locality,region`,
      );

      if (!response.ok) {
        throw new Error(`HTTP ${response.status}`);
      }

      const results: PeliasResponse = await response.json();
      setSuggestions(results.features || []);
      setShowDropdown(true);
      setSelectedIndex(-1);
    } catch (err) {
      console.error("City search error:", err);
      setSuggestions([]);
      setShowDropdown(false);
    } finally {
      setIsLoading(false);
    }
  };

  // Handle input change with debouncing
  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const newValue = e.target.value;
    setQuery(newValue);

    // Clear existing timeout
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current);
    }

    // If input is empty, clear suggestions immediately
    if (!newValue.trim()) {
      setSuggestions([]);
      setShowDropdown(false);
      return;
    }

    // Set new timeout for search
    timeoutRef.current = setTimeout(() => {
      void searchCities(newValue);
    }, 300);
  };

  // Handle city selection
  const handleCitySelect = (city: PeliasFeature) => {
    const [lng, lat] = city.geometry.coordinates;
    onLocationSelect({
      lng,
      lat,
      name: city.properties.label || city.properties.name,
    });

    setQuery(city.properties.label || city.properties.name);
    setShowDropdown(false);
    setSelectedIndex(-1);
    setSuggestions([]);
    inputRef.current?.blur();
  };

  // Handle keyboard navigation
  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (!showDropdown || suggestions.length === 0) return;

    switch (e.key) {
      case "ArrowDown":
        e.preventDefault();
        setSelectedIndex((prev) =>
          prev < suggestions.length - 1 ? prev + 1 : prev,
        );
        break;
      case "ArrowUp":
        e.preventDefault();
        setSelectedIndex((prev) => (prev > 0 ? prev - 1 : -1));
        break;
      case "Enter":
        e.preventDefault();
        if (selectedIndex >= 0) {
          handleCitySelect(suggestions[selectedIndex]);
        }
        break;
      case "Escape":
        setShowDropdown(false);
        setSelectedIndex(-1);
        break;
    }
  };

  // Clear input
  const handleClear = () => {
    setQuery("");
    setSuggestions([]);
    setShowDropdown(false);
    setSelectedIndex(-1);
    inputRef.current?.focus();
  };

  // Close dropdown when clicking outside
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        containerRef.current &&
        !containerRef.current.contains(event.target as Node)
      ) {
        setShowDropdown(false);
        setSelectedIndex(-1);
      }
    };

    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  // Cleanup timeout on unmount
  useEffect(() => {
    return () => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
    };
  }, []);

  return (
    <div>
      <div ref={containerRef} className="relative">
        <div className="relative">
          <Input
            ref={inputRef}
            value={query}
            onChange={handleInputChange}
            onKeyDown={handleKeyDown}
            placeholder={placeholder}
            className="pr-8"
          />
          {query && (
            <Button
              type="button"
              variant="ghost"
              size="sm"
              className="absolute right-1 top-1/2 -translate-y-1/2 h-6 w-6 p-0 hover:bg-gray-100"
              onClick={handleClear}
            >
              <FontAwesomeIcon icon={faTimes} className="h-3 w-3" />
            </Button>
          )}
        </div>

        {/* Loading indicator */}
        {isLoading && (
          <div className="absolute right-10 top-1/2 -translate-y-1/2">
            <div className="animate-spin h-4 w-4 border-2 border-gray-300 border-t-blue-600 rounded-full"></div>
          </div>
        )}

        {/* Search results dropdown */}
        {showDropdown && (
          <div className="absolute z-50 w-full mt-1 bg-white border border-gray-200 rounded-md shadow-lg max-h-60 overflow-y-auto">
            {suggestions.length > 0 ? (
              suggestions.map((city, index) => (
                <button
                  key={city.properties.id}
                  type="button"
                  className={`w-full px-4 py-3 text-left hover:bg-gray-50 focus:bg-gray-50 focus:outline-none border-b border-gray-100 last:border-b-0 ${
                    index === selectedIndex ? "bg-blue-50" : ""
                  }`}
                  onClick={() => handleCitySelect(city)}
                >
                  <div className="flex items-start gap-3">
                    <FontAwesomeIcon
                      icon={faMapMarkerAlt}
                      className="text-gray-400 mt-0.5 flex-shrink-0"
                    />
                    <div>
                      <div className="text-sm font-medium text-gray-900">
                        {city.properties.name}
                      </div>
                      <div className="text-xs text-gray-500">
                        {city.properties.label || city.properties.name}
                      </div>
                    </div>
                  </div>
                </button>
              ))
            ) : query.length >= 2 && !isLoading ? (
              <div className="px-4 py-3 text-sm text-gray-500">
                No cities found
              </div>
            ) : null}
          </div>
        )}
      </div>

      <div className="mt-4 text-xs text-gray-500">
        <p>Search for Canadian cities to quickly navigate the map.</p>
        <p>Click on a city to fly to that location.</p>
      </div>
    </div>
  );
}
