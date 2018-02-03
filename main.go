package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/gentlemanautomaton/signaler"
)

func main() {
	// Capture shutdown signals
	shutdown := signaler.New().Capture(os.Interrupt, syscall.SIGTERM)

	// Parse arguments and environment
	c := DefaultConfig
	c.ParseEnv()
	if len(os.Args) > 0 {
		c.ParseArgs(os.Args[1:], flag.ExitOnError)
	}

	// Prepare an http server
	s := &http.Server{
		Addr:    "0.0.0.0:80",
		Handler: http.StripPrefix(c.URLPrefix, http.FileServer(filesOnlyFilesystem{http.Dir(c.Target)})),
	}

	// Create the mount
	if err := mount(c.Source, c.Target, "cifs", c.MountFlags(), c.MountOptions()); err != nil {
		fmt.Printf("Unable to mount \"%s\" at \"%s\": %v\n", c.Source, c.Target, err)
		os.Exit(1)
	}

	// Tell the server to stop gracefully when a shutdown signal is received
	stopped := shutdown.Then(func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()
		s.Shutdown(ctx)
	})

	// Unmount once the server has been shutdown
	unmounted := stopped.Then(func() {
		unmount(c.Target, DefaultUnmountFlags)
	})

	// Always cleanup and wait until the shutdown has completed
	defer unmounted.Wait()
	defer shutdown.Trigger()

	// Run the server and print the final result
	fmt.Println(s.ListenAndServe())
}

type filesOnlyFilesystem struct {
	fs http.FileSystem
}

func (fs filesOnlyFilesystem) Open(name string) (http.File, error) {
	f, err := fs.fs.Open(name)
	if err != nil {
		return nil, err
	}
	if stat, _ := f.Stat(); stat.IsDir() {
		return nil, os.ErrPermission
	}
	return f, nil
}
