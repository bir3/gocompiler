// Code generated by go tool dist; DO NOT EDIT.

package buildcfg

import "runtime"

const defaultGO386 = `sse2`
const defaultGOAMD64 = `v1`

var defaultGOARM = func(goos string, goarch string) string {
	switch goos + `/` + goarch {

	case "linux/amd64":
		return `5`

	case "darwin/amd64":
		return `5`

	case "darwin/arm64":
		return `5`

	case "freebsd/386":
		return `5`

	case "freebsd/amd64":
		return `5`

	case "linux/386":
		return `5`

	case "linux/arm64":
		return `5`

	case "linux/armv6l":
		return `6`

	case "linux/ppc64le":
		return `5`

	case "linux/s390x":
		return `5`

	case "windows/386":
		return `7`

	case "windows/amd64":
		return `7`

	case "windows/arm64":
		return `7`

	case "js/wasm":
		return `5`

	}
	panic("gocompiler: unknown platform " + goos + `/` + goarch)
}(runtime.GOOS, runtime.GOARCH)

const defaultGOMIPS = `hardfloat`
const defaultGOMIPS64 = `hardfloat`
const defaultGOPPC64 = `power8`
const defaultGOEXPERIMENT = ``
const defaultGO_EXTLINK_ENABLED = ``

var defaultGO_LDSO = func(goos string, goarch string) string {
	switch goos + `/` + goarch {

	case "linux/amd64":
		return `/lib64/ld-linux-x86-64.so.2`

	case "darwin/amd64":
		return ``

	case "darwin/arm64":
		return ``

	case "freebsd/386":
		return `/libexec/ld-elf.so.1`

	case "freebsd/amd64":
		return `/libexec/ld-elf.so.1`

	case "linux/386":
		return `/lib64/ld-linux-x86-64.so.2`

	case "linux/arm64":
		return `/lib/ld-linux-aarch64.so.1`

	case "linux/armv6l":
		return `/lib/ld-linux-armhf.so.3`

	case "linux/ppc64le":
		return `/lib64/ld64.so.2`

	case "linux/s390x":
		return `/lib64/ld-linux-x86-64.so.2`

	case "windows/386":
		return ``

	case "windows/amd64":
		return ``

	case "windows/arm64":
		return ``

	case "js/wasm":
		return `/lib64/ld-linux-x86-64.so.2`

	}
	panic("gocompiler: unknown platform " + goos + `/` + goarch)
}(runtime.GOOS, runtime.GOARCH)

const version = `go1.19.3`
const defaultGOOS = runtime.GOOS
const defaultGOARCH = runtime.GOARCH

//gocompiler patch
