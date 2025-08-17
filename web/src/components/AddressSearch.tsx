import { useState, useEffect, useRef } from "react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faTimes } from "@fortawesome/free-solid-svg-icons";
import mbxGeocoding, {
  GeocodeFeature,
} from "@mapbox/mapbox-sdk/services/geocoding";

interface AddressSearchProps {
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
  id?: string;
  required?: boolean;
}

export function AddressSearch({
  value,
  onChange,
  placeholder = "Search for address...",
  id,
  required = false,
}: AddressSearchProps) {
  const [suggestions, setSuggestions] = useState<GeocodeFeature[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [showDropdown, setShowDropdown] = useState(false);
  const [hasSearched, setHasSearched] = useState(false);
  const [selectedIndex, setSelectedIndex] = useState(-1);
  const [error, setError] = useState<string | null>(null);

  const inputRef = useRef<HTMLInputElement>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const timeoutRef = useRef<NodeJS.Timeout | null>(null);

  // Initialize Mapbox client
  const mapboxToken = import.meta.env.VITE_MAPBOX_ACCESS_TOKEN as string;
  const geocodingClient = mapboxToken
    ? mbxGeocoding({ accessToken: mapboxToken })
    : null;

  // Debounced search function
  const searchAddresses = async (query: string) => {
    if (!geocodingClient || query.length < 3) {
      setSuggestions([]);
      setShowDropdown(false);
      setHasSearched(false);
      return;
    }

    setIsLoading(true);
    setError(null);
    setHasSearched(true);

    try {
      const response = await geocodingClient
        .forwardGeocode({
          query,
          countries: ["ca"], // Limit to Canada
          limit: 5,
          types: ["address", "poi"],
        })
        .send();

      const features = response.body.features || [];
      setSuggestions(features);
      setShowDropdown(true); // Show dropdown when search completes
      setSelectedIndex(-1);
    } catch (err) {
      console.error("Geocoding error:", err);
      setError("Failed to search addresses");
      setSuggestions([]);
      setShowDropdown(true); // Still show to display error/no results
    } finally {
      setIsLoading(false);
    }
  };

  // Handle input change with debouncing
  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const newValue = e.target.value;
    onChange(newValue);

    // Clear existing timeout
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current);
    }

    // If input is empty, clear suggestions immediately
    if (!newValue.trim()) {
      setSuggestions([]);
      setShowDropdown(false);
      setHasSearched(false);
      return;
    }

    // Set new timeout for search
    timeoutRef.current = setTimeout(() => {
      void searchAddresses(newValue);
    }, 300);
  };

  // Handle suggestion selection
  const handleSuggestionSelect = (suggestion: GeocodeFeature) => {
    onChange(suggestion.place_name);
    setShowDropdown(false);
    setSelectedIndex(-1);
    setSuggestions([]);
    setHasSearched(false);
    inputRef.current?.focus();
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
          handleSuggestionSelect(suggestions[selectedIndex]);
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
    onChange("");
    setSuggestions([]);
    setShowDropdown(false);
    setHasSearched(false);
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

  // If no Mapbox token, fallback to regular input
  if (!mapboxToken) {
    return (
      <div>
        <Input
          id={id}
          value={value}
          onChange={(e) => onChange(e.target.value)}
          placeholder={placeholder}
          required={required}
        />
        <p className="text-sm text-yellow-600 mt-1">
          Address search unavailable (missing Mapbox token)
        </p>
      </div>
    );
  }

  return (
    <div ref={containerRef} className="relative">
      <div className="relative">
        <Input
          ref={inputRef}
          id={id}
          value={value}
          onChange={handleInputChange}
          onKeyDown={handleKeyDown}
          placeholder={placeholder}
          required={required}
          className="pr-8"
        />
        {value && (
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

      {/* Custom dropdown */}
      {showDropdown && (
        <div className="absolute z-50 w-full mt-1 bg-white border border-gray-200 rounded-md shadow-lg max-h-60 overflow-y-auto">
          {suggestions.length > 0 ? (
            suggestions.map((suggestion, index) => (
              <button
                key={`${suggestion.center[0]}-${suggestion.center[1]}`}
                type="button"
                className={`w-full px-4 py-3 text-left hover:bg-gray-50 focus:bg-gray-50 focus:outline-none border-b border-gray-100 last:border-b-0 ${
                  index === selectedIndex ? "bg-blue-50" : ""
                }`}
                onClick={() => handleSuggestionSelect(suggestion)}
              >
                <div className="text-sm font-medium text-gray-900">
                  {suggestion.place_name}
                </div>
              </button>
            ))
          ) : hasSearched && !isLoading ? (
            <div className="px-4 py-3 text-sm text-gray-500">
              No results found
            </div>
          ) : null}
        </div>
      )}

      {/* Error message */}
      {error && <p className="text-sm text-red-600 mt-1">{error}</p>}
    </div>
  );
}
