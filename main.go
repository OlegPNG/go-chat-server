package main

import (
	"embed"
	"flag"
	"log"
	"net/http"

	//"path"
	"html/template"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	//"github.com/a-h/templ"
)

var addr = flag.String("addr", ":8080", "http service address")

//go:embed public/*
var files embed.FS

var (
    msgCard = parse("public/message_card.html")
)

func main() {
    flag.Parse()
    hub := newHub()
    go hub.run()


    r := chi.NewRouter()
    r.Use(middleware.Logger)
    r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
        serveWs(hub, w, r)
    })
    r.Get("/chat", func(w http.ResponseWriter, r *http.Request) {
	chat(hub.history).Render(r.Context(), w)
    })
    err := http.ListenAndServe(*addr, r)
    if err != nil {
        log.Fatal("ListenAndServer: ", err)
    }
}

func parse(file ...string) *template.Template {
    return template.Must(
	template.ParseFS(files, file...),
    )
}
