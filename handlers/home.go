package handlers

import (
	"fmt"
	"net/http"
)

// HomeHandler handlerfunc create
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	var data Data
	data.Menu = Menu
	err := templateHelper(
		w,
		data,
		"home.html",
		fmt.Sprintf("%s/assets/home.html", ROOT_FOLDER))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
