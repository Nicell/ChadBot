package helpers

import (
	"net/url"
)

// ValidURL checks whether given string parameter is a valid URL
func ValidURL(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	}
	return true
}
