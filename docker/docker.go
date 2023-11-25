package docker

/**
func main() {
	d, err := docker.NewDockerDash("alpine")
	if err != nil {
		return
	}
	dash := &docker.Dash{
		Command:      []string{"sh", "-c", "while true; do echo hello; sleep 1; done"},
		DockerDasher: d,
	}

	err = docker.PullAndRunAlpine(dash)
	if err != nil {
		log.Fatal(err)
		return
	}
}
*/


import (
	"context"

	"github.com/docker/docker/client"
)

type (
	Dash struct {
		Command      []string // 命令
		DockerDasher *D       // 操作
	}
	// D 包含 Docker 操作的方法
	D struct {
		Image   string          // 镜像名
		Ports   []string        // 映射端口
		ConName string          // 容器名
		ConID   string          // 容器id
		Client  *client.Client  // 连接
		Context context.Context // 上下文
	}
)

type Dasher interface {
	Start() (err error)
	Stop() (err error)
	Delete() (err error)
	Create(command []string) (string, error)
}

func NewDockerDash(image string, opts ...client.Opt) (*D, error) {

	if len(opts) == 0 {
		opts = append(opts, client.WithHost("tcp://localhost:2375"))
	}

	cli, err := client.NewClientWithOpts(opts...)
	if err != nil {
		return nil, err
	}

	return &D{
		Image:   image,
		Client:  cli,
		Context: context.Background(),
	}, nil
}
