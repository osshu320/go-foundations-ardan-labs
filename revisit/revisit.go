package main

import (
	"bufio"
	"cmp"
	"compress/gzip"
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"slices"
	"sort"
	"strings"
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

func sha1Sum(fileName string) (string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return "", nil
	}
	defer file.Close()

	var r io.Reader = file

	if strings.HasSuffix(fileName, ".gz") {
		gz, err := gzip.NewReader(file)
		if err != nil {
			return "", err
		}
		defer gz.Close()
		r = gz
	}

	w := sha1.New()
	if _, err := io.Copy(w, r); err != nil {
		return "", err
	}

	sig := w.Sum(nil)
	return fmt.Sprintf("%x", sig), nil
}

func appendInt(s []int, v int) []int {
	i := len(s)

	if len(s) < cap(s) {
		s = s[:len(s)+1]
	} else {
		fmt.Printf("reallocate: %d->%d\n", len(s), 2*len(s)+1)
		s2 := make([]int, 2*len(s)+1)
		copy(s2, s)
		s = s2[:len(s)+1]
	}

	s[i] = v
	return s
}

func concat(s1 []string, s2 []string) []string {
	s := make([]string, len(s1)+len(s2))
	copy(s[:len(s1)], s1)
	copy(s[len(s1):], s2)
	return s

	// return append(s1, s2...)
}

func median(vs []float64) (float64, error) {
	if len(vs) == 0 {
		return 0, fmt.Errorf("median of empty slice")
	}

	nums := make([]float64, len(vs))
	copy(nums, vs)

	sort.Float64s(nums)
	n := len(nums)
	if n%2 == 1 {
		return nums[n/2], nil
	}

	return (nums[n/2] + nums[n/2-1]) / 2, nil
}

const (
	Jade Key = iota + 1
	Copper
	Crystal
	invalidKey // internal - not exported
)

const (
	maxX = 1000
	maxY = 600
)

type Key byte

type Item struct {
	X int
	Y int
}

func NewItem(x, y int) (*Item, error) {
	if x < 0 || x > maxX || y < 0 || y > maxY {
		return nil, fmt.Errorf("%d/%d out of bounds %d/%d", x, y, maxX, maxY)
	}

	i := Item{
		X: x,
		Y: y,
	}

	// Go compiler does escape analysis and allocate i on heap
	return &i, nil
}

// i is receiver
// if you want to mutate use pointer receiver
func (i *Item) Move(x, y int) {
	i.X = x
	i.Y = y
}

type Player struct {
	Name string
	Item
	Keys []Key
}

func getDist(p Player, x, y int) int {
	dx := p.X - x
	dy := p.Y - y
	if dx < 0 {
		dx = -dx
	}
	if dy < 0 {
		dy = -dy
	}
	return dx + dy
}

func sortByDistance(players []Player, x, y int) {
	slices.SortFunc(players, func(p, q Player) int {
		d1 := getDist(p, x, y)
		d2 := getDist(q, x, y)
		return cmp.Compare(d1, d2)
	})
}

func (p *Player) FoundKey(k Key) error {
	if k < Jade || k >= invalidKey {
		return fmt.Errorf("invalid key: %#v", k)
	}

	// if !containsKey(p.Key, k) {
	if !slices.Contains(p.Keys, k) {
		p.Keys = append(p.Keys, k)
	}

	return nil
}

func containsKey(keys []Key, k Key) bool {
	for _, kk := range keys {
		if kk == k {
			return true
		}
	}
	return false
}

type mover interface {
	Move(x, y int)
}

// Rule of Thumb: Accept interfaces, return types
func moveAll(ms []mover, x, y int) {
	for _, m := range ms {
		m.Move(x, y)
	}
}

func (k Key) String() string {
	switch k {
	case Jade:
		return "jade"
	case Copper:
		return "copper"
	case Crystal:
		return "crystal"
	}

	return fmt.Sprintf("<Key %d>", k)
}

func mostCommon(r io.Reader) (string, error) {
	freqs, err := wordFrequency(r)
	if err != nil {
		return "", err
	}

	return maxWord(freqs)
}

