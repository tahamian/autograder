package server

import (
	log "github.com/sirupsen/logrus"
	"html/template"
	"net/http"
)

func handlemain(w http.ResponseWriter, r *http.Request) {

	log.Info("Got request")

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
