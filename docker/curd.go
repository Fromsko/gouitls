package docker

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

func (dash *D) Create(command []string) (string, error) {

	resp, err := dash.Client.ContainerCreate(
		dash.Context,
		&container.Config{
			Image: dash.Image,
			Cmd:   command,
			Tty:   true,
		},
		nil,
		nil,
		nil,
		dash.ConName,
	)
	if err != nil {
		return "", err
	} else {
		dash.ConID = resp.ID
	}

	return "", nil
}

func (dash *D) Start() (err error) {

	if err = dash.Client.ContainerStart(
		dash.Context, dash.ConID,
		types.ContainerStartOptions{},
	); err != nil {
		return err
	}

	return nil
}

func (dash *D) Delete() (err error) {

	if err = dash.Client.ContainerRemove(
		dash.Context, dash.ConID,
		types.ContainerRemoveOptions{},
	); err != nil {
		return
	}

	return nil
}

func (dash *D) Stop() (err error) {

	if err = dash.Client.ContainerStop(
		dash.Context, dash.ConID,
		container.StopOptions{},
	); err != nil {
		return
	}

	return nil
}
