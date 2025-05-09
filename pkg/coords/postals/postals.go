package postals

import (
	"fmt"
	"io"
	"os"

	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/coords"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Postals = *coords.CoordsRO[*Postal]

var postalCodesMap = map[string]*Postal{}

func New(cfg *config.Config) (Postals, error) {
	file, err := os.Open(cfg.PostalsFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open postals file. %w", err)
	}
	defer file.Close()

	buf, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read postals file. %w", err)
	}

	var codes []*Postal
	if err := json.Unmarshal(buf, &codes); err != nil {
		return nil, fmt.Errorf("failed to unmarshal postals file. %w", err)
	}

	cs, err := coords.NewReadOnly(codes)
	if err != nil {
		return nil, fmt.Errorf("failed to add postals to postals coords map: %w", err)
	}

	for k := range codes {
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
