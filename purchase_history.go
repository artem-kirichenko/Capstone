package main

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

func purchase_history(ctx context.Context, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		tmpl, err := template.ParseFiles("templates/purchase_history.html", "templates/header.html", "templates/footer.html")
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}

		//getting data from mysql
		row, err := db.Query(
			"SELECT CONCAT(users.first_name, ' ', users.last_name) AS `customer_name`, products.name AS `product_name`, bids.bid AS `bide`, `purchase_date` " +
				"FROM `purchase_history` " +
				"JOIN `users` ON users.id = purchase_history.user_id " +
				"JOIN `products` ON products.id = purchase_history.product_id " +
				"JOIN `bids` ON bids.user_id = purchase_history.user_id AND bids.product_id = purchase_history.product_id")

		if err != nil {
			panic(err)
		}

		var purchases []Inventory
		for row.Next() {
			var p Inventory
			var purchaseDate string
			err = row.Scan(&p.Customer, &p.ProductName, &p.Price, &purchaseDate)
			if err != nil {
				panic(err)
			}

			layout := "2006-01-02"
			p.PurchaseDate, _ = time.Parse(layout, purchaseDate)
			p.DisplayPurchaseDate = purchaseDate

			purchases = append(purchases, p)
		}
		if err := row.Err(); err != nil {
			panic(err)
		}

		tmpl.ExecuteTemplate(w, "purchase_history", purchases)
		if err != nil {
			fmt.Fprintf(w, "Error executing template: %s", err.Error())
		}
	}
}
