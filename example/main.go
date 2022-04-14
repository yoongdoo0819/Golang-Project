package main

import (
	"example/chat"
	"example/trace"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})

	t.templ.Execute(w, r)
}

func main() {

	var addr = flag.String("addr", ":3000", "The addr of the application")
	flag.Parse()

	gomniauth.SetSecurityKey("AIzaSyB8LRIELLggQdA9DWK0s6HpdDjePM4p74E")
	gomniauth.WithProviders(
		facebook.New("key", "secret", "http://localhost:3000/auth/callback/facebook"),
		github.New("key", "secret", "http://localhost:3000/auth/callback/github"),
		google.New("827713813492-570jt7bhjjod2h3b8ae6t3309akoaa0f.apps.googleusercontent.com", "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", "http://localhost:3000/auth/callback/google"),
	)

	r := chat.NewRoom()
	r.Tracer = trace.New(os.Stdout)
	http.Handle("/chat", chat.MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", chat.LoginHandler)
	http.Handle("/room", r)

	go r.Run()

	log.Println("Starting web server on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
