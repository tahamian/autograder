package server

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func del_file(id string) {
	err := os.Remove("./files/" + id + ".py")
	if err != nil {
		log.WithFields(log.Fields{"Old File Name": "./files/", "New file Name": "./files/" + id + ".py"}).Info(
			"There was an error when trying to Delete file after tests was sucessful the file")
		return
	}
}

func template_handler(w http.ResponseWriter, r *http.Request, errorname string, logging logerror) {
	t, err := template.ParseFiles(config.TemplatePath + "/error.html")
	err = t.Execute(w, errorname)
	if err != nil {
		log.WithFields(log.Fields{"Error": err}).Info("Template is missing")
		return
	}
	log.WithFields(log.Fields{"Error": logging.goError, "Old File Name": logging.oldFileName,
		"New File Name": logging.newFileName, "Error Type": logging.errortype}).Info(logging.info)
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

func upload(w http.ResponseWriter, r *http.Request) {

	log.Info("Got request")

	if r.URL.Path != "/upload" {
		http.NotFound(w, r)
		return
	}
	if r.Method == "GET" {
		log.Info(w, "404 page not found")
	} else {
		// Set max file size
		r.Body = http.MaxBytesReader(w, r.Body, 20*1024)
		err := r.ParseMultipartForm(20)
		if err != nil {
			template_handler(w, r, err.Error(), (logerror{goError: err, errortype: "", info: "File size is too big",
				oldFileName: "./files/", newFileName: "./files/" + ".py"}))
			return
		}

		//Handle the file
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			template_handler(w, r, "Could Not upload file", (logerror{goError: err, errortype: err.Error(),
				info: "Could not get file from post request", oldFileName: "./files/",
				newFileName: "./files/" + ".py"}))
			return
		}

		err = r.ParseForm()

		if err != nil {
			template_handler(w, r, "Could Not Parse form", (logerror{goError: err, errortype: err.Error(),
				info: "Could not Parse form from post request", oldFileName: "./files/",
				newFileName: "./files/" + ".py"}))
			return
		}

		lab_num := r.Form.Get("labs")
		if err != nil {
			template_handler(w, r, "Could Not upload file", logerror{goError: err, errortype: err.Error(),
				info: "could not get lab number from post request", oldFileName: "./files/" + handler.Filename,
				newFileName: "./files/" + ".py"})
			return
		}

		fmt.Println(lab_num)

		// looks for .py files or other files
		if check_upload_file_extention(handler.Filename, []string{"py"}) {
			template_handler(w, r, "Have to upload a python script",
				logerror{goError: errors.New("can't work with 42"), errortype: "",
					info: "File extention was incorrect", oldFileName: "./files/" + handler.Filename,
					newFileName: "./files/" + ".py"})
			return
		}

		// Open the file create the file and then move it to the docker image

		defer file.Close()
		id := time.Now().Format("20060102150405") + strconv.Itoa(rand.Intn(1000))

		f, err := os.OpenFile("./files/"+handler.Filename+id, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			template_handler(w, r, "File did not load", logerror{goError: err, errortype: "",
				info: "problem opening the file", oldFileName: "./files/" + handler.Filename,
				newFileName: "./files/" + ".py"})
			return
		}

		defer func() {
			err = f.Close()
			if err != nil {
				log.Info("Failed to close file %s")
			}
		}()

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
		t, err := template.ParseFiles(config.TemplatePath + "/result.html")

		err = t.Execute(w, "tests")
		if err != nil {
			log.WithFields(log.Fields{"Error": err}).Info("Template is missing")
			return
		}

	}
}
