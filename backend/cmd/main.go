package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

var db *sql.DB

type URLMapping struct {
	ID          int       `json:"id"`
	OriginalURL string    `json:"original_url"`
	ShortURL    string    `json:"short_url"`
	CreatedAt   time.Time `json:"created_at"`
}

func generateShortURL(n int) (string, error) {
	alphabet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	for {
		b := make([]byte, n)
		for i := range b {
			b[i] = alphabet[time.Now().UnixNano()%int64(len(alphabet))]
		}

		short := string(b)

		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM url_mapping WHERE short_url=$1)", short).Scan(&exists)
		if err != nil {
			return "", err
		}

		if !exists {
			return short, nil
		}
	}
}

func saveURLToDatabaseHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var reqData URLMapping
	err := json.NewDecoder(r.Body).Decode(&reqData)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	short, err := generateShortURL(8)
	if err != nil {
		http.Error(w, "Error generating short URL", http.StatusInternalServerError)
		return
	}

	err = db.QueryRow(
		"INSERT INTO url_mapping (original_url, short_url) VALUES ($1, $2) RETURNING id, created_at",
		reqData.OriginalURL, short,
	).Scan(&reqData.ID, &reqData.CreatedAt)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	baseurl := os.Getenv("BASE_URL")
	reqData.ShortURL = fmt.Sprintf("https://%s/%s", baseurl, short)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reqData)
}

func redirectToOriginalURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	short := r.URL.Path[1:]
	if short == "" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	var originalURL string
	err := db.QueryRow("SELECT original_url FROM url_mapping WHERE short_url=$1", short).Scan(&originalURL)
	if err == sql.ErrNoRows {
		http.Error(w, "Not found", http.StatusNotFound)
		return

	} else if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusFound)
}

func main() {
	if err := godotenv.Load(); err != nil {
        fmt.Println("No .env file found")
    }

	var err error
	connStr := os.Getenv("DATABASE_PUBLIC_URL")
	db, err = sql.Open("postgres", connStr)
	if err != nil {
	    log.Fatal("Database connection error:", err)
	}
	defer db.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/healthy", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("/api/shorten", saveURLToDatabaseHandler)
	mux.HandleFunc("/", redirectToOriginalURLHandler)

	fmt.Println("Server is running on port", port)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, 
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization", "X-Requested-With"},
		AllowCredentials: false,
		MaxAge: 300, 
	})

	handler := c.Handler(mux)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
