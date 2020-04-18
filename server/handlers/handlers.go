package handlers

import (
	"autograder/server/submitor"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"html/template"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var log = logrus.New()

type Lab struct {
	Name             string     `yaml:"name"`
	ID               string     `yaml:"id"`
	ProblemStatement string     `yaml:"problem_statement"`
	Testcase         []Testcase `yaml:"testcase"`
}

type Testcase struct {
	Expected  []Expected `yaml:"expected"`
	Type      string     `yaml:"type"`
	Name      string     `yaml:"name"`
	Functions []Function `yaml:"functions"`
}

type Expected struct {
	Feedback string   `yaml:"feedback"`
	Points   float32  `yaml:"points"`
	Values   []string `yaml:"values"`
}

type Input struct {
	Filename  string     `yaml:"filename" json:"filename"`
	Stdout    bool       `yaml:"stdout" json:"stdout"`
	Functions []Function `yaml:"functions" json:"functions"`
}

type Function struct {
	FunctionName string        `yaml:"function_name" json:"function_name"`
	FunctionArgs []FunctionArg `yaml:"function_args" json:"function_args"`
}

type FunctionArg struct {
	Type  string      `yaml:"type" json:"type"`
	Value interface{} `yaml:"value" json:"value"`
}

type Output struct {
	Stdout    string `json:"stdout"`
	Functions []struct {
		Result       string        `json:"result"`
		Result_Type  string        `json:"result_type"`
		Status       int           `json:"status"`
		FunctionName string        `json:"function_name"`
		FunctionArgs []FunctionArg `json:"function_args"`
		Buffer       string        `json:"buffer"`
	} `json:"functions"`
}

type Result struct {
	Evaluations []Evaluation
	TotalPoints float32
}

type Evaluation struct {
	Type   string
	Actual string
	Status string
	Points float32
}

func output_from_file(filepath string) (*Output, error) {
	o := &Output{}
	jsonFile, err := os.Open(filepath)
	if err != nil {
		log.Warn("unable to load json data")
		return o, err
	}
	jsonData, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Warn("unable to read json data")
		return o, err
	}
	err = json.Unmarshal(jsonData, o)
	if err != nil {
		log.Warn("unable to load json data into struct")
		return o, err
	}
	return o, nil
}

func (l *Lab) to_input() *Input {

	input := Input{}
	for _, test_case := range l.Testcase {
		if test_case.Type == "stdout" {
			input.Stdout = true
		}

		if test_case.Type == "function" {
			input.Functions = append(input.Functions, test_case.Functions...)
		}
	}
	return &input
}

// TODO write the evualtor function
/*
	Is it easier to evaluate function in python ?
	then have a result json
*/
func (l *Lab) evaluate_lab(output *Output) (*Result, error) {
	// for all test cases if the type, and function name matches and parameters check if the values are correct
	var points float32
	var actual string
	var correct bool

	result := Result{}
	for _, test_case := range l.Testcase {
		correct = false
		if test_case.Type == "stdout" {

			for _, expect := range test_case.Expected {
				for _, value := range expect.Values {
					if value == output.Stdout {
						points = expect.Points
						actual = expect.Feedback
						correct = true
					}
				}
			}

			if !correct {
				actual = "Wrong answer try again"
			}

			e := Evaluation{
				Type:   test_case.Type,
				Actual: output.Stdout,
				Status: actual,
				Points: points,
			}

			result.Evaluations = append(result.Evaluations, e)
		} else if test_case.Type == "function" {
			for _, f := range test_case.Functions {
				for _, o := range output.Functions {
					if f.FunctionName == o.FunctionName && check_function_args(f.FunctionArgs, o.FunctionArgs) {

						e := Evaluation{
							Type:   test_case.Type,
							Actual: o.Result,
						}

						for _, expect := range test_case.Expected {
							for _, value := range expect.Values {
								if value == o.Result {
									e.Points = expect.Points
									e.Status = expect.Feedback
									correct = true
								}
							}
						}

						if !correct {
							e.Status = "Wrong function answer"
						}

						result.Evaluations = append(result.Evaluations, e)
					}
				}

			}

		}

	}

	return &result, nil
}

func check_function_args(arg1 []FunctionArg, arg2 []FunctionArg) bool {
	for _, i := range arg1 {
		match := false
		for _, j := range arg2 {
			if i.Value == j.Value {
				match = true
			} else {
				match = false
			}
		}
		if match {
			return true
		}
	}
	return false
}

func (i *Input) to_json(filename string) error {
	input, err := json.Marshal(i)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, input, 0644)
	if err != nil {
		return err
	}
	return nil
}

type logerror struct {
	goError     error
	errortype   string
	info        string
	oldFileName string
	newFileName string
}

type HandlerFunc interface {
	HandleIndex(w http.ResponseWriter, r *http.Request)
	Upload(w http.ResponseWriter, r *http.Request)
}

type Handler struct {
	Template_path string
	Marker        Marker
	Labs          []Lab
}

