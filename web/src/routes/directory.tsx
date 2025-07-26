import { createFileRoute } from '@tanstack/react-router'
import { useState, useEffect } from 'react'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faSearch, faMapMarkerAlt, faFilter, faExternalLinkAlt } from '@fortawesome/free-solid-svg-icons'
import { API_BASE_URL } from '@/constants'

interface BusinessRating {
  business_id: string
  current_rating: string
  confidence_score: number
  report_count: number
  avg_tfw_percentage?: number
  last_updated: string
  is_disputed: boolean
}

interface Business {
  id: string
  name: string
  address?: string
  city?: string
  province?: string
  postal_code?: string
  industry_code?: string
  phone?: string
  website?: string
  size_category?: string
  created_at: string
  updated_at: string
  rating?: BusinessRating
}

interface DirectoryResponse {
  businesses: Business[]
  total: number
  limit: number
  offset: number
}

const getRatingColor = (rating: string) => {
  switch (rating) {
    case 'green':
      return 'bg-green-100 text-green-800 border-green-200'
    case 'yellow':
      return 'bg-yellow-100 text-yellow-800 border-yellow-200'
    case 'red':
      return 'bg-red-100 text-red-800 border-red-200'
    default:
      return 'bg-gray-100 text-gray-800 border-gray-200'
  }
}

const getRatingLabel = (rating: string) => {
  switch (rating) {
    case 'green':
      return '🟢 Canadian-First'
    case 'yellow':
      return '🟡 Mixed'
    case 'red':
      return '🔴 TFW-Heavy'
    default:
      return '⚪ Unverified'
  }
}

