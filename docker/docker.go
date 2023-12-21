package docker

import (
	"context"
	"io"
	"os"
	"strings"
	"time"

	"github.com/Fromsko/gouitls/logs"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

var log = logs.InitLogger()

// Option 结构体用于配置 Docker 容器创建的选项
type Option struct {
	NetworkName    string            // 自定义网络的名称
	IPAddress      string            // 指定IP地址
	Subnet         string            // 子网配置
	Gateway        string            // 网关配置
	Image          string            // Docker 镜像名称
	Command        []string          // 命令参数
	Env            []string          // 容器环境变量
	PortMappings   map[string]string // 端口映射配置
	ContainerName  string            // 容器名称
	VolumeMappings []VolumeMapping   // 数据映射
}

// 添加 VolumeMapping 结构体用于配置数据卷映射
type VolumeMapping struct {
	Source   string // 本地主机路径
	Target   string // 容器内路径
	ReadOnly bool   // 是否为只读
}

// AutoTask 结构体封装了创建 Docker 容器的通用操作
type AutoTask struct {
	cli *client.Client
	ctx context.Context
	opt Option
}

// NewAutoTask 创建一个 AutoTask 结构体的实例
func NewAutoTask(opt Option, ops ...client.Opt) *AutoTask {
	if len(ops) == 0 {
		ops = append(ops, client.WithHost("unix:///var/run/docker.sock"))
	}

	cli, err := client.NewClientWithOpts(ops...)
	if err != nil {
		panic(err)
	}

	return &AutoTask{
		cli: cli,
		opt: opt,
		ctx: context.Background(),
	}
}

// WithNetwork 设置自定义网络的名称
func (at *AutoTask) WithNetwork(networkName string) *AutoTask {
	at.opt.NetworkName = networkName
	return at
}

// WithSubnet 设置子网配置
func (at *AutoTask) WithSubnet(subnet, gateway string) *AutoTask {
	at.opt.Subnet = subnet
	at.opt.Gateway = gateway
	return at
}

// WithImage 设置 Docker 镜像名称
func (at *AutoTask) WithImage(image string) *AutoTask {
	at.opt.Image = image
	return at
}

// WithEnv 设置容器环境变量
func (at *AutoTask) WithEnv(env []string) *AutoTask {
	at.opt.Env = env
	return at
}

// WithPortMappings 设置端口映射配置
func (at *AutoTask) WithPortMappings(portMappings map[string]string) *AutoTask {
	at.opt.PortMappings = portMappings
	return at
}

// WithContainerName 设置容器名称
func (at *AutoTask) WithContainerName(containerName string) *AutoTask {
	at.opt.ContainerName = containerName
	return at
}

// WithStaticIp 设置容器静态IP
func (at *AutoTask) WithStaticIp(IPAddress string) *AutoTask {
	at.opt.IPAddress = IPAddress
	return at
}

func (at *AutoTask) WithCommand(Command string) *AutoTask {
	at.opt.Command = Command
	return at
}

// 添加 WithVolumeMappings 方法用于设置数据卷映射
func (at *AutoTask) WithVolumeMappings(volumeMappings []VolumeMapping) *AutoTask {
	at.opt.VolumeMappings = volumeMappings
	return at
}

// ContainerExists 检测是否存在指定容器名称
func (at *AutoTask) ContainerExists(callBack func()) {
	flag := true

	listOptions := types.ContainerListOptions{
		All: true, // 包括停止的容器
	}

	containers, err := at.cli.ContainerList(at.ctx, listOptions)
	if err != nil {
		log.Errorf("无法获取容器列表：%v", err) // 发生错误时返回 true
	} else {
		for _, container := range containers {
			for _, name := range container.Names {
				if strings.TrimLeft(name, "/") == at.opt.ContainerName {
					log.Warnf("容器存在 => %v", strings.Split(container.Names[0], "/")[1])
					log.Warnf("容器状态 => %v", container.Status)
					log.Warnf("容器ID => %v", container.ID)
					flag = false // 存在指定名称的容器，返回 false
				}
			}
		}
	}
	if flag {
		callBack() // 不存在指定名称的容器，返回 true
	}
}

func (at *AutoTask) checkNetwork() bool {
	// 列出当前系统中的所有网络
	networkList, err := at.cli.NetworkList(at.ctx, types.NetworkListOptions{})
	if err != nil {
		panic(err)
	}

	// 检查是否存在具有相同名称的网络
	networkName := at.opt.NetworkName
	networkExists := false
	for _, network := range networkList {
		if network.Name == networkName {
			networkExists = true
			break
		}
	}

	return networkExists
}

func (at *AutoTask) EchoError(name string, err error, flag bool) {
	if err == nil {
		return
	}
	if !flag {
		log.Warnf("%s %v\n", name, err)
		return
	}
	log.Errorf("%s %v", name, err)
	os.Exit(0)
}

// PullImage 拉取容器
func (at *AutoTask) PullImage(imgName string) error {
	if !strings.Contains(imgName, ":") {
		imgName = imgName + ":latest"
	}

	out, err := at.cli.ImagePull(
		at.ctx, imgName,
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

func (at *AutoTask) Start() (err error) {
	if err = at.cli.ContainerStart(
		at.ctx, at.opt.ContainerName,
		types.ContainerStartOptions{},
	); err != nil {
		return err
	}

	return nil
}

func (at *AutoTask) Delete() (err error) {
	if err = at.cli.ContainerRemove(
		at.ctx, at.opt.ContainerName,
		types.ContainerRemoveOptions{},
	); err != nil {
		return
	}

	return nil
}

func (at *AutoTask) Stop() (err error) {
	if err = at.cli.ContainerStop(
		at.ctx, at.opt.ContainerName,
		container.StopOptions{},
	); err != nil {
		return
	}

	return nil
}

// Run 创建并运行 Docker 容器
func (at *AutoTask) Run() {
	// 如果不存在具有相同名称的网络，则创建新网络
	if !at.checkNetwork() {
		_, err := at.cli.NetworkCreate(at.ctx, at.opt.NetworkName, types.NetworkCreate{
			Driver: "bridge",
			// 设置静态地址

			IPAM: &network.IPAM{
				Config: []network.IPAMConfig{
					{
						Subnet:  at.opt.Subnet,
						Gateway: at.opt.Gateway,
					},
				},
			},
		})
		at.EchoError("创建网络失败: ", err, false)
	}

	// 创建一个新的容器并连接到自定义网络
	containerConfig := &container.Config{
		Image: at.opt.Image,
		Env:   at.opt.Env,
		Cmd:   at.opt.Command,
	}

	// 创建容器的 VolumeConfig 配置
	volumeConfigs := []mount.Mount{}
	for _, volMapping := range at.opt.VolumeMappings {
		volumeConfigs = append(volumeConfigs, mount.Mount{
			Type:     mount.TypeBind,
			Source:   volMapping.Source,
			Target:   volMapping.Target,
			ReadOnly: volMapping.ReadOnly,
		})
	}

	// 主机设置
	hostConfig := &container.HostConfig{
		NetworkMode:  container.NetworkMode(at.opt.NetworkName), // 使用自定义网络
		PortBindings: nat.PortMap{},
		Mounts:       volumeConfigs, // 添加数据卷映射配置
	}

	// 网络配置
	networkConfig := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			at.opt.NetworkName: {
				IPAMConfig: &network.EndpointIPAMConfig{
					IPv4Address: at.opt.IPAddress,
				},
			},
		},
	}

	for containerPort, hostPort := range at.opt.PortMappings {
		hostConfig.PortBindings[nat.Port(containerPort)] = []nat.PortBinding{
			{HostIP: "0.0.0.0", HostPort: hostPort},
		}
	}

	log.Infof("正在创建容器 => %v", at.opt.ContainerName)
	resp, err := at.cli.ContainerCreate(
		at.ctx, containerConfig,
		hostConfig, networkConfig, nil,
		at.opt.ContainerName,
	)
	at.EchoError("创建容器失败: ", err, true)

	// 启动新创建的容器
	err = at.cli.ContainerStart(at.ctx, resp.ID, types.ContainerStartOptions{})
	at.EchoError("启动容器失败: ", err, true)

	// 等待容器启动完成
	time.Sleep(10 * time.Second)

	// 输出容器信息
	containerInfo, err := at.cli.ContainerInspect(at.ctx, resp.ID)
	at.EchoError("容器网络信息获取失败: ", err, true)

	// 获取容器的自定义网络设置
	networkSettings := containerInfo.NetworkSettings.Networks[at.opt.NetworkName]

	// 检查容器是否分配了 IP 地址
	if networkSettings != nil {
		ipAddress := networkSettings.IPAddress
		log.Infof("容器ID => %s", containerInfo.ID[:12])
		log.Infof("IP地址 => %s", ipAddress)
	} else {
		log.Infof("容器ID => %s\n", containerInfo.ID[:12])
		log.Infof("未分配IP地址-使用容器网络 => %s", at.opt.NetworkName)
	}
	defer func() {
		log.Infof("容器创建成功 => %v", containerInfo.State.Status)
	}()
}

/**
func main() {
	// 链式调用构造 AutoTask 结构体
	at := NewAutoTask(Option{}).
		WithImage("mysql:latest").
		WithNetwork("online_MicroService").
		WithContainerName("api-mysql").
		WithStaticIp("172.20.0.10").
		WithSubnet("172.20.0.0/16", "172.20.0.1").
		WithEnv([]string{"MYSQL_ROOT_PASSWORD=root", "MYSQL_PORT=3306"}).
		WithPortMappings(map[string]string{
			"3306/tcp": "3306",
		})

	if at.ContainerExists(at.opt.ContainerName) {
		// 运行 AutoTask
		at.Run()
	}

	// 用户服务
	userAPI := NewAutoTask(Option{}).
		WithImage("service-user:latest").
		WithContainerName("api-user").
		WithNetwork("online_MicroService").
		WithStaticIp("172.20.0.21").
		WithSubnet("172.20.0.0/16", "172.20.0.1").
		WithPortMappings(map[string]string{
			"8000/tcp": "8001",
			"9000/tcp": "9001",
		}).
		WithVolumeMappings([]VolumeMapping{
			{
				Source:   "/home/skong/manage/app/user/configs/",
				Target:   "/data/conf",
				ReadOnly: false, // 根据需要设置为只读或读写
			},
		})

	if userAPI.ContainerExists(userAPI.opt.ContainerName) {
		userAPI.Run()
	}
}
**/
