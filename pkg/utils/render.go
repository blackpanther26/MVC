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
        http.Error(w, "Internal Server Error from render", http.StatusInternalServerError)
        return
    }

    err = t.Execute(w, data)
    if err != nil {
        http.Error(w, "Internal Server Error from render 2", http.StatusInternalServerError)
    }
}

func RenderTemplateWithMessage(w http.ResponseWriter, tmpl, message, messageType string) {
    data := map[string]interface{}{
        "Message":     message,
        "MessageType": messageType,
    }
    RenderTemplate(w, tmpl, data)
}
