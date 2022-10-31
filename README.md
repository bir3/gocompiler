
# gocompiler

The Go compiler as a package



# Example

```bash
# the go compiler version should be the same as the one inside the library, e.g. 1.18.2 in this case
go get github.com/bir3/gocompiler@v0.1.0-go.1.18.2
```

```go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bir3/gocompiler"
)

var goCodeStr string = `
package main

import "fmt"

func main() {
	fmt.Println("This code was compiled standalone")
}
`

func main() {

	// the go toolchain is built into the executable and must be given a chance to run
	// => avoid side effects in init() and global variable initialization
	//    as they will occur multiple times during compilation
	if gocompiler.IsRunToolchainRequest() {
		gocompiler.RunToolchain()
		return
	}

	err := os.WriteFile("temp.go", []byte(goCodeStr), 0666)
	if err != nil {
		log.Fatal(err)
	}

	result, err := gocompiler.Run("go", "run", "temp.go")
	fmt.Fprintf(os.Stdout, "%s", result.Stdout)
	fmt.Fprintf(os.Stderr, "%s", result.Stderr)
	if err != nil {
		log.Fatal(err)
	}
	os.Remove("temp.go")
}
```


# Limitations

- match your Go compiler version with the package Go compiler version
- avoid side effects in init() and global variable initializations

Your executable will serve two purposes: 
- run your application
- run the Go compiler toolchain via `gocompiler.RunToolchain()`

Side effects in init() and global variable initializations occur on every time the executable is started.  
The embedded Go toolchain will repeatedly start the executable during compilation to compile Go source code.  
This means that global side effects like opening a http port or connecting to a database is likely to cause problems.

## limitation : creating a log file in a init() function

The main function may write a few lines to the logfile, then when we compile code, the subprocesses
that are also hosted in the main executable will also open and possibly write or truncate the logfile
creating confusion on why something as simple as writing to a logfile can fail to work !

* => avoid init() functions or limit them to initializing memory structures only

# gocompiler as a package vs. the official Go toolchain

|                 | "github.com/bir3/gocompiler" | official go toolchain | 
| --------------  | ---------------------------- | ------- |
| Number of files | 1                            | 12006   |
| Total size      | 64 MB                        | 462 MB  |
| go build -a     | 7.3 sec                      | 6.9 sec |

Note that this package is only focused on compiling Go source code into an executable, while the official Go toolchain provides many more tools.