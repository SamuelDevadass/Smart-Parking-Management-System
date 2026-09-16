package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"api.com/handlers/wings"
)

func main() {
	log.Println("Attempting to load env ...")
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Failed to load env...\nError: %v", err)
	}
	/*Setup the DB Connection Pool*/
	ctx := context.Background()
	db_url := os.Getenv("DB_URL")
	db_pool, err := pgxpool.New(ctx, db_url)
	if err != nil {
		log.Fatalf("DB Connection failed ...\nError: %v", err)
	}
	defer db_pool.Close()

	/*Initialize Network Mux*/
	mux := http.NewServeMux()

	/*Set up CORS Middleware*/
	frontend_url := os.Getenv("FRONTEND_URL")
	handler_with_CORS := enableCORS(mux, frontend_url)
	port := "8000"
	log.Println("Smart Parking Management System API running on ", os.Getenv("BACKEND_URL"))
	if err := http.ListenAndServe(":"+port, handler_with_CORS); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

	/*Initialize the API*/
	config := huma.DefaultConfig("Smart Parking Management System API", "2.0.0")
	api := humago.New(mux, config)
	//log.Println("API: ", api)

	/*Register Handlers*/
	wings.RegisterHandler(api, db_pool)

}
func enableCORS(next http.Handler, allowedOrigin string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
