package main

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/textproto"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"sync"
	"time"

	"autograder"

	redis "github.com/go-redis/redis"
	mux "github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/ulule/limiter"
	"github.com/ulule/limiter/drivers/middleware/stdlib"
	sredis "github.com/ulule/limiter/drivers/store/redis"
	"gopkg.in/yaml.v2"
)

type ConfigServer struct {
	Redis_server string
	Log_dir      string
	Base_dir     string
}

func upload(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/upload" {
		http.NotFound(w, r)
		return
	}
	if r.Method == "GET" {
		fmt.Fprintf(w, "404 page not found")
	} else {
		// Set max file size
		r.Body = http.MaxBytesReader(w, r.Body, 20*1024)
		err := r.ParseMultipartForm(20)
		if err != nil {
			template_handler(w, r, err.Error(), (logerror{goError: err, errortype: "", info: "File size is too big", oldFileName: "./files/", newFileName: "./files/" + ".py"}))

			return
		}

		//Handle the file
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			template_handler(w, r, "Could Not upload file", (logerror{goError: err, errortype: err.Error(), info: "Could not get file from post request", oldFileName: "./files/", newFileName: "./files/" + ".py"}))

			return
		}

		r.ParseForm()
		lab_num := r.Form.Get("labs")
		if err != nil {
			template_handler(w, r, "Could Not upload file", (logerror{goError: err, errortype: err.Error(), info: "could not get lab number from post request", oldFileName: "./files/" + handler.Filename, newFileName: "./files/" + ".py"}))
			return
		}

		// Check the file name is a .py file
		if handler.Filename[len(handler.Filename)-3:] != ".py" {
			template_handler(w, r, "Have to upload a python script", (logerror{goError: errors.New("can't work with 42"), errortype: "", info: "File extention was incorrect", oldFileName: "./files/" + handler.Filename, newFileName: "./files/" + ".py"}))
			return
		}

		defer file.Close()

		// Open the file
		f, err := os.OpenFile("./files/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			template_handler(w, r, "File did not load", (logerror{goError: err, errortype: "", info: "problem opening the file", oldFileName: "./files/" + handler.Filename, newFileName: "./files/" + ".py"}))
			return
		}

		defer f.Close()
		_, err = io.Copy(f, file)
		if err != nil {
			template_handler(w, r, "Internal Server Error", (logerror{goError: err, errortype: "", info: "There was an error when trying to copy file", oldFileName: "./files/" + handler.Filename, newFileName: "./files/" + ".py"}))
			return
		}

		//Generate unique id
		id := time.Now().Format("20060102150405") + strconv.Itoa(rand.Intn(1000))

		//Rename the file with the unique id
		err = os.Rename("./files/"+handler.Filename, "./files/"+id+".py")
		if err != nil {
			template_handler(w, r, "Internal Server Error", (logerror{goError: err, errortype: "There was an error when trying to Rename file", info: "Rename file error", oldFileName: "./files/" + handler.Filename, newFileName: "./files/" + id + ".py"}))
			return
		}

		// Run the test with the given parameter
		cmd := exec.Command("python3", "./scripts/main.py", "./files/"+id+".py", lab_num)
		// Use a bytes.Buffer to get the output
		var buf bytes.Buffer
		var stderr bytes.Buffer

		cmd.Stderr = &stderr
		cmd.Stdout = &buf

		cmd.Start()
		if err != nil {
			template_handler(w, r, "File Could not run"+stderr.String(), (logerror{goError: err, errortype: "File Did not run", info: "Error when running command", oldFileName: "./files/" + handler.Filename, newFileName: "./files/" + id + ".py"}))
			return
		}
		// Use a channel to signal completion so we can use a select statement
		done := make(chan error)
		go func() { done <- cmd.Wait() }()

		// Start a timer
		timeout := time.After(3 * time.Second)

		// The select statement allows us to execute based on which channel
		// we get a message from first.
		select {
		case <-timeout:
			// Timeout happened first, kill the process and print a message.
			cmd.Process.Kill()
			template_handler(w, r, "Python Script ran for too long", (logerror{goError: err, errortype: "Command timed out", info: "Timeout", oldFileName: "./files/" + handler.Filename, newFileName: "./files/" + ".py"}))
			del_file(id)
			return
		case err := <-done:
			// Command completed before timeout. Print output and error if it exists.
			// fmt.Println("Output:", buf.String())
			if err != nil {
				template_handler(w, r, "Error when running your script"+"\n"+err.Error(), (logerror{goError: err, errortype: "Exit Status non zero", info: err.Error(), oldFileName: "./files/" + handler.Filename, newFileName: "./files/" + ".py"}))
				del_file(id)
				return
			}
		}

		//Remove File
		del_file(id)

		//Give back result
		t, err := template.ParseFiles("/go/templates/result.html")
		t.Execute(w, buf.String())
		if err != nil {
			log.WithFields(log.Fields{"Error": err}).Info("Template is missing")
			return
		}

	}
}

type FileHeader struct {
	Filename string
	Header   textproto.MIMEHeader
}

type logerror struct {
	goError     error
	errortype   string
	info        string
	oldFileName string
	newFileName string
}

