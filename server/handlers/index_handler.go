package handlers

import (
	log "github.com/sirupsen/logrus"
	"html/template"
	"net/http"
)

type Lab struct {
	Name             string `yaml:"name"`
	ID               string `yaml:"id"`
	ProblemStatement string `yaml:"problem_statement"`
	Testcase         []struct {
		Expected []struct {
			Feedback string   `yaml:"feedback"`
			Points   float64  `yaml:"points"`
			Values   []string `yaml:"values"`
		} `yaml:"expected"`
		Type   string   `yaml:"type"`
		Inputs []string `yaml:"inputs"`
	} `yaml:"testcase"`
}

type logerror struct {
	goError     error
	errortype   string
	info        string
	oldFileName string
	newFileName string
}

func Handlemain(w http.ResponseWriter, r *http.Request, template_path string, labs []Lab) {

	log.Info("Got request")

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	t := template.Must(template.ParseFiles(template_path + "/index.html"))
	err := t.ExecuteTemplate(w, "index.html", labs)
	if err != nil {
		log.WithFields(log.Fields{"Error": err}).Info("Template is missing")
		return
	}
}
