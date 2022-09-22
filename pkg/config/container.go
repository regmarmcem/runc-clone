package config

import (
	"errors"
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
	uid      int
	MountDir string
	fd       *os.File
	Hostname string
}

func NewOpts(command string, uid int, mountDir string) (_ *ContainerOpts, sockets [2]*os.File) {
	argv := strings.Split(command, " ")
	sockets, err := util.GenerateSocketPair()

	if err != nil {
		log.Logger.Infof("Failed to generate socke pair. %s", err)
	}

	return &ContainerOpts{
		argv:     argv[1:],
		path:     argv[0],
		uid:      uid,
		MountDir: mountDir,
		fd:       sockets[1],
		Hostname: Hostname(),
	}, sockets
}

type Container struct {
	sockets      [2]*os.File
	config       ContainerOpts
	childProcess *exec.Cmd
}

func NewContainer(ctx *cli.Context) *Container {
	config, sockets := NewOpts(ctx.String("command"), ctx.Int("uid"), ctx.Path("mount"))
	// TODO to pass ipc.go
	// sender := os.NewFile(uintptr(config.fd), "")

	return &Container{
		config:  *config,
		sockets: sockets,
	}
}

func (c *Container) create() (err error) {
	log.Logger.Debugf("c.config.fd.Name() is %s", c.config.fd.Name())
	cmd, err := ExecProcess(&c.config)
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

func (c *Container) cleanExit() (err error) {
	log.Logger.Debug("Exiting...")

	if err := c.sockets[0].Close(); err != nil {
		log.Logger.Infof("Unable to close write socket %s", err)
		return err
	}

	if err := CleanMounts(c.config.MountDir); err != nil {
		log.Logger.Infof("Unable to clean mounts %s", err)
		return err
	}
	return nil
}

func Start(ctx *cli.Context) {
	c := NewContainer(ctx)
	cmd, err := ChildProcess(&c.config)
	if err != nil {
		l.Fatal(err)
	}
	err = supported()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	log.Logger.Info("Architecture is supported")
	err = waitChild(cmd)
	if err != nil {
		log.Logger.Infof("Wait child failed %s", err)
		// os.Exit(1)
	}
}

func Initialize(ctx *cli.Context) {

	c := NewContainer(ctx)
	log.Logger.Infof("Contaienr is %v\n", c)
	if err := c.create(); err != nil {
		log.Logger.Infof("Unable to create child process %s", err)
		os.Exit(1)
	}
	log.Logger.Debug("runc-clone initialize method")
	if err := waitChild(c.childProcess); err != nil {
		log.Logger.Infof("Wait child failed %s", err)
		// os.Exit(1)
	}

	log.Logger.Infof("calling handlechilduidmap: %s", c.sockets[0])
	log.Logger.Infof("c.childProcess.Process.Pid: %s", c.childProcess)
	HandleChildUidMap(c.childProcess.Process.Pid, c.sockets[0])
	log.Logger.Debug("Creation finished")
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
