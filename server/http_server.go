package server

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/gorilla/mux"
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/ulule/limiter"
	"github.com/ulule/limiter/drivers/middleware/stdlib"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

type ConfigServer struct {
	BaseDir        string `yaml:"base_dir"`
	DockerfilePath string `yaml:"dockerfile_path"`
	Host           string `yaml:"host"`
	Labs           []struct {
		Name             string `yaml:"name"`
		ID               string `yaml:"id"`
		ProblemStatement string `yaml:"problem_statement"`
		Testcase         []struct {
			Expected []struct {
				Feedback string   `yaml:"feedback"`
				Points   float64  `yaml:"points"`
				Values   []string `yaml:"values"`
			} `yaml:"expected"`
			Type string `yaml:"type"`
		} `yaml:"testcase"`
	} `yaml:"labs"`
	LogDir       string `yaml:"log_dir"`
	ReadTimeout  int    `yaml:"read_timeout"`
	Redis        Redis  `yaml:"redis"`
	ServerPort   string `yaml:"server_port"`
	TemplatePath string `yaml:"template_path"`
	TestCasePath string `yaml:"test_case_path"`
	WriteTimeout int    `yaml:"write_timeout"`
}

type logerror struct {
	goError     error
	errortype   string
	info        string
	oldFileName string
	newFileName string
}

var config = create_config()

func create_config() ConfigServer {
	c := ConfigServer{}

	c.TemplatePath = "./templates"
	c.LogDir = "./logs"
	c.BaseDir = "./"

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

func stringInSlice(a string, list []string) bool {

	for _, b := range list {
		match := strings.Split(b, ":")[0]
		if a == match {
			return true
		}
	}
	return false
}

func build_image() {

	imageName := "autograder"

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal(err)
	}

	filter := types.ImageListOptions{
		All: true,
	}

	images, err := cli.ImageList(ctx, filter)

	if err != nil {
		log.Fatal("Could not list docker images %v", err)
	}

	for i := range images {
		if stringInSlice(imageName, images[i].RepoTags) {
			removalOptions := types.ImageRemoveOptions{
				Force:         true,
				PruneChildren: true,
			}

			_, err = cli.ImageRemove(ctx, images[i].ID, removalOptions)

			if err != nil {
				log.Fatal("Failed to delete old autograder image and build a new one, %v", err)
			}
		}
	}

	filePath, _ := homedir.Expand("marker")
	dockerBuildContext, _ := archive.TarWithOptions(filePath, &archive.TarOptions{})

	buildOptions := types.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:       []string{imageName},
	}

	buildResponse, err := cli.ImageBuild(ctx, dockerBuildContext, buildOptions)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err = dockerBuildContext.Close()
		if err != nil {
			log.Fatal(err)
		}

		err = buildResponse.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
}

// Initializes the server builds marker docker image and
func StartServer(config_path string) *HTMLServer {

	config.getConf(config_path)

	go build_image()

	store, rate := initalize_redis(config.Redis)
	middleware := stdlib.NewMiddleware(limiter.New(store, rate))

	//_, cancel := context.WithCancel(context.Background())
	//defer cancel()
	router := mux.NewRouter()

	router.Handle("/", middleware.Handler(http.HandlerFunc(handlemain)))
	router.HandleFunc("/upload", upload)
	router.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir("./templates/js/"))))
	router.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("./templates/css/"))))

	htmlServer := HTMLServer{
		server: &http.Server{
			Addr:           config.Host + ":" + config.ServerPort,
			Handler:        router,
			ReadTimeout:    time.Second * 5,
			WriteTimeout:   time.Second * 5,
			MaxHeaderBytes: 1 << 20,
		},
	}

	htmlServer.wg.Add(1)

	go func() {
		log.Info("HTMLServer : Service started : Host=", config.Host, ":", config.ServerPort)
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
	fmt.Printf("\nHTMLServer : Service stopping\n")
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
