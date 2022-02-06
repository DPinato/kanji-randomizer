package jishofetcher

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/dpinato/kanji-randomizer/helper"
)

const JishoURL = "https://jisho.org"

func FetchKanjiList(jlptLevel, destFile string) error {
	log.Printf("Fetching Kanji list for %s to %s\n", jlptLevel, destFile)
	var kanjiList []KanjiCharacter

	jlptStr := "jlpt-" + jlptLevel
	urlStr := JishoURL + "/search/%23" + jlptStr + "%20%23kanji"

	// go through all the pages
	for {
		log.Printf("Fetching %s\n", urlStr)
		resp, err := http.Get(urlStr)
		if err != nil {
			return err
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		page := string(body)
		log.Printf("Success, read %d bytes\n", len(page))

		tmpList, err := ProcessPage(page)
		if err != nil {
			return err
		}
		kanjiList = append(kanjiList, tmpList...)

		// check if there is a next page
		nextPagePos := strings.Index(page, "<a class=\"more\"")
		if nextPagePos == -1 {
			log.Println("No more pages")
			break
		}
		urlStr = "https:" + helper.GetHTMLFieldKeyValue(page[nextPagePos:], "href")
	}

	log.Printf("Got %d kanji\n", len(kanjiList))

	// write JSON file
	jsonFileName := destFile + ".json"
	jsonFile, _ := json.MarshalIndent(kanjiList, "", " ")
	err := ioutil.WriteFile(jsonFileName, jsonFile, 0644)
	if err != nil {
		return fmt.Errorf("failed to write JSON file %s - %s", jsonFileName, err.Error())
	}

	// write CSV file
	csvFileName := destFile + ".csv"
	err = WriteCSVFile(kanjiList, csvFileName)
	if err != nil {
		return fmt.Errorf("failed to write CSV file %s - %s", csvFileName, err.Error())
	}

	return nil
}

func ProcessPage(pageStr string) ([]KanjiCharacter, error) {
	// go through a Jisho page and extract information about the Kanji characters
	// within the page
	var kanjiList []KanjiCharacter
	kanjiSectionStart := "kanji_light_content"
	tmpSection := pageStr

	for {
		pos := strings.Index(tmpSection, kanjiSectionStart)
		if pos == -1 {
			break
		}

		tmp, err := ProcessKanjiSection(tmpSection[pos:])
		if err != nil {
			log.Println(err)
		}
		kanjiList = append(kanjiList, tmp)
		tmpSection = tmpSection[pos+1:]
	}

	return kanjiList, nil
}

func ProcessKanjiSection(sectionStr string) (KanjiCharacter, error) {
	// process the section of the page and extract information about the kanji
	var tmpSection string
	var kanji KanjiCharacter

	jlptPos := strings.Index(sectionStr, "JLPT N")
	kanji.JLPT = sectionStr[jlptPos+5 : jlptPos+7]
	kanji.Joyo = strings.Contains(sectionStr, "Jōyō kanji")
	gradePos := strings.Index(sectionStr, "taught in grade ")
	tmpGrade := sectionStr[gradePos+16 : gradePos+17]
	kanji.Grade, _ = strconv.Atoi(tmpGrade)

	tmpSection = sectionStr[strings.Index(sectionStr, "literal_block"):]
	kanji.KanjiJishoLink = "https:" + helper.GetHTMLFieldKeyValue(tmpSection, "a href")
	kanji.Kanji = helper.GetHTMLFieldValue(tmpSection[strings.Index(tmpSection, "a href"):])

	tmpSection = sectionStr[strings.Index(sectionStr, "meanings english sense"):]
	tmpSection = tmpSection[:strings.Index(tmpSection, "</div>")]
	kanji.Meanings = GetKanjiEnglishMeanings(tmpSection)

	tmpSection = sectionStr[strings.Index(sectionStr, "kun readings"):]
	tmpSection = tmpSection[:strings.Index(tmpSection, "</div>")]
	kanji.Kunyomi = GetKanjiReadings(tmpSection)

	tmpSection = sectionStr[strings.Index(sectionStr, "on readings"):]
	tmpSection = tmpSection[:strings.Index(tmpSection, "</div>")]
	kanji.Onyomi = GetKanjiReadings(tmpSection)

	return kanji, nil
}

func GetKanjiEnglishMeanings(sectionStr string) string {
	// given the english meanings section, will return the meanings in a single string
	var meanings string
	tmpSection := sectionStr[strings.Index(sectionStr, "<span>"):]

	for {
		pos := strings.Index(tmpSection, "<span>")
		if pos == -1 {
			break
		}
		meanings += helper.GetHTMLFieldValue(tmpSection[pos:])
		tmpSection = tmpSection[pos+1:]
	}

	return meanings
}

func GetKanjiReadings(sectionStr string) string {
	// get the readings of the kanji from the section
	var readings string
	tmpSection := sectionStr[strings.Index(sectionStr, "<a href="):]

	for {
		pos := strings.Index(tmpSection, "<a href=")
		if pos == -1 {
			readings = readings[:len(readings)-1]
			break
		}
		readings += helper.GetHTMLFieldValue(tmpSection[pos:]) + " "
		tmpSection = tmpSection[pos+1:]
	}

	return readings
}

func WriteCSVFile(kanjiChars []KanjiCharacter, fileDir string) error {
	csvFile, err := os.Create(fileDir)
	if err != nil {
		log.Printf("Failed creating file\n")
		return err
	}
	defer csvFile.Close()

	csvwriter := csv.NewWriter(csvFile)
	csvwriter.Comma = '\t'
	for _, elem := range kanjiChars {
		var row []string
		row = append(row, elem.Kanji)
		row = append(row, elem.JLPT)
		row = append(row, elem.Kunyomi)
		row = append(row, elem.Onyomi)
		row = append(row, elem.Meanings)
		row = append(row, elem.KanjiJishoLink)
		row = append(row, strconv.FormatBool(elem.Joyo))
		row = append(row, strconv.Itoa(elem.Grade))

		err = csvwriter.Write(row)
		if err != nil {
			return err
		}

	}
	csvwriter.Flush()

	return nil
}
