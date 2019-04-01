package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/conthing/conthing-admin-auth/auth"
)

func startHTTPServer(errChan chan error, port int) {
	go func() {
		http.HandleFunc("/api/v1/login", auth.LoginHandler)
		errChan <- http.ListenAndServe(":"+strconv.Itoa(port), nil)
	}()
}

func listenForInterrupt(errChan chan error) {
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errChan <- fmt.Errorf("%s", <-c)
	}()
}

func main() {
	start := time.Now()
	var port int

	flag.IntVar(&port, "port", 52000, "Specify a port other than default.")
	flag.IntVar(&port, "p", 52000, "Specify a port other than default.")
	flag.Parse()

	// 1.SIGINT 2.httpserver
	errs := make(chan error, 2)

	listenForInterrupt(errs)
	startHTTPServer(errs, port)

	// Time it took to start service
	log.Printf("Admin authorization listening on port %d, started in: %s", port, time.Since(start).String())

	// recv error channel
	c := <-errs
	log.Printf("terminating: %v", c)
	os.Exit(0)
}