func (h *Handler) HandleIndex(w http.ResponseWriter, r *http.Request) {

	log.Info("Got request")

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	t := template.Must(template.ParseFiles(h.Template_path + "/index.html"))
	err := t.ExecuteTemplate(w, "index.html", h.Labs)
	if err != nil {
		log.WithFields(logrus.Fields{"Error": err}).Info("Template Does not exisit")
		return
	}
}

func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/upload" {
		http.NotFound(w, r)
		return
	}

	if r.Method == "GET" {
		log.Info(w, "404 page not found")
		http.NotFound(w, r)
	} else {

		r.Body = http.MaxBytesReader(w, r.Body, 20*1024)
		err := r.ParseMultipartForm(20)
		if err != nil {
			template_handler(w, r, err, "File size too big", h.Template_path)
			return
		}

		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			template_handler(w, r, err, "Could not get file from post request", h.Template_path)
			return
		}

		defer func() {
			err = file.Close()
			if err != nil {
				log.Warn("Failed to close file", err)
			}
		}()

		err = r.ParseForm()

		if err != nil {
			template_handler(w, r, err, "Could Not Parse form", h.Template_path)
			return
		}

		if !check_upload_file_extention(handler.Filename, []string{"py"}) {
			template_handler(w, r, err, "Invalid file extention uploaded", h.Template_path)
			return
		}

		lab_num := r.Form.Get("labs")
		if err != nil {
			template_handler(w, r, err, "Could not get lab number", h.Template_path)
			return
		}

		lab_selected, err := get_lab(h.Labs, lab_num)
		if err != nil {
			template_handler(w, r, err, "failed to get lab id", h.Template_path)
			return
		}

		fmt.Println(lab_selected)
		id := time.Now().Format("20060102150405") + strconv.Itoa(rand.Intn(1000))
		bindedDir := h.Marker.SubmissionFolderPath + id

		err = os.MkdirAll(bindedDir, os.ModePerm)
		if err != nil {
			template_handler(w, r, err, "Could not create a directory", h.Template_path)
			return
		}

		f, err := os.OpenFile(bindedDir+"/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			template_handler(w, r, err, "Failed to save file", h.Template_path)
			return
		}
		defer func() {
			err = f.Close()
			if err != nil {
				log.Warn("Failed to close file", err)
			}
		}()

		_, err = io.Copy(f, file)
		if err != nil {
			log.Warn("could not copy")
		}

		absoluteBindedDir, err := filepath.Abs(bindedDir)
		if err != nil {
			template_handler(w, r, err, "Unable to get abs path of dir", h.Template_path)
			return
		}

		// Write to json
		input := lab_selected.to_input()
		input.Filename = h.Marker.MountPath + handler.Filename
		err = input.to_json(absoluteBindedDir + "/input.json")
		if err != nil {
			template_handler(w, r, err, "unable to get input", h.Template_path)
		}

		var a = &submitor.Submission{
			ContainerName: id,
			ImageName:     h.Marker.ImageName,
			TargetDir:     h.Marker.MountPath,
			BindedDir:     absoluteBindedDir,
		}

		submitor.CreateContainer(a)
		output, err := output_from_file(absoluteBindedDir + "/output.json")
		if err != nil {
			log.Warn("Failed", err)
		}

		evaluation, err := lab_selected.evaluate_lab(output)
		if err != nil {
			log.Warn("Evaluation failed", err)
		}

		t, err := template.ParseFiles(h.Template_path + "/result.html")

		err = t.Execute(w, evaluation)
		if err != nil {
			log.WithFields(logrus.Fields{"Error": err}).Info("Template is missing")
			return
		}

	}
}

func get_lab(labs []Lab, lab_id string) (Lab, error) {
	for _, lab := range labs {
		if lab.ID == lab_id {
			return lab, nil
		}
	}
	err := errors.New("Invalid lab ID")
	return Lab{}, err
}

func SetLogger(logger *logrus.Logger) {
	log = logger
}

type Marker struct {
	DockerfilePath       string `yaml:"dockerfile_path"`
	SubmissionFolderPath string `yaml:"submissions_folder"`
	MountPath            string `yaml:"mount_path"`
	Command              string `yaml:"command"`
	ImageName            string `yaml:"image_name"`
}

func del_file(id string) {
	err := os.Remove("./files/" + id + ".py")
	if err != nil {
		log.WithFields(logrus.Fields{"Old File Name": "./files/", "New file Name": "./files/" + id + ".py"}).Info(
			"There was an error when trying to Delete file after tests was sucessful the file")
		return
	}
}

func template_handler(w http.ResponseWriter, r *http.Request, error error, message string, template_path string) {
	t, err := template.ParseFiles(template_path + "/error.html")
	err = t.Execute(w, err.Error())
	if err != nil {
		log.WithFields(logrus.Fields{"Error": err}).Info("Template is missing")
		return
	}
	log.WithFields(logrus.Fields{"Error": error, "Error Type": error.Error()}).Info(message)

	return
}

// Checks if file exists
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func check_upload_file_extention(filename string, extentions []string) bool {

	for _, value := range extentions {
		split := strings.Split(filename, ".")
		if split[len(split)-1] == value {
			return true
		}
	}

	return false
}

func check_py_file(filename string) {

}
