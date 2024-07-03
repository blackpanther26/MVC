package utils

import (
    "html/template"
    "net/http"
    "path/filepath"
)

func RenderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	
    tmplPath := filepath.Join("templates", tmpl+".html")

    t, err := template.ParseFiles(tmplPath)
    if err != nil {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    err = t.Execute(w, data)
    if err != nil {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    }
}

