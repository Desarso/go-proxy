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
// 		// Extract the 'target' from the query parameter
// 		target := r.URL.Query().Get("target")
// 		log.Printf("Incoming request path: %s, target: %s", r.URL.Path, target)

// 		if target == "" {
// 			// If 'target' is not provided in the query, check the Referer header
// 			referer := r.Referer()
// 			if referer != "" {
// 				// Parse the referer URL
// 				refererURL, err := url.Parse(referer)
// 				if err == nil {
// 					// Extract the 'target' from the referer query parameter if it exists
// 					target = refererURL.Query().Get("target")
// 					log.Printf("Referer detected, new target: %s", target)
// 				}
// 			}
// 		}

// 		if target == "" {
// 			// If still no target, return an error
// 			http.Error(w, "Target parameter is missing or Referer is invalid", http.StatusBadRequest)
// 			return
// 		}

// 		// Parse the target URL
// 		targetURL, err := url.Parse(target)
// 		if err != nil || targetURL.Scheme == "" || targetURL.Host == "" {
// 			http.Error(w, "Invalid target URL", http.StatusBadRequest)
// 			return
// 		}

// 		// Combine the target URL with the requested path and query string
// 		proxyURL := targetURL.ResolveReference(r.URL)

// 		log.Printf("Proxying request to: %s", proxyURL.String())

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
// 	err := http.ListenAndServe(":8080", nil)
// 	if err != nil {
// 		log.Fatalf("Failed to start server: %v", err)
// 	}
// }

package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

func main() {
	// Define the handler function for the proxy
	handler := func(w http.ResponseWriter, r *http.Request) {
		// Extract the 'target' from the query parameter
		target := r.URL.Query().Get("target")
		log.Printf("Incoming request path: %s, target: %s", r.URL.Path, target)

		//if target is set it means is it the first request and we must immediately retun the html page
		if target != "" {
			// Parse the target URL
			targetURL, err := url.Parse(target)
			fmt.Println("targetURL: ", targetURL)
			if err != nil || targetURL.Scheme == "" || targetURL.Host == "" {
				http.Error(w, "Invalid target URL", http.StatusBadRequest)
				return
			}
			proxyURL := targetURL
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
		} else {
			if target == "" {
				// If 'target' is not provided in the query, check the Referer header
				referer := r.Referer()
				if referer != "" {
					// Parse the referer URL
					refererURL, err := url.Parse(referer)
					if err == nil {
						// Extract the 'target' from the referer query parameter if it exists
						target = refererURL.Query().Get("target")
						log.Printf("Referer detected, new target: %s", target)
					}
				}
			}

			if target == "" {
				// If still no target, return an error
				http.Error(w, "Target parameter is missing or Referer is invalid", http.StatusBadRequest)
				return
			}

			fmt.Println("targethere: ", target)

			// Parse the target URL
			targetURL, err := url.Parse(target)
			if err != nil || targetURL.Scheme == "" || targetURL.Host == "" {
				http.Error(w, "Invalid target URL", http.StatusBadRequest)
				return
			}

			// Combine the target URL with the requested path and query string only
			//add r.URL.Path and r.URL.RawQuery to the targetURL
			//newURL := targetURL.String() + r.URL.Path + "?" + r.URL.RawQuery
			proxyURL := targetURL.ResolveReference(r.URL)

			// proxyURL, err := url.Parse(newURL)
			// if err != nil {
			// 	http.Error(w, "Invalid target URL", http.StatusBadRequest)
			// 	return
			// }
			// fmt.Println("proxyURL: ", proxyURL)
			// fmt.Println("targetURL: ", targetURL)
			// fmt.Println("newURL: ", newURL)

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

	}

	// Start the HTTP server with the proxy handler
	http.HandleFunc("/", handler)
	log.Println("Starting proxy server on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
