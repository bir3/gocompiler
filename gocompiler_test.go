package gocompiler_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

func compileAndRunString(t *testing.T, code string) string {
	dir := t.TempDir()

	err := os.WriteFile(filepath.Join(dir, "main.go"), []byte(code), 0666)
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
		s := string(buf)
		if err != nil {
			cmd := strings.Join(args, " ")
			t.Fatalf("cmd %s failed with %s - %s", cmd, err, s)
		}
		return s
	}
	run("go", "mod", "init", "abc")
	run("go", "build")
	cmd := exec.Command(filepath.Join(dir, "abc"))
	buf, err := cmd.CombinedOutput()
	s := string(buf)
	if err != nil {
		t.Fatalf("failed to run - %s - %s", s, err)
	}

	return s
}

func TestCompile(t *testing.T) {
	t.Parallel()
	goCode := `
	package main
	
	import "fmt"
	
	func main() {
			fmt.Printf("magic")
	}`
	output := compileAndRunString(t, goCode)
	exp := "magic"
	if output != exp {
		t.Fatalf("expected string %q but got %q", exp, output)
	}
}

func TestCompileWithC(t *testing.T) {
	t.Parallel()

	goCode := `
	package main
	
	// typedef int (*intFunc) ();
	//
	// int
	// bridge_int_func(intFunc f)
	// {
	//		return f();
	// }
	//
	// int fortytwo()
	// {
	//	    return 42;
	// }
	import "C"
	import "fmt"
	
	func main() {
		f := C.intFunc(C.fortytwo)
		fmt.Println(int(C.bridge_int_func(f)))
		// Output: 42
	}
	
	`
	output := compileAndRunString(t, goCode)
	exp := "42\n"
	if output != exp {
		t.Fatalf("expected string %q but got %q", exp, output)
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
