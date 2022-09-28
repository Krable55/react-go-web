package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

const portNumber = ":8080"

// * Build environemnts
const (
	production  = "prod"
	development = "dev"
)

func init() {
	env := flag.String("env", development, "dev | prod")
	flag.Parse()
	os.Setenv("ENV", *env)
}
func main() {

	r := mux.NewRouter()
	env := os.Getenv("ENV")
	fmt.Printf("Environment: %s\n", env)

	// Handle API routes
	api := r.PathPrefix("/api/").Subrouter()
	api.HandleFunc("/example", func(w http.ResponseWriter, r *http.Request) {
		// * http://localhost:8080/api/example -> "From the API"
		fmt.Fprintln(w, "From the API")
	})

	isProd := env == production
	if isProd {
		// Serve static files if in production
		r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./build/static/"))))

		// Serve index page on all unhandled routes
		r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "./build/index.html")
		})
		fmt.Printf("Serving at: http://localhost%s", portNumber)
	} else {
		// Only serve the API if in development
		fmt.Printf("Serving api at: http://localhost%s/api", portNumber)
	}

	log.Fatal(http.ListenAndServe(portNumber, r))

	// * Alternative implementation with gin router:
	/*
		 router := gin.Default()

		// Serve frontend static files
		router.Use(static.Serve("/", static.LocalFile("./build", true)))

		// Setup route group for the API
		api := router.Group("/api")
		{
			api.GET("/", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "pong",
				})
			})
		}

		// Start and run the server
		fmt.Printf("Staring application on port %s", portNumber)
		router.Run(portNumber)
	*/

}
