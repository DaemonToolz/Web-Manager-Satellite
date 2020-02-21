package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/gorilla/mux"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Printf("[ERROR] - %s: %s", msg, err)
	}
}

func printRequest(addr string) {
	log.Printf("[ %s ] - Request from %s ", time.Now().Format(time.RFC3339), addr)
}
