
# gocompiler

The Go compiler as a package

```bash
# go1.22.0
go get github.com/bir3/gocompiler@v0.9.2202
```


# Example


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

	err := os.WriteFile("temp.go", []byte(goCodeStr), 0644)
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

# Standard library

The standard library is embedded and is extracted on first run
to standard config location;
- `$HOME/.config/bir3-gocompiler` (linux)
- `$HOME/Library/Application Support/bir3-gocompiler` (mac)

# Private Go build cache

To avoid interfering with the standard Go toolchain build cache, the package has a private 
Go build cache;
- `$HOME/.cache/bir3-gocompiler` (linux)
- `$HOME/Library/Caches/bir3-gocompiler` (mac)

# Limitations

- avoid side effects in init() and global variable initializations

Reason: Your executable will serve two purposes: 
- run your application
- run the Go compiler toolchain via `gocompiler.RunToolchain()`

Side effects in init() and global variable initializations occur every time the executable is started.  
The embedded Go toolchain will repeatedly start the executable during compilation to compile Go source code.  
This means that global side effects like opening a http port, writing to a file or connecting to a database is likely to cause problems.


# gocompiler as a package vs. the official Go toolchain

|                      | "github.com/bir3/gocompiler"  | official go toolchain |                           |
| -------------------  | ----------------------------- | --------------------- | ------------------------- |
| Download size        | 26 MB (gzip of executable)    | 62 MB (gzip tarfile)  |                           |
| Size on disk         | 91 MB                         | 237 MB                |                           |
| Compile speed        | 12.9 sec                      | 12.4 sec              | macbook M1, `go build -a` |

Note that this package is only focused on compiling Go source code into an executable, while the official Go toolchain provides many more tools.

# Acknowledgments

* The Go Authors. https://github.com/golang/go 
* Klaus Post, zstd decompression: https://github.com/klauspost/compress