package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"shortfyurl/internal/cache"
	"shortfyurl/internal/database"
	"shortfyurl/internal/models"
	"shortfyurl/internal/utils"
	"time"

	"github.com/gorilla/mux"
)

type Handler struct {
	db    *database.CassandraDB
	cache *cache.RedisCache
}

func NewHandler(db *database.CassandraDB, cache *cache.RedisCache) *Handler {
	return &Handler{db: db, cache: cache}
}

func (h *Handler) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	var req models.CreateURLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Requisição inválida", http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		http.Error(w, "URL é obrigatória", http.StatusBadRequest)
		return
	}

	id := time.Now().UnixNano()
	shortCode := utils.EncodeBase62(id)

	url := &models.URL{
		ID:          id,
		ShortCode:   shortCode,
		OriginalURL: req.URL,
		CreatedAt:   time.Now(),
		Clicks:      0,
	}

	if err := h.db.SaveURL(url); err != nil {
		log.Printf("Erro ao salvar URL: %v", err)
		http.Error(w, "Erro ao criar URL curta", http.StatusInternalServerError)
		return
	}

	h.cache.Set(shortCode, url, 24*time.Hour)

	response := models.CreateURLResponse{
		ShortCode: shortCode,
		ShortURL:  "http://localhost:8080/" + shortCode,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) RedirectURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortCode := vars["shortCode"]

	url, err := h.cache.Get(shortCode)
	if err != nil {
		url, err = h.db.GetURLByShortCode(shortCode)
		if err != nil {
			http.Error(w, "URL não encontrada", http.StatusNotFound)
			return
		}
		h.cache.Set(shortCode, url, 24*time.Hour)
	}

	go h.db.IncrementClicks(shortCode)

	http.Redirect(w, r, url.OriginalURL, http.StatusMovedPermanently)
}

func (h *Handler) GetStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortCode := vars["shortCode"]

	url, err := h.db.GetURLByShortCode(shortCode)
	if err != nil {
		http.Error(w, "URL não encontrada", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(url)
}

func (h *Handler) ListAllURLs(w http.ResponseWriter, r *http.Request) {
	urls, err := h.db.GetAllURLs()
	if err != nil {
		log.Printf("Erro ao buscar URLs: %v", err)
		http.Error(w, "Erro ao buscar URLs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(urls)
}
