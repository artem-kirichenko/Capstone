package main

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

func view_users(ctx context.Context, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("templates/view_users.html", "templates/header.html", "templates/footer.html")
		if err != nil {
			_, _ = fmt.Fprint(w, err.Error())
			return
		}

		//getting data from mysql
		row, err := db.Query(
			"SELECT `id`, `first_name`, `last_name`, `email`, `phone`, `dob`, `status`, `role`, `registration_date` " +
				"FROM `users`")
		if err != nil {
			panic(err)
		}

		var users []User
		for row.Next() {
			var u User
			var dob string
			err = row.Scan(&u.Id, &u.FirstName, &u.LastName, &u.Email, &u.Phone, &dob, &u.Status, &u.Role, &u.RegistrationDate)
			if err != nil {
				panic(err)
			}
			layout := "2006-01-02"
			u.Dob, _ = time.Parse(layout, dob)
			u.DisplayDoB = dob
			users = append(users, u)
		}
		if err := row.Err(); err != nil {
			panic(err)
		}

		tmpl.ExecuteTemplate(w, "view_users", users)
		if err != nil {
			fmt.Fprintf(w, "Error executing template: %s", err.Error())
		}
	}

}
