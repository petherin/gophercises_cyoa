package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
)

type Arc struct {
	Title   string   `json:title`
	Story   []string `json:story`
	Options []Option `json:options`
}

type Option struct {
	Text string `json:text`
	Arc  string `json:arc`
}

var (
	arcs map[string]Arc
)

func handler(w http.ResponseWriter, r *http.Request) {
	var arc string

	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		arc = "intro"
	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			fmt.Printf("%v",err)
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}

		arc = r.FormValue("arc")
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}

	tmpl := template.Must(template.ParseFiles("static/html/layout.html"))

	tmpl.Execute(w, arcs[arc])
}

func main() {
	jsonBytes := loadJson()

	json.Unmarshal([]byte(jsonBytes), &arcs)

	http.HandleFunc("/", handler)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./static/css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./static/js"))))

	//r.PathPrefix("/GolangTraining/blog/static/css").Handler(http.StripPrefix("/GolangTraining/blog/static/css", http.FileServer(http.Dir("./GolangTraining/blog/static/css"))))


	http.ListenAndServe(":8080", nil)
}

func loadJson() []byte {
	jsonFile, err := os.Open("./gopher.json")
	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	return byteValue
}
