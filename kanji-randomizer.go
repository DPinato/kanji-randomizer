package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dpinato/kanji-randomizer/jishofetcher"
)

func handleRequest() (string, error) {
	return "Hello from Go!", nil
}

var EnvKey = "KANJI_RANDOMIZER_ENV"

func main() {
	var err error
	env := os.Getenv(EnvKey)
	switch env {
	case "lambda":
		lambda.Start(handleRequest)
	default:
		// run this locally

	}

	n5Count := flag.Int("n5", 0, "number of N5 kanji to output")
	n4Count := flag.Int("n4", 0, "number of N4 kanji to output")
	n3Count := flag.Int("n3", 0, "number of N3 kanji to output")
	n2Count := flag.Int("n2", 0, "number of N2 kanji to output")
	n1Count := flag.Int("n1", 0, "number of N1 kanji to output")
	flag.Parse()

	kanjiAmounts := map[string]int{
		"n5": *n5Count,
		"n4": *n4Count,
		"n3": *n3Count,
		"n2": *n2Count,
		"n1": *n1Count}

	log.Println(kanjiAmounts)

	// check whether we have the kanji list locally, otherwise we need to fetch it
	mainListDir := "/tmp/kanji-randomizer"
	for k, v := range kanjiAmounts {
		if v != 0 {
			localPath := mainListDir + "/jlpt_" + k // without extension
			fetchRemote := false

			// check if both JSON and CSV files exist
			if _, err := os.Stat(localPath + ".json"); errors.Is(err, os.ErrNotExist) {
				log.Printf("Could not find %s\n", localPath+".json")
				fetchRemote = true
			}
			if _, err := os.Stat(localPath + ".csv"); errors.Is(err, os.ErrNotExist) {
				log.Printf("Could not find %s\n", localPath+".csv")
				fetchRemote = true
			}

			if fetchRemote {
				err = jishofetcher.FetchKanjiList(k, localPath)
				if err != nil {
					log.Fatalf("Failed to fetch kanji list from Jisho - %v\n", err)
				}
				log.Printf("Fetched %s Kanji list, %s\n", k, localPath)
			} else {
				log.Printf("Found %s kanji list locally, %s\n", k, localPath)
			}
		}
	}

	// grab a random number of Kanji from the lists
	randKanjiList := []jishofetcher.KanjiCharacter{}
	for k, v := range kanjiAmounts {
		if v != 0 {
			localPath := mainListDir + "/jlpt_" + k // without extension

			list, err := getRandomNCharacters(localPath+".json", v)
			if err != nil {
				log.Fatalf("Failed to get %d random %s Kanji - %v\n", v, k, err)
			}

			randKanjiList = append(randKanjiList, list...)
		}
	}

	// give a nice useful output
	for _, elem := range randKanjiList {
		log.Printf("%s (%s) - %s\n", elem.Kanji, elem.JLPT, elem.KanjiJishoLink)
	}

}

func getRandomNCharacters(listFile string, n int) ([]jishofetcher.KanjiCharacter, error) {
	outputList := make([]jishofetcher.KanjiCharacter, n)
	listData := []jishofetcher.KanjiCharacter{}

	data, err := ioutil.ReadFile(listFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s - %s", listFile, err)
	}

	err = json.Unmarshal([]byte(data), &listData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON in %s - %s", listFile, err)
	}

	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	for i := 0; i < n; i++ {
		outputList[i] = listData[r.Intn(len(listData))]

		// just remove that randomly selected element from listData
		listData[i] = listData[len(listData)-1]
		listData = listData[:len(listData)-1]
	}

	return outputList, nil
}
