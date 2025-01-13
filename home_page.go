package main

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
)

func home_page(ctx context.Context, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		tmpl, err := template.ParseFiles("templates/home_page.html", "templates/header.html", "templates/footer.html")
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}

		tmpl.ExecuteTemplate(w, "home_page", nil)

	}
}
