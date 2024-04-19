package exec

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type process struct {
	c *Cmd
}
type processState struct {
	c *Cmd
}

type Cmd struct {
	// Path is the path of the command to run.
	//
	// This is the only field that must be set to a non-zero
	// value. If Path is relative, it is evaluated relative
	// to Dir.
	Path string

	// Args holds command line arguments, including the command as Args[0].
	// If the Args field is empty or nil, Run uses {Path}.
	//
	// In typical use, both Path and Args are set by calling Command.
	Dir    string
	Args   []string
	Env    []string
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer

	Cancel       func() error
	WaitDelay    time.Duration
	Process      process
	ProcessState *processState

	realCmd *exec.Cmd

	isCommand        bool
	isCommandContext bool
	ctx              context.Context
	name             string
	arg              []string

	isTool bool // temp debug
}

// type Cmd = exec.Cmd
type ExitError = exec.ExitError
type Signal = os.Signal

var ErrWaitDelay error = exec.ErrWaitDelay

func log(msg string) {
	f, err := os.OpenFile("/tmp/bir3.log", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		panic("exit")
	}
	fmt.Fprintf(f, "%s\n", msg)
	f.Close()

}
func stackTrace() (string, []string) {
	b := make([]byte, 9900) // adjust buffer size to be larger than expected stack
	n := runtime.Stack(b, false)
	s := string(b[:n])
	fstack := []string{}
	lines := strings.Split(s, "\n")
	for _, line := range lines {
		line := strings.TrimSpace(line)
		if strings.HasPrefix(line, "/Volumes") {
			i := strings.Index(line, "/src/")
			//k := strings.Index(line, ":")
			if i > 0 {
				f := line[i+5:]
				fstack = append(fstack, f)
				if len(fstack) > 4 {
					break
				}
			}
		}
	}
	return s, fstack
}

func (p process) Signal(sig Signal) error {
	return p.c.realCmd.Process.Signal(sig)
}
func (p process) Kill() error {
	return p.c.realCmd.Process.Kill()
}

func (p processState) Success() bool {
	return p.c.realCmd.ProcessState.Success()
}
func (p processState) UserTime() time.Duration {
	return p.c.realCmd.ProcessState.UserTime()
}
func (p processState) SystemTime() time.Duration {
	return p.c.realCmd.ProcessState.SystemTime()
}

func Command(name string, arg ...string) *Cmd {
	log("x1")
	if len(arg) > 1 {
		log(fmt.Sprintf("### bir3.exec, arg[0]=%s arg[1]=%s", arg[0], arg[1]))
		//panic("exit")
	}
	//xcmd := exec.Command(name, arg...)
	return &Cmd{name: name, arg: arg, isCommand: true,
		Path: name,
		Args: append([]string{name}, arg...),
	}
}

func CommandContext(ctx context.Context, name string, arg ...string) *Cmd {
	log("x2")
	if len(arg) > 1 {
		log(fmt.Sprintf("### bir3.exec, arg[0]=%s arg[1]=%s", arg[0], arg[1]))
	}
	return &Cmd{ctx: ctx, name: name, arg: arg, isCommandContext: true,
		Path: name,
		Args: append([]string{name}, arg...),
	}

}

func isTool(file string) string {
	// input: $GOROOT/pkg/tool/darwin_arm64/compile
	// - assume GOROOT is absolute path
	if strings.HasPrefix(file, os.Getenv("GOROOT")) {
		dir_os_arch := filepath.Dir(file)     // pkg/tool/darwin_arm64
		dir_tool := filepath.Dir(dir_os_arch) // pkg/tool

		if filepath.Base(dir_tool) == "tool" {
			return filepath.Base(file)
		}
	}
	return ""
}

func LookPath(file string) (string, error) {
	if isTool(file) != "" {
		return file, nil
	}
	return exec.LookPath(file)
}

