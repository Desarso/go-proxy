package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
)

// handleRequestAndRedirect will take an incoming request and forward it
func handleRequestAndRedirect(res http.ResponseWriter, req *http.Request) {
	// Parse the target URL
	targetURL := req.URL.Query().Get("url")
	if targetURL == "" {
		http.Error(res, "Missing 'url' query parameter", http.StatusBadRequest)
		return
	}

	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		http.Error(res, "Invalid URL provided", http.StatusBadRequest)
		return
	}

	// Create a new request to the target URL
	proxyReq, err := http.NewRequest(req.Method, parsedURL.String(), req.Body)
	if err != nil {
		http.Error(res, "Failed to create request", http.StatusInternalServerError)
		return
	}

	// Copy headers from the original request
	proxyReq.Header = req.Header

	// Perform the request
	client := &http.Client{}
	resp, err := client.Do(proxyReq)
	if err != nil {
		http.Error(res, "Failed to perform request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Copy the response headers and status code
	for key, values := range resp.Header {
		for _, value := range values {
			res.Header().Add(key, value)
		}
	}
	res.WriteHeader(resp.StatusCode)

	// Copy the response body
	io.Copy(res, resp.Body)
}

func main() {
	// Handle all requests with the handleRequestAndRedirect function
	http.HandleFunc("/", handleRequestAndRedirect)

	// Start the server on port 8080
	log.Println("Proxy server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
