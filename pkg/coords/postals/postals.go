package postals

import (
	"fmt"
	"io"
	"os"

	"github.com/galexrt/fivenet/pkg/config"
	"github.com/galexrt/fivenet/pkg/coords"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Postals = *coords.Coords[*Postal]

var postalCodesMap = map[string]*Postal{}

func New(cfg *config.Config) (Postals, error) {
	file, err := os.Open(cfg.Game.Livemap.PostalsFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buf, _ := io.ReadAll(file)

	var codes []*Postal
	if err := json.Unmarshal(buf, &codes); err != nil {
		return nil, err
	}

	cs := coords.New[*Postal]()
	for k := range codes {
		if err := cs.Add(codes[k]); err != nil {
			return nil, fmt.Errorf("failed to add postal to postals map: %w", err)
		}

		if codes[k].Code != nil {
			postalCodesMap[*codes[k].Code] = codes[k]
		}
	}

	return cs, nil
}

func ByCode(code string) (*Postal, bool) {
	postalCode, ok := postalCodesMap[code]
	return postalCode, ok
}
