package utils

import (
    "fmt"
    "html/template"
    "net/http"
    "path/filepath"
)

func RenderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	
    tmplPath := filepath.Join("templates", tmpl+".html")
    fmt.Println("Template Path:", tmplPath)

    t, err := template.ParseFiles(tmplPath)
    if err != nil {
        fmt.Println("Error parsing template:", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    err = t.Execute(w, data)
    if err != nil {
        fmt.Println("Error executing template:", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    }
}

