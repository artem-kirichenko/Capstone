package main

import (
	"context"
	"database/sql"
	"net/http"
)

func place_bid(ctx context.Context, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
		r.ParseForm()
		price := r.FormValue("price")
		name := r.FormValue("name")
		description := r.FormValue("description")
		productID := r.FormValue("id")

		_, err := db.Exec(
			"UPDATE `products` SET name = ?, price = ?, description = ? WHERE id = ?", name, price, description, productID)
		if err != nil {
			panic(err)
		}

		//here a function for person who did bid!!!!!!!!!!!!!!!!!!!!!!!!!!!!

		http.Redirect(w, r, "/auction/", http.StatusSeeOther)
	}
}