function DirectoryPage() {
  const [businesses, setBusinesses] = useState<Business[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [total, setTotal] = useState(0)
  
  // Search and filter state
  const [searchQuery, setSearchQuery] = useState('')
  const [cityFilter, setCityFilter] = useState('')
  const [provinceFilter, setProvinceFilter] = useState('')
  const [ratingFilter, setRatingFilter] = useState('')
  const [yearFilter, setYearFilter] = useState('')
  const [currentPage, setCurrentPage] = useState(1)
  const [limit] = useState(20)

  // Generate year options (current year back to 2015)
  const currentYear = new Date().getFullYear()
  const yearOptions = Array.from({length: currentYear - 2014}, (_, i) => currentYear - i)

  const fetchBusinesses = async () => {
    setLoading(true)
    setError(null)
    
    try {
      const params = new URLSearchParams()
      if (searchQuery) params.append('query', searchQuery)
      if (cityFilter) params.append('city', cityFilter)
      if (provinceFilter) params.append('province', provinceFilter)
      if (ratingFilter) params.append('rating', ratingFilter)
      if (yearFilter) params.append('year', yearFilter)
      
      // If no search criteria, default to latest year
      if (!searchQuery && !cityFilter && !provinceFilter && !ratingFilter && !yearFilter) {
        params.append('year', currentYear.toString())
      }
      
      params.append('limit', limit.toString())
      params.append('offset', ((currentPage - 1) * limit).toString())
      
      const response = await fetch(`${API_BASE_URL}/directory?${params}`)
      
      if (!response.ok) {
        throw new Error('Failed to fetch businesses')
      }
      
      const data: DirectoryResponse = await response.json()
      setBusinesses(data.businesses || [])
      setTotal(data.total || 0)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchBusinesses()
  }, [searchQuery, cityFilter, provinceFilter, ratingFilter, yearFilter, currentPage])

  const handleSearch = () => {
    setCurrentPage(1)
    fetchBusinesses()
  }

  const clearFilters = () => {
    setSearchQuery('')
    setCityFilter('')
    setProvinceFilter('')
    setRatingFilter('')
    setYearFilter('')
    setCurrentPage(1)
  }

  const totalPages = Math.ceil(total / limit)

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold mb-2">Business Directory</h1>
        <p className="text-gray-600">
          Community-verified information about business hiring practices in Canada
        </p>
      </div>

      {/* Search and Filters */}
      <div className="bg-white p-6 rounded-lg shadow-sm border mb-6">
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-4 mb-4">
          <div className="relative">
            <FontAwesomeIcon 
              icon={faSearch} 
              className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400" 
            />
            <Input
              placeholder="Search businesses..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="pl-10"
            />
          </div>
          
          <Input
            placeholder="City"
            value={cityFilter}
            onChange={(e) => setCityFilter(e.target.value)}
          />
          
          <Input
            placeholder="Province"
            value={provinceFilter}
            onChange={(e) => setProvinceFilter(e.target.value)}
          />
          
          <select
            value={ratingFilter}
            onChange={(e) => setRatingFilter(e.target.value)}
            className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
          >
            <option value="">All Ratings</option>
            <option value="green">🟢 Canadian-First</option>
            <option value="yellow">🟡 Mixed</option>
            <option value="red">🔴 TFW-Heavy</option>
            <option value="unverified">⚪ Unverified</option>
          </select>
          
          <select
            value={yearFilter}
            onChange={(e) => setYearFilter(e.target.value)}
            className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
          >
            <option value="">All Years</option>
            {yearOptions.map((year) => (
              <option key={year} value={year.toString()}>
                {year}
              </option>
            ))}
          </select>
        </div>
        
        <div className="flex gap-2">
          <Button onClick={handleSearch} size="sm">
            <FontAwesomeIcon icon={faSearch} className="mr-2" />
            Search
          </Button>
          <Button onClick={clearFilters} variant="outline" size="sm">
            <FontAwesomeIcon icon={faFilter} className="mr-2" />
            Clear Filters
          </Button>
        </div>
      </div>

      {/* Results */}
      {loading && (
        <div className="text-center py-8">
          <div className="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
          <p className="mt-2 text-gray-600">Loading businesses...</p>
        </div>
      )}

      {error && (
        <div className="bg-red-50 border border-red-200 rounded-lg p-4 mb-6">
          <p className="text-red-800">Error: {error}</p>
        </div>
      )}

      {!loading && !error && (
        <>
          <div className="mb-4 text-sm text-gray-600">
            Showing {businesses.length} of {total} businesses
          </div>

          <div className="bg-white rounded-lg shadow-sm border mb-8">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Business Name</TableHead>
                  <TableHead>Location</TableHead>
                  <TableHead>Rating</TableHead>
                  <TableHead>Confidence</TableHead>
                  <TableHead>Reports</TableHead>
                  <TableHead>TFW Usage</TableHead>
                  <TableHead>Website</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {businesses.map((business) => (
                  <TableRow key={business.id}>
                    <TableCell className="font-medium">
                      {business.name}
                    </TableCell>
                    <TableCell>
                      {business.address ? (
                        <div className="flex items-center text-sm">
                          <FontAwesomeIcon icon={faMapMarkerAlt} className="mr-1 text-gray-400" />
                          <span>
                            {business.city && `${business.city}`}
                            {business.province && `, ${business.province}`}
                          </span>
                        </div>
                      ) : (
                        <span className="text-gray-400">-</span>
                      )}
                    </TableCell>
                    <TableCell>
                      {business.rating ? (
                        <Badge
                          className={getRatingColor(business.rating.current_rating)}
                          variant="outline"
                        >
                          {getRatingLabel(business.rating.current_rating)}
                        </Badge>
                      ) : (
                        <Badge variant="outline" className="bg-gray-100 text-gray-800 border-gray-200">
                          ⚪ Unverified
                        </Badge>
                      )}
                    </TableCell>
                    <TableCell>
                      {business.rating ? (
                        <span className="font-medium">{business.rating.confidence_score.toFixed(1)}</span>
                      ) : (
                        <span className="text-gray-400">-</span>
                      )}
                    </TableCell>
                    <TableCell>
                      {business.rating ? (
                        business.rating.report_count
                      ) : (
                        <span className="text-gray-400">0</span>
                      )}
                    </TableCell>
                    <TableCell>
                      {business.rating && business.rating.avg_tfw_percentage !== null ? (
                        <span>{business.rating.avg_tfw_percentage.toFixed(1)}%</span>
                      ) : (
                        <span className="text-gray-400">-</span>
                      )}
                    </TableCell>
                    <TableCell>
                      {business.website ? (
                        <a
                          href={business.website}
                          target="_blank"
                          rel="noopener noreferrer"
                          className="text-blue-600 hover:text-blue-800 inline-flex items-center"
                        >
                          <FontAwesomeIcon icon={faExternalLinkAlt} className="w-3 h-3" />
                        </a>
                      ) : (
                        <span className="text-gray-400">-</span>
                      )}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </div>

          {/* Pagination */}
          {totalPages > 1 && (
            <div className="flex justify-center items-center space-x-2">
              <Button
                onClick={() => setCurrentPage(Math.max(1, currentPage - 1))}
                disabled={currentPage === 1}
                variant="outline"
                size="sm"
              >
                Previous
              </Button>
              
              <span className="text-sm text-gray-600">
                Page {currentPage} of {totalPages}
              </span>
              
              <Button
                onClick={() => setCurrentPage(Math.min(totalPages, currentPage + 1))}
                disabled={currentPage === totalPages}
                variant="outline"
                size="sm"
              >
                Next
              </Button>
            </div>
          )}
        </>
      )}

      {!loading && !error && businesses.length === 0 && (
        <div className="text-center py-12">
          <p className="text-gray-500 text-lg">No businesses found matching your criteria.</p>
          <Button onClick={clearFilters} className="mt-4" variant="outline">
            Clear filters to see all businesses
          </Button>
        </div>
      )}
    </div>
  )
}

export const Route = createFileRoute('/directory')({
  component: DirectoryPage,
})