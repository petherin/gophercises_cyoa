package internal

import (
	"html/template"
	"log"
	"net/http"
	"strings"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseFiles("static/html/defaultLayout.html"))
}

type HandlerOption func(h *handler)

func WithTemplate(t *template.Template) HandlerOption {
	return func(h *handler) {
		h.t = t
	}
}

func WithChapterParseFunc(fn func(r *http.Request) (string, error)) HandlerOption {
	return func(h *handler) {
		h.chapterParseFnc = fn
	}
}

func NewHandler(s Story, opts ...HandlerOption) http.Handler {
	h := handler{s, tpl, defaultChapterParseFn}
	for _, opt := range opts {
		opt(&h)
	}
	return h
}

type handler struct {
	s               Story
	t               *template.Template
	chapterParseFnc func(r *http.Request) (string, error)
}

func defaultChapterParseFn(r *http.Request) (string, error) {
	path := strings.TrimSpace(r.URL.Path)
	if path == "" || path == "/" {
		path = "/intro"
	}

	path = path[1:]

	return path, nil
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path, err := h.chapterParseFnc(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if chapter, ok := h.s[path]; ok {
		err := h.t.Execute(w, chapter)
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		}
		return
	}

	http.Error(w, "Chapter not found.", http.StatusNotFound)
}
