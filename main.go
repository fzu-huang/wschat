package main

import (
	"flag"
	"golang.org/x/net/websocket"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
	"wschat/handler"
)

var port *int = flag.Int("p", 23456, "Port to listen.")

func main() {
	flag.Parse()

	http.Handle("/login", websocket.Handler(handler.WSLogin))
	http.Handle("/say", websocket.Handler(handler.WSRandW))
	http.Handle("/logout", websocket.Handler(handler.WSLogout))

	http.ListenAndServe(`:`+strconv.Itoa(*port), nil)
}
