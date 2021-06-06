package controller

import (
	"github.com/jkrus/test_sarama/app/pkg/repository"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"net/http"
	"path/filepath"
)

func GetTable(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	table, err := repository.ReadTable()
	if err != nil {
		http.Error(rw, err.Error(), 400)
		return
	}

	main := filepath.Join("public", "html", "tablePage.html")
	tmpl, err := template.ParseFiles(main)
	if err != nil {
		http.Error(rw, err.Error(), 400)
	}

	err = tmpl.ExecuteTemplate(rw, "table", table)
	if err != nil {
		http.Error(rw, err.Error(), 400)
	}
}
