package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

// EditHandler handlerfunc edit
func EditHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	table := vars["table"]
	key := vars["key"]
	t, _ := Tables[table]
	token := ""
	/*
		token, err := r.Cookie("token")
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
	*/
	listURL := fmt.Sprintf("%s/gocrud/public/%s?id=$eq.%s", PREST_ENDPOINT, table, key)
	resp, err := httpClientHelper(
		token,
		http.MethodGet,
		listURL,
		nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewDecoder(resp.Body).Decode(&t.Data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var data Data
	data.Menu = Menu
	data.Table = t
	err = templateHelper(
		w,
		data,
		"edit.html",
		fmt.Sprintf("%s/assets/edit.html", ROOT_FOLDER))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// EditHandlerPut handlerfunc update
func EditHandlerPut(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	table := vars["table"]
	key := vars["key"]
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
	fmt.Println("data: ", data)
	resp, err := httpClientHelper(
		token,
		http.MethodPut,
		fmt.Sprintf("%s/gocrud/public/%s?id=$eq.%s", PREST_ENDPOINT, table, key),
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
