package postals

import (
	"fmt"
	"io"
	"os"

	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/coords"
	jsoniter "github.com/json-iterator/go"
)

// json is a jsoniter instance configured to be compatible with the standard library.
var json = jsoniter.ConfigCompatibleWithStandardLibrary

// Postals is a type alias for a read-only coordinate store of Postal pointers.
type Postals = *coords.CoordsRO[*Postal]

// postalCodesMap maps postal codes to their corresponding Postal struct.
var postalCodesMap = map[string]*Postal{}

// New loads postal codes from the configured file and returns a read-only coordinate store.
// Returns an error if the file cannot be read or parsed, or if the points cannot be added.
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
		return nil, fmt.Errorf("failed to add postals to postals coords map. %w", err)
	}

	// Populate the postalCodesMap for fast lookup by code
	for k := range codes {
		if codes[k].Code != nil {
			postalCodesMap[*codes[k].Code] = codes[k]
		}
	}

	return cs, nil
}

// ByCode returns the Postal struct for a given code, if it exists.
func ByCode(code string) (*Postal, bool) {
	postalCode, ok := postalCodesMap[code]
	return postalCode, ok
}
