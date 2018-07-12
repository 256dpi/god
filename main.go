package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

var mem = flag.Bool("mem", false, "memory profile")

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

	// run command
	err = cmd.Run()
	if err != nil {
		panic(err)
	}
}
