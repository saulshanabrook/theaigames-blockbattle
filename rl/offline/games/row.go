package games

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type bot struct {
	name     string
	revision int
}

func newBot(sel *goquery.Selection) bot {
	revision, err := strconv.Atoi(
		strings.TrimPrefix(
			strings.TrimSpace(
				sel.Find(".div-revision-gameLog").Text(),
			),
			"v"))
	handleErr(err)
	return bot{
		name:     strings.TrimSpace(sel.Find(".div-botName-gameLog").Text()),
		revision: revision,
	}
}

type row struct {
	bots   [2]bot
	rounds int
	match  *url.URL
}

func newRow(sel *goquery.Selection) *row {
	r := row{
		bots: [2]bot{
			newBot(sel.Find(".cell-table-pointRight")),
			newBot(sel.Find(".cell-table-pointLeft")),
		},
	}
	var err error
	r.rounds, err = strconv.Atoi(
		strings.TrimSuffix(
			strings.TrimSpace(
				sel.Find(".cell-table-score").Text(),
			),
			" Rounds"))
	handleErr(err)
	matchS, exists := sel.Find(".button-reviewMatch").Attr("href")
	if !exists {
		panic("cant find href for match")
	}
	r.match, err = url.Parse(matchS)
	handleErr(err)
	return &r
}

func pageRows(rs chan<- *row, q string, p int) error {
	u := url.URL{
		Scheme:   "http",
		Host:     "theaigames.com",
		Path:     fmt.Sprintf("competitions/ai-block-battle/game-log/%v/%v", url.QueryEscape(q), p),
		RawQuery: `?_pjax=%23page`,
	}
	doc, err := goquery.NewDocument(u.String())
	handleErr(err)

	sels := doc.Find(".row-table")
	if sels.Length() == 0 {
		return errors.New("no rows")
	}
	for i := range sels.Nodes {
		rs <- newRow(sels.Eq(i))
	}
	return nil
}

func queryRowsWorker(rs chan<- *row, q string) {
	defer close(rs)
	var err error
	for p := 1; err == nil; p++ {
		err = pageRows(rs, q, p)
	}
	return
}

func queryRows(q string) <-chan *row {
	rs := make(chan *row)
	go queryRowsWorker(rs, q)
	return rs
}
