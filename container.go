package repairman

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/system"
	"github.com/docker/go-connections/nat"
)

type Container struct {
	c *client.Client
}

func NewContainer() *Container {
	c, err := client.NewEnvClient()
	if err != nil {
		panic("Error: container")
	}

	return &Container{
		c: c,
	}
}

func (d Container) Pull(ctx context.Context, image, user, password string) (string, error) {
	authConfig := types.AuthConfig{
		Username: user,
		Password: password,
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		return "", err
	}

	authStr := base64.URLEncoding.EncodeToString(encodedJSON)
	reader, err := d.c.ImagePull(ctx, image, types.ImagePullOptions{RegistryAuth: authStr})
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(reader)
	return buf.String(), nil
}

func (d Container) Import(ctx context.Context, file string, image string) (string, error) {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return "", err
	}

	reader, err := d.c.ImageImport(ctx, types.ImageImportSource{Source: f, SourceName: "-"}, image, types.ImageImportOptions{})
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(reader)
	return buf.String(), nil
}

func (d Container) Run(ctx context.Context, name string, image string, cmd []string, volumes map[string]struct{}, ports nat.PortSet) error {
	resp, err := d.c.ContainerCreate(ctx, &container.Config{Image: image, Volumes: volumes, ExposedPorts: ports, Cmd: cmd}, nil, nil, nil, name)
	if err != nil {
		return err
	}
	if err = d.c.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}
	return nil
}

func (d Container) Copy(ctx context.Context, file string, dest string, container string) error {
	srcPath := file
	dstPath := dest
	dstInfo := archive.CopyInfo{Path: dstPath}
	dstStat, err := d.c.ContainerStatPath(ctx, container, dstPath)

	if err == nil && dstStat.Mode&os.ModeSymlink != 0 {
		linkTarget := dstStat.LinkTarget
		if !system.IsAbs(linkTarget) {
			dstParent, _ := archive.SplitPathDirEntry(dstPath)
			linkTarget = filepath.Join(dstParent, linkTarget)
		}

		dstInfo.Path = linkTarget
		dstStat, err = d.c.ContainerStatPath(ctx, container, linkTarget)
	}

	if err == nil {
		dstInfo.Exists, dstInfo.IsDir = true, dstStat.Mode.IsDir()
	}

	var (
		content         io.Reader
		resolvedDstPath string
	)

	srcInfo, err := archive.CopyInfoSourcePath(srcPath, true)
	if err != nil {
		return err
	}

	srcArchive, err := archive.TarResource(srcInfo)
	if err != nil {
		return err
	}
	defer srcArchive.Close()

	dstDir, preparedArchive, err := archive.PrepareArchiveCopy(srcArchive, srcInfo, dstInfo)
	if err != nil {
		return err
	}
	defer preparedArchive.Close()

	resolvedDstPath = dstDir
	content = preparedArchive

	options := types.CopyToContainerOptions{
		AllowOverwriteDirWithFile: true,
	}
	return d.c.CopyToContainer(ctx, container, resolvedDstPath, content, options)
}

func (d Container) Start(ctx context.Context, container string) error {
	err := d.c.ContainerStart(ctx, container, types.ContainerStartOptions{})
	return err
}

func (d Container) Stop(ctx context.Context, container string) error {
	timeout := time.Second * 5
	err := d.c.ContainerStop(ctx, container, &timeout)
	return err
}

func (d Container) Rm(ctx context.Context, container string, force bool) error {
	err := d.c.ContainerRemove(ctx, container, types.ContainerRemoveOptions{Force: force})
	return err
}

func (d Container) Push(ctx context.Context, image string, user string, password string) (string, error) {
	authConfig := types.AuthConfig{
		Username: user,
		Password: password,
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		return "", err
	}

	authStr := base64.URLEncoding.EncodeToString(encodedJSON)
	reader, err := d.c.ImagePush(ctx, image, types.ImagePushOptions{RegistryAuth: authStr})
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)
	return buf.String(), err
}

func (d Container) Exec(ctx context.Context, container string, chdir string, cmd []string, env []string) (string, error) {
	id, err := d.c.ContainerExecCreate(ctx, container, types.ExecConfig{Tty: true, WorkingDir: chdir, Cmd: cmd, Env: env, AttachStderr: true, AttachStdout: true})
	if err != nil {
		return "", err
	}
	resp, err := d.c.ContainerExecAttach(ctx, id.ID, types.ExecStartCheck{})
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Reader)
	return buf.String(), err
}

func (d Container) Restart(ctx context.Context, container string) error {
	_ = d.Stop(ctx, container)
	return d.Start(ctx, container)
}

func (d Container) IsRun(ctx context.Context, container string) bool {
	stat, err := d.c.ContainerInspect(ctx, container)
	if err != nil {
		return false
	}
	if !stat.State.Running {
		return false
	}
	return true
}

func (d Container) List(ctx context.Context) ([]types.Container, error) {
	return d.c.ContainerList(ctx, types.ContainerListOptions{
		All: true,
	})
}
