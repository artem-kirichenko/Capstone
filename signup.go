package main

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"net/http"
	"net/smtp"
	"os"
	"path/filepath"
)

func signup_page(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/signup.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	tmpl.ExecuteTemplate(w, "signup", nil)
	if err != nil {
		fmt.Fprintf(w, "Error executing template: %s", err.Error())
	}
}

// hashing the password
func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// generating token
func generateToken() (string, error) {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

type MailConfig struct {
	Gmail struct {
		Host     string `json:"host"`     // SMTP host (e.g., smtp.gmail.com)
		Port     string `json:"port"`     // SMTP port (e.g., 587 or 465)
		UserName string `json:"username"` // Email address used to send emails
		Password string `json:"password"` // Password or application-specific password
	} `json:"gmail"`
}

// sends a confirmation email to the specified recipient
func sendConfirmationEmail(email, link string) error {
	// Get the current working directory
	pwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	// Read the email configuration file
	bytes, err := os.ReadFile(filepath.Join(pwd, "env", "env.json"))
	if err != nil {
		return fmt.Errorf("failed to read configuration file: %w", err)
	}

	// Parse the configuration JSON into the MailConfig struct
	var gmail MailConfig

	err = json.Unmarshal(bytes, &gmail)
	if err != nil {
		return fmt.Errorf("failed to parse configuration file: %w", err)
	}

	//Construct the email message
	subject := "Subject: Account Confirmation\n"
	contentType := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := fmt.Sprintf(`
		<h1>Welcome!</h1>
		<p>Thank you for registering. Please confirm your account by clicking the link below:</p>
		<a href="%s">Confirm Account</a>
	`, link)

	message := subject + contentType + body

	// SMTP server configuration
	smtpHost := gmail.Gmail.Host
	smtpPort := gmail.Gmail.Port

	// Sender data
	from := gmail.Gmail.UserName
	password := gmail.Gmail.Password

	// Set up SMTP authentication
	auth := smtp.PlainAuth("", from, password, smtpHost)
	to := email

	// Sending email
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil

}

func signup(ctx context.Context, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
		r.ParseForm()
		first_name := r.FormValue("first_name")
		last_name := r.FormValue("last_name")
		email := r.FormValue("email")
		phone := r.FormValue("phone")
		dob := r.FormValue("dob")
		pass := r.FormValue("pass")

		// hashing password
		hashedPassword, err := hashPassword(pass)

		_, err = db.Exec("INSERT INTO `users` SET first_name = ?, last_name = ?, email = ?, phone = ?, dob = ?, pass = ?, status = 0, role = 3",
			first_name, last_name, email, phone, dob, hashedPassword)
		if err != nil {
			panic(err)
		}

		// getting token
		token, err := generateToken()
		if err != nil {
			fmt.Println("Error generating token:", err)
			return
		}

		//user_id := 0

		//send token to database
		_, err = db.Exec("INSERT INTO `email_confirmation_tokens` SET user_id = ?, token = ?, expiration = ?, status = 0",
			first_name, last_name, email, phone, dob, hashedPassword)
		if err != nil {
			panic(err)
		}
		confirmationLink := "http://localhost:8080/activate/" + token

		// sending confirmation email
		err = sendConfirmationEmail(email, confirmationLink)
		if err != nil {
			fmt.Printf("Error sending email: %v\n", err)
		}

		http.Redirect(w, r, "/view_users/", http.StatusSeeOther)
	}
}
