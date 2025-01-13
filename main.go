package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"time"
)

var db *sql.DB

type User struct {
	Id               uint16
	FirstName        string
	LastName         string
	Email            string
	Phone            string
	Dob              time.Time
	DisplayDoB       string
	Status           bool
	Role             uint8
	Pass             string
	RegistrationDate string
}

type Product struct {
	Id          uint16
	Name        string
	Price       float64
	Description string
	Quantity    uint16
}

type Inventory struct {
	ProductName         string
	Customer            string
	Price               string
	PurchaseDate        time.Time
	DisplayPurchaseDate string
}

// CHECK
func renderPage(w http.ResponseWriter, r *http.Request, page string) {
	tmpl, err := template.ParseFiles("templates/header.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := struct {
		ActivePage string
	}{
		ActivePage: page,
	}
	tmpl.Execute(w, data)
}

func setupRoutes(db *sql.DB) {
	ctx := context.Background()
	rtr := mux.NewRouter()

	rtr.HandleFunc("/", home_page(ctx, db)).Methods("GET")

	//products
	rtr.HandleFunc("/view_products/", view_products(ctx, db)).Methods("GET")
	rtr.HandleFunc("/delete_product/{id:[0-9]+}", delete_product_page).Methods("GET")
	rtr.HandleFunc("/delete_product/", delete_product).Methods("POST")
	rtr.HandleFunc("/edit_product/{id:[0-9]+}", edit_product_page(ctx, db)).Methods("GET")
	rtr.HandleFunc("/edit_product/", edit_product(ctx, db)).Methods("POST")
	//add
	rtr.HandleFunc("/add_product/", add_product_page).Methods("GET")
	rtr.HandleFunc("/add_product_form/", add_product).Methods("POST")

	//place bid with sessions
	bidRtr := rtr.PathPrefix("/place_bid").Subrouter()
	bidRtr.Use(SessionMiddleware(ctx, db))
	bidRtr.HandleFunc("/{id:[0-9]+}", place_bid(ctx, db)).Methods("POST")

	rtr.HandleFunc("/auction/", view_auction).Methods("GET")
	bidRtr.HandleFunc("/place_bid/", place_bid(ctx, db)).Methods("POST")

	rtr.HandleFunc("/purchase_history/", purchase_history(ctx, db)).Methods("GET")
	//users
	rtr.HandleFunc("/view_users/", view_users(ctx, db)).Methods("GET")
	rtr.HandleFunc("/edit_user/{id:[0-9]+}", edit_user_page(ctx, db)).Methods("GET")
	rtr.HandleFunc("/edit_user/", edit_user(ctx, db)).Methods("POST")
	rtr.HandleFunc("/delete_user/{id:[0-9]+}", delete_user_page).Methods("GET")
	rtr.HandleFunc("/delete_user/", delete_user).Methods("POST")

	rtr.HandleFunc("/signup/", signup_page).Methods("GET")
	rtr.HandleFunc("/signup_submit/", signup(ctx, db)).Methods("POST")
	//login

	rtr.HandleFunc("/login_page/", login_page).Methods("GET")
	rtr.HandleFunc("/login_authorization/", login_authorization(ctx, db)).Methods("POST")

	rtr.HandleFunc("/activate/{token}", activate_page).Methods("GET")

	http.ListenAndServe(":8080", rtr)
}

func main() {
	db, err := mysql_connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("Connected to database ")
	setupRoutes(db)
}
