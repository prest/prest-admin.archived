package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/gorilla/mux"
)

type Table struct {
	Name        string   `json:"name"`
	Columns     []Column `json:"columns"`
	Data        []map[string]interface{}
	ListColumns string `json:"list_columns"`
}
type Column struct {
	Label string `json:"label"`
	Name  string `json:"name"`
	Type  string `json:"type"`
}

var tmplList *template.Template
var tmplCreate *template.Template
var tmplEdit *template.Template

var Tables map[string]Table

// Load templates
func Load() (err error) {
	tmplList = template.Must(template.ParseFiles("./assets/list.html"))
	tmplCreate = template.Must(template.ParseFiles("./assets/create.html"))
	tmplEdit = template.Must(template.ParseFiles("./assets/edit.html"))
	err = loadFilesTable()
	return
}

// CreateHandler handlerfunc create
func CreateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	table := vars["table"]
	data, _ := Tables[table]
	err := templateHelper(
		w,
		data,
		"create.html",
		"./assets/create.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// CreateHandlerPost handlerfunc create
func CreateHandlerPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	table := vars["table"]
	//data, _ := Tables[table]
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
		fmt.Sprintf("http://localhost:3000/gocrud/public/%s", table),
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

// EditHandler handlerfunc edit
func EditHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	table := vars["table"]
	key := vars["key"]
	data, _ := Tables[table]
	/*
		err := templateHelper(
			w,
			data,
			"edit.html",
			"./assets/edit.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	*/
	token := ""
	/*
		token, err := r.Cookie("token")
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
	*/
	listURL := fmt.Sprintf("http://localhost:3000/gocrud/public/%s?id=$eq.%s", table, key)
	resp, err := httpClientHelper(
		token,
		http.MethodGet,
		listURL,
		nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewDecoder(resp.Body).Decode(&data.Data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = templateHelper(
		w,
		data,
		"edit.html",
		"./assets/edit.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

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
	endpoint := fmt.Sprintf("http://localhost:3000/gocrud/public/%s?_select=%s", table.Name, table.ListColumns)
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
	err = templateHelper(
		w,
		table,
		"list.html",
		"./assets/list.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func httpClientHelper(token, method, url string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return
	}
	req.Header.Add("authorization", fmt.Sprintf("Bearer %v", token))

	client := http.Client{}
	resp, err = client.Do(req)
	return
}

func templateHelper(w http.ResponseWriter, table Table, name string, file string) (err error) {
	tpl, err := template.New(name).ParseFiles(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tpl.Execute(w, table)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return
}

func getTablesByJson() (files []string, err error) {
	dir := "./tables"
	ext := "*.json"
	tempFiles, err := filepath.Glob(filepath.Join(dir, ext))
	if err != nil {
		return files, err
	}
	files = append(files, tempFiles...)
	return
}

func loadFilesTable() (err error) {
	files, err := getTablesByJson()
	if err != nil {
		return
	}
	Tables = make(map[string]Table)
	for _, f := range files {
		c, err := ioutil.ReadFile(f)
		if err != nil {
			return errors.New("loadFilesTable: error loaded table " + f)
		}
		var table Table
		table.ListColumns = "*"
		err = json.Unmarshal(c, &table)
		if err != nil {
			return errors.New("loadFilesTable: error unmarshal table " + f)
		}
		Tables[table.Name] = table
	}
	return
}
