package submitor

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"io/ioutil"

	"strings"
)

func stringInSlice(a string, list []string) bool {

	for _, b := range list {
		match := strings.Split(b, ":")[0]
		if a == match {
			return true
		}
	}
	return false
}

type ImageBuildLog struct {
	ImageName string
	Error     string
}

func SetLogger(logger *logrus.Logger) {
	log = logger
}

type BuildImageLog struct {
	state      string
	error      string
	error_type string
}

func BuildImage(imageName string) {

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

			deleted, err := cli.ImageRemove(ctx, images[i].ID, removalOptions)
			if err != nil {
				log.Fatal("Failed to delete old autograder image and build a new one, %v", err)
			}

			log.Info("deleted image: ", deleted)
		}
	}

	filePath, err := homedir.Expand("marker")
	if err != nil {
		log.Fatal("failed to expand file path")
	}

	dockerBuildContext, err := archive.TarWithOptions(filePath, &archive.TarOptions{})
	if err != nil {
		log.Fatal("Unable to create docker context")
	}

	defer func() {
		if err = dockerBuildContext.Close(); err != nil {
			log.Warn("unable to close build context")
		}
	}()

	buildOptions := types.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:       []string{imageName},
	}

	buildResponse, err := cli.ImageBuild(ctx, dockerBuildContext, buildOptions)
	if err != nil {
		log.Fatal(err)
	}

	response, err := ioutil.ReadAll(buildResponse.Body)
	if err != nil {
		log.Warn("Failed to read build response body")
	}

	defer func() {
		if err = buildResponse.Body.Close(); err != nil {
			log.Warn("Failed to close build response body")
		}
	}()

	log.WithFields(logrus.Fields{
		"autograder": "Built image",
	}).Info(strings.TrimSpace(string(response)))
}

type FunctionType string

const (
	String  FunctionType = "str"
	Int     FunctionType = "int"
	Float   FunctionType = "float"
	Complex FunctionType = "complex"
	List    FunctionType = "list"
	Dict    FunctionType = "dict"
)

type SubmitPayload struct {
	Filename  string `json:"filename"`
	Stdout    bool   `json:"stdout"`
	Functions []struct {
		FunctionName string `json:"function_name"`
		FunctionArgs []struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		} `json:"function_args"`
	} `json:"functions"`
}
