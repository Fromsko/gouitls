package docker

import "time"

// PullAndRunAlpine 拉取并运行 Alpine 容器
func PullAndRunAlpine(d *Dash) error {
	hasAlpine, err := d.HasImage()
	if err != nil {
		return err
	}

	if !hasAlpine {
		if err := d.PullImage(d.DockerDasher.Image); err != nil {
			return err
		}
	}

	// 创建并启动 alpine 容器
	_, err = d.DockerDasher.Create(d.Command)
	if err != nil {
		return err
	} else {
		// 启动监控日志的协程
		go d.MonitorContainerLogs()
	}

	// 等待一段时间，模拟容器运行
	time.Sleep(10 * time.Second)

	// 停止和移除容器
	if err := d.DockerDasher.Stop(); err != nil {
		return err
	} else {
		if err := d.DockerDasher.Delete(); err != nil {
			return err
		}
	}

	return nil
}
