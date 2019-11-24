package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

// CreateHandler handlerfunc create
func CreateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	table := vars["table"]
	var data Data
	data.Menu = Menu
	data.Table = Tables[table]
	err := templateHelper(
		w,
		data,
		"create.html",
		fmt.Sprintf("%s/assets/create.html", ROOT_FOLDER))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// CreateHandlerPost handlerfunc create
func CreateHandlerPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	table := vars["table"]
	tableStruct := Tables[table]
	token := ""
	/*
		token, err := r.Cookie("token")
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
	*/
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "form error", http.StatusInternalServerError)
		return
	}
	data := make(map[string]interface{})
	for key, value := range r.PostForm {
		if key == tableStruct.PrimaryKey {
			continue
		}
		data[key] = value[0]
	}
	b, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp, err := httpClientHelper(
		token,
		http.MethodPost,
		fmt.Sprintf("%s/%s/public/%s", PREST_ENDPOINT, DATABASE, table),
		bytes.NewReader(b))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/padmin/list/%s", table), http.StatusSeeOther)
	return
}
