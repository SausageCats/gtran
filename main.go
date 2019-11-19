package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"unicode/utf8"
)

// required to translate language flags
var (
	// source language
	source = flag.String("source", "", "translate source")
	// target language
	target = flag.String("target", "ja", "translate traget")
	// source language text
	text = flag.String("text", "", "translate source text")
	// Use the Google Apps Script to translate language
	endpoint = flag.String("endpoint", "https://script.google.com/macros/s/AKfycbzU7EjH7TAakYcypslv9qvzyRF5yGDdHBG_r3ZXMDSvzqzdYtrn/exec", "translate endpoint")
)

type post struct {
	Text   string `json:"text"`
	Source string `json:"source"`
	Target string `json:"target"`
}

// translate language
func translate(text, source, target string) (string, error) {
	postData, err := json.Marshal(post{text, source, target})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, *endpoint, bytes.NewBuffer([]byte(postData)))

	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func run(args []string) int {
	envEndpoint := os.Getenv("GTRAN_ENDPOINT")
	if envEndpoint != "" {
		*endpoint = envEndpoint
	}

	if *text == "" {
		*text = strings.Join(args, " ")
		if *text == "" {
			flag.Usage()
			return -1
		}
	}

	if *source == "" {
		nEN := 0
		nJA := 0
		for _, char := range *text {
			str := string([]rune{char})
			if str == " " {
				continue
			}
			nByte := len(str)
			nRune := utf8.RuneCountInString(str)
			if nByte == nRune {
				nEN++
			} else {
				nJA++
			}
		}
		if nEN/2 <= nJA {
			*source = "ja"
			*target = "en"
		} else {
			*source = "en"
			*target = "ja"
		}
	}

	result, err := translate(*text, *source, *target)
	if err != nil {
		return -1
	}
	fmt.Println(result)
	return 0
}

func main() {
	flag.Parse()
	os.Exit(run(flag.Args()))
}
