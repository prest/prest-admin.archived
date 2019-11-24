package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// ListHandler handlerfunc list
func ListHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tableName := vars["table"]
	table, _ := Tables[tableName]
	token := ""
	/*
		token, err := r.Cookie("token")
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		resp, err := httpClientHelper(token.Value, http.MethodGet, fmt.Sprintf("http://localhost:3000/gocrud/public/%s", table.Name), nil)
	*/
	endpoint := fmt.Sprintf("%s/%s/public/%s?_select=%s", PREST_ENDPOINT, DATABASE, table.Name, table.ListColumns)
	resp, err := httpClientHelper(token, http.MethodGet, endpoint, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewDecoder(resp.Body).Decode(&table.Data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var data Data
	data.Menu = Menu
	data.Table = table
	err = templateHelper(
		w,
		data,
		"list.html",
		fmt.Sprintf("%s/assets/list.html", ROOT_FOLDER))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