func del_file(id string) {
	err := os.Remove("./files/" + id + ".py")
	if err != nil {
		log.WithFields(log.Fields{"Old File Name": "./files/", "New file Name": "./files/" + id + ".py"}).Info("There was an error when trying to Delete file after test was sucessful the file")
		return
	}
}

func template_handler(w http.ResponseWriter, r *http.Request, errorname string, logging logerror) {
	t, err := template.ParseFiles("/go/templates/error.html")
	t.Execute(w, errorname)
	if err != nil {
		log.WithFields(log.Fields{"Error": err}).Info("Template is missing")
		return
	}
	log.WithFields(log.Fields{"Error": logging.goError, "Old File Name": logging.oldFileName, "New File Name": logging.newFileName, "Error Type": logging.errortype}).Info(logging.info)
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

func handlemain(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	t := template.Must(template.ParseFiles("./templates/index.html"))
	err := t.ExecuteTemplate(w, "index.html", "ssd")
	if err != nil {
		log.WithFields(log.Fields{"Error": err}).Info("Template is missing")
		return
	}
}

func init() {
	// Check if log file exists if not create a log folder
	log_check, err := exists("./logs")
	if err != nil {
		panic(err)
	}
	if log_check == false {
		err = os.Mkdir("./logs", os.FileMode(755))
		if err != nil {
			panic(err)
		}
	}
	logFile, err := os.Create("./logs/" + time.Now().Format("20060102150405") + ".log")
	if err != nil {
		panic(err)
	}

	// Set log format as JSON
	log.SetFormatter(&log.JSONFormatter{})

	// Set log output
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

	// Check if files folder exists and if not create it
	check, err := exists("./files")
	if err != nil {
		log.WithFields(log.Fields{"Error": "File could not be created"}).Info("File could not be created")
		return
	}
	if check == false {
		err = os.Mkdir("./files", os.FileMode(755))
		if err != nil {
			log.WithFields(log.Fields{"Error": "File could not be created"}).Info("File could not be created")
			return
		}
	}

}

func (c *ConfigServer) getConf(path string) *ConfigServer {

	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatal("Unmarshal: ", err)
	}

	return c
}

func main() {
	// Config server with timeouts

	pathptr := flag.String("config", "config.yaml", "path to yaml config file")
	flag.Parse()

	var c ConfigServer
	c.getConf(*pathptr)

	var a autograder.Labs

	a.getConf("test_case.yaml")

	serverCfg := Config{
		Host:         "0.0.0.0:9090",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	// Start the server
	htmlServer := Start(serverCfg, c)

	defer htmlServer.Stop()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan

	fmt.Println("main : shutting down")
}

func Start(cfg Config, c ConfigServer) *HTMLServer {

	// create rate limiter
	rate, err := limiter.NewRateFromFormatted("50-H")
	if err != nil {
		log.Fatal(err)
	}

	// Create a redis client.
	option, err := redis.ParseURL(c.Redis_server + "/0")
	if err != nil {
		log.Fatal(err)
	}

	client := redis.NewClient(option)

	pong, err := client.Ping().Result()
	// redis_server := strings.Replace(c.Redis_server, "redis", "http", 1)
	if err != nil {
		log.Info(err)
		for true {
			pong, err = client.Ping().Result()

			if err == nil {
				log.Info("Successful Ping", pong)
				break
			}
			time.Sleep(10 * time.Second)
			fmt.Println(err)
		}
	}

	// Create a store with the redis client.
	store, err := sredis.NewStoreWithOptions(client, limiter.StoreOptions{
		Prefix:   "limiter_http",
		MaxRetry: 3,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Create a new middleware with the limiter instance.
	middleware := stdlib.NewMiddleware(limiter.New(store, rate))

	_, cancel := context.WithCancel(context.Background())
	defer cancel()
	router := mux.NewRouter()
	router.Handle("/", middleware.Handler(http.HandlerFunc(handlemain)))
	router.HandleFunc("/upload", upload)
	// router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	htmlServer := HTMLServer{
		server: &http.Server{
			Addr:           cfg.Host,
			Handler:        router,
			ReadTimeout:    cfg.ReadTimeout,
			WriteTimeout:   cfg.WriteTimeout,
			MaxHeaderBytes: 1 << 20,
		},
	}

	htmlServer.wg.Add(1)

	go func() {
		fmt.Printf("\nHTMLServer : Service started : Host=%v\n", cfg.Host)
		htmlServer.server.ListenAndServe()
		htmlServer.wg.Done()
	}()

	return &htmlServer
}

func (htmlServer *HTMLServer) Stop() error {
	const timeout = 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	fmt.Printf("\nHTMLServer : Service stopping\n")
	if err := htmlServer.server.Shutdown(ctx); err != nil {
		if err := htmlServer.server.Close(); err != nil {
			fmt.Printf("\nHTMLServer : Service stopping : Error=%v\n", err)
			return err
		}
	}
	htmlServer.wg.Wait()
	fmt.Printf("\nHTMLServer : Stopped\n")
	return nil
}

type Config struct {
	Host         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type HTMLServer struct {
	server *http.Server
	wg     sync.WaitGroup
}
