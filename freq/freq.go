package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
)

// Q: What is the most common word in sherlock.txt?
// word Frequency

func main() {
	file, err := os.Open("sherlock.txt")
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	defer file.Close()

	// w, err := mostCommon(file)
	// if err != nil {
	// 	log.Fatalf("error: %s", err)
	// }
	// fmt.Println(w)

	ws, err := mostCommonN(file, 10)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	fmt.Println(ws)

	// mapDemo()

	// path := "C:\to\new\report.csv"
	// `s` is a "raw" string, \ is just a \
	path := `C:\to\new\report.csv`
	fmt.Println(path)
}

func mostCommonN(r io.Reader, n int) ([]string, error) {
	freqs, err := wordFrequency(r)
	if err != nil {
		return []string{}, err
	}
	if len(freqs) < n {
		return []string{}, fmt.Errorf("only %d distinct words, asked %d", len(freqs), n)
	}
	return maxNWords(freqs, n)
}

type freqStruct struct {
	count int
	word  string
}

func maxNWords(freqs map[string]int, n int) ([]string, error) {
	freqSlice := make([]freqStruct, 0)
	for k, v := range freqs {
		freqSlice = append(freqSlice, freqStruct{count: v, word: k})
	}

	sort.Slice(freqSlice, func(i, j int) bool {
		return freqSlice[i].count > freqSlice[j].count
	})

	// sort.Sort(freqSlice)

	w := make([]string, 0)
	for _, fs := range freqSlice[:n+1] {
		w = append(w, fs.word)
	}
	return w, nil
}

func mostCommon(r io.Reader) (string, error) {
	freqs, err := wordFrequency(r)
	if err != nil {
		return "", err
	}
	return maxWords(freqs)
}

func mapDemo() {
	var stocks map[string]float64 // symbol -> price
	sym := "TTWO"
	price := stocks[sym]
	log.Printf("%s -> $%.2f\n", sym, price)

	if price, ok := stocks[sym]; ok {
		log.Printf("%s -> $%.2f\n", sym, price)
	} else {
		log.Printf("%s not found\n", sym)
	}

	/*
		stocks = make(map[string]float64)
		stocks[sym]=136.73
	*/
	stocks = map[string]float64{
		sym:    137.73,
		"AAPL": 172.35,
	}

	stocks[sym] = 136.73
	if price, ok := stocks[sym]; ok {
		log.Printf("%s -> $%.2f\n", sym, price)
	} else {
		log.Printf("%s not found\n", sym)
	}

	for k := range stocks { // keys
		log.Println(k)
	}

	for k, v := range stocks { // keys and values
		log.Println(k, "->", v)
	}

	for _, v := range stocks {
		log.Println(v)
	}

	delete(stocks, "AAPL")
	log.Println(stocks)
	delete(stocks, "AAPL") // happens nothing
}

/*
// You can also use raw string to create multi line strings
var request = `GET / HTTP/1.1
Host: httpbin.org
Connection: Close

`
*/

// "Who's on first?" -> [Whos on first]
var wordRe = regexp.MustCompile(`[a-zA-Z]+`)

/* Will run before main
func init() {
	// ...
}
*/

func maxWords(freqs map[string]int) (string, error) {
	if len(freqs) == 0 {
		return "", fmt.Errorf("empty map")
	}

	maxN, maxW := 0, ""
	for word, count := range freqs {
		if count > maxN {
			maxN, maxW = count, word
		}
	}

	return maxW, nil
}

func wordFrequency(r io.Reader) (map[string]int, error) {
	s := bufio.NewScanner(r)
	freqs := make(map[string]int) // word -> count
	lnum := 0
	for s.Scan() {
		lnum++
		words := wordRe.FindAllString(s.Text(), -1) // current line
		for _, w := range words {
			freqs[strings.ToLower(w)]++
		}
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	log.Println("num lines:", lnum)

	return freqs, nil
}
