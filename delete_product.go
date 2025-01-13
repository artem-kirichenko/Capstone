package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
)

func delete_product_page(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/delete_product.html", "templates/header.html", "templates/footer.html")
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
	err = db.QueryRow("SELECT `id`, `name` FROM `products` WHERE id = ?", productID).Scan(&product.Id, &product.Name)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	tmpl.ExecuteTemplate(w, "delete_product", product)
	if err != nil {
		fmt.Fprintf(w, "Error executing template: %s", err.Error())
	}
}

func delete_product(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	productID := r.FormValue("id")

	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/auction")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM `products` WHERE id = ?", productID)
	if err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/view_products/", http.StatusSeeOther)
}
