package handlers

import (
	"autograder/server/submitor"
	"errors"
	"github.com/sirupsen/logrus"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var log = logrus.New()

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

		log.Info(lab_num)

		id := time.Now().Format("20060102150405") + strconv.Itoa(rand.Intn(1000))
		log.Info(id)

		f, err := os.OpenFile("./files/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
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

		var a = &submitor.Submission{
			ContainerName: id,
			ImageName:     h.Marker.ImageName,
			TargetDir:     h.Marker.MountPath,
			Command:       []string{"ls"},
			BindedDir:     "/usr",
		}

		submitor.CreateContainer(a)

		//_, err = io.Copy(f, file)
		//if err != nil {
		//	template_handler(w, r, "Internal Server Error", (logerror{goError: err, errortype: "",
		//		info: "There was an error when trying to copy file", oldFileName: "./files/" + handler.Filename,
		//		newFileName: "./files/" + ".py"}))
		//	return
		//}

		//Generate unique id

		//Rename the file with the unique id
		//err = os.Rename("./files/"+handler.Filename, "./files/"+id+".py")
		//if err != nil {
		//	template_handler(w, r, "Internal Server Error", (logerror{goError: err,
		//		errortype: "There was an error when trying to Rename file", info: "Rename file error",
		//		oldFileName: "./files/" + handler.Filename, newFileName: "./files/" + id + ".py"}))
		//	return
		//}

		// Run the tests with the given parameter
		//cmd := exec.Command("python", "./scripts/main.py", "./files/"+id+".py", lab_num)
		// Use a bytes.Buffer to get the output
		//var buf bytes.Buffer
		//var stderr bytes.Buffer
		//
		//cmd.Stderr = &stderr
		//cmd.Stdout = &buf
		//
		//cmd.Start()
		//if err != nil {
		//	template_handler(w, r, "File Could not run"+stderr.String(), (logerror{goError: err,
		//		errortype: "File Did not run", info: "Error when running command",
		//		oldFileName: "./files/" + handler.Filename, newFileName: "./files/" + id + ".py"}))
		//	return
		//}
		//// Use a channel to signal completion so we can use a select statement
		//done := make(chan error)
		//go func() { done <- cmd.Wait() }()
		//
		//// Start a timer
		//timeout := time.After(3 * time.Second)
		//
		//// The select statement allows us to execute based on which channel
		//// we get a message from first.
		//select {
		//case <-timeout:
		//	// Timeout happened first, kill the process and print a message.
		//	cmd.Process.Kill()
		//	template_handler(w, r, "Python Script ran for too long", (logerror{goError: err,
		//		errortype: "Command timed out", info: "Timeout", oldFileName: "./files/" + handler.Filename,
		//		newFileName: "./files/" + ".py"}))
		//	del_file(id)
		//	return
		//case err := <-done:
		//	// Command completed before timeout. Print output and error if it exists.
		//	// fmt.Println("Output:", buf.String())
		//	if err != nil {
		//		template_handler(w, r, "Error when running your script"+"\n"+err.Error(), (logerror{goError: err,
		//			errortype: "Exit Status non zero", info: err.Error(),
		//			oldFileName: "./files/" + handler.Filename, newFileName: "./files/" + ".py"}))
		//		del_file(id)
		//		return
		//	}
		//}
		//
		////Remove File
		//defer del_file(id)

		//Give back result
		t, err := template.ParseFiles(h.Template_path + "/result.html")

		err = t.Execute(w, "tests")
		if err != nil {
			log.WithFields(logrus.Fields{"Error": err}).Info("Template is missing")
			return
		}

	}
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
	err = t.Execute(w, errorname)
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
