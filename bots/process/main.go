package main

import (
	"crypto/rand"
	"fmt"
	"os"
	"reflect"

	"github.com/bouk/monkey"
	"github.com/saulshanabrook/blockbattle/bots"
)

//https://github.com/golang/go/issues/6871#issuecomment-66088766
var errorType = reflect.TypeOf(make([]error, 1)).Elem()

// bad bad bad
// since we are in a chroot-ed env, we can't acccess `/dev/urandom`
// so calls to random will fail when using websockets
func patchRand() {
	fType := reflect.FuncOf(
		[]reflect.Type{reflect.TypeOf(rand.Reader), reflect.TypeOf([]byte{})},
		[]reflect.Type{reflect.TypeOf(1), errorType},
		false,
	)
	fInner := func(args []reflect.Value) (results []reflect.Value) {
		fmt.Fprintf(os.Stderr, "Calling mock reader\n")
		return []reflect.Value{reflect.ValueOf(1), reflect.Zero(errorType)}
	}

	fVal := reflect.MakeFunc(fType, fInner)
	monkey.PatchInstanceMethod(reflect.TypeOf(rand.Reader), "Read", fVal.Interface())
}

var nnHost string

func main() {
	patchRand()

	b := bots.NewDQN(nnHost)
	p := NewPlayer()
	bots.Play(b, p)
}
