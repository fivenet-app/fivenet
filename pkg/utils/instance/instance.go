package instance

import (
	"crypto/rand"
	"encoding/hex"
	"os"
)

var (
	component  string
	instanceId string
)

func init() {
}

func SetComponent(t string) {
	component = t
	setID()
}

func setID() {
	if id := os.Getenv("POD_NAME"); id != "" {
		instanceId = component + "_" + id
		return
	}
	if id := os.Getenv("HOSTNAME"); id != "" {
		instanceId = component + "_" + id
		return
	}

	// Fallback to random ID generation
	b := make([]byte, 4)
	if _, err := rand.Read(b); err != nil {
		panic("failed to generate random instance ID. " + err.Error())
	}
	instanceId = component + "_" + hex.EncodeToString(b)
}

func ID() string {
	return instanceId
}
