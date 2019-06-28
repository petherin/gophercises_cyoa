package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

//"title": "The Little Blue Gopher",
//"story": [
//"Once upon a time, long long ago, there was a little blue gopher. Our little blue friend wanted to go on an adventure, but he wasn't sure where to go. Will you go on an adventure with him?",
//"One of his friends once recommended going to New York to make friends at this mysterious thing called \"GothamGo\". It is supposed to be a big event with free swag and if there is one thing gophers love it is free trinkets. Unfortunately, the gopher once heard a campfire story about some bad fellas named the Sticky Bandits who also live in New York. In the stories these guys would rob toy stores and terrorize young boys, and it sounded pretty scary.",
//"On the other hand, he has always heard great things about Denver. Great ski slopes, a bad hockey team with cheap tickets, and he even heard they have a conference exclusively for gophers like himself. Maybe Denver would be a safer place to visit."
//],
//"options": [
//{
//"text": "That story about the Sticky Bandits isn't real, it is from Home Alone 2! Let's head to New York.",
//"arc": "new-york"
//},
//{
//"text": "Gee, those bandits sound pretty real to me. Let's play it safe and try our luck in Denver.",
//"arc": "denver"
//}
//]

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
	tmpl *template.Template
	arcs map[string]Arc
)

func handler(w http.ResponseWriter, r *http.Request) {
	key := "intro"
	keys, ok := r.URL.Query()["key"]

	if !ok || len(keys[0]) < 1 {
		log.Println("Url Param 'key' is missing, defaulting to intro")
		//return
	}

	if len(keys) > 0 {
		// Query()["key"] will return an array of items,
		// we only want the single item.
		key = keys[0]
	}
	tmpl.Execute(w, arcs[key])
}

func main() {
	// load json x
	// there will be a map like map[arcname]struct x
	// define what that struct looks like x
	// get json decoded into structs x
	// add structs to map x
	// display first arc x
	// show links to options
	// based on option, show next arc
	const introArc = "intro"
	jsonBytes := loadJson()

	json.Unmarshal([]byte(jsonBytes), &arcs)

	tmpl = template.Must(template.ParseFiles("html/layout.html"))

	http.HandleFunc("/", handler)

	http.ListenAndServe(":8080", nil)
}

func loadJson() []byte {
	// Open our jsonFile
	jsonFile, err := os.Open("gopher.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened gophers.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	return byteValue
}
