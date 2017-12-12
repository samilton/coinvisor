package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path"
	"syscall"
)

func main() {
	dir := "."

	flag.StringVar(&dir, "dir", dir, "configuration directory")
	flag.Parse()

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	fmt.Printf("Starting iqvisor [%d]\n", os.Getpid())
	services := path.Join(dir, "services")
	conf := path.Join(dir, "conf")

	if d, e := os.Open(services); e != nil {
		fmt.Printf("Unable to find [%s]\n", services)
		os.Exit(-1)
	} else if files, e := d.Readdirnames(-1); e != nil {
		panic("no services defined")
	} else {
		for _, f := range files {
			serviceFile := path.Join(services, f)
			fmt.Printf("Loading %s\n", serviceFile)
		}
	}

	if d, e := os.Open(conf); e != nil {
		fmt.Printf("Unable to find [%s]\n", conf)
	} else {
		for _, f := range files {
			confFile := path.Join(conf, f)
			fmt.Printf("Loading %s\n", confFile)
		}
	}
	// Set up a handler, so that we shutdown cleanly if possible.
	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Printf("Shutting down from [%s]", sig)
		done <- true
	}()

	// Wait for a termination signal, and shutdown cleanly if we get it.
	<-done
	os.Exit(1)
}
