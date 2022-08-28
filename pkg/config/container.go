package config

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
)

var (
	Contaiener *Container
)

type ContainerOpts struct {
	path     string
	argv     []string
	uid      uint32
	mountDir string
}

func NewOpts(command string, uid uint32, mountDir string) *ContainerOpts {
	argv := strings.Split(command, " ")
	return &ContainerOpts{
		argv:     argv,
		path:     argv[0],
		uid:      uid,
		mountDir: mountDir,
	}
}

type Container struct {
	config ContainerOpts
}

func NewContainer(ctx *cli.Context) *Container {
	config := NewOpts(ctx.String("command"), uint32(ctx.Int("uid")), ctx.Path("mount"))
	return &Container{config: *config}
}

func (c Container) create() {
	fmt.Printf("Create finished")
}

func (c Container) cleanExit() {
	fmt.Printf("Cleaning container")
}

func Start(ctx *cli.Context) {
	Container := NewContainer(ctx)
	fmt.Printf("Container is %v", Container)
	Container.create()
	Container.cleanExit()
}
