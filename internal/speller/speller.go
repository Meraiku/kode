package speller

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func CheckText(text string) (string, error) {

	client := http.Client{Timeout: time.Minute}

	resp, err := client.PostForm(spellerURL, url.Values{
		"text":    {text},
		"lang":    {spellerLang},
		"options": {spellerOpt},
	})
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	bytes.TrimSpace(body)

	outSpell := []SpellerResponse{}

	json.Unmarshal(body, &outSpell)

	text = correctText(text, outSpell)

	return text, nil
}

func correctText(text string, words []SpellerResponse) string {
	for _, v := range words {
		text = strings.ReplaceAll(text, v.Word, v.S[0])
	}

	return text
}
