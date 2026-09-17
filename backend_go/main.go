package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	"api.com/handlers/wings"
	"api.com/repositories"
)

func main() {

	//Load environment variables
	log.Println("Attempting to load env ...")
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Failed to load env...\nError: %v", err)
	}
	backend_url := os.Getenv("BACKEND_URL")
	connection_string := os.Getenv("DB_URL")

	//Initialize connection to DB
	log.Println("Initializing DB ...")
	repositories.Init_DB(context.Background(), connection_string)

	// 1. CREATE CHI ROUTER
	r := chi.NewRouter()

	// 2. SET UP MIDDLEWARES
	//CORS
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*") // frontend URL
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			// Catch preflight browser tests instantly
			if req.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, req)
		})
	})
	// Standard chi utilities (optional but highly recommended)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// 3. INITIALIZE THE API
	config := huma.DefaultConfig("Smart Parking Management System API", "2.0.0")
	api := humachi.New(r, config)

	// 4. REGISTER HANDLERS
	wings.RegisterHandler(api)

	// 5. LISTEN AND SERVE
	log.Printf("Server running on http://%s", backend_url)
	log.Printf("OpenAPI UI available at http://%s/docs", backend_url)
	err = http.ListenAndServe(backend_url, r)
	if err != nil {
		log.Fatalf("Server crashed: %v", err)
	}
}
