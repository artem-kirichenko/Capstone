package main

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"net/http"
	"time"
)

func SessionMiddleware(ctx context.Context, db *sql.DB) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract session token from the header
			token := r.Header.Get("session-token")
			if token == "" {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Missing session token"))
				return
			}

			// Check token in the database
			var sessionID string
			query := "SELECT id FROM sessions WHERE token = ? AND expiration > NOW()"
			err := db.QueryRowContext(ctx, query, token).Scan(&sessionID)
			if err != nil {
				if err == sql.ErrNoRows {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("Invalid or expired session token"))
				} else {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
				return
			}

			// Add session ID to the context
			ctx := context.WithValue(r.Context(), "sessionID", sessionID)
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// generate token for session
func generateSessionToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// inserts a session into the database
func CreateSession(db *sql.DB, token string, duration time.Duration) (string, error) {
	// Generate a unique session ID
	id := uuid.NewString()
	expiration := time.Now().Add(duration)

	// Insert session into the database
	query := `INSERT INTO sessions (user_id, token, expiration) VALUES (?, ?, ?)`
	_, err := db.Exec(query, id, token, expiration)
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}

	return id, err
}

// checking password
func checkPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func login_page(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/login.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	tmpl.ExecuteTemplate(w, "login", nil)
	if err != nil {
		fmt.Fprintf(w, "Error executing template: %s", err.Error())
	}
}

func login_authorization(ctx context.Context, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
		r.ParseForm()
		email := r.Form.Get("email")
		pass := r.Form.Get("pass")

		row := db.QueryRow("SELECT `pass` FROM `users` WHERE email = ?", email)
		var hashedPassword string
		err := row.Scan(&hashedPassword)
		if err != nil {
			if err == sql.ErrNoRows {
				fmt.Println("User not found")
				return
			}
			panic(err)
		}

		isValid := checkPassword(hashedPassword, pass)
		if isValid {
			fmt.Println("Password is correct")
		}

		// Create a session
		sessionDuration := time.Hour * 24

		token, err := generateSessionToken()
		if err != nil {
			fmt.Println("error creating token: ")
			panic(err)
		}
		sessionID, err := CreateSession(db, token, sessionDuration)
		if err != nil {
			fmt.Println("error creating session")
			panic(err)
		}

		fmt.Println("Session created with ID:", sessionID)

		http.Redirect(w, r, "/", http.StatusSeeOther)

	}
}
