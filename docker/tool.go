package docker

import (
	"fmt"
	"github.com/docker/docker/pkg/stdcopy"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
)

func (d *Dash) HasImage() (bool, error) {

	images, err := d.DockerDasher.Client.ImageList(
		d.DockerDasher.Context,
		types.ImageListOptions{},
	)
	if err != nil {
		return false, err
	}

	for _, img := range images {
		for _, tag := range img.RepoTags {
			return strings.Contains(tag, d.DockerDasher.Image), nil
		}
	}

	return false, nil
}

func (d *Dash) PullImage(imgName string) error {

	if !strings.Contains(imgName, ":") {
		imgName = imgName + ":latest"
	}

	d.DockerDasher.Image = imgName

	out, err := d.DockerDasher.Client.ImagePull(
		d.DockerDasher.Context, d.DockerDasher.Image,
		types.ImagePullOptions{},
	)
	if err != nil {
		return err
	}

	defer func(out io.ReadCloser) {
		_ = out.Close()
	}(out)

	_, _ = io.Copy(os.Stdout, out)

	return nil
}

// monitorContainerLogs 监控容器日志
func (d *Dash) MonitorContainerLogs() {

	out, err := d.DockerDasher.Client.ContainerLogs(
		d.DockerDasher.Context,
		d.DockerDasher.ConID,
		types.ContainerLogsOptions{
			ShowStdout: true,
			ShowStderr: true,
			Follow:     true,
		},
	)
	if err != nil {
		fmt.Printf("Error retrieving container logs: %v\n", err)
		return
	}
	defer func() { _ = out.Close() }()

	// 将容器日志输出到标准输出
	_, _ = stdcopy.StdCopy(os.Stdout, os.Stderr, out)
}
