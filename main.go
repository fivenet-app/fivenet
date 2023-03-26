package main

import (
	"math/rand"
	"time"

	"github.com/galexrt/arpanet/cmd"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	cmd.Execute()
}
