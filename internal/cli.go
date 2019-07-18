package internal

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

type StoryTeller struct {
	story Story
}

func NewStoryTeller(s Story) StoryTeller {
	return StoryTeller{s}
}

func (s StoryTeller) Start() {
	chapter := "intro"

	for chapter != "" {
		s.showChapter(chapter)
		chapter = s.getNextChapter(chapter)
	}

	fmt.Print("The End\n")
}

func (s StoryTeller) showChapter(chapterName string) {
	chapter := s.story[chapterName]
	fmt.Println("------------------------------------------------")
	fmt.Printf("\n%s\n\n", chapter.Title)

	for _, para := range chapter.Paragraphs {
		fmt.Printf("%s\n", para)
	}
	fmt.Println("")

	var i int
	for _, opt := range chapter.Options {
		i++
		fmt.Printf("%d) %s\n", i, opt.Text)
	}
}

func (s StoryTeller) getNextChapter(chapterName string) string {
	chapter := s.story[chapterName]
	if len(chapter.Options) == 0 {
		return ""
	}

	fmt.Print("\nPress a number followed by Enter...\n")

	nextChapter := ""
	pattern := regexp.MustCompile(fmt.Sprintf("[1-%d]", len(chapter.Options)))
	scanner := bufio.NewScanner(os.Stdin)

	for nextChapter == "" {
		scanner.Scan()

		if scanner.Err() != nil {
			fmt.Printf("%v\n", scanner.Err())
			continue
		}

		textInput := scanner.Text()
		if !pattern.MatchString(textInput) {
			fmt.Printf("Invalid option\n")
			continue
		}

		input, err := strconv.Atoi(textInput)
		if err != nil {
			fmt.Printf("%v\n", err)
			continue
		}

		nextChapter = chapter.Options[input-1].Chapter
	}

	return nextChapter
}
