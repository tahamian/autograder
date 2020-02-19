package server

import (
	"net/http"
	"html/template"
	log "github.com/sirupsen/logrus"

)

func handlemain(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	t := template.Must(template.ParseFiles(config.TemplatePath + "/index.html"))
	err := t.ExecuteTemplate(w, "index.html", "")
	if err != nil {
		log.WithFields(log.Fields{"Error": err}).Info("Template is missing")
		return
	}
}