func (c *Cmd) Run() error {
	c.mirror()
	log(fmt.Sprintf("isTool=%t Run with Args[0]=%s Path=%s GOROOT=%s", c.isTool, c.realCmd.Args[0], c.Path, os.Getenv("GOROOT")))
	if c.isTool {
		s, fstack := stackTrace()
		if len(fstack) > 0 {
			log(fmt.Sprintf("who: fstack[0]=%s\n stack=%s\n", fstack[0], s))
		}
		///log("")
		//panic("who is it")
	}
	err := c.realCmd.Run()
	log(fmt.Sprintf("- Run result=%s", err))
	return err
}
func (c *Cmd) Environ() []string {
	c.mirror()
	return c.realCmd.Environ()
}
func (c *Cmd) mirror() {
	if c.isCommand {
		c.realCmd = exec.Command(c.name, c.arg...)
	} else if c.isCommandContext {
		c.realCmd = exec.CommandContext(c.ctx, c.name, c.arg...)
	} else {
		// manually constructed by the user
		c.realCmd = &exec.Cmd{}
		c.realCmd.Path = c.Path
		c.realCmd.Args = c.Args
	}

	c.realCmd.Dir = c.Dir
	c.realCmd.Env = c.Env
	c.realCmd.Stdin = c.Stdin
	c.realCmd.Stdout = c.Stdout
	c.realCmd.Stderr = c.Stderr
	c.realCmd.Cancel = c.Cancel
	c.realCmd.WaitDelay = c.WaitDelay
	c.Process.c = c
	if c.ProcessState == nil {
		c.ProcessState = &processState{}
	}
	c.ProcessState.c = c

	// fork/exec $GOROOT/pkg/tool/darwin_arm64/compile: no such file or directory
	tool := isTool(c.Path)

	if tool != "" {
		// BUG: could add multiple times since shared var
		c.realCmd.Env = append(c.realCmd.Env, fmt.Sprintf("BIR3_GOCOMPILER_TOOL=%s", tool))
		c.realCmd.Path, _ = os.Executable()
		c.realCmd.Args[0], _ = os.Executable()
		c.isTool = true
	}
}

func (c *Cmd) Start() error {
	c.mirror()
	log(fmt.Sprintf("isTool=%t Start with Args[0]=%s", c.isTool, c.realCmd.Args[0]))
	return c.realCmd.Start()
}
func (c *Cmd) Wait() error {
	c.mirror()
	return c.realCmd.Wait()
}
func (c *Cmd) CombinedOutput() ([]byte, error) {
	log(fmt.Sprintf("CombinedOutput with Args[0]=%s", c.realCmd.Args[0]))
	c.mirror()
	return c.realCmd.CombinedOutput()

}
func (c *Cmd) Output() ([]byte, error) {
	log(fmt.Sprintf("Output with Args[0]=%s", c.realCmd.Args[0]))
	c.mirror()
	return c.realCmd.Output()

}
func (c *Cmd) StdoutPipe() (io.ReadCloser, error) {
	c.mirror()
	return c.realCmd.StdoutPipe()
}
func (c *Cmd) StdinPipe() (io.WriteCloser, error) {
	c.mirror()
	return c.realCmd.StdinPipe()
}
func (c *Cmd) StderrPipe() (io.ReadCloser, error) {
	c.mirror()
	return c.realCmd.StderrPipe()
}

/*

# github.com/bir3/gocompiler/src/cmd/cgo
../gocompiler/src/cmd/cgo/util.go:64:25: undefined: exec.ExitError
# github.com/bir3/gocompiler/src/cmd/gocmd/internal/cache
../gocompiler/src/cmd/gocmd/internal/cache/prog.go:164:14: undefined: exec.CommandContext
# github.com/bir3/gocompiler/src/cmd/gocmd/internal/tool
../gocompiler/src/cmd/gocmd/internal/tool/tool.go:128:26: undefined: exec.ExitError
# github.com/bir3/gocompiler/src/cmd/gocmd/internal/vcs
../gocompiler/src/cmd/gocmd/internal/vcs/vcs.go:699:28: undefined: exec.ExitError
# github.com/bir3/gocompiler/src/cmd/gocmd/internal/modfetch/codehost
../gocompiler/src/cmd/gocmd/internal/modfetch/codehost/codehost.go:363:12: undefined: exec.CommandContext
../gocompiler/src/cmd/gocmd/internal/modfetch/codehost/git.go:844:42: undefined: exec.ExitError
# github.com/bir3/gocompiler/src/cmd/compile/internal/importer
../gocompiler/src/cmd/compile/internal/importer/gcimporter.go:50:30: undefined: exec.ExitError


*/
