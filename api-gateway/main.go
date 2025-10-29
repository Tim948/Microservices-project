package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	userServiceURL := os.Getenv("USER_SERVICE_URL")
	taskServiceURL := os.Getenv("TASK_SERVICE_URL")

	log.Printf("Configuring API Gateway with:")
	log.Printf("User Service URL: %s", userServiceURL)
	log.Printf("Task Service URL: %s", taskServiceURL)

	// Health check
	r.GET("/health", func(c *gin.Context) {
		services := map[string]string{
			"user_service": userServiceURL,
			"task_service": taskServiceURL,
		}

		status := "OK"
		for name, url := range services {
			if !isServiceHealthy(url) {
				status = "DEGRADED"
				log.Printf("Service %s is not healthy", name)
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"status":    status,
			"services":  services,
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	// User service routes - FIXED: don't trim base path
	r.Any("/users/*path", createProxyHandler(userServiceURL, ""))
	r.Any("/users", createProxyHandler(userServiceURL, ""))

	// Task service routes - FIXED: don't trim base path
	r.Any("/tasks/*path", createProxyHandler(taskServiceURL, ""))
	r.Any("/tasks", createProxyHandler(taskServiceURL, ""))

	// Default route
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "API Gateway is running",
			"endpoints": []string{
				"GET /health",
				"GET /users",
				"GET /tasks",
				"POST /users",
				"POST /tasks",
			},
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("API Gateway running on port %s", port)
	log.Fatal(r.Run(":" + port))
}

func isServiceHealthy(serviceURL string) bool {
	if serviceURL == "" {
		return false
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(serviceURL + "/health")
	if err != nil {
		log.Printf("Health check failed for %s: %v", serviceURL, err)
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

func createProxyHandler(targetURL, basePath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// If targetURL is not set, return error
		if targetURL == "" {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error": "Service is not configured",
			})
			return
		}

		target, err := url.Parse(targetURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Invalid service URL: " + err.Error(),
			})
			return
		}

		// Get the full path from the request
		path := c.Request.URL.Path

		// Build the target URL - append the full path to the target service URL
		proxyURL := target.ResolveReference(&url.URL{Path: path})

		log.Printf("Proxying request: %s %s -> %s", c.Request.Method, c.Request.URL.Path, proxyURL.String())

		// Create new request
		req, err := http.NewRequest(c.Request.Method, proxyURL.String(), c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Copy headers
		for key, values := range c.Request.Header {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}

		// Add original host header if not present
		if req.Header.Get("X-Forwarded-Host") == "" {
			req.Header.Set("X-Forwarded-Host", c.Request.Host)
		}

		// Execute request
		client := &http.Client{Timeout: 30 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{
				"error":       "Cannot connect to service: " + err.Error(),
				"service_url": proxyURL.String(),
			})
			return
		}
		defer resp.Body.Close()

		// Copy response headers
		for key, values := range resp.Header {
			for _, value := range values {
				c.Header(key, value)
			}
		}

		// Copy status code and body
		c.Status(resp.StatusCode)
		io.Copy(c.Writer, resp.Body)
	}
}
