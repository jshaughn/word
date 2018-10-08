// word gets definitions, synonyms and antonyms for a list of words
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Word struct {
	Word         string   `json:"word"`
	Score        int      `json:"score"`
	NumSyllables int      `json:"numSyllables"`
	Defs         []string `json:"defs"`
}

type Words []Word

func main() {
	filename := flag.String("f", "words.txt", "Filename with words, one word per line")
	synonyms := flag.Bool("s", false, "if specified, fetch synonyms")
	antonyms := flag.Bool("a", false, "if specified, fetch antonyms")
	helpShort := flag.Bool("h", false, "get usage help")
	helpLong := flag.Bool("help", false, "get usage help")

	flag.Parse()

	if *helpShort || *helpLong {
		flag.Usage()
		return
	}

	fmt.Printf("\nSearching for the words in [%s] with synonyms=[%v] and antonyms=[%v]\n", *filename, *synonyms, *antonyms)

	file, err := os.Open(*filename)
	check(err, "Invalid file")
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		w := scanner.Text()
		request := fmt.Sprintf("https://api.datamuse.com/words?sl=%v&md=d", w)
		makeRequest(fmt.Sprintf("Getting word: '%s'...", w), request, 1)

		if *synonyms {
			request = fmt.Sprintf("https://api.datamuse.com/words?rel_syn=%v&md=d", w)
			makeRequest(fmt.Sprintf("Synonyms for : '%s'...", w), request, 0)
		}

		if *antonyms {
			request = fmt.Sprintf("https://api.datamuse.com/words?rel_ant=%v&md=d", w)
			makeRequest(fmt.Sprintf("Antonyms for : '%s'...", w), request, 0)
		}
	}
	check(scanner.Err(), "Problem scanning file")
	fmt.Printf("\n\nDone!\n")
}

func makeRequest(title, request string, limit int) {
	response, err := http.Get(request)
	check(err, "The HTTP request failed")

	var words Words
	data, _ := ioutil.ReadAll(response.Body)
	err = json.Unmarshal(data, &words)
	check(err, "The Unmarshal failed")

	if len(words) == 0 {
		fmt.Printf("Sorry, could not find word")
		return
	}
	if "" != title {
		fmt.Printf("\n%s", title)
	}
	for i, word := range words {
		fmt.Printf("\n  %s", word.Word)
		if word.NumSyllables > 0 {
			fmt.Printf(" (%v syllables)", word.NumSyllables)
		}
		printDefs(word.Defs)
		if i == limit-1 {
			break
		}
	}
}

func printDefs(defs []string) {
	for _, def := range defs {
		fmt.Printf("\n\t%s", def)
	}
}

func check(e error, message string) {
	if e != nil {
		fmt.Printf("\n%s\n", message)
		panic(e)
	}
}
