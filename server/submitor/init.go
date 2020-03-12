package submitor

// TODO add json logging

import (
	"context"
	"encoding/json"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/mitchellh/go-homedir"
	"io/ioutil"

	//log "github.com/sirupsen/logrus"
	"fmt"
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

	//log.Info(images)

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

	filePath, err := homedir.Expand("marker")
	if err != nil {
		log.Fatal("failed to expand file path ")
	}

	dockerBuildContext, err := archive.TarWithOptions(filePath, &archive.TarOptions{})
	defer dockerBuildContext.Close()
	if err != nil {
		log.Fatal("Unable to create docker context")
	}

	buildOptions := types.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:       []string{imageName},
	}

	buildResponse, err := cli.ImageBuild(ctx, dockerBuildContext, buildOptions)
	defer buildResponse.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	response, err := ioutil.ReadAll(buildResponse.Body)
	if err != nil {
		fmt.Printf("%s", err.Error())
	}

	var data map[string]interface{}

	if err := json.Unmarshal(response, &data); err != nil {
		log.Info("Failed to parse")
		log.Info(err)
	}

	fmt.Println(data)

	//defer func() {

	//

	//}()

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

type Output struct {
	Output struct {
		Stdout    string `json:"stdout"`
		Functions []struct {
			Output string `json:"output"`
		} `json:"functions"`
	} `json:"output"`
}
