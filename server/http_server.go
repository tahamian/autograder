package server

import (
	"autograder/server/handlers"
	//"autograder/server/handler"
	"autograder/server/submitor"
	"context"
	"github.com/gorilla/mux"

	"github.com/ulule/limiter"
	"github.com/ulule/limiter/drivers/middleware/stdlib"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type ConfigServer struct {
	Labs         []handlers.Lab  `yaml:"labs"`
	Redis        Redis           `yaml:"redis"`
	ServerConfig ServerConfig    `yaml:"server"`
	Logging      LogConfig       `yaml:"logging"`
	Marker       handlers.Marker `yaml:"marker"`
}

type ServerConfig struct {
	ServerPort   string `yaml:"server_port"`
	Host         string `yaml:"host"`
	WriteTimeout int32  `yaml:"write_timeout"`
	ReadTimeout  int32  `yaml:"read_timeout"`
	TemplatePath string `yaml:"template_path"`
}

type LogConfig struct {
	LogDir   string `yaml:"log_dir"`
	LogLevel string `yaml:"log_level"`
}

type LogHTTPSever struct {
	server_status string
	message       string
	error         string
	error_type    string
}

var config = create_config()

func create_config() ConfigServer {
	c := ConfigServer{}

	c.Redis.MaxRetry = 3
	c.Redis.RedisServer = "0.0.0.0:6738"
	c.Redis.RateLimiter = "50-H"
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

// Initializes the server builds marker docker image and
// TODO
func StartServer(config_path string) *HTMLServer {

	config.getConf(config_path)
	//log.Info("Creating docker image .....")

	// make this a go func
	//log.Info(config.Marker.ImageName)
	submitor.BuildImage(config.Marker.ImageName)

	store, rate := initalize_redis(config.Redis)
	middleware := stdlib.NewMiddleware(limiter.New(store, rate))

	router := mux.NewRouter()

	router.Handle("/", middleware.Handler(http.HandlerFunc(indexHander)))
	router.HandleFunc("/upload", upload)
	router.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir("./templates/js/"))))
	router.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("./templates/css/"))))

	htmlServer := HTMLServer{
		server: &http.Server{
			Addr:           config.ServerConfig.Host + ":" + config.ServerConfig.ServerPort,
			Handler:        router,
			ReadTimeout:    time.Second * time.Duration(config.ServerConfig.ReadTimeout),
			WriteTimeout:   time.Second * time.Duration(config.ServerConfig.WriteTimeout),
			MaxHeaderBytes: 1 << 20,
		},
	}

	htmlServer.wg.Add(1)

	go func() {
		log.Info("HTMLServer : Service started : Host=", config.ServerConfig.Host, ":", config.ServerConfig.ServerPort)
		err := htmlServer.server.ListenAndServe()
		if err != nil {
			log.Info("HTTP server failed to start ", err)
		}
		htmlServer.wg.Done()
	}()

	return &htmlServer

}

func (htmlServer *HTMLServer) Stop() error {
	const timeout = 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	//fmt.Printf("\nHTMLServer : Service stopping\n")
	if err := htmlServer.server.Shutdown(ctx); err != nil {
		if err := htmlServer.server.Close(); err != nil {
			log.Info("HTMLServer : Service stopping : Error=", err)
			return err
		}
	}
	htmlServer.wg.Wait()
	log.Info("HTMLServer : Stopped")
	return nil
}

type HTMLServer struct {
	server *http.Server
	wg     sync.WaitGroup
}

func indexHander(w http.ResponseWriter, r *http.Request) {
	handlers.Handlemain(w, r, config.ServerConfig.TemplatePath, config.Labs)
}

func upload(w http.ResponseWriter, r *http.Request) {
	handlers.Upload(w, r, config.ServerConfig.TemplatePath, config.Marker)
}
