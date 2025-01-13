package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
)

func edit_product_page(ctx context.Context, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("templates/edit_product.html", "templates/header.html", "templates/footer.html")
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}

		// Extract the product ID from the URL
		vars := mux.Vars(r)
		productID := vars["id"]

		db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/auction")
		if err != nil {
			panic(err)
		}
		defer db.Close()

		var product Product
		err = db.QueryRow(
			"SELECT `id`, `name`, `price`, `description` "+
				"FROM `products` WHERE id = ?", productID).Scan(&product.Id, &product.Name, &product.Price, &product.Description)
		if err != nil {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}

		tmpl.ExecuteTemplate(w, "edit_product", product)
		if err != nil {
			fmt.Fprintf(w, "Error executing template: %s", err.Error())
		}
	}
}

func edit_product(ctx context.Context, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
		r.ParseForm()
		name := r.FormValue("name")
		price := r.FormValue("price")
		description := r.FormValue("description")
		productID := r.FormValue("id")

		db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/auction")
		if err != nil {
			panic(err)
		}
		defer db.Close()

		_, err = db.Exec(
			"UPDATE `products` SET name = ?, price = ?, description = ? WHERE id = ?", name, price, description, productID)
		if err != nil {
			panic(err)
		}

		http.Redirect(w, r, "/view_products/", http.StatusSeeOther)
	}
}
