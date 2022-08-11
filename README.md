
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
	// => avoid side effects in init() as they will occur multiple times during compilation
	if gocompiler.IsRunToolchainRequest() {
		gocompiler.RunToolchain()
		return
	}

    err := os.WriteFile("temp.go", []byte(goCodeStr), 0666)
	if err != nil {
		log.Fatal(err)
	}

	result, err := gocompiler.Run("go", "run", "temp.go")
    fmt.Fprintf(os.Stdout, result.Stdout)
    fmt.Fprintf(os.Stderr, result.Stderr)
	if err != nil {
		log.Fatal(err)
	}

	os.Remove("temp.go")
}
```


# Limitations


- match your go compiler version with the library go compiler version
- avoid side effects in init() or global `var` initialization

Your executable will serve two purposes: 
- run your application
- run the Go compiler toolchain via `gocompiler.RunToolchain()`

This means that you should avoid side-effects like opening a http port or connecting to a database as they will happen say 100 times when the Go toolchain repeatedly invokes the executable during compilation. 


# gocompiler as a library vs. the official go toolchain

|                 | "github.com/bir3/gocompiler" | official go toolchain | 
| --------------  | ---------------------------- | ------- |
| Number of files | 1                            | 12006   |
| Total size      | 64 MB                        | 462 MB  |
| go build -a     | 7.3 sec                      | 6.9 sec |

Note that this package is only focused on compiling go source code into an executable, while the official go toolchain provides many more tools.