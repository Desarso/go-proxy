// package main

// import (
// 	"io"
// 	"log"
// 	"net/http"
// 	"net/url"
// )

// // handleRequestAndRedirect will take an incoming request and forward it
// func handleRequestAndRedirect(res http.ResponseWriter, req *http.Request) {
// 	// Parse the target URL
// 	targetURL := req.URL.Query().Get("url")
// 	if targetURL == "" {
// 		http.Error(res, "Missing 'url' query parameter", http.StatusBadRequest)
// 		return
// 	}

// 	parsedURL, err := url.Parse(targetURL)
// 	if err != nil {
// 		http.Error(res, "Invalid URL provided", http.StatusBadRequest)
// 		return
// 	}

// 	// Create a new request to the target URL
// 	proxyReq, err := http.NewRequest(req.Method, parsedURL.String(), req.Body)
// 	if err != nil {
// 		http.Error(res, "Failed to create request", http.StatusInternalServerError)
// 		return
// 	}

// 	// Copy headers from the original request
// 	proxyReq.Header = req.Header

// 	// Perform the request
// 	client := &http.Client{}
// 	resp, err := client.Do(proxyReq)
// 	if err != nil {
// 		http.Error(res, "Failed to perform request", http.StatusInternalServerError)
// 		return
// 	}
// 	defer resp.Body.Close()

// 	// Copy the response headers and status code
// 	for key, values := range resp.Header {
// 		for _, value := range values {
// 			res.Header().Add(key, value)
// 		}
// 	}
// 	res.WriteHeader(resp.StatusCode)

// 	// Copy the response body
// 	io.Copy(res, resp.Body)
// }

// func main() {
// 	// Handle all requests with the handleRequestAndRedirect function
// 	http.HandleFunc("/", handleRequestAndRedirect)

// 	// Start the server on port 8080
// 	log.Println("Proxy server is running on port 8080")
// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }

// package main

// import (
// 	"io"
// 	"log"
// 	"net/http"
// 	"net/url"
// )

// func main() {
// 	// Define the target base URL that you want to proxy
// 	target := "https://example.com" // Change this to your target URL

// 	// Parse the target URL
// 	targetURL, err := url.Parse(target)
// 	if err != nil {
// 		log.Fatalf("Failed to parse target URL: %v", err)
// 	}

// 	// Define the handler function for the proxy
// 	handler := func(w http.ResponseWriter, r *http.Request) {
// 		// Log the incoming request URL
// 		log.Printf("Handling request for: %s", r.URL.String())

// 		// Create a new URL using the target base URL and the request URI
// 		proxyURL := targetURL.ResolveReference(r.URL)

// 		// Create a new HTTP request based on the incoming request
// 		req, err := http.NewRequest(r.Method, proxyURL.String(), r.Body)
// 		if err != nil {
// 			http.Error(w, "Failed to create request", http.StatusInternalServerError)
// 			return
// 		}

// 		// Copy headers from the original request to the proxied request
// 		for name, values := range r.Header {
// 			for _, value := range values {
// 				req.Header.Add(name, value)
// 			}
// 		}

// 		// Make the HTTP request to the target server
// 		resp, err := http.DefaultClient.Do(req)
// 		if err != nil {
// 			http.Error(w, "Failed to reach target server", http.StatusBadGateway)
// 			return
// 		}
// 		defer resp.Body.Close()

// 		// Copy headers from the target server's response
// 		for name, values := range resp.Header {
// 			for _, value := range values {
// 				w.Header().Add(name, value)
// 			}
// 		}

// 		// Set the status code from the target server's response
// 		w.WriteHeader(resp.StatusCode)

// 		// Copy the response body from the target server to the client
// 		_, err = io.Copy(w, resp.Body)
// 		if err != nil {
// 			log.Printf("Error copying response body: %v", err)
// 		}
// 	}

// 	// Start the HTTP server with the proxy handler
// 	http.HandleFunc("/", handler)
// 	log.Println("Starting proxy server on port 8080...")
// 	err = http.ListenAndServe(":8080", nil)
// 	if err != nil {
// 		log.Fatalf("Failed to start server: %v", err)
// 	}
// }

// package main

// import (
// 	"io"
// 	"log"
// 	"net/http"
// 	"net/url"
// )

