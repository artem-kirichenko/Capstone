package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
)

func edit_user_page(ctx context.Context, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("templates/edit_user.html", "templates/header.html", "templates/footer.html")
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}

		// Extract the user ID from the URL
		vars := mux.Vars(r)
		userID := vars["id"]

		db, err := mysql_connect()
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		var user User
		err = db.QueryRow(
			"SELECT `id`, `first_name`, `last_name`, `email`, `phone`, `dob` "+
				"FROM `users` WHERE id = ?", userID).Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.Phone, &user.Dob)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		tmpl.ExecuteTemplate(w, "edit_user", user)
		if err != nil {
			fmt.Fprintf(w, "Error executing template: %s", err.Error())
		}
	}

}

func edit_user(ctx context.Context, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
		r.ParseForm()
		firstName := r.Form.Get("first_name")
		lastName := r.Form.Get("last_name")
		email := r.Form.Get("email")
		phone := r.Form.Get("phone")
		dob := r.Form.Get("dob")
		userID := r.Form.Get("id")

		db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/auction")
		if err != nil {
			panic(err)
		}
		defer db.Close()

		_, err = db.Exec(
			"UPDATE `users` SET first_name = ?, last_name = ?, email = ? , phone = ?, dob = ? WHERE id = ?",
			firstName, lastName, email, phone, dob, userID)
		if err != nil {
			panic(err)
		}

		http.Redirect(w, r, "/view_users/", http.StatusSeeOther)
	}
}
