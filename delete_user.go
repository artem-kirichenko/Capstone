package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
)

func delete_user_page(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/delete_user.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	// Extract the user ID from the URL
	vars := mux.Vars(r)
	userID := vars["id"]

	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/auction")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var user User
	err = db.QueryRow("SELECT `id`, `first_name`, `last_name`"+
		" FROM `users` WHERE id = ?", userID).Scan(&user.Id, &user.FirstName, &user.LastName)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	tmpl.ExecuteTemplate(w, "delete_user", user)
	if err != nil {
		fmt.Fprintf(w, "Error executing template: %s", err.Error())
	}
}

func delete_user(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	userID := r.FormValue("id")

	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/auction")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM `users` WHERE id = ?", userID)
	if err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/view_users/", http.StatusSeeOther)
}
