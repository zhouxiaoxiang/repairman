package repairman

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/dixonwille/wmenu/v5"
	"github.com/docker/docker/api/types"
	"github.com/pieterclaerhout/go-log"
)

const (
	DIR       = "BOOT-INF"
	FILETYPE  = "-type f -print"
	LIST      = "/tmp/class.list"
	JAR       = "jar uf0"
	JARFILE   = "/app.jar"
	EXTENSION = ".tar.gz"
	PACKAGE   = "client-*-*[0-9]"
	COPY      = "cp -fr"
	EXTRACT   = "tar -zxf"
	PROMPT    = "选择要修改的微服务容器:"
	NGINX     = "nodejs-server"
	TARGET    = "/gmpcloud/client/"
)

var (
	ctx = context.Background()
	c   = NewContainer()
)

type Repairman struct {
	cs []types.Container
}

func NewRepairman() *Repairman {
	cs, err := c.List(ctx)
	if err != nil {
		panic("No containers")
	}
	return &Repairman{cs: cs}
}

func (r *Repairman) RepairWeb() {
	fullName := ""
	baseName := ""

	fullNames, err := filepath.Glob(PACKAGE + EXTENSION)
	log.CheckError(err)

	baseNames, err := filepath.Glob(PACKAGE)
	log.CheckError(err)

	if len(baseNames) != 0 {
		baseName = baseNames[0]
	}

	if len(fullNames) != 0 {
		fullName = fullNames[0]
		cmd := exec.Command("sh", "-c", fmt.Sprintf("%s %s", EXTRACT, fullName))
		err = cmd.Run()
		log.CheckError(err)
		baseName = strings.TrimSuffix(fullName, EXTENSION)
	}

	_, err = os.Stat(baseName)
	if err != nil {
		return
	}

	log.Info("repair web: " + baseName)
	cmd := exec.Command("sh", "-c", fmt.Sprintf("%s %s/* %s", COPY, baseName, TARGET))
	err = cmd.Run()
	log.CheckError(err)

	c.Restart(ctx, NGINX)
}

func (r *Repairman) RepairJar() {
	_, err := os.Stat(DIR)
	if err != nil {
		return
	}

	log.Info("repair jar")
	mm := mainMenu(r.cs)
	err = mm.Run()
	log.CheckError(err)
}

func mainMenu(cs []types.Container) *wmenu.Menu {
	menu := wmenu.NewMenu(PROMPT)
	for _, container := range cs {
		menu.Option(container.Names[0][1:]+"["+container.State+"]", container.ID, false, nil)
	}

	menu.Action(func(opts []wmenu.Opt) error {
		if opts[0].Value == nil {
			return nil
		}

		id := opts[0].Value.(string)
		c.Copy(ctx, DIR, "/", id)

		script := fmt.Sprintf("find %s %s >%s; %s %s @%s",
			DIR, FILETYPE, LIST, JAR, JARFILE, LIST)
		cmd := []string{"sh", "-c", script}
		result, _ := c.Exec(ctx, id, "", cmd, nil)
		log.Info(result)

		c.Restart(ctx, id)
		return nil
	})
	return menu
}
