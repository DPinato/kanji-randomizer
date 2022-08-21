package helper

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

// for HTTP client
const (
	UserAgent    = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36"
	Accept       = "*/*"
	RedditCookie = "over18=1; Path=/; Domain=.reddit.com; Secure;"
)

var DefaultHeaders = http.Header{"User-Agent": []string{UserAgent},
	"Accept": []string{Accept},
}

func DoHTTPRequest(url, method, jsonBody string, client *http.Client, h http.Header) (resp *http.Response, err error) {
	// log.Printf("DoHTTPRequest()\t")

	req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(jsonBody)))
	// log.Printf("%v\n", req)
	if err != nil {
		return nil, err
	}

	req.Header = h
	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, err
}

func ReadListFromFile(filepath string) ([]string, error) {
	var output []string
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		output = append(output, string(scanner.Bytes()))
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return output, nil
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) (string, error) {
	// source: https://golangcode.com/download-a-file-from-a-url/
	resp, err := http.Get(url) // Get the data
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath) // Create the file
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body) // Write the body to file
	return filepath, err
}

// GetHTMLFieldKeyValue given a string of HTML and a field name
// will return the value associated with the field
func GetHTMLFieldKeyValue(htmlStr, field string) string {
	pos1 := strings.Index(htmlStr, field) + len(field) + 2
	if pos1 == -1 {
		return ""
	}
	pos2 := strings.Index(htmlStr[pos1:], "\"") + pos1
	if pos2 == -1 {
		return ""
	}

	output := htmlStr[pos1:pos2]
	return output
}

// GetHTMLFieldValue given a string of HTML, returns the string between the next > and the following <
func GetHTMLFieldValue(htmlStr string) string {
	pos1 := strings.Index(htmlStr, ">")
	if pos1 == -1 {
		return ""
	}
	pos2 := strings.Index(htmlStr[pos1:], "<")
	if pos2 == -1 {
		return ""
	}
	return htmlStr[pos1+1 : pos1+pos2]
}

// GetWebPage returns a string containing the web page at the URL provided
// using the client and headers provided
func GetWebPage(url string, client *http.Client, h http.Header) (string, error) {
	resp, err := DoHTTPRequest(url, "GET", "", client, h)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("GetWebPage got response code %d\n", resp.StatusCode)
		return "", nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	return string(body), err
}
