package repos

import (
	"canada-hires/models"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type AddressGeocodingCacheRepository interface {
	GetByNormalizedAddress(normalizedAddress string) (*models.AddressGeocodingCache, error)
	Upsert(cache *models.AddressGeocodingCache) error
	DeleteOlderThan(days int) error
}

type addressGeocodingCacheRepository struct {
	db *sqlx.DB
}

func NewAddressGeocodingCacheRepository(db *sqlx.DB) AddressGeocodingCacheRepository {
	return &addressGeocodingCacheRepository{db: db}
}

// NormalizeAddress normalizes an address string for consistent caching
func NormalizeAddress(address string) string {
	// Convert to lowercase and trim whitespace
	normalized := strings.ToLower(strings.TrimSpace(address))
	
	// Remove extra whitespace between words
	words := strings.Fields(normalized)
	normalized = strings.Join(words, " ")
	
	// Remove common punctuation
	replacements := []string{
		",", "",
		".", "",
		"  ", " ", // Replace double spaces with single space
	}
	
	replacer := strings.NewReplacer(replacements...)
	normalized = replacer.Replace(normalized)
	
	return strings.TrimSpace(normalized)
}

func (r *addressGeocodingCacheRepository) GetByNormalizedAddress(normalizedAddress string) (*models.AddressGeocodingCache, error) {
	var cache models.AddressGeocodingCache
	query := `
		SELECT id, address, normalized_address, latitude, longitude, confidence, geocoded_at, created_at, updated_at
		FROM address_geocoding_cache
		WHERE normalized_address = $1`

	err := r.db.Get(&cache, query, normalizedAddress)
	if err != nil {
		return nil, err
	}

	return &cache, nil
}

func (r *addressGeocodingCacheRepository) Upsert(cache *models.AddressGeocodingCache) error {
	query := `
		INSERT INTO address_geocoding_cache (address, normalized_address, latitude, longitude, confidence, geocoded_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (normalized_address) DO UPDATE SET
			address = EXCLUDED.address,
			latitude = EXCLUDED.latitude,
			longitude = EXCLUDED.longitude,
			confidence = EXCLUDED.confidence,
			geocoded_at = EXCLUDED.geocoded_at,
			updated_at = CURRENT_TIMESTAMP
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(query,
		cache.Address,
		cache.NormalizedAddress,
		cache.Latitude,
		cache.Longitude,
		cache.Confidence,
		cache.GeocodedAt,
	).Scan(&cache.ID, &cache.CreatedAt, &cache.UpdatedAt)
}

func (r *addressGeocodingCacheRepository) DeleteOlderThan(days int) error {
	query := `DELETE FROM address_geocoding_cache WHERE geocoded_at < $1`
	cutoffDate := time.Now().AddDate(0, 0, -days)
	
	_, err := r.db.Exec(query, cutoffDate)
	return err
}