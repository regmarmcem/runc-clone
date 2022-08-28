package config

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regmarmcem/runc-clone/pkg/log"
	"regmarmcem/runc-clone/pkg/util"
	"strings"
	"syscall"

	l "log"

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
	fd       int
}

func NewOpts(command string, uid uint32, mountDir string) (_ *ContainerOpts, sockets [2]int) {
	argv := strings.Split(command, " ")
	sockets, err := util.GenerateSocketPair()

	if err != nil {
		log.Logger.Infof("Failed to generate socke pair. %s", err)
	}

	return &ContainerOpts{
		argv:     argv[1:],
		path:     argv[0],
		uid:      uid,
		mountDir: mountDir,
		fd:       sockets[1],
	}, sockets
}

type Container struct {
	sockets      [2]int
	config       ContainerOpts
	childProcess *exec.Cmd
}

func NewContainer(ctx *cli.Context) *Container {
	config, sockets := NewOpts(ctx.String("command"), uint32(ctx.Int("uid")), ctx.Path("mount"))
	// TODO to pass ipc.go
	// sender := os.NewFile(uintptr(config.fd), "")

	return &Container{
		config:  *config,
		sockets: sockets,
	}
}

func (c *Container) create() (err error) {
	fmt.Println("Create finished")
	cmd, err := ChildProcess(c.config)
	if err != nil {
		log.Logger.Infof("Unable to create child process %s", err)
		return err
	}
	c.setProcess(cmd)
	log.Logger.Debug("Creation finished")
	return nil
}

func (c *Container) setProcess(cmd *exec.Cmd) {
	c.childProcess = cmd
}

func (c Container) cleanExit() (err error) {
	fmt.Printf("Cleaning container")
	if err := syscall.Close(c.sockets[0]); err != nil {
		log.Logger.Infof("Unable to close write socket %s", err)
		return err
	}

	if err := syscall.Close(c.sockets[1]); err != nil {
		log.Logger.Infof("Unable to close read socket %s", err)
		return err
	}
	return nil
}

func Start(ctx *cli.Context) {

	if err := log.InitLogger(ctx.Bool("debug")); err != nil {
		l.Fatal(err)
	}
	err := supported()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	log.Logger.Info("Architecture is supported")
	c := NewContainer(ctx)
	fmt.Printf("Container is %v\n", c)
	if err = c.create(); err != nil {
		log.Logger.Infof("Unable to create child process %s", err)
		os.Exit(1)
	}
	log.Logger.Debug("Waiting child process")
	err = waitChild(c.childProcess)
	if err != nil {
		log.Logger.Infof("Wait child failed %s", err)
		os.Exit(1)
	}
	c.cleanExit()
}

func waitChild(cmd *exec.Cmd) (err error) {

	err = cmd.Wait()
	if err != nil {
		log.Logger.Infof("Unable to wait child process %s", err)
		return err
	}
	return nil
}

func supported() (err error) {
	u := syscall.Utsname{}
	err = syscall.Uname(&u)
	if err != nil {
		return err
	}
	var a string
	for _, val := range u.Machine {
		if val := rune(int(val)); val != rune(0) {
			a += string(val)
		}
	}

	if !(strings.Compare(a, "x86_64") == 0) {
		return errors.New("x86_64 is only supported architecture")
	}

	return nil
}
