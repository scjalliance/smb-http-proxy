package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
	go func(s chan os.Signal) {
		<-s
		os.Exit(0)
	}(s)

	prefix := os.Getenv("URLPREFIX")
	http.ListenAndServe("0.0.0.0:80", http.StripPrefix(prefix, http.FileServer(filesOnlyFilesystem{http.Dir("/mnt/smb")})))
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
