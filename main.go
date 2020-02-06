package autograder

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/textproto"
	"os"
	"os/signal"
	"sync"
	"time"

	mux "github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/ulule/limiter"
	"github.com/ulule/limiter/drivers/middleware/stdlib"
	"gopkg.in/yaml.v2"
)

type ConfigServer struct {
	Redis_server string
	Log_dir      string
	Base_dir     string
	Rate_limit   string
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

	pathptr := flag.String("config", "config.yaml", "path to yaml config file")
	flag.Parse()

	var c ConfigServer
	c.getConf(*pathptr)

	// var a autograder.Labs
	// a.getConf("test_case.yaml")

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

	store, rate := initalize_redis(c.Redis_server, c.Rate_limit)
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
