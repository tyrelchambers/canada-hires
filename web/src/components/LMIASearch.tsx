import { useState, useMemo, useEffect } from "react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Pagination } from "@/components/ui/pagination";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faSearch,
  faMapMarkerAlt,
  faExclamationCircle,
  faBuilding,
  faUsers,
  faCalendar,
} from "@fortawesome/free-solid-svg-icons";
import { useLMIASearch, useLMIALocation } from "@/hooks/useLMIA";

export function LMIASearch() {
  const [searchQuery, setSearchQuery] = useState("");
  const [city, setCity] = useState("");
  const [province, setProvince] = useState("");
  const [year, setYear] = useState("");
  const [searchType, setSearchType] = useState<"name" | "location">("location");
  const [currentPage, setCurrentPage] = useState(1);
  const itemsPerPage = 25;

  // Generate year options (current year back to 2015)
  const currentYear = new Date().getFullYear();
  const yearOptions = Array.from(
    { length: currentYear - 2014 },
    (_, i) => currentYear - i,
  );

  const {
    data: searchResults,
    isLoading: isSearchLoading,
    error: searchError,
  } = useLMIASearch(
    searchType === "name" ? searchQuery : "",
    year || undefined,
  );

  const {
    data: locationResults,
    isLoading: isLocationLoading,
    error: locationError,
  } = useLMIALocation(
    searchType === "location" ? city : "",
    searchType === "location" ? province : "",
    year || undefined,
  );

  const isLoading = isSearchLoading || isLocationLoading;
  const error = searchError || locationError;
  const rawResults = searchType === "name" ? searchResults : locationResults;

  // Pagination logic
  const paginatedResults = useMemo(() => {
    if (!rawResults || !rawResults.employers) return null;

    const totalItems = rawResults.employers.length;
    const totalPages = Math.ceil(totalItems / itemsPerPage);
    const startIndex = (currentPage - 1) * itemsPerPage;
    const endIndex = startIndex + itemsPerPage;
    const employers = rawResults.employers.slice(startIndex, endIndex);

    return {
      ...rawResults,
      employers,
      totalItems,
      totalPages,
    };
  }, [rawResults, currentPage, itemsPerPage]);

  const results = paginatedResults;

  // Reset to first page when search parameters change
  useEffect(() => {
    setCurrentPage(1);
  }, [searchQuery, city, province, year, searchType]);

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    // The queries will automatically trigger due to dependency changes
  };

  const handlePageChange = (page: number) => {
    setCurrentPage(page);
  };

  const canadianProvinces = [
    "Alberta",
    "British Columbia",
    "Manitoba",
    "New Brunswick",
    "Newfoundland and Labrador",
    "Northwest Territories",
    "Nova Scotia",
    "Nunavut",
    "Ontario",
    "Prince Edward Island",
    "Quebec",
    "Saskatchewan",
    "Yukon",
  ];

  return (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <FontAwesomeIcon icon={faSearch} className="w-5 h-5" />
            Search LMIA Employers
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            {/* Search Type Toggle */}
            <div className="flex space-x-2">
              <Button
                type="button"
                variant={searchType === "name" ? "default" : "outline"}
                onClick={() => setSearchType("name")}
                size="sm"
              >
                <FontAwesomeIcon icon={faSearch} className="w-4 h-4 mr-1" />
                By Name
              </Button>
              <Button
                type="button"
                variant={searchType === "location" ? "default" : "outline"}
                onClick={() => setSearchType("location")}
                size="sm"
              >
                <FontAwesomeIcon
                  icon={faMapMarkerAlt}
                  className="w-4 h-4 mr-1"
                />
                By Location
              </Button>
            </div>

            <form onSubmit={handleSearch} className="space-y-4">
              {searchType === "name" ? (
                <div>
                  <Input
                    placeholder="Search by employer or business name..."
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                    className="w-full"
                  />
                  <p className="text-sm text-gray-500 mt-1">
                    Enter at least 2 characters to search
                  </p>
                </div>
              ) : (
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div>
                    <Input
                      placeholder="City (optional)"
                      value={city}
                      onChange={(e) => setCity(e.target.value)}
                    />
                  </div>
                  <div>
                    <select
                      value={province}
                      onChange={(e) => setProvince(e.target.value)}
                      className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                    >
                      <option value="">Select Province (optional)</option>
                      {canadianProvinces.map((prov) => (
                        <option key={prov} value={prov}>
                          {prov}
                        </option>
                      ))}
                    </select>
                  </div>
                </div>
              )}

              {/* Year filter - available for both search types */}
              <div>
                <select
                  value={year}
                  onChange={(e) => setYear(e.target.value)}
                  className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                >
                  <option value="">All Years</option>
                  {yearOptions.map((yearOption) => (
                    <option key={yearOption} value={yearOption.toString()}>
                      {yearOption}
                    </option>
                  ))}
                </select>
                <p className="text-sm text-gray-500 mt-1">
                  Filter by year (defaults to {currentYear} when no search
                  criteria)
                </p>
              </div>
            </form>
          </div>
        </CardContent>
      </Card>

      {/* Results */}
      {isLoading && (
        <Card>
          <CardContent className="py-8 text-center">
            <div className="flex items-center justify-center space-x-2">
              <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-gray-900"></div>
              <span>Searching...</span>
            </div>
          </CardContent>
        </Card>
      )}

      {error && (
        <Card>
          <CardContent className="py-8 text-center">
            <FontAwesomeIcon
              icon={faExclamationCircle}
              className="w-8 h-8 text-red-500 mx-auto mb-2"
            />
            <p className="text-red-600">
              Error loading results. Please try again.
            </p>
          </CardContent>
        </Card>
      )}

      {results && results.employers && results.employers?.length > 0 && (
        <div className="space-y-4">
          <div className="flex items-center justify-between">
            <h3 className="text-lg font-semibold">
              Search Results ({results.totalItems || results.count} found)
            </h3>
          </div>

          <div className="bg-white rounded-lg shadow-sm border">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Employer</TableHead>
                  <TableHead>Location</TableHead>
                  <TableHead>Occupation</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>Positions</TableHead>
                  <TableHead>LMIAs</TableHead>
                  <TableHead>Program</TableHead>
                  <TableHead>Year</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {results.employers.map((employer) => (
                  <TableRow key={employer.id}>
                    <TableCell className="font-medium min-w-[200px]">
                      <div className="w-[300px]">
                        <div
                          className="font-semibold text-gray-900 truncate"
                          title={employer.employer}
                        >
                          {employer.employer}
                        </div>
                        {employer.address && (
                          <div
                            className="text-xs text-gray-500 truncate"
                            title={employer.address}
                          >
                            {employer.address}
                          </div>
                        )}
                      </div>
                    </TableCell>
                    <TableCell>
                      {employer.province_territory ? (
                        <div className="flex items-center text-sm">
                          <FontAwesomeIcon
                            icon={faMapMarkerAlt}
                            className="mr-1 text-gray-400 w-3 h-3"
                          />
                          <span className="truncate">
                            {employer.province_territory}
                          </span>
                        </div>
                      ) : (
                        <span className="text-gray-400">-</span>
                      )}
                    </TableCell>
                    <TableCell>
                      {employer.occupation ? (
                        <Badge
                          variant="outline"
                          className="text-xs truncate max-w-32"
                          title={employer.occupation}
                        >
                          {employer.occupation}
                        </Badge>
                      ) : (
                        <span className="text-gray-400">-</span>
                      )}
                    </TableCell>
                    <TableCell>
                      {employer.incorporate_status ? (
                        <div className="flex items-center text-sm">
                          <FontAwesomeIcon
                            icon={faBuilding}
                            className="mr-1 text-gray-400 w-3 h-3"
                          />
                          <span
                            className="text-xs truncate max-w-24"
                            title={employer.incorporate_status}
                          >
                            {employer.incorporate_status}
                          </span>
                        </div>
                      ) : (
                        <span className="text-gray-400">-</span>
                      )}
                    </TableCell>
                    <TableCell>
                      {employer.approved_positions ? (
                        <div className="flex items-center">
                          <FontAwesomeIcon
                            icon={faUsers}
                            className="mr-1 text-gray-400 w-3 h-3"
                          />
                          <span className="font-medium">
                            {employer.approved_positions}
                          </span>
                        </div>
                      ) : (
                        <span className="text-gray-400">-</span>
                      )}
                    </TableCell>
                    <TableCell>
                      {employer.approved_lmias ? (
                        <div className="flex items-center">
                          <FontAwesomeIcon
                            icon={faCalendar}
                            className="mr-1 text-gray-400 w-3 h-3"
                          />
                          <span>{employer.approved_lmias}</span>
                        </div>
                      ) : (
                        <span className="text-gray-400">-</span>
                      )}
                    </TableCell>
                    <TableCell>
                      {employer.program_stream ? (
                        <Badge
                          variant="outline"
                          className="text-xs truncate max-w-24"
                          title={employer.program_stream}
                        >
                          {employer.program_stream}
                        </Badge>
                      ) : (
                        <span className="text-gray-400">-</span>
                      )}
                    </TableCell>
                    <TableCell>
                      <span className="text-xs text-gray-500">
                        {employer.year}
                      </span>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </div>

          {/* Pagination */}
          {results.totalPages && results.totalPages > 1 && (
            <Pagination
              currentPage={currentPage}
              totalPages={results.totalPages}
              totalItems={results.totalItems || results.count}
              itemsPerPage={itemsPerPage}
              onPageChange={handlePageChange}
              className="mt-6"
            />
          )}
        </div>
      )}

      {results &&
        results.employers &&
        results.employers.length === 0 &&
        !isLoading &&
        (searchQuery.length >= 2 ||
          city ||
          province ||
          year ||
          (!searchQuery && !city && !province && !year)) && (
          <Card>
            <CardContent className="py-8 text-center">
              <FontAwesomeIcon
                icon={faSearch}
                className="w-8 h-8 text-gray-400 mx-auto mb-2"
              />
              <p className="text-gray-600">
                No employers found matching your search.
              </p>
              <p className="text-sm text-gray-500 mt-1">
                Try different search terms or check the spelling.
              </p>
            </CardContent>
          </Card>
        )}
    </div>
  );
}
