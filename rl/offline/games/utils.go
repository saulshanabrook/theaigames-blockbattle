package games

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
	return
}
