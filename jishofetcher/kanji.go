package jishofetcher

type KanjiCharacter struct {
	Kanji          string `json:"kanji"`
	KanjiJishoLink string `json:"kanji_jisho_link"`
	Strokes        int    `json:"strokes"`
	Kunyomi        string `json:"kunyomi"`
	Onyomi         string `json:"onyomi"`
	Meanings       string `json:"meaning"`
	Joyo           bool   `json:"joyo_kanji"`
	Grade          int    `json:"grade"`
	JLPT           string `json:"jlpt"`
}
