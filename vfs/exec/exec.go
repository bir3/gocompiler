// Copyright 2022 Bergur Ragnarsson.
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package exec

import (
	"context"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/bir3/gocompiler/vfs"
)

type ExitError = exec.ExitError
type Error = exec.Error

var ErrNotFound = exec.ErrNotFound

type Cmd struct {
	osCmd exec.Cmd

	Path   string
	Args   []string
	Env    []string
	Dir    string
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer

	Process *os.Process // after Start

	// info about exited process, available after a call to Wait or Run:
	ProcessState *os.ProcessState

	ctx context.Context // nil means none
}

func CommandContext(ctx context.Context, name string, arg ...string) *Cmd {
	if ctx == nil {
		panic("nil Context")
	}
	osCmd := exec.CommandContext(ctx, name, arg...)
	cmd := Cmd{osCmd: *osCmd}
	return &cmd
}

func wrap(cmd *exec.Cmd) *Cmd {
	return &Cmd{osCmd: *cmd}
}

func copyAttr(c *Cmd) {
	c.osCmd.Path = c.Path
	c.osCmd.Args = c.Args
	c.osCmd.Env = c.Env
	c.osCmd.Dir = c.Dir
	c.osCmd.Stdin = c.Stdin
	c.osCmd.Stdout = c.Stdout
	c.osCmd.Stderr = c.Stderr

}
func addExitAttr(c *Cmd) {
	c.ProcessState = c.osCmd.ProcessState
	c.Process = c.osCmd.Process
}
func LookPath(file string) (string, error) {
	return exec.LookPath(file)
}
func Command(name string, arg ...string) *Cmd {
	osCmd := exec.Command(name, arg...)
	return &Cmd{osCmd: *osCmd, Path: osCmd.Path, Args: osCmd.Args}
}

func detectGoTool(c *Cmd) {
	if strings.HasPrefix(c.Path, vfs.GorootTool) {
		c.osCmd.Env = append(c.osCmd.Env, "GOCOMPILER_TOOL="+filepath.Base(c.Path))
		c.osCmd.Path = vfs.SharedExe
		c.osCmd.Args[0] = vfs.SharedExe
	}

}

func (c *Cmd) Run() error {
	copyAttr(c)
	detectGoTool(c)

	err := c.osCmd.Run()
	addExitAttr(c)
	return err
}
func (c *Cmd) Output() ([]byte, error) {
	// = Run
	copyAttr(c)
	detectGoTool(c)
	buf, err := c.osCmd.Output()
	addExitAttr(c)
	return buf, err
}
func (c *Cmd) CombinedOutput() ([]byte, error) {
	// = Run with combined output
	copyAttr(c)
	detectGoTool(c)
	buf, err := c.osCmd.CombinedOutput()
	addExitAttr(c)
	return buf, err
}

func (c *Cmd) Start() error {
	copyAttr(c)
	detectGoTool(c)
	err := c.osCmd.Start()
	c.Process = c.osCmd.Process
	return err
}

func (c *Cmd) Wait() error {
	err := c.osCmd.Wait()
	addExitAttr(c)
	return err
}

func (c *Cmd) StdinPipe() (io.WriteCloser, error) {
	return c.osCmd.StdinPipe()
}

func (c *Cmd) StderrPipe() (io.ReadCloser, error) {
	return c.osCmd.StderrPipe() // go vet
}

func (c *Cmd) StdoutPipe() (io.ReadCloser, error) {
	return c.osCmd.StdoutPipe() // go vet
}

//
// optional code by go version must be at end of file
// and in increasing go version order

//gocompiler: go1.19
func (c *Cmd) Environ() []string {
	return c.osCmd.Environ()
}
