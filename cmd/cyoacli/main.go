package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/petherin/gophercises_cyoa/internal"
)

func main() {
	filename := flag.String("file", "gopher.json", "the JSON file with the CYOA story")
	flag.Parse()
	fmt.Printf("Using the story in %s.\n", *filename)

	f, err := os.Open(*filename)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	story, err := internal.JsonStory(f)
	if err != nil {
		panic(err)
	}

	f.Close()

	storyTeller := internal.NewStoryTeller(story)
	storyTeller.Start()

}
