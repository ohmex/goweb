package util

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
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
