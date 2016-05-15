package games

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type data struct {
	Settings struct {
		Field struct {
			Width  int
			Height int
		}
		Players struct {
			Names []string
			Count int
		}
	}
	States []struct {
		NextShape string
		Round     int
		Winner    string
		Players   []struct {
			Move  string
			Skips string
			Field string
			Combo int
		}
	}
}

func newData(b []byte) *data {
	var d data
	handleErr(json.Unmarshal(b, d))
	return &d
}

func rowMatchBytes(r *row) []byte {
	r.match.Path += "/data"
	resp, err := http.Get(r.match.String())
	handleErr(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	handleErr(err)
	return body
}

func queryDataWorker(rs <-chan *row, ds chan<- *data) {
	defer close(ds)
	for r := range rs {
		ds <- newData(rowMatchBytes(r))
	}
	return
}

func queryData(rs <-chan *row) <-chan *data {
	ds := make(chan *data)
	go queryDataWorker(rs, ds)
	return ds
}
