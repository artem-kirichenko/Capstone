package main

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
)

func view_products(ctx context.Context, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		tmpl, err := template.ParseFiles("templates/view_products.html", "templates/header.html", "templates/footer.html")
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}

		//getting data from mysql
		row, err := db.Query("SELECT `id`, `name`, `price`, `description` FROM `products`")

		if err != nil {
			panic(err)
		}

		var products []Product
		for row.Next() {
			var p Product
			err = row.Scan(&p.Id, &p.Name, &p.Price, &p.Description)
			if err != nil {
				panic(err)
			}
			products = append(products, p)
		}
		if err := row.Err(); err != nil {
			panic(err)
		}

		tmpl.ExecuteTemplate(w, "view_products", products)
		if err != nil {
			fmt.Fprintf(w, "Error executing template: %s", err.Error())
		}
	}
}
