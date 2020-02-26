package submitor

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/mitchellh/go-homedir"
	//"log"
	log "github.com/sirupsen/logrus"
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

func BuildImage(imageName string) {

	//imageName := "autograder"

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
