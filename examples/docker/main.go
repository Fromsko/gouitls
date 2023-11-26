package main

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

func createMySQLContainer() {
	// 创建 Docker 客户端实例
	cli, err := client.NewClientWithOpts(client.WithHost("unix:///var/run/docker.sock"))
	if err != nil {
		panic(err)
	}

	// 创建一个自定义的网络
	ctx := context.Background()
	_, err = cli.NetworkCreate(ctx, "MicroService", types.NetworkCreate{
		Driver: "bridge",
		IPAM: &network.IPAM{
			Config: []network.IPAMConfig{
				{
					Subnet:  "172.18.0.0/16", // 使用不同的子网
					Gateway: "172.18.0.1",
				},
			},
		},
	})
	if err != nil {
		panic(err)
	}

	// 创建一个新的 MySQL 容器并连接到自定义网络
	containerConfig := &container.Config{
		Image: "mysql:latest",
		Env: []string{
			"MYSQL_ROOT_PASSWORD=root",
			"MYSQL_PORT=3306",
		},
	}

	hostConfig := &container.HostConfig{
		NetworkMode: "MicroService", // 使用自定义网络
		PortBindings: map[nat.Port][]nat.PortBinding{
			"3306/tcp": {
				{HostIP: "0.0.0.0", HostPort: "3700"},
			},
		},
	}

	resp, err := cli.ContainerCreate(
		ctx, containerConfig, hostConfig, nil, nil, "api-mysql")
	if err != nil {
		panic(err)
	}

	// 启动新创建的 MySQL 容器
	err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		panic(err)
	}

	// 等待容器启动完成
	time.Sleep(10 * time.Second)

	// 输出容器信息
	containerInfo, err := cli.ContainerInspect(ctx, resp.ID)
	if err != nil {
		panic(err)
	}
	fmt.Printf("MySQL container created with ID: %s and IP: %s\n", containerInfo.ID[:12], "172.18.0.5")
}

func main() {
	createMySQLContainer()
}
