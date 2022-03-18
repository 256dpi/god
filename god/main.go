package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

var duration = flag.Int("duration", 2, "trace duration")

func main() {
	// parse flags
	flag.Parse()

	// get port
	port := flag.Arg(0)
	if port == "" {
		port = "6060"
	}

	// prepare files
	fmt.Println("==> creating temp files...")
	mem := tempFile("heap")
	allocs := tempFile("allocs")
	full := tempFile("full")
	cpu := tempFile("cpu")
	block := tempFile("block")
	mutex := tempFile("mutex")
	trace := tempFile("trace")

	// ensure cleanups
	defer cleanup(mem)
	defer cleanup(allocs)
	defer cleanup(full)
	defer cleanup(cpu)
	defer cleanup(block)
	defer cleanup(mutex)
	defer cleanup(trace)

	// download profiles
	fmt.Println("==> downloading profiles...")
	download(mem, fmt.Sprintf("http://localhost:%s/debug/pprof/heap", port))
	download(allocs, fmt.Sprintf("http://localhost:%s/debug/pprof/allocs", port))
	download(full, fmt.Sprintf("http://localhost:%s/debug/fgprof?seconds=%d", port, *duration))
	download(cpu, fmt.Sprintf("http://localhost:%s/debug/pprof/profile?seconds=%d", port, *duration))
	download(block, fmt.Sprintf("http://localhost:%s/debug/pprof/block?seconds=%d", port, *duration))
	download(mutex, fmt.Sprintf("http://localhost:%s/debug/pprof/mutex?seconds=%d", port, *duration))
	download(trace, fmt.Sprintf("http://localhost:%s/debug/pprof/trace?seconds=%d", port, *duration))

	// make sure trace command does not open a browser
	_ = os.Setenv("BROWSER", "/bin/echo")

	// run servers
	fmt.Println("==> running servers...")
	kill0 := run("go", "tool", "pprof", "-http=0.0.0.0:3797", "-no_browser", full.Name())
	kill1 := run("go", "tool", "pprof", "-http=0.0.0.0:3790", "-no_browser", cpu.Name())
	kill2 := run("go", "tool", "pprof", "-http=0.0.0.0:3791", "-no_browser", mem.Name())
	kill3 := run("go", "tool", "pprof", "-http=0.0.0.0:3792", "-no_browser", allocs.Name())
	kill4 := run("go", "tool", "pprof", "-http=0.0.0.0:3793", "-no_browser", block.Name())
	kill5 := run("go", "tool", "pprof", "-http=0.0.0.0:3794", "-no_browser", mutex.Name())
	kill6 := run("go", "tool", "trace", "-http=0.0.0.0:3796", trace.Name())

	// ensure kills
	defer kill0()
	defer kill1()
	defer kill2()
	defer kill3()
	defer kill4()
	defer kill5()
	defer kill6()

	// add handler
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(index(port)))
	}))

	// run server
	go func() {
		_ = http.ListenAndServe("0.0.0.0:3795", nil)
	}()

	// run browser
	err := exec.Command("open", "-a", "Google Chrome", "http://0.0.0.0:3795").Run()
	if err != nil {
		panic(err)
	}

	// prepare exit
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	<-exit

	// just exit
	os.Exit(0)
}

func tempFile(name string) *os.File {
	// create temp file
	f, err := ioutil.TempFile("", name)
	if err != nil {
		panic(err)
	}

	// print file
	fmt.Println(f.Name())

	return f
}

func cleanup(f *os.File) {
	err := os.Remove(f.Name())
	if err != nil {
		panic(err)
	}
}

func download(f *os.File, url string) {
	// print url
	fmt.Println(url)

	// download profile
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	// ensure close
	defer func() {
		_ = res.Body.Close()
	}()

	// copy data
	_, err = io.Copy(f, res.Body)
	if err != nil {
		panic(err)
	}
}

func run(bin string, args ...string) func() {
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

	fmt.Printf("%s %s\n", bin, strings.Join(args, " "))

	// run command
	err = cmd.Start()
	if err != nil {
		panic(err)
	}

	return func() {
		_ = cmd.Process.Kill()
	}
}
