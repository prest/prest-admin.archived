package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
)

var PREST_ENDPOINT string
var ROOT_FOLDER string

type Data struct {
	Menu  map[string][]string
	Table Table
}

type Table struct {
	PrimaryKey  string   `json:"primary_key"`
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

var Tables map[string]Table

var Menu map[string][]string

// Load templates
func Load() (err error) {
	PREST_ENDPOINT = os.Getenv("PREST_ENDPOINT")
	if PREST_ENDPOINT == "" {
		PREST_ENDPOINT = "http://localhost:8001"
	}
	ROOT_FOLDER = "."
	err = loadFilesTable()
	return
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

func templateHelper(w http.ResponseWriter, data Data, name string, file string) (err error) {
	base := fmt.Sprintf("%s/assets/base.html", ROOT_FOLDER)
	tpl, err := template.New(name).ParseFiles(base, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tpl.ExecuteTemplate(w, "menu", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = tpl.ExecuteTemplate(w, "content", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return
}

func getTablesByJson() (files []string, err error) {
	dir := fmt.Sprintf("%s/tables", ROOT_FOLDER)
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
	Menu = make(map[string][]string)
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
		Menu["Tables"] = append(Menu["Tables"], table.Name)
	}
	fmt.Println(Menu)
	return
}
