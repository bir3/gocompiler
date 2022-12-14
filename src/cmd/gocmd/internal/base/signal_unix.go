// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build unix || js

package base

import (
	       "github.com/bir3/gocompiler/vfs/os"
	"syscall"
)

var signalsToIgnore = []os.Signal{os.Interrupt, syscall.SIGQUIT}

// SignalTrace is the signal to send to make a Go program
// crash with a stack trace.
var SignalTrace os.Signal = syscall.SIGQUIT
