package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

type Reply struct {
	Name string
	// Public_Repos int
	NumRepos int `json:"public_repos"`
}

// githubInfo returns name and number of public repos for login
func githubInfo(ctx context.Context, login string) (string, int, error) {
	link := "https://api.github.com/users/" + url.PathEscape(login)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, link, nil)
	// resp, err := http.Get(link)
	if err != nil {
		return "", 0, fmt.Errorf("error: %s", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", 0, err
	}

	if resp.StatusCode != http.StatusOK {
		return "", 0, fmt.Errorf("error: %s", resp.Status)
	}

	var r struct {
		Name     string
		NumRepos int `json:"public_repos"`
	}
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&r); err != nil {
		return "", 0, fmt.Errorf("error: can't decode - %s", err)
	}
	return r.Name, r.NumRepos, nil
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Millisecond)
	defer cancel()
	fmt.Println(githubInfo(ctx, "shubhamchemate003"))
}

func demo() {
	resp, err := http.Get("https://api.github.com/users/shubhamchemate003")
	if err != nil {
		log.Fatalf("error: %s", err)
		/*
			log.Printf("error: %s", err)
			os.Exit()
		*/
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("error: %s", resp.Status)
	}

	fmt.Printf("Content-Type: %s\n", resp.Header.Get("Content-Type"))
	/*
		if _, err := io.Copy(os.Stdout, resp.Body); err != nil {
			log.Fatalf("error: can't copy - %s", err)
		}
	*/
	var r Reply
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&r); err != nil {
		log.Fatalf("error: can't decode - %s", err)
	}
	fmt.Printf("%#v\n", r)
}

/*
JSON 		<-> Go
true/false 	<-> true/false
string 		<-> string
null 		<-> nil
number 		<-> float64, float32, int8, int16, int32, int64, int, uint8, ...
arrays 		<-> []any ([]interface{})
object		<-> map[string]any, struct

Every language has it's own data types, we need to make arrangements to convert
one to another

JSON -> io.Reader -> GO: json.Decoder
JSON -> []byte -> GO: json.Unmarshal
Go -> io.Writer -> JSON: json.Encoder
Go -> []byte -> JSON: json.Marshal
*/
