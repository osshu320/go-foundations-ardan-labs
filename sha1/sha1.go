package main

import (
	"compress/gzip"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"strings"
)

func sha1Sum(fileName string) (string, error) {
	// idiom: acquire a resource, check for error, defer release
	file, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close() // defer are called in LIFO order

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

func main() {
	sig, err := sha1Sum("http.log.gz")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(sig)

	sig, err = sha1Sum("sha1.go")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(sig)
}
