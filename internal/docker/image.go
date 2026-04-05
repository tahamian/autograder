package docker

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/mitchellh/go-homedir"
	archive "github.com/moby/go-archive"
	"github.com/moby/moby/client"
	"github.com/sirupsen/logrus"
)

// EnsureImage checks if the marker image exists. If it does, it's used as-is.
// If not, it builds it from the marker directory.
func EnsureImage(log *logrus.Logger, cli Client, imageName string) error {
	ctx := context.Background()

	result, err := cli.ImageList(ctx, client.ImageListOptions{All: true})
	if err != nil {
		return fmt.Errorf("listing images: %w", err)
	}

	for _, img := range result.Items {
		if imageTagMatches(imageName, img.RepoTags) {
			log.WithField("image", imageName).Info("marker image already exists, skipping build")
			return nil
		}
	}

	log.WithField("image", imageName).Info("marker image not found, building...")
	return buildImage(log, cli, imageName)
}

func buildImage(log *logrus.Logger, cli Client, imageName string) error {
	ctx := context.Background()

	filePath, err := homedir.Expand("marker")
	if err != nil {
		return fmt.Errorf("expanding marker path: %w", err)
	}

	buildCtx, err := archive.TarWithOptions(filePath, &archive.TarOptions{})
	if err != nil {
		return fmt.Errorf("creating build context: %w", err)
	}
	defer buildCtx.Close()

	resp, err := cli.ImageBuild(ctx, buildCtx, client.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:       []string{imageName},
	})
	if err != nil {
		return fmt.Errorf("building image: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading build output: %w", err)
	}

	log.WithField("image", imageName).Info("built marker image")
	log.Debug(strings.TrimSpace(string(body)))
	return nil
}

func imageTagMatches(name string, tags []string) bool {
	for _, tag := range tags {
		if strings.Split(tag, ":")[0] == name {
			return true
		}
	}
	return false
}