// func main() {
// 	// Define the handler function for the proxy
// 	handler := func(w http.ResponseWriter, r *http.Request) {
// 		// Extract the 'target' query parameter from the incoming request
// 		target := r.URL.Query().Get("target")
// 		if target == "" {
// 			http.Error(w, "Target parameter is missing", http.StatusBadRequest)
// 			return
// 		}

// 		// Parse the target URL
// 		targetURL, err := url.Parse(target)
// 		if err != nil || targetURL.Scheme == "" || targetURL.Host == "" {
// 			http.Error(w, "Invalid target URL", http.StatusBadRequest)
// 			return
// 		}

// 		// Resolve the full URL based on the target and the original request
// 		proxyURL := targetURL.ResolveReference(r.URL)

// 		// Create a new HTTP request based on the incoming request
// 		req, err := http.NewRequest(r.Method, proxyURL.String(), r.Body)
// 		if err != nil {
// 			http.Error(w, "Failed to create request", http.StatusInternalServerError)
// 			return
// 		}

// 		// Copy headers from the original request to the proxied request
// 		for name, values := range r.Header {
// 			for _, value := range values {
// 				req.Header.Add(name, value)
// 			}
// 		}

// 		// Make the HTTP request to the target server
// 		resp, err := http.DefaultClient.Do(req)
// 		if err != nil {
// 			http.Error(w, "Failed to reach target server", http.StatusBadGateway)
// 			return
// 		}
// 		defer resp.Body.Close()

// 		// Copy headers from the target server's response
// 		for name, values := range resp.Header {
// 			for _, value := range values {
// 				w.Header().Add(name, value)
// 			}
// 		}

// 		// Set the status code from the target server's response
// 		w.WriteHeader(resp.StatusCode)

// 		// Copy the response body from the target server to the client
// 		_, err = io.Copy(w, resp.Body)
// 		if err != nil {
// 			log.Printf("Error copying response body: %v", err)
// 		}
// 	}

//		// Start the HTTP server with the proxy handler
//		http.HandleFunc("/", handler)
//		log.Println("Starting proxy server on port 8080...")
//		err := http.ListenAndServe(":8080", nil)
//		if err != nil {
//			log.Fatalf("Failed to start server: %v", err)
//		}
//	}
package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
)

func main() {
	// Define the handler function for the proxy
	handler := func(w http.ResponseWriter, r *http.Request) {
		// First, try to extract the 'target' from the query parameter
		target := r.URL.Query().Get("target")

		if target == "" {
			// If 'target' is not provided in the query, check the Referer header
			referer := r.Referer()
			if referer != "" {
				// Parse the referer URL
				refererURL, err := url.Parse(referer)
				if err == nil {
					// Extract the 'target' from the referer query parameter if it exists
					target = refererURL.Query().Get("target")
				}
			}
		}

		if target == "" {
			// If still no target, return an error
			http.Error(w, "Target parameter is missing or Referer is invalid", http.StatusBadRequest)
			return
		}

		// Parse the target URL
		targetURL, err := url.Parse(target)
		if err != nil || targetURL.Scheme == "" || targetURL.Host == "" {
			http.Error(w, "Invalid target URL", http.StatusBadRequest)
			return
		}

		// Resolve the full URL based on the target and the original request path
		// This preserves the path for assets and other resources
		proxyURL := targetURL.ResolveReference(r.URL)

		// Create a new HTTP request based on the incoming request
		req, err := http.NewRequest(r.Method, proxyURL.String(), r.Body)
		if err != nil {
			http.Error(w, "Failed to create request", http.StatusInternalServerError)
			return
		}

		// Copy headers from the original request to the proxied request
		for name, values := range r.Header {
			for _, value := range values {
				req.Header.Add(name, value)
			}
		}

		// Make the HTTP request to the target server
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			http.Error(w, "Failed to reach target server", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		// Copy headers from the target server's response
		for name, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(name, value)
			}
		}

		// Set the status code from the target server's response
		w.WriteHeader(resp.StatusCode)

		// Copy the response body from the target server to the client
		_, err = io.Copy(w, resp.Body)
		if err != nil {
			log.Printf("Error copying response body: %v", err)
		}
	}

	// Start the HTTP server with the proxy handler
	http.HandleFunc("/", handler)
	log.Println("Starting proxy server on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
