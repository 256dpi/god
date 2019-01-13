package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var mem = flag.Bool("mem", false, "memory profile")
var memAll = flag.Bool("mem-all", false, "full memory profile")
var trace = flag.Bool("trace", false, "trace profile")
var duration = flag.Int("duration", 5, "trace duration")

func main() {
	// parse flags
	flag.Parse()

	// get args
	arg := flag.Arg(0)

	// prepare port
	port := 6060
	if arg != "" {
		if n, err := strconv.Atoi(arg); err == nil {
			port += n
		}
	}

	// profile
	if *mem {
		fmt.Printf("mem: %d\n", port)
		profileMemory(port)
	} else if *memAll {
		fmt.Printf("mem-all: %d\n", port)
		profileMemoryAll(port)
	} else if *trace {
		fmt.Printf("trace: %d\n", port)
		profileTrace(port)
	} else {
		fmt.Printf("cpu: %d\n", port)
		profileCPU(port)
	}
}

func profileCPU(port int) {
	loc := fmt.Sprintf("http://localhost:%d/debug/pprof/profile", port)
	run("go", "tool", "pprof", "-pdf", "-output", "cpu.pdf", loc)
	run("open", "cpu.pdf")
}

func profileMemory(port int) {
	loc := fmt.Sprintf("http://localhost:%d/debug/pprof/heap", port)
	run("go", "tool", "pprof", "-pdf", "-output", "mem.pdf", loc)
	run("open", "mem.pdf")
}

func profileMemoryAll(port int) {
	loc := fmt.Sprintf("http://localhost:%d/debug/pprof/heap", port)
	run("go", "tool", "pprof", "--alloc_space", "-pdf", "-output", "mem.pdf", loc)
	run("open", "mem.pdf")
}

func profileTrace(port int) {
	loc := fmt.Sprintf("http://localhost:%d/debug/pprof/trace?seconds=%d", port, *duration)
	run("wget", "-O", "trace.out", loc)
	run("go", "tool", "trace", "trace.out")
}

func run(bin string, args ...string) {
	// create command
	cmd := exec.Command(bin, args...)

	// set working directory
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	cmd.Dir = wd

	// connect output
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// inherit current environment
	cmd.Env = os.Environ()

	fmt.Printf("=> %s %s\n", bin, strings.Join(args, " "))

	// run command
	err = cmd.Run()
	if err != nil {
		panic(err)
	}
}
