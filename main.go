package main

import (
	"math/rand"
	"time"

	"github.com/galexrt/fivenet/cmd"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	cmd.Execute()
}
