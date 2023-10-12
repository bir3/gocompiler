package gocompiler_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/bir3/gocompiler"
	"github.com/bir3/gocompiler/extra"
)

func TestMain(m *testing.M) {

	// the go toolchain is built into the executable and must be given a chance to run
	// => avoid side effects in init() as they will occur multiple times during compilation
	if gocompiler.IsRunToolchainRequest() {
		gocompiler.RunToolchain()
		return
	}

	os.Exit(m.Run())
}

func TestCompileStdin(t *testing.T) {
	t.Parallel()
	goSimple := `
	package main
	
	import "fmt"
	
	func main() {
			fmt.Printf("magic")
	}`

	dir := t.TempDir()

	err := os.WriteFile(filepath.Join(dir, "main.go"), []byte(goSimple), 0666)
	if err != nil {
		t.Fatalf("%s", err)
	}

	run := func(args ...string) string {
		cmd, err := gocompiler.Command(os.Environ(), args...)
		if err != nil {
			t.Fatalf("%s", err)
		}
		cmd.Dir = dir
		buf, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("%s", err)
		}
		s := string(buf)
		return s
	}
	run("go", "mod", "init", "abc")
	run("go", "build")

	cmd := exec.Command(filepath.Join(dir, "abc"))
	buf, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%s", err)
	}
	s := string(buf)
	exp := "magic"
	if s != exp {
		t.Fatalf("expected string %q but got %q", exp, s)
	}
}

func TestMkdirAllRace(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	dir = filepath.Join(dir, "a", "b", "c")
	err := extra.MkdirAllRace(dir, 0777)
	if err != nil {
		t.Fatalf("%s", err)
	}
	info, err := os.Stat(dir)
	if err != nil {
		t.Fatalf("%s", err)
	}
	if !info.IsDir() {
		t.Fatal("failed")
	}
}