func maxWord(freqs map[string]int) (string, error) {
	if len(freqs) == 0 {
		return "", fmt.Errorf("empty map")
	}

	maxN, maxW := 0, ""
	for word, count := range freqs {
		if count > maxN {
			maxN = count
			maxW = word
		}
	}

	return maxW, nil
}

func wordFrequency(r io.Reader) (map[string]int, error) {
	s := bufio.NewScanner(r)
	freqs := make(map[string]int)
	for s.Scan() {
		words := wordRe.FindAllString(s.Text(), -1)
		for _, w := range words {
			freqs[strings.ToLower(w)]++
		}
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	return freqs, nil
}

var wordRe = regexp.MustCompile(`[a-zA-Z]+`)

// Runs Before main
// func init() {
// 	fmt.Println("Hello World (init)")
// }

func maxInts(nums []int) int {
	if len(nums) == 0 {
		return 0
	}

	max := nums[0]
	for _, n := range nums[1:] {
		if n > max {
			max = n
		}
	}
	return max
}

func maxFloat64s(nums []float64) float64 {
	if len(nums) == 0 {
		return 0
	}

	max := nums[0]
	for _, n := range nums[1:] {
		if n > max {
			max = n
		}
	}
	return max
}

type Number interface {
	int | float64
}

// func max[T int | float64](nums []T) T {
func max[T Number](nums []T) T {
	if len(nums) == 0 {
		return 0
	}

	max := nums[0]
	for _, n := range nums[1:] {
		if n > max {
			max = n
		}
	}
	return max
}

func safeDiv(a, b int) (q int, err error) {
	defer func() {
		if e := recover(); e != nil {
			log.Println("ERROR:", e)
			err = fmt.Errorf("%v", e)
		}
	}()

	return a / b, nil
}

func mostCommonN(r io.Reader, N int) error {
	freqs, err := wordFrequency(r)
	if err != nil {
		return err
	}

	type wf struct {
		word string
		freq int
	}

	var fs []wf
	for k, v := range freqs {
		fs = append(fs, wf{word: k, freq: v})
	}

	slices.SortFunc(fs, func(wf1, wf2 wf) int {
		if n := cmp.Compare(wf1.freq, wf2.freq); n != 0 {
			if n == -1 {
				return 1
			} else {
				return -1
			}
		}

		return cmp.Compare(wf1.word, wf2.word)
	})

	for i := 0; i < N; i++ {
		fmt.Println(fs[i])
	}

	return nil
}

func main() {
	slices_go_main()
}

func mostCommonN_demo() {
	file, err := os.Open("sherlock.txt")
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	defer file.Close()

	mostCommonN(file, 10)

	// w, err := mostCommon(file)
	// if err != nil {
	// 	log.Fatalf("error: %s", err)
	// }
	// fmt.Println(w)
}

func sortByDistance_demo() {
	p1 := Player{
		Name: "Parzival",
		Item: Item{500, 300},
	}
	p2 := Player{
		Name: "Parzival",
		Item: Item{400, 300},
	}
	p3 := Player{
		Name: "Parzival",
		Item: Item{100, 100},
	}
	players := []Player{p1, p2, p3}
	sortByDistance(players, 0, 0)
	fmt.Println(players)
}

func div_go() {
	fmt.Println(safeDiv(1, 0))
	fmt.Println(safeDiv(7, 2))
}

func empty_go() {
	var i any
	// var i interface{}

	i = 7
	fmt.Println(i)

	i = "hi"
	fmt.Println(i)

	// Rule of Thumb: Don't use any

	s := i.(string)
	fmt.Println(s)

	n, ok := i.(int)
	if ok {
		fmt.Println(n)
	} else {
		fmt.Println("not an int")
	}

	switch i.(type) {
	case int:
		fmt.Println("int")
	case string:
		fmt.Println("string")
	default:
		fmt.Printf("unknown type %T\n", i)
	}

	// fmt.Println(maxInts([]int{1, 2, 3}))
	// fmt.Println(maxFloat64s([]float64{3, 2, 1}))
	fmt.Println(max([]float64{3, 2, 1}))
	fmt.Println(max([]int{3, 2, 1}))
}

func freq_go() {
	var stocks map[string]float64
	sym := "TTWO"
	price := stocks[sym]
	fmt.Printf("%s -> $%.2f\n", sym, price)

	if price, ok := stocks[sym]; ok {
		fmt.Printf("%s -> $%.2f\n", sym, price)
	} else {
		fmt.Printf("%s not found\n", sym)
	}

	stocks = map[string]float64{
		sym:    137.73,
		"AAPL": 172.35,
	}
	if price, ok := stocks[sym]; ok {
		fmt.Printf("%s -> $%.2f\n", sym, price)
	} else {
		fmt.Printf("%s not found\n", sym)
	}

	for k := range stocks {
		fmt.Println(k)
	}

	for k, v := range stocks {
		fmt.Println(k, "->", v)
	}

	for _, v := range stocks {
		fmt.Println(v)
	}

	delete(stocks, "AAPL")
	fmt.Println(stocks)
	delete(stocks, "AAPL") // no panic

	file, err := os.Open("sherlock.txt")
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	defer file.Close()

	w, err := mostCommon(file)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	fmt.Println(w)
}

func game_go() {
	var i1 Item
	fmt.Println(i1)
	fmt.Printf("i1: %#v\n", i1)

	i2 := Item{1, 2}
	fmt.Printf("i2: %#v\n", i2)

	i3 := Item{
		Y: 10,
	}
	fmt.Printf("i3: %#v\n", i3)
	fmt.Println(NewItem(10, 20))
	fmt.Println(NewItem(10, -20))

	i3.Move(100, 200)
	fmt.Printf("i3 (move): %#v\n", i3)

	p1 := Player{
		Name: "Parzival",
		Item: Item{500, 300},
	}
	fmt.Printf("p1: %#v\n", p1)
	fmt.Printf("p1.X: %#v\n", p1.X)
	fmt.Printf("p1.Item.X: %#v\n", p1.Item.X)
	p1.Move(400, 500)
	fmt.Printf("p1 (move): %#v", p1)

	ms := []mover{
		&i1,
		&p1,
		&i2,
	}
	moveAll(ms, 0, 0)
	for _, m := range ms {
		fmt.Println(m)
	}

	k := Jade
	fmt.Println("k:", k)
	fmt.Println("key:", Key(17))

	p1.FoundKey(Jade)
	fmt.Println(p1.Keys)
	p1.FoundKey(Jade)
	fmt.Println(p1.Keys)
}

func slices_go() {
	var s []int
	fmt.Println("len", len(s))
	if s == nil {
		fmt.Println("nil slice")
	}

	s2 := []int{1, 2, 3, 4, 5, 6, 7}
	fmt.Printf("s2 = %#v\n", s2)

	s3 := s2[1:4]
	fmt.Printf("s2 = %#v\n", s3)

	s3 = append(s3, 100)
	fmt.Printf("s3 (append) = %#v\n", s3)
	fmt.Printf("s2 (append) = %#v\n", s2)
	fmt.Printf("s2: len=%d, cap=%d\n", len(s2), cap(s2))
	fmt.Printf("s3: len=%d, cap=%d\n", len(s3), cap(s3))

	var s4 []int
	for i := 0; i < 1000; i++ {
		s4 = appendInt(s4, i)
	}
	fmt.Println("s4", len(s4), cap(s4))

	// s4[1001] = 5 // panic

	fmt.Println(concat([]string{"A", "B"}, []string{"C", "D", "E"}))

	vs := []float64{2, 1, 3}
	fmt.Println(median(vs))
	vs = []float64{2, 1, 3, 4}
	fmt.Println(median(vs))
	fmt.Println(vs)

	fmt.Println(median(nil))
}

func sha1_go() {
	sig, err := sha1Sum("http.log.gz")
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	fmt.Println(sig)

	sig, err = sha1Sum("revisit.go")
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	fmt.Println(sig)
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
