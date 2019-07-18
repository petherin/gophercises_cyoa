package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	cyoa "github.com/petherin/gophercises_cyoa"
)

func main() {
	port := flag.Int("port", 3000, "the port to start CYOA web application on")
	filename := flag.String("file", "gopher.json", "the JSON file with the CYOA story")
	parser := flag.String("parser", "default", "how to parse the chapters")
	flag.Parse()
	fmt.Printf("Using the story in %s.\n", *filename)

	f, err := os.Open(*filename)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	story, err := cyoa.JsonStory(f)
	if err != nil {
		panic(err)
	}

	var h http.Handler
	h = cyoa.NewHandler(story)

	if *parser == "form" {
		pathTpl := template.Must(template.ParseFiles("static/html/formLayout.html"))
		h = cyoa.NewHandler(story,
			cyoa.WithTemplate(pathTpl),
			cyoa.WithChapterParseFunc(formChapterParseFn))
	}

	fmt.Printf("Starting the server on port: %d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), h))
}

func formChapterParseFn(r *http.Request) (string, error) {
	if r.URL.Path != "/" {
		return "", fmt.Errorf("404 not found.")
	}

	var path string

	switch r.Method {
	case "GET":
		path = "intro"
	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			fmt.Printf("%v", err)
			return "", fmt.Errorf("ParseForm() err: %v", err)
		}

		path = r.FormValue("arc")
	default:
		return "", fmt.Errorf("Sorry, only GET and POST methods are supported.")
	}

	return path, nil
}
