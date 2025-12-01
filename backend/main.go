package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"shortfyurl/internal/cache"
	"shortfyurl/internal/database"
	"shortfyurl/internal/handlers"

	"github.com/gorilla/mux"
)

func main() {
	cassandraHosts := os.Getenv("CASSANDRA_HOSTS")
	if cassandraHosts == "" {
		cassandraHosts = "127.0.0.1"
	}
	hosts := strings.Split(cassandraHosts, ",")

	db, err := database.NewCassandraDB(hosts, "shortfy")
	if err != nil {
		log.Fatalf("Erro ao conectar no Cassandra: %v", err)
	}
	defer db.Close()

	if err := db.InitSchema(); err != nil {
		log.Fatalf("Erro ao inicializar schema: %v", err)
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	redisCache := cache.NewRedisCache(redisAddr)
	defer redisCache.Close()

	handler := handlers.NewHandler(db, redisCache)

	r := mux.NewRouter()
	
	// API endpoints
	r.HandleFunc("/api/shorten", handler.CreateShortURL).Methods("POST")
	r.HandleFunc("/api/stats/{shortCode}", handler.GetStats).Methods("GET")
	r.HandleFunc("/api/urls", handler.ListAllURLs).Methods("GET")
	
	// Redirect
	r.HandleFunc("/{shortCode}", handler.RedirectURL).Methods("GET")
	
	// Frontend
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./frontend")))

	log.Println("Servidor rodando na porta 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
