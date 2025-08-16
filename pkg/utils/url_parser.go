package utils

import "net/url"

// GetDomainFromURL mengekstrak hostname (domain) dari sebuah string URL.
func GetDomainFromURL(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	return parsedURL.Hostname(), nil
}
