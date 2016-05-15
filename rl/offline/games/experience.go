package games

import "github.com/saulshanabrook/blockbattle/dqn"

func parseSingleData(d *data, es chan<- *dqn.Experience) {
	return
}

func parseDataWorker(ds <-chan *data, es chan<- *dqn.Experience) {
	defer close(es)
	for d := range ds {
		parseSingleData(d, es)
	}
	return
}

func parseData(ds <-chan *data) <-chan *dqn.Experience {
	es := make(chan *dqn.Experience)
	go parseDataWorker(ds, es)
	return es
}
