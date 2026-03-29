package postals

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/pkg/coords"
)

type postalStore struct {
	*coords.CoordsRO[*Postal]

	byCode map[string]*Postal
}

// Postals is a read-only postal coordinate store with additional postal code lookups.
type Postals = *postalStore

// ByCode returns the postal entry for the given code, if present.
func (p *postalStore) ByCode(code string) (*Postal, bool) {
	if p == nil {
		return nil, false
	}

	postal, ok := p.byCode[code]
	return postal, ok
}

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

	byCode := make(map[string]*Postal, len(codes))
	for i := range codes {
		if codes[i] == nil || codes[i].Code == nil {
			continue
		}

		byCode[*codes[i].Code] = codes[i]
	}

	return &postalStore{
		CoordsRO: cs,
		byCode:   byCode,
	}, nil
}
