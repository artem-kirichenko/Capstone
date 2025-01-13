package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
)

func add_product_page(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/add_product.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	tmpl.ExecuteTemplate(w, "add_product", nil)
	if err != nil {
		fmt.Fprintf(w, "Error executing template: %s", err.Error())
	}
}

func add_product(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	r.ParseForm()
	name := r.FormValue("name")
	price := r.FormValue("price")
	description := r.FormValue("description")

	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/auction")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO `products` SET name = ?, price = ?, description = ?", name, price, description)
	if err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/view_products/", http.StatusSeeOther) // Перенаправление на страницу с продуктами
}
