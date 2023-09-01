package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func check(e error) {
	if(e != nil) {
		panic(e)
	}
}

func is_proper_noun(word string) bool {
	for pos, char := range word {
		if pos == 0 && unicode.IsUpper(char) {
			return true
		} else {
			break
		}
	}

	return false
}

var quote = false
func remove_punctuation(r rune) rune {
	if r == '\'' || r == '"' {
		quote = !quote
		return -1
	} else if strings.ContainsRune(`.,?!:;"'()-&/`, r) {
		if quote == true {
			return r
		} else {
			return -1
		}
	} else {
		return r
	}
}

func find_cloze(sentence string, frequencies map[string]int) string {
	quote = false
	
	sentence = strings.Map(remove_punctuation, sentence)
	
	words := strings.Fields(sentence)
	valid_words := make([]string, 0, 1)
	min_word, min_freq := "", math.MaxInt
	
	for _, word := range words {
		if len(word) <= 3 {
			continue
		}

		valid_words = append(valid_words, word)
		
		freq := frequencies[strings.ToLower(word)]
		if freq < min_freq {
			min_word = word
			min_freq = freq
		}
	}

	if len(min_word) > 0 {
		return min_word
	} else if len(valid_words) > 0 {
		return valid_words[int(math.Round(rand.Float64()*float64(len(valid_words)-1)))]
	} else {
		return ""
	}
}

func main() {
	stdout := bufio.NewWriter(os.Stdout)
	defer stdout.Flush()

	//stdout.WriteString("Initialising...\n")
	//stdout.Flush()
	
	sf, err := os.Open(os.Args[1])
	check(err)
	defer sf.Close()
	sentences := csv.NewReader(sf)
	sentences.Comma = '\t'
	
	ff, err := os.Open(os.Args[2])
	check(err)
	defer ff.Close()
	freqs := csv.NewReader(ff)
	freqs.Comma = ' '
	
	frequencies := make(map[string]int)
	for {
		row, eof := freqs.Read()
		
		if eof == io.EOF {
			break
		}

		if eof != nil {
			fmt.Fprintln(os.Stderr, eof)
			continue
		}
		
		v, err := strconv.Atoi(row[1])
		check(err)
		frequencies[row[0]] = v
	}

	//stdout.WriteString("Processing...\n")
	//stdout.Flush()

	for {
		row, err := sentences.Read()
		
		if err == io.EOF {
			break
		}
		
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		
		original := row[1]
		translation := row[3]

		cloze_word := find_cloze(original, frequencies)
		if len(cloze_word) == 0 {
			continue
		}

		original = strings.Replace(original, cloze_word, fmt.Sprintf("{{c1::%s}}", cloze_word), -1)

		stdout.WriteString(fmt.Sprintf("\"%s\",\"%s\"\n", original, translation))
		
		
		row, err = sentences.Read()
	}
	
}

