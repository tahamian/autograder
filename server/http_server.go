package server

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ulule/limiter"
	"github.com/ulule/limiter/drivers/middleware/stdlib"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	log "github.com/sirupsen/logrus"
	"net/http"
	"sync"
	"time"

	//"grader/server/handlers"
)

type ConfigServer struct {
	Redis		 Redis
	Log_dir      string
	Base_dir     string
	Labs         Labs
	Template_path string
}

type Labs struct {
	Labs []Lab
}

type Lab struct {
	name        string
	LabTestCase []Testcase
}

type Testcase struct {
	Type     string
	Expected []Expected
}

type Expected struct {
	Values   []string
	Points   float32
	Feedback string
}

type logerror struct {
	goError     error
	errortype   string
	info        string
	oldFileName string
	newFileName string
}

var config = create_config()

// create config with default values if yaml file is not intialized
func create_config() ConfigServer{
	c :=  ConfigServer{}

	c.Template_path = "./templates"
	c.Log_dir = "./logs"
	c.Base_dir = "./"

	c.Redis.Max_retry = 3
	c.Redis.Redis_server = "0.0.0.0:6738"
	c.Redis.Rate_limit = "50-H"


	return c
}


func (c *ConfigServer) getConf(path string) *ConfigServer {

	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Info("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Info("Unmarshal: ", err)
	}

	return c
}

func StartServer(config_path string) *HTMLServer {

	// create and validate config
	config.getConf(config_path)

	store, rate := initalize_redis(config.Redis)

	// Create a new middleware with the limiter instance.
	middleware := stdlib.NewMiddleware(limiter.New(store, rate))

	_, cancel := context.WithCancel(context.Background())
	defer cancel()
	router := mux.NewRouter()


	router.Handle("/", middleware.Handler(http.HandlerFunc(handlemain)))
	router.HandleFunc("/upload", upload)


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