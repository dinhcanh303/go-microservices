package utils

import "net/http"

// GetUserAgent helps in getting the user agent from the request
func GetUserAgent(r *http.Request) string {
	return r.UserAgent()
}

// GetIP helps in getting the IP address from the request
func GetIP(r *http.Request) string {
	ipAddress := r.Header.Get("X-Read-Ip")
	if ipAddress == "" {
		ipAddress = r.Header.Get("X-Forwarded-For")
	}
	if ipAddress == "" {
		ipAddress = r.RemoteAddr
	}
	return ipAddress
}
