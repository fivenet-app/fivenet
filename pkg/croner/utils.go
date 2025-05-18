package croner

import (
	"fmt"
	"net"
	"os"
)

func getNodeName(adminListen string) (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", err
	}
	_, port, err := net.SplitHostPort(adminListen)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s-%s", hostname, port), nil
}
