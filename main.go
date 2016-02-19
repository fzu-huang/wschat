package main

import (
	"flag"
	"github.com/Unknwon/macaron"
	//"github.com/go-macaron/pongo2"
	"golang.org/x/net/websocket"
	"net/http"
	"strconv"
	"wschat/handler"
	. "wschat/util"
)

var port *int = flag.Int("p", 23456, "Port to listen.")

func main() {

	flag.Parse()
	flag.Lookup("logtostderr").Value.Set("true")

	WSSERVER += ":" + strconv.Itoa(*port) + `/join`

	m := macaron.Classic()
	m.Use(macaron.Logger())
	m.Use(macaron.Recovery())
	m.Use(macaron.Static("public"))
	m.Use(macaron.Renderer(macaron.RenderOptions{
		// 模板文件目录，默认为 "templates"
		Directory: "templates",
		// 模板文件后缀，默认为 [".tmpl", ".html"]
		Extensions: []string{".tmpl"},
	}))

	m.Get("/", handler.Index)
	m.Post("/user/login", handler.WCLogin)
	m.Get("/logout", handler.WCLogout)

	m.Get("/roomusers", handler.ListusersinRoom)
	m.Get("/rooms", handler.ListRoom)
	m.Get("/createroom", handler.CreateRoom)

	m.Get("/tab", func(ctx *macaron.Context) { ctx.HTML(200, "tab") })

	//m.Use(websocket.Handler(handler.WSEnterRoom))
	//m.Handle("POST", "/say", )

	//http.HandleFunc("/index", handler.Index)

	http.Handle("/join", websocket.Handler(handler.WSEnterRoom))
	http.Handle("/createroom", websocket.Handler(handler.WSCreateRoom))
	go http.ListenAndServe(":"+strconv.Itoa(*port), nil)
	m.Run(*port + 1)

}
