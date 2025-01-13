package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func activate_page(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/activate.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	//vars := mux.Vars(r)
	//token := vars["token"]

	db, err := mysql_connect()
	if err != nil {
		w.Write([]byte(err.Error())) // writing directly to the response body (to be processed by the browser)
		return                       // возвращаем пустоту т.к. функия ничего не возвращает и нам не нужно обрабатывать после
	}
	defer db.Close()

	var user User
	err = db.QueryRow(
		"SELECT `id`, `first_name`, `last_name`, `email`, `phone`, `dob`, `status`"+
			"FROM `users` WHERE id = ?", 1).Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.Phone, &user.Dob, &user.Status)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	tmpl.ExecuteTemplate(w, "activate", user.Status)
	if err != nil {
		fmt.Fprintf(w, "Error executing template: %s", err.Error())
	}
}
