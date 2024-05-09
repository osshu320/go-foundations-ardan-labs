package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
	"unicode/utf8"
)

func banner(text string, width int) {
	padding := (width - utf8.RuneCountInString(text)) / 2
	// padding := (width - len(text)) / 2
	for i := 0; i < padding; i++ {
		fmt.Print(" ")
	}
	fmt.Println(text)
	for i := 0; i < width; i++ {
		fmt.Print("-")
	}
	fmt.Println()
}

func isPalindrome(s string) bool {
	rs := []rune(s)
	for i := 0; i < len(rs)/2; i++ {
		if rs[i] != rs[len(rs)-1-i] {
			return false
		}
	}
	return true
}

func githubInfo(ctx context.Context, login string) (string, int, error) {
	url := "https://api.github.com/users/" + url.PathEscape(login)
	log.Println(url)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", 0, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", 0, err
	}

	if res.StatusCode != http.StatusOK {
		return "", 0, fmt.Errorf("%#v - %s", url, res.Status)
	}

	defer res.Body.Close()

	var r struct {
		Name     string
		NumRepos int `json:"public_repos"`
	}

	dec := json.NewDecoder(res.Body)
	if err := dec.Decode(&r); err != nil {
		return "", 0, err
	}

	return r.Name, r.NumRepos, nil
}

func main() {

}

func githubInfo_go() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	fmt.Println(githubInfo(ctx, "osshu320"))
}

func hw_go() {
	fmt.Println("Hello World! ☺")
}

func banner_go() {
	fmt.Println(len("☺"))
	banner("Go", 6)
	banner("G☺", 6)

	// s := "Go"
	s := "G☺"
	fmt.Println("len: ", len(s))

	for i, r := range s {
		fmt.Println(i, r)
		// this gives you rune
		fmt.Printf("%c of type %T\n", r, r)
	}

	b := s[0] // this gives you byte
	fmt.Printf("%c of type %T\n", b, b)

	x, y := 1, "1"
	fmt.Printf("x=%v, y=%v\n", x, y)
	fmt.Printf("x=%#v, y=%#v\n", x, y)

	fmt.Printf("%5s\n", s)

	fmt.Println("g", isPalindrome("g"))
	fmt.Println("go", isPalindrome("go"))
	fmt.Println("gog", isPalindrome("gog"))
	fmt.Println("g☺g", isPalindrome("g☺g"))
}
