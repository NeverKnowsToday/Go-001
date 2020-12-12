package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
)

var sErr chan error
var sigC chan os.Signal
var fd http.Server

func main() {
	g := errgroup.Group{}
	sErr = make(chan error, 1)
	sigC = make(chan os.Signal, 1)
	fd = http.Server{Addr:":8080"}

	g.Go(serve)
	g.Go(signalListen)
	select{}
}

func serve() error {
	sErr <- fd.ListenAndServe()
	select {
	case err := <- sErr:
		close(sigC)
		close(sErr)
		return err
	}
}

func signalListen() error{
	signal.Notify(sigC)
	select {
	case s:= <- sigC:
		fmt.Println("get signal:",s)
		signal.Stop(sigC)
		return fd.Shutdown(context.TODO())
	}
}