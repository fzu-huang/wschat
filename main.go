package main

import (
	"flag"
	//"github.com/Unknwon/macaron"
	//"github.com/go-macaron/pongo2"
	"golang.org/x/net/websocket"
	"net/http"
	"strconv"
	"wschat/handler"
)

var port *int = flag.Int("p", 23456, "Port to listen.")

func main() {

	flag.Parse()
	flag.Lookup("logtostderr").Value.Set("true")
	//	m := macaron.Classic()
	//	m.Use(macaron.Logger())
	//	m.Use(macaron.Recovery())
	//	m.Use(pongo2.Pongoer(pongo2.Options{
	//		Directory:  "templates",
	//		Extensions: []string{".html"},
	//	}))

	//	m.Get("/", handler.Index)
	//	m.Post("/login", handler.WCLogin)
	//	m.Get("/logout", handler.WCLogout)

	//m.Use(websocket.Handler(handler.WSEnterRoom))
	//m.Handle("POST", "/say", )

	http.Handle("/join", websocket.Handler(handler.WSEnterRoom))
	http.Handle("/createroom", websocket.Handler(handler.WSCreateRoom))
	http.ListenAndServe(":"+strconv.Itoa(*port), nil)
	//go m.Run(*port + 1)

}
