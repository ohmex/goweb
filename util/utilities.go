package util

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// function to check given string is in array or not
func Contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

func GenerateAPIKeySecret() (string, string, string) {
	key := make([]byte, 16)
	secret := make([]byte, 32)
	rand.Read(key)
	rand.Read(secret)

	apiKey := hex.EncodeToString(key)
	apiSecret := hex.EncodeToString(secret)
	hashedSecret, _ := bcrypt.GenerateFromPassword([]byte(apiSecret), bcrypt.DefaultCost)

	return apiKey, apiSecret, string(hashedSecret)
}

// Generic helper to bind and validate request
func BindAndValidate[T any](c echo.Context) (*T, error) {
	obj := new(T)
	if err := c.Bind(obj); err != nil {
		return nil, err
	}
	if validator, ok := any(obj).(interface{ Validate() error }); ok {
		if err := validator.Validate(); err != nil {
			return nil, err
		}
	}
	return obj, nil
}

// Extracts domain from echo.Context
func ExtractDomain(e echo.Context) (interface{}, error) {
	domain := e.Get("domain")
	if domain == nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Invalid domain context")
	}
	return domain, nil
}

// Checks if user domains are loaded
func DomainsLoaded(user interface{ GetDomains() []interface{} }) bool {
	domains := user.GetDomains()
	return len(domains) > 0 && domains[0] != nil
}

// Skips Gzip middleware for health and swagger endpoints
func GzipSkipper(c echo.Context) bool {
	p := c.Path()
	return p == "/health" || len(p) >= 8 && p[:8] == "/swagger"
}

// HashPassword hashes a plain password using bcrypt
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

// CheckPasswordHash compares a bcrypt hashed password with its possible plaintext equivalent
func CheckPasswordHash(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// GetUUIDParam extracts the 'uuid' parameter from the Echo context and returns an error if missing
func GetUUIDParam(e echo.Context) (string, error) {
	uuid := e.Param("uuid")
	if uuid == "" {
		return "", echo.NewHTTPError(http.StatusBadRequest, "Missing or invalid UUID parameter")
	}
	return uuid, nil
}

// GeneratePartitionName generates a safe partition name for a domain UUID
func GeneratePartitionName(domainUUID string) string {
	// Replace hyphens with underscores to make it PostgreSQL compatible
	safeUUID := strings.ReplaceAll(domainUUID, "-", "_")
	return "posts_" + safeUUID
}

// IsPartitioningEnabled checks if database partitioning is enabled via environment variable
func IsPartitioningEnabled() bool {
	return os.Getenv("DB_PARTITIONING_ENABLED") == "true"
}

// IsDatabasePartitioningSupported checks if the current database supports partitioning
func IsDatabasePartitioningSupported(db *gorm.DB) bool {
	return isPostgreSQLCompatible(db)
}

// isPostgreSQLCompatible checks if the database is PostgreSQL or YugabyteDB
// Both support the same partitioning syntax
func isPostgreSQLCompatible(db *gorm.DB) bool {
	// Check both the dialector name and the configured driver name
	dialectorName := db.Dialector.Name()
	configuredDriver := os.Getenv("DB_DRIVER")
	
	// PostgreSQL dialector or YugabyteDB configured driver both support partitioning
	return dialectorName == "postgres" || configuredDriver == "yugabytedb"
}

// DatabasePoolStats represents connection pool statistics
type DatabasePoolStats struct {
	MaxOpenConnections int
	OpenConnections    int
	InUse              int
	Idle               int
	WaitCount          int64
	WaitDuration       time.Duration
	MaxIdleClosed      int64
	MaxLifetimeClosed  int64
}

// GetDatabasePoolStats retrieves current database connection pool statistics
func GetDatabasePoolStats(db *gorm.DB) (*DatabasePoolStats, error) {
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	return &DatabasePoolStats{
		MaxOpenConnections: sqlDB.Stats().MaxOpenConnections,
		OpenConnections:    sqlDB.Stats().OpenConnections,
		InUse:              sqlDB.Stats().InUse,
		Idle:               sqlDB.Stats().Idle,
		WaitCount:          sqlDB.Stats().WaitCount,
		WaitDuration:       sqlDB.Stats().WaitDuration,
		MaxIdleClosed:      sqlDB.Stats().MaxIdleClosed,
		MaxLifetimeClosed:  sqlDB.Stats().MaxLifetimeClosed,
	}, nil
}

// LogDatabasePoolStats logs database connection pool statistics
func LogDatabasePoolStats(db *gorm.DB) {
	stats, err := GetDatabasePoolStats(db)
	if err != nil {
		return
	}

	log.Info().
		Int("max_open", stats.MaxOpenConnections).
		Int("open", stats.OpenConnections).
		Int("in_use", stats.InUse).
		Int("idle", stats.Idle).
		Int64("wait_count", stats.WaitCount).
		Dur("wait_duration", stats.WaitDuration).
		Int64("max_idle_closed", stats.MaxIdleClosed).
		Int64("max_lifetime_closed", stats.MaxLifetimeClosed).
		Msg("Database connection pool stats")
}
