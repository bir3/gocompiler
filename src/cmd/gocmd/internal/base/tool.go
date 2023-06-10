// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package base

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/bir3/gocompiler/src/go/build"

	"github.com/bir3/gocompiler/src/cmd/gocmd/internal/cfg"
)

// Tool returns the path to the named tool (for example, "vet").
// If the tool cannot be found, Tool exits the process.
func Tool(toolName string) string {
	toolPath := filepath.Join(build.ToolDir, toolName) + cfg.ToolExeSuffix()
	if len(cfg.BuildToolexec) > 0 {
		return toolPath
	}
	// Give a nice message if there is no tool with that name.

	s := "gocompiler:" + toolName + ":" + toolPath
	return s
}

func ToolCommand(exe string, args ...string) *exec.Cmd {
	var cmd *exec.Cmd
	if strings.HasPrefix(exe, "gocompiler:") {
		// format: "gocompiler:<tool>:<original-path-is-ignored>"
		s := exe[len("gocompiler:"):]
		k := strings.Index(s, ":")
		tool := s[0:k]
		exe, err := os.Executable()
		if err != nil {
			panic(fmt.Sprintf("ToolCommand[gocompiler] failed to get self executable: %s", err))
		}
		cmd = exec.Command(exe, args...)
		cmd.Env = append(cmd.Environ(), "GOCOMPILER_TOOL="+tool)
	} else {
		cmd = exec.Command(exe, args...)
	}
	return cmd
}
