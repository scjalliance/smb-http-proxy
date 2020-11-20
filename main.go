package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/alecthomas/kong"
	"github.com/gentlemanautomaton/signaler"
)

func main() {
	// Capture shutdown signals
	shutdown := signaler.New().Capture(os.Interrupt, syscall.SIGTERM)

	// Parse arguments and environment
	var c Config
	kong.Parse(&c,
		kong.Description("Serves files from an SMB share over HTTP."),
		kong.UsageOnError())

	// Announce startup
	fmt.Printf("The process has started with this configuration:\n  %s\n", strings.Join(strings.Split(c.Summary(), "\n"), "\n  "))

	// Connect to the remote file system
	fs := NewFS(c)
	if err := fs.Connect(shutdown.Context()); err != nil {
		fmt.Printf("Failed to connect to \"%s\": %v\n", c.Source, err)
		os.Exit(1)
	}

	// Disconnect from the remote file system when done
	defer fs.Close()

	// Prepare an http server
	s := &http.Server{
		Addr:    "0.0.0.0:80",
		Handler: http.StripPrefix(c.URLPrefix, http.FileServer(filesOnlyFilesystem{fs})),
	}

	// Tell the server to stop gracefully when a shutdown signal is received
	stopped := shutdown.Then(func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()
		s.Shutdown(ctx)
	})

	// Always cleanup and wait until the shutdown has completed
	defer stopped.Wait()
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
