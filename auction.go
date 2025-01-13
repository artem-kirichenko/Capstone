package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
)

func view_auction(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles("templates/auction.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/auction")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	//getting data from mysql
	row, err := db.Query(
		"SELECT `id`, `name`, `price`, inventory.quantity AS `quantity`, products.description AS `description` " +
			"FROM `products` " +
			"JOIN inventory ON products.id = inventory.product_id")

	if err != nil {
		panic(err)
	}

	var products []Product
	for row.Next() {
		var p Product
		err = row.Scan(&p.Id, &p.Name, &p.Price, &p.Quantity, &p.Description)
		if err != nil {
			panic(err)
		}
		products = append(products, p)
	}
	if err := row.Err(); err != nil {
		panic(err)
	}

	tmpl.ExecuteTemplate(w, "auction", products)
	if err != nil {
		fmt.Fprintf(w, "Error executing template: %s", err.Error())
	}

}
